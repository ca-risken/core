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
				mock.ExpectExec(deleteProject).WillReturnError(c.mockErr)
			} else {
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
		wantErr bool
		mockErr error
	}{
		{
			name: "OK",
			mockSQL: []string{
				"delete tbl from alert tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from alert_history tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from rel_alert_finding tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from alert_condition tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from alert_cond_rule tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from alert_rule tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from alert_cond_notification tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from notification tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from finding tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from finding_tag tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from resource tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from resource_tag tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from pend_finding tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from finding_setting tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from report_finding tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from recommend_finding tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from access_token tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from access_token_role tbl where not exists(select * from access_token at where at.access_token_id = tbl.access_token_id)",
				"delete tbl from role tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from user_role tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from policy tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from role_policy tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
				"delete tbl from project_tag tbl where tbl.project_id is not null and not exists(select * from project p where p.project_id = tbl.project_id) ",
			},
			wantErr: false,
		},
		{
			name:    "NG DB error",
			wantErr: true,
			mockErr: errors.New("DB error"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			for _, sql := range c.mockSQL {
				mock.ExpectExec(regexp.QuoteMeta(sql)).WillReturnResult(sqlmock.NewResult(int64(1), int64(1)))
			}
			if c.mockErr != nil {
				mock.ExpectExec(regexp.QuoteMeta(`delete tbl from`)).WillReturnError(c.mockErr)
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
