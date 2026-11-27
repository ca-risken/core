package db

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ca-risken/core/proto/finding"
)

func TestListFindingAssetsCountScope(t *testing.T) {
	cases := []struct {
		name string
		view bool
		tags []string
	}{
		{name: "asset list", view: true},
		{name: "tag duplicates", view: true, tags: []string{"a", "b"}},
		{name: "expanded finding count"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatal(err)
			}
			selection := "SELECT "
			if len(c.tags) > 0 {
				selection += "DISTINCT "
			}
			prefix := `(?s)^SELECT COUNT\(\*\) FROM \(`
			if c.view {
				prefix = `(?s)^SELECT COUNT\(\*\) AS total,.*FROM \(SELECT COUNT\(\*\) AS finding_count FROM \(`
			}
			pattern := prefix + selection + `finding\.finding_id, finding\.project_id, finding\.asset_type, CASE .* AS asset_key FROM finding`
			if c.view {
				mock.ExpectQuery(pattern).WillReturnRows(sqlmock.NewRows([]string{"total", "finding_count"}).AddRow(1, 4))
				mock.ExpectQuery(`(?s)^SELECT project_id, asset_type, asset_key, MIN\(asset_name\).*FROM \(` + selection + `finding\.finding_id,.*finding\.asset_name.*ORDER BY max_score desc, project_id ASC, asset_type ASC, asset_key ASC LIMIT 20 OFFSET 0`).WillReturnRows(sqlmock.NewRows([]string{"project_id", "asset_type", "asset_key", "asset_name", "max_score", "finding_count", "updated_at", "data_sources_json"}).AddRow(1001, "VM", "key", "asset", 0.9, 4, 123, `["google:scc"]`))
			} else {
				mock.ExpectQuery(pattern).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(4))
			}
			assets, total, count, err := client.ListFindingAssets(context.Background(), &finding.ListFindingRequest{ProjectId: 1001, FromScore: 0.8, ToScore: 1, Limit: 20, AssetView: c.view, Tag: c.tags}, 0)
			if err != nil {
				t.Fatal(err)
			}
			if count != 4 {
				t.Fatalf("finding count=%d", count)
			}
			if c.view && (total != 1 || len(assets) != 1 || assets[0].FindingCount != 4 || assets[0].AssetName != "asset") {
				t.Fatalf("asset result changed: %+v total=%d", assets, total)
			}
			if !c.view && (total != 4 || assets != nil) {
				t.Fatalf("expanded count changed: %d", total)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
