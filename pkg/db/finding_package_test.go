package db

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ca-risken/core/proto/finding"
	"github.com/stretchr/testify/require"
)

func TestListFindingPackages(t *testing.T) {
	cases := []struct {
		name  string
		org   uint32
		empty bool
	}{{name: "project"}, {name: "organization", org: 7}, {name: "empty", empty: true}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, m, err := newMockClient()
			require.NoError(t, err)
			rows := sqlmock.NewRows([]string{"data_source", "method", "finding_count"})
			if !c.empty {
				rows.AddRow("google:scc", "vulnerability", 5).AddRow("google:asset", "review", 1)
			}
			m.ExpectQuery("SELECT data_source, method").WillReturnRows(rows)
			groups, err := client.ListFindingPackages(context.Background(), &finding.ListFindingRequest{ProjectId: 1, AssetType: "Image", AssetKey: "key", AssetProjectId: 1, ToScore: 1}, c.org)
			require.NoError(t, err)
			if c.empty {
				require.Empty(t, groups)
			} else {
				require.Len(t, groups, 2)
				require.Equal(t, uint64(5), groups[0].FindingCount)
				require.Equal(t, "vulnerability", groups[0].WorkType)
				require.Equal(t, `["google:scc","","","vulnerability:"]`, groups[0].PackageKey)
				require.Empty(t, groups[0].Versions)
			}
			require.NoError(t, m.ExpectationsWereMet())
		})
	}
}
