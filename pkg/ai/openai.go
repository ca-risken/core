package ai

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/ai"
	"github.com/ca-risken/core/proto/finding"
	"github.com/coocood/freecache"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const (
	CACHE_SIZE       = 200 * 1024 * 1024 // 200MB
	CACHE_EXPIRE_SEC = 3600
	CACHE_KEY_FORMAT = "OpenAICache/%s"
)

type AIService interface {
	ChatAI(ctx context.Context, req *ai.ChatAIRequest) (*ai.ChatAIResponse, error)
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
	openaiClient := openai.NewClient(
		option.WithAPIKey(token),
	)
	client := AIClient{
		openaiClient: &openaiClient,
		logger:       logger,
		chatGPTModel: model,
		cache:        freecache.NewCache(CACHE_SIZE),
	}
	return &client
}

func generateCacheKey(content string) []byte {
	hash := md5.Sum([]byte(content))
	return []byte(hash[:])
}

func (a *AIClient) getAICache(ctx context.Context, key string) string {
	cacheKey := generateCacheKey(fmt.Sprintf(CACHE_KEY_FORMAT, key))
	if got, err := a.cache.Get(cacheKey); err == nil {
		a.logger.Infof(ctx, "Cache HIT: key=%s", key)
		return string(got)
	}
	return ""
}

func (a *AIClient) setAICache(key string, answer string) error {
	cacheKey := generateCacheKey(fmt.Sprintf(CACHE_KEY_FORMAT, key))
	if err := a.cache.Set(cacheKey, []byte(answer), CACHE_EXPIRE_SEC); err != nil {
		return err
	}
	return nil
}
