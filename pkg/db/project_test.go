package db

import (
	"context"
	"errors"
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
