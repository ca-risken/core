package db

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	"github.com/stretchr/testify/require"
)

func TestFindingPackagesMySQL(t *testing.T) {
	if os.Getenv("FINDING_ASSET_TEST_DSN") == "" {
		t.Skip("dedicated MySQL fixture not configured")
	}
	client, err := NewClient(&Config{MasterHost: "db.middleware.svc.cluster.local", SlaveHost: "db.middleware.svc.cluster.local", MasterUser: "root", SlaveUser: "root", Schema: "scc_asset_t003", Port: 3306}, logging.NewLogger())
	require.NoError(t, err)
	ctx := context.Background()
	sum := sha256.Sum256([]byte("package-integration"))
	key := base64.RawURLEncoding.EncodeToString(sum[:])
	type pkg struct {
		PackageName    string `json:"package_name"`
		PackageType    string `json:"package_type"`
		PackageVersion string `json:"package_version"`
		CPE            string `json:"cpe_uri"`
	}
	records := []struct{ typ, name, platform, current, fixed, raw string }{
		{"NPM", "next", "alpine", "1", "3", ""}, {"NPM", "next", "debian", "2", "4", ""},
		{"OS", "next", "alpine", "1", "", ""}, {"OS", "next", "debian", "1", "2", ""},
		{raw: `{}`}, {raw: `"invalid json"`},
	}
	cleanup := func() {
		client.Master.Exec("DELETE FROM finding_tag WHERE finding_id BETWEEN 90001 AND 90006")
		client.Master.Exec("DELETE FROM pend_finding WHERE finding_id BETWEEN 90001 AND 90006")
		client.Master.Exec("DELETE FROM finding WHERE finding_id BETWEEN 90001 AND 90006")
	}
	cleanup()
	defer cleanup()
	for i, r := range records {
		raw := r.raw
		if raw == "" {
			data, err := json.Marshal(map[string]any{"finding": map[string]any{"vulnerability": map[string]any{"offending_package": pkg{r.name, r.typ, r.current, r.platform}, "fixed_package": pkg{r.name, r.typ, r.fixed, r.platform}}}})
			require.NoError(t, err)
			raw = string(data)
		}
		require.NoError(t, client.BulkUpsertFinding(ctx, []*model.Finding{{FindingID: uint64(90001 + i), ProjectID: 1, DataSource: "google:scc", DataSourceID: string(rune('a' + i)), ResourceName: "fixture", Score: 0.9, Data: raw, AssetType: "Image", AssetKey: key, AssetName: "fixture"}}))
	}
	require.NoError(t, client.Master.Exec("INSERT INTO finding_tag VALUES(90001,90001,1,'package-a'),(90002,90001,1,'package-b')").Error)
	require.NoError(t, client.Master.Exec("INSERT INTO pend_finding(finding_id,project_id) VALUES(90002,1)").Error)
	for _, filter := range []struct {
		name   string
		tags   []string
		status finding.FindingStatus
		from   float32
		want   uint64
	}{
		{name: "duplicate tag matches", tags: []string{"package-a", "package-b"}, want: 1},
		{name: "active", status: finding.FindingStatus_FINDING_ACTIVE, want: 5},
		{name: "pending", status: finding.FindingStatus_FINDING_PENDING, want: 1},
		{name: "score", from: 0.95, want: 0},
	} {
		t.Run(filter.name, func(t *testing.T) {
			groups, err := client.ListFindingPackages(ctx, &finding.ListFindingRequest{ProjectId: 1, AssetProjectId: 1, AssetType: "Image", AssetKey: key, ToScore: 1, FromScore: filter.from, Tag: filter.tags, Status: filter.status}, 0)
			require.NoError(t, err)
			var count uint64
			for _, g := range groups {
				count += g.FindingCount
			}
			require.Equal(t, filter.want, count)
		})
	}
	cases := []struct {
		name    string
		org     uint32
		project uint32
		groups  int
		count   uint64
	}{
		{"project", 0, 1, 4, 6}, {"organization", 7, 1, 4, 6}, {"other project", 0, 2, 0, 0}, {"outside organization", 8, 1, 0, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := &finding.ListFindingRequest{ProjectId: c.project, AssetProjectId: c.project, AssetType: "Image", AssetKey: key, ToScore: 1, Limit: 20, Sort: "finding_id", Direction: "asc"}
			groups, err := client.ListFindingPackages(ctx, req, c.org)
			require.NoError(t, err)
			require.Len(t, groups, c.groups)
			var count uint64
			for _, g := range groups {
				count += g.FindingCount
				req.PackageKey = g.PackageKey
				if c.org == 0 {
					items, err := client.ListFinding(ctx, req)
					require.NoError(t, err)
					require.Len(t, *items, int(g.FindingCount))
				} else {
					items, err := client.ListFindingForOrg(ctx, &finding.ListFindingForOrgRequest{OrganizationId: c.org, AssetProjectId: c.project, AssetType: "Image", AssetKey: key, PackageKey: g.PackageKey, ToScore: 1, Limit: 20, Sort: "finding_id", Direction: "asc"})
					require.NoError(t, err)
					require.Len(t, *items, int(g.FindingCount))
				}
			}
			require.Equal(t, c.count, count)
		})
	}
	pathCases := []struct {
		name, typ      string
		paths          []string
		kind, location string
	}{
		{"app", "NPM", []string{"/app/node_modules/next/package.json"}, "app", "/app"},
		{"pnpm", "NPM", []string{"/app/node_modules/.pnpm/next@16.2.6/node_modules/next/package.json"}, "app", "/app"},
		{"same app", "NPM", []string{"/app/node_modules/next/package.json", "/app/node_modules/.pnpm/next@16.2.9/node_modules/next/package.json"}, "app", "/app"},
		{"different apps", "NPM", []string{"/app/node_modules/next/package.json", "/other/node_modules/next/package.json"}, "location", ""},
		{"missing path", "NPM", nil, "location", ""},
		{"app root", "NPM", []string{"/app/package.json"}, "app", "/app"},
		{"same app different layouts", "NPM", []string{"/app/node_modules/next/package.json", "/app/package.json"}, "app", "/app"},
		{"bundled npm", "NPM", []string{"/usr/local/lib/node_modules/npm/node_modules/next/package.json"}, "location", "/usr/local/lib/node_modules"},
		{"bundled corepack", "NPM", []string{"/usr/lib/node_modules/corepack/node_modules/next/package.json"}, "location", "/usr/lib/node_modules"},
		{"other global package", "NPM", []string{"/usr/lib/node_modules/other/node_modules/next/package.json"}, "location", "/usr/lib/node_modules"},
		{"relative traversal", "NPM", []string{"/usr/local/lib/node_modules/npm/../other/package.json"}, "location", ""},
		{"java app", "MAVEN", []string{"/app/example.jar"}, "app", "/app"},
		{"python app", "PYPI", []string{"/app/venv/lib/python3.12/site-packages/requests/__init__.py"}, "app", "/app"},
		{"python shared", "PYPI", []string{"/usr/local/lib/python3.12/site-packages/pip/_vendor/requests/__init__.py"}, "location", "/usr/local/lib/python3.12/site-packages"},
		{"python dist", "PYPI", []string{"/usr/lib/python3/dist-packages/requests/__init__.py"}, "location", "/usr/lib/python3/dist-packages"},
		{"java shared", "MAVEN", []string{"/opt/service/lib/example.jar"}, "location", "/opt/service/lib"},
		{"other npm app directory", "NPM", []string{"/service/node_modules/.pnpm/next@1.0/node_modules/next/package.json"}, "location", "/service/node_modules"},
		{"different app prefixes", "PYPI", []string{"/app2/site-packages/requests/__init__.py"}, "location", "/app2/site-packages"},
		{"unknown package type with app path", "", []string{"/app/lib/example.jar"}, "app", "/app"},
		{"relative path", "PYPI", []string{"venv/lib/site-packages/a.py"}, "location", "/venv/lib/site-packages"},
		{"different directories", "MAVEN", []string{"/opt/a/lib/a.jar", "/opt/b/lib/a.jar"}, "location", ""},
		{"OS overrides app location", "OS", []string{"/app/node_modules/example/package.json"}, "os", "alpine"},

		{"os", "OS", nil, "os", "alpine"},
	}
	for _, c := range pathCases {
		t.Run("path/"+c.name, func(t *testing.T) {
			files := []map[string]string{}
			for _, p := range c.paths {
				files = append(files, map[string]string{"path": p})
			}
			data, err := json.Marshal(map[string]any{"finding": map[string]any{"files": files, "vulnerability": map[string]any{"offending_package": pkg{"next", c.typ, "1", "alpine"}, "fixed_package": pkg{"next", c.typ, "2", "alpine"}}}})
			require.NoError(t, err)
			require.NoError(t, client.Master.Exec("UPDATE finding SET data=? WHERE finding_id=90001", string(data)).Error)
			req := &finding.ListFindingRequest{ProjectId: 1, AssetProjectId: 1, AssetType: "Image", AssetKey: key, FindingId: 90001, ToScore: 1, Sort: "finding_id", Direction: "asc", Limit: 20}
			groups, err := client.ListFindingPackages(ctx, req, 0)
			require.NoError(t, err)
			require.Len(t, groups, 1)
			require.Equal(t, c.kind, groups[0].WorkType)
			require.Equal(t, c.location, groups[0].WorkLocation)
			req.PackageKey = groups[0].PackageKey
			items, err := client.ListFinding(ctx, req)
			require.NoError(t, err)
			require.Len(t, *items, 1)
		})
	}

}
