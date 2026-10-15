package db

import (
	"context"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ca-risken/core/proto/finding"
)

func TestAssetConditions(t *testing.T) {
	cases := []struct {
		name, typ, key string
		project        uint32
		args           int
	}{{"scope", "", "", 0, 0}, {"exact", "Workload", "key", 5, 3}, {"unknown", "Unknown", "99", 0, 2}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			q, args := assetConditions(&finding.ListFindingRequest{AssetType: c.typ, AssetKey: c.key, AssetProjectId: c.project})
			if len(args) != c.args || !strings.Contains(q, "google:scc") || !strings.Contains(q, "asset_type IS NOT NULL") {
				t.Fatalf("%s %v", q, args)
			}
			if strings.Contains(q, "BINARY") || strings.Contains(q, "CASE") {
				t.Fatal("expression in asset filter")
			}
			if c.key != "" {
				column := "finding.asset_key = ?"
				if c.typ == "Unknown" {
					column = "finding.finding_id = ?"
				}
				if !strings.Contains(q, column) {
					t.Fatal(q)
				}
			}
		})
	}
}
func TestListFindingAssets(t *testing.T) {
	cases := []struct {
		name string
		org  uint32
		view bool
	}{{"project", 0, true}, {"org", 7, true}, {"expansion", 7, false}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, m, err := newMockClient()
			if err != nil {
				t.Fatal(err)
			}
			m.ExpectQuery("SELECT COUNT.*SELECT DISTINCT").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			if c.view {
				m.ExpectQuery("SELECT project_id.*GROUP BY.*ORDER BY max_score").WillReturnRows(sqlmock.NewRows([]string{"project_id", "asset_type", "asset_key", "asset_name", "max_score", "finding_count", "updated_at"}).AddRow(1, "VM", "vm", "vm", 0.9, 2, 100))
			}
			a, n, err := client.ListFindingAssets(context.Background(), &finding.ListFindingRequest{ProjectId: 1, ToScore: 1, AssetView: c.view, Limit: 20}, c.org)
			if err != nil || n != 1 {
				t.Fatalf("%v %d", err, n)
			}
			if c.view && len(a) != 1 {
				t.Fatal(a)
			}
			if err = m.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
