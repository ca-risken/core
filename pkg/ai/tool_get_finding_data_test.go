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
		wantErr    bool
	}{
		{
			name: "Simple SELECT with WHERE clause",
			args: args{
				sql:       "SELECT * FROM finding WHERE score > 0.5",
				projectID: 1001,
				limit:     100,
				offset:    0,
			},
			wantSQL:    "SELECT * FROM (SELECT * FROM finding WHERE finding.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding pf_scope WHERE pf_scope.project_id = finding.project_id AND pf_scope.finding_id = finding.finding_id AND (pf_scope.expired_at IS NULL OR pf_scope.expired_at > NOW())) AND (score > 0.5)) as t LIMIT ? OFFSET ?",
			wantParams: []any{uint32(1001), uint32(100), uint32(0)},
		},
		{
			name: "OR condition is wrapped inside scoped WHERE clause",
			args: args{
				sql:       "SELECT * FROM finding WHERE score >= 0.8 OR score <= 0.3",
				projectID: 2002,
				limit:     50,
				offset:    10,
			},
			wantSQL:    "SELECT * FROM (SELECT * FROM finding WHERE finding.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding pf_scope WHERE pf_scope.project_id = finding.project_id AND pf_scope.finding_id = finding.finding_id AND (pf_scope.expired_at IS NULL OR pf_scope.expired_at > NOW())) AND (score >= 0.8 OR score <= 0.3)) as t LIMIT ? OFFSET ?",
			wantParams: []any{uint32(2002), uint32(50), uint32(10)},
		},
		{
			name: "Complex SELECT with GROUP BY and ORDER BY",
			args: args{
				sql:       "SELECT data_source, COUNT(*) as count FROM finding WHERE data_source LIKE 'aws%' GROUP BY data_source ORDER BY count DESC",
				projectID: 3003,
				limit:     200,
				offset:    5,
			},
			wantSQL:    "SELECT * FROM (SELECT data_source, COUNT(*) as count FROM finding WHERE finding.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding pf_scope WHERE pf_scope.project_id = finding.project_id AND pf_scope.finding_id = finding.finding_id AND (pf_scope.expired_at IS NULL OR pf_scope.expired_at > NOW())) AND (data_source LIKE 'aws%') GROUP BY data_source ORDER BY count DESC) as t LIMIT ? OFFSET ?",
			wantParams: []any{uint32(3003), uint32(200), uint32(5)},
		},
		{
			name: "lowercase where clause is accepted",
			args: args{
				sql:       "select * from finding where score > 0.5",
				projectID: 4004,
				limit:     10,
				offset:    0,
			},
			wantSQL:    "SELECT * FROM (select * from finding WHERE finding.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding pf_scope WHERE pf_scope.project_id = finding.project_id AND pf_scope.finding_id = finding.finding_id AND (pf_scope.expired_at IS NULL OR pf_scope.expired_at > NOW())) AND (score > 0.5)) as t LIMIT ? OFFSET ?",
			wantParams: []any{uint32(4004), uint32(10), uint32(0)},
		},
		{
			name: "SELECT with alias and OR conditions",
			args: args{
				sql:       "SELECT f.finding_id FROM finding f WHERE f.score > 0.5 OR f.data_source LIKE 'aws%'",
				projectID: 5005,
				limit:     25,
				offset:    3,
			},
			wantSQL:    "SELECT * FROM (SELECT f.finding_id FROM finding f WHERE f.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding pf_scope WHERE pf_scope.project_id = f.project_id AND pf_scope.finding_id = f.finding_id AND (pf_scope.expired_at IS NULL OR pf_scope.expired_at > NOW())) AND (f.score > 0.5 OR f.data_source LIKE 'aws%')) as t LIMIT ? OFFSET ?",
			wantParams: []any{uint32(5005), uint32(25), uint32(3)},
		},
		{
			name: "SELECT with JSON field extraction",
			args: args{
				sql:       "SELECT finding_id, JSON_EXTRACT(data, '$.severity') as severity FROM finding WHERE updated_at > '2024-01-01'",
				projectID: 6006,
				limit:     10,
				offset:    0,
			},
			wantSQL:    "SELECT * FROM (SELECT finding_id, JSON_EXTRACT(data, '$.severity') as severity FROM finding WHERE finding.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding pf_scope WHERE pf_scope.project_id = finding.project_id AND pf_scope.finding_id = finding.finding_id AND (pf_scope.expired_at IS NULL OR pf_scope.expired_at > NOW())) AND (updated_at > '2024-01-01')) as t LIMIT ? OFFSET ?",
			wantParams: []any{uint32(6006), uint32(10), uint32(0)},
		},
		{
			name: "SELECT with aggregate functions",
			args: args{
				sql:       "SELECT data_source, AVG(score) as avg_score, MAX(score) as max_score FROM finding WHERE (score >= 0.8 OR score <= 0.3) AND resource_name LIKE '%bucket%' GROUP BY data_source",
				projectID: 7007,
				limit:     25,
				offset:    3,
			},
			wantSQL:    "SELECT * FROM (SELECT data_source, AVG(score) as avg_score, MAX(score) as max_score FROM finding WHERE finding.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding pf_scope WHERE pf_scope.project_id = finding.project_id AND pf_scope.finding_id = finding.finding_id AND (pf_scope.expired_at IS NULL OR pf_scope.expired_at > NOW())) AND ((score >= 0.8 OR score <= 0.3) AND resource_name LIKE '%bucket%') GROUP BY data_source) as t LIMIT ? OFFSET ?",
			wantParams: []any{uint32(7007), uint32(25), uint32(3)},
		},
		{
			name: "SQL with semicolon is rejected",
			args: args{
				sql:       "SELECT * FROM finding WHERE score > 0.5; SELECT * FROM finding WHERE score < 0.5",
				projectID: 8008,
				limit:     10,
				offset:    0,
			},
			wantErr: true,
		},
		{
			name: "JOIN is rejected",
			args: args{
				sql:       "SELECT f.finding_id FROM finding f JOIN project p ON f.project_id = p.project_id WHERE f.score > 0.5",
				projectID: 9009,
				limit:     10,
				offset:    0,
			},
			wantErr: true,
		},
		{
			name: "subquery is rejected",
			args: args{
				sql:       "SELECT * FROM finding WHERE finding_id IN (SELECT finding_id FROM finding WHERE score > 0.8)",
				projectID: 9010,
				limit:     10,
				offset:    0,
			},
			wantErr: true,
		},
		{
			name: "internal pend_finding alias avoids collision with finding alias",
			args: args{
				sql:       "SELECT pf_scope.finding_id FROM finding pf_scope WHERE pf_scope.score > 0.5",
				projectID: 9011,
				limit:     10,
				offset:    0,
			},
			wantSQL:    "SELECT * FROM (SELECT pf_scope.finding_id FROM finding pf_scope WHERE pf_scope.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding pf_scope_1 WHERE pf_scope_1.project_id = pf_scope.project_id AND pf_scope_1.finding_id = pf_scope.finding_id AND (pf_scope_1.expired_at IS NULL OR pf_scope_1.expired_at > NOW())) AND (pf_scope.score > 0.5)) as t LIMIT ? OFFSET ?",
			wantParams: []any{uint32(9011), uint32(10), uint32(0)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL, gotParams, err := formatSQL(tt.args.sql, tt.args.projectID, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Fatalf("formatSQL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

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
			name: "Valid SELECT with alias and WHERE",
			args: args{
				sql: "SELECT f.finding_id FROM finding f WHERE f.score > 0.8",
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
			name: "Case insensitive FROM and WHERE",
			args: args{
				sql: "select finding_id from finding Where score > 0.5",
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
			name: "JOIN is rejected",
			args: args{
				sql: "SELECT f.*, p.* FROM finding f JOIN project p ON f.project_id = p.project_id WHERE f.score > 0.8",
			},
			wantErr: true,
		},
		{
			name: "subquery is rejected",
			args: args{
				sql: "SELECT * FROM finding WHERE finding_id IN (SELECT finding_id FROM finding WHERE score > 0.5)",
			},
			wantErr: true,
		},
		{
			name: "Finding only source is required",
			args: args{
				sql: "SELECT * FROM project WHERE project_id = 1",
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
