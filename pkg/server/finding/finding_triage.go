package finding

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ca-risken/core/pkg/model"
)

type RiskenTriage struct {
	BaseScore     *float32          `json:"base_score,omitempty"`
	AdjustedScore *float32          `json:"adjusted_score,omitempty"`
	Source        *TriageSource     `json:"source,omitempty"`
	Assessment    *TriageAssessment `json:"assessment,omitempty"`
}

type TriageSource struct {
	Exploitation   *Exploitation   `json:"exploitation,omitempty"`
	SystemExposure *SystemExposure `json:"system_exposure,omitempty"`
	Utility        *Utility        `json:"utility,omitempty"`
	HumanImpact    *HumanImpact    `json:"human_impact,omitempty"`
}

type Exploitation struct {
	HasCVE    *bool    `json:"has_cve,omitempty"`
	EpssScore *float32 `json:"epss_score,omitempty"`
	HasKEV    *bool    `json:"has_kev,omitempty"`
	PublicPOC *bool    `json:"public_poc,omitempty"`
}

type SystemExposure struct {
	PublicFacing  *string `json:"public_facing,omitempty"`
	AccessControl *string `json:"access_control,omitempty"`
}

type Utility struct {
	Automatable  *string `json:"automatable,omitempty"`
	ValueDensity *string `json:"value_density,omitempty"`
}

type HumanImpact struct {
	SafetyImpact  *string `json:"safety_impact,omitempty"`
	MissionImpact *string `json:"mission_impact,omitempty"`
}

type TriageAssessment struct {
	Exploitation   *AssessmentDetail `json:"exploitation,omitempty"`
	SystemExposure *AssessmentDetail `json:"system_exposure,omitempty"`
	Utility        *AssessmentDetail `json:"utility,omitempty"`
	HumanImpact    *AssessmentDetail `json:"human_impact,omitempty"`
}

type AssessmentDetail struct {
	Result *string  `json:"result,omitempty"`
	Score  *float32 `json:"score,omitempty"`
}

func (f *FindingService) TriageFinding(ctx context.Context, finding *model.Finding) (*model.Finding, error) {
	findingData := finding.Data
	riskenTriage := RiskenTriage{}
	if err := json.Unmarshal([]byte(findingData), &riskenTriage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal finding data: %w", err)
	}
	if riskenTriage.Source == nil {
		return finding, nil // no triage data
	}

	source := riskenTriage.Source
	assessment := TriageAssessment{}

	// 1. Exploitation
	if source.Exploitation != nil {
		assessment.Exploitation = evaluateExploitation(source.Exploitation)
	}

	// 2. SystemExposure
	if source.SystemExposure != nil {
		assessment.SystemExposure = evaluateSystemExposure(source.SystemExposure)
	}

	// 3. Utility
	if source.Utility != nil {
		assessment.Utility = evaluateUtility(source.Utility)
	}

	// 4. HumanImpact
	if source.HumanImpact != nil {
		assessment.HumanImpact = evaluateHumanImpact(source.HumanImpact)
	}

	// Adjust score
	triaged := adjustScore(&riskenTriage, &assessment, finding)
	updatedFinding, err := updateFindingData(finding, triaged)
	if err != nil {
		return nil, fmt.Errorf("failed to update finding data: %w", err)
	}
	return updatedFinding, nil
}

// evaluateExploitation return ACTIVE, PUBLIC_POC, UNKNOWN
// @ref https://certcc.github.io/SSVC/reference/decision_points/exploitation/
func evaluateExploitation(source *Exploitation) *AssessmentDetail {
	assessment := AssessmentDetail{
		Result: Ptr("UNKNOWN"),
		Score:  Ptr(float32(0)),
	}
	if source.HasCVE == nil || !*source.HasCVE {
		// no CVE
		return &assessment
	}
	if source.HasKEV != nil && *source.HasKEV {
		// has KEV
		assessment.Result = Ptr("ACTIVE")
		return &assessment
	}
	if source.PublicPOC != nil && *source.PublicPOC {
		// has public POC
		assessment.Result = Ptr("PUBLIC_POC")
		assessment.Score = Ptr(float32(-0.1))
		if source.EpssScore != nil && *source.EpssScore >= 0.01 {
			// epss >= 1%
			assessment.Result = Ptr("ACTIVE")
			assessment.Score = Ptr(float32(0.0))
			return &assessment
		}
		return &assessment
	}
	assessment.Result = Ptr("NONE")
	assessment.Score = Ptr(float32(-0.1))
	return &assessment
}

// evaluateSystemExposure return OPEN, CONTROLLED, SMALL, UNKNOWN
// @ref https://certcc.github.io/SSVC/reference/decision_points/system_exposure/
func evaluateSystemExposure(source *SystemExposure) *AssessmentDetail {
	assessment := AssessmentDetail{
		Result: Ptr("UNKNOWN"),
		Score:  Ptr(float32(0)),
	}

	// UNKNOWN
	if source.PublicFacing == nil || source.AccessControl == nil {
		return &assessment
	}
	if *source.PublicFacing == "UNKNOWN" || *source.AccessControl == "UNKNOWN" {
		return &assessment
	}

	// OPEN
	if *source.PublicFacing == "OPEN" && *source.AccessControl == "NONE" {
		assessment.Result = Ptr("OPEN")
		return &assessment
	}

	// CONTROLLED
	if *source.PublicFacing == "OPEN" && *source.AccessControl == "LIMITED_IP" {
		assessment.Result = Ptr("CONTROLLED")
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}
	if *source.PublicFacing == "OPEN" && *source.AccessControl == "AUTHENTICATED" {
		assessment.Result = Ptr("CONTROLLED")
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}
	if *source.PublicFacing == "INTERNAL" && *source.AccessControl == "NONE" {
		assessment.Result = Ptr("CONTROLLED")
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}

	// SMALL
	if *source.PublicFacing == "INTERNAL" && *source.AccessControl == "LIMITED_IP" {
		assessment.Result = Ptr("SMALL")
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}
	if *source.PublicFacing == "INTERNAL" && *source.AccessControl == "AUTHENTICATED" {
		assessment.Result = Ptr("SMALL")
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}

	// Other
	return &assessment
}

// evaluateUtility return SUPER_EFFICIENT, EFFICIENT, LABORIOUS, UNKNOWN
// @ref https://certcc.github.io/SSVC/reference/decision_points/utility/
func evaluateUtility(source *Utility) *AssessmentDetail {
	assessment := AssessmentDetail{
		Result: Ptr("UNKNOWN"),
		Score:  Ptr(float32(0)),
	}
	// UNKNOWN
	if source.Automatable == nil || source.ValueDensity == nil {
		return &assessment
	}

	// SUPER_EFFICIENT
	if *source.Automatable == "YES" && *source.ValueDensity == "CONCENTRATED" {
		assessment.Result = Ptr("SUPER_EFFICIENT")
		return &assessment
	}

	// EFFICIENT
	if *source.Automatable == "YES" {
		assessment.Result = Ptr("EFFICIENT")
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}
	if *source.ValueDensity == "CONCENTRATED" {
		assessment.Result = Ptr("EFFICIENT")
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}

	// LABORIOUS
	if *source.Automatable == "NO" && *source.ValueDensity == "DIFFUSE" {
		assessment.Result = Ptr("LABORIOUS")
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}

	// UNKNOWN
	return &assessment
}

// evaluateHumanImpact return VERY_HIGH, HIGH, MEDIUM, LOW, UNKNOWN
// @ref https://certcc.github.io/SSVC/reference/decision_points/human_impact/
func evaluateHumanImpact(source *HumanImpact) *AssessmentDetail {
	assessment := AssessmentDetail{
		Result: Ptr("UNKNOWN"),
		Score:  Ptr(float32(0)),
	}

	// UNKNOWN
	if source.SafetyImpact == nil || source.MissionImpact == nil {
		return &assessment
	}

	// VERY_HIGH
	if *source.SafetyImpact == "CATASTROPHIC" || *source.MissionImpact == "MISSION_FAILURE" {
		assessment.Result = Ptr("VERY_HIGH")
		return &assessment
	}

	// HIGH
	if *source.SafetyImpact == "CRITICAL" && (*source.MissionImpact == "NONE" || *source.MissionImpact == "DEGRADED" || *source.MissionImpact == "CRIPPLED") {
		assessment.Result = Ptr("HIGH")
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}
	if *source.SafetyImpact == "MARGINAL" && *source.MissionImpact == "MEF_FAILURE" {
		assessment.Result = Ptr("HIGH")
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}

	// MEDIUM
	if *source.SafetyImpact == "NEGLIGIBLE" && *source.MissionImpact == "MEF_FAILURE" {
		assessment.Result = Ptr("MEDIUM")
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}
	if *source.SafetyImpact == "MARGINAL" && (*source.MissionImpact == "NONE" || *source.MissionImpact == "DEGRADED" || *source.MissionImpact == "CRIPPLED") {
		assessment.Result = Ptr("MEDIUM")
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}

	// LOW
	if *source.SafetyImpact == "NEGLIGIBLE" && (*source.MissionImpact == "NONE" || *source.MissionImpact == "DEGRADED" || *source.MissionImpact == "CRIPPLED") {
		assessment.Result = Ptr("LOW")
		return &assessment
	}

	// OTHER
	return &assessment
}

func adjustScore(triage *RiskenTriage, assessment *TriageAssessment, finding *model.Finding) *RiskenTriage {
	baseScore := finding.Score
	adjustment := float32(0.0)
	if assessment.Exploitation != nil {
		adjustment += *assessment.Exploitation.Score
	}
	if assessment.SystemExposure != nil {
		adjustment += *assessment.SystemExposure.Score
	}
	if assessment.Utility != nil {
		adjustment += *assessment.Utility.Score
	}
	if assessment.HumanImpact != nil {
		adjustment += *assessment.HumanImpact.Score
	}

	// overwrite
	triage.BaseScore = Ptr(baseScore)
	triage.AdjustedScore = Ptr(baseScore + adjustment)
	if *triage.AdjustedScore < 0 {
		*triage.AdjustedScore = 0
	}
	return triage
}

func updateFindingData(finding *model.Finding, triaged *RiskenTriage) (*model.Finding, error) {
	// Keep original data
	var originalData map[string]interface{}
	if err := json.Unmarshal([]byte(finding.Data), &originalData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal original finding.Data: %w", err)
	}

	// Convert recalculated triaged to JSON and then to map
	triagedJSON, err := json.Marshal(triaged)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal recalculated risken triage: %w", err)
	}
	var triagedData map[string]interface{}
	if err := json.Unmarshal(triagedJSON, &triagedData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recalculated risken triage JSON: %w", err)
	}

	// Replace the risken_triarge part of the original data with the recalculated triaged data
	originalData["risken_triarge"] = triagedData

	// Encode the updated map to JSON
	updatedJSON, err := json.Marshal(originalData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated data: %w", err)
	}
	finding.Data = string(updatedJSON)
	return finding, nil
}

func Ptr[T any](v T) *T {
	return &v
}
