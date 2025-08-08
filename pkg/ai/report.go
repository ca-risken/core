package ai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
)

func (a *AIClient) GenerateReport(ctx context.Context, projectID uint32, prompt string) (string, error) {
	instruction := generatePrompt(projectID)
	tools := DefaultTools
	tools = append(tools, GetFindingDataTool())

	resp, err := a.responsesAPI(ctx, a.reasoningModel, instruction, responses.ResponseNewParamsInputUnion{
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
	}, tools, openai.ReasoningEffortHigh)
	if err != nil {
		return "", err
	}
	return resp.OutputText(), nil
}

const (
	REPORT_PROMPT = `
Please generate RISKEN Finding report for the project.
The report language should match the user's language preference.

## Report details
- Project ID: %d
- Min score: 0.4

## Strict Rules
- NO HALLUCINATION: All content MUST be based on actual data and facts only
- DATA VALIDATION: Every insight MUST be backed by concrete evidence and verifiable data sources

## RISKEN Scoring System Reference
RISKEN implements a scoring system where severity levels are determined by the following score ranges:
- HIGH: 1.0 ~ 0.8 (Requires immediate review and response)
- MEDIUM: 0.7 ~ 0.4 (Not urgent but should be reviewed)
- LOW: 0.3 ~ 0.0 (INFO level, no review or response required)

## Output format
You MUST output in Markdown format.
Also, please include the following sections:

- Report conditions and periods
- Executive Summary
- Sections
  - Section name (Categorize the analysis results and summarize the report for each section.)
  - Graph (Generate a graph using mermaid format)
  - Fact (Organize the facts and list specific data if available)
  - Insight (Pattern analysis, predictions, data trends, etc.)
  - Example Finding (List example findings) (optional)
  - Recommendation (Improvement measures, countermeasures) (optional)

### Graph examples

When inserting mermaid graphs, you MUST specify mermaid in the code block.

#### Example) table (not mermaid)

| Category | Score | Count |
|----------|-------|-------|
| AWS      | 0.8   | 10    |
| Google Cloud | 0.5 | 5   |
| Code     | 0.3   | 3     |
| Other    | 0.2   | 2     |

#### Example) xychart-beta

xychart-beta
    title "AWS Findings count by month"
    x-axis ["Jan", "Feb", "Mar", "Apr", "May"]
    y-axis "Finding count" 0 --> 50
    bar [10, 20, 30, 40, 50]
    line [10, 20, 30, 40, 50]


#### Example) flowchart

flowchart LR
    A[Implementation<br>Option]
    B[AWS<br>(Hoge Service)]
    C[Google Cloud<br>(Fuga Service)]
    D[Azure<br>(Piyo Service)]
    A --> B
    A --> C
    A --> D
    B -.- B1[Large-scale development<br>Many achievements]
    C -.- C1[Low cost<br>Advanced technology]
    D -.- D1[Security<br>]
    style B fill:#f9f,stroke:#333
    style C fill:#bbf,stroke:#333
    style D fill:#bfb,stroke:#333

#### Example) pie

pie
    title "Finding category"
    "AWS" : 35
    "Google Cloud" : 15
    "Code" : 10
    "Other" : 40

#### Example) mindmap

mindmap
    root("Center")
        Pattern A
            A-1 (Full-width parentheses)
            A-2
        Pattern B
            B-1
            B-2
            B-3
`
)

func generatePrompt(projectID uint32) string {
	return fmt.Sprintf(REPORT_PROMPT, projectID)
}
