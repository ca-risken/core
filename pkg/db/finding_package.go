package db

import (
	"context"
	"encoding/json"

	"github.com/ca-risken/core/proto/finding"
)

// Missing, stale, and cross-source recommendations stay visible as review.
const remediationMethodSQL = "COALESCE(NULLIF(asset_recommend.remediation_type, ''), 'review')"
const remediationJoinSQL = " LEFT JOIN recommend_finding asset_rf ON asset_rf.finding_id = finding.finding_id AND asset_rf.project_id = finding.project_id LEFT JOIN recommend asset_recommend ON asset_recommend.recommend_id = asset_rf.recommend_id AND asset_recommend.data_source = finding.data_source "

func (c *Client) ListFindingPackages(ctx context.Context, req *finding.ListFindingRequest, organizationID uint32) ([]*finding.FindingPackage, error) {
	var cond string
	var args []any
	if organizationID == 0 {
		cond, args = generateProjectFindingQuery(req.ProjectId, req.AlertId, req.FromScore, req.ToScore, req.FindingId, req.DataSource, req.ResourceName, req.Tag, req.Status)
	} else {
		cond, args = generateOrganizationFindingQuery(organizationID, 0, req.FromScore, req.ToScore, req.FindingId, req.DataSource, req.ResourceName, req.Tag, req.Status)
	}
	extra, params := assetConditions(req)
	cond += extra
	args = append(args, params...)
	// Only IDs, source, and the stored method are read; no Finding JSON or
	// recommendation text is transferred or parsed to count groups.
	scoped := "SELECT DISTINCT finding.finding_id, finding.project_id, finding.data_source FROM finding" + cond
	base := "SELECT finding.finding_id, finding.data_source, " + remediationMethodSQL + " AS method FROM (" + scoped + ") finding" + remediationJoinSQL
	var rows []struct {
		DataSource   string
		Method       string
		FindingCount uint64
	}
	query := "SELECT data_source, method, COUNT(*) AS finding_count FROM (" + base + ") scoped GROUP BY data_source, method ORDER BY finding_count DESC, data_source, method"
	if err := c.Slave.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	groups := make([]*finding.FindingPackage, 0, len(rows))
	for _, row := range rows {
		// Keep the existing four-string wire format; its meaning is now source +
		// remediation method, not package name, version, or filesystem location.
		key, _ := json.Marshal([4]string{row.DataSource, "", "", row.Method + ":"})
		groups = append(groups, &finding.FindingPackage{PackageKey: string(key), WorkType: row.Method, WorkLocation: row.DataSource, FindingCount: row.FindingCount})
	}
	return groups, nil
}
