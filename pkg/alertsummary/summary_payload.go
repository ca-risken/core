package alertsummary

import (
	"encoding/json"
	"strings"
)

const (
	BlockTypeText = "text"
	BlockTypeLink = "link"
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
	if payload, ok := Parse(raw); ok {
		if normalized, ok := marshalPayload(payload); ok {
			return normalized
		}
	}
	return ""
}

func Parse(raw string) (Payload, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Payload{}, false
	}
	raw = stripCodeFence(raw)
	return parsePayload(raw)
}

func stripCodeFence(raw string) string {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	if len(lines) < 3 {
		return raw
	}
	if !strings.HasPrefix(strings.TrimSpace(lines[0]), "```") {
		return raw
	}
	if strings.TrimSpace(lines[len(lines)-1]) != "```" {
		return raw
	}
	return strings.TrimSpace(strings.Join(lines[1:len(lines)-1], "\n"))
}

func parsePayload(raw string) (Payload, bool) {
	var payload Payload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return Payload{}, false
	}
	payload = sanitizePayload(payload)
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
			url := sanitizeLinkURL(strings.TrimSpace(block.URL))
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

func sanitizeLinkURL(url string) string {
	if !isSlackSafeURL(url) {
		return ""
	}
	if strings.ContainsAny(url, "<>|") {
		return ""
	}
	return url
}
