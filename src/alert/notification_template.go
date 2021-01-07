package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/kelseyhightower/envconfig"
	"github.com/vikyd/zero"
)

type slackWebhookConfig struct {
	Channel              string
	NotificationAlertURL string `split_words:"true"`
}

func newslackWebhookConfig(channel string) (*slackWebhookConfig, error) {
	config := &slackWebhookConfig{}
	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}
	config.Channel = channel
	return config, nil
}

func (t *slackWebhookConfig) GetPayload(alert *model.Alert, projectName string) (string, error) {
	now := time.Now().Unix()
	text := fmt.Sprintf("%vアラートを検知しました。", getMention(alert.Severity))
	attachments := []interface{}{
		map[string]interface{}{
			"color": getColor(alert.Severity),
			"fields": []interface{}{
				map[string]string{
					"title": "Project",
					"value": projectName,
					"short": "true",
				},
				map[string]string{
					"title": "Severity",
					"value": alert.Severity,
					"short": "true",
				},
				map[string]string{
					"title": "Link",
					"value": fmt.Sprintf("<%s|詳細はこちらから>", t.NotificationAlertURL),
					"short": "true",
				},
				map[string]string{
					"title": "Description",
					"value": alert.Description,
					"short": "true",
				},
			},
			"footer": "Send from RISKEN",
			"ts":     now,
		},
	}
	payload := map[string]interface{}{}
	payload["text"] = text

	payload["attachments"] = attachments
	if !zero.IsZeroVal(t.Channel) {
		payload["channel"] = t.Channel
	}
	p, err := json.Marshal(payload)
	appLogger.Infof("json: %v", string(p))
	if err != nil {
		return "", err
	}
	return string(p), err
}

func getColor(severity string) string {
	switch severity {
	case "high":
		return "danger"
	case "medium":
		return "warning"
	case "low":
		return "good"
	default:
		return "good"
	}
}

func getMention(severity string) string {
	switch severity {
	case "high":
		return "<!channel> "
	case "medium":
		return "<!here> "
	case "low":
		return ""
	default:
		return ""
	}
}
