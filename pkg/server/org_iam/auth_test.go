package org_iam

import (
	"context"
	"reflect"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/iam"
	iammock "github.com/ca-risken/core/proto/iam/mocks"
	"github.com/ca-risken/core/proto/organization"
	organizationmock "github.com/ca-risken/core/proto/organization/mocks"
	"github.com/ca-risken/core/proto/org_iam"
	"github.com/ca-risken/core/proto/project"
	projectmock "github.com/ca-risken/core/proto/project/mocks"
	"gorm.io/gorm"
)

func TestIsAuthorizedByOrgPolicy(t *testing.T) {
	validPolicies := &[]model.OrganizationPolicy{
		{PolicyID: 1, Name: "organization-admin", OrganizationID: 1, ActionPtn: "organization/.*", ProjectPtn: ".*"},
		{PolicyID: 2, Name: "organizatino-viewer", OrganizationID: 1, ActionPtn: "project/(get|list)", ProjectPtn: ".*"},
	}
	cases := []struct {
		name        string
		action      string
		projectName string
		policy      *[]model.OrganizationPolicy
		want        bool
		wantErr     bool
	}{
		{
			name:   "OK Authorized organization get without project context",
			action: "organization/get-organization",
			policy: validPolicies,
			want:   true,
		},
		{
			name:   "OK Unauthorized action not allowed",
			action: "organization/delete-organization",
			policy: &[]model.OrganizationPolicy{{PolicyID: 2, Name: "organization-viewer", OrganizationID: 1, ActionPtn: "organization/(get|list)", ProjectPtn: ".*"}},
			want:   false,
		},
		{
			name:        "OK Authorized with project name match",
			action:      "finding/put-finding",
			projectName: "team-a-prod",
			policy: &[]model.OrganizationPolicy{
				{PolicyID: 3, Name: "team-a-editor", OrganizationID: 1, ActionPtn: "finding/.*", ProjectPtn: "team-a-.*"},
			},
			want: true,
		},
		{
			name:        "OK Unauthorized when project name does not match project_ptn",
			action:      "finding/put-finding",
			projectName: "team-b-prod",
			policy: &[]model.OrganizationPolicy{
				{PolicyID: 3, Name: "team-a-editor", OrganizationID: 1, ActionPtn: "finding/.*", ProjectPtn: "team-a-.*"},
			},
			want: false,
		},
		{
			name:    "NG Error invalid action regex pattern",
			action:  "organization/get",
			policy:  &[]model.OrganizationPolicy{{PolicyID: 1, Name: "invalid-pattern", OrganizationID: 1, ActionPtn: "[invalid regex", ProjectPtn: ".*"}},
			wantErr: true,
		},
		{
			name:        "NG Error invalid project regex pattern",
			action:      "finding/put-finding",
			projectName: "team-a-prod",
			policy:      &[]model.OrganizationPolicy{{PolicyID: 1, Name: "invalid-project-pattern", OrganizationID: 1, ActionPtn: "finding/.*", ProjectPtn: "[invalid regex"}},
			wantErr:     true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := isAuthorizedByOrgPolicy(c.action, c.projectName, c.policy)
			if (err != nil) != c.wantErr {
				t.Errorf("isAuthorizedByOrgPolicy() error = %v, wantErr %v", err, c.wantErr)
				return
			}
			if got != c.want {
				t.Errorf("isAuthorizedByOrgPolicy() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestIsAuthorizedOrganization(t *testing.T) {
	cases := []struct {
		name              string
		input             *org_iam.IsAuthorizedOrgRequest
		want              *org_iam.IsAuthorizedOrgResponse
		wantErr           bool
		mockResponse      *[]model.OrganizationPolicy
		mockError         error
		mockIsAdminResp   bool
		mockIsAdminErr    error
		expectIsAdminCall bool
		callListProject   bool
		listProjectResp   *project.ListProjectResponse
		listProjectErr    error
	}{
		{
			name: "OK Admin user - immediate authorization",
			input: &org_iam.IsAuthorizedOrgRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "organization/update-organization",
			},
			want:              &org_iam.IsAuthorizedOrgResponse{Ok: true},
			mockIsAdminResp:   true,
			expectIsAdminCall: true,
		},
		{
			name: "OK Authorized - non-admin user with matching policy",
			input: &org_iam.IsAuthorizedOrgRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "organization/update-organization",
			},
			want:              &org_iam.IsAuthorizedOrgResponse{Ok: true},
			mockIsAdminResp:   false,
			expectIsAdminCall: true,
			mockResponse: &[]model.OrganizationPolicy{
				{PolicyID: 101, Name: "organization-admin", OrganizationID: 1001, ActionPtn: "organization/.*"},
			},
		},
		{
			name: "OK Unauthorized - non-admin user with no matching policy",
			input: &org_iam.IsAuthorizedOrgRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "organization/delete-organization",
			},
			want:              &org_iam.IsAuthorizedOrgResponse{Ok: false},
			mockIsAdminResp:   false,
			expectIsAdminCall: true,
			mockResponse: &[]model.OrganizationPolicy{
				{PolicyID: 102, Name: "organization-viewer", OrganizationID: 1001, ActionPtn: "organization/(get|list)"},
			},
		},
		{
			name: "OK Unauthorized - non-admin user with no policies found",
			input: &org_iam.IsAuthorizedOrgRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "organization/create-organization",
			},
			want:              &org_iam.IsAuthorizedOrgResponse{Ok: false},
			mockIsAdminResp:   false,
			expectIsAdminCall: true,
			mockError:         gorm.ErrRecordNotFound,
		},
		{
			name: "NG IsAdmin check error",
			input: &org_iam.IsAuthorizedOrgRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "organization/update-organization",
			},
			mockIsAdminErr:    gorm.ErrInvalidDB,
			expectIsAdminCall: true,
			wantErr:           true,
		},
		{
			name: "NG Invalid params - organization_id is zero",
			input: &org_iam.IsAuthorizedOrgRequest{
				UserId:         111,
				OrganizationId: 0,
				ActionName:     "organization/create-organization",
			},
			wantErr: true,
		},
		{
			name: "NG Invalid params - action_name is invalid",
			input: &org_iam.IsAuthorizedOrgRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "",
			},
			wantErr: true,
		},
		{
			name: "NG Invalid DB error in policy check",
			input: &org_iam.IsAuthorizedOrgRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "organization/create-organization",
			},
			mockIsAdminResp:   false,
			expectIsAdminCall: true,
			mockError:         gorm.ErrInvalidDB,
			wantErr:           true,
		},
		{
			name: "OK Authorized with project_ptn matching",
			input: &org_iam.IsAuthorizedOrgRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "finding/put-finding",
				ProjectId:      3001,
			},
			want:              &org_iam.IsAuthorizedOrgResponse{Ok: true},
			mockIsAdminResp:   false,
			expectIsAdminCall: true,
			callListProject:   true,
			listProjectResp: &project.ListProjectResponse{
				Project: []*project.Project{{ProjectId: 3001, Name: "team-a-prod"}},
			},
			mockResponse: &[]model.OrganizationPolicy{
				{PolicyID: 101, Name: "team-a-editor", OrganizationID: 1001, ActionPtn: "finding/.*", ProjectPtn: "team-a-.*"},
			},
		},
		{
			name: "OK Unauthorized when project_ptn does not match",
			input: &org_iam.IsAuthorizedOrgRequest{
				UserId:         111,
				OrganizationId: 1001,
				ActionName:     "finding/put-finding",
				ProjectId:      3002,
			},
			want:              &org_iam.IsAuthorizedOrgResponse{Ok: false},
			mockIsAdminResp:   false,
			expectIsAdminCall: true,
			callListProject:   true,
			listProjectResp: &project.ListProjectResponse{
				Project: []*project.Project{{ProjectId: 3002, Name: "team-b-prod"}},
			},
			mockResponse: &[]model.OrganizationPolicy{
				{PolicyID: 101, Name: "team-a-editor", OrganizationID: 1001, ActionPtn: "finding/.*", ProjectPtn: "team-a-.*"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := mocks.NewOrgIAMRepository(t)
			logger := logging.NewLogger()
			mockIAM := iammock.NewIAMServiceClient(t)
			mockOrg := organizationmock.NewOrganizationServiceClient(t)
			mockProject := projectmock.NewProjectServiceClient(t)
			svc := NewOrgIAMService(mockRepo, mockOrg, mockIAM, mockProject, logger)

			if c.expectIsAdminCall {
				if c.mockIsAdminErr != nil {
					mockIAM.On("IsAdmin", test.RepeatMockAnything(2)...).Return(nil, c.mockIsAdminErr).Once()
				} else {
					mockIAM.On("IsAdmin", test.RepeatMockAnything(2)...).Return(&iam.IsAdminResponse{Ok: c.mockIsAdminResp}, nil).Once()
				}
			}

			if c.callListProject {
				mockProject.On("ListProject", test.RepeatMockAnything(2)...).Return(c.listProjectResp, c.listProjectErr).Once()
			}

			if !c.mockIsAdminResp && c.mockIsAdminErr == nil && c.expectIsAdminCall && (c.mockResponse != nil || c.mockError != nil) {
				mockRepo.On("GetOrgPolicyByUserID", test.RepeatMockAnything(3)...).Return(c.mockResponse, c.mockError).Once()
			}

			result, err := svc.IsAuthorizedOrg(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatal("Expected error but got nil")
			}
			if !c.wantErr && !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestIsAuthorizedOrgToken(t *testing.T) {
	cases := []struct {
		name            string
		input           *org_iam.IsAuthorizedOrgTokenRequest
		want            *org_iam.IsAuthorizedOrgTokenResponse
		wantErr         bool
		callOrgList     bool
		orgListResp     *organization.ListOrganizationResponse
		orgListErr      error
		callMaintainer  bool
		maintainerResp  bool
		maintainerErr   error
		callListProject bool
		listProjectResp *project.ListProjectResponse
		listProjectErr  error
		callGetPolicy   bool
		getPolicyResp   *[]model.OrganizationPolicy
		getPolicyErr    error
	}{
		{
			name: "OK Authorized",
			input: &org_iam.IsAuthorizedOrgTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "organization/update",
			},
			want:           &org_iam.IsAuthorizedOrgTokenResponse{Ok: true},
			callMaintainer: true,
			maintainerResp: true,
			callGetPolicy:  true,
			getPolicyResp: &[]model.OrganizationPolicy{
				{PolicyID: 1, OrganizationID: 1001, ActionPtn: "organization/.*"},
			},
		},
		{
			name: "OK Unauthorized policy not found",
			input: &org_iam.IsAuthorizedOrgTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "organization/delete",
			},
			want:           &org_iam.IsAuthorizedOrgTokenResponse{Ok: false},
			callMaintainer: true,
			maintainerResp: true,
			callGetPolicy:  true,
			getPolicyResp: &[]model.OrganizationPolicy{
				{PolicyID: 1, OrganizationID: 1001, ActionPtn: "organization/(get|list)"},
			},
		},
		{
			name: "OK Record not found",
			input: &org_iam.IsAuthorizedOrgTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "organization/update",
			},
			want:           &org_iam.IsAuthorizedOrgTokenResponse{Ok: false},
			callMaintainer: true,
			maintainerResp: true,
			callGetPolicy:  true,
			getPolicyErr:   gorm.ErrRecordNotFound,
		},
		{
			name: "NG Invalid parameter - action",
			input: &org_iam.IsAuthorizedOrgTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "",
			},
			wantErr: true,
		},
		{
			name: "NG Maintainer check error",
			input: &org_iam.IsAuthorizedOrgTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "organization/update",
			},
			wantErr:        true,
			callMaintainer: true,
			maintainerErr:  gorm.ErrInvalidDB,
		},
		{
			name: "NG Maintainer not found or token expired",
			input: &org_iam.IsAuthorizedOrgTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "organization/update",
			},
			want:           &org_iam.IsAuthorizedOrgTokenResponse{Ok: false},
			callMaintainer: true,
			maintainerResp: false,
		},
		{
			name: "NG Get policy error",
			input: &org_iam.IsAuthorizedOrgTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "organization/update",
			},
			wantErr:        true,
			callMaintainer: true,
			maintainerResp: true,
			callGetPolicy:  true,
			getPolicyErr:   gorm.ErrInvalidDB,
		},
		{
			name: "OK Authorized with project",
			input: &org_iam.IsAuthorizedOrgTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "organization/update",
				ProjectId:      3001,
			},
			want:        &org_iam.IsAuthorizedOrgTokenResponse{Ok: true},
			callOrgList: true,
			orgListResp: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{
					{OrganizationId: 1001},
				},
			},
			callMaintainer:  true,
			maintainerResp:  true,
			callListProject: true,
			listProjectResp: &project.ListProjectResponse{
				Project: []*project.Project{{ProjectId: 3001, Name: "team-a-prod"}},
			},
			callGetPolicy: true,
			getPolicyResp: &[]model.OrganizationPolicy{
				{PolicyID: 1, OrganizationID: 1001, ActionPtn: "organization/.*", ProjectPtn: ".*"},
			},
		},
		{
			name: "OK Unauthorized by project_ptn mismatch",
			input: &org_iam.IsAuthorizedOrgTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "organization/update",
				ProjectId:      3002,
			},
			want:        &org_iam.IsAuthorizedOrgTokenResponse{Ok: false},
			callOrgList: true,
			orgListResp: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{
					{OrganizationId: 1001},
				},
			},
			callMaintainer:  true,
			maintainerResp:  true,
			callListProject: true,
			listProjectResp: &project.ListProjectResponse{
				Project: []*project.Project{{ProjectId: 3002, Name: "team-b-prod"}},
			},
			callGetPolicy: true,
			getPolicyResp: &[]model.OrganizationPolicy{
				{PolicyID: 1, OrganizationID: 1001, ActionPtn: "organization/.*", ProjectPtn: "team-a-.*"},
			},
		},
		{
			name: "OK Unauthorized project not linked",
			input: &org_iam.IsAuthorizedOrgTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "organization/update",
				ProjectId:      3001,
			},
			want:        &org_iam.IsAuthorizedOrgTokenResponse{Ok: false},
			callOrgList: true,
			orgListResp: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{
					{OrganizationId: 9999},
				},
			},
		},
		{
			name: "NG Project check error",
			input: &org_iam.IsAuthorizedOrgTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "organization/update",
				ProjectId:      3001,
			},
			want:        &org_iam.IsAuthorizedOrgTokenResponse{Ok: false},
			callOrgList: true,
			orgListErr:  gorm.ErrInvalidDB,
		},
		{
			name: "OK Unauthorized policy not found with project",
			input: &org_iam.IsAuthorizedOrgTokenRequest{
				OrganizationId: 1001,
				AccessTokenId:  2001,
				ActionName:     "organization/delete",
				ProjectId:      3001,
			},
			want:        &org_iam.IsAuthorizedOrgTokenResponse{Ok: false},
			callOrgList: true,
			orgListResp: &organization.ListOrganizationResponse{
				Organization: []*organization.Organization{
					{OrganizationId: 1001},
				},
			},
			callMaintainer:  true,
			maintainerResp:  true,
			callListProject: true,
			listProjectResp: &project.ListProjectResponse{
				Project: []*project.Project{{ProjectId: 3001, Name: "team-a-prod"}},
			},
			callGetPolicy: true,
			getPolicyResp: &[]model.OrganizationPolicy{
				{PolicyID: 1, OrganizationID: 1001, ActionPtn: "organization/(get|list)", ProjectPtn: ".*"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := mocks.NewOrgIAMRepository(t)
			logger := logging.NewLogger()
			mockIAM := iammock.NewIAMServiceClient(t)
			mockOrg := organizationmock.NewOrganizationServiceClient(t)
			mockProject := projectmock.NewProjectServiceClient(t)
			svc := NewOrgIAMService(mockRepo, mockOrg, mockIAM, mockProject, logger)

			if c.callOrgList {
				mockOrg.On("ListOrganization", test.RepeatMockAnything(2)...).Return(c.orgListResp, c.orgListErr).Once()
			}
			if c.callMaintainer {
				mockRepo.On("ExistsOrgAccessTokenMaintainer", test.RepeatMockAnything(3)...).Return(c.maintainerResp, c.maintainerErr).Once()
			}
			if c.callListProject {
				mockProject.On("ListProject", test.RepeatMockAnything(2)...).Return(c.listProjectResp, c.listProjectErr).Once()
			}
			if c.callGetPolicy {
				mockRepo.On("GetOrgTokenPolicy", test.RepeatMockAnything(3)...).Return(c.getPolicyResp, c.getPolicyErr).Once()
			}

			got, err := svc.IsAuthorizedOrgToken(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if !c.wantErr && !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
