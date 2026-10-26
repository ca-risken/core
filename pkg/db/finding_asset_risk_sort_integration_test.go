package db

import (
	"context"
	"os"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/finding"
	"github.com/stretchr/testify/require"
)

func TestFindingAssetRiskSortMySQL(t *testing.T) {
	if os.Getenv("FINDING_ASSET_TEST_DSN") == "" {
		t.Skip("dedicated MySQL fixture not configured")
	}
	client, err := NewClient(&Config{MasterHost: "db.middleware.svc.cluster.local", SlaveHost: "db.middleware.svc.cluster.local", MasterUser: "root", SlaveUser: "root", Schema: "scc_asset_t003", Port: 3306}, logging.NewLogger())
	require.NoError(t, err)
	cleanup := func() {
		client.Master.Exec("DELETE FROM finding_tag WHERE finding_id BETWEEN 94501 AND 94509")
		client.Master.Exec("DELETE FROM finding WHERE finding_id BETWEEN 94501 AND 94509")
	}
	cleanup()
	defer cleanup()
	require.NoError(t, client.Master.Exec("INSERT IGNORE INTO organization_project VALUES(7,1)").Error)
	fixtures := []struct {
		id        int
		key, data string
	}{
		{94501, "a-none", `{"finding":{"category":"SOFTWARE_VULNERABILITY"}}`},
		{94502, "b-public", `{"finding":{"category":"PUBLIC_SQL_INSTANCE"}}`},
		{94503, "c-exploit", `{"finding":{"vulnerability":{"cve":{"id":"CVE-1","exploitation_activity":3}}}}`},
		{94504, "d-public-exploit", `{"finding":{"category":"PUBLIC_SQL_INSTANCE","vulnerability":{"cve":{"id":"CVE-1","exploitation_activity":3}}}}`},
		{94505, "e-kev", `{"finding":{"vulnerability":{"cve":{"id":"CVE-1"}}},"vulnerability":{"kev":{"cveID":"CVE-1"}}}`},
		{94506, "f-public-kev", `{"finding":{"category":"PUBLIC_SQL_INSTANCE","vulnerability":{"cve":{"id":"CVE-1"}}},"vulnerability":{"kev":{"cveID":"CVE-1"}}}`},
		{94507, "g-kev", `{"finding":{"vulnerability":{"cve":{"id":"CVE-1"}}},"vulnerability":{"kev":{"cveID":"CVE-1"}}}`},
	}
	for _, f := range fixtures {
		require.NoError(t, client.Master.Exec("INSERT INTO finding(finding_id,project_id,data_source,data_source_id,resource_name,asset_type,asset_key,asset_name,original_score,score,data,created_at,updated_at) VALUES(?,1,'google:scc',?,?,'Workload',?,'risk-sort-fixture',0,0.9,?,NOW(),NOW())", f.id, f.key, f.key, f.key, f.data).Error)
	}
	require.NoError(t, client.Master.Exec("INSERT INTO finding(finding_id,project_id,data_source,data_source_id,resource_name,asset_type,asset_key,asset_name,original_score,score,data,created_at,updated_at) VALUES(94508,1,'google:scc','low-kev','a-none','Workload','a-none','risk-sort-fixture',0,0.1,?,NOW(),NOW())", fixtures[4].data).Error)
	cases := []struct {
		name, dir     string
		offset, limit int32
		org           uint32
		want          []string
	}{
		{"descending project", "desc", 0, 3, 0, []string{"f-public-kev", "e-kev", "g-kev"}},
		{"descending second page", "desc", 3, 3, 0, []string{"d-public-exploit", "c-exploit", "b-public"}},
		{"ascending org", "asc", 0, 3, 7, []string{"a-none", "b-public", "c-exploit"}},
		{"descending org", "desc", 0, 3, 7, []string{"f-public-kev", "e-kev", "g-kev"}},
		{"past last page", "desc", 100, 3, 7, []string{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assets, total, _, err := client.ListFindingAssets(context.Background(), &finding.ListFindingRequest{ProjectId: 1, AssetView: true, AssetName: "risk-sort-fixture", FromScore: .8, ToScore: 1, Sort: "risk_information", Direction: tc.dir, Offset: tc.offset, Limit: tc.limit}, tc.org)
			require.NoError(t, err)
			require.Equal(t, int64(7), total)
			got := []string{}
			for _, a := range assets {
				got = append(got, a.AssetKey)
			}
			require.Equal(t, tc.want, got)
		})
	}
	t.Run("tag excludes stronger evidence", func(t *testing.T) {
		require.NoError(t, client.Master.Exec("INSERT INTO finding_tag(finding_tag_id,finding_id,project_id,tag) VALUES(94501,94501,1,'sort-selected'),(94503,94503,1,'sort-selected')").Error)
		assets, _, _, err := client.ListFindingAssets(context.Background(), &finding.ListFindingRequest{ProjectId: 1, AssetView: true, AssetName: "risk-sort-fixture", FromScore: 0, ToScore: 1, Tag: []string{"sort-selected"}, Sort: "risk_information", Direction: "desc", Limit: 20}, 7)
		require.NoError(t, err)
		require.Len(t, assets, 2)
		require.Equal(t, "c-exploit", assets[0].AssetKey)
		require.Equal(t, "a-none", assets[1].AssetKey)
		require.Empty(t, assets[1].RiskInformation)
	})

}
