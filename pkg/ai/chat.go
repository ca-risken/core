package ai

import (
	"context"
	"fmt"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/ai"
	"github.com/openai/openai-go"
)

func (a *AIClient) ChatAI(ctx context.Context, req *ai.ChatAIRequest) (*ai.ChatAIResponse, error) {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(
			"You are an AI chat assistant. Please respond to user questions politely and concisely." +
				"Please generate appropriate responses according to the user's language."),
	}
	for _, h := range req.ChatHistory {
		if h.Role == *ai.ChatRole_CHAT_ROLE_ASSISTANT.Enum() {
			messages = append(messages, openai.AssistantMessage(h.Content))
		} else {
			messages = append(messages, openai.UserMessage(h.Content))
		}
	}
	messages = append(messages, openai.UserMessage(req.Question))
	answer, err := a.chatOpenAI(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("ChatAI API error: err=%w", err)
	}
	return &ai.ChatAIResponse{Answer: answer}, nil
}

func (a *AIClient) chatOpenAI(ctx context.Context, messages []openai.ChatCompletionMessageParamUnion) (string, error) {
	resp, err := a.openaiClient.Chat.Completions.New(ctx,
		openai.ChatCompletionNewParams{
			Model:    a.chatGPTModel,
			Messages: messages,
			Seed:     openai.Int(1),
		})
	if err != nil {
		return "", fmt.Errorf("CreateChatCompletion API error: err=%w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("CreateChatCompletion API: no response")
	}
	fields := map[string]any{
		"openai_token": resp.Usage.TotalTokens,
	}
	a.logger.WithItemsf(ctx, logging.InfoLevel, fields, "OpenAI usage: %+v", resp.Usage)
	answer := resp.Choices[0].Message.Content
	return answer, nil
}
