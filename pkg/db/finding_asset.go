package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ca-risken/core/proto/finding"
)

const assetKeySQL = "CASE WHEN finding.asset_type = 'Unknown' THEN CAST(finding.finding_id AS CHAR) ELSE finding.asset_key END"

// assetConditions is appended after the existing permission/search/status scope.
func assetConditions(req *finding.ListFindingRequest) (string, []interface{}) {
	cond := " and finding.data_source = 'google:scc' and finding.asset_type <> ''"
	var args []interface{}
	if req.AssetType != "" {
		cond += " and finding.asset_type = ?"
		args = append(args, req.AssetType)
	}
	if req.AssetKey != "" {
		if req.AssetType == "Unknown" {
			id, _ := strconv.ParseUint(req.AssetKey, 10, 64)
			cond += " and finding.finding_id = ?"
			args = append(args, id)
		} else {
			cond += " and finding.asset_key = ?"
			args = append(args, req.AssetKey)
		}
	}
	if req.AssetProjectId != 0 {
		cond += " and finding.project_id = ?"
		args = append(args, req.AssetProjectId)
	}
	return cond, args
}

// ListFindingAssets also returns the exact expanded finding count when asset_view
// is false. DISTINCT removes duplicate tag matches before grouping/pagination.
func (c *Client) ListFindingAssets(ctx context.Context, req *finding.ListFindingRequest, organizationID uint32) ([]*finding.FindingAsset, int64, error) {
	var cond string
	var args []interface{}
	if organizationID == 0 {
		cond, args = generateProjectFindingQuery(req.ProjectId, req.AlertId, req.FromScore, req.ToScore, req.FindingId, req.DataSource, req.ResourceName, req.Tag, req.Status)
	} else {
		cond, args = generateOrganizationFindingQuery(organizationID, 0, req.FromScore, req.ToScore, req.FindingId, req.DataSource, req.ResourceName, req.Tag, req.Status)
	}
	extra, params := assetConditions(req)
	cond += extra
	args = append(args, params...)
	base := "SELECT DISTINCT finding.finding_id, finding.project_id, finding.asset_type, " + assetKeySQL + " AS asset_key, finding.asset_name, finding.score, finding.updated_at FROM finding" + cond
	var total int64
	if !req.AssetView {
		err := c.Slave.WithContext(ctx).Raw("SELECT COUNT(*) FROM ("+base+") scoped", args...).Scan(&total).Error
		return nil, total, err
	}
	grouped := "SELECT project_id, asset_type, asset_key, MIN(asset_name) AS asset_name, MAX(score) AS max_score, COUNT(*) AS finding_count, UNIX_TIMESTAMP(MAX(updated_at)) AS updated_at FROM (" + base + ") scoped GROUP BY project_id, asset_type, asset_key"
	if err := c.Slave.WithContext(ctx).Raw("SELECT COUNT(*) FROM ("+grouped+") assets", args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}
	sort := req.Sort
	switch sort {
	case "score":
		sort = "max_score"
	case "updated_at", "asset_name", "asset_type", "finding_count", "project_id":
	default:
		sort = "max_score"
	}
	direction := "desc"
	if req.Direction == "asc" {
		direction = "asc"
	}
	query := grouped + fmt.Sprintf(" ORDER BY %s %s, project_id ASC, asset_type ASC, asset_key ASC LIMIT %d OFFSET %d", sort, direction, req.Limit, req.Offset)
	var assets []*finding.FindingAsset
	err := c.Slave.WithContext(ctx).Raw(query, args...).Scan(&assets).Error
	return assets, total, err
}
