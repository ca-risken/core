package db

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	"github.com/stretchr/testify/require"
)

func TestFindingAssetNameMySQL(t *testing.T) {
	if os.Getenv("FINDING_ASSET_TEST_DSN") == "" {
		t.Skip("dedicated MySQL fixture not configured")
	}
	client, err := NewClient(&Config{MasterHost: "db.middleware.svc.cluster.local", SlaveHost: "db.middleware.svc.cluster.local", MasterUser: "root", SlaveUser: "root", Schema: "scc_asset_t003", Port: 3306}, logging.NewLogger())
	require.NoError(t, err)
	ctx := context.Background()
	cleanup := func() { client.Master.Exec("DELETE FROM finding WHERE finding_id BETWEEN 92001 AND 92028") }
	cleanup()
	defer cleanup()
	require.NoError(t, client.Master.Exec("INSERT IGNORE INTO organization_project VALUES(7,1),(7,2)").Error)
	hash := func(s string) string {
		sum := sha256.Sum256([]byte(s))
		return base64.RawURLEncoding.EncodeToString(sum[:])
	}
	key := hash("asset-name-24")
	entries := []*model.Finding{}
	for i := 0; i < 28; i++ {
		name := fmt.Sprintf("sample-%02d", i)
		project := uint32(1)
		assetKey := hash(fmt.Sprintf("asset-name-%d", i))
		if i >= 24 {
			name = "stage-25%_日本語"
			assetKey = key
		}
		if i == 26 {
			project = 2
		}
		if i == 27 {
			project = 3
		}
		entries = append(entries, &model.Finding{FindingID: uint64(92001 + i), ProjectID: project, DataSource: "google:scc", DataSourceID: fmt.Sprintf("name-search-%d", i), ResourceName: "asset-name-fixture", Score: .9, AssetType: "VM", AssetName: name, AssetKey: assetKey, Data: `{"finding":{"files":[{"path":"/app/package.json"}],"vulnerability":{"offending_package":{"package_name":"example","package_type":"NPM","package_version":"1"},"fixed_package":{"package_version":"2"}}}}`})
	}
	require.NoError(t, client.BulkUpsertFinding(ctx, entries))
	cases := []struct {
		name, search string
		org, project uint32
		groups       int64
		count        uint64
	}{
		{"all pages", "", 0, 0, 25, 26}, {"partial last page", "stage", 0, 0, 1, 2}, {"literal percent", "%", 0, 0, 1, 2}, {"literal underscore", "_", 0, 0, 1, 2}, {"unicode", "日本語", 0, 0, 1, 2}, {"no match", "absent", 0, 0, 0, 0}, {"quote is literal", "' OR 1=1", 0, 0, 0, 0}, {"organization", "stage", 7, 0, 2, 3}, {"organization project", "stage", 7, 1, 1, 2}, {"outside organization", "stage", 7, 3, 0, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := &finding.ListFindingRequest{ProjectId: 1, AssetView: true, AssetName: c.search, AssetProjectId: c.project, ResourceName: []string{"asset-name-fixture"}, ToScore: 1, Sort: "asset_name", Direction: "asc", Limit: 20}
			assets, total, count, err := client.ListFindingAssets(ctx, req, c.org)
			require.NoError(t, err)
			require.Equal(t, c.groups, total)
			require.Equal(t, c.count, count)
			require.Len(t, assets, min(int(c.groups), 20))
			if c.search == "" {
				for _, a := range assets {
					require.NotEqual(t, "stage-25%_日本語", a.AssetName)
				}
			}
		})
	}
	for _, c := range []struct {
		name, search string
		want         int
	}{{"matching details", "stage", 2}, {"nonmatching details", "absent", 0}} {
		t.Run(c.name, func(t *testing.T) {
			req := &finding.ListFindingRequest{ProjectId: 1, AssetType: "VM", AssetKey: key, AssetProjectId: 1, AssetName: c.search, ResourceName: []string{"asset-name-fixture"}, ToScore: 1, Sort: "finding_id", Direction: "asc", Limit: 20}
			items, err := client.ListFinding(ctx, req)
			require.NoError(t, err)
			require.Len(t, *items, c.want)
			packages, err := client.ListFindingPackages(ctx, req, 0)
			require.NoError(t, err)
			var count uint64
			for _, p := range packages {
				count += p.FindingCount
			}
			require.Equal(t, uint64(c.want), count)
			orgItems, err := client.ListFindingForOrg(ctx, &finding.ListFindingForOrgRequest{OrganizationId: 7, AssetType: "VM", AssetKey: key, AssetProjectId: 1, AssetName: c.search, ResourceName: []string{"asset-name-fixture"}, ToScore: 1, Sort: "finding_id", Direction: "asc", Limit: 20})
			require.NoError(t, err)
			require.Len(t, *orgItems, c.want)
		})
	}
}
