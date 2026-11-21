package ai

import (
	"context"
	"errors"
	"strings"
	"testing"

	aiservice "github.com/ca-risken/core/pkg/ai"
	"github.com/ca-risken/core/proto/ai"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
)

type classifyContentFunc func(context.Context, *ai.ClassifySensitiveContentRequest) (*ai.ClassifySensitiveContentResponse, error)

func (f classifyContentFunc) Classify(ctx context.Context, req *ai.ClassifySensitiveContentRequest) (*ai.ClassifySensitiveContentResponse, error) {
	return f(ctx, req)
}

func TestClassifySensitiveContent(t *testing.T) {
	t.Parallel()

	const providerError = "provider echoed the sample"
	classified := &ai.ClassifySensitiveContentResponse{Sensitivity: ai.Sensitivity_SENSITIVITY_FALSE}
	tests := []struct {
		name        string
		input       string
		image       []byte
		fileName    string
		prompt      string
		disabled    bool
		providerErr error
		canceled    bool
		want        codes.Code
	}{
		{name: "bounded sample", input: strings.Repeat("a", 32768), want: codes.OK},
		{name: "image byte limit inclusive", image: make([]byte, 2<<20), want: codes.OK},
		{name: "image over byte limit", image: make([]byte, (2<<20)+1), want: codes.InvalidArgument},
		{name: "image and text are mutually exclusive", input: "sample", image: []byte{1}, want: codes.InvalidArgument},
		{name: "invalid image", image: []byte{1}, providerErr: aiservice.ErrInvalidSensitiveContentImage, want: codes.InvalidArgument},
		{name: "valid file name", input: "sample", fileName: "test_customers.csv", want: codes.OK},
		{name: "file name byte limit", input: "sample", fileName: strings.Repeat("a", 1025), want: codes.InvalidArgument},
		{name: "custom prompt byte limit inclusive", input: "sample", prompt: strings.Repeat("a", 4096), want: codes.OK},
		{name: "custom prompt over byte limit", input: "sample", prompt: strings.Repeat("あ", 1366), want: codes.InvalidArgument},
		{name: "missing OpenAI token", disabled: true, want: codes.FailedPrecondition},
		{name: "blank input", input: " \n ", want: codes.InvalidArgument},
		{name: "input byte limit", input: strings.Repeat("あ", 10923), want: codes.InvalidArgument},
		{name: "provider error", input: "sample", providerErr: errors.New(providerError), want: codes.Unavailable},
		{name: "canceled request", input: "sample", canceled: true, want: codes.Canceled},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := NewSensitiveContentService("")
			calls := 0
			if !tt.disabled {
				service.classifier = classifyContentFunc(func(ctx context.Context, req *ai.ClassifySensitiveContentRequest) (*ai.ClassifySensitiveContentResponse, error) {
					calls++
					if req.GetText() != tt.input || req.GetFileName() != tt.fileName || req.GetCustomPrompt() != tt.prompt || len(req.GetImage()) != len(tt.image) {
						t.Error("classifier did not receive the request fields")
					}
					if _, ok := ctx.Deadline(); !ok {
						t.Error("classifier context needs a deadline")
					}
					if tt.providerErr != nil {
						return nil, tt.providerErr
					}
					return classified, nil
				})
			}
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if tt.canceled {
				cancel()
			}
			got, err := service.ClassifySensitiveContent(ctx, &ai.ClassifySensitiveContentRequest{
				Text: tt.input, FileName: tt.fileName, Image: tt.image, CustomPrompt: tt.prompt,
			})
			if code := status.Code(err); code != tt.want {
				t.Fatalf("status = %v, want %v", code, tt.want)
			}
			var want *ai.ClassifySensitiveContentResponse
			if tt.want == codes.OK {
				want = classified
			}
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("response mismatch (-want +got):\n%s", diff)
			}
			if (tt.want == codes.OK || tt.providerErr != nil || tt.canceled) != (calls == 1) {
				t.Errorf("unexpected classifier call count = %d", calls)
			}
			if err != nil && strings.Contains(err.Error(), providerError) {
				t.Error("gRPC status included the upstream error")
			}
		})
	}
}
