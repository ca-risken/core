package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
)

// GetFindingDataTool returns the get_finding_data function tool definition
func GetFindingDataTool() responses.ToolUnionParam {
	return responses.ToolUnionParam{
		OfFunction: &responses.FunctionToolParam{
			Name: "get_finding_data",
			Description: openai.String(
				"Get RISKEN findings in a project. " +
					"Use this when a request include \"finding\", \"issue\", \"ファインディング\", \"問題\"..." +
					"Only data that can be retrieved with a single SQL query per execution will be returned.",
			),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"project_id": map[string]any{
						"type":        "integer",
						"description": "Project ID to get findings for",
					},
					"prompt": map[string]any{
						"type":        "string",
						"description": "Prompt to get findings. Generate SQL from natural language. e.g. \"Summarize aws findings by data_source. and sort by data_source name.\"",
					},
					"limit": map[string]any{
						"type":        "integer",
						"description": "Maximum number of results to return. (default: 100)",
					},
					"offset": map[string]any{
						"type":        "integer",
						"description": "Offset of the results to return. (default: 0)",
					},
				},
				"required": []string{"project_id", "prompt"},
			},
		},
	}
}

// GetFindingDataParams represents the parameters for get_finding_data function
type GetFindingDataParams struct {
	ProjectID uint32 `json:"project_id"`
	Prompt    string `json:"prompt"`
	Limit     uint32 `json:"limit"`
	Offset    uint32 `json:"offset"`
}

// GetFindingDataFunction handles the get_finding_data function call
func (a *AIClient) GetFindingDataFunction(ctx context.Context, params GetFindingDataParams) ([]map[string]any, error) {
	// Default params
	if params.Limit == 0 {
		params.Limit = 100
	}

	// Retry conf
	retryer := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3) // Max 3 retries
	retryLogger := func(err error, t time.Duration) {
		a.logger.Warnf(ctx, "[RetryLogger] GetFindingDataFunction error: duration=%+v, err=%+v", t, err)
	}

	// Operation
	var results []map[string]any
	var err error
	operation := func() error {
		results, err = a.getFindingDataFunction(ctx, params)
		if err != nil {
			a.logger.Warnf(ctx, "[RetryLogger] GetFindingDataFunction error: err=%+v", err)
		}
		return err
	}

	// Exec with retry
	if err := backoff.RetryNotify(operation, retryer, retryLogger); err != nil {
		return nil, err
	}
	return results, nil
}

func (a *AIClient) getFindingDataFunction(ctx context.Context, params GetFindingDataParams) ([]map[string]any, error) {
	sql, sqlParams, err := a.generateSQL(ctx, params.Prompt, params.ProjectID, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	a.logger.Infof(ctx, "Generated SQL: sql=%s, params=%v", sql, sqlParams)

	if err := validateSQL(sql); err != nil {
		return nil, err
	}

	data, err := a.executeSQL(ctx, sql, sqlParams)
	if err != nil {
		return nil, err
	}
	a.logger.Infof(ctx, "Successfully executed SQL: len(data)=%d", len(data))
	return data, nil
}

const TOOL_GENERATE_SQL_INSTRUCTION = `
You are an AI that generates SQL queries to retrieve data from the RISKEN Finding table.

## Table schema
- DDL
CREATE TABLE finding (
  finding_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  description VARCHAR(200) NULL,
  data_source VARCHAR(64) NOT NULL,
  data_source_id VARCHAR(255) NOT NULL,
  resource_name VARCHAR(512) NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  original_score FLOAT(5,2) UNSIGNED NOT NULL,
  score FLOAT(3,2) UNSIGNED NULL,
  data JSON NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(finding_id),
  UNIQUE KEY uidx_data_source (project_id, data_source, data_source_id),
  INDEX idx_score(project_id, score, updated_at)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

- data_source: This is the identifier of the modified source data.
	- e.g. "aws:*", "google:*", "azure:*", "code:*", "osint:*", "diagnosis:*", 
- resource_name: This is a name of cloud resource.
	- e.g. "arn:aws:iam::123456789012:user/john.doe", "//bigquery.googleapis.com/projects/pj-name/datasets/ds-name"
- score: Findings always have a score. Although score evaluation varies across data sources, when registered in RISKEN, scores are standardized to a number between 0.0 ~ 1.0. The score helps filter high-risk data and support prioritization decisions.
	- 1.0 ~ 0.8: High risk
	- 0.7 ~ 0.4: Medium risk
	- 0.3 ~ 0.0: Low risk
- data: This is a JSON object that contains the full data of the finding.
  - **IMPORTANT**: The data structure varies depending on the data_source or finding types.

## Strict Rules
- Must return a valid JSON object. No other text.
- Must return a valid SQL string.
- Do NOT include project_id condition in the WHERE clause(it will be automatically inserted when retrieving data).
- Must has prefix "SELECT" and include "WHERE" keywords.
- Do NOT use dangerous keywords like "INSERT", "UPDATE", "DELETE" etc.
- Do NOT include semicolons (;) in the SQL(multiple statements are not allowed).
- Generate only ONE SELECT statement per request.

## Output Schema
{
	"sql": string
}

## Output Example
{"sql": "SELECT * FROM finding WHERE project_id = 1001 AND score >= 0.8 AND data_source like 'aws%' GROUP BY data_source ORDER BY data_source"}
`

type GenerateSQLOutput struct {
	SQL string `json:"sql"`
}

func (a *AIClient) generateSQL(ctx context.Context, prompt string, projectID, limit, offset uint32) (string, []any, error) {
	resp, err := a.callResponsesAPI(ctx,
		a.chatGPTModel,
		TOOL_GENERATE_SQL_INSTRUCTION,
		responses.ResponseNewParamsInputUnion{
			OfInputItemList: []responses.ResponseInputItemUnionParam{
				{
					OfMessage: &responses.EasyInputMessageParam{
						Role: responses.EasyInputMessageRoleUser,
						Content: responses.EasyInputMessageContentUnionParam{
							OfString: openai.String(prompt),
						},
					},
				},
			},
		}, DefaultTools)
	if err != nil {
		return "", nil, err
	}
	jsonOutput, err := ExtractJSONString(resp.OutputText())
	if err != nil {
		return "", nil, fmt.Errorf("failed to extract JSON string: err=%w, output=%s", err, resp.OutputText())
	}
	output, err := ConvertSchema(jsonOutput, GenerateSQLOutput{})
	if err != nil {
		return "", nil, err
	}

	sql, params := formatSQL(output.SQL, projectID, limit, offset)
	if err := validateSQL(sql); err != nil {
		return "", nil, err
	}
	return sql, params, nil
}

func formatSQL(sql string, projectID, limit, offset uint32) (string, []any) {
	params := []any{}

	// Trim SQL after semicolon if exists (defensive programming)
	if idx := strings.Index(sql, ";"); idx != -1 {
		sql = strings.TrimSpace(sql[:idx])
	}

	// strict project_id filter & ignore pend_findings
	sql = strings.ReplaceAll(sql, "WHERE", `WHERE
	project_id = ? 
	AND not exists (
		SELECT 1 
		FROM pend_finding
		WHERE 
		  pend_finding.finding_id = finding.finding_id
			and (pend_finding.expired_at is NULL or pend_finding.expired_at > NOW())
	)
	AND`)
	params = append(params, projectID)

	// add limit and offset
	sql = fmt.Sprintf(`SELECT * FROM (%s) as t LIMIT ? OFFSET ?`, sql)
	params = append(params, limit, offset)
	return sql, params
}

func validateSQL(sql string) error {
	// keyword check
	sqlUpper := strings.ToUpper(sql)
	if !strings.HasPrefix(sqlUpper, "SELECT") || !strings.Contains(sqlUpper, "WHERE") || !strings.Contains(sqlUpper, "FINDING") {
		return fmt.Errorf("sql must contain SELECT and WHERE keywords, sql=%s", sql)
	}

	// Check for semicolons (multiple SQL statements)
	if strings.Contains(sql, ";") {
		return fmt.Errorf("sql must not contain semicolons (multiple statements not allowed), sql=%s", sql)
	}

	// Dangerous keywords check (case-insensitive)
	dangerousKeywords := []string{
		"INSERT ", "UPDATE ", "DELETE ", "DROP TABLE ", "CREATE TABLE ", "ALTER TABLE ",
		"TRUNCATE ", "EXEC ", "EXECUTE ", "INTO OUTFILE ",
		"LOAD_FILE ", "SYSTEM ", "SHUTDOWN ",
	}
	for _, keyword := range dangerousKeywords {
		if strings.Contains(sqlUpper, keyword) {
			return fmt.Errorf("sql must not contain dangerous keywords, sql=%s", sql)
		}
	}
	return nil
}

func (a *AIClient) executeSQL(ctx context.Context, sql string, params []any) ([]map[string]any, error) {
	findings, err := a.findingRepo.ExecSQL(ctx, sql, params)
	if err != nil {
		return nil, err
	}
	return findings, nil
}
