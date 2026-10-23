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
		name         string
		organization uint32
		empty        bool
	}{
		{name: "project"}, {name: "organization", organization: 7}, {name: "empty", empty: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, m, err := newMockClient()
			require.NoError(t, err)
			rows := sqlmock.NewRows([]string{"package_key", "current_version", "fixed_version", "finding_count"})
			if !c.empty {
				rows.AddRow(`["NPM", "next", ""]`, "16.2.6", "16.3.6", 2).
					AddRow(`["NPM", "next", ""]`, "16.2.9", "16.3.8", 3).
					AddRow(`["OS", "next", "alpine"]`, "1", "", 1).
					AddRow(`["", "", ""]`, "", "", 1)
			}
			m.ExpectQuery("SELECT package_key, current_version, fixed_version, COUNT").WillReturnRows(rows)
			groups, err := client.ListFindingPackages(context.Background(), &finding.ListFindingRequest{ProjectId: 1, AssetType: "Image", AssetKey: "key", AssetProjectId: 1, ToScore: 1}, c.organization)
			require.NoError(t, err)
			if c.empty {
				require.Empty(t, groups)
			} else {
				require.Len(t, groups, 3)
				require.Equal(t, "next", groups[0].PackageName)
				require.Equal(t, uint64(5), groups[0].FindingCount)
				require.Len(t, groups[0].Versions, 2)
				var total uint64
				for _, g := range groups {
					total += g.FindingCount
				}
				require.Equal(t, uint64(7), total)
			}
			require.NoError(t, m.ExpectationsWereMet())
		})
	}
}
