package db

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ca-risken/core/proto/finding"
)

func TestListFindingForOrgPage(t *testing.T) {
	cases := []struct {
		name, sort, direction, assetType string
		offset                           int32
		fail                             bool
	}{
		{name: "ID page", sort: "finding_id", direction: "asc"},
		{name: "score page with offset", sort: "score", direction: "desc", offset: 20},
		{name: "asset filter retains existing query", sort: "finding_id", direction: "asc", assetType: "google.compute.Instance"},
		{name: "database error", sort: "finding_id", direction: "asc", fail: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatal(err)
			}
			req := &finding.ListFindingForOrgRequest{OrganizationId: 1019, FromScore: 0.8, ToScore: 1, Status: finding.FindingStatus_FINDING_ACTIVE, Sort: c.sort, Direction: c.direction, Offset: c.offset, Limit: 20, AssetType: c.assetType}
			pattern := `(?s)^select finding\.\*,.*from \(select finding\.finding_id, f_alias\.` + regexp.QuoteMeta(c.sort) + ` as page_sort.*inner join finding f_alias using\(finding_id\).*op\.organization_id = \?.*pf\.expired_at <= NOW\(\).*` + regexp.QuoteMeta(fmt.Sprintf("order by f_alias.%s %s limit %d, 20", c.sort, c.direction, c.offset)) + `\) finding_page.*inner join finding on finding\.finding_id = finding_page\.finding_id.*left join pend_finding.*left join user.*order by finding_page\.page_sort ` + c.direction + `$`
			if c.assetType != "" {
				pattern = `(?s)^select distinct finding\.\*,.*order by finding\.finding_id asc limit 0, 20$`
			}
			expected := mock.ExpectQuery(pattern)
			if c.fail {
				expected.WillReturnError(errors.New("database unavailable"))
			} else {
				expected.WillReturnRows(sqlmock.NewRows([]string{"finding_id", "project_id", "description", "data", "pend_user_id", "pend_note", "pend_user_name"}).AddRow(12, 1001, "Finding detail", `{"key":"value"}`, 7, "note", "Example"))
			}
			got, err := client.ListFindingForOrg(context.Background(), req)
			if (err != nil) != c.fail {
				t.Fatalf("unexpected error: %v", err)
			}
			if !c.fail {
				if got == nil || len(*got) != 1 {
					t.Fatalf("unexpected findings: %+v", got)
				}
				row := (*got)[0]
				if row.FindingID != 12 || row.ProjectID != 1001 || row.Description != "Finding detail" || row.Data != `{"key":"value"}` || row.PendUserID == nil || *row.PendUserID != 7 || row.PendNote == nil || *row.PendNote != "note" || row.PendUserName == nil || *row.PendUserName != "Example" {
					t.Fatalf("full detail or pending fields lost: %+v", row)
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
