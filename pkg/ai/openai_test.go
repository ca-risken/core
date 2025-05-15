package ai

import (
	"context"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/coocood/freecache"
)

func TestGetSetAICache(t *testing.T) {
	logger := logging.NewLogger()
	tests := []struct {
		name        string
		key         string
		value       string
		wantCache   bool
		wantCached  string
		presetCache bool
	}{
		{
			name:       "Basic cache operation",
			key:        "test-key-1",
			value:      "test-value-1",
			wantCache:  true,
			wantCached: "test-value-1",
		},
		{
			name:       "Empty value",
			key:        "test-key-2",
			value:      "",
			wantCache:  true,
			wantCached: "",
		},
		{
			name:       "Long key",
			key:        "test-key-with-very-long-string-that-exceeds-normal-limits-but-should-still-work-with-md5-hash",
			value:      "test-value-3",
			wantCache:  true,
			wantCached: "test-value-3",
		},
		{
			name:        "Get non-existing cache",
			key:         "not-exist-key",
			presetCache: false,
			wantCache:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &AIClient{
				cache:  freecache.NewCache(CACHE_SIZE),
				logger: logger,
			}

			ctx := context.Background()
			if tt.presetCache && tt.value != "" {
				if err := client.setAICache(tt.key, tt.value); err != nil {
					t.Fatalf("Failed to preset cache: %v", err)
				}
			}
			if !tt.presetCache && tt.wantCache {
				if err := client.setAICache(tt.key, tt.value); err != nil {
					t.Errorf("setAICache() error = %v", err)
				}
			}

			got := client.getAICache(ctx, tt.key)
			if tt.wantCache {
				if got != tt.wantCached {
					t.Errorf("getAICache() = %v, want %v", got, tt.wantCached)
				}
			} else {
				if got != "" {
					t.Errorf("getAICache() for non-existing key returned %v, want empty string", got)
				}
			}
		})
	}
}

func TestGenerateCacheKey(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantLen int
	}{
		{
			name:    "Normal string",
			content: "test-content",
			wantLen: 16, // MD5ハッシュは16バイト
		},
		{
			name:    "Empty string",
			content: "",
			wantLen: 16,
		},
		{
			name:    "Long string",
			content: "this-is-a-very-long-string-that-should-still-produce-a-fixed-length-md5-hash-value-no-matter-how-long-the-input-is",
			wantLen: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateCacheKey(tt.content)
			if len(got) != tt.wantLen {
				t.Errorf("generateCacheKey() result length = %v, want %v", len(got), tt.wantLen)
			}

			gotAgain := generateCacheKey(tt.content)
			for i := 0; i < len(got); i++ {
				if got[i] != gotAgain[i] {
					t.Errorf("generateCacheKey() is not deterministic at position %d", i)
					break
				}
			}
		})
	}
}
