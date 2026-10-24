package db

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"github.com/ca-risken/core/proto/finding"
)

func packageFieldSQL(pkg, field string) string {
	document := "IF(JSON_VALID(finding.data), finding.data, '{}')"
	return "COALESCE(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(" + document + ", '$.finding.vulnerability." + pkg + "." + field + "')), 'null'), '')"
}

func workScopeSQL() string {
	document := "IF(JSON_VALID(finding.data), finding.data, '{}')"
	files := "JSON_MERGE_PRESERVE(COALESCE(JSON_EXTRACT(" + document + ", '$.finding.files'), JSON_ARRAY()), COALESCE(JSON_EXTRACT(" + document + ", '$.finding.vulnerability.offending_package.files'), JSON_ARRAY()))"
	filePath := "COALESCE(NULLIF(file_path, ''), disk_relative_path)"
	path := "(CONVERT(CASE WHEN LEFT(" + filePath + ", 1) = '/' THEN " + filePath + " ELSE CONCAT('/', " + filePath + ") END USING utf8mb4) COLLATE utf8mb4_bin)"
	root := "CASE WHEN " + path + " IS NULL OR " + path + " NOT LIKE '/%/%' OR " + path + " LIKE '%/../%' OR " + path + " LIKE '%/./%' OR " + path + " LIKE '%/..' OR " + path + " LIKE '%/.' OR " + path + " LIKE '%/' THEN '' WHEN " + path + " LIKE '/app/%' THEN 'app:/app' WHEN " + path + " LIKE '%/node_modules/%' THEN CONCAT('location:', SUBSTRING_INDEX(" + path + ", '/node_modules/', 1), '/node_modules') WHEN " + path + " LIKE '%/site-packages/%' THEN CONCAT('location:', SUBSTRING_INDEX(" + path + ", '/site-packages/', 1), '/site-packages') WHEN " + path + " LIKE '%/dist-packages/%' THEN CONCAT('location:', SUBSTRING_INDEX(" + path + ", '/dist-packages/', 1), '/dist-packages') ELSE CONCAT('location:', REGEXP_REPLACE(" + path + ", '/[^/]+$', '')) END"

	locations := "SELECT " + root + " AS location FROM JSON_TABLE(" + files + ", '$[*]' COLUMNS(file_path VARCHAR(8192) PATH '$.path' NULL ON EMPTY NULL ON ERROR, disk_relative_path VARCHAR(8192) PATH '$.disk_path.relative_path' NULL ON EMPTY NULL ON ERROR)) paths"
	pathScope := "(SELECT IF(COUNT(*) > 0 AND COUNT(DISTINCT location) = 1 AND MIN(location) <> '', MIN(location), '') FROM (" + locations + ") locations)"
	typ := packageFieldSQL("offending_package", "package_type")
	platform := packageFieldSQL("offending_package", "cpe_uri")
	return "CASE WHEN " + typ + " = 'OS' THEN CONCAT('os:', " + platform + ") WHEN " + packageFieldSQL("offending_package", "package_name") + " = '' THEN '' ELSE COALESCE(NULLIF(" + pathScope + ", ''), 'location:') END"
}

func packageKeySQL() string {
	name := packageFieldSQL("offending_package", "package_name")
	typ := packageFieldSQL("offending_package", "package_type")
	platform := packageFieldSQL("offending_package", "cpe_uri")
	return "JSON_ARRAY(IF(" + name + " = '', '', " + typ + "), " + name + ", IF(" + name + " <> '' AND " + typ + " = 'OS', " + platform + ", ''), " + workScopeSQL() + ")"
}

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
	base := "SELECT DISTINCT finding.finding_id, CAST(" + packageKeySQL() + " AS BINARY) AS package_key, CAST(" + packageFieldSQL("offending_package", "package_version") + " AS BINARY) AS current_version, CAST(" + packageFieldSQL("fixed_package", "package_version") + " AS BINARY) AS fixed_version FROM finding" + cond
	var rows []struct {
		PackageKey     string
		CurrentVersion string
		FixedVersion   string
		FindingCount   uint64
	}
	query := "SELECT package_key, current_version, fixed_version, COUNT(*) AS finding_count FROM (" + base + ") scoped GROUP BY package_key, current_version, fixed_version ORDER BY package_key, current_version, fixed_version"
	if err := c.Slave.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	groups := []*finding.FindingPackage{}
	byKey := map[string]*finding.FindingPackage{}
	for _, row := range rows {
		group := byKey[row.PackageKey]
		if group == nil {
			var parts []string
			if err := json.Unmarshal([]byte(row.PackageKey), &parts); err != nil {
				return nil, err
			}
			group = &finding.FindingPackage{PackageKey: row.PackageKey, PackageType: parts[0], PackageName: parts[1], Platform: parts[2]}
			if len(parts) > 3 && parts[3] != "" {
				group.WorkType, group.WorkLocation, _ = strings.Cut(parts[3], ":")
			}
			byKey[row.PackageKey] = group
			groups = append(groups, group)
		}
		group.FindingCount += row.FindingCount
		group.Versions = append(group.Versions, &finding.FindingPackageVersion{CurrentVersion: row.CurrentVersion, FixedVersion: row.FixedVersion, FindingCount: row.FindingCount})
	}
	sort.Slice(groups, func(i, j int) bool {
		if groups[i].FindingCount != groups[j].FindingCount {
			return groups[i].FindingCount > groups[j].FindingCount
		}
		return groups[i].PackageKey < groups[j].PackageKey
	})
	return groups, nil
}
