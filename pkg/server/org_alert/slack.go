package org_alert

import (
	"context"
	"encoding/json"
	"fmt"

	riskenslack "github.com/ca-risken/core/pkg/slack"
	"github.com/slack-go/slack"
)

const (
	LocaleJa                       = "ja"
	LocaleEn                       = "en"
	slackTestNotificationMessageJa = "RISKENからのテスト通知です"
	slackTestNotificationMessageEn = "This is a test notification from RISKEN"
)

func (s *OrgAlertService) sendSlackTestNotification(ctx context.Context, notifySetting string) error {
	var setting riskenslack.NotifySetting
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
