package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/findingasset"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	"github.com/stretchr/testify/require"
)

func TestAWSAssetsMySQL(t *testing.T) {
	dsn := os.Getenv("FINDING_ASSET_TEST_DSN")
	if dsn == "" {
		t.Skip("dedicated MySQL fixture not configured")
	}
	raw, err := sql.Open("mysql", dsn)
	require.NoError(t, err)
	defer raw.Close()
	c, err := NewClient(&Config{MasterHost: "db.middleware.svc.cluster.local", SlaveHost: "db.middleware.svc.cluster.local", MasterUser: "root", SlaveUser: "root", Schema: "scc_asset_t003", Port: 3306}, logging.NewLogger())
	require.NoError(t, err)
	for _, q := range []string{"DELETE FROM finding_tag", "DELETE FROM pend_finding", "DELETE FROM organization_project", "DELETE FROM finding", "INSERT INTO organization_project VALUES(7,1),(7,2)"} {
		_, err = raw.Exec(q)
		require.NoError(t, err)
	}
	const arn = "arn:aws:iam::123456789012:user/team/alice"
	inputs := []struct {
		source, data string
		project      uint32
		score        float32
	}{
		{"aws:access-analyzer", `{"Resource":"` + arn + `"}`, 1, .9},
		{"aws:admin-checker", `{"user_arn":"` + arn + `"}`, 1, .6},
		{"aws:cloudsploit", `{"data":{"Resource":"` + arn + `"}}`, 1, .8},
		{"aws:admin-checker", `{"user_arn":"` + arn + `"}`, 2, .7},
		{"aws:admin-checker", `{"user_arn":"` + arn + `"}`, 3, .7},
		{"aws:guard-duty", `{"Resource":{"ResourceType":"AccessKey"}}`, 1, .9},
		{"aws:guard-duty", `{}`, 1, .9},
		{"aws:portscan", `{"data":{"Resource":"` + arn + `"}}`, 1, .9},
	}
	var entries []*model.Finding
	for i, input := range inputs {
		a := findingasset.Classify(input.source, "", input.data)
		entries = append(entries, &model.Finding{FindingID: uint64(i + 1), ProjectID: input.project, DataSource: input.source, DataSourceID: fmt.Sprint(i + 1), ResourceName: arn, Score: input.score, OriginalScore: input.score, Data: input.data, AssetType: a.Type, AssetKey: a.Key, AssetName: a.Name})
	}
	require.NoError(t, c.BulkUpsertFinding(context.Background(), entries))
	_, err = raw.Exec("INSERT INTO finding_tag VALUES(1,1,1,'a'),(2,1,1,'b');")
	require.NoError(t, err)
	_, err = raw.Exec("INSERT INTO pend_finding(finding_id,project_id,expired_at) VALUES(2,1,NULL)")
	require.NoError(t, err)
	key := entries[0].AssetKey
	cases := []struct {
		name   string
		org    uint32
		req    *finding.ListFindingRequest
		groups int64
		count  uint64
	}{
		{"project cross source", 0, &finding.ListFindingRequest{ProjectId: 1, AssetType: "IAMUser", AssetKey: key}, 1, 3},
		{"org separated and scoped", 7, &finding.ListFindingRequest{AssetType: "IAMUser", AssetKey: key}, 2, 4},
		{"org forbidden intersection", 7, &finding.ListFindingRequest{AssetType: "IAMUser", AssetKey: key, AssetProjectId: 3}, 0, 0},
		{"score", 0, &finding.ListFindingRequest{ProjectId: 1, AssetType: "IAMUser", AssetKey: key, FromScore: .85}, 1, 1},
		{"tags deduplicated", 0, &finding.ListFindingRequest{ProjectId: 1, Tag: []string{"a", "b"}}, 1, 1},
		{"source", 0, &finding.ListFindingRequest{ProjectId: 1, DataSource: []string{"aws:admin-checker"}}, 1, 1},
		{"active", 0, &finding.ListFindingRequest{ProjectId: 1, AssetType: "IAMUser", AssetKey: key, Status: finding.FindingStatus_FINDING_ACTIVE}, 1, 2},
		{"unknown not lost", 0, &finding.ListFindingRequest{ProjectId: 1, AssetType: "Unknown"}, 2, 2},
		{"unsupported excluded", 0, &finding.ListFindingRequest{ProjectId: 1, DataSource: []string{"aws:portscan"}}, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := tc.req
			req.ToScore = 1
			req.Limit = 20
			req.AssetView = true
			a, n, count, err := c.ListFindingAssets(context.Background(), req, tc.org)
			require.NoError(t, err)
			require.Equal(t, tc.groups, n)
			require.Equal(t, tc.count, count)
			require.Len(t, a, int(tc.groups))
			for _, asset := range a {
				require.Empty(t, asset.RiskInformation)
			}
		})
	}
	t.Run("categories and exact inner counts", func(t *testing.T) {
		req := &finding.ListFindingRequest{ProjectId: 1, AssetType: "IAMUser", AssetKey: key, AssetProjectId: 1, ToScore: 1, Limit: 20, Sort: "finding_id", Direction: "asc"}
		groups, err := c.ListFindingPackages(context.Background(), req, 0)
		require.NoError(t, err)
		require.Len(t, groups, 3)
		for _, g := range groups {
			require.Contains(t, []string{"review"}, g.WorkType)
			req.PackageKey = g.PackageKey
			rows, err := c.ListFinding(context.Background(), req)
			require.NoError(t, err)
			require.Equal(t, g.FindingCount, uint64(len(*rows)))
			_, n, count, err := c.ListFindingAssets(context.Background(), req, 0)
			require.NoError(t, err)
			require.Equal(t, int64(g.FindingCount), n)
			require.Equal(t, g.FindingCount, count)
		}
	})
	t.Run("organization tagged details", func(t *testing.T) {
		rows, err := c.ListFindingForOrg(context.Background(), &finding.ListFindingForOrgRequest{OrganizationId: 7, Tag: []string{"a", "b"}, AssetType: "IAMUser", AssetKey: key, AssetProjectId: 1, ToScore: 1, Limit: 20, Sort: "finding_id", Direction: "asc"})
		require.NoError(t, err)
		require.Len(t, *rows, 1)
	})
	t.Run("unknown target retains category", func(t *testing.T) {
		req := &finding.ListFindingRequest{ProjectId: 1, AssetType: "Unknown", AssetKey: "7", AssetProjectId: 1, ToScore: 1}
		groups, err := c.ListFindingPackages(context.Background(), req, 0)
		require.NoError(t, err)
		require.Len(t, groups, 1)
		require.Equal(t, "review", groups[0].WorkType)
	})
}
