package finding

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ca-risken/core/pkg/model"
	"github.com/google/go-cmp/cmp"
)

func TestEvaluateExploitation(t *testing.T) {
	tests := []struct {
		name  string
		input *Exploitation
		want  *AssessmentDetail
	}{
		{
			name:  "nil",
			input: &Exploitation{},
			want: &AssessmentDetail{
				Result: Ptr(TRIAGE_UNKNOWN),
				Score:  Ptr(float32(0)),
			},
		},
		{
			name: "false_cve",
			input: &Exploitation{
				HasCVE: Ptr(false),
			},
			want: &AssessmentDetail{
				Result: Ptr(EXPLOITATION_RESULT_NONE),
				Score:  Ptr(float32(-0.1)),
			},
		},
		{
			name: "true_cve_no_flags",
			input: &Exploitation{
				HasCVE: Ptr(true),
			},
			want: &AssessmentDetail{
				Result: Ptr(EXPLOITATION_RESULT_NONE),
				Score:  Ptr(float32(-0.1)),
			},
		},
		{
			name: "with_kev",
			input: &Exploitation{
				HasCVE: Ptr(true),
				HasKEV: Ptr(true),
			},
			want: &AssessmentDetail{
				Result: Ptr(EXPLOITATION_RESULT_ACTIVE),
				Score:  Ptr(float32(0.0)),
			},
		},
		{
			name: "publicpoc_high_epss",
			input: &Exploitation{
				HasCVE:    Ptr(true),
				PublicPOC: Ptr(true),
				EpssScore: Ptr(float32(0.01)),
			},
			want: &AssessmentDetail{
				Result: Ptr(EXPLOITATION_RESULT_ACTIVE),
				Score:  Ptr(float32(0.0)),
			},
		},
		{
			name: "publicpoc_threshold_epss",
			input: &Exploitation{
				HasCVE:    Ptr(true),
				PublicPOC: Ptr(true),
				EpssScore: Ptr(float32(0.009)),
			},
			want: &AssessmentDetail{
				Result: Ptr(EXPLOITATION_RESULT_PUBLIC_POC),
				Score:  Ptr(float32(-0.1)),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := evaluateExploitation(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestEvaluateSystemExposure(t *testing.T) {
	tests := []struct {
		name  string
		input *SystemExposure
		want  *AssessmentDetail
	}{
		{
			name:  "nil",
			input: &SystemExposure{},
			want: &AssessmentDetail{
				Result: Ptr(TRIAGE_UNKNOWN),
				Score:  Ptr(float32(0)),
			},
		},
		{
			name: "missing_access",
			input: &SystemExposure{
				PublicFacing: Ptr("OPEN"),
			},
			want: &AssessmentDetail{
				Result: Ptr(TRIAGE_UNKNOWN),
				Score:  Ptr(float32(0)),
			},
		},
		{
			name: "missing_public",
			input: &SystemExposure{
				AccessControl: Ptr(ACCESS_CONTROL_NONE),
			},
			want: &AssessmentDetail{
				Result: Ptr(TRIAGE_UNKNOWN),
				Score:  Ptr(float32(0)),
			},
		},
		{
			name: "unknown_access",
			input: &SystemExposure{
				PublicFacing:  Ptr(PUBLIC_FACING_OPEN),
				AccessControl: Ptr(TRIAGE_UNKNOWN),
			},
			want: &AssessmentDetail{
				Result: Ptr(TRIAGE_UNKNOWN),
				Score:  Ptr(float32(0)),
			},
		},
		{
			name: "open",
			input: &SystemExposure{
				PublicFacing:  Ptr(PUBLIC_FACING_OPEN),
				AccessControl: Ptr(ACCESS_CONTROL_NONE),
			},
			want: &AssessmentDetail{
				Result: Ptr(SYSTEM_EXPOSURE_OPEN),
				Score:  Ptr(float32(0)),
			},
		},
		{
			name: "controlled_from_open",
			input: &SystemExposure{
				PublicFacing:  Ptr(PUBLIC_FACING_OPEN),
				AccessControl: Ptr(ACCESS_CONTROL_AUTHENTICATED),
			},
			want: &AssessmentDetail{
				Result: Ptr(SYSTEM_EXPOSURE_CONTROLLED),
				Score:  Ptr(float32(-0.1)),
			},
		},
		{
			name: "internal_controlled",
			input: &SystemExposure{
				PublicFacing:  Ptr(PUBLIC_FACING_INTERNAL),
				AccessControl: Ptr(ACCESS_CONTROL_NONE),
			},
			want: &AssessmentDetail{
				Result: Ptr(SYSTEM_EXPOSURE_CONTROLLED),
				Score:  Ptr(float32(-0.1)),
			},
		},
		{
			name: "internal_small",
			input: &SystemExposure{
				AccessControl: Ptr(ACCESS_CONTROL_LIMITED_IP),
				PublicFacing:  Ptr(PUBLIC_FACING_INTERNAL),
			},
			want: &AssessmentDetail{
				Result: Ptr(SYSTEM_EXPOSURE_SMALL),
				Score:  Ptr(float32(-0.1)),
			},
		},
		{
			name: "default",
			input: &SystemExposure{
				AccessControl: Ptr("HOGE"),
				PublicFacing:  Ptr("FUGA"),
			},
			want: &AssessmentDetail{
				Result: Ptr(TRIAGE_UNKNOWN),
				Score:  Ptr(float32(0)),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := evaluateSystemExposure(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestEvaluateUtility(t *testing.T) {
	tests := []struct {
		name  string
		input *Utility
		want  *AssessmentDetail
	}{
		{
			name:  "nil",
			input: &Utility{},
			want: &AssessmentDetail{
				Result: Ptr(TRIAGE_UNKNOWN),
				Score:  Ptr(float32(0)),
			},
		},
		{
			name: "missing_automatable",
			input: &Utility{
				ValueDensity: Ptr(VALUE_DENSITY_CONCENTRATED),
			},
			want: &AssessmentDetail{
				Result: Ptr(TRIAGE_UNKNOWN),
				Score:  Ptr(float32(0)),
			},
		},
		{
			name: "missing_valdensity",
			input: &Utility{
				Automatable: Ptr(AUTOMATABLE_YES),
			},
			want: &AssessmentDetail{
				Result: Ptr(TRIAGE_UNKNOWN),
				Score:  Ptr(float32(0)),
			},
		},
		{
			name: "super_efficient",
			input: &Utility{
				Automatable:  Ptr(AUTOMATABLE_YES),
				ValueDensity: Ptr(VALUE_DENSITY_CONCENTRATED),
			},
			want: &AssessmentDetail{
				Result: Ptr(UTILITY_SUPER_EFFICIENT),
				Score:  Ptr(float32(0)),
			},
		},
		{
			name: "efficient_from_automatable",
			input: &Utility{
				Automatable:  Ptr(AUTOMATABLE_YES),
				ValueDensity: Ptr(VALUE_DENSITY_DIFFUSE),
			},
			want: &AssessmentDetail{
				Result: Ptr(UTILITY_EFFICIENT),
				Score:  Ptr(float32(-0.1)),
			},
		},
		{
			name: "efficient_from_valdensity",
			input: &Utility{
				Automatable:  Ptr(AUTOMATABLE_NO),
				ValueDensity: Ptr(VALUE_DENSITY_CONCENTRATED),
			},
			want: &AssessmentDetail{
				Result: Ptr(UTILITY_EFFICIENT),
				Score:  Ptr(float32(-0.1)),
			},
		},
		{
			name: "laborious",
			input: &Utility{
				Automatable:  Ptr(AUTOMATABLE_NO),
				ValueDensity: Ptr(VALUE_DENSITY_DIFFUSE),
			},
			want: &AssessmentDetail{
				Result: Ptr(UTILITY_LABORIOUS),
				Score:  Ptr(float32(-0.1)),
			},
		},
		{
			name: "default",
			input: &Utility{
				Automatable:  Ptr("HOGE"),
				ValueDensity: Ptr("FUGA"),
			},
			want: &AssessmentDetail{
				Result: Ptr(TRIAGE_UNKNOWN),
				Score:  Ptr(float32(0)),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := evaluateUtility(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestEvaluateHumanImpact(t *testing.T) {
	tests := []struct {
		name  string
		input *HumanImpact
		want  *AssessmentDetail
	}{
		{
			name:  "nil",
			input: &HumanImpact{},
			want: &AssessmentDetail{
				Result: Ptr(TRIAGE_UNKNOWN),
				Score:  Ptr(float32(0)),
			},
		},
		{
			name: "very_high_from_safety",
			input: &HumanImpact{
				SafetyImpact:  Ptr(SAFETY_IMPACT_CATASTROPHIC),
				MissionImpact: Ptr(TRIAGE_UNKNOWN),
			},
			want: &AssessmentDetail{
				Result: Ptr(HUMAN_IMPACT_VERY_HIGH),
				Score:  Ptr(float32(0)),
			},
		},
		{
			name: "very_high_from_mission",
			input: &HumanImpact{
				SafetyImpact:  Ptr(TRIAGE_UNKNOWN),
				MissionImpact: Ptr(MISSION_IMPACT_MISSION_FAILURE),
			},
			want: &AssessmentDetail{
				Result: Ptr(HUMAN_IMPACT_VERY_HIGH),
				Score:  Ptr(float32(0)),
			},
		},
		{
			name: "high_from_critical",
			input: &HumanImpact{
				SafetyImpact:  Ptr(SAFETY_IMPACT_CRITICAL),
				MissionImpact: Ptr(MISSION_IMPACT_NONE),
			},
			want: &AssessmentDetail{
				Result: Ptr(HUMAN_IMPACT_HIGH),
				Score:  Ptr(float32(-0.1)),
			},
		},
		{
			name: "high_from_marginal_mef",
			input: &HumanImpact{
				SafetyImpact:  Ptr(SAFETY_IMPACT_MARGINAL),
				MissionImpact: Ptr(MISSION_IMPACT_MEF_FAILURE),
			},
			want: &AssessmentDetail{
				Result: Ptr(HUMAN_IMPACT_HIGH),
				Score:  Ptr(float32(-0.1)),
			},
		},
		{
			name: "medium_from_negligible",
			input: &HumanImpact{
				SafetyImpact:  Ptr(SAFETY_IMPACT_NEGLIGIBLE),
				MissionImpact: Ptr(MISSION_IMPACT_MEF_FAILURE),
			},
			want: &AssessmentDetail{
				Result: Ptr(HUMAN_IMPACT_MEDIUM),
				Score:  Ptr(float32(-0.1)),
			},
		},
		{
			name: "medium_from_marginal",
			input: &HumanImpact{
				SafetyImpact:  Ptr(SAFETY_IMPACT_MARGINAL),
				MissionImpact: Ptr(MISSION_IMPACT_NONE),
			},
			want: &AssessmentDetail{
				Result: Ptr(HUMAN_IMPACT_MEDIUM),
				Score:  Ptr(float32(-0.1)),
			},
		},
		{
			name: "low",
			input: &HumanImpact{
				SafetyImpact:  Ptr(SAFETY_IMPACT_NEGLIGIBLE),
				MissionImpact: Ptr(MISSION_IMPACT_DEGRADED),
			},
			want: &AssessmentDetail{
				Result: Ptr(HUMAN_IMPACT_LOW),
				Score:  Ptr(float32(0)),
			},
		},
		{
			name: "default",
			input: &HumanImpact{
				SafetyImpact:  Ptr("HOGE"),
				MissionImpact: Ptr("FUGA"),
			},
			want: &AssessmentDetail{
				Result: Ptr(TRIAGE_UNKNOWN),
				Score:  Ptr(float32(0)),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := evaluateHumanImpact(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAdjustScore(t *testing.T) {
	type args struct {
		assessment *TriageAssessment
		finding    *model.Finding
	}
	tests := []struct {
		name                  string
		input                 args
		expectedBaseScore     float32
		expectedAdjustedScore float32
	}{
		{
			name: "all_nil",
			input: args{
				assessment: &TriageAssessment{},
				finding:    &model.Finding{Score: 1.0},
			},
			expectedBaseScore:     1.0,
			expectedAdjustedScore: 1.0,
		},
		{
			name: "only_exploitation",
			input: args{
				assessment: &TriageAssessment{
					Exploitation: &AssessmentDetail{Score: Ptr(float32(-0.1))},
				},
				finding: &model.Finding{Score: 1.0},
			},
			expectedBaseScore:     1.0,
			expectedAdjustedScore: 0.9,
		},
		{
			name: "exploitation_and_utility",
			input: args{
				assessment: &TriageAssessment{
					Exploitation: &AssessmentDetail{Score: Ptr(float32(-0.1))},
					Utility:      &AssessmentDetail{Score: Ptr(float32(-0.1))},
				},
				finding: &model.Finding{Score: 1.0},
			},
			expectedBaseScore:     1.0,
			expectedAdjustedScore: 0.8,
		},
		{
			name: "all_assessments",
			input: args{
				assessment: &TriageAssessment{
					Exploitation:   &AssessmentDetail{Score: Ptr(float32(0))},
					SystemExposure: &AssessmentDetail{Score: Ptr(float32(-0.1))},
					Utility:        &AssessmentDetail{Score: Ptr(float32(-0.1))},
					HumanImpact:    &AssessmentDetail{Score: Ptr(float32(-0.1))},
				},
				finding: &model.Finding{Score: 1.0},
			},
			expectedBaseScore:     1.0,
			expectedAdjustedScore: 0.7,
		},
		{
			name: "only_utility",
			input: args{
				assessment: &TriageAssessment{
					Utility: &AssessmentDetail{Score: Ptr(float32(-0.2))},
				},
				finding: &model.Finding{Score: 1.0},
			},
			expectedBaseScore:     1.0,
			expectedAdjustedScore: 0.8,
		},
		{
			name: "min_score",
			input: args{
				assessment: &TriageAssessment{
					Exploitation: &AssessmentDetail{Score: Ptr(float32(-1.1))},
				},
				finding: &model.Finding{Score: 1.0},
			},
			expectedBaseScore:     1.0,
			expectedAdjustedScore: 0.0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create RiskenTriage
			triage := &RiskenTriage{
				Assessment: tc.input.assessment,
			}
			result := adjustScore(triage, tc.input.finding)

			if result.BaseScore == nil {
				t.Errorf("[%s] base score is nil; expected %v", tc.name, tc.expectedBaseScore)
			} else if diff := cmp.Diff(*result.BaseScore, tc.expectedBaseScore); diff != "" {
				t.Errorf("[%s] base score mismatch (-want +got):\n%s", tc.name, diff)
			}
			if result.AdjustedScore == nil {
				t.Errorf("[%s] adjusted score is nil; expected %v", tc.name, tc.expectedAdjustedScore)
			} else if diff := cmp.Diff(*result.AdjustedScore, tc.expectedAdjustedScore); diff != "" {
				t.Errorf("[%s] adjusted score mismatch (-want +got):\n%s", tc.name, diff)
			}
		})
	}
}

func TestUpdateFindingData(t *testing.T) {
	type args struct {
		finding *model.Finding
		triaged *RiskenTriage
	}
	tests := []struct {
		name         string
		input        args
		expectedData map[string]interface{}
	}{
		{
			name: "missing_risken_triage",
			input: args{
				finding: &model.Finding{Data: `{"other_field": "value"}`},
				triaged: &RiskenTriage{
					BaseScore:     Ptr(float32(1.0)),
					AdjustedScore: Ptr(float32(0.5)),
				},
			},
			expectedData: map[string]interface{}{
				"other_field": "value",
				"risken_triage": map[string]interface{}{
					"base_score":     1.0,
					"adjusted_score": 0.5,
				},
			},
		},
		{
			name: "existing_risken_triage",
			input: args{
				finding: &model.Finding{Data: `{"other_field": "value", "risken_triage": {"dummy": "old"}}`},
				triaged: &RiskenTriage{
					BaseScore:     Ptr(float32(2.0)),
					AdjustedScore: Ptr(float32(1.5)),
				},
			},
			expectedData: map[string]interface{}{
				"other_field": "value",
				"risken_triage": map[string]interface{}{
					"base_score":     2.0,
					"adjusted_score": 1.5,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			updatedFinding, err := updateFindingData(tc.input.finding, tc.input.triaged)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			var updatedMap map[string]interface{}
			if err := json.Unmarshal([]byte(updatedFinding.Data), &updatedMap); err != nil {
				t.Fatalf("failed to unmarshal updated data: %v", err)
			}
			if diff := cmp.Diff(tc.expectedData, updatedMap); diff != "" {
				t.Errorf("updateFindingData mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTriageFinding(t *testing.T) {
	cases := []struct {
		name    string
		input   *model.Finding
		want    *model.Finding
		wantErr bool
	}{
		{
			name: "success to triage finding",
			input: &model.Finding{
				Score: 0.8,
				Data: `{
					"data": {
						"key1": "value1",
						"key2": "value2"
					},
					"risken_triage": {
						"source": {
							"utility": {
								"automatable": "NO",
								"value_density": "UNKNOWN"
							},
							"exploitation": {
								"has_cve": true,
								"has_kev": false,
								"epss_score": 0.00116,
								"public_poc": false
							}
						}
					}
				}`,
			},
			want: &model.Finding{
				Score: 0.6,
				Data: `{
					"data": {
						"key1": "value1",
						"key2": "value2"
					},
					"risken_triage": {
						"base_score": 0.8,
						"adjusted_score": 0.6,
						"source": {
							"utility": {
								"automatable": "NO",
								"value_density": "UNKNOWN"
							},
							"exploitation": {
								"has_cve": true,
								"has_kev": false,
								"epss_score": 0.00116,
								"public_poc": false
							}
						},
						"assessment": {
							"exploitation": {
								"result": "NONE",
								"score": -0.1
							},
							"utility": {
								"result": "LABORIOUS",
								"score": -0.1
							}
						}
					}
				}`,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			svc := FindingService{}
			got, err := svc.TriageFinding(context.Background(), c.input)
			if err != nil {
				if !c.wantErr {
					t.Fatalf("unexpected error: %+v", err)
				}
				return
			}
			if c.wantErr {
				t.Fatal("expected error but got nil")
			}

			var gotData, wantData map[string]interface{}
			if err := json.Unmarshal([]byte(got.Data), &gotData); err != nil {
				t.Fatalf("failed to unmarshal got data: %+v", err)
			}
			if err := json.Unmarshal([]byte(c.want.Data), &wantData); err != nil {
				t.Fatalf("failed to unmarshal want data: %+v", err)
			}

			if diff := cmp.Diff(wantData, gotData); diff != "" {
				t.Errorf("Data mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
