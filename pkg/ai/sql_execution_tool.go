package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ca-risken/core/pkg/db"
)

// SQLExecutionTool handles SQL execution for AI responses
type SQLExecutionTool struct {
	sqlExecutor *db.SQLExecutor
	projectID   uint32
}

// NewSQLExecutionTool creates a new SQL execution tool
func NewSQLExecutionTool(findingRepo db.FindingRepository, projectID uint32) *SQLExecutionTool {
	return &SQLExecutionTool{
		sqlExecutor: db.NewSQLExecutor(findingRepo.(*db.Client)),
		projectID:   projectID,
	}
}

// ExecuteSQL executes SQL query and returns results
func (t *SQLExecutionTool) ExecuteSQL(ctx context.Context, prompt string, limit int32, offset int32) (*SQLToolResponse, error) {
	// Set defaults
	if limit <= 0 {
		limit = DEFAULT_LIMIT
	}
	if offset < 0 {
		offset = 0
	}

	// For now, generate a basic SQL query based on the prompt
	// In a real implementation, this would use AI to generate the SQL
	sqlQuery := t.generateSQLFromPrompt(prompt)
	
	// Execute the SQL query
	results, err := t.sqlExecutor.ExecuteSelectQuery(ctx, sqlQuery, t.projectID, limit, offset)
	if err != nil {
		return &SQLToolResponse{
			SQL:       sqlQuery,
			ResultSet: nil,
			Error:     fmt.Sprintf("SQL execution failed: %v", err),
		}, nil
	}

	return &SQLToolResponse{
		SQL:       sqlQuery,
		ResultSet: results,
		Error:     "",
	}, nil
}

// generateSQLFromPrompt generates SQL based on the prompt
// This is a simplified implementation - in practice, you'd use AI to generate SQL
func (t *SQLExecutionTool) generateSQLFromPrompt(_ string) string {
	// Basic SQL query for security findings analysis
	return fmt.Sprintf(`
SELECT 
    data_source,
    COUNT(*) as finding_count,
    AVG(score) as avg_score,
    MAX(score) as max_score,
    MIN(score) as min_score,
    description
FROM finding 
WHERE project_id = %d 
    AND NOT EXISTS (
        SELECT 1 FROM pend_finding pf 
        WHERE pf.finding_id = finding.finding_id 
        AND pf.project_id = %d 
        AND (pf.expired_at IS NULL OR pf.expired_at > NOW())
    )
GROUP BY data_source, description
ORDER BY finding_count DESC
LIMIT 100`, t.projectID, t.projectID)
}

// HandleToolCall handles tool calls from AI responses
// Note: This is a placeholder implementation since OpenAI Go SDK tool handling is not yet fully implemented
func (t *SQLExecutionTool) HandleToolCall(ctx context.Context, functionName string, arguments string) (string, error) {
	if functionName != SQL_TOOL_NAME {
		return "", fmt.Errorf("unexpected function name: %s", functionName)
	}

	// Parse the arguments
	var args SQLToolRequest
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse tool arguments: %w", err)
	}

	// Execute the SQL
	result, err := t.ExecuteSQL(ctx, args.Prompt, int32(args.Limit), int32(args.Offset))
	if err != nil {
		return "", fmt.Errorf("SQL execution failed: %w", err)
	}

	// Return the result as JSON
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(resultJSON), nil
}