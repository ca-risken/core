package db

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ca-risken/core/proto/finding"
)

func TestSCCAssetRisks(t *testing.T) {
	cases := []struct {
		name, data, kind string
		count            int
	}{
		{"KEV only", `{"finding":{"vulnerability":{"cve":{"id":"CVE-2026-1"}}},"vulnerability":{"kev":{"cveID":"CVE-2026-1"}}}`, "exploit", 1},
		{"KEV and exploit retain both", `{"finding":{"vulnerability":{"cve":{"id":"CVE-2026-1","exploitation_activity":3}}},"vulnerability":{"kev":{"cveID":"CVE-2026-1"}}}`, "exploit", 2},
		{"different CVE ignored", `{"finding":{"vulnerability":{"cve":{"id":"CVE-2026-1"}}},"vulnerability":{"kev":{"cveID":"CVE-2026-2"}}}`, "", 0},
		{"empty KEV ignored", `{"finding":{"vulnerability":{"cve":{"id":"CVE-2026-1"}}},"vulnerability":{"kev":{}}}`, "", 0},
		{"null KEV ignored", `{"finding":{"vulnerability":{"cve":{"id":"CVE-2026-1"}}},"vulnerability":{"kev":null}}`, "", 0},
		{"numeric available", `{"finding":{"vulnerability":{"cve":{"id":"CVE-2026-1","exploitation_activity":3}}}}`, "exploit", 1},
		{"named available", `{"finding":{"vulnerability":{"cve":{"id":"CVE-2026-1","exploitation_activity":"AVAILABLE"}}}}`, "exploit", 1},
		{"release timestamp", `{"finding":{"vulnerability":{"cve":{"id":"CVE-2026-1","exploit_release_date":{"seconds":1770000000}}}}}`, "exploit", 1},
		{"release date", `{"finding":{"vulnerability":{"cve":{"id":"CVE-2026-1","exploit_release_date":"2026-01-01T00:00:00Z"}}}}`, "exploit", 1},
		{"invalid date", `{"finding":{"vulnerability":{"cve":{"id":"CVE-2026-1","exploit_release_date":{}}}}}`, "", 0},
		{"no known", `{"finding":{"vulnerability":{"cve":{"id":"CVE-2026-1","exploitation_activity":5}}}}`, "", 0},
		{"confirmed is not public exploit proof", `{"finding":{"vulnerability":{"cve":{"id":"CVE-2026-1","exploitation_activity":2}}}}`, "", 0},
		{"public sql", `{"finding":{"category":"PUBLIC_SQL_INSTANCE"}}`, "internet", 1},
		{"public bucket", `{"finding":{"category":"PUBLIC_BUCKET_ACL"}}`, "internet", 1},
		{"IP insufficient", `{"finding":{"category":"PUBLIC_IP_ADDRESS"}}`, "", 0},
		{"score insufficient", `{"finding":{"attack_exposure":{"state":1,"score":10}}}`, "", 0},
		{"inactive", `{"finding":{"state":2,"category":"PUBLIC_BUCKET_ACL"}}`, "", 0},
		{"invalid", `{`, "", 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			risks := sccAssetRisks(assetRiskFinding{FindingID: 42, UpdatedAt: 100, Data: c.data})
			if len(risks) != c.count {
				t.Fatalf("risks=%+v", risks)
			}
			if c.name == "KEV and exploit retain both" && (risks[0].Detail != "KEV" || risks[1].Detail != "AVAILABLE") {
				t.Fatal(risks)
			}
			if c.count > 0 && (risks[0].Kind != c.kind || risks[0].FindingId != 42 || risks[0].UpdatedAt != 100) {
				t.Fatal(risks[0])
			}
		})
	}
}

func TestLoadAssetRisks(t *testing.T) {
	cases := []struct {
		name, typ, key string
		org            uint32
		expected       int
	}{
		{"project", "VM", "key", 0, 2},
		{"organization", "VM", "key", 7, 2},
		{"workload does not inherit cluster exposure", "Workload", "key", 7, 1},
		{"unknown is per finding", "Unknown", "11", 7, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatal(err)
			}
			rows := sqlmock.NewRows([]string{"finding_id", "asset_index", "same_asset", "data"}).
				AddRow(11, 0, true, `{"finding":{"vulnerability":{"cve":{"id":"CVE-1","exploitation_activity":3}}}}`)
			if tc.typ == "VM" {
				rows.AddRow(12, 0, false, `{"finding":{"category":"PUBLIC_SQL_INSTANCE"}}`)
			}
			query := "SELECT DISTINCT finding.finding_id.*finding.data.*"
			if tc.org != 0 {
				query += "organization_project.*op.organization_id"
			} else {
				query += "finding.project_id"
			}
			mock.ExpectQuery(query).WillReturnRows(rows)
			assets := []*finding.FindingAsset{{ProjectId: 1, AssetType: tc.typ, AssetKey: tc.key, AssetName: "resource"}}
			err = client.loadAssetRisks(context.Background(), &finding.ListFindingRequest{ProjectId: 1, FromScore: 0.9, Tag: []string{"filtered"}, Status: finding.FindingStatus_FINDING_ACTIVE}, tc.org, assets)
			if err != nil {
				t.Fatal(err)
			}
			if len(assets[0].RiskInformation) != tc.expected {
				t.Fatal(assets[0].RiskInformation)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestRiskScopeDoesNotDependOnSearchFilters(t *testing.T) {
	cases := []struct {
		name string
		org  uint32
	}{{"project", 0}, {"org", 7}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatal(err)
			}
			mock.ExpectQuery("SELECT DISTINCT finding.finding_id.*").WithArgs(0, uint32(1), "VM", "key", float32(0), float32(1), tcScope(tc.org)).WillReturnRows(sqlmock.NewRows([]string{"finding_id"}))
			err = client.loadAssetRisks(context.Background(), &finding.ListFindingRequest{ProjectId: 1, FromScore: 0.9, Tag: []string{"filtered"}}, tc.org, []*finding.FindingAsset{{ProjectId: 1, AssetType: "VM", AssetKey: "key", AssetName: "resource"}})
			if err != nil || mock.ExpectationsWereMet() != nil {
				t.Fatal(err)
			}
		})
	}
}

func tcScope(org uint32) uint32 {
	if org != 0 {
		return org
	}
	return 1
}
