package org_alert

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	riskenslack "github.com/ca-risken/core/pkg/slack"
	"github.com/cenkalti/backoff/v4"
	goslack "github.com/slack-go/slack"
)

func (s *OrgAlertService) postMessageSlack(channelID string, msg ...goslack.MsgOption) error {
	if _, _, err := s.slackClient.PostMessage(channelID, msg...); err != nil {
		var rateLimitError *goslack.RateLimitedError
		if errors.As(err, &rateLimitError) {
			time.Sleep(rateLimitError.RetryAfter)
		}
		return err
	}
	return nil
}

func (s *OrgAlertService) postMessageSlackWithRetry(ctx context.Context, channelID string, msg ...goslack.MsgOption) error {
	operation := func() error {
		return s.postMessageSlack(channelID, msg...)
	}
	return backoff.RetryNotify(operation, s.retryer, s.newRetryLogger(ctx, "postMessageSlack"))
}

func (s *OrgAlertService) sendTestOrgNotification(ctx context.Context, notifySettingJSON string) error {
	var setting riskenslack.NotifySetting
	if err := json.Unmarshal([]byte(notifySettingJSON), &setting); err != nil {
		return err
	}

	locale := riskenslack.GetLocale(setting.Locale, s.defaultLocale)
	msg := riskenslack.GetOrgTestMessageText(locale)

	if setting.WebhookURL != "" {
		webhookMsg := &goslack.WebhookMessage{Text: msg}
		if setting.Data.Channel != "" {
			webhookMsg.Channel = setting.Data.Channel
		}
		if err := goslack.PostWebhook(setting.WebhookURL, webhookMsg); err != nil {
			return fmt.Errorf("failed to send slack(webhookurl): %w", err)
		}
	} else if setting.ChannelID != "" {
		if err := s.postMessageSlackWithRetry(ctx, setting.ChannelID, goslack.MsgOptionText(msg, false)); err != nil {
			return fmt.Errorf("failed to send slack(postmessage): %w", err)
		}
	}
	return nil
}
