package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/responses"
)

var DefaultTools = []responses.ToolUnionParam{
	{
		OfWebSearchPreview: &responses.WebSearchToolParam{
			Type:              responses.WebSearchToolTypeWebSearchPreview,
			SearchContextSize: responses.WebSearchToolSearchContextSizeMedium,
		},
	},
}

// FunctionCallInfo represents a function call extracted from response
type FunctionCallInfo struct {
	Name      string
	Arguments string
	CallID    string
}

// extractFunctionCalls extracts all function calls from response output
func extractFunctionCalls(output []responses.ResponseOutputItemUnion) []FunctionCallInfo {
	var functionCalls []FunctionCallInfo
	for _, outputItem := range output {
		if outputItem.Type == "function_call" {
			functionCall := outputItem.AsFunctionCall()
			functionCalls = append(functionCalls, FunctionCallInfo{
				Name:      functionCall.Name,
				Arguments: functionCall.Arguments,
				CallID:    functionCall.CallID,
			})
		}
	}
	return functionCalls
}

// handleFunctionCalls executes function calls and updates conversation history
func (a *AIClient) handleFunctionCalls(
	ctx context.Context,
	currentInputs responses.ResponseNewParamsInputUnion,
	functionCalls []FunctionCallInfo,
) (responses.ResponseNewParamsInputUnion, error) {
	// Get current input items
	var inputItems responses.ResponseInputParam
	if currentInputs.OfInputItemList != nil {
		inputItems = currentInputs.OfInputItemList
	}

	// Process each function call: add request, execute, then add output
	for _, funcCall := range functionCalls {
		functionCallRequest := responses.ResponseInputItemParamOfFunctionCall(
			funcCall.Arguments,
			funcCall.CallID,
			funcCall.Name,
		)
		inputItems = append(inputItems, functionCallRequest)
		a.logger.Infof(ctx, "Executing function: call_id=%s, name=%s, args=%s", funcCall.CallID, funcCall.Name, funcCall.Arguments)

		// Execute the function
		result, err := executeFunctionCall(funcCall.Name, funcCall.Arguments)
		if err != nil {
			return responses.ResponseNewParamsInputUnion{}, fmt.Errorf("function call execution failed: %w", err)
		}
		functionCallOutput := responses.ResponseInputItemParamOfFunctionCallOutput(
			funcCall.CallID,
			result,
		)
		inputItems = append(inputItems, functionCallOutput)
		a.logger.Debugf(ctx, "Function call result: call_id=%s, result=%s", funcCall.CallID, result)
	}

	return responses.ResponseNewParamsInputUnion{OfInputItemList: inputItems}, nil
}

// executeFunctionCall executes the specified function call
func executeFunctionCall(functionName string, arguments string) (string, error) {
	switch functionName {
	case "get_findings":
		var params GetFindingsParams
		if err := json.Unmarshal([]byte(arguments), &params); err != nil {
			return "", fmt.Errorf("failed to parse get_findings parameters: %w", err)
		}

		findings, err := GetFindingsFunction(params)
		if err != nil {
			return "", fmt.Errorf("failed to get findings: %w", err)
		}

		result, err := json.Marshal(findings)
		if err != nil {
			return "", fmt.Errorf("failed to marshal findings result: %w", err)
		}

		return string(result), nil
	default:
		return "", fmt.Errorf("unknown function: %s", functionName)
	}
}
