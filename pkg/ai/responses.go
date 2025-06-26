package ai

import (
	"context"
	"fmt"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/openai/openai-go"
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

func (a *AIClient) responsesAPI(
	ctx context.Context,
	instruction string,
	inputs responses.ResponseNewParamsInputUnion,
	tools []responses.ToolUnionParam,
) (string, error) {
	resp, err := a.openaiClient.Responses.New(ctx,
		responses.ResponseNewParams{
			Model:        a.chatGPTModel,
			Instructions: openai.String(instruction),
			Input:        inputs,
			Tools:        tools,
		},
	)
	if err != nil {
		return "", fmt.Errorf("Responses API error: err=%w", err)
	}
	if resp.OutputText() == "" {
		return "", fmt.Errorf("Responses API: no response (instruction=%q, model=%q)", instruction, a.chatGPTModel)
	}
	a.logger.Infof(ctx, "Responses API Finished: %+v", resp.Usage)
	answer := resp.OutputText()
	return answer, nil
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
