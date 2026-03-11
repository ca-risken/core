package org_alert

import (
	"context"
	"encoding/json"
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
