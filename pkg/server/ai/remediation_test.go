package ai

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/pkg/test"
	"github.com/ca-risken/core/proto/ai"
	"gorm.io/gorm"
)

func TestGetRemediationProposal(t *testing.T) {
	now := time.Now()
	plan := `{"summary": "test"}`
	cases := []struct {
		name         string
		input        *ai.GetRemediationProposalRequest
		want         *ai.GetRemediationProposalResponse
		mockResponse *model.RemediationProposal
		mockError    error
		wantErr      bool
	}{
		{
			name:  "OK",
			input: &ai.GetRemediationProposalRequest{ProjectId: 1, RequestId: "req-0001"},
			want: &ai.GetRemediationProposalResponse{
				RemediationProposal: &ai.RemediationProposal{
					RequestId:       "req-0001",
					FindingId:       1001,
					ProjectId:       1,
					Status:          model.RemediationProposalStatusSucceeded,
					RemediationPlan: plan,
					GeneratedAt:     now.Unix(),
					CreatedAt:       now.Unix(),
					UpdatedAt:       now.Unix(),
				},
			},
			mockResponse: &model.RemediationProposal{
				RequestID:       "req-0001",
				FindingID:       1001,
				ProjectID:       1,
				Status:          model.RemediationProposalStatusSucceeded,
				RemediationPlan: &plan,
				GeneratedAt:     &now,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
		},
		{
			name:      "OK Record not found",
			input:     &ai.GetRemediationProposalRequest{ProjectId: 1, RequestId: "req-9999"},
			want:      &ai.GetRemediationProposalResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG Validation error",
			input:   &ai.GetRemediationProposalRequest{ProjectId: 0, RequestId: "req-0001"},
			want:    nil,
			wantErr: true,
		},
		{
			name:      "NG DB error",
			input:     &ai.GetRemediationProposalRequest{ProjectId: 1, RequestId: "req-0001"},
			want:      nil,
			mockError: errors.New("DB error"),
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewAIRepository(t)
			svc := AIService{repository: mockDB}

			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("GetRemediationProposal", test.RepeatMockAnything(3)...).Return(c.mockResponse, c.mockError).Once()
			}
			result, err := svc.GetRemediationProposal(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestListRemediationProposal(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *ai.ListRemediationProposalRequest
		want         *ai.ListRemediationProposalResponse
		mockResponse []*model.RemediationProposal
		mockError    error
		wantErr      bool
	}{
		{
			name:  "OK",
			input: &ai.ListRemediationProposalRequest{ProjectId: 1, FindingId: 1001},
			want: &ai.ListRemediationProposalResponse{
				RemediationProposal: []*ai.RemediationProposal{
					{
						RequestId: "req-0002",
						FindingId: 1001,
						ProjectId: 1,
						Status:    model.RemediationProposalStatusPending,
						CreatedAt: now.Unix(),
						UpdatedAt: now.Unix(),
					},
					{
						RequestId:    "req-0001",
						FindingId:    1001,
						ProjectId:    1,
						Status:       model.RemediationProposalStatusFailed,
						ErrorMessage: "some error",
						CreatedAt:    now.Add(-time.Hour).Unix(),
						UpdatedAt:    now.Add(-time.Hour).Unix(),
					},
				},
			},
			mockResponse: []*model.RemediationProposal{
				{
					RequestID: "req-0002",
					FindingID: 1001,
					ProjectID: 1,
					Status:    model.RemediationProposalStatusPending,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					RequestID:    "req-0001",
					FindingID:    1001,
					ProjectID:    1,
					Status:       model.RemediationProposalStatusFailed,
					ErrorMessage: func() *string { s := "some error"; return &s }(),
					CreatedAt:    now.Add(-time.Hour),
					UpdatedAt:    now.Add(-time.Hour),
				},
			},
		},
		{
			name:         "OK empty",
			input:        &ai.ListRemediationProposalRequest{ProjectId: 1, FindingId: 1001, Status: []string{model.RemediationProposalStatusPending}},
			want:         &ai.ListRemediationProposalResponse{},
			mockResponse: []*model.RemediationProposal{},
		},
		{
			name:    "NG Validation error(no finding_id)",
			input:   &ai.ListRemediationProposalRequest{ProjectId: 1},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "NG Validation error(invalid status)",
			input:   &ai.ListRemediationProposalRequest{ProjectId: 1, FindingId: 1001, Status: []string{"unknown"}},
			want:    nil,
			wantErr: true,
		},
		{
			name:      "NG DB error",
			input:     &ai.ListRemediationProposalRequest{ProjectId: 1, FindingId: 1001},
			want:      nil,
			mockError: errors.New("DB error"),
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewAIRepository(t)
			svc := AIService{repository: mockDB}

			if c.mockResponse != nil || c.mockError != nil {
				mockDB.On("ListRemediationProposal", test.RepeatMockAnything(4)...).Return(c.mockResponse, c.mockError).Once()
			}
			result, err := svc.ListRemediationProposal(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}

func TestPutRemediationProposal(t *testing.T) {
	now := time.Now()
	plan := `{"summary": "test"}`
	pendingUpsert := &ai.RemediationProposalForUpsert{
		RequestId: "req-0001",
		FindingId: 1001,
		ProjectId: 1,
		Status:    model.RemediationProposalStatusPending,
	}
	succeededUpsert := &ai.RemediationProposalForUpsert{
		RequestId:       "req-0001",
		FindingId:       1001,
		ProjectId:       1,
		Status:          model.RemediationProposalStatusSucceeded,
		RemediationPlan: plan,
		GeneratedAt:     now.Unix(),
	}
	pendingModel := &model.RemediationProposal{
		RequestID: "req-0001",
		FindingID: 1001,
		ProjectID: 1,
		Status:    model.RemediationProposalStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	succeededModel := &model.RemediationProposal{
		RequestID:       "req-0001",
		FindingID:       1001,
		ProjectID:       1,
		Status:          model.RemediationProposalStatusSucceeded,
		RemediationPlan: &plan,
		GeneratedAt:     &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	cases := []struct {
		name           string
		input          *ai.PutRemediationProposalRequest
		want           *ai.PutRemediationProposalResponse
		mockGetResp    *model.RemediationProposal
		mockGetErr     error
		mockCreateResp *model.RemediationProposal
		mockCreateErr  error
		mockUpdateResp *model.RemediationProposal
		mockUpdateErr  error
		wantErr        bool
	}{
		{
			name:           "OK create new proposal",
			input:          &ai.PutRemediationProposalRequest{ProjectId: 1, RemediationProposal: pendingUpsert},
			mockGetErr:     gorm.ErrRecordNotFound,
			mockCreateResp: pendingModel,
			want: &ai.PutRemediationProposalResponse{
				RemediationProposal: &ai.RemediationProposal{
					RequestId: "req-0001",
					FindingId: 1001,
					ProjectId: 1,
					Status:    model.RemediationProposalStatusPending,
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),
				},
			},
		},
		{
			name:           "OK update existing proposal",
			input:          &ai.PutRemediationProposalRequest{ProjectId: 1, RemediationProposal: succeededUpsert},
			mockGetResp:    pendingModel,
			mockUpdateResp: succeededModel,
			want: &ai.PutRemediationProposalResponse{
				RemediationProposal: &ai.RemediationProposal{
					RequestId:       "req-0001",
					FindingId:       1001,
					ProjectId:       1,
					Status:          model.RemediationProposalStatusSucceeded,
					RemediationPlan: plan,
					GeneratedAt:     now.Unix(),
					CreatedAt:       now.Unix(),
					UpdatedAt:       now.Unix(),
				},
			},
		},
		{
			name:    "NG Validation error(no remediation_proposal)",
			input:   &ai.PutRemediationProposalRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG Validation error(project_id mismatch)",
			input:   &ai.PutRemediationProposalRequest{ProjectId: 2, RemediationProposal: pendingUpsert},
			wantErr: true,
		},
		{
			name: "NG Validation error(invalid status)",
			input: &ai.PutRemediationProposalRequest{ProjectId: 1, RemediationProposal: &ai.RemediationProposalForUpsert{
				RequestId: "req-0001",
				FindingId: 1001,
				ProjectId: 1,
				Status:    "unknown",
			}},
			wantErr: true,
		},
		{
			name:       "NG DB error(get)",
			input:      &ai.PutRemediationProposalRequest{ProjectId: 1, RemediationProposal: pendingUpsert},
			mockGetErr: errors.New("DB error"),
			wantErr:    true,
		},
		{
			name:          "NG DB error(create)",
			input:         &ai.PutRemediationProposalRequest{ProjectId: 1, RemediationProposal: pendingUpsert},
			mockGetErr:    gorm.ErrRecordNotFound,
			mockCreateErr: errors.New("DB error"),
			wantErr:       true,
		},
		{
			name:          "NG DB error(update)",
			input:         &ai.PutRemediationProposalRequest{ProjectId: 1, RemediationProposal: succeededUpsert},
			mockGetResp:   pendingModel,
			mockUpdateErr: errors.New("DB error"),
			wantErr:       true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewAIRepository(t)
			svc := AIService{repository: mockDB}

			if c.mockGetResp != nil || c.mockGetErr != nil {
				mockDB.On("GetRemediationProposal", test.RepeatMockAnything(3)...).Return(c.mockGetResp, c.mockGetErr).Once()
			}
			if c.mockCreateResp != nil || c.mockCreateErr != nil {
				mockDB.On("CreateRemediationProposal", test.RepeatMockAnything(2)...).Return(c.mockCreateResp, c.mockCreateErr).Once()
			}
			if c.mockUpdateResp != nil || c.mockUpdateErr != nil {
				mockDB.On("UpdateRemediationProposalStatus", test.RepeatMockAnything(7)...).Return(c.mockUpdateResp, c.mockUpdateErr).Once()
			}
			result, err := svc.PutRemediationProposal(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("expected error but got nil")
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, result)
			}
		})
	}
}
