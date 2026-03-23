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
			name:  "payload wrapped by json code fence",
			input: "```json\n{\"blocks\":[{\"type\":\"text\",\"text\":\"summary\"}]}\n```",
			want:  `{"blocks":[{"type":"text","text":"summary"}]}`,
		},
		{
			name:  "payload wrapped by plain code fence",
			input: "```\n{\"blocks\":[{\"type\":\"text\",\"text\":\"summary\"}]}\n```",
			want:  `{"blocks":[{"type":"text","text":"summary"}]}`,
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
