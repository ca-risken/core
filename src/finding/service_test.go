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
			name:         "ok",
			input:        &finding.ListFindingRequest{ProjectId: []uint32{123}, DataSource: []string{"aws:guardduty"}, ResourceName: []string{"hoge"}, FromScore: 0.0, ToScore: 1.0},
			want:         &finding.ListFindingResponse{FindingId: []uint64{111, 222}},
			mockResponce: &[]findingIds{{FindingID: 111}, {FindingID: 222}},
			mockError:    nil,
		},
		{
			name:         "record not found",
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
			name:         "ok",
			input:        &finding.GetFindingRequest{FindingId: 1001},
			want:         &finding.GetFindingResponse{Finding: &finding.Finding{FindingId: 1001, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.Finding{FindingID: 1001, CreatedAt: now, UpdatedAt: now},
			mockError:    nil,
		},
		{
			name:         "record not found",
			input:        &finding.GetFindingRequest{},
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
		mockGetResp *model.Finding
		mockGetErr  error
		mockInsResp *model.Finding
		mockInsErr  error
		mockUpResp  *model.Finding
		mockUpErr   error
	}{
		{
			name:        "ok:insert",
			input:       &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{DataSource: "aws:hoge", ResourceName: "resource", OriginalScore: 100.00, OriginalMaxScore: 100.00}},
			want:        &finding.PutFindingResponse{Finding: &finding.Finding{FindingId: 1001, DataSource: "aws:hoge", ResourceName: "resource", OriginalScore: 100.00, Score: 1.0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: nil,
			mockGetErr:  gorm.ErrRecordNotFound,
			mockInsResp: &model.Finding{FindingID: 1001, DataSource: "aws:hoge", ResourceName: "resource", OriginalScore: 100.00, Score: 1.0, CreatedAt: now, UpdatedAt: now},
			mockInsErr:  nil,
			mockUpResp:  nil,
			mockUpErr:   nil,
		},
		{
			name:        "ok:update",
			input:       &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{FindingId: 1001, DataSource: "aws:hoge-2", ResourceName: "resource-2", OriginalScore: 20.00, OriginalMaxScore: 100.00}},
			want:        &finding.PutFindingResponse{Finding: &finding.Finding{FindingId: 1001, DataSource: "aws:hoge-2", ResourceName: "resource-2", OriginalScore: 20.00, Score: 0.2, CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetResp: &model.Finding{FindingID: 1001, DataSource: "aws:hoge-1", ResourceName: "resource-1", OriginalScore: 10.00, Score: 0.1, CreatedAt: now, UpdatedAt: now},
			mockGetErr:  nil,
			mockInsResp: nil,
			mockInsErr:  nil,
			mockUpResp:  &model.Finding{FindingID: 1001, DataSource: "aws:hoge-2", ResourceName: "resource-2", OriginalScore: 20.00, Score: 0.2, CreatedAt: now, UpdatedAt: now},
			mockUpErr:   nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGetResp != nil || c.mockGetErr != nil {
				mock.On("GetFinding").Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockInsResp != nil || c.mockInsErr != nil {
				mock.On("InsertFinding").Return(c.mockInsResp, c.mockInsErr).Once()
			}
			if c.mockUpResp != nil || c.mockUpErr != nil {
				mock.On("UpdateFinding").Return(c.mockUpResp, c.mockUpErr).Once()
			}

			got, err := svc.PutFinding(ctx, c.input)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestConvertFinding(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		input *model.Finding
		want  *finding.Finding
	}{
		{
			name:  "convert unix time",
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
			name:  "ok:0.01",
			input: [2]float32{1.0, 100.0},
			want:  0.01,
		},
		{
			name:  "ok:100",
			input: [2]float32{100.0, 100.0},
			want:  1.00,
		},
		{
			name:  "ok:0",
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

type mockFindingRepository struct {
	mock.Mock
}

func (m *mockFindingRepository) ListFinding(req *finding.ListFindingRequest) (*[]findingIds, error) {
	args := m.Called()
	return args.Get(0).(*[]findingIds), args.Error(1)
}

func (m *mockFindingRepository) GetFinding(findingID uint64) (*model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*model.Finding), args.Error(1)
}

func (m *mockFindingRepository) InsertFinding(*model.Finding) (*model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*model.Finding), args.Error(1)
}

func (m *mockFindingRepository) UpdateFinding(*model.Finding) (*model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*model.Finding), args.Error(1)
}
