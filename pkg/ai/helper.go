package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ExtractJSONString(output string) (string, error) {
	start := strings.Index(output, "{")
	end := strings.LastIndex(output, "}")
	if start == -1 || end == -1 || start >= end {
		return "", fmt.Errorf("invalid JSON string")
	}
	return output[start : end+1], nil // return {...} string
}

// ConvertSchema
func ConvertSchema[T any](output string, schema T) (*T, error) {
	output = strings.TrimSpace(output)
	jsonStr, err := ExtractJSONString(output)
	if err != nil {
		return nil, fmt.Errorf("failed to extract JSON string: %w", err)
	}
	if err := json.Unmarshal([]byte(jsonStr), &schema); err != nil {
		return nil, err
	}
	return &schema, nil
}
