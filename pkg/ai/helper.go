package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

func extractSingleJSONObject(output string) (string, error) {
	start := strings.Index(output, "{")
	end := strings.LastIndex(output, "}")
	if start == -1 || end == -1 || start >= end {
		return "", fmt.Errorf("no JSON object found")
	}
	jsonStr := output[start : end+1]
	var jsonObj any
	if err := json.Unmarshal([]byte(jsonStr), &jsonObj); err != nil {
		return "", fmt.Errorf("invalid JSON string: err=%w, jsonStr=%s", err, jsonStr)
	}
	return string(jsonStr), nil
}

// ConvertSchema
func ConvertSchema[T any](output string, schema T) (*T, error) {
	output = strings.TrimSpace(output)
	jsonStr, err := extractSingleJSONObject(output)
	if err != nil {
		return nil, fmt.Errorf("failed to extract JSON string: %w", err)
	}
	if err := json.Unmarshal([]byte(jsonStr), &schema); err != nil {
		return nil, err
	}
	return &schema, nil
}
