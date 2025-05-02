package ai

import (
	"testing"

	"github.com/ca-risken/core/pkg/model"
)

func TestGenerateFindingDataForAI(t *testing.T) {
	// テスト用の固定データ
	finding := &model.Finding{
		Score:       0.7,
		DataSource:  "AWS",
		Description: "Insecure configuration of security group",
		Data:        `{"resource_name": "sg-xxxxxx"}`,
	}

	recommend := &model.Recommend{
		Risk:           "risk description",
		Recommendation: "recommendation description",
	}

	tests := []struct {
		name      string
		finding   *model.Finding
		recommend *model.Recommend
		want      string
	}{
		{
			name:      "Only Finding",
			finding:   finding,
			recommend: nil,
			want: `The RISKEN tool detected the following issue related to cloud security.
Score: 
0.7

Type: 
AWS

Description: 
Insecure configuration of security group

ScanResult(json):
{"resource_name": "sg-xxxxxx"}
`,
		},
		{
			name:      "Finding and Recommend",
			finding:   finding,
			recommend: recommend,
			want: `The RISKEN tool detected the following issue related to cloud security.
Score: 
0.7

Type: 
AWS

Description: 
Insecure configuration of security group

ScanResult(json):
{"resource_name": "sg-xxxxxx"}

Detail: risk description

Recommendation: recommendation description
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateFindingDataForAI(tt.finding, tt.recommend)
			if got != tt.want {
				t.Errorf("generateFindingDataForAI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateCacheKeyForFinding(t *testing.T) {
	tests := []struct {
		name      string
		findingID uint64
		lang      string
		want      string
	}{
		{
			name:      "English language",
			findingID: 123456,
			lang:      "en",
			want:      "123456/en",
		},
		{
			name:      "Japanese language",
			findingID: 789012,
			lang:      "ja",
			want:      "789012/ja",
		},
		{
			name:      "Zero ID with language",
			findingID: 0,
			lang:      "en",
			want:      "0/en",
		},
		{
			name:      "Empty language",
			findingID: 123456,
			lang:      "",
			want:      "123456/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateCacheKeyForFinding(tt.findingID, tt.lang)
			if got != tt.want {
				t.Errorf("generateCacheKeyForFinding() = %v, want %v", got, tt.want)
			}
		})
	}
}
