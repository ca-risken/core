package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ca-risken/core/pkg/model"
	projectproto "github.com/ca-risken/core/proto/project"
	"github.com/slack-go/slack"
)

type slackNotifySetting struct {
	WebhookURL string            `json:"webhook_url"`
	Data       slackNotifyOption `json:"data"`
}

type slackNotifyOption struct {
	Channel string `json:"channel,omitempty"`
	Message string `json:"message,omitempty"`
}

func sendSlackNotification(ctx context.Context, notifyURL, notifySetting string, alert *model.Alert, project *projectproto.Project, rules *[]model.AlertRule) error {
	var setting slackNotifySetting
	if err := json.Unmarshal([]byte(notifySetting), &setting); err != nil {
		return err
	}
	if setting.WebhookURL == "" {
		appLogger.Warn(ctx, "Unset webhook_url")
		return nil
	}

	payload := getPayload(ctx, setting.Data.Channel, setting.Data.Message, notifyURL, alert, project, rules)
	// TODO http tracing
	if err := slack.PostWebhook(setting.WebhookURL, payload); err != nil {
		appLogger.Errorf(ctx, "Failed to send slack, err=%+v", err)
		return err
	}
	return nil
}

func sendSlackTestNotification(ctx context.Context, notifyURL, notifySetting string) error {
	var setting slackNotifySetting
	if err := json.Unmarshal([]byte(notifySetting), &setting); err != nil {
		return err
	}
	if setting.WebhookURL == "" {
		appLogger.Warn(ctx, "Unset webhook_url")
		return nil
	}

	payload := getTestPayload(setting.Data.Channel)
	// TODO http tracing
	if err := slack.PostWebhook(setting.WebhookURL, payload); err != nil {
		appLogger.Errorf(ctx, "Failed to send slack, err=%+v", err)
		return err
	}
	return nil
}

func getPayload(
	ctx context.Context,
	channel string,
	message string,
	notifyURL string,
	alert *model.Alert,
	project *projectproto.Project,
	rules *[]model.AlertRule,
) *slack.WebhookMessage {

	// attachments
	attachment := slack.Attachment{
		Color: getColor(alert.Severity),
		Fields: []slack.AttachmentField{
			{
				Title: "Project",
				Value: project.Name,
				Short: true,
			},
			{
				Title: "Severity",
				Value: alert.Severity,
				Short: true,
			},
			{
				Title: "Description",
				Value: alert.Description,
				Short: true,
			},
			{
				Title: "Link",
				Value: fmt.Sprintf("<%s?project_id=%d&from=slack|詳細はこちらから>", notifyURL, project.ProjectId),
				Short: true,
			},
			{
				Title: "Rules",
				Value: generateRuleList(rules),
			},
		},
		Ts: json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := slack.WebhookMessage{
		Text:        fmt.Sprintf("%vアラートを検知しました。", getMention(alert.Severity)),
		Attachments: []slack.Attachment{attachment},
	}

	// override message
	if message != "" {
		msg.Text = message // update text
	}
	if channel != "" {
		msg.Channel = channel // add channel
	}

	appLogger.Debugf(ctx, "Slack Webhook contents: %+v", msg)
	return &msg
}

func getTestPayload(channel string) *slack.WebhookMessage {
	msg := slack.WebhookMessage{
		Text: "RISKENからのテスト通知です",
	}
	// override message
	if channel != "" {
		msg.Channel = channel
	}
	return &msg
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
