package ai

import (
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
)

// GetFindingsTool returns the get_findings function tool definition
func GetFindingsTool() responses.ToolUnionParam {
	return responses.ToolUnionParam{
		OfFunction: &responses.FunctionToolParam{
			Name: "get_findings",
			Description: openai.String(
				"Get RISKEN findings in a project. " +
					"Use this when a request include \"finding\", \"issue\", \"ファインディング\", \"問題\"..."),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"project_id": map[string]any{
						"type":        "integer",
						"description": "Project ID to get findings for",
					},
				},
				"required": []string{"project_id"},
			},
		},
	}
}

// TODO: delete
type Finding struct {
	ID          uint32    `json:"id"`
	Description string    `json:"description"`
	DataSource  string    `json:"data_source"`
	ResourceID  string    `json:"resource_id"`
	ProjectID   uint32    `json:"project_id"`
	Score       float32   `json:"score"`
	Data        string    `json:"data"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetFindingsParams represents the parameters for get_findings function
type GetFindingsParams struct {
	ProjectID uint32 `json:"project_id"`
}

// GetFindingsFunction handles the get_findings function call
func GetFindingsFunction(params GetFindingsParams) ([]Finding, error) {
	// TODO: implement
	findings := []Finding{
		{
			ID:          1,
			Description: "Suspicious network activity detected",
			DataSource:  "aws-guard-duty",
			ResourceID:  "i-1234567890abcdef0",
			ProjectID:   params.ProjectID,
			Score:       8.5,
			Data:        `{"source_ip": "192.168.1.100", "destination_ip": "10.0.0.1", "port": 443}`,
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			UpdatedAt:   time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          2,
			Description: "Unauthorized access attempt",
			DataSource:  "aws-cloud-trail",
			ResourceID:  "arn:aws:s3:::my-bucket",
			ProjectID:   params.ProjectID,
			Score:       7.2,
			Data:        `{"user": "unknown", "action": "GetObject", "result": "Failed"}`,
			CreatedAt:   time.Now().Add(-12 * time.Hour),
			UpdatedAt:   time.Now().Add(-12 * time.Hour),
		},
		{
			ID:          3,
			Description: "Potential SQL injection attempt",
			DataSource:  "web-application-firewall",
			ResourceID:  "arn:aws:waf::123456789012:webacl/example-web-acl",
			ProjectID:   params.ProjectID,
			Score:       9.1,
			Data:        `{"request_uri": "/api/users", "payload": "' OR 1=1 --", "blocked": true}`,
			CreatedAt:   time.Now().Add(-6 * time.Hour),
			UpdatedAt:   time.Now().Add(-6 * time.Hour),
		},
		{
			ID:          4,
			Description: "High CPU usage detected",
			DataSource:  "aws-cloudwatch",
			ResourceID:  "i-0987654321fedcba0",
			ProjectID:   params.ProjectID,
			Score:       6.8,
			Data:        `{"cpu_utilization": 95.2, "duration": "30min", "threshold": 80}`,
			CreatedAt:   time.Now().Add(-3 * time.Hour),
			UpdatedAt:   time.Now().Add(-3 * time.Hour),
		},
	}

	return findings, nil
}
