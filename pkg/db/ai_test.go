package db

import (
	"context"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ca-risken/core/pkg/model"
	"gorm.io/gorm"
)

func newRemediationProposalRows(data ...*model.RemediationProposal) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"request_id", "finding_id", "project_id", "status", "error_message", "remediation_plan", "generated_at", "created_at", "updated_at",
	})
	for _, d := range data {
		rows.AddRow(d.RequestID, d.FindingID, d.ProjectID, d.Status, d.ErrorMessage, d.RemediationPlan, d.GeneratedAt, d.CreatedAt, d.UpdatedAt)
	}
	return rows
}

func TestCreateRemediationProposal(t *testing.T) {
	now := time.Now()
	proposal := &model.RemediationProposal{
		RequestID: "req-0001",
		FindingID: 1001,
		ProjectID: 1,
		Status:    model.RemediationProposalStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	cases := []struct {
		name       string
		input      *model.RemediationProposal
		want       *model.RemediationProposal
		wantErr    bool
		mockErr    error
		mockGetErr error
	}{
		{
			name:    "OK",
			input:   proposal,
			want:    proposal,
			wantErr: false,
		},
		{
			name:    "NG DB error(insert)",
			input:   proposal,
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
		{
			name:       "NG DB error(select)",
			input:      proposal,
			wantErr:    true,
			mockGetErr: errors.New("DB error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("Failed to open mock sql db, error: %+v", err)
			}
			ctx := context.Background()
			if c.mockErr != nil {
				mock.ExpectExec(regexp.QuoteMeta(insertCreateRemediationProposal)).WillReturnError(c.mockErr)
			} else {
				mock.ExpectExec(regexp.QuoteMeta(insertCreateRemediationProposal)).
					WithArgs(c.input.RequestID, c.input.FindingID, c.input.ProjectID, c.input.Status, nil, nil, nil).
					WillReturnResult(sqlmock.NewResult(1, 1))
				if c.mockGetErr != nil {
					mock.ExpectQuery(regexp.QuoteMeta(selectGetRemediationProposal)).WillReturnError(c.mockGetErr)
				} else {
					mock.ExpectQuery(regexp.QuoteMeta(selectGetRemediationProposal)).
						WithArgs(c.input.ProjectID, c.input.RequestID).
						WillReturnRows(newRemediationProposalRows(c.want))
				}
			}

			got, err := client.CreateRemediationProposal(ctx, c.input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("No error")
			}
			if !c.wantErr && got.RequestID != c.want.RequestID {
				t.Fatalf("Unexpected result: got=%+v, want=%+v", got, c.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %+v", err)
			}
		})
	}
}

func TestGetRemediationProposal(t *testing.T) {
	now := time.Now()
	proposal := &model.RemediationProposal{
		RequestID: "req-0001",
		FindingID: 1001,
		ProjectID: 1,
		Status:    model.RemediationProposalStatusSucceeded,
		CreatedAt: now,
		UpdatedAt: now,
	}
	cases := []struct {
		name        string
		mockRows    *sqlmock.Rows
		mockErr     error
		wantErr     bool
		wantErrType error
	}{
		{
			name:     "OK",
			mockRows: newRemediationProposalRows(proposal),
			wantErr:  false,
		},
		{
			name:        "NG record not found",
			mockRows:    newRemediationProposalRows(),
			wantErr:     true,
			wantErrType: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG DB error",
			mockErr: errors.New("DB error"),
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("Failed to open mock sql db, error: %+v", err)
			}
			ctx := context.Background()
			if c.mockErr != nil {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetRemediationProposal)).WillReturnError(c.mockErr)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetRemediationProposal)).
					WithArgs(proposal.ProjectID, proposal.RequestID).
					WillReturnRows(c.mockRows)
			}

			got, err := client.GetRemediationProposal(ctx, proposal.ProjectID, proposal.RequestID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("No error")
			}
			if c.wantErrType != nil && !errors.Is(err, c.wantErrType) {
				t.Fatalf("Unexpected error type: got=%+v, want=%+v", err, c.wantErrType)
			}
			if !c.wantErr && got.RequestID != proposal.RequestID {
				t.Fatalf("Unexpected result: got=%+v, want=%+v", got, proposal)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %+v", err)
			}
		})
	}
}

func TestListRemediationProposal(t *testing.T) {
	now := time.Now()
	proposal1 := &model.RemediationProposal{RequestID: "req-0002", FindingID: 1001, ProjectID: 1, Status: model.RemediationProposalStatusPending, CreatedAt: now, UpdatedAt: now}
	proposal2 := &model.RemediationProposal{RequestID: "req-0001", FindingID: 1001, ProjectID: 1, Status: model.RemediationProposalStatusFailed, CreatedAt: now.Add(-time.Hour), UpdatedAt: now.Add(-time.Hour)}
	cases := []struct {
		name      string
		status    []string
		wantQuery string
		wantArgs  []driver.Value
		mockRows  *sqlmock.Rows
		mockErr   error
		wantCount int
		wantErr   bool
	}{
		{
			name:      "OK",
			wantQuery: `select * from remediation_proposal where project_id = ? and finding_id = ? order by created_at desc`,
			wantArgs:  []driver.Value{uint32(1), uint64(1001)},
			mockRows:  newRemediationProposalRows(proposal1, proposal2),
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "OK with status filter",
			status:    []string{model.RemediationProposalStatusPending, model.RemediationProposalStatusSucceeded},
			wantQuery: `select * from remediation_proposal where project_id = ? and finding_id = ? and status in (?,?) order by created_at desc`,
			wantArgs:  []driver.Value{uint32(1), uint64(1001), model.RemediationProposalStatusPending, model.RemediationProposalStatusSucceeded},
			mockRows:  newRemediationProposalRows(proposal1),
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "OK empty",
			wantQuery: `select * from remediation_proposal where project_id = ? and finding_id = ? order by created_at desc`,
			wantArgs:  []driver.Value{uint32(1), uint64(1001)},
			mockRows:  newRemediationProposalRows(),
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "NG DB error",
			wantQuery: `select * from remediation_proposal where project_id = ? and finding_id = ? order by created_at desc`,
			mockErr:   errors.New("DB error"),
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("Failed to open mock sql db, error: %+v", err)
			}
			ctx := context.Background()
			if c.mockErr != nil {
				mock.ExpectQuery(regexp.QuoteMeta(c.wantQuery)).WillReturnError(c.mockErr)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(c.wantQuery)).
					WithArgs(c.wantArgs...).
					WillReturnRows(c.mockRows)
			}

			got, err := client.ListRemediationProposal(ctx, 1, 1001, c.status)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("No error")
			}
			if !c.wantErr && len(got) != c.wantCount {
				t.Fatalf("Unexpected count: got=%d, want=%d", len(got), c.wantCount)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %+v", err)
			}
		})
	}
}

func TestGetLatestRemediationProposal(t *testing.T) {
	now := time.Now()
	proposal := &model.RemediationProposal{
		RequestID: "req-0002",
		FindingID: 1001,
		ProjectID: 1,
		Status:    model.RemediationProposalStatusSucceeded,
		CreatedAt: now,
		UpdatedAt: now,
	}
	cases := []struct {
		name        string
		mockRows    *sqlmock.Rows
		mockErr     error
		wantErr     bool
		wantErrType error
	}{
		{
			name:     "OK",
			mockRows: newRemediationProposalRows(proposal),
			wantErr:  false,
		},
		{
			name:        "NG record not found",
			mockRows:    newRemediationProposalRows(),
			wantErr:     true,
			wantErrType: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG DB error",
			mockErr: errors.New("DB error"),
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("Failed to open mock sql db, error: %+v", err)
			}
			ctx := context.Background()
			if c.mockErr != nil {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetLatestRemediationProposal)).WillReturnError(c.mockErr)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetLatestRemediationProposal)).
					WithArgs(proposal.ProjectID, proposal.FindingID).
					WillReturnRows(c.mockRows)
			}

			got, err := client.GetLatestRemediationProposal(ctx, proposal.ProjectID, proposal.FindingID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("No error")
			}
			if c.wantErrType != nil && !errors.Is(err, c.wantErrType) {
				t.Fatalf("Unexpected error type: got=%+v, want=%+v", err, c.wantErrType)
			}
			if !c.wantErr && got.RequestID != proposal.RequestID {
				t.Fatalf("Unexpected result: got=%+v, want=%+v", got, proposal)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %+v", err)
			}
		})
	}
}

func TestUpdateRemediationProposalStatus(t *testing.T) {
	now := time.Now()
	plan := `{"summary": "test"}`
	errMsg := "some error"
	succeeded := &model.RemediationProposal{
		RequestID:       "req-0001",
		FindingID:       1001,
		ProjectID:       1,
		Status:          model.RemediationProposalStatusSucceeded,
		RemediationPlan: &plan,
		GeneratedAt:     &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	failed := &model.RemediationProposal{
		RequestID:    "req-0001",
		FindingID:    1001,
		ProjectID:    1,
		Status:       model.RemediationProposalStatusFailed,
		ErrorMessage: &errMsg,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	cases := []struct {
		name    string
		want    *model.RemediationProposal
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK succeeded",
			want:    succeeded,
			wantErr: false,
		},
		{
			name:    "OK failed",
			want:    failed,
			wantErr: false,
		},
		{
			name:    "NG DB error",
			want:    succeeded,
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("Failed to open mock sql db, error: %+v", err)
			}
			ctx := context.Background()
			if c.mockErr != nil {
				mock.ExpectExec(regexp.QuoteMeta(updateUpdateRemediationProposalStatus)).WillReturnError(c.mockErr)
			} else {
				mock.ExpectExec(regexp.QuoteMeta(updateUpdateRemediationProposalStatus)).
					WithArgs(c.want.Status, c.want.ErrorMessage, c.want.RemediationPlan, c.want.GeneratedAt, c.want.ProjectID, c.want.RequestID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectGetRemediationProposal)).
					WithArgs(c.want.ProjectID, c.want.RequestID).
					WillReturnRows(newRemediationProposalRows(c.want))
			}

			got, err := client.UpdateRemediationProposalStatus(ctx, c.want.ProjectID, c.want.RequestID, c.want.Status, c.want.ErrorMessage, c.want.RemediationPlan, c.want.GeneratedAt)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("No error")
			}
			if !c.wantErr && got.Status != c.want.Status {
				t.Fatalf("Unexpected result: got=%+v, want=%+v", got, c.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %+v", err)
			}
		})
	}
}

func TestGetActiveRemediationProposal(t *testing.T) {
	now := time.Now()
	since := now.Add(-time.Hour)
	proposal := &model.RemediationProposal{
		RequestID: "req-0001",
		FindingID: 1001,
		ProjectID: 1,
		Status:    model.RemediationProposalStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	cases := []struct {
		name        string
		mockRows    *sqlmock.Rows
		mockErr     error
		wantErr     bool
		wantErrType error
	}{
		{
			name:     "OK",
			mockRows: newRemediationProposalRows(proposal),
			wantErr:  false,
		},
		{
			name:        "NG record not found",
			mockRows:    newRemediationProposalRows(),
			wantErr:     true,
			wantErrType: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG DB error",
			mockErr: errors.New("DB error"),
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatalf("Failed to open mock sql db, error: %+v", err)
			}
			ctx := context.Background()
			if c.mockErr != nil {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetActiveRemediationProposal)).WillReturnError(c.mockErr)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(selectGetActiveRemediationProposal)).
					WithArgs(proposal.ProjectID, proposal.FindingID, model.RemediationProposalStatusPending, model.RemediationProposalStatusSucceeded, since).
					WillReturnRows(c.mockRows)
			}

			got, err := client.GetActiveRemediationProposal(ctx, proposal.ProjectID, proposal.FindingID, since)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("No error")
			}
			if c.wantErrType != nil && !errors.Is(err, c.wantErrType) {
				t.Fatalf("Unexpected error type: got=%+v, want=%+v", err, c.wantErrType)
			}
			if !c.wantErr && got.RequestID != proposal.RequestID {
				t.Fatalf("Unexpected result: got=%+v, want=%+v", got, proposal)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %+v", err)
			}
		})
	}
}
