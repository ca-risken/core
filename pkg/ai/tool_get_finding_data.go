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
				"Get RISKEN findings in the current authorized project. " +
					"Use this when a request include \"finding\", \"issue\", \"ファインディング\", \"問題\"..." +
					"Only data that can be retrieved with a single SQL query per execution will be returned.",
			),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"prompt": map[string]any{
						"type":        "string",
						"description": "Prompt to get findings in the current authorized project. Generate SQL from natural language. e.g. \"Summarize aws findings by data_source. and sort by data_source name.\"",
					},
					"limit": map[string]any{
						"type":        "integer",
						"description": "Maximum number of results to return. (default: 20, max: 50)",
					},
					"offset": map[string]any{
						"type":        "integer",
						"description": "Offset of the results to return. (default: 0)",
					},
				},
				"required": []string{"prompt"},
			},
		},
	}
}

// GetFindingDataParams represents the parameters for get_finding_data function
type GetFindingDataParams struct {
	Prompt string `json:"prompt"`
	Limit  uint32 `json:"limit"`
	Offset uint32 `json:"offset"`
}

// GetFindingDataFunction handles the get_finding_data function call
func (a *AIClient) GetFindingDataFunction(ctx context.Context, projectID uint32, params GetFindingDataParams) ([]map[string]any, error) {
	// Default params
	if params.Limit == 0 {
		params.Limit = 20
	} else if params.Limit > 50 {
		params.Limit = 50
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
		results, err = a.getFindingDataFunction(ctx, projectID, params)
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

func (a *AIClient) getFindingDataFunction(ctx context.Context, projectID uint32, params GetFindingDataParams) ([]map[string]any, error) {
	sql, sqlParams, err := a.generateSQL(ctx, params.Prompt, projectID, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	a.logger.Infof(ctx, "Generated SQL: sql=%s, params=%v", sql, sqlParams)

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
- Do NOT include project_id condition in the WHERE clause(it will be automatically inserted from the authorized project when retrieving data).
- Must has prefix "SELECT" and include "WHERE" keywords.
- Use only a single direct FROM finding clause (optionally with an alias).
- Do NOT use JOIN, subqueries, UNION, or reference tables other than finding.
- Do NOT use dangerous keywords like "INSERT", "UPDATE", "DELETE" etc.
- Do NOT include semicolons (;) in the SQL(multiple statements are not allowed).
- Generate only ONE SELECT statement per request.

## Output Schema
{
	"sql": string
}

## Output Example
{"sql": "SELECT * FROM finding WHERE score >= 0.8 AND data_source like 'aws%' GROUP BY data_source ORDER BY data_source"}
`

type GenerateSQLOutput struct {
	SQL string `json:"sql"`
}

func (a *AIClient) generateSQL(ctx context.Context, prompt string, projectID, limit, offset uint32) (string, []any, error) {
	resp, err := a.callResponsesAPI(ctx,
		a.reasoningModel,
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
		}, DefaultTools, openai.ReasoningEffortMedium)
	if err != nil {
		return "", nil, err
	}
	output, err := ConvertSchema(resp.OutputText(), GenerateSQLOutput{})
	if err != nil {
		return "", nil, err
	}

	if err := validateSQL(output.SQL); err != nil {
		return "", nil, err
	}

	sql, params, err := formatSQL(output.SQL, projectID, limit, offset)
	if err != nil {
		return "", nil, err
	}
	return sql, params, nil
}

func formatSQL(sql string, projectID, limit, offset uint32) (string, []any, error) {
	if err := validateSQLSafety(sql); err != nil {
		return "", nil, err
	}
	parsed, err := parseFindingSQL(sql)
	if err != nil {
		return "", nil, err
	}

	qualifier := parsed.alias
	if qualifier == "" {
		qualifier = "finding"
	}
	pendFindingAlias := "pf_scope"
	if strings.EqualFold(qualifier, pendFindingAlias) {
		pendFindingAlias = "pf_scope_1"
	}
	scopeCondition := fmt.Sprintf(
		"%s.project_id = ? AND NOT EXISTS (SELECT 1 FROM pend_finding %s WHERE %s.project_id = %s.project_id AND %s.finding_id = %s.finding_id AND (%s.expired_at IS NULL OR %s.expired_at > NOW()))",
		qualifier,
		pendFindingAlias,
		pendFindingAlias,
		qualifier,
		pendFindingAlias,
		qualifier,
		pendFindingAlias,
		pendFindingAlias,
	)

	scopedSQL := fmt.Sprintf("%s WHERE %s AND (%s)", parsed.selectFromClause, scopeCondition, parsed.whereClause)
	if parsed.suffixClause != "" {
		scopedSQL = fmt.Sprintf("%s %s", scopedSQL, parsed.suffixClause)
	}

	params := []any{projectID, limit, offset}
	return fmt.Sprintf("SELECT * FROM (%s) as t LIMIT ? OFFSET ?", scopedSQL), params, nil
}

func validateSQL(sql string) error {
	if err := validateSQLSafety(sql); err != nil {
		return err
	}
	_, err := parseFindingSQL(sql)
	return err
}

type parsedFindingSQL struct {
	selectFromClause string
	whereClause      string
	suffixClause     string
	alias            string
}

type sqlWord struct {
	text  string
	start int
	end   int
	depth int
}

func parseFindingSQL(sql string) (*parsedFindingSQL, error) {
	trimmedSQL := strings.TrimSpace(sql)
	if trimmedSQL == "" {
		return nil, fmt.Errorf("sql must not be empty")
	}

	words, err := scanSQLWords(trimmedSQL)
	if err != nil {
		return nil, err
	}
	if len(words) == 0 || words[0].depth != 0 || words[0].text != "SELECT" || words[0].start != 0 {
		return nil, fmt.Errorf("sql must start with SELECT")
	}

	var fromWord, whereWord *sqlWord
	for i := range words {
		word := words[i]
		if word.text == "SELECT" {
			if word.depth > 0 {
				return nil, fmt.Errorf("subqueries are not allowed")
			}
			if i > 0 {
				return nil, fmt.Errorf("multiple SELECT statements are not allowed")
			}
		}
		if word.depth != 0 {
			continue
		}
		switch word.text {
		case "UNION", "INTERSECT", "EXCEPT":
			return nil, fmt.Errorf("set operators are not allowed")
		case "JOIN":
			return nil, fmt.Errorf("JOIN is not allowed")
		case "FROM":
			if fromWord != nil {
				return nil, fmt.Errorf("multiple FROM clauses are not allowed")
			}
			fromWord = &words[i]
		case "WHERE":
			if whereWord != nil {
				return nil, fmt.Errorf("multiple WHERE clauses are not allowed")
			}
			whereWord = &words[i]
		}
	}
	if fromWord == nil || whereWord == nil || fromWord.start >= whereWord.start {
		return nil, fmt.Errorf("sql must contain SELECT, FROM finding, and WHERE")
	}

	fromSegment := strings.TrimSpace(trimmedSQL[fromWord.end:whereWord.start])
	alias, err := parseFindingSource(fromSegment)
	if err != nil {
		return nil, err
	}

	suffixStart := len(trimmedSQL)
	for _, word := range words {
		if word.depth != 0 || word.start <= whereWord.end {
			continue
		}
		if isSQLSuffixKeyword(word.text) {
			suffixStart = word.start
			break
		}
	}

	whereClause := strings.TrimSpace(trimmedSQL[whereWord.end:suffixStart])
	if whereClause == "" {
		return nil, fmt.Errorf("WHERE clause is required")
	}

	suffixClause := strings.TrimSpace(trimmedSQL[suffixStart:])
	return &parsedFindingSQL{
		selectFromClause: strings.TrimSpace(trimmedSQL[:whereWord.start]),
		whereClause:      whereClause,
		suffixClause:     suffixClause,
		alias:            alias,
	}, nil
}

func validateSQLSafety(sql string) error {
	sqlUpper := strings.ToUpper(strings.TrimSpace(sql))
	if sqlUpper == "" {
		return fmt.Errorf("sql must not be empty")
	}
	if strings.Contains(sql, ";") {
		return fmt.Errorf("sql must not contain semicolons (multiple statements not allowed)")
	}

	dangerousKeywords := []string{
		"INSERT ", "UPDATE ", "DELETE ", "DROP TABLE ", "CREATE TABLE ", "ALTER TABLE ",
		"TRUNCATE ", "EXEC ", "EXECUTE ", "INTO OUTFILE ",
		"LOAD_FILE ", "SYSTEM ", "SHUTDOWN ",
	}
	for _, keyword := range dangerousKeywords {
		if strings.Contains(sqlUpper, keyword) {
			return fmt.Errorf("sql must not contain dangerous keywords")
		}
	}
	return nil
}

func parseFindingSource(source string) (string, error) {
	if source == "" {
		return "", fmt.Errorf("FROM finding clause is required")
	}
	if strings.ContainsAny(source, ",()") {
		return "", fmt.Errorf("only direct finding table references are allowed")
	}

	tokens := strings.Fields(source)
	switch len(tokens) {
	case 1:
		if !isFindingIdentifier(tokens[0]) {
			return "", fmt.Errorf("only finding table is allowed")
		}
		return "", nil
	case 2:
		if !isFindingIdentifier(tokens[0]) {
			return "", fmt.Errorf("only finding table is allowed")
		}
		alias := normalizeIdentifier(tokens[1])
		if !isSafeIdentifier(alias) {
			return "", fmt.Errorf("invalid finding alias")
		}
		return alias, nil
	case 3:
		if !isFindingIdentifier(tokens[0]) || !strings.EqualFold(tokens[1], "AS") {
			return "", fmt.Errorf("only finding table is allowed")
		}
		alias := normalizeIdentifier(tokens[2])
		if !isSafeIdentifier(alias) {
			return "", fmt.Errorf("invalid finding alias")
		}
		return alias, nil
	default:
		return "", fmt.Errorf("only finding table is allowed")
	}
}

func scanSQLWords(sql string) ([]sqlWord, error) {
	words := make([]sqlWord, 0, 16)
	depth := 0
	var quote byte

	for i := 0; i < len(sql); i++ {
		ch := sql[i]
		if quote != 0 {
			if ch == quote {
				if quote != '`' && i+1 < len(sql) && sql[i+1] == ch {
					i++
					continue
				}
				if i > 0 && sql[i-1] == '\\' {
					continue
				}
				quote = 0
			}
			continue
		}

		switch ch {
		case '\'', '"', '`':
			quote = ch
		case '(':
			depth++
		case ')':
			if depth == 0 {
				return nil, fmt.Errorf("sql has unbalanced parentheses")
			}
			depth--
		default:
			if !isSQLWordChar(ch) {
				continue
			}
			start := i
			for i+1 < len(sql) && isSQLWordChar(sql[i+1]) {
				i++
			}
			words = append(words, sqlWord{
				text:  strings.ToUpper(sql[start : i+1]),
				start: start,
				end:   i + 1,
				depth: depth,
			})
		}
	}

	if quote != 0 || depth != 0 {
		return nil, fmt.Errorf("sql has unbalanced quotes or parentheses")
	}
	return words, nil
}

func isSQLSuffixKeyword(word string) bool {
	switch word {
	case "GROUP", "ORDER", "HAVING", "LIMIT", "OFFSET":
		return true
	default:
		return false
	}
}

func isFindingIdentifier(identifier string) bool {
	return strings.EqualFold(normalizeIdentifier(identifier), "finding")
}

func normalizeIdentifier(identifier string) string {
	return strings.Trim(identifier, "`")
}

func isSafeIdentifier(identifier string) bool {
	if identifier == "" {
		return false
	}
	for i := 0; i < len(identifier); i++ {
		ch := identifier[i]
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' {
			continue
		}
		if i > 0 && ch >= '0' && ch <= '9' {
			continue
		}
		return false
	}
	return true
}

func isSQLWordChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '_'
}

func (a *AIClient) executeSQL(ctx context.Context, sql string, params []any) ([]map[string]any, error) {
	findings, err := a.findingRepo.ExecSQL(ctx, sql, params)
	if err != nil {
		return nil, err
	}
	return findings, nil
}
