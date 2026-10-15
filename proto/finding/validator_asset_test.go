package finding

import (
	"strings"
	"testing"
)

func TestAssetKeyByteLimit(t *testing.T) {
	cases := []struct {
		name    string
		key     string
		invalid bool
	}{
		{"at limit", strings.Repeat("a", 8192), false},
		{"over limit", strings.Repeat("a", 8193), true},
		{"multibyte at limit", strings.Repeat("あ", 2730) + "ab", false},
		{"multibyte over limit", strings.Repeat("あ", 2731), true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			requests := []interface{ Validate() error }{
				&ListFindingRequest{ProjectId: 1, AssetType: "VM", AssetKey: c.key},
				&ListFindingForOrgRequest{OrganizationId: 1, AssetType: "VM", AssetKey: c.key},
			}
			for _, req := range requests {
				if err := req.Validate(); (err != nil) != c.invalid {
					t.Errorf("%T Validate() = %v, invalid = %v", req, err, c.invalid)
				}
			}
		})
	}
}

func TestOrgAssetPaginationLimit(t *testing.T) {
	cases := []struct {
		name    string
		view    bool
		key     string
		limit   int32
		offset  int32
		invalid bool
	}{
		{"view at limit", true, "", 200, 0, false},
		{"view over limit", true, "", 201, 0, true},
		{"view default", true, "", 0, 0, false},
		{"view negative limit", true, "", -1, 0, true},
		{"view negative offset", true, "", 20, -1, true},
		{"detail at limit", false, "vm", 200, 0, false},
		{"detail over limit", false, "vm", 201, 0, true},
		{"detail negative limit", false, "vm", -1, 0, true},
		{"detail negative offset", false, "vm", 20, -1, true},
		{"ordinary list unchanged", false, "", 201, 0, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := &ListFindingForOrgRequest{OrganizationId: 1, AssetView: c.view, AssetKey: c.key, Limit: c.limit, Offset: c.offset}
			if c.key != "" {
				req.AssetType = "VM"
			}
			if err := req.Validate(); (err != nil) != c.invalid {
				t.Errorf("Validate() = %v, invalid = %v", err, c.invalid)
			}
		})
	}
}

func TestOrgAssetProjectFilterLimit(t *testing.T) {
	cases := []struct {
		name    string
		limit   int32
		invalid bool
	}{
		{"at limit", 200, false},
		{"over limit", 201, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := &ListFindingForOrgRequest{OrganizationId: 1, AssetProjectId: 1, Limit: c.limit}
			if err := req.Validate(); (err != nil) != c.invalid {
				t.Errorf("Validate() = %v, invalid = %v", err, c.invalid)
			}
		})
	}
}
