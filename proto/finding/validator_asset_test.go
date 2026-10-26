package finding

import (
	"strings"
	"testing"
)

func TestAssetKeyFormat(t *testing.T) {
	cases := []struct {
		name, typ, key string
		invalid        bool
	}{
		{"hash", "VM", "W86Y9z8-0Mg38nKe2VCbOOpmoVbbf2UzVstv43s2boU", false},
		{"short", "VM", strings.Repeat("a", 42), true},
		{"long", "VM", strings.Repeat("a", 44), true},
		{"unicode", "VM", strings.Repeat("あ", 43), true},
		{"padding", "VM", "W86Y9z8-0Mg38nKe2VCbOOpmoVbbf2UzVstv43s2boU=", true},
		{"invalid alphabet", "VM", strings.Repeat("+", 43), true},
		{"noncanonical", "VM", strings.Repeat("A", 42) + "B", true},
		{"unknown", "Unknown", "99", false},
		{"unknown zero", "Unknown", "0", true},
		{"unknown invalid", "Unknown", "99abc", true},
		{"unknown overflow", "Unknown", "18446744073709551616", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			requests := []interface{ Validate() error }{&ListFindingRequest{ProjectId: 1, AssetType: c.typ, AssetKey: c.key}, &ListFindingForOrgRequest{OrganizationId: 1, AssetType: c.typ, AssetKey: c.key}}
			for _, req := range requests {
				if err := req.Validate(); (err != nil) != c.invalid {
					t.Errorf("%T: %v invalid=%v", req, err, c.invalid)
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
		{"detail at limit", false, "W86Y9z8-0Mg38nKe2VCbOOpmoVbbf2UzVstv43s2boU", 200, 0, false},
		{"detail over limit", false, "W86Y9z8-0Mg38nKe2VCbOOpmoVbbf2UzVstv43s2boU", 201, 0, true},
		{"detail negative limit", false, "W86Y9z8-0Mg38nKe2VCbOOpmoVbbf2UzVstv43s2boU", -1, 0, true},
		{"detail negative offset", false, "W86Y9z8-0Mg38nKe2VCbOOpmoVbbf2UzVstv43s2boU", 20, -1, true},
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

func TestAssetSortSupportedColumns(t *testing.T) {
	cases := []struct {
		name, sort string
		invalid    bool
	}{
		{"score", "score", false},
		{"finding count", "finding_count", false},
		{"risk sorting disabled", "risk_information", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			requests := []interface{ Validate() error }{
				&ListFindingRequest{ProjectId: 1, AssetView: true, Sort: tc.sort},
				&ListFindingForOrgRequest{OrganizationId: 1, AssetView: true, Sort: tc.sort},
			}
			for _, req := range requests {
				if err := req.Validate(); (err != nil) != tc.invalid {
					t.Fatalf("%T: error=%v, invalid=%v", req, err, tc.invalid)
				}
			}
		})
	}
}
