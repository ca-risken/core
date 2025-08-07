package ai

import (
	"context"
	"fmt"

	"github.com/ca-risken/core/proto/ai"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
)

func (a *AIClient) ChatAI(ctx context.Context, req *ai.ChatAIRequest) (*ai.ChatAIResponse, error) {
	instruction := "You are an AI chat assistant. Please respond to user questions politely and concisely." +
		"Please generate appropriate responses according to the user's language."
	inputParam := responses.ResponseInputParam{}
	for _, h := range req.ChatHistory {
		role := responses.EasyInputMessageRoleUser
		if h.Role == *ai.ChatRole_CHAT_ROLE_ASSISTANT.Enum() {
			role = responses.EasyInputMessageRoleAssistant
		}
		inputParam = append(inputParam, responses.ResponseInputItemUnionParam{
			OfMessage: &responses.EasyInputMessageParam{
				Role: role,
				Content: responses.EasyInputMessageContentUnionParam{
					OfString: openai.String(h.Content),
				},
			},
		})
	}
	inputParam = append(inputParam, responses.ResponseInputItemUnionParam{
		OfMessage: &responses.EasyInputMessageParam{
			Role: responses.EasyInputMessageRoleUser,
			Content: responses.EasyInputMessageContentUnionParam{
				OfString: openai.String(req.Question),
			},
		},
	})
	inputs := responses.ResponseNewParamsInputUnion{OfInputItemList: inputParam}
	tools := DefaultTools
	tools = append(tools, GetFindingDataTool())
	answer, err := a.responsesAPI(ctx, a.chatGPTModel, instruction, inputs, tools)
	if err != nil {
		return nil, fmt.Errorf("ChatAI API error: err=%w", err)
	}
	return &ai.ChatAIResponse{Answer: answer.OutputText()}, nil
}
