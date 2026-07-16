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
	aipb "github.com/ca-risken/core/proto/ai"
)

func TestCreateRemediationProposal(t *testing.T) {
	now := time.Now()
	pendingModel := &model.RemediationProposal{
		RemediationProposalID: 1001,
		FindingID:             2001,
		ProjectID:             1,
		Status:                model.RemediationProposalStatusPending,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	cases := []struct {
		name           string
		input          *aipb.CreateRemediationProposalRequest
		want           *aipb.CreateRemediationProposalResponse
		mockCreateResp *model.RemediationProposal
		mockCreateErr  error
		wantErr        bool
	}{
		{
			name: "OK create pending proposal",
			input: &aipb.CreateRemediationProposalRequest{
				ProjectId: 1,
				FindingId: 2001,
			},
			mockCreateResp: pendingModel,
			want: &aipb.CreateRemediationProposalResponse{
				RemediationProposal: &aipb.RemediationProposal{
					RemediationProposalId: 1001,
					FindingId:             2001,
					ProjectId:             1,
					Status:                model.RemediationProposalStatusPending,
					CreatedAt:             now.Unix(),
					UpdatedAt:             now.Unix(),
				},
			},
		},
		{
			name: "NG validation error",
			input: &aipb.CreateRemediationProposalRequest{
				ProjectId: 1,
			},
			wantErr: true,
		},
		{
			name: "NG DB error(create)",
			input: &aipb.CreateRemediationProposalRequest{
				ProjectId: 1,
				FindingId: 2001,
			},
			mockCreateErr: errors.New("DB error"),
			wantErr:       true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAIRepository(t)
			svc := AIService{repository: mockDB}
			if c.mockCreateResp != nil || c.mockCreateErr != nil {
				mockDB.On("CreateRemediationProposal", test.RepeatMockAnything(2)...).Return(c.mockCreateResp, c.mockCreateErr).Once()
			}
			result, err := svc.CreateRemediationProposal(context.Background(), c.input)
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

func TestUpdateRemediationProposalStatus(t *testing.T) {
	now := time.Now()
	plan := `{"summary":"test"}`
	succeededModel := &model.RemediationProposal{
		RemediationProposalID: 1001,
		FindingID:             2001,
		ProjectID:             1,
		Status:                model.RemediationProposalStatusSucceeded,
		RemediationPlan:       &plan,
		GeneratedAt:           &now,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	cases := []struct {
		name           string
		input          *aipb.UpdateRemediationProposalStatusRequest
		want           *aipb.UpdateRemediationProposalStatusResponse
		mockUpdateResp *model.RemediationProposal
		mockUpdateErr  error
		wantErr        bool
	}{
		{
			name: "OK update succeeded proposal",
			input: &aipb.UpdateRemediationProposalStatusRequest{
				ProjectId:             1,
				RemediationProposalId: 1001,
				Status:                model.RemediationProposalStatusSucceeded,
				RemediationPlan:       plan,
				GeneratedAt:           now.Unix(),
			},
			mockUpdateResp: succeededModel,
			want: &aipb.UpdateRemediationProposalStatusResponse{
				RemediationProposal: &aipb.RemediationProposal{
					RemediationProposalId: 1001,
					FindingId:             2001,
					ProjectId:             1,
					Status:                model.RemediationProposalStatusSucceeded,
					RemediationPlan:       plan,
					GeneratedAt:           now.Unix(),
					CreatedAt:             now.Unix(),
					UpdatedAt:             now.Unix(),
				},
			},
		},
		{
			name: "NG validation error",
			input: &aipb.UpdateRemediationProposalStatusRequest{
				ProjectId:             1,
				RemediationProposalId: 1001,
				Status:                "unknown",
			},
			wantErr: true,
		},
		{
			name: "NG missing remediation_proposal_id",
			input: &aipb.UpdateRemediationProposalStatusRequest{
				ProjectId: 1,
				Status:    model.RemediationProposalStatusSucceeded,
			},
			wantErr: true,
		},
		{
			name: "NG DB error(update)",
			input: &aipb.UpdateRemediationProposalStatusRequest{
				ProjectId:             1,
				RemediationProposalId: 1001,
				Status:                model.RemediationProposalStatusSucceeded,
			},
			mockUpdateErr: errors.New("DB error"),
			wantErr:       true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockDB := mocks.NewAIRepository(t)
			svc := AIService{repository: mockDB}
			if c.mockUpdateResp != nil || c.mockUpdateErr != nil {
				mockDB.On("UpdateRemediationProposalStatus", test.RepeatMockAnything(7)...).Return(c.mockUpdateResp, c.mockUpdateErr).Once()
			}
			result, err := svc.UpdateRemediationProposalStatus(context.Background(), c.input)
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
