package organization_alert

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/slack-go/slack"
)

type slackNotifySetting struct {
	WebhookURL string            `json:"webhook_url"`
	ChannelID  string            `json:"channel_id"`
	Data       slackNotifyOption `json:"data"`
	Locale     string            `json:"locale"`
}

type slackNotifyOption struct {
	Channel string `json:"channel,omitempty"`
	Message string `json:"message,omitempty"`
}

func (s *OrgAlertService) replaceSlackNotifySetting(ctx context.Context, jsonNotifySettingExist, jsonNotifySettingUpdate string) (slackNotifySetting, error) {
	var notifySettingUpdate slackNotifySetting
	if err := json.Unmarshal([]byte(jsonNotifySettingUpdate), &notifySettingUpdate); err != nil {
		s.logger.Errorf(ctx, "Error occured when unmarshal update.NotifySetting. err: %v", err)
		return slackNotifySetting{}, err
	}
	var notifySettingExist slackNotifySetting
	if err := json.Unmarshal([]byte(jsonNotifySettingExist), &notifySettingExist); err != nil {
		s.logger.Errorf(ctx, "Error occured when unmarshal exist.NotifySetting. err: %v", err)
		return slackNotifySetting{}, err
	}

	// webhookURL and ChannelID are mutually exclusive
	if notifySettingUpdate.WebhookURL != "" {
		notifySettingUpdate.ChannelID = ""
		return notifySettingUpdate, nil
	}
	if notifySettingUpdate.ChannelID != "" {
		notifySettingUpdate.WebhookURL = ""
		return notifySettingUpdate, nil
	}

	// No update options
	notifySettingUpdate = notifySettingExist
	return notifySettingUpdate, nil
}

func maskingNotifySetting(notificationType, notifySetting string) (string, error) {
	switch notificationType {
	case "slack":
		var setting slackNotifySetting
		if err := json.Unmarshal([]byte(notifySetting), &setting); err != nil {
			return "", err
		}
		if setting.WebhookURL == "" {
			return notifySetting, nil
		}
		setting.WebhookURL = maskRight(setting.WebhookURL, len(setting.WebhookURL)/2)
		ret, err := json.Marshal(setting)
		if err != nil {
			return "", err
		}
		return string(ret), nil
	default:
		return notifySetting, nil
	}
}

func maskRight(s string, num int) string {
	rs := []rune(s)
	for i := num; i < len(rs); i++ {
		rs[i] = '*'
	}
	return string(rs)
}

const (
	LocaleJa                         = "ja"
	LocaleEn                         = "en"
	slackTestNotificationMessageJa   = "RISKENからのテスト通知です"
	slackTestNotificationMessageEn   = "This is a test notification from RISKEN"
)

func (s *OrgAlertService) sendSlackTestNotification(ctx context.Context, notifySetting string) error {
	var setting slackNotifySetting
	if err := json.Unmarshal([]byte(notifySetting), &setting); err != nil {
		return err
	}

	locale := s.getLocale(setting.Locale)
	msg := getTestSlackMessageText(locale)

	if setting.WebhookURL != "" {
		webhookMsg := &slack.WebhookMessage{Text: msg}
		if setting.Data.Channel != "" {
			webhookMsg.Channel = setting.Data.Channel
		}
		if err := slack.PostWebhook(setting.WebhookURL, webhookMsg); err != nil {
			return fmt.Errorf("failed to send slack(webhookurl): %w", err)
		}
	} else if setting.ChannelID != "" {
		if _, _, err := s.slackClient.PostMessage(setting.ChannelID, slack.MsgOptionText(msg, false)); err != nil {
			return fmt.Errorf("failed to send slack(postmessage): %w", err)
		}
	}
	return nil
}

func (s *OrgAlertService) getLocale(settingLocale string) string {
	switch settingLocale {
	case LocaleJa:
		return LocaleJa
	case LocaleEn:
		return LocaleEn
	default:
		return s.defaultLocale
	}
}

func getTestSlackMessageText(locale string) string {
	switch locale {
	case LocaleJa:
		return slackTestNotificationMessageJa
	default:
		return slackTestNotificationMessageEn
	}
}
