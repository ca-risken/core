package finding

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/finding"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestGetRecommend(t *testing.T) {
	var ctx context.Context
	now := time.Now()
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
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}

			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetRecommend", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
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
			mockDB := mocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB, logger: logging.NewLogger()}

			if c.mockGetFindingResp != nil || c.mockGetFindingErr != nil {
				mockDB.On("GetFinding", test.RepeatMockAnything(4)...).Return(c.mockGetFindingResp, c.mockGetFindingErr).Once()
			}
			if c.mockUpsertRecommendResp != nil || c.mockUpsertRecommendErr != nil {
				mockDB.On("UpsertRecommend", test.RepeatMockAnything(2)...).Return(c.mockUpsertRecommendResp, c.mockUpsertRecommendErr).Once()
			}
			if c.mockUpsertRecommendFindingResp != nil || c.mockUpsertRecommendFindingErr != nil {
				mockDB.On("UpsertRecommendFinding", test.RepeatMockAnything(2)...).Return(c.mockUpsertRecommendFindingResp, c.mockUpsertRecommendFindingErr).Once()
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

func TestPutRecommendSCCMethod(t *testing.T) {
	cases := []struct{ name, data, want string }{
		{"OS", `{"finding":{"vulnerability":{"offending_package":{"package_type":"OS"}}}}`, "os_update"},
		{"application", `{"finding":{"files":[{"path":"app/lib/a"}],"vulnerability":{"offending_package":{"package_name":"a"}}}}`, "app_dependency_update"},
		{"binary", `{"finding":{"files":[{"disk_path":{"relative_path":"usr/bin/a"}}],"vulnerability":{"offending_package":{"package_name":"a"}}}}`, "library_update"},
		{"missing package", `{}`, "review"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := mocks.NewFindingRepository(t)
			repo.On("GetFinding", mock.Anything, uint32(1), uint64(2), true).Return(&model.Finding{FindingID: 2, DataSource: "google:scc", Data: c.data}, nil).Once()
			repo.On("UpsertRecommend", mock.Anything, mock.MatchedBy(func(r *model.Recommend) bool {
				return r.RemediationType == c.want && r.Recommendation == "original steps"
			})).Return(&model.Recommend{RecommendID: 3}, nil).Once()
			repo.On("UpsertRecommendFinding", mock.Anything, mock.Anything).Return(&model.RecommendFinding{FindingID: 2, RecommendID: 3, ProjectID: 1}, nil).Once()
			service := FindingService{repository: repo, logger: logging.NewLogger()}
			_, err := service.PutRecommend(context.Background(), &finding.PutRecommendRequest{ProjectId: 1, FindingId: 2, DataSource: "google:scc", Type: "unique-finding-name", Risk: "SOFTWARE_VULNERABILITY", Recommendation: "original steps"})
			require.NoError(t, err)
		})
	}
}
