package main

import (
	"encoding/json"
	"fmt"

	"github.com/vikyd/zero"
)

type slackWebhookSetting struct {
	Channel string
	AlertID uint32
}

func (t *slackWebhookSetting) GetPayload() (string, error) {
	text := fmt.Sprintf(`設定されたAlertに合致する結果を検知しました。
AlertID: %v
実際にはAlertに飛べるリンクなどつけるのが良さそう(画面できてから設定かな)`, t.AlertID)
	payload := map[string]string{}
	payload["text"] = text
	if !zero.IsZeroVal(t.Channel) {
		payload["channel"] = t.Channel
	}
	p, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(p), err
}
