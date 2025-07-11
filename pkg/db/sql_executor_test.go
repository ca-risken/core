package db

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSQLExecutor_validateSQL(t *testing.T) {
	cases := []struct {
		name     string
		inputSQL string
		wantErr  bool
	}{
		{
			name:     "OK - Valid SELECT statement with WHERE",
			inputSQL: "SELECT * FROM finding WHERE status = 'active'",
			wantErr:  false,
		},
		{
			name:     "OK - SELECT with column names containing forbidden words",
			inputSQL: "SELECT updated_at, created_at FROM finding WHERE 1=1",
			wantErr:  false,
		},
		{
			name:     "NG - SELECT without WHERE clause",
			inputSQL: "SELECT * FROM finding",
			wantErr:  true,
		},
		{
			name:     "NG - Non-SELECT statement",
			inputSQL: "UPDATE finding SET status = 'inactive'",
			wantErr:  true,
		},
		{
			name:     "NG - INSERT statement",
			inputSQL: "INSERT INTO finding (data_source) VALUES ('test')",
			wantErr:  true,
		},
		{
			name:     "NG - DELETE statement",
			inputSQL: "DELETE FROM finding WHERE id = 1",
			wantErr:  true,
		},
		{
			name:     "NG - CREATE statement",
			inputSQL: "CREATE TABLE test (id int)",
			wantErr:  true,
		},
		{
			name:     "NG - DROP statement",
			inputSQL: "DROP TABLE finding",
			wantErr:  true,
		},
		{
			name:     "NG - ALTER statement",
			inputSQL: "ALTER TABLE finding ADD COLUMN new_field VARCHAR(255)",
			wantErr:  true,
		},
		{
			name:     "NG - TRUNCATE statement",
			inputSQL: "TRUNCATE TABLE finding",
			wantErr:  true,
		},
		{
			name:     "OK - Empty input",
			inputSQL: "",
			wantErr:  true,
		},
		{
			name:     "OK - Whitespace only",
			inputSQL: "   ",
			wantErr:  true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			executor := &SQLExecutor{}
			err := executor.validateSQL(c.inputSQL)

			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("No error")
			}
		})
	}
}

func TestSQLExecutor_buildSecureSQL(t *testing.T) {
	type args struct {
		sql       string
		projectID uint32
		limit     int32
		offset    int32
	}
	
	cases := []struct {
		name     string
		args     args
		wantSQL  string
		wantArgs []any
	}{
		{
			name: "Simple SELECT with default limit/offset",
			args: args{
				sql:       "SELECT * FROM finding WHERE status = 'active'",
				projectID: 123,
				limit:     0,
				offset:    0,
			},
			wantSQL:  "SELECT * FROM (SELECT * FROM finding WHERE finding.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding pf WHERE pf.finding_id = finding.finding_id AND pf.project_id = ? AND (pf.expired_at IS NULL OR pf.expired_at > NOW())) AND (status = 'active')) AS subquery LIMIT ? OFFSET ?",
			wantArgs: []any{uint32(123), uint32(123), int32(1000), int32(0)},
		},
		{
			name: "SELECT with custom limit/offset",
			args: args{
				sql:       "SELECT * FROM finding WHERE 1=1",
				projectID: 456,
				limit:     50,
				offset:    10,
			},
			wantSQL:  "SELECT * FROM (SELECT * FROM finding WHERE finding.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding pf WHERE pf.finding_id = finding.finding_id AND pf.project_id = ? AND (pf.expired_at IS NULL OR pf.expired_at > NOW())) AND (1=1)) AS subquery LIMIT ? OFFSET ?",
			wantArgs: []any{uint32(456), uint32(456), int32(50), int32(10)},
		},
		{
			name: "SELECT with limit exceeding max (should be capped)",
			args: args{
				sql:       "SELECT * FROM finding WHERE 1=1",
				projectID: 789,
				limit:     2000,
				offset:    0,
			},
			wantSQL:  "SELECT * FROM (SELECT * FROM finding WHERE finding.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding pf WHERE pf.finding_id = finding.finding_id AND pf.project_id = ? AND (pf.expired_at IS NULL OR pf.expired_at > NOW())) AND (1=1)) AS subquery LIMIT ? OFFSET ?",
			wantArgs: []any{uint32(789), uint32(789), int32(1000), int32(0)},
		},
		{
			name: "SELECT with negative offset (should be reset to 0)",
			args: args{
				sql:       "SELECT * FROM finding WHERE 1=1",
				projectID: 123,
				limit:     100,
				offset:    -5,
			},
			wantSQL:  "SELECT * FROM (SELECT * FROM finding WHERE finding.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding pf WHERE pf.finding_id = finding.finding_id AND pf.project_id = ? AND (pf.expired_at IS NULL OR pf.expired_at > NOW())) AND (1=1)) AS subquery LIMIT ? OFFSET ?",
			wantArgs: []any{uint32(123), uint32(123), int32(100), int32(0)},
		},
		{
			name: "SELECT with PEND_FINDING (exclusion always applied)",
			args: args{
				sql:       "SELECT * FROM finding f JOIN pend_finding pf ON f.finding_id = pf.finding_id WHERE 1=1",
				projectID: 123,
				limit:     200,
				offset:    20,
			},
			wantSQL:  "SELECT * FROM (SELECT * FROM finding f JOIN pend_finding pf ON f.finding_id = pf.finding_id WHERE finding.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding pf WHERE pf.finding_id = finding.finding_id AND pf.project_id = ? AND (pf.expired_at IS NULL OR pf.expired_at > NOW())) AND (1=1)) AS subquery LIMIT ? OFFSET ?",
			wantArgs: []any{uint32(123), uint32(123), int32(200), int32(20)},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			executor := &SQLExecutor{}
			result, err := executor.buildSecureSQL(c.args.sql, c.args.projectID, c.args.limit, c.args.offset)

			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if result.SQL != c.wantSQL {
				t.Fatalf("Unexpected SQL:\nwant: %s\ngot:  %s", c.wantSQL, result.SQL)
			}
			if len(result.Args) != len(c.wantArgs) {
				t.Fatalf("Unexpected args length: want=%d, got=%d", len(c.wantArgs), len(result.Args))
			}
			for i, arg := range result.Args {
				if arg != c.wantArgs[i] {
					t.Fatalf("Unexpected arg[%d]: want=%v, got=%v", i, c.wantArgs[i], arg)
				}
			}
		})
	}
}

func TestSQLExecutor_ExecuteSelectQuery(t *testing.T) {
	client, mock, err := newMockClient()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, error: %+v", err)
	}

	type args struct {
		sql       string
		projectID uint32
		limit     int32
		offset    int32
	}
	
	cases := []struct {
		name        string
		args        args
		want        []map[string]any
		wantErr     bool
		mockClosure func(mock sqlmock.Sqlmock)
	}{
		{
			name: "OK - Simple query execution",
			args: args{
				sql:       "SELECT finding_id, data_source FROM finding WHERE status = 'active'",
				projectID: 123,
				limit:     100,
				offset:    0,
			},
			want: []map[string]any{
				{"finding_id": int64(1), "data_source": "aws"},
				{"finding_id": int64(2), "data_source": "gcp"},
			},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				expectedSQL := "SELECT * FROM (SELECT finding_id, data_source FROM finding WHERE finding.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding pf WHERE pf.finding_id = finding.finding_id AND pf.project_id = ? AND (pf.expired_at IS NULL OR pf.expired_at > NOW())) AND (status = 'active')) AS subquery LIMIT ? OFFSET ?"
				rows := sqlmock.NewRows([]string{"finding_id", "data_source"}).
					AddRow(1, "aws").
					AddRow(2, "gcp")
				mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
					WithArgs(123, 123, 100, 0).
					WillReturnRows(rows)
			},
		},
		{
			name: "NG - DB error",
			args: args{
				sql:       "SELECT * FROM finding WHERE 1=1",
				projectID: 123,
				limit:     50,
				offset:    0,
			},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				expectedSQL := "SELECT * FROM (SELECT * FROM finding WHERE finding.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding pf WHERE pf.finding_id = finding.finding_id AND pf.project_id = ? AND (pf.expired_at IS NULL OR pf.expired_at > NOW())) AND (1=1)) AS subquery LIMIT ? OFFSET ?"
				mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
					WithArgs(123, 123, 50, 0).
					WillReturnError(errors.New("DB connection error"))
			},
		},
		{
			name: "NG - Invalid SQL (no WHERE)",
			args: args{
				sql:       "SELECT * FROM finding",
				projectID: 123,
				limit:     100,
				offset:    0,
			},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				// No mock expectation needed as it should fail at validation
			},
		},
		{
			name: "NG - Invalid SQL (forbidden operation)",
			args: args{
				sql:       "INSERT INTO finding VALUES (1, 'test')",
				projectID: 123,
				limit:     100,
				offset:    0,
			},
			want:    nil,
			wantErr: true,
			mockClosure: func(mock sqlmock.Sqlmock) {
				// No mock expectation needed as it should fail at validation
			},
		},
		{
			name: "OK - Query with no results",
			args: args{
				sql:       "SELECT * FROM finding WHERE status = 'nonexistent'",
				projectID: 456,
				limit:     200,
				offset:    10,
			},
			want:    []map[string]any{},
			wantErr: false,
			mockClosure: func(mock sqlmock.Sqlmock) {
				expectedSQL := "SELECT * FROM (SELECT * FROM finding WHERE finding.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding pf WHERE pf.finding_id = finding.finding_id AND pf.project_id = ? AND (pf.expired_at IS NULL OR pf.expired_at > NOW())) AND (status = 'nonexistent')) AS subquery LIMIT ? OFFSET ?"
				rows := sqlmock.NewRows([]string{"finding_id", "data_source"})
				mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
					WithArgs(456, 456, 200, 10).
					WillReturnRows(rows)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			executor := NewSQLExecutor(client)

			if c.mockClosure != nil {
				c.mockClosure(mock)
			}

			got, err := executor.ExecuteSelectQuery(ctx, c.args.sql, c.args.projectID, c.args.limit, c.args.offset)

			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if err == nil && c.wantErr {
				t.Fatal("No error")
			}
			if !c.wantErr && len(got) != len(c.want) {
				t.Fatalf("Unexpected result count: want=%d, got=%d", len(c.want), len(got))
			}
			if !c.wantErr {
				for i, row := range got {
					expectedRow := c.want[i]
					for key, expectedValue := range expectedRow {
						if gotValue, exists := row[key]; !exists || gotValue != expectedValue {
							t.Fatalf("Unexpected value for key %s: want=%v, got=%v", key, expectedValue, gotValue)
						}
					}
				}
			}
		})
	}
}