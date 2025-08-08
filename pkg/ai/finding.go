package ai

import (
	"context"
	"fmt"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
)

const (
	LANG_JP              = "ja"
	PROMPT_SYSTEM_MSG_EN = "You are a helpful security advisor. Please explain this in a way that a non-security expert can understand."
	PROMPT_SYSTEM_MSG_JP = "あなたは役に立つセキュリティアドバイザーです。セキュリティの専門家ではない人にも理解できるように説明をお願いします。"
	PROMPT_SUMMARY_EN    = `I have detected the following security issue in my cloud environment. Please summarize the contents.
Also, please include any ways to address the issue.

The definition of the score is as follows. If the score is low, the problem may be resolved.
<score>
- 0.0 ~ 0.3 [Low]: Already low risk. If you are curious, check it.
- 0.4 ~ 0.7 [Midium]: There is a risk of being connected. It's not necessary to do it immediately, but check it when you have time.
- 0.8 ~ 1.0 [High]: It's dangerous. Check it immediately if necessary.
</score>

Use the following markdown format for your response.
<format>
## Summary

## Detection content
- aaa
- bbb

## How to fix
- aaa
- bbb
</format>
`
	PROMPT_SUMMARY_JP = `クラウド環境で以下のセキュリティの問題を検知しました。日本語で内容を要約してください。
また、問題の対処方法もあれば含めてください。

スコアの定義は以下を参考にしてください。スコアが低い場合には問題が解決されてる可能性があります。
<score>
- 0.0 ~ 0.3 [低]: 無視しても大丈夫です。既に十分リスクは低い状態ですが気になる場合は確認してください。
- 0.4 ~ 0.7 [中]: リスクにつながる可能性があります。すぐにではないですが時間がある時に確認してください。
- 0.8 ~ 1.0 [高]: 危険という判定です。すぐに確認し必要があればインシデント発生前に対処すべきです。
</score>

回答は以下のMarkdownフォーマットでお願いします。(URLリンクは前後に半角スペースを入れてください)
<format>
## 要約

## 検出内容
・aaa
・bbb

## 対処方法
・aaa
・bbb
</format>
`
	FINDING_FORMAT_FOR_AI = `The RISKEN tool detected the following issue related to cloud security.
Score: 
%.1f

Type: 
%s

Description: 
%s

ScanResult(json):
%s
`
	RECOMMEND_FORMAT_FOR_AI = `
Detail: %s

Recommendation: %s
`
)

func (a *AIClient) AskAISummaryFromFinding(ctx context.Context, f *model.Finding, r *model.Recommend, lang string) (string, error) {
	if summaryCache := a.getAICache(ctx, generateCacheKeyForFinding(f.FindingID, lang)); summaryCache != "" {
		a.logger.Infof(ctx, "Cache HIT: finding_id=%d, lang=%s", f.FindingID, lang)
		return summaryCache, nil
	}

	instruction, inputs := generateAskAISummaryInputs(f, r, lang)
	answer, err := a.responsesAPI(ctx, a.chatGPTModel, instruction, inputs, DefaultTools, "")
	if err != nil {
		return "", fmt.Errorf("openai API error: finding_id=%d, err=%w", f.FindingID, err)
	}
	if err := a.setAICache(generateCacheKeyForFinding(f.FindingID, lang), answer.OutputText()); err != nil {
		return "", fmt.Errorf("cache set error: err=%w", err)
	}
	return answer.OutputText(), nil
}

func (a *AIClient) AskAISummaryStreamFromFinding(
	ctx context.Context, f *model.Finding, r *model.Recommend, lang string, stream finding.FindingService_GetAISummaryStreamServer,
) error {
	if summaryCache := a.getAICache(ctx, generateCacheKeyForFinding(f.FindingID, lang)); summaryCache != "" {
		a.logger.Infof(ctx, "Cache HIT: finding_id=%d, lang=%s", f.FindingID, lang)
		if sendErr := stream.Send(&finding.GetAISummaryResponse{Answer: summaryCache}); sendErr != nil {
			return sendErr
		}
		return nil
	}

	instruction, inputs := generateAskAISummaryInputs(f, r, lang)
	answer, err := responsesStreamingAPI(
		ctx,
		a.openaiClient,
		a.chatGPTModel,
		instruction,
		inputs,
		DefaultTools,
		stream,
		func(s string) *finding.GetAISummaryResponse { return &finding.GetAISummaryResponse{Answer: s} },
		a.logger,
	)
	if err != nil {
		return fmt.Errorf("openai API error: finding_id=%d, err=%w", f.FindingID, err)
	}

	if err := a.setAICache(generateCacheKeyForFinding(f.FindingID, lang), answer); err != nil {
		return fmt.Errorf("cache set error: err=%w", err)
	}
	return nil
}

func generateFindingDataForAI(f *model.Finding, r *model.Recommend) string {
	text := fmt.Sprintf(FINDING_FORMAT_FOR_AI, f.Score, f.DataSource, f.Description, f.Data)
	if r != nil {
		text += fmt.Sprintf(RECOMMEND_FORMAT_FOR_AI, r.Risk, r.Recommendation)
	}
	return text
}

func generateCacheKeyForFinding(findingID uint64, lang string) string {
	return fmt.Sprintf("%d/%s", findingID, lang)
}

func getAskAISummaryPrompt(lang string) (promptSystem, promptSummary string) {
	promptSystem = PROMPT_SYSTEM_MSG_EN
	promptSummary = PROMPT_SUMMARY_EN
	if lang == LANG_JP {
		promptSystem = PROMPT_SYSTEM_MSG_JP
		promptSummary = PROMPT_SUMMARY_JP
	}
	return
}

func generateAskAISummaryInputs(f *model.Finding, r *model.Recommend, lang string) (string, responses.ResponseNewParamsInputUnion) {
	promptSystem, promptSummary := getAskAISummaryPrompt(lang)
	inputParam := responses.ResponseInputParam{
		responses.ResponseInputItemUnionParam{
			OfMessage: &responses.EasyInputMessageParam{
				Role: responses.EasyInputMessageRoleAssistant,
				Content: responses.EasyInputMessageContentUnionParam{
					OfString: openai.String(promptSummary),
				},
			},
		},
		responses.ResponseInputItemUnionParam{
			OfMessage: &responses.EasyInputMessageParam{
				Role: responses.EasyInputMessageRoleUser,
				Content: responses.EasyInputMessageContentUnionParam{
					OfString: openai.String(generateFindingDataForAI(f, r)),
				},
			},
		},
	}
	return promptSystem, responses.ResponseNewParamsInputUnion{OfInputItemList: inputParam}
}
