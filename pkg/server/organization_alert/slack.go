package organization_alert

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ca-risken/common/pkg/logging"
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

const (
	localeJa                         = "ja"
	localeEn                         = "en"
	slackNotificationTestMessageJa   = "RISKENからのテスト通知です"
	slackNotificationTestMessageEn   = "This is a test notification from RISKEN"
)

func (s *OrganizationAlertService) sendSlackTestNotification(ctx context.Context, notifySetting, defaultLocale string) error {
	var setting slackNotifySetting
	if err := json.Unmarshal([]byte(notifySetting), &setting); err != nil {
		return err
	}

	locale := resolveLocale(setting.Locale, defaultLocale)

	if setting.WebhookURL != "" {
		webhookMsg := getTestWebhookMessage(setting.Data.Channel, locale)
		if err := slack.PostWebhook(setting.WebhookURL, webhookMsg); err != nil {
			return fmt.Errorf("failed to send slack(webhookurl): %w", err)
		}
	} else if setting.ChannelID != "" {
		if err := s.postMessageSlackWithRetry(ctx,
			setting.ChannelID, slack.MsgOptionText(getTestSlackMessageText(locale), false)); err != nil {
			return fmt.Errorf("failed to send slack(postmessage): %w", err)
		}
	}
	return nil
}

func (s *OrganizationAlertService) postMessageSlack(channelID string, msg ...slack.MsgOption) error {
	if _, _, err := s.slackClient.PostMessage(channelID, msg...); err != nil {
		var rateLimitError *slack.RateLimitedError
		if errors.As(err, &rateLimitError) {
			time.Sleep(rateLimitError.RetryAfter)
		}
		return err
	}
	return nil
}

func (s *OrganizationAlertService) postMessageSlackWithRetry(ctx context.Context, channelID string, msg ...slack.MsgOption) error {
	const maxRetries = 3
	var err error
	for i := 0; i < maxRetries; i++ {
		err = s.postMessageSlack(channelID, msg...)
		if err == nil {
			return nil
		}
		s.logger.Warnf(ctx, "[RetryLogger] postMessageSlack error: attempt=%d, err=%+v", i+1, err)
	}
	return err
}

func replaceSlackNotifySetting(ctx context.Context, logger logging.Logger, jsonNotifySettingExist, jsonNotifySettingUpdate string) (slackNotifySetting, error) {
	var notifySettingUpdate slackNotifySetting
	if err := json.Unmarshal([]byte(jsonNotifySettingUpdate), &notifySettingUpdate); err != nil {
		logger.Errorf(ctx, "Error occured when unmarshal update.NotifySetting. err: %v", err)
		return slackNotifySetting{}, err
	}
	var notifySettingExist slackNotifySetting
	if err := json.Unmarshal([]byte(jsonNotifySettingExist), &notifySettingExist); err != nil {
		logger.Errorf(ctx, "Error occured when unmarshal exist.NotifySetting. err: %v", err)
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

func resolveLocale(settingLocale, defaultLocale string) string {
	switch settingLocale {
	case localeJa:
		return localeJa
	case localeEn:
		return localeEn
	default:
		return defaultLocale
	}
}

func getTestWebhookMessage(channel, locale string) *slack.WebhookMessage {
	msg := slack.WebhookMessage{
		Text: getTestSlackMessageText(locale),
	}
	if channel != "" {
		msg.Channel = channel
	}
	return &msg
}

func getTestSlackMessageText(locale string) string {
	switch locale {
	case localeJa:
		return slackNotificationTestMessageJa
	default:
		return slackNotificationTestMessageEn
	}
}
