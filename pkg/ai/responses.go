package ai

import (
	"context"
	"fmt"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
)

const MAX_TOOL_USE_COUNT = 30

func (a *AIClient) responsesAPI(
	ctx context.Context,
	model string,
	instruction string,
	inputs responses.ResponseNewParamsInputUnion,
	tools []responses.ToolUnionParam,
) (*responses.Response, error) {
	currentInputs := inputs
	for range MAX_TOOL_USE_COUNT {
		resp, err := a.callResponsesAPI(ctx, model, instruction, currentInputs, tools)
		if err != nil {
			return nil, err
		}
		functionCalls := extractFunctionCalls(resp.Output)
		if len(functionCalls) == 0 {
			// No function calls found, return the final response
			a.logger.Infof(ctx, "Responses API Finished: %+v", resp.Usage)
			return resp, nil
		}
		currentInputs, err = a.handleFunctionCalls(ctx, currentInputs, functionCalls)
		if err != nil {
			return nil, fmt.Errorf("function call handling error: %w", err)
		}
	}
	return nil, fmt.Errorf("maximum function call iterations (%d) exceeded", MAX_TOOL_USE_COUNT)
}

func (a *AIClient) callResponsesAPI(
	ctx context.Context,
	model string,
	instruction string,
	inputs responses.ResponseNewParamsInputUnion,
	tools []responses.ToolUnionParam,
) (*responses.Response, error) {
	resp, err := a.openaiClient.Responses.New(ctx,
		responses.ResponseNewParams{
			Model:           model,
			Instructions:    openai.String(instruction),
			Input:           inputs,
			Tools:           tools,
			MaxToolCalls:    openai.Int(MAX_TOOL_USE_COUNT),
			MaxOutputTokens: openai.Int(25000),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Responses API error: err=%w", err)
	}
	return resp, nil
}

// StreamSender is a generic interface for sending data in a stream
type StreamSender[T any] interface {
	Send(T) error
}

func responsesStreamingAPI[T any](
	ctx context.Context,
	openaiClient *openai.Client,
	model string,
	instruction string,
	inputs responses.ResponseNewParamsInputUnion,
	tools []responses.ToolUnionParam,
	stream StreamSender[T],
	sendContentFunc func(string) T,
	logger logging.Logger,
) (string, error) {
	streamResp := openaiClient.Responses.NewStreaming(ctx, responses.ResponseNewParams{
		Model:        model,
		Instructions: openai.String(instruction),
		Input:        inputs,
		Tools:        tools,
	},
	)
	defer streamResp.Close()

	var answer string
	for streamResp.Next() {
		data := streamResp.Current()
		if data.Delta.OfString != "" {
			response := sendContentFunc(data.Delta.OfString)
			if sendErr := stream.Send(response); sendErr != nil {
				return "", sendErr
			}
			if data.JSON.Text.Valid() {
				// finish
				answer = data.Text
				break
			}
		}
	}
	if err := streamResp.Err(); err != nil {
		return "", fmt.Errorf("stream error: err=%w", err)
	}

	logger.Infof(ctx, "Responses API(streaming) Finished: type=%T", stream)
	return answer, nil
}
