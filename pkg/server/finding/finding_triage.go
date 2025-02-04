package finding

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ca-risken/core/pkg/model"
)

const (
	TRIAGE_UNKNOWN = "UNKNOWN"
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
	if strings.TrimSpace(finding.Data) == "" {
		return finding, nil // no triage data
	}
	riskenTriage := RiskenTriage{}
	if err := json.Unmarshal([]byte(finding.Data), &riskenTriage); err != nil {
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

const (
	// Exploitation
	EXPLOITATION_RESULT_ACTIVE     = "ACTIVE"
	EXPLOITATION_RESULT_PUBLIC_POC = "PUBLIC_POC"
	EXPLOITATION_RESULT_NONE       = "NONE"
)

// evaluateExploitation return ACTIVE, PUBLIC_POC, UNKNOWN
// @ref https://certcc.github.io/SSVC/reference/decision_points/exploitation/
func evaluateExploitation(source *Exploitation) *AssessmentDetail {
	assessment := AssessmentDetail{
		Result: Ptr(TRIAGE_UNKNOWN),
		Score:  Ptr(float32(0)),
	}
	if source.HasCVE == nil {
		// no CVE data
		return &assessment
	}
	if !*source.HasCVE {
		// no vulnerability
		assessment.Result = Ptr(EXPLOITATION_RESULT_NONE)
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}
	if source.HasKEV != nil && *source.HasKEV {
		// has KEV
		assessment.Result = Ptr(EXPLOITATION_RESULT_ACTIVE)
		return &assessment
	}
	if source.PublicPOC != nil && *source.PublicPOC {
		// has public POC
		assessment.Result = Ptr(EXPLOITATION_RESULT_PUBLIC_POC)
		assessment.Score = Ptr(float32(-0.1))
		if source.EpssScore != nil && *source.EpssScore >= 0.01 {
			// epss >= 1%
			assessment.Result = Ptr(EXPLOITATION_RESULT_ACTIVE)
			assessment.Score = Ptr(float32(0.0))
			return &assessment
		}
		return &assessment
	}
	assessment.Result = Ptr(EXPLOITATION_RESULT_NONE)
	assessment.Score = Ptr(float32(-0.1))
	return &assessment
}

const (
	// PublicFacing
	PUBLIC_FACING_OPEN     = "OPEN"
	PUBLIC_FACING_INTERNAL = "INTERNAL"

	// AccessControl
	ACCESS_CONTROL_NONE          = "NONE"
	ACCESS_CONTROL_LIMITED_IP    = "LIMITED_IP"
	ACCESS_CONTROL_AUTHENTICATED = "AUTHENTICATED"

	// SystemExposure
	SYSTEM_EXPOSURE_OPEN       = "OPEN"
	SYSTEM_EXPOSURE_CONTROLLED = "CONTROLLED"
	SYSTEM_EXPOSURE_SMALL      = "SMALL"
)

// evaluateSystemExposure return OPEN, CONTROLLED, SMALL, UNKNOWN
// @ref https://certcc.github.io/SSVC/reference/decision_points/system_exposure/
func evaluateSystemExposure(source *SystemExposure) *AssessmentDetail {
	assessment := AssessmentDetail{
		Result: Ptr(TRIAGE_UNKNOWN),
		Score:  Ptr(float32(0)),
	}

	// UNKNOWN
	if source.PublicFacing == nil || source.AccessControl == nil {
		return &assessment
	}
	if *source.PublicFacing == TRIAGE_UNKNOWN || *source.AccessControl == TRIAGE_UNKNOWN {
		return &assessment
	}

	// OPEN
	if *source.PublicFacing == PUBLIC_FACING_OPEN && *source.AccessControl == ACCESS_CONTROL_NONE {
		assessment.Result = Ptr(SYSTEM_EXPOSURE_OPEN)
		return &assessment
	}

	// CONTROLLED
	if *source.PublicFacing == PUBLIC_FACING_OPEN && *source.AccessControl == ACCESS_CONTROL_LIMITED_IP {
		assessment.Result = Ptr(SYSTEM_EXPOSURE_CONTROLLED)
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}
	if *source.PublicFacing == PUBLIC_FACING_OPEN && *source.AccessControl == ACCESS_CONTROL_AUTHENTICATED {
		assessment.Result = Ptr(SYSTEM_EXPOSURE_CONTROLLED)
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}
	if *source.PublicFacing == PUBLIC_FACING_INTERNAL && *source.AccessControl == ACCESS_CONTROL_NONE {
		assessment.Result = Ptr(SYSTEM_EXPOSURE_CONTROLLED)
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}

	// SMALL
	if *source.PublicFacing == PUBLIC_FACING_INTERNAL && *source.AccessControl == ACCESS_CONTROL_LIMITED_IP {
		assessment.Result = Ptr(SYSTEM_EXPOSURE_SMALL)
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}
	if *source.PublicFacing == PUBLIC_FACING_INTERNAL && *source.AccessControl == ACCESS_CONTROL_AUTHENTICATED {
		assessment.Result = Ptr(SYSTEM_EXPOSURE_SMALL)
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}

	// Other
	return &assessment
}

const (
	// Automatable
	AUTOMATABLE_YES = "YES"
	AUTOMATABLE_NO  = "NO"

	// ValueDensity
	VALUE_DENSITY_CONCENTRATED = "CONCENTRATED"
	VALUE_DENSITY_DIFFUSE      = "DIFFUSE"

	// Utility
	UTILITY_SUPER_EFFICIENT = "SUPER_EFFICIENT"
	UTILITY_EFFICIENT       = "EFFICIENT"
	UTILITY_LABORIOUS       = "LABORIOUS"
)

// evaluateUtility return SUPER_EFFICIENT, EFFICIENT, LABORIOUS, UNKNOWN
// @ref https://certcc.github.io/SSVC/reference/decision_points/utility/
func evaluateUtility(source *Utility) *AssessmentDetail {
	assessment := AssessmentDetail{
		Result: Ptr(TRIAGE_UNKNOWN),
		Score:  Ptr(float32(0)),
	}
	// UNKNOWN
	if source.Automatable == nil || source.ValueDensity == nil {
		return &assessment
	}

	// SUPER_EFFICIENT
	if *source.Automatable == AUTOMATABLE_YES && *source.ValueDensity == VALUE_DENSITY_CONCENTRATED {
		assessment.Result = Ptr(UTILITY_SUPER_EFFICIENT)
		return &assessment
	}

	// EFFICIENT
	if *source.Automatable == AUTOMATABLE_YES {
		assessment.Result = Ptr(UTILITY_EFFICIENT)
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}
	if *source.ValueDensity == VALUE_DENSITY_CONCENTRATED {
		assessment.Result = Ptr(UTILITY_EFFICIENT)
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}

	// LABORIOUS
	if *source.Automatable == AUTOMATABLE_NO && *source.ValueDensity == VALUE_DENSITY_DIFFUSE {
		assessment.Result = Ptr(UTILITY_LABORIOUS)
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}

	// UNKNOWN
	return &assessment
}

const (
	// SafetyImpact
	SAFETY_IMPACT_CATASTROPHIC = "CATASTROPHIC"
	SAFETY_IMPACT_CRITICAL     = "CRITICAL"
	SAFETY_IMPACT_MARGINAL     = "MARGINAL"
	SAFETY_IMPACT_NEGLIGIBLE   = "NEGLIGIBLE"

	// MissionImpact
	MISSION_IMPACT_MISSION_FAILURE = "MISSION_FAILURE"
	MISSION_IMPACT_MEF_FAILURE     = "MEF_FAILURE"
	MISSION_IMPACT_NONE            = "NONE"
	MISSION_IMPACT_DEGRADED        = "DEGRADED"
	MISSION_IMPACT_CRIPPLED        = "CRIPPLED"

	// HumanImpact
	HUMAN_IMPACT_VERY_HIGH = "VERY_HIGH"
	HUMAN_IMPACT_HIGH      = "HIGH"
	HUMAN_IMPACT_MEDIUM    = "MEDIUM"
	HUMAN_IMPACT_LOW       = "LOW"
)

// evaluateHumanImpact return VERY_HIGH, HIGH, MEDIUM, LOW, UNKNOWN
// @ref https://certcc.github.io/SSVC/reference/decision_points/human_impact/
func evaluateHumanImpact(source *HumanImpact) *AssessmentDetail {
	assessment := AssessmentDetail{
		Result: Ptr(TRIAGE_UNKNOWN),
		Score:  Ptr(float32(0)),
	}

	// UNKNOWN
	if source.SafetyImpact == nil || source.MissionImpact == nil {
		return &assessment
	}

	// VERY_HIGH
	if *source.SafetyImpact == SAFETY_IMPACT_CATASTROPHIC || *source.MissionImpact == MISSION_IMPACT_MISSION_FAILURE {
		assessment.Result = Ptr(HUMAN_IMPACT_VERY_HIGH)
		return &assessment
	}

	// HIGH
	if *source.SafetyImpact == SAFETY_IMPACT_CRITICAL &&
		(*source.MissionImpact == MISSION_IMPACT_NONE ||
			*source.MissionImpact == MISSION_IMPACT_DEGRADED ||
			*source.MissionImpact == MISSION_IMPACT_CRIPPLED) {
		assessment.Result = Ptr(HUMAN_IMPACT_HIGH)
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}
	if *source.SafetyImpact == SAFETY_IMPACT_MARGINAL && *source.MissionImpact == MISSION_IMPACT_MEF_FAILURE {
		assessment.Result = Ptr(HUMAN_IMPACT_HIGH)
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}

	// MEDIUM
	if *source.SafetyImpact == SAFETY_IMPACT_NEGLIGIBLE && *source.MissionImpact == MISSION_IMPACT_MEF_FAILURE {
		assessment.Result = Ptr(HUMAN_IMPACT_MEDIUM)
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}
	if *source.SafetyImpact == SAFETY_IMPACT_MARGINAL &&
		(*source.MissionImpact == MISSION_IMPACT_NONE ||
			*source.MissionImpact == MISSION_IMPACT_DEGRADED ||
			*source.MissionImpact == MISSION_IMPACT_CRIPPLED) {
		assessment.Result = Ptr(HUMAN_IMPACT_MEDIUM)
		assessment.Score = Ptr(float32(-0.1))
		return &assessment
	}

	// LOW
	if *source.SafetyImpact == SAFETY_IMPACT_NEGLIGIBLE &&
		(*source.MissionImpact == MISSION_IMPACT_NONE ||
			*source.MissionImpact == MISSION_IMPACT_DEGRADED ||
			*source.MissionImpact == MISSION_IMPACT_CRIPPLED) {
		assessment.Result = Ptr(HUMAN_IMPACT_LOW)
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
