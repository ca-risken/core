package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	is "github.com/go-ozzo/ozzo-validation/v4/is"
	goslack "github.com/slack-go/slack"
)

const (
	localeJa = "ja"
	localeEn = "en"

	testNotificationMessageJa    = "RISKENからのテスト通知です"
	testNotificationMessageEn    = "This is a test notification from RISKEN"
	testOrgNotificationMessageJa = "RISKENからのテスト通知です（Organization）"
	testOrgNotificationMessageEn = "This is a test notification from RISKEN (Organization)"
)

// NotifySetting is the Slack notification setting shared by alert and org_alert services.
type NotifySetting struct {
	WebhookURL string       `json:"webhook_url"`
	ChannelID  string       `json:"channel_id"`
	Data       NotifyOption `json:"data"`
	Locale     string       `json:"locale"`
}

// NotifyOption is the optional Slack notification parameters.
type NotifyOption struct {
	Channel string `json:"channel,omitempty"`
	Message string `json:"message,omitempty"`
}

// ReplaceNotifySetting merges update into exist, enforcing mutual exclusivity of WebhookURL and ChannelID.
func ReplaceNotifySetting(existJSON, updateJSON string) (NotifySetting, error) {
	var update NotifySetting
	if err := json.Unmarshal([]byte(updateJSON), &update); err != nil {
		return NotifySetting{}, err
	}
	var exist NotifySetting
	if err := json.Unmarshal([]byte(existJSON), &exist); err != nil {
		return NotifySetting{}, err
	}

	// webhookURL and ChannelID are mutually exclusive
	if update.WebhookURL != "" {
		update.ChannelID = ""
		return update, nil
	}
	if update.ChannelID != "" {
		update.WebhookURL = ""
		return update, nil
	}

	// No update options
	return exist, nil
}

// MaskNotifySetting masks sensitive fields in the notify setting JSON.
func MaskNotifySetting(notificationType, notifySetting string) (string, error) {
	switch notificationType {
	case "slack":
		var setting NotifySetting
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

// ValidateNewNotifySetting validates notify setting for new notification.
// Compatible with validation.By (func(any) error).
func ValidateNewNotifySetting(value any) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("notify setting is not string, %v", value)
	}
	var setting NotifySetting
	if err := json.Unmarshal([]byte(s), &setting); err != nil {
		return fmt.Errorf("invalid json, %w", err)
	}
	if strings.TrimSpace(setting.WebhookURL) == "" && strings.TrimSpace(setting.ChannelID) == "" {
		return errors.New("required webhook_url or channel_id in json")
	}
	if err := validation.Validate(strings.TrimSpace(setting.WebhookURL), is.URL); err != nil {
		return err
	}
	return nil
}

// ValidateExistingNotifySetting validates notify setting for existing notification.
// Compatible with validation.By (func(any) error).
func ValidateExistingNotifySetting(value any) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("notify setting is not string, %v", value)
	}
	var setting NotifySetting
	if err := json.Unmarshal([]byte(s), &setting); err != nil {
		return fmt.Errorf("invalid json, %w", err)
	}
	if strings.TrimSpace(setting.WebhookURL) != "" {
		if err := validation.Validate(strings.TrimSpace(setting.WebhookURL), validation.Required, is.URL); err != nil {
			return err
		}
	}
	return nil
}

// getLocale returns a valid locale string, falling back to defaultLocale.
func getLocale(settingLocale, defaultLocale string) string {
	switch settingLocale {
	case localeJa:
		return localeJa
	case localeEn:
		return localeEn
	default:
		return defaultLocale
	}
}

// getTestMessageText returns the test notification message for the given locale.
func getTestMessageText(locale string) string {
	switch locale {
	case localeJa:
		return testNotificationMessageJa
	default:
		return testNotificationMessageEn
	}
}

// getOrgTestMessageText returns the Organization test notification message for the given locale.
func getOrgTestMessageText(locale string) string {
	switch locale {
	case localeJa:
		return testOrgNotificationMessageJa
	default:
		return testOrgNotificationMessageEn
	}
}

// SendTestNotification sends a test Slack notification using the given settings.
func SendTestNotification(slackClient *goslack.Client, notifySettingJSON, defaultLocale string) error {
	return sendTestNotificationWithMessage(slackClient, notifySettingJSON, defaultLocale, getTestMessageText)
}

// SendOrgTestNotification sends a test Slack notification for Organization using the given settings.
func SendOrgTestNotification(slackClient *goslack.Client, notifySettingJSON, defaultLocale string) error {
	return sendTestNotificationWithMessage(slackClient, notifySettingJSON, defaultLocale, getOrgTestMessageText)
}

func sendTestNotificationWithMessage(slackClient *goslack.Client, notifySettingJSON, defaultLocale string, msgFunc func(string) string) error {
	var setting NotifySetting
	if err := json.Unmarshal([]byte(notifySettingJSON), &setting); err != nil {
		return err
	}

	locale := getLocale(setting.Locale, defaultLocale)
	msg := msgFunc(locale)

	if setting.WebhookURL != "" {
		webhookMsg := &goslack.WebhookMessage{Text: msg}
		if setting.Data.Channel != "" {
			webhookMsg.Channel = setting.Data.Channel
		}
		if err := goslack.PostWebhook(setting.WebhookURL, webhookMsg); err != nil {
			return fmt.Errorf("failed to send slack(webhookurl): %w", err)
		}
	} else if setting.ChannelID != "" {
		if _, _, err := slackClient.PostMessage(setting.ChannelID, goslack.MsgOptionText(msg, false)); err != nil {
			return fmt.Errorf("failed to send slack(postmessage): %w", err)
		}
	}
	return nil
}

func maskRight(s string, num int) string {
	rs := []rune(s)
	for i := num; i < len(rs); i++ {
		rs[i] = '*'
	}
	return string(rs)
}
