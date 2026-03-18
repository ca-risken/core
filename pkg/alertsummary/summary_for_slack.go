package alertsummary

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	BlockTypeText = "text"
	BlockTypeLink = "link"
)

var slackMrkdwnReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
)

type Payload struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Type  string `json:"type"`
	Text  string `json:"text,omitempty"`
	Label string `json:"label,omitempty"`
	URL   string `json:"url,omitempty"`
}

func Normalize(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if payload, ok := parsePayload(raw); ok {
		if normalized, ok := marshalPayload(sanitizePayload(payload)); ok {
			return normalized
		}
	}
	return ""
}

func RenderSlack(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if payload, ok := parsePayload(raw); ok {
		lines := []string{}
		for _, block := range sanitizePayload(payload).Blocks {
			switch block.Type {
			case BlockTypeText:
				if text := escapeSlackMrkdwn(strings.TrimSpace(block.Text)); text != "" {
					lines = append(lines, text)
				}
			case BlockTypeLink:
				url := sanitizeSlackLinkURL(strings.TrimSpace(block.URL))
				if url == "" {
					continue
				}
				label := sanitizeSlackLinkLabel(strings.TrimSpace(block.Label))
				if label == "" {
					label = sanitizeSlackLinkLabel(url)
				}
				lines = append(lines, fmt.Sprintf("<%s|%s>", url, label))
			}
		}
		if len(lines) > 0 {
			return strings.Join(lines, "\n")
		}
	}
	return ""
}

func parsePayload(raw string) (Payload, bool) {
	var payload Payload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return Payload{}, false
	}
	if len(payload.Blocks) == 0 {
		return Payload{}, false
	}
	return payload, true
}

func sanitizePayload(payload Payload) Payload {
	blocks := make([]Block, 0, len(payload.Blocks))
	for _, block := range payload.Blocks {
		switch block.Type {
		case BlockTypeText:
			text := strings.TrimSpace(block.Text)
			if text == "" {
				continue
			}
			blocks = append(blocks, Block{
				Type: BlockTypeText,
				Text: text,
			})
		case BlockTypeLink:
			url := sanitizeSlackLinkURL(strings.TrimSpace(block.URL))
			if url == "" {
				continue
			}
			label := strings.TrimSpace(block.Label)
			if label == "" {
				label = url
			}
			blocks = append(blocks, Block{
				Type:  BlockTypeLink,
				Label: label,
				URL:   url,
			})
		}
	}
	return Payload{Blocks: blocks}
}

func marshalPayload(payload Payload) (string, bool) {
	if len(payload.Blocks) == 0 {
		return "", false
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return "", false
	}
	return string(b), true
}

func isSlackSafeURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func sanitizeSlackLinkURL(url string) string {
	if !isSlackSafeURL(url) {
		return ""
	}
	if strings.ContainsAny(url, "<>|") {
		return ""
	}
	return url
}

func escapeSlackMrkdwn(text string) string {
	return slackMrkdwnReplacer.Replace(text)
}

func sanitizeSlackLinkLabel(label string) string {
	label = escapeSlackMrkdwn(label)
	return strings.ReplaceAll(label, "|", "¦")
}
