package main

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"gorm.io/gorm"
)

func TestListFinding(t *testing.T) {
	var ctx context.Context
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
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
			mockDB.On("ListFindingCount").Return(c.totalCount, nil).Once()
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListFinding").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
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
			mockDB.On("ListFindingCount").Return(c.totalCount, nil).Once()
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("BatchListFinding").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetFinding").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
	// Resource関連, FindingSettingのupdateは別テストで実施。ここでは一律カラを返す
	mockDB.On("GetFindingSettingByResource").Return(&model.FindingSetting{}, nil)
	mockDB.On("GetResourceByName").Return(&model.Resource{}, nil)
	mockDB.On("UpsertResource").Return(&model.Resource{}, nil)

	cases := []struct {
		name        string
		input       *finding.PutFindingRequest
		want        *finding.PutFindingResponse
		wantErr     bool
		mockGetResp *model.Finding
		mockGetErr  error
		mockUpResp  *model.Finding
		mockUpErr   error
	}{
		{
			name:       "OK Insert",
			input:      &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, OriginalMaxScore: 100.00}},
			want:       &finding.PutFindingResponse{Finding: &finding.Finding{FindingId: 1001, DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, Score: 1.0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr: gorm.ErrRecordNotFound,
			mockUpResp: &model.Finding{FindingID: 1001, DataSource: "ds", DataSourceID: "ds-001", ResourceName: "rn", OriginalScore: 100.00, Score: 1.0, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Update",
			input:       &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 20.00, OriginalMaxScore: 100.00}},
			want:        &finding.PutFindingResponse{Finding: &finding.Finding{FindingId: 1001, DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 20.00, Score: 0.2, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.Finding{FindingID: 1001, DataSource: "ds", DataSourceID: "ds-001", ResourceName: "rn", OriginalScore: 10.00, Score: 0.1, CreatedAt: now, UpdatedAt: now},
			mockUpResp:  &model.Finding{FindingID: 1001, DataSource: "ds", DataSourceID: "ds-001", ResourceName: "rn", OriginalScore: 20.00, Score: 0.2, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:    "NG Invalid request(no data_source)",
			input:   &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{DataSource: "", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, OriginalMaxScore: 100.00}},
			wantErr: true,
		},
		{
			name:       "NG GetFindingByDataSource error",
			input:      &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, OriginalMaxScore: 100.00}},
			wantErr:    true,
			mockGetErr: gorm.ErrInvalidDB,
		},
		{
			name:        "NG UpsertFinding error",
			input:       &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, OriginalMaxScore: 100.00}},
			wantErr:     true,
			mockGetResp: &model.Finding{FindingID: 1001, DataSource: "ds", DataSourceID: "ds-001", ResourceName: "rn", OriginalScore: 10.00, Score: 0.1, CreatedAt: now, UpdatedAt: now},
			mockUpErr:   gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mockDB.On("GetFindingByDataSource").Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertFinding").Return(c.mockUpResp, c.mockUpErr).Once()
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
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
	cases := []struct {
		name    string
		input   *finding.DeleteFindingRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &finding.DeleteFindingRequest{ProjectId: 1, FindingId: 1001},
			wantErr: false,
			mockErr: nil,
		},
		{
			name:    "NG validation error",
			input:   &finding.DeleteFindingRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "Invalid DB error",
			input:   &finding.DeleteFindingRequest{ProjectId: 1, FindingId: 1001},
			wantErr: true,
			mockErr: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeleteFinding").Return(c.mockErr).Once()
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
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
	cases := []struct {
		name       string
		input      *finding.ListFindingTagRequest
		want       *finding.ListFindingTagResponse
		wantErr    bool
		mockResp   *[]model.FindingTag
		mockErr    error
		totalCount int64
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
			totalCount: 999,
		},
		{
			name:       "OK Record Not Found",
			input:      &finding.ListFindingTagRequest{ProjectId: 1, FindingId: 1001},
			want:       &finding.ListFindingTagResponse{Tag: []*finding.FindingTag{}, Count: 0, Total: 0},
			totalCount: 0,
		},
		{
			name:       "NG Invalid Request",
			input:      &finding.ListFindingTagRequest{ProjectId: 1, FindingId: 0},
			wantErr:    true,
			totalCount: 999,
		},
		{
			name:       "Invalid DB error",
			input:      &finding.ListFindingTagRequest{ProjectId: 1, FindingId: 1001},
			wantErr:    true,
			mockErr:    gorm.ErrInvalidDB,
			totalCount: 999,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("ListFindingTagCount").Return(c.totalCount, nil).Once()
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("ListFindingTag").Return(c.mockResp, c.mockErr).Once()
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
	now := time.Now()
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
	cases := []struct {
		name       string
		input      *finding.ListFindingTagNameRequest
		want       *finding.ListFindingTagNameResponse
		wantErr    bool
		mockResp   *[]tagName
		mockErr    error
		totalCount int64
	}{
		{
			name:  "OK",
			input: &finding.ListFindingTagNameRequest{ProjectId: 1, FromAt: 0, ToAt: now.Unix()},
			want: &finding.ListFindingTagNameResponse{Count: 2, Total: 999,
				Tag: []string{"tag1", "tag2"}},
			mockResp:   &[]tagName{{Tag: "tag1"}, {Tag: "tag2"}},
			totalCount: 999,
		},
		{
			name:       "OK Record Not Found",
			input:      &finding.ListFindingTagNameRequest{ProjectId: 1},
			want:       &finding.ListFindingTagNameResponse{Tag: []string{}, Count: 0, Total: 0},
			totalCount: 0,
		},
		{
			name:       "NG Invalid Request",
			input:      &finding.ListFindingTagNameRequest{},
			wantErr:    true,
			totalCount: 999,
		},
		{
			name:       "Invalid DB error",
			input:      &finding.ListFindingTagNameRequest{ProjectId: 1},
			wantErr:    true,
			mockErr:    gorm.ErrInvalidDB,
			totalCount: 999,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("ListFindingTagNameCount").Return(c.totalCount, nil).Once()
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("ListFindingTagName").Return(c.mockResp, c.mockErr).Once()
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
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
	cases := []struct {
		name        string
		input       *finding.TagFindingRequest
		want        *finding.TagFindingResponse
		wantErr     bool
		mockGetResp *model.FindingTag
		mockGetErr  error
		mockUpResp  *model.FindingTag
		mockUpErr   error
	}{
		{
			name:       "OK Insert",
			input:      &finding.TagFindingRequest{ProjectId: 1, Tag: &finding.FindingTagForUpsert{FindingId: 1001, ProjectId: 1, Tag: "tag"}},
			want:       &finding.TagFindingResponse{Tag: &finding.FindingTag{FindingTagId: 10011, ProjectId: 1, FindingId: 1001, Tag: "tag", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr: gorm.ErrRecordNotFound,
			mockUpResp: &model.FindingTag{FindingTagID: 10011, FindingID: 1001, ProjectID: 1, Tag: "tag", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Update",
			input:       &finding.TagFindingRequest{ProjectId: 1, Tag: &finding.FindingTagForUpsert{FindingId: 1001, ProjectId: 1, Tag: "tag"}},
			want:        &finding.TagFindingResponse{Tag: &finding.FindingTag{FindingTagId: 10011, ProjectId: 1, FindingId: 1001, Tag: "tag", CreatedAt: now.Add(-1 * 24 * time.Hour).Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.FindingTag{FindingTagID: 10011, FindingID: 1001, ProjectID: 1, Tag: "tag", CreatedAt: now.Add(-1 * 24 * time.Hour), UpdatedAt: now.Add(-1 * 24 * time.Hour)},
			mockUpResp:  &model.FindingTag{FindingTagID: 10011, FindingID: 1001, ProjectID: 1, Tag: "tag", CreatedAt: now.Add(-1 * 24 * time.Hour), UpdatedAt: now},
		},
		{
			name:    "NG Invalid request",
			input:   &finding.TagFindingRequest{ProjectId: 1, Tag: &finding.FindingTagForUpsert{FindingId: 1001, Tag: ""}},
			wantErr: true,
		},
		{
			name:       "NG GetFindingTagByKey error",
			input:      &finding.TagFindingRequest{ProjectId: 1, Tag: &finding.FindingTagForUpsert{FindingId: 1001, Tag: "tag"}},
			wantErr:    true,
			mockGetErr: gorm.ErrInvalidDB,
		},
		{
			name:        "NG TagFinding error",
			input:       &finding.TagFindingRequest{ProjectId: 1, Tag: &finding.FindingTagForUpsert{FindingId: 1001, Tag: "tag"}},
			wantErr:     true,
			mockGetResp: &model.FindingTag{FindingTagID: 10011, FindingID: 1001, Tag: "tag", CreatedAt: now, UpdatedAt: now},
			mockUpErr:   gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mockDB.On("GetFindingTagByKey").Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("TagFinding").Return(c.mockUpResp, c.mockUpErr).Once()
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
