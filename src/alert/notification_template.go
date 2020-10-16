package main

import (
	"encoding/json"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/vikyd/zero"
)

type slackWebhookConfig struct {
	Channel              string
	NotificationAlertUrl string `split_words:"true"`
}

func newslackWebhookConfig(channel string) (*slackWebhookConfig, error) {
	config := &slackWebhookConfig{}
	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}
	config.Channel = channel
	return config, nil
}

func (t *slackWebhookConfig) GetPayload() (string, error) {
	text := fmt.Sprintf(`設定されたAlertに合致する結果を検知しました。
以下のリンクからご確認ください。
%v
`, t.NotificationAlertUrl)
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
