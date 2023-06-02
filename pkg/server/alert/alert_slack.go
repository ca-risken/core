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
	Locale     string            `json:"locale"`
}

type slackNotifyOption struct {
	Channel string `json:"channel,omitempty"`
	Message string `json:"message,omitempty"`
}

const (
	LocaleJa                       = "ja"
	LocaleEn                       = "en"
	slackNotificationMessageJa     = "%vアラートを検知しました。"
	slackNotificationMessageEn     = "%vDetected alerts."
	slackNotificationAttachmentJa  = "その他、%d件すべてのFindingは <%s/#/alert/alert?project_id=%d&from=slack|アラート画面> からご確認ください。"
	slackNotificationAttachmentEn  = "Please check all %d Findings from <%s/#/alert/alert?project_id=%d&from=slack|Alert screen>."
	slackNotificationTestMessageJa = "RISKENからのテスト通知です"
	slackNotificationTestMessageEn = "This is a test notification from RISKEN"
)

func sendSlackNotification(
	ctx context.Context, url, notifySetting string,
	alert *model.Alert,
	project *projectproto.Project,
	rules *[]model.AlertRule,
	findings *findingDetail,
	defaultLocale string,
) error {
	var setting slackNotifySetting
	if err := json.Unmarshal([]byte(notifySetting), &setting); err != nil {
		return err
	}
	if setting.WebhookURL == "" {
		return nil
	}
	var locale string
	switch setting.Locale {
	case LocaleJa:
		locale = LocaleJa
	case LocaleEn:
		locale = LocaleEn
	default:
		locale = defaultLocale
	}

	payload := getPayload(ctx, setting.Data.Channel, setting.Data.Message, url, alert, project, rules, findings, locale)
	// TODO http tracing
	if err := slack.PostWebhook(setting.WebhookURL, payload); err != nil {
		return fmt.Errorf("failed to send slack: %w", err)
	}
	return nil
}

func sendSlackTestNotification(ctx context.Context, url, notifySetting string) error {
	var setting slackNotifySetting
	if err := json.Unmarshal([]byte(notifySetting), &setting); err != nil {
		return err
	}
	if setting.WebhookURL == "" {
		return nil
	}

	payload := getTestPayload(setting.Data.Channel)
	// TODO http tracing
	if err := slack.PostWebhook(setting.WebhookURL, payload); err != nil {
		return fmt.Errorf("failed to send slack: %w", err)
	}
	return nil
}

func getPayload(
	ctx context.Context,
	channel string,
	message string,
	url string,
	alert *model.Alert,
	project *projectproto.Project,
	rules *[]model.AlertRule,
	findings *findingDetail,
	locale string,
) *slack.WebhookMessage {

	// attachments
	attachment := slack.Attachment{
		Color: getColor(alert.Severity),
		Fields: []slack.AttachmentField{
			{
				Value: fmt.Sprintf("<%s/#/alert/alert?project_id=%d&from=slack|%s>", url, project.ProjectId, alert.Description),
			},
			{
				Title: "Rules",
				Value: generateRuleList(rules),
			},
			{
				Title: "Project",
				Value: project.Name,
				Short: true,
			},
			{
				Title: "Findings",
				Value: fmt.Sprint(findings.FindingCount),
				Short: true,
			},
		},
	}

	var msgText string
	switch locale {
	case LocaleJa:
		msgText = slackNotificationMessageJa
	default:
		msgText = slackNotificationMessageEn
	}
	msg := slack.WebhookMessage{
		Text:        fmt.Sprintf(msgText, getMention(alert.Severity)),
		Attachments: []slack.Attachment{attachment},
	}
	msg.Attachments = append(msg.Attachments, *getFindingAttachment(url, project.ProjectId, findings, locale)...)

	// override message
	if message != "" {
		msg.Text = message // update text
	}
	if channel != "" {
		msg.Channel = channel // add channel
	}
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

const (
	MAX_NOTIFY_RULE_NUM = 3
)

func generateRuleList(rules *[]model.AlertRule) string {
	if rules == nil {
		return ""
	}
	list := ""
	for idx, rule := range *rules {
		if idx >= MAX_NOTIFY_RULE_NUM {
			list = fmt.Sprintf("%s\n- %s", list, "...")
			return list
		}
		if idx != 0 {
			list += "\n"
		}
		list += fmt.Sprintf("- %s", rule.Name)
	}
	return list
}

func getFindingAttachment(url string, projectID uint32, findings *findingDetail, locale string) *[]slack.Attachment {
	attachments := []slack.Attachment{}
	for _, f := range findings.Exampls {
		a := slack.Attachment{
			Color: getColorByScore(f.Score),
			Fields: []slack.AttachmentField{
				{
					Value: fmt.Sprintf("<%s/#/finding/finding?project_id=%d&finding_id=%d&from_score=0&status=1&from=slack|%s>", url, projectID, f.FindingID, f.Description),
				},
				{
					Title: "DataSource",
					Value: f.DataSource,
					Short: true,
				},
				{
					Title: "ResourceName",
					Value: f.ResourceName,
					Short: true,
				},
				{
					Title: "Tags",
					Value: generateTagContentByFinding(f.Tags),
				},
			},
		}
		attachments = append(attachments, a)
	}
	if findings.FindingCount > len(findings.Exampls) {
		var attachmentText string
		switch locale {
		case LocaleJa:
			attachmentText = slackNotificationAttachmentJa
		default:
			attachmentText = slackNotificationAttachmentEn
		}
		attachments = append(attachments, slack.Attachment{
			Color: "grey",
			Fields: []slack.AttachmentField{
				{
					Value: fmt.Sprintf(attachmentText, findings.FindingCount, url, projectID),
				},
			},
			Ts: json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
		})
	}
	return &attachments
}

func getColorByScore(score float32) string {
	switch {
	case score >= 0.8:
		return "danger"
	case score >= 0.6:
		return "warning"
	default:
		return "good"
	}
}

const (
	MAX_NOTIFY_FINDING_TAG_NUM = 15
)

func generateTagContentByFinding(tags []string) string {
	content := ""
	for idx, t := range tags {
		if content != "" {
			content += " "
		}
		content += fmt.Sprintf("`%s`", t)
		if idx+1 >= MAX_NOTIFY_FINDING_TAG_NUM {
			content += " ..."
			return content
		}
	}
	return content
}
