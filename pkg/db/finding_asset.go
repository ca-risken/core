package db

import (
	"context"
	"encoding/json"
	"fmt"
	sortutil "sort"
	"strconv"

	"github.com/ca-risken/core/proto/finding"
)

const assetKeySQL = "CASE WHEN finding.asset_type = 'Unknown' THEN CAST(finding.finding_id AS CHAR) ELSE finding.asset_key END"

// assetConditions is appended after the existing permission/search/status scope.
func assetConditions(req *finding.ListFindingRequest) (string, []any) {
	cond := " and finding.data_source = 'google:scc' and finding.asset_type IS NOT NULL"
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
		cond += " and CAST(" + packageKeySQL() + " AS BINARY) = CAST(? AS BINARY)"
		args = append(args, req.PackageKey)
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
	case "updated_at", "asset_name", "asset_type", "finding_count", "project_id", "risk_information":
	default:
		sort = "max_score"
	}
	direction := "desc"
	if req.Direction == "asc" {
		direction = "asc"
	}
	query := grouped + fmt.Sprintf(" ORDER BY %s %s, project_id ASC, asset_type ASC, asset_key ASC LIMIT %d OFFSET %d", sort, direction, req.Limit, req.Offset)
	if req.Sort == "risk_information" {
		query = grouped + " ORDER BY project_id, asset_type, asset_key"
	}
	var rows []struct {
		RiskRank                       int
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
	if req.Sort == "risk_information" {
		identities := make([]struct {
			ProjectID           uint32
			AssetType, AssetKey string
		}, len(rows))
		for i, row := range rows {
			identities[i].ProjectID = row.ProjectID
			identities[i].AssetType = row.AssetType
			identities[i].AssetKey = row.AssetKey
		}
		pageJSON, err := json.Marshal(identities)
		if err != nil {
			return nil, 0, 0, err
		}
		pageSQL := `SELECT ordinal - 1 AS asset_index, project_id, asset_type, asset_key FROM JSON_TABLE(?, '$[*]' COLUMNS (ordinal FOR ORDINALITY, project_id INT PATH '$.ProjectID', asset_type VARCHAR(32) PATH '$.AssetType', asset_key VARCHAR(64) PATH '$.AssetKey')) ranked`
		if err := c.scanAssetRisks(ctx, req, organizationID, pageSQL, []any{string(pageJSON)}, func(index int, risk *finding.FindingAssetRisk) {
			rows[index].RiskRank = assetRiskSortRank(rows[index].RiskRank, risk)
		}); err != nil {
			return nil, 0, 0, err
		}
		sortutil.SliceStable(rows, func(i, j int) bool {
			if direction == "asc" {
				return rows[i].RiskRank < rows[j].RiskRank
			}
			return rows[i].RiskRank > rows[j].RiskRank
		})
		start := min(int(req.Offset), len(rows))
		end := min(start+int(req.Limit), len(rows))
		rows = rows[start:end]
	}
	assets := make([]*finding.FindingAsset, 0, len(rows))
	for _, row := range rows {
		var sources []string
		if err := json.Unmarshal([]byte(row.DataSourcesJSON), &sources); err != nil {
			return nil, 0, 0, err
		}
		assets = append(assets, &finding.FindingAsset{ProjectId: row.ProjectID, AssetType: row.AssetType, AssetKey: row.AssetKey, AssetName: row.AssetName, MaxScore: row.MaxScore, FindingCount: row.FindingCount, UpdatedAt: row.UpdatedAt, DataSources: sources})
	}
	if err := c.loadAssetRisks(ctx, req, organizationID, assets); err != nil {
		return nil, 0, 0, err
	}
	return assets, total, counts.FindingCount, nil
}

func assetRiskSortRank(current int, risk *finding.FindingAssetRisk) int {
	if risk.Kind == "internet" {
		return current | 1
	}
	if risk.Kind != "exploit" {
		return current
	}
	level := 4
	if risk.Detail == "KEV" {
		level = 6
	}
	return max(current&^1, level) | (current & 1)
}
