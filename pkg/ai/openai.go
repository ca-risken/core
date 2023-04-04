package ai

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/model"
	"github.com/coocood/freecache"
	"github.com/sashabaranov/go-openai"
)

const (
	CACHE_SIZE       = 500 * 1024 * 1024 // 500MB
	CACHE_EXPIRE_SEC = 3600
	CACHE_KEY_FORMAT = "OpenAICache/%d/%s"

	LANG_JP               = "ja"
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
	cache        *freecache.Cache
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
		cache:        freecache.NewCache(CACHE_SIZE),
	}
	return &client
}

func (a *AIClient) AskAISummaryFromFinding(ctx context.Context, f *model.Finding, r *model.Recommend, lang string) (string, error) {
	cacheKey := generateCacheKey(fmt.Sprintf(CACHE_KEY_FORMAT, f.FindingID, lang))
	if got, err := a.cache.Get(cacheKey); err == nil {
		a.logger.Infof(ctx, "Cache HIT: finding_id=%d, lang=%s", f.FindingID, lang)
		return string(got), nil
	}
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
	fields := map[string]interface{}{
		"openai_token": resp.Usage.TotalTokens,
	}
	a.logger.WithItemsf(ctx, logging.InfoLevel, fields, "OpenAI usage: %+v", resp.Usage)
	answer := resp.Choices[0].Message.Content
	if err := a.cache.Set(cacheKey, []byte(answer), CACHE_EXPIRE_SEC); err != nil {
		return "", fmt.Errorf("cache set error: err=%w", err)
	}
	return answer, nil
}

func generateFindingDataForAI(f *model.Finding, r *model.Recommend) string {
	text := fmt.Sprintf(FINDING_FORMAT_FOR_AI, f.DataSource, f.Description)
	if r != nil {
		text += fmt.Sprintf(RECOMMEND_FORMAT_FOR_AI, r.Risk, r.Recommendation)
	}
	return text
}

func generateCacheKey(content string) []byte {
	hash := md5.Sum([]byte(content))
	return []byte(hash[:])
}
