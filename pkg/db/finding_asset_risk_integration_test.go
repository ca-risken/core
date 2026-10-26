package db

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/finding"
	"github.com/stretchr/testify/require"
)

func TestFindingAssetRiskMySQL(t *testing.T) {
	if os.Getenv("FINDING_ASSET_TEST_DSN") == "" {
		t.Skip("dedicated MySQL fixture not configured")
	}
	client, err := NewClient(&Config{MasterHost: "db.middleware.svc.cluster.local", SlaveHost: "db.middleware.svc.cluster.local", MasterUser: "root", SlaveUser: "root", Schema: "scc_asset_t003", Port: 3306}, logging.NewLogger())
	require.NoError(t, err)
	cleanup := func() {
		client.Master.Exec("DELETE FROM pend_finding WHERE finding_id BETWEEN 94001 AND 94004")
		client.Master.Exec("DELETE FROM finding WHERE finding_id BETWEEN 94001 AND 94004")
	}
	cleanup()
	defer cleanup()
	require.NoError(t, client.Master.Exec("INSERT IGNORE INTO organization_project VALUES(7,1),(7,2)").Error)
	insert := func(id int, project int, typ, key, name string, score float64, data string) {
		require.NoError(t, client.Master.Exec("INSERT INTO finding(finding_id,project_id,data_source,data_source_id,resource_name,asset_type,asset_key,asset_name,original_score,score,data,created_at,updated_at) VALUES(?,?,'google:scc',?,'risk-sql',NULLIF(?,''),NULLIF(?,''),NULLIF(?,''),0,?,?,NOW(),NOW())", id, project, name+strconv.Itoa(id), typ, key, name, score, data).Error)
	}
	insert(94001, 1, "Resource", "risk-key", "risk-sql", .9, `{"finding":{"category":"SOFTWARE_VULNERABILITY","vulnerability":{"cve":{"id":"CVE-2026-1","exploitation_activity":3}}}}`)
	insert(94002, 1, "", "", "public-setting", .1, `{"finding":{"category":"PUBLIC_SQL_INSTANCE"}}`)
	insert(94003, 3, "Resource", "risk-key", "risk-sql", .9, `{"finding":{"category":"PUBLIC_BUCKET_ACL"}}`)
	insert(94004, 1, "Resource", "different-key", "other-sql", .9, `{"finding":{"vulnerability":{"cve":{"id":"CVE-2026-2","exploitation_activity":3}}}}`)
	// Pending exposure must not appear in the active view.
	cases := []struct {
		name     string
		org      uint32
		pending  bool
		expected int
	}{
		{"project lower score evidence", 0, false, 2},
		{"organization excludes other project", 7, false, 2},
		{"active excludes pending evidence", 7, true, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.pending {
				require.NoError(t, client.Master.Exec("INSERT INTO pend_finding(finding_id,project_id,note) VALUES(94002,1,'risk test')").Error)
			}
			assets, total, _, err := client.ListFindingAssets(context.Background(), &finding.ListFindingRequest{ProjectId: 1, AssetView: true, AssetName: "risk-sql", FromScore: .8, ToScore: 1, Status: finding.FindingStatus_FINDING_ACTIVE, Limit: 20}, tc.org)
			require.NoError(t, err)
			require.Equal(t, int64(1), total)
			require.Len(t, assets, 1)
			require.Len(t, assets[0].RiskInformation, tc.expected)
			for _, r := range assets[0].RiskInformation {
				require.NotEqual(t, uint64(94003), r.FindingId)
				require.NotEqual(t, uint64(94004), r.FindingId)
			}
		})
	}
}
