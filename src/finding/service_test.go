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

type mockFindingRepository struct {
	mock.Mock
}

func (m *mockFindingRepository) ListFinding(req *finding.ListFindingRequest) (*[]findingIds, error) {
	args := m.Called()
	return args.Get(0).(*[]findingIds), args.Error(1)
}

func (m *mockFindingRepository) GetFinding(req *finding.GetFindingRequest) (*model.Finding, error) {
	args := m.Called()
	return args.Get(0).(*model.Finding), args.Error(1)
}

func newFindingMockRepository() findingRepoInterface {
	return &mockFindingRepository{}
}

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
			name:         "normal",
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
			name:         "normal",
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
