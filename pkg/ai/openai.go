package ai

import (
	"context"
	"fmt"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/model"
	"github.com/sashabaranov/go-openai"
)

const (
	LANG_JP               = "jp"
	PROMPT_SYSTEM_MSG_EN  = "You are a helpful security assistant."
	PROMPT_SYSTEM_MSG_JP  = "あなたは役に立つセキュリティアシスタントです。"
	PROMPT_SUMMARY_EN     = "I have detected the following security issue in my cloud environment. Please summarize in 500 characters or less. Also, please include any ways to address the issue. By the way, I am not a security expert."
	PROMPT_SUMMARY_JP     = "私のクラウド環境で以下のセキュリティの問題を検知しました。500字以内の日本語で内容を要約してください。また、問題の対処方法もあれば含めてください。ちなみに私はセキュリティの専門家ではありません。"
	FINDING_FORMAT_FOR_AI = `The RISKEN tool detected the following issue related to cloud security.
Type: %s
Description: %s
`
	RECOMMEND_FORMAT_FOR_AI = `Detail: %s
Recommendation: %s
`
)

type AIService interface {
	AskAISummaryFromFinding(ctx context.Context, f *model.Finding, r *model.Recommend, lang string) (string, error)
}

type AIClient struct {
	openaiClient *openai.Client
	logger       logging.Logger
}

var _ AIService = (*AIClient)(nil)

func NewAIClient(token string, logger logging.Logger) AIService {
	if token == "" {
		return nil
	}
	client := AIClient{
		openaiClient: openai.NewClient(token),
		logger:       logger,
	}
	return &client
}

func (a *AIClient) AskAISummaryFromFinding(ctx context.Context, f *model.Finding, r *model.Recommend, lang string) (string, error) {
	promptSystem := PROMPT_SYSTEM_MSG_EN
	promptSummary := PROMPT_SUMMARY_EN
	if lang == LANG_JP {
		promptSystem = PROMPT_SYSTEM_MSG_JP
		promptSummary = PROMPT_SUMMARY_JP
	}
	resp, err := a.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: promptSystem,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: promptSummary,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: generateFindingDataForAI(f, r),
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("openai API error: err=%w", err)
	}
	a.logger.Infof(ctx, "OpenAI response: usage=%+v, resp=%+v", resp.Usage, resp)
	return resp.Choices[0].Message.Content, nil
}

func generateFindingDataForAI(f *model.Finding, r *model.Recommend) string {
	text := fmt.Sprintf(FINDING_FORMAT_FOR_AI, f.DataSource, f.Description)
	if r != nil {
		text += fmt.Sprintf(RECOMMEND_FORMAT_FOR_AI, r.Risk, r.Recommendation)
	}
	return text
}
