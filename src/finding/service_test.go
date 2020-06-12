package main

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

func TestListFinding(t *testing.T) {
	var ctx context.Context
	mock := mockFindingRepository{}
	svc := newFindingService(&mock)

	cases := []struct {
		name         string
		input        *finding.ListFindingRequest
		want         *finding.ListFindingResponse
		mockResponce *[]findingIds
		mockError    error
	}{
		{
			name:         "OK",
			input:        &finding.ListFindingRequest{ProjectId: []uint32{123}, DataSource: []string{"aws:guardduty"}, ResourceName: []string{"hoge"}, FromScore: 0.0, ToScore: 1.0},
			want:         &finding.ListFindingResponse{FindingId: []uint64{111, 222}},
			mockResponce: &[]findingIds{{FindingID: 111}, {FindingID: 222}},
		},
		{
			name:      "NG Record not found",
			input:     &finding.ListFindingRequest{ProjectId: []uint32{123}, DataSource: []string{"aws:guardduty"}, ResourceName: []string{"hoge"}, FromScore: 0.0, ToScore: 1.0},
			want:      &finding.ListFindingResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.On("ListFinding").Return(c.mockResponce, c.mockError).Once()
			result, err := svc.ListFinding(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestConvertListFindingRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *finding.ListFindingRequest
		want  *finding.ListFindingRequest
	}{
		{
			name:  "OK full-set",
			input: &finding.ListFindingRequest{ProjectId: []uint32{111, 222}, DataSource: []string{"ds"}, ResourceName: []string{"rn"}, FromScore: 0.3, ToScore: 0.9, FromAt: now.Unix(), ToAt: now.Unix()},
			want:  &finding.ListFindingRequest{ProjectId: []uint32{111, 222}, DataSource: []string{"ds"}, ResourceName: []string{"rn"}, FromScore: 0.3, ToScore: 0.9, FromAt: now.Unix(), ToAt: now.Unix()},
		},
		{
			name:  "OK convert ToScore",
			input: &finding.ListFindingRequest{ToAt: now.Unix()},
			want:  &finding.ListFindingRequest{ToScore: 1.0, ToAt: now.Unix()},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertListFindingRequest(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected convert: got=%+v, want: %+v", got, c.want)
			}
		})
	}
}

func TestGetFinding(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mock := mockFindingRepository{}
	svc := newFindingService(&mock)

	cases := []struct {
		name         string
		input        *finding.GetFindingRequest
		want         *finding.GetFindingResponse
		mockResponce *model.Finding
		mockError    error
	}{
		{
			name:         "OK",
			input:        &finding.GetFindingRequest{FindingId: 1001},
			want:         &finding.GetFindingResponse{Finding: &finding.Finding{FindingId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.Finding{FindingID: 1001, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "NG record not found",
			input:     &finding.GetFindingRequest{FindingId: 9999},
			want:      &finding.GetFindingResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.On("GetFinding").Return(c.mockResponce, c.mockError).Once()
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
	mock := mockFindingRepository{}
	svc := newFindingService(&mock)

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
			input:       &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{FindingId: 1001, DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 20.00, OriginalMaxScore: 100.00}},
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
			mockGetErr: gorm.ErrInvalidSQL,
		},
		{
			name:        "NG Invalid finding_id",
			input:       &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{FindingId: 9999, DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, OriginalMaxScore: 100.00}},
			wantErr:     true,
			mockGetResp: &model.Finding{FindingID: 1001, DataSource: "ds", DataSourceID: "ds-001", ResourceName: "rn", OriginalScore: 10.00, Score: 0.1, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "NG UpsertFinding error",
			input:       &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, OriginalMaxScore: 100.00}},
			wantErr:     true,
			mockGetResp: &model.Finding{FindingID: 1001, DataSource: "ds", DataSourceID: "ds-001", ResourceName: "rn", OriginalScore: 10.00, Score: 0.1, CreatedAt: now, UpdatedAt: now},
			mockUpErr:   gorm.ErrInvalidSQL,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Resource関連のupdateは別テストで実施。ここでは一律カラを返す
			mock.On("GetResourceByName").Return(&model.Resource{}, nil)
			mock.On("UpsertResource").Return(&model.Resource{}, nil)
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mock.On("GetFindingByDataSource").Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mock.On("UpsertFinding").Return(c.mockUpResp, c.mockUpErr).Once()
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
	mock := mockFindingRepository{}
	svc := newFindingService(&mock)

	cases := []struct {
		name    string
		input   *finding.DeleteFindingRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &finding.DeleteFindingRequest{FindingId: 1001},
			wantErr: false,
		},
		{
			name:    "NG validation error",
			input:   &finding.DeleteFindingRequest{},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &finding.DeleteFindingRequest{FindingId: 1001},
			wantErr: true,
			mockErr: gorm.ErrCantStartTransaction,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.On("DeleteFinding").Return(c.mockErr).Once()
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
	mock := mockFindingRepository{}
	svc := newFindingService(&mock)
	cases := []struct {
		name     string
		input    *finding.ListFindingTagRequest
		want     *finding.ListFindingTagResponse
		wantErr  bool
		mockResp *[]model.FindingTag
		mockErr  error
	}{
		{
			name:  "OK",
			input: &finding.ListFindingTagRequest{FindingId: 1001},
			want: &finding.ListFindingTagResponse{Tag: []*finding.FindingTag{
				{FindingTagId: 1, FindingId: 111, TagKey: "k-1", TagValue: "v-1", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{FindingTagId: 2, FindingId: 111, TagKey: "k-2", TagValue: "v-2", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResp: &[]model.FindingTag{
				{FindingTagID: 1, FindingID: 111, TagKey: "k-1", TagValue: "v-1", CreatedAt: now, UpdatedAt: now},
				{FindingTagID: 2, FindingID: 111, TagKey: "k-2", TagValue: "v-2", CreatedAt: now, UpdatedAt: now},
			},
			mockErr: nil,
		},
		{
			name:    "OK Record Not Found",
			input:   &finding.ListFindingTagRequest{FindingId: 1001},
			want:    &finding.ListFindingTagResponse{},
			wantErr: false,
			mockErr: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid Request",
			input:   &finding.ListFindingTagRequest{FindingId: 0},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &finding.ListFindingTagRequest{FindingId: 0},
			wantErr: true,
			mockErr: gorm.ErrInvalidSQL,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.On("ListFindingTag").Return(c.mockResp, c.mockErr).Once()
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

func TestTagFinding(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mock := mockFindingRepository{}
	svc := newFindingService(&mock)
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
			input:      &finding.TagFindingRequest{Tag: &finding.FindingTagForUpsert{FindingId: 1001, TagKey: "k", TagValue: "v"}},
			want:       &finding.TagFindingResponse{Tag: &finding.FindingTag{FindingTagId: 10011, FindingId: 1001, TagKey: "k", TagValue: "v", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetErr: gorm.ErrRecordNotFound,
			mockUpResp: &model.FindingTag{FindingTagID: 10011, FindingID: 1001, TagKey: "k", TagValue: "v", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "OK Update",
			input:       &finding.TagFindingRequest{Tag: &finding.FindingTagForUpsert{FindingId: 1001, TagKey: "k", TagValue: "v"}},
			want:        &finding.TagFindingResponse{Tag: &finding.FindingTag{FindingTagId: 10011, FindingId: 1001, TagKey: "k", TagValue: "v", CreatedAt: now.Add(-1 * 24 * time.Hour).Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.FindingTag{FindingTagID: 10011, FindingID: 1001, TagKey: "k", TagValue: "v", CreatedAt: now.Add(-1 * 24 * time.Hour), UpdatedAt: now.Add(-1 * 24 * time.Hour)},
			mockUpResp:  &model.FindingTag{FindingTagID: 10011, FindingID: 1001, TagKey: "k", TagValue: "v", CreatedAt: now.Add(-1 * 24 * time.Hour), UpdatedAt: now},
		},
		{
			name:    "NG Invalid request",
			input:   &finding.TagFindingRequest{Tag: &finding.FindingTagForUpsert{FindingId: 1001, TagKey: "", TagValue: "v"}},
			wantErr: true,
		},
		{
			name:       "NG GetFindingTagByKey error",
			input:      &finding.TagFindingRequest{Tag: &finding.FindingTagForUpsert{FindingId: 1001, TagKey: "k", TagValue: "v"}},
			wantErr:    true,
			mockGetErr: gorm.ErrInvalidSQL,
		},
		{
			name:        "NG Invalid findingTagID",
			input:       &finding.TagFindingRequest{Tag: &finding.FindingTagForUpsert{FindingTagId: 10011, FindingId: 1001, TagKey: "k", TagValue: "v"}},
			wantErr:     true,
			mockGetResp: &model.FindingTag{FindingTagID: 99999, FindingID: 1001, TagKey: "k", TagValue: "v", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "NG TagFinding error",
			input:       &finding.TagFindingRequest{Tag: &finding.FindingTagForUpsert{FindingTagId: 10011, FindingId: 1001, TagKey: "k", TagValue: "v"}},
			wantErr:     true,
			mockGetResp: &model.FindingTag{FindingTagID: 10011, FindingID: 1001, TagKey: "k", TagValue: "v", CreatedAt: now, UpdatedAt: now},
			mockUpErr:   gorm.ErrInvalidSQL,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mock.On("GetFindingTagByKey").Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mock.On("TagFinding").Return(c.mockUpResp, c.mockUpErr).Once()
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

func TestUntagFinding(t *testing.T) {
	var ctx context.Context
	mock := mockFindingRepository{}
	svc := newFindingService(&mock)
	cases := []struct {
		name    string
		input   *finding.UntagFindingRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &finding.UntagFindingRequest{FindingTagId: 1001},
			wantErr: false,
		},
		{
			name:    "NG validation error",
			input:   &finding.UntagFindingRequest{},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &finding.UntagFindingRequest{FindingTagId: 1001},
			wantErr: true,
			mockErr: gorm.ErrCantStartTransaction,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.On("UntagFinding").Return(c.mockErr).Once()
			_, err := svc.UntagFinding(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
		})
	}
}

func TestListResource(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mock := mockFindingRepository{}
	svc := newFindingService(&mock)
	cases := []struct {
		name     string
		input    *finding.ListResourceRequest
		want     *finding.ListResourceResponse
		wantErr  bool
		mockResp *[]resourceIds
		mockErr  error
	}{
		{
			name:     "OK",
			input:    &finding.ListResourceRequest{ProjectId: []uint32{111}, ResourceName: []string{"rn"}, FromSumScore: 0.0, ToSumScore: 100.0, FromAt: now.Unix(), ToAt: now.Unix()},
			want:     &finding.ListResourceResponse{ResourceId: []uint64{1001, 1002}},
			mockResp: &[]resourceIds{{ResourceID: 1001}, {ResourceID: 1002}},
		},
		{
			name:    "OK Not found",
			input:   &finding.ListResourceRequest{},
			want:    &finding.ListResourceResponse{},
			mockErr: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid request",
			input:   &finding.ListResourceRequest{FromSumScore: -0.1},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &finding.ListResourceRequest{},
			wantErr: true,
			mockErr: gorm.ErrUnaddressable,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mock.On("ListResource").Return(c.mockResp, c.mockErr).Once()
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

func TestConvertListResourceRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *finding.ListResourceRequest
		want  *finding.ListResourceRequest
	}{
		{
			name:  "OK full-set",
			input: &finding.ListResourceRequest{ProjectId: []uint32{111, 222}, ResourceName: []string{"rn"}, FromSumScore: 0.3, ToSumScore: 0.9, FromAt: now.Unix(), ToAt: now.Unix()},
			want:  &finding.ListResourceRequest{ProjectId: []uint32{111, 222}, ResourceName: []string{"rn"}, FromSumScore: 0.3, ToSumScore: 0.9, FromAt: now.Unix(), ToAt: now.Unix()},
		},
		{
			name:  "OK convert ToSumScore",
			input: &finding.ListResourceRequest{ToAt: now.Unix()},
			want:  &finding.ListResourceRequest{ToSumScore: maxSumScore, ToAt: now.Unix()},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertListResourceRequest(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected convert: got=%+v, want: %+v", got, c.want)
			}
		})
	}
}

func TestGetResource(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mock := mockFindingRepository{}
	svc := newFindingService(&mock)
	cases := []struct {
		name         string
		input        *finding.GetResourceRequest
		want         *finding.GetResourceResponse
		mockResponce *model.Resource
		mockError    error
	}{
		{
			name:         "OK",
			input:        &finding.GetResourceRequest{ResourceId: 1001},
			want:         &finding.GetResourceResponse{Resource: &finding.Resource{ResourceId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.Resource{ResourceID: 1001, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "NG Record not found",
			input:     &finding.GetResourceRequest{ResourceId: 9999},
			want:      &finding.GetResourceResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.On("GetResource").Return(c.mockResponce, c.mockError).Once()
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
	mock := mockFindingRepository{}
	svc := newFindingService(&mock)

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
			input:       &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceId: 1001, ResourceName: "rn-2", ProjectId: 999}},
			want:        &finding.PutResourceResponse{Resource: &finding.Resource{ResourceId: 1001, ResourceName: "rn-2", ProjectId: 999, CreatedAt: now.Add(-1 * 24 * time.Hour).Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.Resource{ResourceID: 1001, ResourceName: "rn-2", ProjectID: 111, CreatedAt: now.Add(-1 * 24 * time.Hour), UpdatedAt: now.Add(-1 * 24 * time.Hour)},
			mockUpResp:  &model.Resource{ResourceID: 1001, ResourceName: "rn-2", ProjectID: 999, CreatedAt: now.Add(-1 * 24 * time.Hour), UpdatedAt: now},
		},
		{
			name:    "NG Invalid request",
			input:   &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceId: 1001, ResourceName: "", ProjectId: 111}},
			wantErr: true,
		},
		{
			name:       "NG GetResourceByName error",
			input:      &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceId: 1001, ResourceName: "", ProjectId: 111}},
			wantErr:    true,
			mockGetErr: gorm.ErrCantStartTransaction,
		},
		{
			name:        "NG Invalid resource_id error",
			input:       &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceId: 1001, ResourceName: "rn", ProjectId: 111}},
			wantErr:     true,
			mockGetResp: &model.Resource{ResourceID: 9999, ResourceName: "rn", ProjectID: 111, CreatedAt: now, UpdatedAt: now},
		},
		{
			name:        "NG UpsertResource error",
			input:       &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceId: 1001, ResourceName: "rn", ProjectId: 111}},
			wantErr:     true,
			mockGetResp: &model.Resource{ResourceID: 1001, ResourceName: "rn", ProjectID: 111, CreatedAt: now, UpdatedAt: now},
			mockUpErr:   gorm.ErrInvalidSQL,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mock.On("GetResourceByName").Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mock.On("UpsertResource").Return(c.mockUpResp, c.mockUpErr).Once()
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
}

func TestListResourceTag(t *testing.T) {
}

func TestTagResource(t *testing.T) {
}

func TestUntagResource(t *testing.T) {
}

func TestConvertFinding(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *model.Finding
		want  *finding.Finding
	}{
		{
			name:  "OK convert unix time",
			input: &model.Finding{FindingID: 10001, CreatedAt: now, UpdatedAt: now},
			want:  &finding.Finding{FindingId: 10001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
		},
		{
			name:  "OK empty",
			input: nil,
			want:  &finding.Finding{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertFinding(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestConvertFindingTag(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *model.FindingTag
		want  *finding.FindingTag
	}{
		{
			name:  "OK convert unix time",
			input: &model.FindingTag{FindingTagID: 11111, FindingID: 10001, TagKey: "k", TagValue: "v", CreatedAt: now, UpdatedAt: now},
			want:  &finding.FindingTag{FindingTagId: 11111, FindingId: 10001, TagKey: "k", TagValue: "v", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
		},
		{
			name:  "OK empty",
			input: nil,
			want:  &finding.FindingTag{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertFindingTag(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestCalculateScore(t *testing.T) {
	cases := []struct {
		name  string
		input [2]float32
		want  float32
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
			name:  "ok Score 0%",
			input: [2]float32{0, 100.0},
			want:  0.00,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := calculateScore(c.input[0], c.input[1])
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected result: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestConvertResource(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *model.Resource
		want  *finding.Resource
	}{
		{
			name:  "OK convert unix time",
			input: &model.Resource{ResourceID: 10001, ResourceName: "rn", ProjectID: 111, CreatedAt: now, UpdatedAt: now},
			want:  &finding.Resource{ResourceId: 10001, ResourceName: "rn", ProjectId: 111, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
		},
		{
			name:  "OK empty",
			input: nil,
			want:  &finding.Resource{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := convertResource(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

/*
 * Mock Repository
 */
type mockFindingRepository struct {
	mock.Mock
}

func (m *mockFindingRepository) ListFinding(*finding.ListFindingRequest) (*[]findingIds, error) {
	args := m.Called()
	return args.Get(0).(*[]findingIds), args.Error(1)
}

func (m *mockFindingRepository) GetFinding(uint64) (*model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*model.Finding), args.Error(1)
}

func (m *mockFindingRepository) UpsertFinding(*model.Finding) (*model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*model.Finding), args.Error(1)
}

func (m *mockFindingRepository) GetFindingByDataSource(uint32, string, string) (*model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*model.Finding), args.Error(1)
}

func (m *mockFindingRepository) DeleteFinding(uint64) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockFindingRepository) ListFindingTag(uint64) (*[]model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*[]model.FindingTag), args.Error(1)
}

func (m *mockFindingRepository) TagFinding(*model.FindingTag) (*model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingTag), args.Error(1)
}

func (m *mockFindingRepository) GetFindingTagByKey(uint64, string) (*model.FindingTag, error) {
	args := m.Called()
	return args.Get(0).(*model.FindingTag), args.Error(1)
}

func (m *mockFindingRepository) UntagFinding(uint64) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockFindingRepository) ListResource(*finding.ListResourceRequest) (*[]resourceIds, error) {
	args := m.Called()
	return args.Get(0).(*[]resourceIds), args.Error(1)
}

func (m *mockFindingRepository) UpsertResource(*model.Resource) (*model.Resource, error) {
	args := m.Called()
	return args.Get(0).(*model.Resource), args.Error(1)
}

func (m *mockFindingRepository) GetResourceByName(uint32, string) (*model.Resource, error) {
	args := m.Called()
	return args.Get(0).(*model.Resource), args.Error(1)
}

func (m *mockFindingRepository) GetResource(uint64) (*model.Resource, error) {
	args := m.Called()
	return args.Get(0).(*model.Resource), args.Error(1)
}
