package ai

import (
	"context"
	"fmt"

	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/ai"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
)

const (
	REPORT_SYSTEM_PROMPT = `You are a security analyst expert. You have access to a SQL query executor tool that can analyze security findings data.

Use the sql_query_executor tool to:
1. Generate and execute SQL queries based on the user's request
2. Analyze the returned data
3. Create a comprehensive security report

Your report should include:
- Executive Summary
- Key Findings
- Data Analysis
- Security Recommendations
- Risk Assessment

Always explain your methodology and provide actionable insights.`
)

func (a *AIClient) GenerateReport(ctx context.Context, req *ai.GenerateReportRequest, findingRepo db.FindingRepository) (*ai.GenerateReportResponse, error) {
	// Use unified English system prompt
	systemPrompt := REPORT_SYSTEM_PROMPT

	// Create SQL execution tool
	sqlTool := NewSQLExecutionTool(findingRepo, req.ProjectId)

	// Create tools array with SQL executor
	tools := []responses.ToolUnionParam{
		GetSQLExecutorTool(),
	}

	// Prepare user message with context
	userMessage := fmt.Sprintf(`Project ID: %d
Analysis Request: %s

Please analyze the security findings data for this project and generate a comprehensive report based on the request above.`, req.ProjectId, req.Prompt)

	// Create input parameters
	inputParam := responses.ResponseInputParam{
		responses.ResponseInputItemUnionParam{
			OfMessage: &responses.EasyInputMessageParam{
				Role: responses.EasyInputMessageRoleUser,
				Content: responses.EasyInputMessageContentUnionParam{
					OfString: openai.String(userMessage),
				},
			},
		},
	}

	inputs := responses.ResponseNewParamsInputUnion{OfInputItemList: inputParam}

	// Generate report using AI with SQL tool handling
	report, err := a.responsesAPIWithSQLTool(ctx, systemPrompt, inputs, tools, sqlTool)
	if err != nil {
		return nil, fmt.Errorf("report generation failed: %w", err)
	}

	return &ai.GenerateReportResponse{
		Report: report,
	}, nil
}

func (a *AIClient) responsesAPIWithSQLTool(
	ctx context.Context,
	instruction string,
	inputs responses.ResponseNewParamsInputUnion,
	tools []responses.ToolUnionParam,
	sqlTool *SQLExecutionTool,
) (string, error) {
	// TODO: Implement proper tool handling when OpenAI Go SDK supports it
	// For now, the AI will use the SQL tool through function calling
	_ = sqlTool // Acknowledge the parameter to avoid linter warning
	return a.responsesAPI(ctx, instruction, inputs, tools)
}
