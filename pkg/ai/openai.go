package ai

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	"github.com/coocood/freecache"
	"github.com/sashabaranov/go-openai"
)

const (
	CACHE_SIZE       = 200 * 1024 * 1024 // 200MB
	CACHE_EXPIRE_SEC = 3600
	CACHE_KEY_FORMAT = "OpenAICache/%d/%s"

	LANG_JP              = "ja"
	PROMPT_SYSTEM_MSG_EN = "You are a helpful security advisor. Please explain this in a way that a non-security expert can understand."
	PROMPT_SYSTEM_MSG_JP = "あなたは役に立つセキュリティアドバイザーです。セキュリティの専門家ではない人にも理解できるように説明をお願いします。"
	PROMPT_SUMMARY_EN    = `I have detected the following security issue in my cloud environment. Please summarize the contents.
Also, please include any ways to address the issue.
Use the following format for your response.
<Summary>

<Detection content>
- aaa
- bbb

<How to fix>
- aaa
- bbb
`
	PROMPT_SUMMARY_JP = `クラウド環境で以下のセキュリティの問題を検知しました。日本語で内容を要約してください。
また、問題の対処方法もあれば含めてください。
回答は以下のフォーマットでお願いします。
＜要約＞

＜検出内容＞
・aaa
・bbb

＜対処方法＞
・aaa
・bbb
`
	FINDING_FORMAT_FOR_AI = `The RISKEN tool detected the following issue related to cloud security.
Type: 
%s

Description: 
%s
`
	RECOMMEND_FORMAT_FOR_AI = `Detail: %s
Recommendation: %s
`
)

type AIService interface {
	AskAISummaryFromFinding(ctx context.Context, f *model.Finding, r *model.Recommend, lang string) (string, error)
	AskAISummaryStreamFromFinding(
		ctx context.Context,
		f *model.Finding,
		r *model.Recommend,
		lang string,
		stream finding.FindingService_GetAISummaryStreamServer,
	) error
}

type AIClient struct {
	openaiClient *openai.Client
	cache        *freecache.Cache
	chatGPTModel string // https://platform.openai.com/docs/models/overview
	logger       logging.Logger
}

var _ AIService = (*AIClient)(nil)

func NewAIClient(token, model string, logger logging.Logger) AIService {
	if token == "" {
		return nil
	}
	if model == "" {
		return nil
	}
	client := AIClient{
		openaiClient: openai.NewClient(token),
		logger:       logger,
		chatGPTModel: model,
		cache:        freecache.NewCache(CACHE_SIZE),
	}
	return &client
}

func (a *AIClient) AskAISummaryFromFinding(ctx context.Context, f *model.Finding, r *model.Recommend, lang string) (string, error) {
	if summaryCache := a.getAISummaryCache(ctx, f.FindingID, lang); summaryCache != "" {
		a.logger.Infof(ctx, "Cache HIT: finding_id=%d, lang=%s", f.FindingID, lang)
		return summaryCache, nil
	}

	promptSystem, promptSummary := getPrompt(lang)
	resp, err := a.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: a.chatGPTModel,
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
		return "", fmt.Errorf("openai API error: finding_id=%d, err=%w", f.FindingID, err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai API no response: finding_id=%d", f.FindingID)
	}
	fields := map[string]interface{}{
		"openai_token": resp.Usage.TotalTokens,
	}
	a.logger.WithItemsf(ctx, logging.InfoLevel, fields, "OpenAI usage: %+v", resp.Usage)
	answer := resp.Choices[0].Message.Content
	if err := a.setAISummaryCache(ctx, f.FindingID, lang, answer); err != nil {
		return "", fmt.Errorf("cache set error: err=%w", err)
	}
	return answer, nil
}

func (a *AIClient) AskAISummaryStreamFromFinding(
	ctx context.Context, f *model.Finding, r *model.Recommend, lang string, stream finding.FindingService_GetAISummaryStreamServer,
) error {
	if summaryCache := a.getAISummaryCache(ctx, f.FindingID, lang); summaryCache != "" {
		a.logger.Infof(ctx, "Cache HIT: finding_id=%d, lang=%s", f.FindingID, lang)
		if sendErr := stream.Send(&finding.GetAISummaryResponse{Answer: summaryCache}); sendErr != nil {
			return sendErr
		}
		return nil
	}

	promptSystem, promptSummary := getPrompt(lang)
	streamResp, err := a.openaiClient.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model: a.chatGPTModel,
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
			Stream: true,
		},
	)
	if err != nil {
		return fmt.Errorf("openai API error: finding_id=%d, err=%w", f.FindingID, err)
	}
	defer streamResp.Close()

	var usageTokens int
	var answer string
	for {
		resp, err := streamResp.Recv()
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("openai API streaming error: finding_id=%d, err=%w", f.FindingID, err)
		}

		if resp.Choices != nil || len(resp.Choices) > 0 {
			if sendErr := stream.Send(&finding.GetAISummaryResponse{Answer: resp.Choices[0].Delta.Content}); sendErr != nil {
				return sendErr
			}
			usageTokens += resp.Usage.TotalTokens
			answer += resp.Choices[0].Delta.Content
		}
		if errors.Is(err, io.EOF) {
			break
		}
	}

	fields := map[string]interface{}{
		"openai_token": usageTokens,
	}
	a.logger.WithItemsf(ctx, logging.InfoLevel, fields, "OpenAI usage tokens: %d", usageTokens)
	if err := a.setAISummaryCache(ctx, f.FindingID, lang, answer); err != nil {
		return fmt.Errorf("cache set error: err=%w", err)
	}
	return nil
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

func (a *AIClient) getAISummaryCache(ctx context.Context, findingID uint64, lang string) string {
	cacheKey := generateCacheKey(fmt.Sprintf(CACHE_KEY_FORMAT, findingID, lang))
	if got, err := a.cache.Get(cacheKey); err == nil {
		a.logger.Infof(ctx, "Cache HIT: finding_id=%d, lang=%s", findingID, lang)
		return string(got)
	}
	return ""
}

func (a *AIClient) setAISummaryCache(ctx context.Context, findingID uint64, lang, answer string) error {
	cacheKey := generateCacheKey(fmt.Sprintf(CACHE_KEY_FORMAT, findingID, lang))
	if err := a.cache.Set(cacheKey, []byte(answer), CACHE_EXPIRE_SEC); err != nil {
		return err
	}
	return nil
}

func getPrompt(lang string) (promptSystem, promptSummary string) {
	promptSystem = PROMPT_SYSTEM_MSG_EN
	promptSummary = PROMPT_SUMMARY_EN
	if lang == LANG_JP {
		promptSystem = PROMPT_SYSTEM_MSG_JP
		promptSummary = PROMPT_SUMMARY_JP
	}
	return
}
