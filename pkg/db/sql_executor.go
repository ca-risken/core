package db

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

const (
	MAX_SQL_QUERY_LIMIT = 1000
)


type SQLExecutor struct {
	client *Client
}

func NewSQLExecutor(client *Client) *SQLExecutor {
	return &SQLExecutor{client: client}
}

func (s *SQLExecutor) ExecuteSelectQuery(ctx context.Context, sql string, projectID uint32, limit int32, offset int32) ([]map[string]any, error) {
	// Validate SQL first
	if err := s.validateSQL(sql); err != nil {
		return nil, fmt.Errorf("SQL validation failed: %w", err)
	}

	// Build secure SQL with parameters
	sqlResult, err := s.buildSecureSQL(sql, projectID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("SQL security build failed: %w", err)
	}

	// Execute the query and return results as a slice of maps
	var results []map[string]any
	rows, err := s.client.Slave.WithContext(ctx).Raw(sqlResult.SQL, sqlResult.Args...).Rows()
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	// Process each row
	for rows.Next() {
		// Create a slice of any to hold the values
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Create a map for this row
		row := make(map[string]any)
		for i, column := range columns {
			row[column] = values[i]
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func (s *SQLExecutor) validateSQL(sql string) error {
	sql = strings.TrimSpace(sql)

	// Ensure it's a SELECT statement
	if !strings.HasPrefix(strings.ToUpper(sql), "SELECT") {
		return fmt.Errorf("only SELECT statements are allowed")
	}

	// Ensure WHERE clause is present
	if !strings.Contains(strings.ToUpper(sql), " WHERE ") {
		return fmt.Errorf("WHERE clause is required. If no conditions are needed, use 'WHERE 1=1'")
	}

	// Check for forbidden operations (word boundaries to avoid false positives)
	forbidden := []string{"INSERT", "UPDATE", "DELETE", "CREATE", "DROP", "ALTER", "TRUNCATE", "GRANT", "REVOKE"}
	upperSQL := strings.ToUpper(sql)
	for _, op := range forbidden {
		// Use word boundaries to avoid matching within column names
		pattern := `\b` + op + `\b`
		matched, err := regexp.MatchString(pattern, upperSQL)
		if err != nil {
			return fmt.Errorf("regex error checking forbidden operation %s: %w", op, err)
		}
		if matched {
			return fmt.Errorf("operation %s is not allowed", op)
		}
	}

	return nil
}

type SQLResult struct {
	SQL  string
	Args []any
}

func (s *SQLExecutor) buildSecureSQL(sql string, projectID uint32, limit int32, offset int32) (*SQLResult, error) {
	sql = strings.TrimSpace(sql)

	// Build security constraints with placeholders
	securityConstraints := "finding.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding pf WHERE pf.finding_id = finding.finding_id AND pf.project_id = ? AND (pf.expired_at IS NULL OR pf.expired_at > NOW()))"
	args := []any{projectID, projectID}

	// Replace WHERE clause (validation ensures WHERE exists)
	whereRegex := regexp.MustCompile(`(?i)\bWHERE\s+`)
	sql = whereRegex.ReplaceAllString(sql, fmt.Sprintf("WHERE %s AND (", securityConstraints)) + ")"

	// Apply limit with max constraint
	if limit <= 0 || limit > MAX_SQL_QUERY_LIMIT {
		limit = MAX_SQL_QUERY_LIMIT
	}

	// Apply offset (ensure non-negative)
	if offset < 0 {
		offset = 0
	}

	// Wrap with LIMIT and OFFSET using subquery with prepared statement placeholders
	secureSQL := fmt.Sprintf("SELECT * FROM (%s) AS subquery LIMIT ? OFFSET ?", sql)
	args = append(args, limit, offset)

	return &SQLResult{
		SQL:  secureSQL,
		Args: args,
	}, nil
}