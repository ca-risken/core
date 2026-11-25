package finding

import (
	"context"
	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/proto/finding"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListFindingPackagesPagination(t *testing.T) {
	cases := []struct {
		name          string
		offset, limit int32
		want          int
		org           uint32
	}{
		{"first page", 0, 1, 2, 0}, {"second page", 1, 1, 1, 0}, {"out of range", 5, 1, 0, 0}, {"org", 0, 20, 3, 7},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := mocks.NewFindingRepository(t)
			repo.On("ListFindingPackages", mock.Anything, mock.Anything, c.org).Return([]*finding.FindingPackage{{PackageName: "next", FindingCount: 3, WorkType: "vulnerability", WorkLocation: "google:scc"}, {PackageName: "openssl", FindingCount: 2, WorkType: "public_access", WorkLocation: "google:asset"}, {PackageName: "undici", FindingCount: 4, WorkType: "vulnerability", WorkLocation: "google:scc"}}, nil).Once()
			service := &FindingService{repository: repo}
			groups, total, count, err := service.listFindingPackages(context.Background(), &finding.ListFindingRequest{Offset: c.offset, Limit: c.limit}, c.org)
			require.NoError(t, err)
			require.Len(t, groups, c.want)
			require.Equal(t, uint32(2), total)
			require.Equal(t, uint64(9), count)
		})
	}
}

func TestReviewWorkGroupLast(t *testing.T) {
	cases := []struct {
		name   string
		offset int32
		want   string
	}{
		{"known work first", 0, "next"}, {"review last", 1, "unknown"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := mocks.NewFindingRepository(t)
			repo.On("ListFindingPackages", mock.Anything, mock.Anything, uint32(0)).Return([]*finding.FindingPackage{{PackageName: "unknown", FindingCount: 100, WorkType: "review"}, {PackageName: "next", FindingCount: 1, WorkType: "vulnerability", WorkLocation: "google:scc"}}, nil).Once()
			service := &FindingService{repository: repo}
			groups, total, count, err := service.listFindingPackages(context.Background(), &finding.ListFindingRequest{Offset: c.offset, Limit: 1}, 0)
			require.NoError(t, err)
			require.Len(t, groups, 1)
			require.Equal(t, c.want, groups[0].PackageName)
			require.Equal(t, uint32(2), total)
			require.Equal(t, uint64(101), count)
		})
	}
}

func TestSourceMethodWorkGroups(t *testing.T) {
	cases := []struct {
		name, kind               string
		wantGroups, wantPackages int
	}{
		{"same method stays separate across sources", "public_access", 2, 1},
		{"unknown methods stay separate across sources", "review", 2, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := mocks.NewFindingRepository(t)
			repo.On("ListFindingPackages", mock.Anything, mock.Anything, uint32(0)).Return([]*finding.FindingPackage{{PackageName: "tar", FindingCount: 1, WorkType: c.kind, WorkLocation: "google:asset"}, {PackageName: "undici", FindingCount: 1, WorkType: c.kind, WorkLocation: "google:cloudsploit"}}, nil).Once()
			service := &FindingService{repository: repo}
			groups, total, count, err := service.listFindingPackages(context.Background(), &finding.ListFindingRequest{Limit: 1}, 0)
			require.NoError(t, err)
			require.Len(t, groups, c.wantPackages)
			require.Equal(t, uint32(c.wantGroups), total)
			require.Equal(t, uint64(2), count)
		})
	}
}
