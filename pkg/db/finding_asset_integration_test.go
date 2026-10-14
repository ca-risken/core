package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/findingasset"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	_ "github.com/go-sql-driver/mysql"
)

// Dedicated disposable schema only; enabled explicitly in Kubernetes verification.
func TestFindingAssetsMySQL(t *testing.T) {
	dsn := os.Getenv("FINDING_ASSET_TEST_DSN")
	if dsn == "" {
		t.Skip("dedicated MySQL fixture not configured")
	}
	raw, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer raw.Close()
	client, err := NewClient(&Config{MasterHost: "db.middleware.svc.cluster.local", SlaveHost: "db.middleware.svc.cluster.local", MasterUser: "root", SlaveUser: "root", Schema: "scc_asset_t003", Port: 3306}, logging.NewLogger())
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	for _, q := range []string{"DELETE FROM finding_tag", "DELETE FROM pend_finding", "DELETE FROM organization_project", "DELETE FROM finding", "INSERT INTO organization_project VALUES(7,1),(7,2)"} {
		if _, err = raw.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	data := `{"asset":{"type":"google.compute.Instance","full_name":"vm"},"finding":{"category":"OS_VULNERABILITY","vulnerability":{"cve":{"id":"CVE-2023-1234"}}}}`
	mk := func(id uint64, project uint32, key string) *model.Finding {
		return &model.Finding{FindingID: id, ProjectID: project, DataSource: "google:scc", DataSourceID: fmt.Sprint(id), ResourceName: "reported-resource", Score: 0.9, OriginalScore: 0.9, Data: data, AssetType: "VM", AssetKey: key, AssetName: key}
	}
	entries := []*model.Finding{mk(1, 1, "vm"), mk(2, 1, "vm"), mk(3, 2, "vm"), mk(4, 3, "vm"), mk(5, 1, strings.Repeat("a", 1500)+"X"), mk(6, 1, strings.Repeat("a", 1500)+"Y"), mk(7, 1, "vm ")}
	u1, u2 := mk(8, 1, ""), mk(9, 1, "")
	u1.AssetType = "Unknown"
	u2.AssetType = "Unknown"
	entries = append(entries, u1, u2)
	if err = client.BulkUpsertFinding(ctx, entries); err != nil {
		t.Fatal(err)
	}
	if _, err = raw.Exec("INSERT INTO finding_tag VALUES(1,1,1,'a'),(2,1,1,'b');"); err != nil {
		t.Fatal(err)
	}
	if _, err = raw.Exec("INSERT INTO pend_finding(finding_id,project_id,expired_at) VALUES(2,1,NULL),(3,2,'2000-01-01')"); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name       string
		org        uint32
		req        *finding.ListFindingRequest
		groups     int64
		firstCount uint32
	}{
		{"project", 0, &finding.ListFindingRequest{ProjectId: 1, AssetType: "VM", AssetKey: "vm"}, 1, 2},
		{"org project separated", 7, &finding.ListFindingRequest{AssetType: "VM", AssetKey: "vm"}, 2, 2},
		{"org intersection", 7, &finding.ListFindingRequest{AssetType: "VM", AssetKey: "vm", AssetProjectId: 3}, 0, 0},
		{"active expired pend", 7, &finding.ListFindingRequest{AssetType: "VM", AssetKey: "vm", Status: finding.FindingStatus_FINDING_ACTIVE}, 2, 1},
		{"pending", 7, &finding.ListFindingRequest{AssetType: "VM", AssetKey: "vm", Status: finding.FindingStatus_FINDING_PENDING}, 1, 1},
		{"duplicate tags", 0, &finding.ListFindingRequest{ProjectId: 1, Tag: []string{"a", "b"}}, 1, 1},
		{"unknown separate", 0, &finding.ListFindingRequest{ProjectId: 1, AssetType: "Unknown"}, 2, 1},
		{"binary long keys", 0, &finding.ListFindingRequest{ProjectId: 1, AssetType: "VM"}, 4, 2},
		{"search source intersection", 0, &finding.ListFindingRequest{ProjectId: 1, DataSource: []string{"aws:"}}, 0, 0},
		{"search resource intersection", 0, &finding.ListFindingRequest{ProjectId: 1, ResourceName: []string{"other"}}, 0, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.req.AssetView = true
			c.req.ToScore = 1
			c.req.Limit = 200
			c.req.Sort = "finding_count"
			c.req.Direction = "desc"
			a, n, err := client.ListFindingAssets(ctx, c.req, c.org)
			if err != nil || n != c.groups {
				t.Fatalf("groups=%d want=%d err=%v", n, c.groups, err)
			}
			if n > 0 && a[0].FindingCount != c.firstCount {
				t.Fatalf("count=%d want=%d", a[0].FindingCount, c.firstCount)
			}
		})
	}
	t.Run("expansion exact", func(t *testing.T) {
		r := &finding.ListFindingRequest{ProjectId: 1, ToScore: 1, Limit: 200, Sort: "finding_id", Direction: "asc", AssetType: "VM", AssetKey: "vm", Tag: []string{"a", "b"}}
		_, n, err := client.ListFindingAssets(ctx, r, 0)
		if err != nil || n != 1 {
			t.Fatalf("%d %v", n, err)
		}
		list, err := client.ListFinding(ctx, r)
		if err != nil || len(*list) != 1 {
			t.Fatalf("%v %v", list, err)
		}
	})
	t.Run("save clears old classification", func(t *testing.T) {
		e := mk(10, 1, "old")
		a := findingasset.Classify(e.DataSource, e.ResourceName, data)
		e.AssetType, e.AssetKey, e.AssetName = a.Type, a.Key, a.Name
		if _, err = client.UpsertFinding(ctx, e); err != nil {
			t.Fatal(err)
		}
		e.Data = `{"finding":{"category":"PUBLIC_BUCKET"}}`
		a = findingasset.Classify(e.DataSource, e.ResourceName, e.Data)
		e.AssetType, e.AssetKey, e.AssetName = a.Type, a.Key, a.Name
		if _, err = client.UpsertFinding(ctx, e); err != nil {
			t.Fatal(err)
		}
		got, err := client.GetFinding(ctx, 1, 10, true)
		if err != nil || got.AssetType != "" {
			t.Fatalf("%+v %v", got, err)
		}
	})
	t.Run("backfill preserves timestamp and state", func(t *testing.T) {
		if _, err = raw.Exec("UPDATE finding SET updated_at='2020-01-01' WHERE finding_id=1"); err != nil {
			t.Fatal(err)
		}
		last, n, err := findingasset.BackfillChunk(ctx, raw, 0, 2)
		if err != nil || last != 2 || n != 2 {
			t.Fatalf("%d %d %v", last, n, err)
		}
		var stamp time.Time
		if err = raw.QueryRow("SELECT updated_at FROM finding WHERE finding_id=1").Scan(&stamp); err != nil || stamp.Year() != 2020 {
			t.Fatalf("%v %v", stamp, err)
		}
		var pending int
		if err = raw.QueryRow("SELECT COUNT(*) FROM pend_finding").Scan(&pending); err != nil || pending != 2 {
			t.Fatal(pending, err)
		}
	})
	t.Run("backfill waits for latest row", func(t *testing.T) {
		tx, err := raw.Begin()
		if err != nil {
			t.Fatal(err)
		}
		newData := strings.Replace(data, `"full_name":"vm"`, `"full_name":"new-vm"`, 1)
		if _, err = tx.Exec("UPDATE finding SET data=? WHERE finding_id=1", newData); err != nil {
			t.Fatal(err)
		}
		done := make(chan error, 1)
		go func() { _, _, err := findingasset.BackfillChunk(ctx, raw, 0, 1); done <- err }()
		select {
		case err := <-done:
			_ = tx.Rollback()
			t.Fatalf("lock bypassed: %v", err)
		case <-time.After(100 * time.Millisecond):
		}
		if err = tx.Commit(); err != nil {
			t.Fatal(err)
		}
		if err = <-done; err != nil {
			t.Fatal(err)
		}
		var key string
		if err = raw.QueryRow("SELECT asset_key FROM finding WHERE finding_id=1").Scan(&key); err != nil || key != "new-vm" {
			t.Fatalf("%s %v", key, err)
		}
	})
}
