package finding

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	aimocks "github.com/ca-risken/core/pkg/ai/mocks"
	"github.com/ca-risken/core/pkg/alertsummary"
	dbmocks "github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/finding"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func TestAskAISummary(t *testing.T) {
	savedAlertSummary := "saved alert summary"
	type MockGetFinding struct {
		Resp *model.Finding
		Err  error
	}
	type MockGetRecommend struct {
		Resp *model.Recommend
		Err  error
	}
	type MockAskAI struct {
		Resp string
		Err  error
	}
	cases := []struct {
		name             string
		input            *finding.GetAISummaryRequest
		want             *finding.GetAISummaryResponse
		wantErr          bool
		mockGetFinding   *MockGetFinding
		mockGetRecommend *MockGetRecommend
		mockAskAI        *MockAskAI
	}{
		{
			name:  "OK",
			input: &finding.GetAISummaryRequest{ProjectId: 1, FindingId: 1, Lang: "en"},
			want:  &finding.GetAISummaryResponse{Answer: "answer"},
			mockGetFinding: &MockGetFinding{
				Resp: &model.Finding{},
			},
			mockGetRecommend: &MockGetRecommend{
				Resp: &model.Recommend{},
			},
			mockAskAI: &MockAskAI{
				Resp: "answer",
			},
		},
		{
			name:  "OK ignore persisted alert summary",
			input: &finding.GetAISummaryRequest{ProjectId: 1, FindingId: 1, Lang: "en"},
			want:  &finding.GetAISummaryResponse{Answer: "fresh answer"},
			mockGetFinding: &MockGetFinding{
				Resp: &model.Finding{AISummary: &savedAlertSummary},
			},
			mockGetRecommend: &MockGetRecommend{
				Resp: &model.Recommend{},
			},
			mockAskAI: &MockAskAI{
				Resp: "fresh answer",
			},
		},
		{
			name:    "NG Invalid param",
			input:   &finding.GetAISummaryRequest{FindingId: 1, Lang: "en"},
			wantErr: true,
		},
		{
			name:    "NG DB error(GetFinding)",
			input:   &finding.GetAISummaryRequest{ProjectId: 1, FindingId: 1, Lang: "en"},
			wantErr: true,
			mockGetFinding: &MockGetFinding{
				Err: errors.New("some error"),
			},
		},
		{
			name:    "NG DB error(GetRecommend)",
			input:   &finding.GetAISummaryRequest{ProjectId: 1, FindingId: 1, Lang: "en"},
			wantErr: true,
			mockGetFinding: &MockGetFinding{
				Resp: &model.Finding{},
			},
			mockGetRecommend: &MockGetRecommend{
				Err: errors.New("some error"),
			},
		},
		{
			name:    "NG OpenAI API error",
			input:   &finding.GetAISummaryRequest{ProjectId: 1, FindingId: 1, Lang: "en"},
			wantErr: true,
			mockGetFinding: &MockGetFinding{
				Resp: &model.Finding{},
			},
			mockGetRecommend: &MockGetRecommend{
				Resp: &model.Recommend{},
			},
			mockAskAI: &MockAskAI{
				Err: errors.New("some error"),
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := dbmocks.NewFindingRepository(t)
			mockAI := aimocks.NewAIService(t)
			svc := FindingService{repository: mockDB, ai: mockAI, logger: logging.NewLogger()}

			if c.mockGetFinding != nil {
				mockDB.
					On("GetFinding", mock.Anything, c.input.ProjectId, c.input.FindingId, false).
					Return(c.mockGetFinding.Resp, c.mockGetFinding.Err).
					Once()
			}
			if c.mockGetRecommend != nil {
				mockDB.
					On("GetRecommend", test.RepeatMockAnything(3)...).
					Return(c.mockGetRecommend.Resp, c.mockGetRecommend.Err).
					Once()
			}
			if c.mockAskAI != nil {
				mockAI.
					On("AskAISummaryFromFinding", test.RepeatMockAnything(4)...).
					Return(c.mockAskAI.Resp, c.mockAskAI.Err).
					Once()
			}
			got, err := svc.GetAISummary(context.TODO(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetAlertAISummary(t *testing.T) {
	savedSummary := `{"blocks":[{"type":"text","text":"saved"}]}`
	type MockGetFinding struct {
		Resp *model.Finding
		Err  error
	}
	type MockGetRecommend struct {
		Resp *model.Recommend
		Err  error
	}
	type MockAskAI struct {
		Resp string
		Err  error
	}
	cases := []struct {
		name             string
		input            *finding.GetAlertAISummaryRequest
		want             *finding.GetAlertAISummaryResponse
		wantErr          bool
		wantErrCode      codes.Code
		mockGetFinding   *MockGetFinding
		mockGetRecommend *MockGetRecommend
		mockAskAI        *MockAskAI
		mockUpdateErr    error
	}{
		{
			name:  "OK generate and save",
			input: &finding.GetAlertAISummaryRequest{ProjectId: 1, FindingId: 1, Lang: "ja"},
			want:  &finding.GetAlertAISummaryResponse{AiSummary: `{"blocks":[{"type":"text","text":"summary"}]}`},
			mockGetFinding: &MockGetFinding{
				Resp: &model.Finding{},
			},
			mockGetRecommend: &MockGetRecommend{
				Resp: &model.Recommend{},
			},
			mockAskAI: &MockAskAI{
				Resp: `{"blocks":[{"type":"text","text":"summary"}]}`,
			},
		},
		{
			name:    "OK return saved summary",
			input:   &finding.GetAlertAISummaryRequest{ProjectId: 1, FindingId: 1, Lang: "ja"},
			want:    &finding.GetAlertAISummaryResponse{AiSummary: savedSummary},
			wantErr: false,
			mockGetFinding: &MockGetFinding{
				Resp: &model.Finding{AISummary: &savedSummary},
			},
		},
		{
			name:    "OK save failure still returns summary",
			input:   &finding.GetAlertAISummaryRequest{ProjectId: 1, FindingId: 1, Lang: "ja"},
			want:    &finding.GetAlertAISummaryResponse{AiSummary: `{"blocks":[{"type":"text","text":"summary"}]}`},
			wantErr: false,
			mockGetFinding: &MockGetFinding{
				Resp: &model.Finding{},
			},
			mockGetRecommend: &MockGetRecommend{
				Resp: &model.Recommend{},
			},
			mockAskAI: &MockAskAI{
				Resp: `{"blocks":[{"type":"text","text":"summary"}]}`,
			},
			mockUpdateErr: errors.New("save error"),
		},
		{
			name:    "NG invalid saved payload regenerates and fails when ai returns invalid payload",
			input:   &finding.GetAlertAISummaryRequest{ProjectId: 1, FindingId: 1, Lang: "ja"},
			wantErr: true,
			mockGetFinding: &MockGetFinding{
				Resp: &model.Finding{AISummary: ptr("saved")},
			},
			mockGetRecommend: &MockGetRecommend{
				Resp: &model.Recommend{},
			},
			mockAskAI: &MockAskAI{
				Resp: "summary",
			},
		},
		{
			name:    "NG invalid param",
			input:   &finding.GetAlertAISummaryRequest{FindingId: 1, Lang: "ja"},
			wantErr: true,
		},
		{
			name:    "NG DB error(GetFinding)",
			input:   &finding.GetAlertAISummaryRequest{ProjectId: 1, FindingId: 1, Lang: "ja"},
			wantErr: true,
			mockGetFinding: &MockGetFinding{
				Err: errors.New("some error"),
			},
		},
		{
			name:        "NG finding not found",
			input:       &finding.GetAlertAISummaryRequest{ProjectId: 1, FindingId: 1, Lang: "ja"},
			wantErr:     true,
			wantErrCode: codes.NotFound,
			mockGetFinding: &MockGetFinding{
				Err: gorm.ErrRecordNotFound,
			},
		},
		{
			name:    "NG DB error(GetRecommend)",
			input:   &finding.GetAlertAISummaryRequest{ProjectId: 1, FindingId: 1, Lang: "ja"},
			wantErr: true,
			mockGetFinding: &MockGetFinding{
				Resp: &model.Finding{},
			},
			mockGetRecommend: &MockGetRecommend{
				Err: errors.New("some error"),
			},
		},
		{
			name:    "NG OpenAI API error",
			input:   &finding.GetAlertAISummaryRequest{ProjectId: 1, FindingId: 1, Lang: "ja"},
			wantErr: true,
			mockGetFinding: &MockGetFinding{
				Resp: &model.Finding{},
			},
			mockGetRecommend: &MockGetRecommend{
				Resp: &model.Recommend{},
			},
			mockAskAI: &MockAskAI{
				Err: errors.New("some error"),
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := dbmocks.NewFindingRepository(t)
			mockAI := aimocks.NewAIService(t)
			svc := FindingService{repository: mockDB, ai: mockAI, logger: logging.NewLogger()}

			if c.mockGetFinding != nil {
				mockDB.
					On("GetFinding", mock.Anything, c.input.ProjectId, c.input.FindingId, true).
					Return(c.mockGetFinding.Resp, c.mockGetFinding.Err).
					Once()
			}
			if c.mockGetRecommend != nil {
				mockDB.
					On("GetRecommend", test.RepeatMockAnything(3)...).
					Return(c.mockGetRecommend.Resp, c.mockGetRecommend.Err).
					Once()
			}
			if c.mockAskAI != nil {
				mockAI.
					On("AskAlertAISummaryFromFinding", test.RepeatMockAnything(4)...).
					Return(c.mockAskAI.Resp, c.mockAskAI.Err).
					Once()
			}
			if c.mockAskAI != nil && c.mockAskAI.Err == nil && alertsummary.Normalize(c.mockAskAI.Resp) != "" {
				mockDB.
					On("UpdateFindingAISummary", test.RepeatMockAnything(5)...).
					Return(c.mockUpdateErr).
					Once()
			}
			got, err := svc.GetAlertAISummary(context.TODO(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatalf("Expected error, got nil")
			}
			if c.wantErrCode != 0 && status.Code(err) != c.wantErrCode {
				t.Fatalf("Unexpected gRPC code: want=%v, got=%v err=%v", c.wantErrCode, status.Code(err), err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}

func TestUpdateFindingAISummary(t *testing.T) {
	cases := []struct {
		name        string
		input       *finding.UpdateFindingAISummaryRequest
		wantErr     bool
		wantErrCode codes.Code
		mockErr     error
	}{
		{
			name:  "OK",
			input: &finding.UpdateFindingAISummaryRequest{ProjectId: 1, FindingId: 1, AiSummary: "summary", AiSummaryCreatedAt: 1735689600},
		},
		{
			name:    "NG invalid param",
			input:   &finding.UpdateFindingAISummaryRequest{ProjectId: 1, FindingId: 1, AiSummaryCreatedAt: 1735689600},
			wantErr: true,
		},
		{
			name:    "NG DB error",
			input:   &finding.UpdateFindingAISummaryRequest{ProjectId: 1, FindingId: 1, AiSummary: "summary", AiSummaryCreatedAt: 1735689600},
			wantErr: true,
			mockErr: errors.New("some error"),
		},
		{
			name:        "NG not found",
			input:       &finding.UpdateFindingAISummaryRequest{ProjectId: 1, FindingId: 1, AiSummary: "summary", AiSummaryCreatedAt: 1735689600},
			wantErr:     true,
			wantErrCode: codes.NotFound,
			mockErr:     gorm.ErrRecordNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := dbmocks.NewFindingRepository(t)
			svc := FindingService{repository: mockDB}
			if !c.wantErr || c.mockErr != nil {
				mockDB.
					On("UpdateFindingAISummary", test.RepeatMockAnything(5)...).
					Return(c.mockErr).
					Once()
			}
			_, err := svc.UpdateFindingAISummary(context.TODO(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatalf("Expected error, got nil")
			}
			if c.wantErrCode != 0 && status.Code(err) != c.wantErrCode {
				t.Fatalf("Unexpected gRPC code: want=%v, got=%v err=%v", c.wantErrCode, status.Code(err), err)
			}
		})
	}
}
