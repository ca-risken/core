package db

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/ca-risken/core/proto/finding"
)

type assetRiskFinding struct {
	FindingID                               uint64
	ProjectID                               uint32
	AssetType, AssetKey, ResourceName, Data string
	UpdatedAt                               int64
}

// Risk evidence is independent of score/tag filters, but retains the caller's
// authorization and finding status scope. Only the displayed assets are loaded.
func (c *Client) loadAssetRisks(ctx context.Context, req *finding.ListFindingRequest, orgID uint32, assets []*finding.FindingAsset) error {
	if len(assets) == 0 {
		return nil
	}
	var cond string
	var args []any
	if orgID == 0 {
		cond, args = generateProjectFindingQuery(req.ProjectId, 0, 0, 1, 0, nil, nil, nil, req.Status)
	} else {
		cond, args = generateOrganizationFindingQuery(orgID, 0, 0, 1, 0, nil, nil, nil, req.Status)
	}
	matches := []string{}
	for _, a := range assets {
		if a.AssetType == "Unknown" {
			matches = append(matches, "(finding.project_id = ? and finding.finding_id = ?)")
			args = append(args, a.ProjectId, a.AssetKey)
		} else {
			match := "(finding.project_id = ? and ((finding.asset_type = ? and finding.asset_key = ?)"
			args = append(args, a.ProjectId, a.AssetType, a.AssetKey)
			// A VM/Resource name is a full resource ID. Do not transfer a cluster's or
			// image repository's exposure to a workload or an image build.
			if a.AssetType == "VM" || a.AssetType == "Resource" {
				match += " or finding.resource_name = ?"
				args = append(args, a.AssetName)
			}
			matches = append(matches, match+"))")
		}
	}
	query := "SELECT DISTINCT finding.finding_id, finding.project_id, finding.asset_type, finding.asset_key, finding.resource_name, finding.data, UNIX_TIMESTAMP(finding.updated_at) AS updated_at FROM finding" + cond + " and finding.data_source = 'google:scc' and (" + strings.Join(matches, " or ") + ") ORDER BY finding.finding_id"
	var rows []assetRiskFinding
	if err := c.Slave.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return err
	}
	for _, a := range assets {
		for _, row := range rows {
			sameAsset := row.ProjectID == a.ProjectId && row.AssetType == a.AssetType && row.AssetKey == a.AssetKey
			if a.AssetType == "Unknown" {
				sameAsset = row.ProjectID == a.ProjectId && a.AssetKey == strconv.FormatUint(row.FindingID, 10)
			}
			sameResource := row.ProjectID == a.ProjectId && (a.AssetType == "VM" || a.AssetType == "Resource") && row.ResourceName == a.AssetName
			if !sameAsset && !sameResource {
				continue
			}
			for _, risk := range sccAssetRisks(row) {
				if sameAsset || (sameResource && risk.Kind == "internet") {
					a.RiskInformation = append(a.RiskInformation, risk)
				}
			}
		}
	}
	return nil
}

func sccAssetRisks(row assetRiskFinding) []*finding.FindingAssetRisk {
	var data struct {
		URL           string `json:"scc_detail_url"`
		Vulnerability struct {
			KEV *struct {
				CVEID string `json:"cveID"`
			} `json:"kev"`
		} `json:"vulnerability"`
		Finding struct {
			Category      string          `json:"category"`
			State         json.RawMessage `json:"state"`
			Vulnerability struct {
				CVE struct {
					ID       string          `json:"id"`
					Activity json.RawMessage `json:"exploitation_activity"`
					Release  json.RawMessage `json:"exploit_release_date"`
				} `json:"cve"`
			} `json:"vulnerability"`
		} `json:"finding"`
	}
	if json.Unmarshal([]byte(row.Data), &data) != nil {
		return nil
	}
	f := data.Finding
	// SCC protobuf JSON is persisted using encoding/json (numeric enums).
	if string(f.State) == "2" || string(f.State) == `"INACTIVE"` {
		return nil
	}
	risks := []*finding.FindingAssetRisk{}
	makeRisk := func(kind, detail string) *finding.FindingAssetRisk {
		return &finding.FindingAssetRisk{Kind: kind, FindingId: row.FindingID, Category: f.Category, CveId: f.Vulnerability.CVE.ID, Detail: detail, UpdatedAt: row.UpdatedAt, SccDetailUrl: data.URL}
	}
	// These categories explicitly report public access settings. External IPs,
	// firewall rules and attack exposure scores alone are insufficient evidence.
	switch f.Category {
	case "PUBLIC_BUCKET_ACL", "PUBLIC_LOG_BUCKET", "PUBLIC_DATASET", "PUBLIC_SQL_INSTANCE":
		risks = append(risks, makeRisk("internet", f.Category))
	}
	cve := f.Vulnerability.CVE
	if cve.ID != "" {
		if data.Vulnerability.KEV != nil && data.Vulnerability.KEV.CVEID == cve.ID {
			risks = append(risks, makeRisk("exploit", "KEV"))
		}
		if string(cve.Activity) == "3" || string(cve.Activity) == `"AVAILABLE"` {
			risks = append(risks, makeRisk("exploit", "AVAILABLE"))
		} else if hasExploitRelease(cve.Release) {
			risks = append(risks, makeRisk("exploit", "exploit_release_date"))
		}
	}
	return risks
}

func hasExploitRelease(raw json.RawMessage) bool {
	var date string
	if json.Unmarshal(raw, &date) == nil {
		_, err := time.Parse(time.RFC3339Nano, date)
		return err == nil
	}
	var timestamp struct {
		Seconds int64 `json:"seconds"`
	}
	return json.Unmarshal(raw, &timestamp) == nil && timestamp.Seconds > 0
}
