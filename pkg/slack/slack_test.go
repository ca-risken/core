package slack

import (
	"testing"
)

func TestReplaceNotifySetting(t *testing.T) {
	cases := []struct {
		name       string
		existJSON  string
		updateJSON string
		want       NotifySetting
		wantErr    bool
	}{
		{
			name:       "Update WebhookURL clears ChannelID",
			existJSON:  `{"webhook_url":"https://old.com","channel_id":"ch1"}`,
			updateJSON: `{"webhook_url":"https://new.com","channel_id":"ch2"}`,
			want:       NotifySetting{WebhookURL: "https://new.com"},
		},
		{
			name:       "Update ChannelID clears WebhookURL",
			existJSON:  `{"webhook_url":"https://old.com"}`,
			updateJSON: `{"channel_id":"ch1"}`,
			want:       NotifySetting{ChannelID: "ch1"},
		},
		{
			name:       "No update returns exist",
			existJSON:  `{"webhook_url":"https://old.com","channel_id":"ch1"}`,
			updateJSON: `{}`,
			want:       NotifySetting{WebhookURL: "https://old.com", ChannelID: "ch1"},
		},
		{
			name:       "Invalid exist JSON",
			existJSON:  `{invalid`,
			updateJSON: `{}`,
			wantErr:    true,
		},
		{
			name:       "Invalid update JSON",
			existJSON:  `{}`,
			updateJSON: `{invalid`,
			wantErr:    true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ReplaceNotifySetting(c.existJSON, c.updateJSON)
			if (err != nil) != c.wantErr {
				t.Fatalf("ReplaceNotifySetting() error = %v, wantErr %v", err, c.wantErr)
			}
			if !c.wantErr && got != c.want {
				t.Fatalf("ReplaceNotifySetting() = %+v, want %+v", got, c.want)
			}
		})
	}
}

func TestMaskNotifySetting(t *testing.T) {
	cases := []struct {
		name             string
		notificationType string
		notifySetting    string
		want             string
		wantErr          bool
	}{
		{
			name:             "Slack with WebhookURL",
			notificationType: "slack",
			notifySetting:    `{"webhook_url":"https://example.com/hook","channel_id":""}`,
			want:             `{"webhook_url":"https://exam************","channel_id":"","data":{},"locale":""}`,
		},
		{
			name:             "Slack without WebhookURL",
			notificationType: "slack",
			notifySetting:    `{"channel_id":"ch1"}`,
			want:             `{"channel_id":"ch1"}`,
		},
		{
			name:             "Unknown type returns as-is",
			notificationType: "email",
			notifySetting:    `{"to":"user@example.com"}`,
			want:             `{"to":"user@example.com"}`,
		},
		{
			name:             "Slack invalid JSON",
			notificationType: "slack",
			notifySetting:    `{invalid`,
			wantErr:          true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := MaskNotifySetting(c.notificationType, c.notifySetting)
			if (err != nil) != c.wantErr {
				t.Fatalf("MaskNotifySetting() error = %v, wantErr %v", err, c.wantErr)
			}
			if !c.wantErr && got != c.want {
				t.Fatalf("MaskNotifySetting() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestValidateNewNotifySetting(t *testing.T) {
	cases := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name:  "OK - webhook_url",
			input: `{"webhook_url":"https://hooks.slack.com/xxx"}`,
		},
		{
			name:  "OK - channel_id",
			input: `{"channel_id":"C12345"}`,
		},
		{
			name:    "NG - both empty",
			input:   `{}`,
			wantErr: true,
		},
		{
			name:    "NG - invalid JSON",
			input:   `{invalid`,
			wantErr: true,
		},
		{
			name:    "NG - invalid URL",
			input:   `{"webhook_url":"not-a-url"}`,
			wantErr: true,
		},
		{
			name:    "NG - not string",
			input:   123,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateNewNotifySetting(c.input)
			if (err != nil) != c.wantErr {
				t.Fatalf("ValidateNewNotifySetting() error = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}

func TestValidateExistingNotifySetting(t *testing.T) {
	cases := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name:  "OK - valid webhook_url",
			input: `{"webhook_url":"https://hooks.slack.com/xxx"}`,
		},
		{
			name:  "OK - empty webhook_url",
			input: `{"channel_id":"C12345"}`,
		},
		{
			name:  "OK - all empty",
			input: `{}`,
		},
		{
			name:    "NG - invalid URL",
			input:   `{"webhook_url":"not-a-url"}`,
			wantErr: true,
		},
		{
			name:    "NG - invalid JSON",
			input:   `{invalid`,
			wantErr: true,
		},
		{
			name:    "NG - not string",
			input:   42,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateExistingNotifySetting(c.input)
			if (err != nil) != c.wantErr {
				t.Fatalf("ValidateExistingNotifySetting() error = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}

func TestMaskRight(t *testing.T) {
	cases := []struct {
		name  string
		input string
		num   int
		want  string
	}{
		{
			name:  "Mask half",
			input: "abcdefgh",
			num:   4,
			want:  "abcd****",
		},
		{
			name:  "Mask all",
			input: "abcd",
			num:   0,
			want:  "****",
		},
		{
			name:  "Mask none",
			input: "abcd",
			num:   4,
			want:  "abcd",
		},
		{
			name:  "Empty string",
			input: "",
			num:   0,
			want:  "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := maskRight(c.input, c.num)
			if got != c.want {
				t.Fatalf("maskRight() = %v, want %v", got, c.want)
			}
		})
	}
}
