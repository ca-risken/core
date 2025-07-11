package ai

import (
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared/constant"
)

const (
	SQL_TOOL_NAME        = "sql_query_executor"
	SQL_TOOL_DESCRIPTION = "Execute SQL queries against the finding and pend_finding tables to analyze security data. This tool generates SQL based on natural language prompts and returns the results."
	MAX_QUERY_LIMIT      = 1000
	DEFAULT_LIMIT        = 100
)

type SQLToolRequest struct {
	Prompt    string `json:"prompt"`
	ProjectID uint32 `json:"project_id"`
	Limit     int    `json:"limit,omitempty"`
	Offset    int    `json:"offset,omitempty"`
}

type SQLToolResponse struct {
	SQL       string `json:"sql"`
	ResultSet any    `json:"resultset"`
	Error     string `json:"error,omitempty"`
}

// GetSQLExecutorTool returns a function calling tool for SQL execution
func GetSQLExecutorTool() responses.ToolUnionParam {
	parameters := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"prompt": map[string]any{
				"type":        "string",
				"description": "Natural language prompt describing the data analysis query",
			},
			"project_id": map[string]any{
				"type":        "integer",
				"description": "Project ID to scope the query",
			},
			"limit": map[string]any{
				"type":        "integer",
				"description": "Maximum number of results to return (default: 100, max: 1000)",
				"minimum":     1,
				"maximum":     1000,
			},
			"offset": map[string]any{
				"type":        "integer",
				"description": "Number of results to skip (default: 0)",
				"minimum":     0,
			},
		},
		"required": []string{"prompt", "project_id"},
	}

	return responses.ToolUnionParam{
		OfFunction: &responses.FunctionToolParam{
			Name:        SQL_TOOL_NAME,
			Description: param.NewOpt(SQL_TOOL_DESCRIPTION),
			Parameters:  parameters,
			Strict:      param.NewOpt(true),
			Type:        constant.Function("function"),
		},
	}
}