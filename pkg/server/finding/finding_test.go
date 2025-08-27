package finding

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/finding"
	"gorm.io/gorm"
)

func TestListFinding(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name         string
		input        *finding.ListFindingRequest
		want         *finding.ListFindingResponse
		wantErr      bool
		totalCount   int64
		mockResponce *[]model.Finding
		mockError    error
	}{
		{
			name:         "OK",
			input:        &finding.ListFindingRequest{ProjectId: 1, DataSource: []string{"aws:guardduty"}, ResourceName: []string{"hoge"}, FromScore: 0.0, ToScore: 1.0},
			want:         &finding.ListFindingResponse{FindingId: []uint64{111, 222}, Count: 2, Total: 999},
			totalCount:   999,
			mockResponce: &[]model.Finding{{FindingID: 111}, {FindingID: 222}},
		},
		{
			name:       "OK zero list",
			input:      &finding.ListFindingRequest{ProjectId: 1, DataSource: []string{"aws:guardduty"}, ResourceName: []string{"hoge"}, FromScore: 0.0, ToScore: 1.0},
			want:       &finding.ListFindingResponse{FindingId: []uint64{}, Count: 0, Total: 0},
			totalCount: 0,
		},
		{
			name:       "Invalid DB error",
			input:      &finding.ListFindingRequest{ProjectId: 1, DataSource: []string{"aws:guardduty"}, ResourceName: []string{"hoge"}, FromScore: 0.0, ToScore: 1.0},
			totalCount: 999,
			mockError:  gorm.ErrInvalidDB,
			wantErr:    true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			mockDB.On("ListFindingCount", test.RepeatMockAnything(12)...).Return(c.totalCount, nil).Once()
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListFinding", test.RepeatMockAnything(2)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListFinding(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBatchListFinding(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name         string
		input        *finding.BatchListFindingRequest
		want         *finding.BatchListFindingResponse
		wantErr      bool
		totalCount   int64
		mockResponce *[]model.Finding
		mockError    error
	}{
		{
			name:         "OK",
			input:        &finding.BatchListFindingRequest{ProjectId: 1, DataSource: []string{"aws:guardduty"}, ResourceName: []string{"hoge"}, FromScore: 0.0, ToScore: 1.0},
			want:         &finding.BatchListFindingResponse{FindingId: []uint64{111, 222}, Count: 2, Total: 2},
			totalCount:   2,
			mockResponce: &[]model.Finding{{FindingID: 111}, {FindingID: 222}},
		},
		{
			name:       "OK zero list",
			input:      &finding.BatchListFindingRequest{ProjectId: 1, DataSource: []string{"aws:guardduty"}, ResourceName: []string{"hoge"}, FromScore: 0.0, ToScore: 1.0},
			want:       &finding.BatchListFindingResponse{FindingId: []uint64{}, Count: 0, Total: 0},
			totalCount: 0,
		},
		{
			name:       "Invalid DB error",
			input:      &finding.BatchListFindingRequest{ProjectId: 1, DataSource: []string{"aws:guardduty"}, ResourceName: []string{"hoge"}, FromScore: 0.0, ToScore: 1.0},
			totalCount: 999,
			mockError:  gorm.ErrInvalidDB,
			wantErr:    true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			mockDB.On("ListFindingCount", test.RepeatMockAnything(12)...).Return(c.totalCount, nil).Once()
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("BatchListFinding", test.RepeatMockAnything(2)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.BatchListFinding(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetFinding(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name         string
		input        *finding.GetFindingRequest
		want         *finding.GetFindingResponse
		mockResponce *model.Finding
		mockError    error
	}{
		{
			name:         "OK",
			input:        &finding.GetFindingRequest{ProjectId: 1, FindingId: 1001},
			want:         &finding.GetFindingResponse{Finding: &finding.Finding{FindingId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.Finding{FindingID: 1001, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "NG record not found",
			input:     &finding.GetFindingRequest{ProjectId: 1, FindingId: 9999},
			want:      &finding.GetFindingResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetFinding", test.RepeatMockAnything(4)...).Return(c.mockResponce, c.mockError).Once()
			}
			result, err := svc.GetFinding(ctx, c.input)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestPutFinding(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name        string
		input       *finding.PutFindingRequest
		want        *finding.PutFindingResponse
		wantErr     bool
		mockGetResp *model.Finding
		mockGetErr  error
		mockUpResp  *model.Finding
		mockUpErr   error

		callListFindingSetting bool
		callGetResourceByName  bool
		callUpsertResource     bool
	}{
		{
			name:                   "OK Insert",
			input:                  &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, OriginalMaxScore: 100.00}},
			want:                   &finding.PutFindingResponse{Finding: &finding.Finding{FindingId: 1001, DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, Score: 1.0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr:             gorm.ErrRecordNotFound,
			mockUpResp:             &model.Finding{FindingID: 1001, DataSource: "ds", DataSourceID: "ds-001", ResourceName: "rn", OriginalScore: 100.00, Score: 1.0, CreatedAt: now, UpdatedAt: now},
			callListFindingSetting: true,
			callGetResourceByName:  true,
			callUpsertResource:     true,
		},
		{
			name:                   "OK Update",
			input:                  &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 20.00, OriginalMaxScore: 100.00}},
			want:                   &finding.PutFindingResponse{Finding: &finding.Finding{FindingId: 1001, DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 20.00, Score: 0.2, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp:            &model.Finding{FindingID: 1001, DataSource: "ds", DataSourceID: "ds-001", ResourceName: "rn", OriginalScore: 10.00, Score: 0.1, CreatedAt: now, UpdatedAt: now},
			mockUpResp:             &model.Finding{FindingID: 1001, DataSource: "ds", DataSourceID: "ds-001", ResourceName: "rn", OriginalScore: 20.00, Score: 0.2, CreatedAt: now, UpdatedAt: now},
			callListFindingSetting: true,
			callGetResourceByName:  true,
			callUpsertResource:     true,
		},
		{
			name:                   "NG Invalid request(no data_source)",
			input:                  &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{DataSource: "", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, OriginalMaxScore: 100.00}},
			wantErr:                true,
			callListFindingSetting: false,
			callGetResourceByName:  false,
			callUpsertResource:     false,
		},
		{
			name:                   "NG GetFindingByDataSource error",
			input:                  &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, OriginalMaxScore: 100.00}},
			wantErr:                true,
			mockGetErr:             gorm.ErrInvalidDB,
			callListFindingSetting: false,
			callGetResourceByName:  false,
			callUpsertResource:     false,
		},
		{
			name:                   "NG UpsertFinding error",
			input:                  &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, OriginalMaxScore: 100.00}},
			wantErr:                true,
			mockGetResp:            &model.Finding{FindingID: 1001, DataSource: "ds", DataSourceID: "ds-001", ResourceName: "rn", OriginalScore: 10.00, Score: 0.1, CreatedAt: now, UpdatedAt: now},
			mockUpErr:              gorm.ErrInvalidDB,
			callListFindingSetting: true,
			callGetResourceByName:  false,
			callUpsertResource:     false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.mockGetResp != nil || c.mockGetErr != nil {
				mockDB.On("GetFindingByDataSource", test.RepeatMockAnything(4)...).Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertFinding", test.RepeatMockAnything(2)...).Return(c.mockUpResp, c.mockUpErr).Once()
			}
			if c.callListFindingSetting {
				mockDB.On("ListFindingSetting", test.RepeatMockAnything(3)...).Return(&[]model.FindingSetting{{ResourceName: "rn", Setting: `{"score_coefficient": 0.1}`}}, nil) // fixed response
			}
			if c.callGetResourceByName {
				mockDB.On("GetResourceByName", test.RepeatMockAnything(3)...).Return(&model.Resource{}, nil) // fixed response
			}
			if c.callUpsertResource {
				mockDB.On("UpsertResource", test.RepeatMockAnything(2)...).Return(&model.Resource{}, nil) // fixed response
			}

			got, err := svc.PutFinding(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteFinding(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name     string
		input    *finding.DeleteFindingRequest
		wantErr  bool
		mockCall bool
		mockErr  error
	}{
		{
			name:     "OK",
			input:    &finding.DeleteFindingRequest{ProjectId: 1, FindingId: 1001},
			wantErr:  false,
			mockCall: true,
			mockErr:  nil,
		},
		{
			name:     "NG validation error",
			input:    &finding.DeleteFindingRequest{ProjectId: 1},
			mockCall: false,
			wantErr:  true,
		},
		{
			name:     "Invalid DB error",
			input:    &finding.DeleteFindingRequest{ProjectId: 1, FindingId: 1001},
			wantErr:  true,
			mockCall: true,
			mockErr:  gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.mockCall {
				mockDB.On("DeleteFinding", test.RepeatMockAnything(3)...).Return(c.mockErr).Once()
			}
			_, err := svc.DeleteFinding(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestListFindingTag(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name       string
		input      *finding.ListFindingTagRequest
		want       *finding.ListFindingTagResponse
		wantErr    bool
		mockResp   *[]model.FindingTag
		mockErr    error
		totalCount *int64
	}{
		{
			name:  "OK",
			input: &finding.ListFindingTagRequest{ProjectId: 1, FindingId: 1001},
			want: &finding.ListFindingTagResponse{Count: 2, Total: 999,
				Tag: []*finding.FindingTag{
					{FindingTagId: 1, FindingId: 111, Tag: "tag1", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
					{FindingTagId: 2, FindingId: 111, Tag: "tag2", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				},
			},
			mockResp: &[]model.FindingTag{
				{FindingTagID: 1, FindingID: 111, Tag: "tag1", CreatedAt: now, UpdatedAt: now},
				{FindingTagID: 2, FindingID: 111, Tag: "tag2", CreatedAt: now, UpdatedAt: now},
			},
			totalCount: test.Int64(999),
		},
		{
			name:       "OK Record Not Found",
			input:      &finding.ListFindingTagRequest{ProjectId: 1, FindingId: 1001},
			want:       &finding.ListFindingTagResponse{Tag: []*finding.FindingTag{}, Count: 0, Total: 0},
			totalCount: test.Int64(0),
		},
		{
			name:       "NG Invalid Request",
			input:      &finding.ListFindingTagRequest{ProjectId: 1, FindingId: 0},
			wantErr:    true,
			totalCount: nil,
		},
		{
			name:       "Invalid DB error",
			input:      &finding.ListFindingTagRequest{ProjectId: 1, FindingId: 1001},
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
				mockDB.On("ListFindingTagCount", test.RepeatMockAnything(2)...).Return(*c.totalCount, nil).Once()
			}
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("ListFindingTag", test.RepeatMockAnything(2)...).Return(c.mockResp, c.mockErr).Once()
			}
			got, err := svc.ListFindingTag(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestListFindingTagName(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name       string
		input      *finding.ListFindingTagNameRequest
		want       *finding.ListFindingTagNameResponse
		wantErr    bool
		mockResp   *[]db.TagName
		mockErr    error
		totalCount *int64
	}{
		{
			name:       "OK",
			input:      &finding.ListFindingTagNameRequest{ProjectId: 1},
			want:       &finding.ListFindingTagNameResponse{Count: 2, Total: 999, Tag: []string{"tag1", "tag2"}},
			mockResp:   &[]db.TagName{{Tag: "tag1"}, {Tag: "tag2"}},
			totalCount: test.Int64(999),
		},
		{
			name:       "OK Record Not Found",
			input:      &finding.ListFindingTagNameRequest{ProjectId: 1},
			want:       &finding.ListFindingTagNameResponse{Tag: []string{}, Count: 0, Total: 0},
			totalCount: test.Int64(0),
		},
		{
			name:       "NG Invalid Request",
			input:      &finding.ListFindingTagNameRequest{},
			wantErr:    true,
			totalCount: nil,
		},
		{
			name:       "Invalid DB error",
			input:      &finding.ListFindingTagNameRequest{ProjectId: 1},
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
				mockDB.On("ListFindingTagNameCount", test.RepeatMockAnything(2)...).Return(*c.totalCount, nil).Once()
			}
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("ListFindingTagName", test.RepeatMockAnything(2)...).Return(c.mockResp, c.mockErr).Once()
			}
			got, err := svc.ListFindingTagName(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestTagFinding(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	cases := []struct {
		name    string
		input   *finding.TagFindingRequest
		want    *finding.TagFindingResponse
		wantErr bool

		// mock
		callGetResourceTagByKey    bool
		callTagResource            bool
		callGetFinding             bool
		mockGetFindingResp         *model.Finding
		mockGetFindingErr          error
		callGetResourceByName      bool
		mockGetResourceByNameResp  *model.Resource
		mockGetResourceByNameErr   error
		callGetFindingTagByKey     bool
		mockGetFindingTagByKeyResp *model.FindingTag
		mockGetFindingTagByKeyErr  error
		callTagFinding             bool
		mockTagFindingResp         *model.FindingTag
		mockTagFindingErr          error
	}{
		{
			name:  "OK Insert",
			input: &finding.TagFindingRequest{ProjectId: 1, Tag: &finding.FindingTagForUpsert{FindingId: 1001, ProjectId: 1, Tag: "tag"}},
			want:  &finding.TagFindingResponse{Tag: &finding.FindingTag{FindingTagId: 10011, ProjectId: 1, FindingId: 1001, Tag: "tag", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},

			callGetResourceTagByKey:   false,
			callTagResource:           true,
			callGetFinding:            true,
			mockGetFindingResp:        &model.Finding{FindingID: 1001, ProjectID: 1, ResourceName: "rn"},
			callGetResourceByName:     true,
			mockGetResourceByNameResp: &model.Resource{ResourceID: 1001, ProjectID: 1, ResourceName: "rn"},
			callGetFindingTagByKey:    true,
			mockGetFindingTagByKeyErr: gorm.ErrRecordNotFound,
			callTagFinding:            true,
			mockTagFindingResp:        &model.FindingTag{FindingTagID: 10011, FindingID: 1001, ProjectID: 1, Tag: "tag", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:  "OK Update",
			input: &finding.TagFindingRequest{ProjectId: 1, Tag: &finding.FindingTagForUpsert{FindingId: 1001, ProjectId: 1, Tag: "tag"}},
			want:  &finding.TagFindingResponse{Tag: &finding.FindingTag{FindingTagId: 10011, ProjectId: 1, FindingId: 1001, Tag: "tag", CreatedAt: now.Add(-1 * 24 * time.Hour).Unix(), UpdatedAt: now.Unix()}},

			callGetResourceTagByKey:    false,
			callTagResource:            true,
			callGetFinding:             true,
			mockGetFindingResp:         &model.Finding{FindingID: 1001, ProjectID: 1, ResourceName: "rn"},
			callGetResourceByName:      true,
			mockGetResourceByNameResp:  &model.Resource{ResourceID: 1001, ProjectID: 1, ResourceName: "rn"},
			callGetFindingTagByKey:     true,
			mockGetFindingTagByKeyResp: &model.FindingTag{FindingTagID: 10011, FindingID: 1001, ProjectID: 1, Tag: "tag", CreatedAt: now.Add(-1 * 24 * time.Hour), UpdatedAt: now.Add(-1 * 24 * time.Hour)},
			callTagFinding:             true,
			mockTagFindingResp:         &model.FindingTag{FindingTagID: 10011, FindingID: 1001, ProjectID: 1, Tag: "tag", CreatedAt: now.Add(-1 * 24 * time.Hour), UpdatedAt: now},
		},
		{
			name:    "NG Invalid request",
			input:   &finding.TagFindingRequest{ProjectId: 1, Tag: &finding.FindingTagForUpsert{FindingId: 1001, Tag: ""}},
			wantErr: true,
		},
		{
			name:    "NG GetFinding error",
			input:   &finding.TagFindingRequest{ProjectId: 1, Tag: &finding.FindingTagForUpsert{FindingId: 1001, ProjectId: 1, Tag: "tag"}},
			wantErr: true,

			callGetFinding:    true,
			mockGetFindingErr: gorm.ErrInvalidDB,
		},
		{
			name:    "NG GetFindingTagByKey error",
			input:   &finding.TagFindingRequest{ProjectId: 1, Tag: &finding.FindingTagForUpsert{FindingId: 1001, ProjectId: 1, Tag: "tag"}},
			wantErr: true,

			callGetFinding:            true,
			mockGetFindingResp:        &model.Finding{FindingID: 1001, ProjectID: 1, ResourceName: "rn"},
			callGetFindingTagByKey:    true,
			mockGetFindingTagByKeyErr: gorm.ErrInvalidDB,
		},
		{
			name:    "NG TagFinding error",
			input:   &finding.TagFindingRequest{ProjectId: 1, Tag: &finding.FindingTagForUpsert{FindingId: 1001, ProjectId: 1, Tag: "tag"}},
			wantErr: true,

			callGetFinding:             true,
			mockGetFindingResp:         &model.Finding{FindingID: 1001, ProjectID: 1, ResourceName: "rn"},
			callGetFindingTagByKey:     true,
			mockGetFindingTagByKeyResp: &model.FindingTag{FindingTagID: 10011, FindingID: 1001, Tag: "tag", CreatedAt: now, UpdatedAt: now},
			callTagFinding:             true,
			mockTagFindingErr:          gorm.ErrInvalidDB,
		},
		{
			name:    "NG GetResourceByName error",
			input:   &finding.TagFindingRequest{ProjectId: 1, Tag: &finding.FindingTagForUpsert{FindingId: 1001, ProjectId: 1, Tag: "tag"}},
			wantErr: true,

			callGetFinding:            true,
			mockGetFindingResp:        &model.Finding{FindingID: 1001, ProjectID: 1, ResourceName: "rn"},
			callGetFindingTagByKey:    true,
			callTagFinding:            true,
			mockTagFindingResp:        &model.FindingTag{FindingTagID: 10011, FindingID: 1001, ProjectID: 1, Tag: "tag", CreatedAt: now, UpdatedAt: now},
			mockGetFindingTagByKeyErr: gorm.ErrRecordNotFound,
			callGetResourceByName:     true,
			mockGetResourceByNameErr:  gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.callGetResourceTagByKey {
				mockDB.On("GetResourceTagByKey", test.RepeatMockAnything(2)...).Return(&model.ResourceTag{}, gorm.ErrRecordNotFound)
			}
			if c.callTagResource {
				mockDB.On("TagResource", test.RepeatMockAnything(2)...).Return(&model.ResourceTag{}, nil)
			}

			if c.callGetFinding {
				mockDB.On("GetFinding", test.RepeatMockAnything(4)...).Return(c.mockGetFindingResp, c.mockGetFindingErr)
			}
			if c.callGetResourceByName {
				mockDB.On("GetResourceByName", test.RepeatMockAnything(3)...).Return(c.mockGetResourceByNameResp, c.mockGetResourceByNameErr)
			}
			if c.callGetFindingTagByKey {
				mockDB.On("GetFindingTagByKey", test.RepeatMockAnything(4)...).Return(c.mockGetFindingTagByKeyResp, c.mockGetFindingTagByKeyErr)
			}
			if c.callTagFinding {
				mockDB.On("TagFinding", test.RepeatMockAnything(2)...).Return(c.mockTagFindingResp, c.mockTagFindingErr)
			}
			got, err := svc.TagFinding(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestClearSocre(t *testing.T) {
	var ctx context.Context
	cases := []struct {
		name     string
		input    *finding.ClearScoreRequest
		wantErr  bool
		execMock bool
		mockResp error
	}{
		{
			name:     "OK",
			input:    &finding.ClearScoreRequest{DataSource: "ds", ProjectId: 1, Tag: []string{"tag1", "tag2"}, FindingId: 1},
			execMock: true,
			wantErr:  false,
			mockResp: nil,
		},
		{
			name:     "NG Invalid request",
			input:    &finding.ClearScoreRequest{}, // Required param error
			execMock: false,
			wantErr:  true,
		},
		{
			name:     "NG DB error",
			input:    &finding.ClearScoreRequest{DataSource: "ds", ProjectId: 1, Tag: []string{"tag1", "tag2"}, FindingId: 1},
			execMock: true,
			wantErr:  true,
			mockResp: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB, logger: logging.NewLogger()}

			if c.execMock {
				mockDB.On("ClearScoreFinding", test.RepeatMockAnything(2)...).Return(c.mockResp).Once()
			}
			_, err := svc.ClearScore(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestCalculateScore(t *testing.T) {
	cases := []struct {
		name    string
		input   [2]float32
		setting *findingSetting
		want    float32
	}{
		{
			name:  "OK Score 1%",
			input: [2]float32{1.0, 100.0},
			want:  0.01,
		},
		{
			name:  "OK Score 100%",
			input: [2]float32{100.0, 100.0},
			want:  1.00,
		},
		{
			name:  "OK Score 0%",
			input: [2]float32{0, 100.0},
			want:  0.00,
		},
		{
			name:    "OK Setting x1",
			input:   [2]float32{0.1, 1.0},
			setting: &findingSetting{ScoreCoefficient: 1.0},
			want:    0.1,
		},
		{
			name:    "OK Setting x1.5",
			input:   [2]float32{0.1, 1.0},
			setting: &findingSetting{ScoreCoefficient: 1.5},
			want:    0.15,
		},
		{
			name:    "OK Setting x100",
			input:   [2]float32{0.1, 1.0},
			setting: &findingSetting{ScoreCoefficient: 100},
			want:    1.0,
		},
		{
			name:    "OK Setting x-1",
			input:   [2]float32{0.1, 1.0},
			setting: &findingSetting{ScoreCoefficient: -1},
			want:    0.0,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := calculateScore(c.input[0], c.input[1], c.setting)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected result: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
