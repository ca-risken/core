package db

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ca-risken/core/pkg/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func newRemediationProposalRows(data ...*model.RemediationProposal) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"remediation_proposal_id", "finding_id", "project_id", "status", "status_detail", "remediation_plan", "generated_at", "created_at", "updated_at",
	})
	for _, d := range data {
		rows.AddRow(d.RemediationProposalID, d.FindingID, d.ProjectID, d.Status, d.StatusDetail, d.RemediationPlan, d.GeneratedAt, d.CreatedAt, d.UpdatedAt)
	}
	return rows
}

func newSplitMockClient() (*Client, sqlmock.Sqlmock, sqlmock.Sqlmock, func(), error) {
	master, masterMock, closeMaster, err := newMockGormDB()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	slave, slaveMock, closeSlave, err := newMockGormDB()
	if err != nil {
		closeMaster()
		return nil, nil, nil, nil, err
	}
	close := func() {
		closeMaster()
		closeSlave()
	}
	return &Client{Master: master, Slave: slave}, masterMock, slaveMock, close, nil
}

func newMockGormDB() (*gorm.DB, sqlmock.Sqlmock, func(), error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, err
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		_ = sqlDB.Close()
		return nil, nil, nil, err
	}
	return gormDB, mock, func() { _ = sqlDB.Close() }, nil
}

func TestCreateRemediationProposal(t *testing.T) {
	now := time.Now()
	proposal := &model.RemediationProposal{
		RemediationProposalID: 1001,
		FindingID:             1001,
		ProjectID:             1,
		Status:                model.RemediationProposalStatusPending,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	cases := []struct {
		name       string
		want       *model.RemediationProposal
		wantErr    bool
		mockErr    error
		mockGetErr error
	}{
		{
			name:    "OK",
			want:    proposal,
			wantErr: false,
		},
		{
			name:    "NG DB error(insert)",
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
		{
			name:       "NG DB error(select)",
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
			input := &model.RemediationProposal{
				FindingID: 1001,
				ProjectID: 1,
				Status:    model.RemediationProposalStatusPending,
			}
			insertQuery := "INSERT INTO `remediation_proposal`"
			if c.mockErr != nil {
				mock.ExpectBegin()
				mock.ExpectExec(insertQuery).WillReturnError(c.mockErr)
				mock.ExpectRollback()
			} else {
				mock.ExpectBegin()
				mock.ExpectExec(insertQuery).
					WithArgs(input.FindingID, input.ProjectID, input.Status, nil, nil, nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1001, 1))
				mock.ExpectCommit()
				if c.mockGetErr != nil {
					mock.ExpectQuery(regexp.QuoteMeta(selectRemediationProposalByID)).WillReturnError(c.mockGetErr)
				} else {
					mock.ExpectQuery(regexp.QuoteMeta(selectRemediationProposalByID)).
						WithArgs(input.ProjectID, c.want.RemediationProposalID).
						WillReturnRows(newRemediationProposalRows(c.want))
				}
			}

			got, err := client.CreateRemediationProposal(ctx, input)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("No error")
			}
			if !c.wantErr && got.RemediationProposalID != c.want.RemediationProposalID {
				t.Fatalf("Unexpected result: got=%+v, want=%+v", got, c.want)
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
		RemediationProposalID: 1001,
		FindingID:             1001,
		ProjectID:             1,
		Status:                model.RemediationProposalStatusSucceeded,
		RemediationPlan:       &plan,
		GeneratedAt:           &now,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	failed := &model.RemediationProposal{
		RemediationProposalID: 1001,
		FindingID:             1001,
		ProjectID:             1,
		Status:                model.RemediationProposalStatusFailed,
		StatusDetail:          &errMsg,
		CreatedAt:             now,
		UpdatedAt:             now,
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
					WithArgs(c.want.Status, c.want.StatusDetail, c.want.RemediationPlan, c.want.GeneratedAt, c.want.ProjectID, c.want.RemediationProposalID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectQuery(regexp.QuoteMeta(selectRemediationProposalByID)).
					WithArgs(c.want.ProjectID, c.want.RemediationProposalID).
					WillReturnRows(newRemediationProposalRows(c.want))
			}

			got, err := client.UpdateRemediationProposalStatus(ctx, c.want.ProjectID, c.want.RemediationProposalID, c.want.Status, c.want.StatusDetail, c.want.RemediationPlan, c.want.GeneratedAt)
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
