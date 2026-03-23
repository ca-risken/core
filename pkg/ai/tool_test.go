package ai

import (
	"errors"
	"reflect"
	"testing"

	"github.com/openai/openai-go/responses"
)

type getFindingDataToolSchema struct {
	name           string
	hasProjectID   bool
	requiredFields []string
}

type chatAIToolsResult struct {
	names []string
}

var (
	errToolFunctionNil       = errors.New("function tool is nil")
	errToolPropertiesInvalid = errors.New("properties is not defined")
	errToolRequiredInvalid   = errors.New("required is not defined")
)

func TestGetFindingDataToolSchema(t *testing.T) {
	cases := []struct {
		name string
		want getFindingDataToolSchema
	}{
		{
			name: "schema is restricted to authorized project context",
			want: getFindingDataToolSchema{
				name:           "get_finding_data",
				hasProjectID:   false,
				requiredFields: []string{"prompt"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := getFindingDataToolSchemaResult()
			if err != nil {
				t.Fatalf("getFindingDataToolSchemaResult() error = %v", err)
			}
			if !reflect.DeepEqual(c.want, *got) {
				t.Fatalf("getFindingDataToolSchemaResult() = %+v, want %+v", got, c.want)
			}
		})
	}
}

func TestBuildChatAITools(t *testing.T) {
	cases := []struct {
		name      string
		projectID uint32
		want      chatAIToolsResult
	}{
		{
			name:      "without project context",
			projectID: 0,
			want: chatAIToolsResult{
				names: []string{string(responses.WebSearchToolTypeWebSearchPreview)},
			},
		},
		{
			name:      "with project context",
			projectID: 1001,
			want: chatAIToolsResult{
				names: []string{string(responses.WebSearchToolTypeWebSearchPreview), "get_finding_data"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getChatAIToolsResult(c.projectID)
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("getChatAIToolsResult() = %+v, want %+v", got, c.want)
			}
		})
	}
}

func getFindingDataToolSchemaResult() (*getFindingDataToolSchema, error) {
	tool := GetFindingDataTool()
	if tool.OfFunction == nil {
		return nil, errToolFunctionNil
	}

	properties, ok := tool.OfFunction.Parameters["properties"].(map[string]any)
	if !ok {
		return nil, errToolPropertiesInvalid
	}

	required, ok := tool.OfFunction.Parameters["required"].([]string)
	if !ok {
		return nil, errToolRequiredInvalid
	}

	_, hasProjectID := properties["project_id"]
	return &getFindingDataToolSchema{
		name:           tool.OfFunction.Name,
		hasProjectID:   hasProjectID,
		requiredFields: required,
	}, nil
}

func getChatAIToolsResult(projectID uint32) chatAIToolsResult {
	tools := buildChatAITools(projectID)
	names := make([]string, 0, len(tools))
	for _, tool := range tools {
		switch {
		case tool.OfFunction != nil:
			names = append(names, tool.OfFunction.Name)
		case tool.OfWebSearchPreview != nil:
			names = append(names, string(tool.OfWebSearchPreview.Type))
		}
	}
	return chatAIToolsResult{names: names}
}
