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

func TestListFindingSelectIDs(t *testing.T) {
	cases := []struct {
		name      string
		sort      string
		direction string
		offset    int32
		assetType string
		fail      bool
	}{
		{name: "ascending IDs", sort: "finding_id", direction: "asc"},
		{name: "descending scores with offset", sort: "score", direction: "desc", offset: 200},
		{name: "asset search retains full findings", sort: "finding_id", direction: "asc", assetType: "google.compute.Instance"},
		{name: "database error", sort: "finding_id", direction: "asc", fail: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, mock, err := newMockClient()
			if err != nil {
				t.Fatal(err)
			}
			req := &finding.ListFindingRequest{ProjectId: 1001, FromScore: 0.8, ToScore: 1, Status: finding.FindingStatus_FINDING_ACTIVE, Sort: c.sort, Direction: c.direction, Offset: c.offset, Limit: 200, AssetType: c.assetType}
			query := "select finding.finding_id from finding inner join finding f_alias using(finding_id) left join pend_finding pf using(finding_id) where finding.score between ? and ? and finding.project_id = ? and (pf.finding_id is NULL or (pf.expired_at is not NULL and pf.expired_at <= NOW()))"
			query += fmt.Sprintf(" order by f_alias.%s %s limit %d, 200", c.sort, c.direction, c.offset)
			var expect *sqlmock.ExpectedQuery
			if c.assetType != "" {
				expect = mock.ExpectQuery("^" + regexp.QuoteMeta("select distinct finding.* from finding"))
			} else {
				expect = mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(float32(0.8), float32(1), uint32(1001))
			}
			if c.fail {
				expect.WillReturnError(errors.New("database unavailable"))
			} else {
				expect.WillReturnRows(sqlmock.NewRows([]string{"finding_id", "project_id"}).AddRow(12, 1001).AddRow(34, 1001))
			}
			got, err := client.ListFinding(context.Background(), req)
			if (err != nil) != c.fail {
				t.Fatalf("unexpected error: %v", err)
			}
			if !c.fail && (got == nil || len(*got) != 2 || (*got)[0].FindingID != 12 || (*got)[1].FindingID != 34) {
				t.Fatalf("unexpected findings: %+v", got)
			}
			if !c.fail && c.assetType != "" && (*got)[0].ProjectID != 1001 {
				t.Fatal("asset search must retain full Finding fields")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
