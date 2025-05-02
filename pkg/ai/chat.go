package ai

import (
	"context"
	"fmt"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/ai"
	"github.com/sashabaranov/go-openai"
)

func (a *AIClient) ChatAI(ctx context.Context, req *ai.ChatAIRequest) (*ai.ChatAIResponse, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "あなたはAIチャットアシスタントです。ユーザの質問に対して丁寧かつ簡潔に回答してください。",
		},
	}
	for _, h := range req.ChatHistory {
		role := openai.ChatMessageRoleUser
		if h.Role == *ai.ChatRole_CHAT_ROLE_ASSISTANT.Enum() {
			role = openai.ChatMessageRoleAssistant
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: h.Content,
		})
	}
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: req.Question,
	})
	answer, err := a.chatOpenAI(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("ChatAI API error: err=%w", err)
	}
	return &ai.ChatAIResponse{Answer: answer}, nil
}

func (a *AIClient) chatOpenAI(ctx context.Context, messages []openai.ChatCompletionMessage) (string, error) {
	resp, err := a.openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    a.chatGPTModel,
		Messages: messages,
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
