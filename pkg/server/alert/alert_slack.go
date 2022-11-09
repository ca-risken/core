package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ca-risken/core/pkg/model"
	projectproto "github.com/ca-risken/core/proto/project"
	"github.com/vikyd/zero"
)

func sendSlackNotification(ctx context.Context, notifyURL, notifySetting string, alert *model.Alert, project *projectproto.Project, rules *[]model.AlertRule) error {
	var setting slackNotifySetting
	if err := json.Unmarshal([]byte(notifySetting), &setting); err != nil {
		return err
	}
	if zero.IsZeroVal(setting.WebhookURL) {
		appLogger.Warn(ctx, "Unset webhook_url")
		return nil
	}
	channel := ""
	if !zero.IsZeroVal(setting.Data["channel"]) {
		channel = setting.Data["channel"]
	}
	message := ""
	if !zero.IsZeroVal(setting.Data["message"]) {
		message = setting.Data["message"]
	}

	slackConfig := &slackWebhookConfig{NotificationAlertURL: notifyURL}
	payload, err := slackConfig.GetPayload(ctx, channel, message, alert, project, rules)
	if err != nil {
		return err
	}
	// TODO http tracing
	resp, err := http.PostForm(setting.WebhookURL, url.Values{"payload": {string(payload)}})
	if err != nil {
		appLogger.Errorf(ctx, "Failed to send slack, resp=%+v, err=%+v", resp, err)
		return err
	}
	defer resp.Body.Close()
	return nil
}

func sendSlackTestNotification(ctx context.Context, notifyURL, notifySetting string) error {
	var setting slackNotifySetting
	if err := json.Unmarshal([]byte(notifySetting), &setting); err != nil {
		return err
	}
	if zero.IsZeroVal(setting.WebhookURL) {
		appLogger.Warn(ctx, "Unset webhook_url")
		return nil
	}
	channel := ""
	if !zero.IsZeroVal(setting.Data["channel"]) {
		channel = setting.Data["channel"]
	}

	slackConfig := slackWebhookConfig{NotificationAlertURL: notifyURL}
	payload, err := slackConfig.GetTestPayload(channel)
	if err != nil {
		return err
	}
	// TODO http tracing
	resp, err := http.PostForm(setting.WebhookURL, url.Values{"payload": {string(payload)}})
	if err != nil {
		appLogger.Errorf(ctx, "Failed to send slack, resp=%+v, err=%+v", resp, err)
		return err
	}
	defer resp.Body.Close()
	return nil
}

type slackNotifySetting struct {
	WebhookURL string            `json:"webhook_url"`
	Data       map[string]string `json:"data"`
}

type slackWebhookConfig struct {
	NotificationAlertURL string
}

func (s *slackWebhookConfig) GetPayload(ctx context.Context, channel, message string, alert *model.Alert, project *projectproto.Project, rules *[]model.AlertRule) (string, error) {
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
					"value": fmt.Sprintf("<%s?project_id=%d&from=slack|詳細はこちらから>", s.NotificationAlertURL, project.ProjectId),
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
	appLogger.Debugf(ctx, "Slack Webhook contents: %s", string(buf))
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
