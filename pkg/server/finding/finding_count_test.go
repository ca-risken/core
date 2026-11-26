package finding

import (
	"context"
	"errors"
	"testing"

	"github.com/ca-risken/core/pkg/db/mocks"
	"github.com/ca-risken/core/proto/finding"
	"github.com/stretchr/testify/require"
)

func TestGetFindingCount(t *testing.T) {
	cases := []struct {
		name    string
		req     *finding.GetFindingCountRequest
		toScore float32
		total   int64
		dbErr   error
		invalid bool
	}{
		{name: "default upper score", req: &finding.GetFindingCountRequest{ProjectId: 1, FromScore: 0.8, DataSource: []string{"aws:"}, Status: finding.FindingStatus_FINDING_ACTIVE}, toScore: 1, total: 345},
		{name: "all filters", req: &finding.GetFindingCountRequest{ProjectId: 2, FromScore: 0.4, ToScore: 0.7, DataSource: []string{"google:"}, ResourceName: []string{"resource"}, Tag: []string{"tag"}, FindingId: 9, AlertId: 3, Status: finding.FindingStatus_FINDING_PENDING}, toScore: 0.7, total: 1},
		{name: "zero", req: &finding.GetFindingCountRequest{ProjectId: 1}, toScore: 1},
		{name: "database error", req: &finding.GetFindingCountRequest{ProjectId: 1}, toScore: 1, dbErr: errors.New("db failed")},
		{name: "missing project", req: &finding.GetFindingCountRequest{}, invalid: true},
		{name: "invalid score", req: &finding.GetFindingCountRequest{ProjectId: 1, FromScore: 1.1}, invalid: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := mocks.NewFindingRepository(t)
			svc := FindingService{repository: repo}
			if !c.invalid {
				repo.On("ListFindingCount", context.Background(), c.req.ProjectId, uint32(0), c.req.AlertId,
					c.req.FromScore, c.toScore, c.req.FindingId, c.req.DataSource, c.req.ResourceName, c.req.Tag, c.req.Status).Return(c.total, c.dbErr).Once()
			}
			got, err := svc.GetFindingCount(context.Background(), c.req)
			if c.invalid || c.dbErr != nil {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.Equal(t, uint32(c.total), got.Total)
			repo.AssertNotCalled(t, "ListFinding")
		})
	}
}
