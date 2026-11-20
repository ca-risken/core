package ai

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ca-risken/core/proto/ai"
	"github.com/google/go-cmp/cmp"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"google.golang.org/protobuf/testing/protocmp"
)

type decisionsTestClient func(*http.Request) (*http.Response, error)

func (f decisionsTestClient) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestClassifySensitiveContentWithDecisions(t *testing.T) {
	t.Parallel()

	const sample = "client_id,contact\n42,[REDACTED]"
	const fileName = "test_customers.csv"
	categoryNames := []string{"pii", "payment_card", "credential", "client_data", "other"}
	kindNames := []string{"log", "tabular_data", "structured_record", "document", "source_code", "configuration", "image", "other"}
	tests := []struct {
		name        string
		first       string
		categories  []string
		kindAnswer  string
		testData    string
		imageFormat string
		want        *ai.ClassifySensitiveContentResponse
		failed      bool
	}{
		{
			name:       "sensitive sample with multiple categories",
			first:      "0.8",
			categories: []string{"0.8", "0.1", "0.2", "0.91", "0.1"},
			kindAnswer: `{"type":"choice","name":"content_kind","choice":"tabular_data","confidence":0.9}`,
			testData:   "0.9",
			want: &ai.ClassifySensitiveContentResponse{
				Sensitivity: ai.Sensitivity_SENSITIVITY_TRUE,
				Categories: []ai.SensitiveContentCategory{
					ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_PII,
					ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_CLIENT_DATA,
				},
				ContentKind: ai.ContentKind_CONTENT_KIND_TABULAR_DATA,
				TestData:    ai.TestData_TEST_DATA_TRUE,
			},
		},
		{
			name:        "sensitive PNG with file name for test-data decision",
			imageFormat: "png",
			first:       "0.95",
			categories:  []string{"0", "0", "0", "0", "0.95"},
			kindAnswer:  `{"type":"choice","name":"content_kind","choice":"image","confidence":0.95}`,
			testData:    "0.9",
			want: &ai.ClassifySensitiveContentResponse{
				Sensitivity: ai.Sensitivity_SENSITIVITY_TRUE,
				Categories:  []ai.SensitiveContentCategory{ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_OTHER},
				ContentKind: ai.ContentKind_CONTENT_KIND_IMAGE,
				TestData:    ai.TestData_TEST_DATA_TRUE,
			},
		},
		{
			name:        "non-sensitive JPEG skips classification",
			imageFormat: "jpeg",
			first:       "0.2",
			want:        &ai.ClassifySensitiveContentResponse{Sensitivity: ai.Sensitivity_SENSITIVITY_FALSE},
		},
		{
			name:  "non-sensitive sample skips classification",
			first: "0.2",
			want:  &ai.ClassifySensitiveContentResponse{Sensitivity: ai.Sensitivity_SENSITIVITY_FALSE},
		},
		{
			name:  "uncertain sensitivity skips classification",
			first: "0.65",
			want:  &ai.ClassifySensitiveContentResponse{Sensitivity: ai.Sensitivity_SENSITIVITY_UNKNOWN},
		},
		{
			name:  "refusal skips classification",
			first: "refusal",
			want:  &ai.ClassifySensitiveContentResponse{Sensitivity: ai.Sensitivity_SENSITIVITY_UNKNOWN},
		},
		{
			name:       "later refusal and uncertain kind preserve detection",
			first:      "0.95",
			categories: []string{"0.1", "0.1", "0.8", "refusal", "0.3"},
			kindAnswer: `{"type":"choice","name":"content_kind","choice":"log","confidence":0.5}`,
			testData:   "refusal",
			want: &ai.ClassifySensitiveContentResponse{
				Sensitivity: ai.Sensitivity_SENSITIVITY_TRUE,
				Categories:  []ai.SensitiveContentCategory{ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_CREDENTIAL},
			},
		},
		{
			name:       "real data despite a test file name",
			first:      "0.95",
			categories: []string{"0.8", "0", "0", "0", "0"},
			kindAnswer: `{"type":"choice","name":"content_kind","choice":"tabular_data","confidence":0.9}`,
			testData:   "0.2",
			want: &ai.ClassifySensitiveContentResponse{
				Sensitivity: ai.Sensitivity_SENSITIVITY_TRUE,
				Categories:  []ai.SensitiveContentCategory{ai.SensitiveContentCategory_SENSITIVE_CONTENT_CATEGORY_PII},
				ContentKind: ai.ContentKind_CONTENT_KIND_TABULAR_DATA,
				TestData:    ai.TestData_TEST_DATA_FALSE,
			},
		},
		{
			name:       "unexpected kind from provider",
			first:      "0.95",
			categories: []string{"0.8", "0", "0", "0", "0"},
			kindAnswer: `{"type":"choice","name":"content_kind","choice":"unexpected","confidence":0.9}`,
			testData:   "0.5",
			failed:     true,
		},
		{name: "missing probability", first: "null", failed: true},
		{name: "out-of-range probability", first: "1.01", failed: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			input := &ai.ClassifySensitiveContentRequest{Text: sample, FileName: fileName}
			var imagePart decisionInputPart
			if tt.imageFormat != "" {
				input.Text = ""
				input.Image = encodeSensitiveContentTestImage(t, tt.imageFormat, 8, 4)
				imagePart = decisionInputPart{
					Type:     "input_image",
					ImageURL: "data:image/" + tt.imageFormat + ";base64," + base64.StdEncoding.EncodeToString(input.Image),
				}
			}
			replies := [][]string{{predicateAnswerJSON("sensitive", tt.first)}}
			if len(tt.categories) > 0 {
				var details []string
				for i, probability := range tt.categories {
					details = append(details, predicateAnswerJSON(categoryNames[i], probability))
				}
				replies = append(replies, append(details, tt.kindAnswer, predicateAnswerJSON("test_data", tt.testData)))
			}
			calls := 0
			transport := decisionsTestClient(func(req *http.Request) (*http.Response, error) {
				if calls >= len(replies) {
					t.Fatal("unexpected extra Decisions API call")
				}
				if req.Method != http.MethodPost || req.URL.Path != "/v1/decisions" || req.Header.Get("Authorization") != "Bearer test-key" {
					t.Error("unexpected Decisions API request route or authentication")
				}
				var request struct {
					Model              string             `json:"model"`
					Input              json.RawMessage    `json:"input"`
					Questions          []decisionQuestion `json:"questions"`
					PromptCacheOptions map[string]string  `json:"prompt_cache_options"`
				}
				if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
					return nil, err
				}
				if request.Model != "gpt-6-luna" || request.PromptCacheOptions["mode"] != "explicit" {
					t.Error("unexpected model or caching mode")
				}
				wantNames := []string{"sensitive"}
				if calls == 1 || tt.imageFormat != "" {
					var gotInput []decisionInputMessage
					if err := json.Unmarshal(request.Input, &gotInput); err != nil {
						t.Fatalf("decode decision evidence: %v", err)
					}
					var parts []decisionInputPart
					if calls == 1 {
						wantNames = append(append([]string{}, categoryNames...), "content_kind", "test_data")
						parts = append(parts, decisionInputPart{Type: "input_json", Data: new(contentSample{FileName: fileName, Text: input.Text})})
					}
					if tt.imageFormat != "" {
						parts = append(parts, imagePart)
					}
					wantInput := []decisionInputMessage{{Role: "user", Content: parts}}
					if diff := cmp.Diff(wantInput, gotInput); diff != "" {
						t.Errorf("decision evidence mismatch (-want +got):\n%s", diff)
					}
				} else {
					var gotInput string
					if err := json.Unmarshal(request.Input, &gotInput); err != nil || gotInput != sample {
						t.Error("sensitivity request must contain only the sample text")
					}
				}
				var gotNames []string
				for _, question := range request.Questions {
					gotNames = append(gotNames, question.Name)
					if question.Name != "content_kind" {
						if question.Type != "predicate" || question.Instructions == "" || len(question.Choices) != 0 {
							t.Errorf("invalid predicate question %q", question.Name)
						}
						if question.Name != "test_data" && !strings.Contains(question.Instructions, "including fictional or synthetic samples") {
							t.Errorf("predicate %q must consider sensitive values in synthetic data", question.Name)
						}
						continue
					}
					if question.Type != "choice" || question.Instructions == "" {
						t.Error("expected a content-kind choice question")
					}
					var gotKinds []string
					for _, choice := range question.Choices {
						gotKinds = append(gotKinds, choice.Value)
					}
					if diff := cmp.Diff(kindNames, gotKinds); diff != "" {
						t.Errorf("content kind choices mismatch (-want +got):\n%s", diff)
					}
				}
				if diff := cmp.Diff(wantNames, gotNames); diff != "" {
					t.Errorf("question names mismatch (-want +got):\n%s", diff)
				}
				answers := replies[calls]
				calls++
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"answers":[` + strings.Join(answers, ",") + `]}`)),
				}, nil
			})
			client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL("https://example.test/v1/"), option.WithHTTPClient(transport))
			classifier := &SensitiveContentClassifier{client: &client}

			got, err := classifier.Classify(t.Context(), input)
			if (err != nil) != tt.failed {
				t.Fatalf("Classify() error = %v, want error = %v", err, tt.failed)
			}
			if calls != len(replies) {
				t.Errorf("Decisions API calls = %d, want %d", calls, len(replies))
			}
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("classification mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func predicateAnswerJSON(name, probability string) string {
	if probability == "refusal" {
		return fmt.Sprintf(`{"type":"refusal","name":%q}`, name)
	}
	return fmt.Sprintf(`{"type":"predicate","name":%q,"probability":%s}`, name, probability)
}
