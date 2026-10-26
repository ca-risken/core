package db

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/ca-risken/core/proto/finding"
)

type assetRiskFinding struct {
	FindingID  uint64
	Data       string
	UpdatedAt  int64
	AssetIndex int
	SameAsset  bool
}

// Read compact evidence incrementally; do not load complete Finding JSON into memory.
func (c *Client) loadAssetRisks(ctx context.Context, req *finding.ListFindingRequest, orgID uint32, assets []*finding.FindingAsset) error {
	if len(assets) == 0 {
		return nil
	}
	var cond string
	var scopeArgs []any
	if orgID == 0 {
		cond, scopeArgs = generateProjectFindingQuery(req.ProjectId, 0, req.FromScore, req.ToScore, 0, nil, nil, req.Tag, req.Status)
	} else {
		cond, scopeArgs = generateOrganizationFindingQuery(orgID, 0, req.FromScore, req.ToScore, 0, nil, nil, req.Tag, req.Status)
	}
	pages := []string{}
	args := []any{}
	for i, a := range assets {
		pages = append(pages, "SELECT ? AS asset_index, ? AS project_id, ? AS asset_type, ? AS asset_key")
		args = append(args, i, a.ProjectId, a.AssetType, a.AssetKey)
	}
	args = append(args, scopeArgs...)
	sameAsset := "((page.asset_type = 'Unknown' AND CAST(finding.finding_id AS CHAR) = page.asset_key) OR (page.asset_type <> 'Unknown' AND finding.asset_type = page.asset_type AND finding.asset_key = page.asset_key))"
	// Resolve the full resource name from the original Finding, never its display name
	// or hashed asset_key. The original Finding is used only for identity matching.
	relatedResource := "(page.asset_type IN ('VM','Resource') AND EXISTS (SELECT 1 FROM finding reference_finding WHERE reference_finding.project_id = page.project_id AND reference_finding.asset_type = page.asset_type AND reference_finding.asset_key = page.asset_key AND reference_finding.resource_name = finding.resource_name))"
	compactData := `JSON_OBJECT('finding', JSON_OBJECT(
 'category', JSON_EXTRACT(finding.data, '$.finding.category'),
 'state', JSON_EXTRACT(finding.data, '$.finding.state'),
 'vulnerability', JSON_OBJECT('cve', JSON_OBJECT(
 'id', JSON_EXTRACT(finding.data, '$.finding.vulnerability.cve.id'),
 'exploitation_activity', JSON_EXTRACT(finding.data, '$.finding.vulnerability.cve.exploitation_activity'),
 'exploit_release_date', JSON_EXTRACT(finding.data, '$.finding.vulnerability.cve.exploit_release_date')))),
 'vulnerability', JSON_OBJECT('kev', JSON_OBJECT('cveID', JSON_EXTRACT(finding.data, '$.vulnerability.kev.cveID'))))`
	query := "SELECT DISTINCT finding.finding_id, page.asset_index, " + sameAsset + " AS same_asset, " + compactData + " AS data FROM finding JOIN (" + strings.Join(pages, " UNION ALL ") + ") page ON finding.project_id = page.project_id AND (" + sameAsset + " OR " + relatedResource + ")" + cond + " and finding.data_source = 'google:scc' ORDER BY finding.finding_id, page.asset_index"
	rows, err := c.Slave.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var row assetRiskFinding
		if err := c.Slave.ScanRows(rows, &row); err != nil {
			return err
		}
		// One Finding ID per kind. KEV precedes public Exploit in sccAssetRisks,
		// so a combined detection retains the stronger icon without duplicate IDs.
		seen := map[string]bool{}
		for _, risk := range sccAssetRisks(row) {
			if !row.SameAsset && risk.Kind != "internet" {
				continue
			}
			if seen[risk.Kind] {
				continue
			}
			seen[risk.Kind] = true
			assets[row.AssetIndex].RiskInformation = append(assets[row.AssetIndex].RiskInformation, &finding.FindingAssetRisk{Kind: risk.Kind, FindingId: risk.FindingId, Detail: risk.Detail})
		}
	}
	return rows.Err()
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
