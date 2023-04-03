package finding

import (
	"context"
	"errors"
	"reflect"
	"testing"

	aimocks "github.com/ca-risken/core/pkg/ai/mocks"
	dbmocks "github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/finding"
)

func TestAskAISummary(t *testing.T) {
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
		input            *finding.AskAISummaryRequest
		want             *finding.AskAISummaryResponse
		wantErr          bool
		mockGetFinding   *MockGetFinding
		mockGetRecommend *MockGetRecommend
		mockAskAI        *MockAskAI
	}{
		{
			name:  "OK",
			input: &finding.AskAISummaryRequest{ProjectId: 1, FindingId: 1, Lang: "en"},
			want:  &finding.AskAISummaryResponse{Answer: "answer"},
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
			name:    "NG Invalid param",
			input:   &finding.AskAISummaryRequest{FindingId: 1, Lang: "en"},
			wantErr: true,
		},
		{
			name:    "NG DB error(GetFinding)",
			input:   &finding.AskAISummaryRequest{ProjectId: 1, FindingId: 1, Lang: "en"},
			wantErr: true,
			mockGetFinding: &MockGetFinding{
				Err: errors.New("some error"),
			},
		},
		{
			name:    "NG DB error(GetRecommend)",
			input:   &finding.AskAISummaryRequest{ProjectId: 1, FindingId: 1, Lang: "en"},
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
			input:   &finding.AskAISummaryRequest{ProjectId: 1, FindingId: 1, Lang: "en"},
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
			svc := FindingService{repository: mockDB, ai: mockAI}

			if c.mockGetFinding != nil {
				mockDB.
					On("GetFinding", test.RepeatMockAnything(4)...).
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
			got, err := svc.AskAISummary(context.TODO(), c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected response: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
