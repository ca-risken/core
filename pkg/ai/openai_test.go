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
