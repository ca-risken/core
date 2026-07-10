package ai

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	aipb "github.com/ca-risken/core/proto/ai"
	"gorm.io/gorm"
)

func TestGetRemediationProposal(t *testing.T) {
	now := time.Now()
	plan := `{"summary":"test"}`
	cases := []struct {
		name         string
		input        *aipb.GetRemediationProposalRequest
		want         *aipb.GetRemediationProposalResponse
		mockResponse *model.RemediationProposal
		mockError    error
		wantErr      bool
	}{
		{
			name:  "OK",
			input: &aipb.GetRemediationProposalRequest{ProjectId: 1, RemediationProposalId: 1001},
			want: &aipb.GetRemediationProposalResponse{
				RemediationProposal: &aipb.RemediationProposal{
					RemediationProposalId: 1001,
					FindingId:             2001,
					ProjectId:             1,
					Status:                model.RemediationProposalStatusSucceeded,
					RemediationPlan:       plan,
					GeneratedAt:           now.Unix(),
					CreatedAt:             now.Unix(),
					UpdatedAt:             now.Unix(),
				},
			},
			mockResponse: &model.RemediationProposal{
				RemediationProposalID: 1001,
				FindingID:             2001,
				ProjectID:             1,
				Status:                model.RemediationProposalStatusSucceeded,
				RemediationPlan:       &plan,
				GeneratedAt:           &now,
				CreatedAt:             now,
				UpdatedAt:             now,
			},
		},
		{
			name:      "OK not found",
			input:     &aipb.GetRemediationProposalRequest{ProjectId: 1, RemediationProposalId: 9999},
			want:      &aipb.GetRemediationProposalResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG validation error",
			input:   &aipb.GetRemediationProposalRequest{ProjectId: 0, RemediationProposalId: 1001},
			wantErr: true,
		},
		{
			name:      "NG DB error",
			input:     &aipb.GetRemediationProposalRequest{ProjectId: 1, RemediationProposalId: 1001},
			mockError: errors.New("DB error"),
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAIRepository(t)
			svc := AIService{repository: mockDB}
			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("GetRemediationProposal", test.RepeatMockAnything(3)...).Return(c.mockResponse, c.mockError).Once()
			}
			result, err := svc.GetRemediationProposal(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestListRemediationProposal(t *testing.T) {
	now := time.Now()
	statusDetail := "failed"
	cases := []struct {
		name         string
		input        *aipb.ListRemediationProposalRequest
		want         *aipb.ListRemediationProposalResponse
		mockResponse []*model.RemediationProposal
		mockError    error
		wantErr      bool
	}{
		{
			name:  "OK",
			input: &aipb.ListRemediationProposalRequest{ProjectId: 1, FindingId: 2001},
			want: &aipb.ListRemediationProposalResponse{
				RemediationProposal: []*aipb.RemediationProposal{
					{
						RemediationProposalId: 1002,
						FindingId:             2001,
						ProjectId:             1,
						Status:                model.RemediationProposalStatusPending,
						CreatedAt:             now.Unix(),
						UpdatedAt:             now.Unix(),
					},
					{
						RemediationProposalId: 1001,
						FindingId:             2001,
						ProjectId:             1,
						Status:                model.RemediationProposalStatusFailed,
						StatusDetail:          statusDetail,
						CreatedAt:             now.Add(-time.Hour).Unix(),
						UpdatedAt:             now.Add(-time.Hour).Unix(),
					},
				},
			},
			mockResponse: []*model.RemediationProposal{
				{
					RemediationProposalID: 1002,
					FindingID:             2001,
					ProjectID:             1,
					Status:                model.RemediationProposalStatusPending,
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					RemediationProposalID: 1001,
					FindingID:             2001,
					ProjectID:             1,
					Status:                model.RemediationProposalStatusFailed,
					StatusDetail:          &statusDetail,
					CreatedAt:             now.Add(-time.Hour),
					UpdatedAt:             now.Add(-time.Hour),
				},
			},
		},
		{
			name:         "OK empty",
			input:        &aipb.ListRemediationProposalRequest{ProjectId: 1, FindingId: 2001, Status: []string{model.RemediationProposalStatusPending}},
			want:         &aipb.ListRemediationProposalResponse{},
			mockResponse: []*model.RemediationProposal{},
		},
		{
			name:    "NG validation error",
			input:   &aipb.ListRemediationProposalRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG invalid status",
			input:   &aipb.ListRemediationProposalRequest{ProjectId: 1, FindingId: 2001, Status: []string{"unknown"}},
			wantErr: true,
		},
		{
			name:      "NG DB error",
			input:     &aipb.ListRemediationProposalRequest{ProjectId: 1, FindingId: 2001},
			mockError: errors.New("DB error"),
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAIRepository(t)
			svc := AIService{repository: mockDB}
			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("ListRemediationProposal", test.RepeatMockAnything(4)...).Return(c.mockResponse, c.mockError).Once()
			}
			result, err := svc.ListRemediationProposal(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestPutRemediationProposal(t *testing.T) {
	now := time.Now()
	plan := `{"summary":"test"}`
	pendingModel := &model.RemediationProposal{
		RemediationProposalID: 1001,
		FindingID:             2001,
		ProjectID:             1,
		Status:                model.RemediationProposalStatusPending,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	succeededModel := &model.RemediationProposal{
		RemediationProposalID: 1001,
		FindingID:             2001,
		ProjectID:             1,
		Status:                model.RemediationProposalStatusSucceeded,
		RemediationPlan:       &plan,
		GeneratedAt:           &now,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	cases := []struct {
		name           string
		input          *aipb.PutRemediationProposalRequest
		want           *aipb.PutRemediationProposalResponse
		mockCreateResp *model.RemediationProposal
		mockCreateErr  error
		mockUpdateResp *model.RemediationProposal
		mockUpdateErr  error
		wantErr        bool
	}{
		{
			name: "OK create pending proposal",
			input: &aipb.PutRemediationProposalRequest{
				ProjectId: 1,
				FindingId: 2001,
				Status:    model.RemediationProposalStatusPending,
			},
			mockCreateResp: pendingModel,
			want: &aipb.PutRemediationProposalResponse{
				RemediationProposal: &aipb.RemediationProposal{
					RemediationProposalId: 1001,
					FindingId:             2001,
					ProjectId:             1,
					Status:                model.RemediationProposalStatusPending,
					CreatedAt:             now.Unix(),
					UpdatedAt:             now.Unix(),
				},
			},
		},
		{
			name: "OK update succeeded proposal",
			input: &aipb.PutRemediationProposalRequest{
				ProjectId:             1,
				RemediationProposalId: 1001,
				FindingId:             2001,
				Status:                model.RemediationProposalStatusSucceeded,
				RemediationPlan:       plan,
				GeneratedAt:           now.Unix(),
			},
			mockUpdateResp: succeededModel,
			want: &aipb.PutRemediationProposalResponse{
				RemediationProposal: &aipb.RemediationProposal{
					RemediationProposalId: 1001,
					FindingId:             2001,
					ProjectId:             1,
					Status:                model.RemediationProposalStatusSucceeded,
					RemediationPlan:       plan,
					GeneratedAt:           now.Unix(),
					CreatedAt:             now.Unix(),
					UpdatedAt:             now.Unix(),
				},
			},
		},
		{
			name: "NG validation error",
			input: &aipb.PutRemediationProposalRequest{
				ProjectId: 1,
				FindingId: 2001,
				Status:    "unknown",
			},
			wantErr: true,
		},
		{
			name: "NG completed status without remediation_proposal_id",
			input: &aipb.PutRemediationProposalRequest{
				ProjectId: 1,
				FindingId: 2001,
				Status:    model.RemediationProposalStatusSucceeded,
			},
			wantErr: true,
		},
		{
			name: "NG DB error(create)",
			input: &aipb.PutRemediationProposalRequest{
				ProjectId: 1,
				FindingId: 2001,
				Status:    model.RemediationProposalStatusPending,
			},
			mockCreateErr: errors.New("DB error"),
			wantErr:       true,
		},
		{
			name: "NG DB error(update)",
			input: &aipb.PutRemediationProposalRequest{
				ProjectId:             1,
				RemediationProposalId: 1001,
				FindingId:             2001,
				Status:                model.RemediationProposalStatusSucceeded,
			},
			mockUpdateErr: errors.New("DB error"),
			wantErr:       true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAIRepository(t)
			svc := AIService{repository: mockDB}
			if c.mockCreateResp != nil || c.mockCreateErr != nil {
				mockDB.On("CreateRemediationProposal", test.RepeatMockAnything(2)...).Return(c.mockCreateResp, c.mockCreateErr).Once()
			}
			if c.mockUpdateResp != nil || c.mockUpdateErr != nil {
				mockDB.On("UpdateRemediationProposalStatus", test.RepeatMockAnything(7)...).Return(c.mockUpdateResp, c.mockUpdateErr).Once()
			}
			result, err := svc.PutRemediationProposal(context.Background(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}
