package main

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/jinzhu/gorm"
)

func TestListResource(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
	cases := []struct {
		name     string
		input    *finding.ListResourceRequest
		want     *finding.ListResourceResponse
		wantErr  bool
		mockResp *[]model.Resource
		mockErr  error
	}{
		{
			name:     "OK",
			input:    &finding.ListResourceRequest{ProjectId: 1, ResourceName: []string{"rn"}, FromSumScore: 0.0, ToSumScore: 100.0, FromAt: now.Unix(), ToAt: now.Unix()},
			want:     &finding.ListResourceResponse{ResourceId: []uint64{1001, 1002}},
			mockResp: &[]model.Resource{{ResourceID: 1001}, {ResourceID: 1002}},
		},
		{
			name:    "OK Not found",
			input:   &finding.ListResourceRequest{ProjectId: 1},
			want:    &finding.ListResourceResponse{},
			mockErr: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid request",
			input:   &finding.ListResourceRequest{ProjectId: 1, FromSumScore: -0.1},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &finding.ListResourceRequest{ProjectId: 1},
			wantErr: true,
			mockErr: gorm.ErrUnaddressable,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("ListResource").Return(c.mockResp, c.mockErr).Once()
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
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
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
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetResource").Return(c.mockResponce, c.mockError).Once()
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
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
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
			input:      &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceName: "", ProjectId: 111}},
			wantErr:    true,
			mockGetErr: gorm.ErrCantStartTransaction,
		},
		{
			name:        "NG UpsertResource error",
			input:       &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceName: "rn", ProjectId: 111}},
			wantErr:     true,
			mockGetResp: &model.Resource{ResourceID: 1001, ResourceName: "rn", ProjectID: 111, CreatedAt: now, UpdatedAt: now},
			mockUpErr:   gorm.ErrInvalidSQL,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mockDB.On("GetResourceByName").Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("UpsertResource").Return(c.mockUpResp, c.mockUpErr).Once()
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
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
	cases := []struct {
		name    string
		input   *finding.DeleteResourceRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &finding.DeleteResourceRequest{ProjectId: 1, ResourceId: 1001},
			wantErr: false,
		},
		{
			name:    "NG validation error",
			input:   &finding.DeleteResourceRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &finding.DeleteResourceRequest{ProjectId: 1, ResourceId: 1001},
			wantErr: true,
			mockErr: gorm.ErrCantStartTransaction,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("DeleteResource").Return(c.mockErr).Once()
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
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
	cases := []struct {
		name     string
		input    *finding.ListResourceTagRequest
		want     *finding.ListResourceTagResponse
		wantErr  bool
		mockResp *[]model.ResourceTag
		mockErr  error
	}{
		{
			name:  "OK",
			input: &finding.ListResourceTagRequest{ProjectId: 1, ResourceId: 1001},
			want: &finding.ListResourceTagResponse{Tag: []*finding.ResourceTag{
				{ResourceTagId: 1, ResourceId: 111, Tag: "tag1", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{ResourceTagId: 2, ResourceId: 111, Tag: "tag2", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResp: &[]model.ResourceTag{
				{ResourceTagID: 1, ResourceID: 111, Tag: "tag1", CreatedAt: now, UpdatedAt: now},
				{ResourceTagID: 2, ResourceID: 111, Tag: "tag2", CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:    "OK Record Not Found",
			input:   &finding.ListResourceTagRequest{ProjectId: 1, ResourceId: 1001},
			want:    &finding.ListResourceTagResponse{},
			wantErr: false,
			mockErr: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid Request",
			input:   &finding.ListResourceTagRequest{ProjectId: 1, ResourceId: 0},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &finding.ListResourceTagRequest{ProjectId: 1, ResourceId: 1001},
			wantErr: true,
			mockErr: gorm.ErrInvalidSQL,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("ListResourceTag").Return(c.mockResp, c.mockErr).Once()
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
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
	cases := []struct {
		name     string
		input    *finding.ListResourceTagNameRequest
		want     *finding.ListResourceTagNameResponse
		wantErr  bool
		mockResp *[]tagName
		mockErr  error
	}{
		{
			name:     "OK",
			input:    &finding.ListResourceTagNameRequest{ProjectId: 1, FromAt: 0, ToAt: now.Unix()},
			want:     &finding.ListResourceTagNameResponse{Tag: []string{"tag1", "tag2"}},
			mockResp: &[]tagName{{Tag: "tag1"}, {Tag: "tag2"}},
		},
		{
			name:    "OK Record Not Found",
			input:   &finding.ListResourceTagNameRequest{ProjectId: 1},
			want:    &finding.ListResourceTagNameResponse{},
			mockErr: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Invalid Request",
			input:   &finding.ListResourceTagNameRequest{},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &finding.ListResourceTagNameRequest{ProjectId: 1},
			wantErr: true,
			mockErr: gorm.ErrInvalidSQL,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResp != nil || c.mockErr != nil {
				mockDB.On("ListResourceTagName").Return(c.mockResp, c.mockErr).Once()
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
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
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
			input:      &finding.TagResourceRequest{ProjectId: 1, Tag: &finding.ResourceTagForUpsert{ResourceId: 1001, Tag: "tag"}},
			wantErr:    true,
			mockGetErr: gorm.ErrInvalidSQL,
		},
		{
			name:        "NG TagFinding error",
			input:       &finding.TagResourceRequest{ProjectId: 1, Tag: &finding.ResourceTagForUpsert{ResourceId: 1001, Tag: "tag"}},
			wantErr:     true,
			mockGetResp: &model.ResourceTag{ResourceTagID: 10011, ResourceID: 1001, Tag: "tag", CreatedAt: now, UpdatedAt: now},
			mockUpErr:   gorm.ErrInvalidSQL,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mockDB.On("GetResourceTagByKey").Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mockDB.On("TagResource").Return(c.mockUpResp, c.mockUpErr).Once()
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
	mockDB := mockFindingRepository{}
	svc := findingService{repository: &mockDB}
	cases := []struct {
		name    string
		input   *finding.UntagResourceRequest
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   &finding.UntagResourceRequest{ProjectId: 1, ResourceTagId: 1001},
			wantErr: false,
		},
		{
			name:    "NG validation error",
			input:   &finding.UntagResourceRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &finding.UntagResourceRequest{ProjectId: 1, ResourceTagId: 1001},
			wantErr: true,
			mockErr: gorm.ErrCantStartTransaction,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB.On("UntagResource").Return(c.mockErr).Once()
			_, err := svc.UntagResource(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
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
			input: &finding.ListResourceRequest{ProjectId: 1, ResourceName: []string{"rn"}, FromSumScore: 0.3, ToSumScore: 0.9, FromAt: now.Unix(), ToAt: now.Unix()},
			want:  &finding.ListResourceRequest{ProjectId: 1, ResourceName: []string{"rn"}, FromSumScore: 0.3, ToSumScore: 0.9, FromAt: now.Unix(), ToAt: now.Unix()},
		},
		{
			name:  "OK convert ToSumScore",
			input: &finding.ListResourceRequest{ProjectId: 1, ToAt: now.Unix()},
			want:  &finding.ListResourceRequest{ProjectId: 1, ToSumScore: maxSumScore, ToAt: now.Unix()},
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