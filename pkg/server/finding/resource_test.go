package finding

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/finding"
	"gorm.io/gorm"
)

func TestListResource(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name       string
		input      *finding.ListResourceRequest
		want       *finding.ListResourceResponse
		wantErr    bool
		mockResp   *[]model.Resource
		mockErr    error
		totalCount *int64
	}{
		{
			name:       "OK",
			input:      &finding.ListResourceRequest{ProjectId: 1, ResourceName: []string{"rn"}, FromAt: now.Unix(), ToAt: now.Unix()},
			want:       &finding.ListResourceResponse{ResourceId: []uint64{1001, 1002}, Count: 2, Total: 999},
			mockResp:   &[]model.Resource{{ResourceID: 1001}, {ResourceID: 1002}},
			totalCount: test.Int64(999),
		},
		{
			name:       "OK Not found",
			input:      &finding.ListResourceRequest{ProjectId: 1},
			want:       &finding.ListResourceResponse{ResourceId: []uint64{}, Count: 0, Total: 0},
			totalCount: test.Int64(0),
		},
		{
			name:       "NG Invalid request",
			input:      &finding.ListResourceRequest{},
			wantErr:    true,
			totalCount: nil,
		},
		{
			name:       "Invalid DB error",
			input:      &finding.ListResourceRequest{ProjectId: 1},
			wantErr:    true,
			mockErr:    gorm.ErrInvalidDB,
			totalCount: test.Int64(999),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.totalCount != nil {
				mockDB.On("ListResourceCount", test.RepeatMockAnything(2)...).Return(*c.totalCount, nil).Once()
			}
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("ListResource", test.RepeatMockAnything(2)...).Return(c.mockResp, c.mockErr).Once()
			}
			got, err := svc.ListResource(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetResource(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *finding.GetResourceRequest
		want         *finding.GetResourceResponse
		mockResponce *model.Resource
		mockError    error
	}{
		{
			name:         "OK",
			input:        &finding.GetResourceRequest{ProjectId: 1, ResourceId: 1001},
			want:         &finding.GetResourceResponse{Resource: &finding.Resource{ResourceId: 1001, ProjectId: 1, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.Resource{ResourceID: 1001, ProjectID: 1, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "NG Record not found",
			input:     &finding.GetResourceRequest{ProjectId: 1, ResourceId: 9999},
			want:      &finding.GetResourceResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetResource", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.GetResource(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestPutResource(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name        string
		input       *finding.PutResourceRequest
		want        *finding.PutResourceResponse
		wantErr     bool
		mockGetResp *model.Resource
		mockGetErr  error
		mockUpResp  *model.Resource
		mockUpErr   error
	}{
		{
			name:       "OK Insert",
			input:      &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceName: "rn", ProjectId: 111}},
			want:       &finding.PutResourceResponse{Resource: &finding.Resource{ResourceId: 1001, ResourceName: "rn", ProjectId: 111, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr: gorm.ErrRecordNotFound,
			mockUpResp: &model.Resource{ResourceID: 1001, ResourceName: "rn", ProjectID: 111, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Update",
			input:       &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceName: "rn-2", ProjectId: 999}},
			want:        &finding.PutResourceResponse{Resource: &finding.Resource{ResourceId: 1001, ResourceName: "rn-2", ProjectId: 999, CreatedAt: now.Add(-1 * 24 * time.Hour).Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.Resource{ResourceID: 1001, ResourceName: "rn-2", ProjectID: 111, CreatedAt: now.Add(-1 * 24 * time.Hour), UpdatedAt: now.Add(-1 * 24 * time.Hour)},
			mockUpResp:  &model.Resource{ResourceID: 1001, ResourceName: "rn-2", ProjectID: 999, CreatedAt: now.Add(-1 * 24 * time.Hour), UpdatedAt: now},
		},
		{
			name:    "NG Invalid request",
			input:   &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceName: "", ProjectId: 111}},
			wantErr: true,
		},
		{
			name:       "NG GetResourceByName error",
			input:      &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceName: "rn", ProjectId: 111}},
			wantErr:    true,
			mockGetErr: gorm.ErrInvalidDB,
		},
		{
			name:        "NG UpsertResource error",
			input:       &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceName: "rn", ProjectId: 111}},
			wantErr:     true,
			mockGetResp: &model.Resource{ResourceID: 1001, ResourceName: "rn", ProjectID: 111, CreatedAt: now, UpdatedAt: now},
			mockUpErr:   gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.mockGetResp != nil || c.mockGetErr != nil {
				mockDB.On("GetResourceByName", test.RepeatMockAnything(3)...).Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertResource", test.RepeatMockAnything(2)...).Return(c.mockUpResp, c.mockUpErr).Once()
			}
			got, err := svc.PutResource(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteResource(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name               string
		input              *finding.DeleteResourceRequest
		wantErr            bool
		mockErr            error
		callDeleteResource bool
	}{
		{
			name:               "OK",
			input:              &finding.DeleteResourceRequest{ProjectId: 1, ResourceId: 1001},
			wantErr:            false,
			callDeleteResource: true,
		},
		{
			name:               "NG validation error",
			input:              &finding.DeleteResourceRequest{ProjectId: 1},
			wantErr:            true,
			callDeleteResource: false,
		},
		{
			name:               "Invalid DB error",
			input:              &finding.DeleteResourceRequest{ProjectId: 1, ResourceId: 1001},
			wantErr:            true,
			mockErr:            gorm.ErrInvalidDB,
			callDeleteResource: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.callDeleteResource {
				mockDB.On("DeleteResource", test.RepeatMockAnything(3)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteResource(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestListResourceTag(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name       string
		input      *finding.ListResourceTagRequest
		want       *finding.ListResourceTagResponse
		wantErr    bool
		mockResp   *[]model.ResourceTag
		mockErr    error
		totalCount *int64
	}{
		{
			name:  "OK",
			input: &finding.ListResourceTagRequest{ProjectId: 1, ResourceId: 1001},
			want: &finding.ListResourceTagResponse{Count: 2, Total: 999,
				Tag: []*finding.ResourceTag{
					{ResourceTagId: 1, ResourceId: 111, Tag: "tag1", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
					{ResourceTagId: 2, ResourceId: 111, Tag: "tag2", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			},
			mockResp: &[]model.ResourceTag{
				{ResourceTagID: 1, ResourceID: 111, Tag: "tag1", CreatedAt: now, UpdatedAt: now},
				{ResourceTagID: 2, ResourceID: 111, Tag: "tag2", CreatedAt: now, UpdatedAt: now},
			},
			totalCount: test.Int64(999),
		},
		{
			name:       "OK Record Not Found",
			input:      &finding.ListResourceTagRequest{ProjectId: 1, ResourceId: 1001},
			want:       &finding.ListResourceTagResponse{Tag: []*finding.ResourceTag{}, Count: 0, Total: 0},
			totalCount: test.Int64(0),
		},
		{
			name:       "NG Invalid Request",
			input:      &finding.ListResourceTagRequest{ProjectId: 1, ResourceId: 0},
			wantErr:    true,
			totalCount: nil,
		},
		{
			name:       "Invalid DB error",
			input:      &finding.ListResourceTagRequest{ProjectId: 1, ResourceId: 1001},
			wantErr:    true,
			mockErr:    gorm.ErrInvalidDB,
			totalCount: test.Int64(999),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.totalCount != nil {
				mockDB.On("ListResourceTagCount", test.RepeatMockAnything(2)...).Return(*c.totalCount, nil).Once()
			}
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("ListResourceTag", test.RepeatMockAnything(2)...).Return(c.mockResp, c.mockErr).Once()
			}
			got, err := svc.ListResourceTag(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestListResourceTagName(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name       string
		input      *finding.ListResourceTagNameRequest
		want       *finding.ListResourceTagNameResponse
		wantErr    bool
		mockResp   *[]db.TagName
		mockErr    error
		totalCount *int64
	}{
		{
			name:       "OK",
			input:      &finding.ListResourceTagNameRequest{ProjectId: 1, FromAt: 0, ToAt: now.Unix()},
			want:       &finding.ListResourceTagNameResponse{Tag: []string{"tag1", "tag2"}, Count: 2, Total: 999},
			mockResp:   &[]db.TagName{{Tag: "tag1"}, {Tag: "tag2"}},
			totalCount: test.Int64(999),
		},
		{
			name:       "OK Record Not Found",
			input:      &finding.ListResourceTagNameRequest{ProjectId: 1},
			want:       &finding.ListResourceTagNameResponse{Tag: []string{}, Count: 0, Total: 0},
			totalCount: test.Int64(0),
		},
		{
			name:       "NG Invalid Request",
			input:      &finding.ListResourceTagNameRequest{},
			wantErr:    true,
			totalCount: nil,
		},
		{
			name:       "Invalid DB error",
			input:      &finding.ListResourceTagNameRequest{ProjectId: 1},
			wantErr:    true,
			mockErr:    gorm.ErrInvalidDB,
			totalCount: test.Int64(999),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.totalCount != nil {
				mockDB.On("ListResourceTagNameCount", test.RepeatMockAnything(2)...).Return(*c.totalCount, nil).Once()
			}
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("ListResourceTagName", test.RepeatMockAnything(2)...).Return(c.mockResp, c.mockErr).Once()
			}
			got, err := svc.ListResourceTagName(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestTagResource(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name        string
		input       *finding.TagResourceRequest
		want        *finding.TagResourceResponse
		wantErr     bool
		mockGetResp *model.ResourceTag
		mockGetErr  error
		mockUpResp  *model.ResourceTag
		mockUpErr   error
	}{
		{
			name:       "OK Insert",
			input:      &finding.TagResourceRequest{ProjectId: 1, Tag: &finding.ResourceTagForUpsert{ResourceId: 1001, ProjectId: 1, Tag: "tag"}},
			want:       &finding.TagResourceResponse{Tag: &finding.ResourceTag{ResourceTagId: 10011, ResourceId: 1001, ProjectId: 1, Tag: "tag", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr: gorm.ErrRecordNotFound,
			mockUpResp: &model.ResourceTag{ResourceTagID: 10011, ResourceID: 1001, ProjectID: 1, Tag: "tag", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Update",
			input:       &finding.TagResourceRequest{ProjectId: 1, Tag: &finding.ResourceTagForUpsert{ResourceId: 1001, ProjectId: 1, Tag: "tag"}},
			want:        &finding.TagResourceResponse{Tag: &finding.ResourceTag{ResourceTagId: 10011, ResourceId: 1001, ProjectId: 1, Tag: "tag", CreatedAt: now.Add(-1 * 24 * time.Hour).Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.ResourceTag{ResourceTagID: 10011, ResourceID: 1001, ProjectID: 1, Tag: "tag", CreatedAt: now.Add(-1 * 24 * time.Hour), UpdatedAt: now.Add(-1 * 24 * time.Hour)},
			mockUpResp:  &model.ResourceTag{ResourceTagID: 10011, ResourceID: 1001, ProjectID: 1, Tag: "tag", CreatedAt: now.Add(-1 * 24 * time.Hour), UpdatedAt: now},
		},
		{
			name:    "NG Invalid request",
			input:   &finding.TagResourceRequest{ProjectId: 1, Tag: &finding.ResourceTagForUpsert{ResourceId: 1001, Tag: "tag"}},
			wantErr: true,
		},
		{
			name:       "NG GetFindingTagByKey error",
			input:      &finding.TagResourceRequest{ProjectId: 1, Tag: &finding.ResourceTagForUpsert{ResourceId: 1001, ProjectId: 1, Tag: "tag"}},
			wantErr:    true,
			mockGetErr: gorm.ErrInvalidDB,
		},
		{
			name:        "NG TagFinding error",
			input:       &finding.TagResourceRequest{ProjectId: 1, Tag: &finding.ResourceTagForUpsert{ResourceId: 1001, ProjectId: 1, Tag: "tag"}},
			wantErr:     true,
			mockGetResp: &model.ResourceTag{ResourceTagID: 10011, ResourceID: 1001, Tag: "tag", CreatedAt: now, UpdatedAt: now},
			mockUpErr:   gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mockDB.On("GetResourceTagByKey", test.RepeatMockAnything(4)...).Return(c.mockGetResp, c.mockGetErr)
			}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("TagResource", test.RepeatMockAnything(2)...).Return(c.mockUpResp, c.mockUpErr)
			}
			got, err := svc.TagResource(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestUntagResource(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name              string
		input             *finding.UntagResourceRequest
		wantErr           bool
		mockErr           error
		callUntagResource bool
	}{
		{
			name:              "OK",
			input:             &finding.UntagResourceRequest{ProjectId: 1, ResourceTagId: 1001},
			wantErr:           false,
			callUntagResource: true,
		},
		{
			name:              "NG validation error",
			input:             &finding.UntagResourceRequest{ProjectId: 1},
			wantErr:           true,
			callUntagResource: false,
		},
		{
			name:              "Invalid DB error",
			input:             &finding.UntagResourceRequest{ProjectId: 1, ResourceTagId: 1001},
			wantErr:           true,
			mockErr:           gorm.ErrInvalidDB,
			callUntagResource: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.callUntagResource {
				mockDB.On("UntagResource", test.RepeatMockAnything(3)...).Return(c.mockErr).Once()
			}
			_, err := svc.UntagResource(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}
