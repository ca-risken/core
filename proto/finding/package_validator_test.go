package finding

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestValidatePackageRequest(t *testing.T) {
	cases := []struct {
		name            string
		view            bool
		key, typ, asset string
		project         uint32
		invalid         bool
	}{
		{name: "default"},
		{name: "view", view: true, typ: "Image", asset: "key", project: 1},
		{name: "details", key: `["NPM", "next", "", "app:/app"]`, typ: "Image", asset: "key", project: 1},
		{name: "no scope", view: true, invalid: true},
		{name: "malformed", key: `{}`, typ: "Image", asset: "key", project: 1, invalid: true},
		{name: "wrong parts", key: `["NPM"]`, typ: "Image", asset: "key", project: 1, invalid: true},
		{name: "too long", key: strings.Repeat("x", 8193), typ: "Image", asset: "key", project: 1, invalid: true},
		{name: "conflicting", view: true, key: `["NPM", "next", "", "app:/app"]`, typ: "Image", asset: "key", project: 1, invalid: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validatePackageRequest(c.view, c.key, false, c.typ, c.asset, c.project)
			if c.invalid {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
