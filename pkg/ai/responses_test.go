package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/responses"
)

type responseTestSender struct{ messages []string }

func (s *responseTestSender) Send(message string) error {
	s.messages = append(s.messages, message)
	return nil
}

func TestResponsesReasoning(t *testing.T) {
	cases := []struct {
		name       string
		model      string
		effort     openai.ReasoningEffort
		stream     bool
		wantEffort string
	}{
		{name: "luna default", model: "gpt-6-luna", wantEffort: "none"},
		{name: "luna analysis", model: "gpt-6-luna", effort: openai.ReasoningEffortMedium, wantEffort: "medium"},
		{name: "luna report", model: "gpt-6-luna", effort: openai.ReasoningEffortHigh, wantEffort: "high"},
		{name: "luna stream", model: "gpt-6-luna", stream: true, wantEffort: "none"},
		{name: "mini rollback", model: "gpt-5.4-mini"},
		{name: "mini stream rollback", model: "gpt-5.4-mini", stream: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var request struct {
				Model     string `json:"model"`
				Stream    bool   `json:"stream"`
				Reasoning *struct {
					Effort string `json:"effort"`
				} `json:"reasoning"`
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					t.Error(err)
				}
				if tc.stream {
					w.Header().Set("Content-Type", "text/event-stream")
					fmt.Fprint(w, "event: response.output_text.delta\ndata: {\"type\":\"response.output_text.delta\",\"delta\":\"ok\",\"output_index\":0,\"content_index\":0}\n\n")
				} else {
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprint(w, `{"id":"resp_test","object":"response","status":"completed","output":[]}`)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithAPIKey("test"), option.WithBaseURL(server.URL), option.WithMaxRetries(0))
			input := responses.ResponseNewParamsInputUnion{OfString: openai.String("test")}
			if tc.stream {
				sender := &responseTestSender{}
				_, err := responsesStreamingAPI(context.Background(), &client, tc.model, "test", input, nil, sender, func(s string) string { return s }, logging.NewLogger())
				if err != nil {
					t.Fatal(err)
				}
				if len(sender.messages) != 1 || sender.messages[0] != "ok" {
					t.Fatalf("stream messages = %v", sender.messages)
				}
			} else {
				ai := &AIClient{openaiClient: &client}
				if _, err := ai.callResponsesAPI(context.Background(), tc.model, "test", input, nil, tc.effort); err != nil {
					t.Fatal(err)
				}
			}
			if request.Model != tc.model || request.Stream != tc.stream {
				t.Fatalf("unexpected model/stream: %+v", request)
			}
			if tc.wantEffort == "" {
				if request.Reasoning != nil {
					t.Fatalf("reasoning must remain omitted: %+v", request.Reasoning)
				}
			} else if request.Reasoning == nil || request.Reasoning.Effort != tc.wantEffort {
				t.Fatalf("reasoning = %+v, want %s", request.Reasoning, tc.wantEffort)
			}
		})
	}
}
