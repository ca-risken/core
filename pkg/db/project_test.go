package db

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDeleteProject(t *testing.T) {
	client, mock, err := newMockClient()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}
	type args struct {
		projectID uint32
	}
	cases := []struct {
		name    string
		input   args
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
			input:   args{projectID: 1},
			wantErr: false,
		},
		{
			name:    "NG DB error",
			input:   args{projectID: 1},
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			if c.mockErr != nil {
				mock.ExpectExec(deleteOrganizationInvitationByProject).WillReturnError(c.mockErr)
			} else {
				mock.ExpectExec(deleteOrganizationInvitationByProject).WillReturnResult(sqlmock.NewResult(int64(1), int64(1)))
				mock.ExpectExec(deleteOrganizationProject).WillReturnResult(sqlmock.NewResult(int64(1), int64(1)))
				mock.ExpectExec(deleteProject).WillReturnResult(sqlmock.NewResult(int64(1), int64(1)))
			}

			err := client.DeleteProject(ctx, c.input.projectID)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("No error")
			}
		})
	}
}

func TestCleanWithNoProject(t *testing.T) {
	client, mock, err := newMockClient()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}
	cases := []struct {
		name    string
		mockSQL []string
		target  *sqlmock.Rows
		wantErr bool
		mockErr error
	}{
		{
			name: "OK",
			target: sqlmock.NewRows(
				[]string{"project_id"}).
				AddRow(1).
				AddRow(2),
			mockSQL: []string{
				"delete from alert where project_id in",
				"delete from alert_history where project_id in",
				"delete from rel_alert_finding where project_id in",
				"delete from alert_condition where project_id in",
				"delete from alert_cond_rule where project_id in",
				"delete from alert_rule where project_id in",
				"delete from alert_cond_notification where project_id in",
				"delete from notification where project_id in",
				"delete from finding where project_id in",
				"delete from finding_tag where project_id in",
				"delete from resource where project_id in",
				"delete from resource_tag where project_id in",
				"delete from pend_finding where project_id in",
				"delete from finding_setting where project_id in",
				"delete from report_finding where project_id in",
				"delete from recommend_finding where project_id in",
				"delete from access_token where project_id in",
				cleanAccessTokenRole,
				"delete from role where project_id in",
				"delete from user_role where project_id in",
				"delete from policy where project_id in",
				"delete from role_policy where project_id in",
				"delete from project_tag where project_id in",
			},
			wantErr: false,
		},
		{
			name: "No target",
			target: sqlmock.NewRows(
				[]string{"project_id"}),
			wantErr: false,
		},
		{
			name: "NG DB error",
			target: sqlmock.NewRows(
				[]string{"project_id"}).
				AddRow(1).
				AddRow(2),
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()

			// target project
			mock.ExpectQuery(regexp.QuoteMeta(selectListCleanProjectTarget)).WillReturnRows(c.target)
			for _, sql := range c.mockSQL {
				mock.ExpectExec(regexp.QuoteMeta(sql)).WillReturnResult(sqlmock.NewResult(int64(1), int64(1)))
			}
			if c.mockErr != nil {
				mock.ExpectExec(regexp.QuoteMeta(`delete from`)).WillReturnError(c.mockErr)
			}

			err := client.CleanWithNoProject(ctx)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("No error")
			}
		})
	}
}
