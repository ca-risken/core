package findingremediation

import (
	"encoding/json"
	"strings"
)

// sccMethod preserves the previous SCC work-category rules, without retaining
// package, platform, directory, or version as aggregation keys.
func sccMethod(data string) string {
	var document struct {
		Finding struct {
			Files         json.RawMessage `json:"files"`
			Vulnerability struct {
				Offending struct {
					Name  string          `json:"package_name"`
					Type  string          `json:"package_type"`
					Files json.RawMessage `json:"files"`
				} `json:"offending_package"`
			} `json:"vulnerability"`
		} `json:"finding"`
	}
	if json.Unmarshal([]byte(data), &document) != nil {
		return "review"
	}
	p := document.Finding.Vulnerability.Offending
	if p.Type == "OS" {
		return "os_update"
	}
	if p.Name == "" {
		return "review"
	}
	files := append(sccFiles(document.Finding.Files), sccFiles(p.Files)...)
	if len(files) == 0 {
		return "library_update"
	}
	for _, f := range files {
		path := f.Path
		if path == "" {
			path = f.DiskPath.RelativePath
		}
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		// The old SQL rejects incomplete paths and directory traversal before
		// applying /app. Invalid or mixed scopes fall back to library update.
		if !strings.HasPrefix(path, "/app/") || strings.HasSuffix(path, "/") {
			return "library_update"
		}
		for _, part := range strings.Split(path, "/") {
			if part == "." || part == ".." {
				return "library_update"
			}
		}
	}
	return "app_dependency_update"
}

type sccFile struct {
	Path     string `json:"path"`
	DiskPath struct {
		RelativePath string `json:"relative_path"`
	} `json:"disk_path"`
}

func sccFiles(raw json.RawMessage) []sccFile {
	if len(raw) == 0 {
		return nil
	}
	// An explicit null contributed an unknown path in the previous JSON_TABLE.
	var files []sccFile
	if string(raw) == "null" || json.Unmarshal(raw, &files) != nil {
		return []sccFile{{}}
	}
	return files
}
