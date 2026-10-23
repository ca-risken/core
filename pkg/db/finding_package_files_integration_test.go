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

func TestFindingPackageFileFormatsMySQL(t *testing.T) {
	if os.Getenv("FINDING_ASSET_TEST_DSN") == "" {
		t.Skip("dedicated MySQL fixture not configured")
	}
	client, err := NewClient(&Config{MasterHost: "db.middleware.svc.cluster.local", SlaveHost: "db.middleware.svc.cluster.local", MasterUser: "root", SlaveUser: "root", Schema: "scc_asset_t003", Port: 3306}, logging.NewLogger())
	require.NoError(t, err)
	ctx := context.Background()
	sum := sha256.Sum256([]byte("file-format-integration"))
	key := base64.RawURLEncoding.EncodeToString(sum[:])
	defer client.Master.Exec("DELETE FROM finding WHERE finding_id=90007")
	cases := []struct {
		name, assetType, files, kind, location string
		offendingFiles                         bool
	}{
		{"VM disk path", "VM", `[{"disk_path":{"relative_path":"home/kubernetes/bin/kubelet"}}]`, "location", "/home/kubernetes/bin", false},
		{"Workload relative path", "Workload", `[{"path":"usr/local/bin/gosu","disk_path":{}}]`, "location", "/usr/local/bin", false},
		{"relative application path", "Workload", `[{"path":"app/node_modules/next/package.json"}]`, "app", "/app", false},
		{"absolute path preserved", "Image", `[{"path":"/app/node_modules/next/package.json"}]`, "app", "/app", false},
		{"empty path uses disk path", "VM", `[{"path":"","disk_path":{"relative_path":"home/kubernetes/bin/kubelet"}}]`, "location", "/home/kubernetes/bin", false},
		{"direct path takes precedence", "Workload", `[{"path":"usr/local/bin/gosu","disk_path":{"relative_path":"other/bin/gosu"}}]`, "location", "/usr/local/bin", false},
		{"offending package files", "VM", `[{"disk_path":{"relative_path":"home/kubernetes/bin/kubelet"}}]`, "location", "/home/kubernetes/bin", true},
		{"equivalent paths", "VM", `[{"path":"/home/kubernetes/bin/kubelet"},{"disk_path":{"relative_path":"home/kubernetes/bin/crictl"}}]`, "location", "/home/kubernetes/bin", false},
		{"different directories stay unknown", "VM", `[{"disk_path":{"relative_path":"home/kubernetes/bin/kubelet"}},{"path":"var/lib/containerd/snapshots/1/fs/bin/tool"}]`, "", "", false},
		{"missing paths", "VM", `[{"disk_path":{}}]`, "", "", false},
		{"empty path", "Workload", `[{"path":""}]`, "", "", false},
		{"parent traversal", "VM", `[{"disk_path":{"relative_path":"../home/kubernetes/bin/kubelet"}}]`, "", "", false},
		{"embedded parent traversal", "Workload", `[{"path":"app/../other/package.json"}]`, "", "", false},
		{"current directory traversal", "Workload", `[{"path":"./app/package.json"}]`, "", "", false},
		{"trailing parent directory", "VM", `[{"path":"home/kubernetes/.."}]`, "", "", false},
		{"trailing current directory", "VM", `[{"path":"home/kubernetes/."}]`, "", "", false},
		{"directory path", "VM", `[{"path":"home/kubernetes/bin/"}]`, "", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			pkg := map[string]any{"package_name": "example", "package_type": "GO", "package_version": "1"}
			report := map[string]any{"vulnerability": map[string]any{"offending_package": pkg}}
			if c.offendingFiles {
				pkg["files"] = json.RawMessage(c.files)
			} else {
				report["files"] = json.RawMessage(c.files)
			}
			raw, err := json.Marshal(map[string]any{"finding": report})
			require.NoError(t, err)
			require.NoError(t, client.BulkUpsertFinding(ctx, []*model.Finding{{FindingID: 90007, ProjectID: 1, DataSource: "google:scc", DataSourceID: "file-formats", ResourceName: "fixture", Score: 0.9, Data: string(raw), AssetType: c.assetType, AssetKey: key, AssetName: "fixture"}}))
			req := &finding.ListFindingRequest{ProjectId: 1, AssetProjectId: 1, AssetType: c.assetType, AssetKey: key, FindingId: 90007, ToScore: 1, Sort: "finding_id", Direction: "asc", Limit: 20}
			groups, err := client.ListFindingPackages(ctx, req, 0)
			require.NoError(t, err)
			require.Len(t, groups, 1)
			require.Equal(t, c.kind, groups[0].WorkType)
			require.Equal(t, c.location, groups[0].WorkLocation)
			req.PackageKey = groups[0].PackageKey
			items, err := client.ListFinding(ctx, req)
			require.NoError(t, err)
			require.Len(t, *items, 1)
			require.Equal(t, uint64(90007), (*items)[0].FindingID)
		})
	}
}
