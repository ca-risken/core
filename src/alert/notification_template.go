package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/core/src/alert/model"
	"github.com/vikyd/zero"
)

type slackWebhookConfig struct {
	NotificationAlertURL string `split_words:"true" default:"http://localhost"`
}

func (s *slackWebhookConfig) GetPayload(channel, message string, alert *model.Alert, project *project.Project, rules *[]model.AlertRule) (string, error) {
	payload := map[string]interface{}{}
	// text
	text := fmt.Sprintf("%vアラートを検知しました。", getMention(alert.Severity))
	if !zero.IsZeroVal(message) {
		text = message // update message
	}
	payload["text"] = text

	// channel
	if !zero.IsZeroVal(channel) {
		payload["channel"] = channel
	}

	// attachments
	now := time.Now().Unix()
	attachments := []interface{}{
		map[string]interface{}{
			"color": getColor(alert.Severity),
			"fields": []interface{}{
				map[string]string{
					"title": "Project",
					"value": project.Name,
					"short": "true",
				},
				map[string]string{
					"title": "Severity",
					"value": alert.Severity,
					"short": "true",
				},
				map[string]string{
					"title": "Description",
					"value": alert.Description,
					"short": "true",
				},
				map[string]string{
					"title": "Link",
					"value": fmt.Sprintf("<%s?project_id=%d|詳細はこちらから>", s.NotificationAlertURL, project.ProjectId),
					"short": "true",
				},
				map[string]string{
					"title": "Rules",
					"value": generateRuleList(rules),
				},
			},
			"footer": "Send from RISKEN",
			"ts":     now,
		},
	}
	payload["attachments"] = attachments
	buf, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	appLogger.Debugf("Slack Webhook contents: %s", string(buf))
	return string(buf), err
}

func (s *slackWebhookConfig) GetTestPayload(channel string) (string, error) {
	payload := map[string]interface{}{}
	// text
	text := "RISKENからのテスト通知です"
	payload["text"] = text

	// channel
	if !zero.IsZeroVal(channel) {
		payload["channel"] = channel
	}

	buf, err := json.Marshal(payload)
	return string(buf), err
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

func generateRuleList(rules *[]model.AlertRule) string {
	if rules == nil {
		return ""
	}
	list := ""
	for idx, rule := range *rules {
		if idx == 0 {
			list = fmt.Sprintf("- %s", rule.Name)
			continue
		}
		list = fmt.Sprintf("%s\n- %s", list, rule.Name)
		if idx >= 4 {
			list = fmt.Sprintf("%s\n- %s", list, "...")
			break
		}
	}
	return list
}
