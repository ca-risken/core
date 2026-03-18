package alertsummary

import "testing"

func TestNormalize(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "existing payload",
			input: `{"blocks":[{"type":"text","text":"summary"},{"type":"link","label":"GitHub","url":"https://example.com"}]}`,
			want:  `{"blocks":[{"type":"text","text":"summary"},{"type":"link","label":"GitHub","url":"https://example.com"}]}`,
		},
		{
			name:  "invalid payload",
			input: "summary",
			want:  "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Normalize(c.input)
			if got != c.want {
				t.Fatalf("Unexpected normalized value: got=%q want=%q", got, c.want)
			}
		})
	}
}

func TestRenderSlack(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "payload",
			input: `{"blocks":[{"type":"text","text":"summary"},{"type":"link","label":"GitHub","url":"https://example.com"}]}`,
			want:  "summary\n<https://example.com|GitHub>",
		},
		{
			name:  "escape mrkdwn in text and label",
			input: `{"blocks":[{"type":"text","text":"notify <!here> & review <this>"},{"type":"link","label":"Git|Hub > docs","url":"https://example.com"}]}`,
			want:  "notify &lt;!here&gt; &amp; review &lt;this&gt;\n<https://example.com|Git¦Hub &gt; docs>",
		},
		{
			name:  "drop unsafe url in link block",
			input: `{"blocks":[{"type":"text","text":"summary"},{"type":"link","label":"malicious","url":"https://example.com/a><!channel>"}]}`,
			want:  "summary",
		},
		{
			name:  "invalid payload",
			input: "確認: [GitHubリンク](https://example.com/path)",
			want:  "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := RenderSlack(c.input)
			if got != c.want {
				t.Fatalf("Unexpected rendered value: got=%q want=%q", got, c.want)
			}
		})
	}
}
