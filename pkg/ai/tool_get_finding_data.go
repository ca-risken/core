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
- Do NOT include project_id condition in the WHERE clause(it will be automatically inserted from the authorized project when retrieving data).
- Must has prefix "SELECT" and include "WHERE" keywords.
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

	sql, params, err := formatSQL(output.SQL, projectID, limit, offset)
	if err != nil {
		return "", nil, err
	}
	if err := validateSQL(sql); err != nil {
		return "", nil, err
	}
	return sql, params, nil
}

func formatSQL(sql string, projectID, limit, offset uint32) (string, []any, error) {
	parsed, err := parseFindingSelectSQL(sql)
	if err != nil {
		return "", nil, err
	}

	// strict project_id filter & ignore pend_findings
	scopedSQL := fmt.Sprintf(`%s WHERE
	project_id = ?
	AND not exists (
		SELECT 1
		FROM pend_finding
		WHERE
		  pend_finding.finding_id = finding.finding_id
			and (pend_finding.expired_at is NULL or pend_finding.expired_at > NOW())
	)
	AND (%s)`, parsed.head, parsed.whereCondition)
	if parsed.tail != "" {
		scopedSQL = fmt.Sprintf("%s %s", scopedSQL, parsed.tail)
	}

	// add limit and offset
	formattedSQL := fmt.Sprintf(`SELECT * FROM (%s) as t LIMIT ? OFFSET ?`, scopedSQL)
	params := []any{projectID, limit, offset}
	return formattedSQL, params, nil
}

type sqlToken struct {
	word  string
	start int
	end   int
	depth int
}

type findingSelectQuery struct {
	head           string
	whereCondition string
	tail           string
}

func parseFindingSelectSQL(sql string) (*findingSelectQuery, error) {
	trimmedSQL := strings.TrimSpace(sql)
	if trimmedSQL == "" {
		return nil, fmt.Errorf("sql must not be empty")
	}
	if strings.Contains(trimmedSQL, ";") {
		return nil, fmt.Errorf("sql must not contain semicolons (multiple statements not allowed), sql=%s", trimmedSQL)
	}

	tokens, hasSubQuery, err := scanSQLTokens(trimmedSQL)
	if err != nil {
		return nil, err
	}
	if hasSubQuery {
		return nil, fmt.Errorf("sql must be single-level select without subquery, sql=%s", trimmedSQL)
	}
	if len(tokens) == 0 || tokens[0].word != "SELECT" {
		return nil, fmt.Errorf("sql must start with SELECT, sql=%s", trimmedSQL)
	}

	fromIndex := -1
	whereIndexes := make([]int, 0, 2)
	for i, token := range tokens {
		switch token.word {
		case "FROM":
			if fromIndex == -1 {
				fromIndex = i
			}
		case "WHERE":
			whereIndexes = append(whereIndexes, i)
		case "UNION", "INTERSECT", "EXCEPT":
			return nil, fmt.Errorf("set operators are not allowed, sql=%s", trimmedSQL)
		}
	}
	if fromIndex == -1 {
		return nil, fmt.Errorf("sql must contain FROM clause, sql=%s", trimmedSQL)
	}
	if len(whereIndexes) != 1 {
		return nil, fmt.Errorf("sql must contain exactly one top-level WHERE clause, sql=%s", trimmedSQL)
	}
	whereIndex := whereIndexes[0]
	if whereIndex <= fromIndex {
		return nil, fmt.Errorf("sql must contain WHERE after FROM, sql=%s", trimmedSQL)
	}

	fromClause := strings.TrimSpace(trimmedSQL[tokens[fromIndex].end:tokens[whereIndex].start])
	if err := validateFromClause(fromClause); err != nil {
		return nil, fmt.Errorf("invalid FROM clause: %w", err)
	}

	tailStart := len(trimmedSQL)
	for i := whereIndex + 1; i < len(tokens); i++ {
		switch tokens[i].word {
		case "GROUP":
			if i+1 < len(tokens) && tokens[i+1].word == "BY" {
				tailStart = tokens[i].start
			}
		case "ORDER":
			if i+1 < len(tokens) && tokens[i+1].word == "BY" {
				tailStart = tokens[i].start
			}
		case "HAVING", "LIMIT":
			tailStart = tokens[i].start
		}
		if tailStart != len(trimmedSQL) {
			break
		}
	}
	whereCondition := strings.TrimSpace(trimmedSQL[tokens[whereIndex].end:tailStart])
	if whereCondition == "" {
		return nil, fmt.Errorf("sql must contain non-empty WHERE condition, sql=%s", trimmedSQL)
	}

	return &findingSelectQuery{
		head:           strings.TrimSpace(trimmedSQL[:tokens[whereIndex].start]),
		whereCondition: whereCondition,
		tail:           strings.TrimSpace(trimmedSQL[tailStart:]),
	}, nil
}

func scanSQLTokens(sql string) ([]sqlToken, bool, error) {
	tokens := make([]sqlToken, 0, 16)
	depth := 0
	hasSubQuery := false
	inSingleQuote := false
	inDoubleQuote := false
	inBacktick := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(sql); {
		ch := sql[i]

		if inLineComment {
			if ch == '\n' {
				inLineComment = false
			}
			i++
			continue
		}
		if inBlockComment {
			if ch == '*' && i+1 < len(sql) && sql[i+1] == '/' {
				inBlockComment = false
				i += 2
				continue
			}
			i++
			continue
		}
		if inSingleQuote {
			if ch == '\\' && i+1 < len(sql) {
				i += 2
				continue
			}
			if ch == '\'' {
				if i+1 < len(sql) && sql[i+1] == '\'' {
					i += 2
					continue
				}
				inSingleQuote = false
			}
			i++
			continue
		}
		if inDoubleQuote {
			if ch == '\\' && i+1 < len(sql) {
				i += 2
				continue
			}
			if ch == '"' {
				if i+1 < len(sql) && sql[i+1] == '"' {
					i += 2
					continue
				}
				inDoubleQuote = false
			}
			i++
			continue
		}
		if inBacktick {
			if ch == '`' {
				inBacktick = false
			}
			i++
			continue
		}

		if ch == '-' && i+1 < len(sql) && sql[i+1] == '-' && (i+2 == len(sql) || isSQLWhitespace(sql[i+2])) {
			inLineComment = true
			i += 2
			continue
		}
		if ch == '#' {
			inLineComment = true
			i++
			continue
		}
		if ch == '/' && i+1 < len(sql) && sql[i+1] == '*' {
			inBlockComment = true
			i += 2
			continue
		}

		switch ch {
		case '\'':
			inSingleQuote = true
			i++
			continue
		case '"':
			inDoubleQuote = true
			i++
			continue
		case '`':
			inBacktick = true
			i++
			continue
		case '(':
			depth++
			i++
			continue
		case ')':
			if depth == 0 {
				return nil, false, fmt.Errorf("sql has unbalanced parentheses, sql=%s", sql)
			}
			depth--
			i++
			continue
		}

		if isSQLWordStart(ch) {
			start := i
			i++
			for i < len(sql) && isSQLWordPart(sql[i]) {
				i++
			}
			word := strings.ToUpper(sql[start:i])
			if word == "SELECT" && depth > 0 {
				hasSubQuery = true
			}
			if depth == 0 {
				tokens = append(tokens, sqlToken{
					word:  word,
					start: start,
					end:   i,
					depth: depth,
				})
			}
			continue
		}

		i++
	}

	if depth != 0 || inSingleQuote || inDoubleQuote || inBacktick || inBlockComment {
		return nil, false, fmt.Errorf("sql has unclosed expression, sql=%s", sql)
	}
	return tokens, hasSubQuery, nil
}

func validateFromClause(fromClause string) error {
	trimmed := strings.TrimSpace(fromClause)
	if trimmed == "" {
		return fmt.Errorf("from clause must not be empty")
	}
	if strings.ContainsAny(trimmed, "()") || strings.Contains(trimmed, ",") {
		return fmt.Errorf("from clause must target single finding table")
	}

	fields := strings.Fields(trimmed)
	if len(fields) == 0 {
		return fmt.Errorf("from clause must not be empty")
	}
	switch len(fields) {
	case 1:
		// finding
	case 2:
		// finding alias
		if strings.EqualFold(fields[1], "AS") {
			return fmt.Errorf("alias is required after AS: %s", trimmed)
		}
		if !isSQLIdentifier(fields[1]) {
			return fmt.Errorf("invalid table alias: %s", fields[1])
		}
	case 3:
		// finding AS alias
		if !strings.EqualFold(fields[1], "AS") || !isSQLIdentifier(fields[2]) {
			return fmt.Errorf("invalid table alias expression: %s", trimmed)
		}
	default:
		return fmt.Errorf("from clause must target single finding table: %s", trimmed)
	}

	tableName := trimSQLIdentifier(fields[0])
	if !strings.EqualFold(tableName, "finding") {
		return fmt.Errorf("from clause must target finding table: %s", trimmed)
	}
	return nil
}

func trimSQLIdentifier(identifier string) string {
	trimmed := strings.TrimSpace(identifier)
	if len(trimmed) >= 2 && trimmed[0] == '`' && trimmed[len(trimmed)-1] == '`' {
		return trimmed[1 : len(trimmed)-1]
	}
	return trimmed
}

func isSQLIdentifier(value string) bool {
	trimmed := trimSQLIdentifier(value)
	if trimmed == "" {
		return false
	}
	if !isSQLWordStart(trimmed[0]) {
		return false
	}
	for i := 1; i < len(trimmed); i++ {
		if !isSQLWordPart(trimmed[i]) {
			return false
		}
	}
	return true
}

func isSQLWordStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isSQLWordPart(ch byte) bool {
	return isSQLWordStart(ch) || (ch >= '0' && ch <= '9')
}

func isSQLWhitespace(ch byte) bool {
	switch ch {
	case ' ', '\t', '\n', '\r', '\f':
		return true
	default:
		return false
	}
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
