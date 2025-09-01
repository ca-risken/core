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

## Scope
- Project ID: %d

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
When using xychart-beta with Japanese text, you MUST use double quotes for the title.
You SHOULD customize the chart colors using the config block with header configuration, as the default colors may not provide sufficient contrast or visual appeal.

---
config:
    themeVariables:
        xyChart:
            plotColorPalette: "teal, pink, cyan"
---
xychart-beta
    title "日本の年間出生数の推移"
    x-axis ["1950年", "1980年", "2000年", "2020年", "2022年"]
    y-axis "出生数（万人）" 0 --> 250
    bar [233.7, 157.7, 119.9, 84.0, 77.0]
    line [233.7, 157.7, 119.9, 84.0, 77.0]

---
config:
    themeVariables:
        xyChart:
            plotColorPalette: "grey, indigo, pink, purple, cyan, teal, green"
---
xychart-beta
    title "Different colors in xyChart"
    x-axis "numbersX" ["Group A", "Group B", "Group C", "Group D"]
    y-axis "numbersY" 0 --> 40
    bar [20,30,25,30]
    bar [10,20,20,20]
    bar [5,10,5,5]
    line [0,5,40,2]

#### quadrantChart
When using quadrantChart with Japanese text, you MUST use double quotes for the x-axis. (But **not for the title**)

quadrantChart
    title 製品評価マトリクス
    x-axis "コスト" --> "効果"
    y-axis "実現性（低）" --> "実現性（高）"
    quadrant-1 "優先度高"
    quadrant-2 "要検討"
    quadrant-3 "保留"
    quadrant-4 "簡単な改善"
    "製品A": [0.3, 0.6]
    "製品B": [0.45, 0.23]
    "製品C": [0.57, 0.69]
    "製品D": [0.78, 0.34]
    "製品E": [0.40, 0.34]

#### Example) flowchart
When using flowchart with Japanese text, you MUST NOT use double quotes for the label.
But **don't use <br />** in the label.(Use <br> instead.)

flowchart LR
    A[実装<br>オプション]
    B[AWS<br>（Hoge Service）]
    C[Google Cloud<br>（Fuga Service）]
    D[Azure<br>（Piyo Service）]
    A --> B
    A --> C
    A --> D
    B -.- B1[大規模開発<br>実績多数]
    C -.- C1[低コスト<br>先進技術]
    D -.- D1[セキュリティ<br>]

    style B fill:#f9f,stroke:#333
    style C fill:#bbf,stroke:#333
    style D fill:#bfb,stroke:#333

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
When using pie with Japanese text, you MUST use double quotes for the labels (But **not for the title**).

pie
    title 2023年度 部門別売上構成比
    "製品A部門" : 35.7
    "製品B部門" : 24.3
    "サービス部門" : 20.5
    "コンサルティング" : 15.2
    "その他" : 4.3

pie
    title Finding category
    "AWS" : 35
    "Google Cloud" : 15
    "Code" : 10
    "Other" : 40

#### Example) mindmap
When using mindmap with Japanese text, you MUST NOT use double quotes for the label.
But if you use parentheses, you MUST use full-width parentheses.

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
