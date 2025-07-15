package ai

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFormatSQL(t *testing.T) {
	type args struct {
		sql       string
		projectID uint32
		limit     uint32
		offset    uint32
	}
	tests := []struct {
		name       string
		args       args
		wantSQL    string
		wantParams []any
	}{
		{
			name: "Simple SELECT with WHERE clause",
			args: args{
				sql:       "SELECT * FROM finding WHERE score > 0.5",
				projectID: 1001,
				limit:     100,
				offset:    0,
			},
			wantSQL: `SELECT * FROM (SELECT * FROM finding WHERE
	project_id = ? 
	AND not exists (
		SELECT 1 
		FROM pend_finding
		WHERE 
		  pend_finding.finding_id = finding.finding_id
			and (pend_finding.expired_at is NULL or pend_finding.expired_at > NOW())
	)
	AND score > 0.5) as t LIMIT ? OFFSET ?`,
			wantParams: []any{uint32(1001), uint32(100), uint32(0)},
		},
		{
			name: "Complex SELECT with GROUP BY and ORDER BY",
			args: args{
				sql:       "SELECT data_source, COUNT(*) as count FROM finding WHERE data_source LIKE 'aws%' GROUP BY data_source ORDER BY count DESC",
				projectID: 2002,
				limit:     50,
				offset:    10,
			},
			wantSQL: `SELECT * FROM (SELECT data_source, COUNT(*) as count FROM finding WHERE
	project_id = ? 
	AND not exists (
		SELECT 1 
		FROM pend_finding
		WHERE 
		  pend_finding.finding_id = finding.finding_id
			and (pend_finding.expired_at is NULL or pend_finding.expired_at > NOW())
	)
	AND data_source LIKE 'aws%' GROUP BY data_source ORDER BY count DESC) as t LIMIT ? OFFSET ?`,
			wantParams: []any{uint32(2002), uint32(50), uint32(10)},
		},
		{
			name: "SELECT with multiple WHERE conditions",
			args: args{
				sql:       "SELECT * FROM finding WHERE score >= 0.8 AND data_source = 'aws:guardduty'",
				projectID: 3003,
				limit:     200,
				offset:    5,
			},
			wantSQL: `SELECT * FROM (SELECT * FROM finding WHERE
	project_id = ? 
	AND not exists (
		SELECT 1 
		FROM pend_finding
		WHERE 
		  pend_finding.finding_id = finding.finding_id
			and (pend_finding.expired_at is NULL or pend_finding.expired_at > NOW())
	)
	AND score >= 0.8 AND data_source = 'aws:guardduty') as t LIMIT ? OFFSET ?`,
			wantParams: []any{uint32(3003), uint32(200), uint32(5)},
		},
		{
			name: "SELECT with JSON field extraction",
			args: args{
				sql:       "SELECT finding_id, JSON_EXTRACT(data, '$.severity') as severity FROM finding WHERE updated_at > '2024-01-01'",
				projectID: 4004,
				limit:     10,
				offset:    0,
			},
			wantSQL: `SELECT * FROM (SELECT finding_id, JSON_EXTRACT(data, '$.severity') as severity FROM finding WHERE
	project_id = ? 
	AND not exists (
		SELECT 1 
		FROM pend_finding
		WHERE 
		  pend_finding.finding_id = finding.finding_id
			and (pend_finding.expired_at is NULL or pend_finding.expired_at > NOW())
	)
	AND updated_at > '2024-01-01') as t LIMIT ? OFFSET ?`,
			wantParams: []any{uint32(4004), uint32(10), uint32(0)},
		},
		{
			name: "SELECT with aggregate functions",
			args: args{
				sql:       "SELECT data_source, AVG(score) as avg_score, MAX(score) as max_score FROM finding WHERE resource_name LIKE '%bucket%' GROUP BY data_source",
				projectID: 5005,
				limit:     25,
				offset:    3,
			},
			wantSQL: `SELECT * FROM (SELECT data_source, AVG(score) as avg_score, MAX(score) as max_score FROM finding WHERE
	project_id = ? 
	AND not exists (
		SELECT 1 
		FROM pend_finding
		WHERE 
		  pend_finding.finding_id = finding.finding_id
			and (pend_finding.expired_at is NULL or pend_finding.expired_at > NOW())
	)
	AND resource_name LIKE '%bucket%' GROUP BY data_source) as t LIMIT ? OFFSET ?`,
			wantParams: []any{uint32(5005), uint32(25), uint32(3)},
		},
		{
			name: "SQL with semicolon - should trim after semicolon",
			args: args{
				sql:       "SELECT * FROM finding WHERE score > 0.5; SELECT * FROM finding WHERE score < 0.5",
				projectID: 6006,
				limit:     10,
				offset:    0,
			},
			wantSQL: `SELECT * FROM (SELECT * FROM finding WHERE
	project_id = ? 
	AND not exists (
		SELECT 1 
		FROM pend_finding
		WHERE 
		  pend_finding.finding_id = finding.finding_id
			and (pend_finding.expired_at is NULL or pend_finding.expired_at > NOW())
	)
	AND score > 0.5) as t LIMIT ? OFFSET ?`,
			wantParams: []any{uint32(6006), uint32(10), uint32(0)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL, gotParams := formatSQL(tt.args.sql, tt.args.projectID, tt.args.limit, tt.args.offset)

			if diff := cmp.Diff(tt.wantSQL, gotSQL); diff != "" {
				t.Errorf("formatSQL() SQL mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantParams, gotParams); diff != "" {
				t.Errorf("formatSQL() params mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestValidateSQL(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid SELECT with WHERE",
			args: args{
				sql: "SELECT * FROM finding WHERE score > 0.5",
			},
			wantErr: false,
		},
		{
			name: "Valid complex SELECT with WHERE",
			args: args{
				sql: "SELECT data_source, COUNT(*) FROM finding WHERE data_source LIKE 'aws%' GROUP BY data_source ORDER BY COUNT(*) DESC",
			},
			wantErr: false,
		},
		{
			name: "Valid SELECT with JOIN and WHERE",
			args: args{
				sql: "SELECT f.*, p.* FROM finding f JOIN project p ON f.project_id = p.project_id WHERE f.score > 0.8",
			},
			wantErr: false,
		},
		{
			name: "Valid SELECT with subquery",
			args: args{
				sql: "SELECT * FROM (SELECT * FROM finding WHERE score > 0.5) subquery WHERE created_at > '2024-01-01'",
			},
			wantErr: false,
		},
		{
			name: "Valid SELECT with CASE statement",
			args: args{
				sql: `SELECT 
					finding_id,
					CASE 
						WHEN score > 0.8 THEN 'High'
						WHEN score > 0.4 THEN 'Medium'
						ELSE 'Low'
					END as risk_level
				FROM finding WHERE project_id = 1001`,
			},
			wantErr: false,
		},
		{
			name: "Missing SELECT keyword",
			args: args{
				sql: "* FROM finding WHERE score > 0.5",
			},
			wantErr: true,
		},
		{
			name: "Missing WHERE keyword",
			args: args{
				sql: "SELECT * FROM finding",
			},
			wantErr: true,
		},
		{
			name: "Case insensitive SELECT",
			args: args{
				sql: "select * from finding where score > 0.5",
			},
			wantErr: false,
		},
		{
			name: "Case insensitive WHERE",
			args: args{
				sql: "SELECT * FROM finding Where score > 0.5",
			},
			wantErr: false,
		},
		{
			name: "Empty SQL string",
			args: args{
				sql: "",
			},
			wantErr: true,
		},
		{
			name: "SQL with only whitespace",
			args: args{
				sql: "   \n\t   ",
			},
			wantErr: true,
		},
		{
			name: "Non-SELECT statement (INSERT)",
			args: args{
				sql: "INSERT INTO finding (score) VALUES (0.5) WHERE project_id = 1",
			},
			wantErr: true,
		},
		{
			name: "Non-SELECT statement (UPDATE)",
			args: args{
				sql: "UPDATE finding SET score = 0.5 WHERE project_id = 1",
			},
			wantErr: true,
		},
		{
			name: "Non-SELECT statement (DELETE)",
			args: args{
				sql: "DELETE FROM finding WHERE project_id = 1",
			},
			wantErr: true,
		},
		{
			name: "SQL with semicolon - multiple statements",
			args: args{
				sql: "SELECT * FROM finding WHERE score > 0.5; SELECT * FROM finding WHERE score < 0.5",
			},
			wantErr: true,
		},
		{
			name: "SQL with semicolon at end",
			args: args{
				sql: "SELECT * FROM finding WHERE score > 0.5;",
			},
			wantErr: true,
		},
		{
			name: "SQL with semicolon in string literal should still fail",
			args: args{
				sql: "SELECT * FROM finding WHERE description = 'test; value' AND score > 0.5",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSQL(tt.args.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateSQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
