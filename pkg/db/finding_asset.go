package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ca-risken/core/proto/finding"
)

const assetKeySQL = "CASE WHEN finding.asset_type = 'Unknown' THEN CAST(finding.finding_id AS CHAR) ELSE finding.asset_key END"

// assetConditions is appended after the existing permission/search/status scope.
func assetConditions(req *finding.ListFindingRequest) (string, []any) {
	cond := " and finding.data_source IN ('google:scc', 'google:asset', 'google:cloudsploit', 'aws:access-analyzer', 'aws:guard-duty', 'aws:cloudsploit', 'aws:admin-checker') and finding.asset_type IS NOT NULL"
	var args []any
	if req.AssetName != "" {
		cond += " and INSTR(finding.asset_name, ?) > 0"
		args = append(args, req.AssetName)
	}
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
	if req.PackageKey != "" {
		var parts []string
		if json.Unmarshal([]byte(req.PackageKey), &parts) != nil || len(parts) != 4 {
			return cond + " and 1=0", args
		}
		if parts[0] == "" || parts[1] != "" || parts[2] != "" || len(parts[3]) < 2 || parts[3][len(parts[3])-1] != ':' {
			return cond + " and 1=0", args
		}
		method := parts[3][:len(parts[3])-1]
		cond += " and finding.data_source = ? and COALESCE((SELECT NULLIF(r.remediation_type, '') FROM recommend_finding rf INNER JOIN recommend r ON r.recommend_id = rf.recommend_id AND r.data_source = finding.data_source WHERE rf.finding_id = finding.finding_id AND rf.project_id = finding.project_id), 'review') = ?"
		args = append(args, parts[0], method)
	}

	return cond, args
}

// ListFindingAssets also returns the exact expanded finding count when asset_view
// is false. DISTINCT removes duplicate tag matches before grouping/pagination.
func (c *Client) ListFindingAssets(ctx context.Context, req *finding.ListFindingRequest, organizationID uint32) ([]*finding.FindingAsset, int64, uint64, error) {
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
	base := "SELECT DISTINCT finding.finding_id, finding.project_id, finding.asset_type, " + assetKeySQL + " AS asset_key, finding.asset_name, finding.data_source, finding.score, finding.updated_at FROM finding" + cond
	var total int64
	if !req.AssetView {
		err := c.Slave.WithContext(ctx).Raw("SELECT COUNT(*) FROM ("+base+") scoped", args...).Scan(&total).Error
		return nil, total, uint64(total), err
	}
	grouped := "SELECT project_id, asset_type, asset_key, MIN(asset_name) AS asset_name, MAX(score) AS max_score, COUNT(*) AS finding_count, CONCAT('[', GROUP_CONCAT(DISTINCT JSON_QUOTE(data_source) ORDER BY data_source SEPARATOR ','), ']') AS data_sources_json, UNIX_TIMESTAMP(MAX(updated_at)) AS updated_at FROM (" + base + ") scoped GROUP BY project_id, asset_type, asset_key"
	var counts struct {
		Total        int64
		FindingCount uint64
	}
	countsSQL := "SELECT COUNT(*) AS finding_count FROM (" + base + ") scoped GROUP BY project_id, asset_type, asset_key"
	if err := c.Slave.WithContext(ctx).Raw("SELECT COUNT(*) AS total, COALESCE(SUM(finding_count), 0) AS finding_count FROM ("+countsSQL+") assets", args...).Scan(&counts).Error; err != nil {
		return nil, 0, 0, err
	}
	total = counts.Total
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
	var rows []struct {
		ProjectID                      uint32
		AssetType, AssetKey, AssetName string
		MaxScore                       float32
		FindingCount                   uint32
		UpdatedAt                      int64
		DataSourcesJSON                string
	}
	if err := c.Slave.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, 0, 0, err
	}
	assets := make([]*finding.FindingAsset, 0, len(rows))
	for _, row := range rows {
		var sources []string
		if err := json.Unmarshal([]byte(row.DataSourcesJSON), &sources); err != nil {
			return nil, 0, 0, err
		}
		assets = append(assets, &finding.FindingAsset{ProjectId: row.ProjectID, AssetType: row.AssetType, AssetKey: row.AssetKey, AssetName: row.AssetName, MaxScore: row.MaxScore, FindingCount: row.FindingCount, UpdatedAt: row.UpdatedAt, DataSources: sources})
	}
	return assets, total, counts.FindingCount, nil
}
