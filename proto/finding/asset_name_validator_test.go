package finding

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestAssetNameValidation(t *testing.T) {
	cases := []struct {
		name    string
		size    int
		invalid bool
	}{{"empty", 0, false}, {"boundary", 512, false}, {"too long", 513, true}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			name := strings.Repeat("x", c.size)
			errs := []error{(&ListFindingRequest{ProjectId: 1, AssetView: true, AssetName: name, Limit: 20}).Validate(), (&ListFindingForOrgRequest{OrganizationId: 7, AssetView: true, AssetName: name, Limit: 20}).Validate()}
			for _, err := range errs {
				if c.invalid {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			}
		})
	}
	cases2 := []struct {
		name          string
		limit, offset int32
		invalid       bool
	}{{"valid", 200, 0, false}, {"over limit", 201, 0, true}, {"negative offset", 20, -1, true}}
	for _, c := range cases2 {
		t.Run(c.name, func(t *testing.T) {
			err := (&ListFindingForOrgRequest{OrganizationId: 7, AssetName: "stage", Limit: c.limit, Offset: c.offset}).Validate()
			if c.invalid {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
