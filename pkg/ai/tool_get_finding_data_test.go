package ai

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func buildExpectedFormattedSQL(head, whereCondition, tail string) string {
	scopedSQL := fmt.Sprintf(`%s WHERE
	project_id = ?
	AND not exists (
		SELECT 1
		FROM pend_finding
		WHERE
		  pend_finding.finding_id = finding.finding_id
			and (pend_finding.expired_at is NULL or pend_finding.expired_at > NOW())
	)
	AND (%s)`, head, whereCondition)
	if tail != "" {
		scopedSQL = fmt.Sprintf("%s %s", scopedSQL, tail)
	}
	return fmt.Sprintf(`SELECT * FROM (%s) as t LIMIT ? OFFSET ?`, scopedSQL)
}

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
			name: "OR condition is enclosed by parentheses",
			args: args{
				sql:       "SELECT * FROM finding WHERE score >= 0.8 OR score <= 0.3",
				projectID: 1001,
				limit:     20,
				offset:    0,
			},
			wantSQL:    buildExpectedFormattedSQL("SELECT * FROM finding", "score >= 0.8 OR score <= 0.3", ""),
			wantParams: []any{uint32(1001), uint32(20), uint32(0)},
		},
		{
			name: "AND and OR mixed condition is enclosed as one expression",
			args: args{
				sql:       "SELECT * FROM finding WHERE score >= 0.8 OR score <= 0.3 AND data_source LIKE 'aws%'",
				projectID: 2002,
				limit:     50,
				offset:    10,
			},
			wantSQL:    buildExpectedFormattedSQL("SELECT * FROM finding", "score >= 0.8 OR score <= 0.3 AND data_source LIKE 'aws%'", ""),
			wantParams: []any{uint32(2002), uint32(50), uint32(10)},
		},
		{
			name: "lowercase where is accepted",
			args: args{
				sql:       "select * from finding where score > 0.5",
				projectID: 3003,
				limit:     15,
				offset:    5,
			},
			wantSQL:    buildExpectedFormattedSQL("select * from finding", "score > 0.5", ""),
			wantParams: []any{uint32(3003), uint32(15), uint32(5)},
		},
		{
			name: "mixed-case Where is accepted",
			args: args{
				sql:       "SELECT * FROM finding Where score > 0.5",
				projectID: 4004,
				limit:     10,
				offset:    0,
			},
			wantSQL:    buildExpectedFormattedSQL("SELECT * FROM finding", "score > 0.5", ""),
			wantParams: []any{uint32(4004), uint32(10), uint32(0)},
		},
		{
			name: "backticked finding table is accepted",
			args: args{
				sql:       "SELECT * FROM `finding` WHERE score > 0.5",
				projectID: 4504,
				limit:     12,
				offset:    1,
			},
			wantSQL:    buildExpectedFormattedSQL("SELECT * FROM `finding`", "score > 0.5", ""),
			wantParams: []any{uint32(4504), uint32(12), uint32(1)},
		},
		{
			name: "GROUP BY and ORDER BY are preserved",
			args: args{
				sql:       "SELECT data_source, AVG(score) as avg_score FROM finding WHERE resource_name LIKE '%bucket%' GROUP BY data_source ORDER BY avg_score DESC",
				projectID: 5005,
				limit:     25,
				offset:    3,
			},
			wantSQL:    buildExpectedFormattedSQL("SELECT data_source, AVG(score) as avg_score FROM finding", "resource_name LIKE '%bucket%'", "GROUP BY data_source ORDER BY avg_score DESC"),
			wantParams: []any{uint32(5005), uint32(25), uint32(3)},
		},
		{
			name: "where and order by in string literal do not break parsing",
			args: args{
				sql:       "SELECT * FROM finding WHERE description = 'where order by' AND score > 0.5 ORDER BY score DESC",
				projectID: 6006,
				limit:     10,
				offset:    0,
			},
			wantSQL:    buildExpectedFormattedSQL("SELECT * FROM finding", "description = 'where order by' AND score > 0.5", "ORDER BY score DESC"),
			wantParams: []any{uint32(6006), uint32(10), uint32(0)},
		},
		{
			name: "comment and newline in where clause are preserved",
			args: args{
				sql: `SELECT * FROM finding
WHERE score > 0.5 -- where in comment
AND data_source LIKE 'aws%'
ORDER BY score DESC`,
				projectID: 6106,
				limit:     11,
				offset:    2,
			},
			wantSQL: buildExpectedFormattedSQL("SELECT * FROM finding", `score > 0.5 -- where in comment
AND data_source LIKE 'aws%'`, "ORDER BY score DESC"),
			wantParams: []any{uint32(6106), uint32(11), uint32(2)},
		},
		{
			name: "fail closed when where clause is missing",
			args: args{
				sql:       "SELECT * FROM finding",
				projectID: 7007,
				limit:     20,
				offset:    0,
			},
			wantErr: true,
		},
		{
			name: "fail closed when multiple statements are included",
			args: args{
				sql:       "SELECT * FROM finding WHERE score > 0.5; SELECT * FROM finding WHERE score < 0.5",
				projectID: 8008,
				limit:     20,
				offset:    0,
			},
			wantErr: true,
		},
		{
			name: "fail closed when subquery is used",
			args: args{
				sql:       "SELECT * FROM finding WHERE finding_id IN (SELECT finding_id FROM finding WHERE score > 0.5)",
				projectID: 9009,
				limit:     20,
				offset:    0,
			},
			wantErr: true,
		},
		{
			name: "fail closed when join is used",
			args: args{
				sql:       "SELECT f.* FROM finding f JOIN project p ON f.project_id = p.project_id WHERE f.score > 0.8",
				projectID: 10010,
				limit:     20,
				offset:    0,
			},
			wantErr: true,
		},
		{
			name: "fail closed when alias declaration is invalid",
			args: args{
				sql:       "SELECT * FROM finding AS WHERE score > 0.5",
				projectID: 10510,
				limit:     20,
				offset:    0,
			},
			wantErr: true,
		},
		{
			name: "fail closed when top-level where appears multiple times",
			args: args{
				sql:       "SELECT * FROM finding WHERE score > 0.5 WHERE created_at > '2024-01-01'",
				projectID: 11011,
				limit:     20,
				offset:    0,
			},
			wantErr: true,
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
