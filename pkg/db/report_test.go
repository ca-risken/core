package db

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPurgeReportFinding(t *testing.T) {
	client, mock, err := newMockClient()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}
	cases := []struct {
		name    string
		wantErr bool
		mockErr error
	}{
		{
			name:    "OK",
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
			if c.mockErr != nil {
				mock.ExpectExec("delete from report_finding").WillReturnError(c.mockErr)
			} else {
				mock.ExpectExec("delete from report_finding").WillReturnResult(sqlmock.NewResult(int64(1), int64(1)))
			}
			err := client.PurgeReportFinding(ctx)
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("No error")
			}
		})
	}
}
