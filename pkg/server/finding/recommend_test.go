package finding

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	"gorm.io/gorm"
)

func TestGetRecommend(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mocks.MockFindingRepository{}
	svc := FindingService{repository: &mockDB}
	cases := []struct {
		name         string
		input        *finding.GetRecommendRequest
		want         *finding.GetRecommendResponse
		mockResponce *model.Recommend
		mockError    error
		wantErr      bool
	}{
		{
			name:         "OK",
			input:        &finding.GetRecommendRequest{ProjectId: 1, FindingId: 1},
			want:         &finding.GetRecommendResponse{Recommend: &finding.Recommend{FindingId: 1, RecommendId: 1, DataSource: "ds", Type: "a", Risk: "risk", Recommendation: "comment", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.Recommend{RecommendID: 1, DataSource: "ds", Type: "a", Risk: "risk", Recommendation: "comment", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK record not found",
			input:     &finding.GetRecommendRequest{ProjectId: 1, FindingId: 1},
			want:      &finding.GetRecommendResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid param",
			input:   &finding.GetRecommendRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &finding.GetRecommendRequest{ProjectId: 1, FindingId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetRecommend").Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.GetRecommend(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutRecommend(t *testing.T) {
	var ctx context.Context
	now := time.Now()
	mockDB := mocks.MockFindingRepository{}
	svc := FindingService{repository: &mockDB}
	cases := []struct {
		name    string
		input   *finding.PutRecommendRequest
		want    *finding.PutRecommendResponse
		wantErr bool

		// Mock responses
		mockGetFindingResp             *model.Finding
		mockGetFindingErr              error
		mockUpsertRecommendResp        *model.Recommend
		mockUpsertRecommendErr         error
		mockUpsertRecommendFindingResp *model.RecommendFinding
		mockUpsertRecommendFindingErr  error
	}{
		{
			name:                           "OK",
			input:                          &finding.PutRecommendRequest{ProjectId: 1, FindingId: 1, DataSource: "ds", Type: "a", Risk: "risk", Recommendation: "comment"},
			want:                           &finding.PutRecommendResponse{Recommend: &finding.Recommend{FindingId: 1, RecommendId: 1, DataSource: "ds", Type: "a", Risk: "risk", Recommendation: "comment", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockGetFindingResp:             &model.Finding{FindingID: 1},
			mockUpsertRecommendResp:        &model.Recommend{RecommendID: 1, DataSource: "ds", Type: "a", Risk: "risk", Recommendation: "comment", CreatedAt: now, UpdatedAt: now},
			mockUpsertRecommendFindingResp: &model.RecommendFinding{RecommendID: 1, FindingID: 1},
		},
		{
			name:    "NG Invalid request",
			input:   &finding.PutRecommendRequest{FindingId: 1, DataSource: "ds", Type: "a", Risk: "risk", Recommendation: "comment"},
			wantErr: true,
		},
		{
			name:              "Invalid DB error(GetFinding)",
			input:             &finding.PutRecommendRequest{ProjectId: 1, FindingId: 1, DataSource: "ds", Type: "a", Risk: "risk", Recommendation: "comment"},
			wantErr:           true,
			mockGetFindingErr: gorm.ErrInvalidDB,
		},
		{
			name:                   "Invalid DB error(UpsertRecommend)",
			input:                  &finding.PutRecommendRequest{ProjectId: 1, FindingId: 1, DataSource: "ds", Type: "a", Risk: "risk", Recommendation: "comment"},
			wantErr:                true,
			mockGetFindingResp:     &model.Finding{FindingID: 1},
			mockUpsertRecommendErr: gorm.ErrInvalidDB,
		},
		{
			name:                          "Invalid DB error(UpsertRecommendFinding)",
			input:                         &finding.PutRecommendRequest{ProjectId: 1, FindingId: 1, DataSource: "ds", Type: "a", Risk: "risk", Recommendation: "comment"},
			wantErr:                       true,
			mockGetFindingResp:            &model.Finding{FindingID: 1},
			mockUpsertRecommendResp:       &model.Recommend{RecommendID: 1, DataSource: "ds", Type: "a", Risk: "risk", Recommendation: "comment"},
			mockUpsertRecommendFindingErr: gorm.ErrInvalidDB,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mockGetFindingResp != nil || c.mockGetFindingErr != nil {
				mockDB.On("GetFinding").Return(c.mockGetFindingResp, c.mockGetFindingErr).Once()
			}
			if c.mockUpsertRecommendResp != nil || c.mockUpsertRecommendErr != nil {
				mockDB.On("UpsertRecommend").Return(c.mockUpsertRecommendResp, c.mockUpsertRecommendErr).Once()
			}
			if c.mockUpsertRecommendFindingResp != nil || c.mockUpsertRecommendFindingErr != nil {
				mockDB.On("UpsertRecommendFinding").Return(c.mockUpsertRecommendFindingResp, c.mockUpsertRecommendFindingErr).Once()
			}
			got, err := svc.PutRecommend(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
