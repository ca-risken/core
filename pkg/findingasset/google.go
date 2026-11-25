package findingasset

import (
	"encoding/json"
	"regexp"
	"strings"
)

var googleFullName = regexp.MustCompile(`^//[a-z0-9.-]+\.googleapis\.com/[^[:space:]]+$`)

func classifyGoogle(source, data string) Asset {
	var record struct {
		Asset struct {
			Name string `json:"name"`
		} `json:"asset"`
		Resource string `json:"resource"`
	}
	if json.Unmarshal([]byte(data), &record) != nil {
		return Asset{Type: "Unknown", Name: "Unknown"}
	}
	full := record.Resource
	if source == "google:asset" {
		full = record.Asset.Name
	}
	return googleResource(full)
}

func googleResource(full string) Asset {
	if !googleFullName.MatchString(full) {
		return Asset{Type: "Unknown", Name: "Unknown"}
	}
	typ := "Resource"
	switch {
	case strings.HasPrefix(full, "//storage.googleapis.com/"):
		typ = "Resource"
	case strings.HasPrefix(full, "//compute.googleapis.com/") && strings.Contains(full, "/instances/"):
		typ = "VM"
	}
	return Asset{Type: typ, Key: full, Name: full}
}
