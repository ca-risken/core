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
			mockError:    nil,
		},
		{
			name:         "NG Record not found",
			input:        &finding.ListFindingRequest{ProjectId: []uint32{123}, DataSource: []string{"aws:guardduty"}, ResourceName: []string{"hoge"}, FromScore: 0.0, ToScore: 1.0},
			want:         &finding.ListFindingResponse{},
			mockResponce: nil,
			mockError:    gorm.ErrRecordNotFound,
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
			mockError:    nil,
		},
		{
			name:         "NG record not found",
			input:        &finding.GetFindingRequest{FindingId: 9999},
			want:         &finding.GetFindingResponse{},
			mockResponce: nil,
			mockError:    gorm.ErrRecordNotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.On("GetFinding").Return(c.mockResponce, c.mockError).Once()
			result, err := svc.GetFinding(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
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
			name:        "OK Insert",
			input:       &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, OriginalMaxScore: 100.00}},
			want:        &finding.PutFindingResponse{Finding: &finding.Finding{FindingId: 1001, DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, Score: 1.0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: nil,
			mockGetErr:  gorm.ErrRecordNotFound,
			mockUpResp:  &model.Finding{FindingID: 1001, DataSource: "ds", DataSourceID: "ds-001", ResourceName: "rn", OriginalScore: 100.00, Score: 1.0, CreatedAt: now, UpdatedAt: now},
			mockUpErr:   nil,
		},
		{
			name:        "OK Update",
			input:       &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{FindingId: 1001, DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 20.00, OriginalMaxScore: 100.00}},
			want:        &finding.PutFindingResponse{Finding: &finding.Finding{FindingId: 1001, DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 20.00, Score: 0.2, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.Finding{FindingID: 1001, DataSource: "ds", DataSourceID: "ds-001", ResourceName: "rn", OriginalScore: 10.00, Score: 0.1, CreatedAt: now, UpdatedAt: now},
			mockGetErr:  nil,
			mockUpResp:  &model.Finding{FindingID: 1001, DataSource: "ds", DataSourceID: "ds-001", ResourceName: "rn", OriginalScore: 20.00, Score: 0.2, CreatedAt: now, UpdatedAt: now},
			mockUpErr:   nil,
		},
		{
			name:    "NG Invalid request(no data_source)",
			input:   &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{DataSource: "", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, OriginalMaxScore: 100.00}},
			want:    nil,
			wantErr: true,
		},
		{
			name:        "NG GetFindingByDataSource error",
			input:       &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, OriginalMaxScore: 100.00}},
			want:        nil,
			wantErr:     true,
			mockGetResp: nil,
			mockGetErr:  gorm.ErrInvalidSQL,
		},
		{
			name:        "NG Invalid finding_id",
			input:       &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{FindingId: 9999, DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, OriginalMaxScore: 100.00}},
			want:        nil,
			wantErr:     true,
			mockGetResp: &model.Finding{FindingID: 1001, DataSource: "ds", DataSourceID: "ds-001", ResourceName: "rn", OriginalScore: 10.00, Score: 0.1, CreatedAt: now, UpdatedAt: now},
			mockGetErr:  nil,
		},
		{
			name:        "NG UpsertFinding error",
			input:       &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", OriginalScore: 100.00, OriginalMaxScore: 100.00}},
			want:        nil,
			wantErr:     true,
			mockGetResp: &model.Finding{FindingID: 1001, DataSource: "ds", DataSourceID: "ds-001", ResourceName: "rn", OriginalScore: 10.00, Score: 0.1, CreatedAt: now, UpdatedAt: now},
			mockGetErr:  nil,
			mockUpResp:  nil,
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
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteFinding(t *testing.T) {
}

func TestListFindingTag(t *testing.T) {
}

func TestTagFinding(t *testing.T) {
}

func TestUntagFinding(t *testing.T) {
}

func TestListResource(t *testing.T) {
}

func TestGetResource(t *testing.T) {
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
			name:        "OK Insert",
			input:       &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceName: "rn", ProjectId: 111}},
			want:        &finding.PutResourceResponse{Resource: &finding.Resource{ResourceId: 1001, ResourceName: "rn", ProjectId: 111, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: nil,
			mockGetErr:  gorm.ErrRecordNotFound,
			mockUpResp:  &model.Resource{ResourceID: 1001, ResourceName: "rn", ProjectID: 111, CreatedAt: now, UpdatedAt: now},
			mockUpErr:   nil,
		},
		{
			name:        "OK Update",
			input:       &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceId: 1001, ResourceName: "rn-2", ProjectId: 999}},
			want:        &finding.PutResourceResponse{Resource: &finding.Resource{ResourceId: 1001, ResourceName: "rn-2", ProjectId: 999, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.Resource{ResourceID: 1001, ResourceName: "rn-2", ProjectID: 111, CreatedAt: now, UpdatedAt: now},
			mockGetErr:  nil,
			mockUpResp:  &model.Resource{ResourceID: 1001, ResourceName: "rn-2", ProjectID: 999, CreatedAt: now, UpdatedAt: now},
			mockUpErr:   nil,
		},
		{
			name:    "NG Invalid request",
			input:   &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceId: 1001, ResourceName: "", ProjectId: 111}},
			want:    nil,
			wantErr: true,
		},
		{
			name:        "NG GetResourceByName error",
			input:       &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceId: 1001, ResourceName: "", ProjectId: 111}},
			want:        nil,
			wantErr:     true,
			mockGetResp: nil,
			mockGetErr:  gorm.ErrCantStartTransaction,
		},
		{
			name:        "NG Invalid resource_id error",
			input:       &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceId: 1001, ResourceName: "rn", ProjectId: 111}},
			want:        nil,
			wantErr:     true,
			mockGetResp: &model.Resource{ResourceID: 9999, ResourceName: "rn", ProjectID: 111, CreatedAt: now, UpdatedAt: now},
			mockGetErr:  nil,
		},
		{
			name:        "NG UpsertResource error",
			input:       &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceId: 1001, ResourceName: "rn", ProjectId: 111}},
			want:        nil,
			wantErr:     true,
			mockGetResp: &model.Resource{ResourceID: 1001, ResourceName: "rn", ProjectID: 111, CreatedAt: now, UpdatedAt: now},
			mockGetErr:  nil,
			mockUpResp:  nil,
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

func (m *mockFindingRepository) UpsertResource(*model.Resource) (*model.Resource, error) {
	args := m.Called()
	return args.Get(0).(*model.Resource), args.Error(1)
}

func (m *mockFindingRepository) GetResourceByName(uint32, string) (*model.Resource, error) {
	args := m.Called()
	return args.Get(0).(*model.Resource), args.Error(1)
}
