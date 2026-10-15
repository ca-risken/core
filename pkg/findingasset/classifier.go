package findingasset

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"regexp"
)

type Asset struct{ Type, Key, Name string }
type object struct {
	NS         string      `json:"ns"`
	Kind       string      `json:"kind"`
	Name       string      `json:"name"`
	Containers []container `json:"containers"`
}
type container struct {
	ImageID string `json:"image_id"`
	URI     string `json:"uri"`
}
type report struct {
	Asset struct {
		Type        string `json:"type"`
		FullName    string `json:"full_name"`
		DisplayName string `json:"display_name"`
	} `json:"asset"`
	Finding struct {
		Category      string `json:"category"`
		ResourceName  string `json:"resource_name"`
		Vulnerability struct {
			CVE struct {
				ID string `json:"id"`
			} `json:"cve"`
		} `json:"vulnerability"`
		Kubernetes struct {
			Objects []object `json:"objects"`
		} `json:"kubernetes"`
		Container container `json:"container"`
	} `json:"finding"`
}

var identifier = regexp.MustCompile(`^(CVE-[0-9]{4}-[0-9]+|GHSA-[a-zA-Z0-9-]+)$`)
var digest = regexp.MustCompile(`^(.+)@sha256:[a-fA-F0-9]{64}$`)

// Unknown keys are resolved from finding_id on read, avoiding a second insert write.
func Classify(source, resource, data string) Asset {
	a := classify(source, resource, data)
	if a.Key != "" {
		sum := sha256.Sum256([]byte(a.Key))
		a.Key = base64.RawURLEncoding.EncodeToString(sum[:])
	}
	return a
}

func classify(source, resource, data string) Asset {
	if source != "google:scc" {
		return Asset{}
	}
	var r report
	if json.Unmarshal([]byte(data), &r) != nil {
		return Asset{}
	}
	switch r.Finding.Category {
	case "OS_VULNERABILITY", "SOFTWARE_VULNERABILITY", "CONTAINER_IMAGE_VULNERABILITY":
	default:
		return Asset{}
	}
	if !identifier.MatchString(r.Finding.Vulnerability.CVE.ID) {
		return Asset{}
	}
	full := r.Asset.FullName
	if full == "" {
		full = r.Finding.ResourceName
	}
	if full == "" {
		full = resource
	}
	workloads := map[string]Asset{}
	images := map[string]Asset{}
	for _, o := range r.Finding.Kubernetes.Objects {
		switch o.Kind {
		case "Deployment", "StatefulSet", "DaemonSet", "ReplicaSet", "Job", "CronJob", "ReplicationController":
			if r.Asset.Type == "google.container.Cluster" && full != "" && o.NS != "" && o.Name != "" {
				key, _ := json.Marshal([]string{full, o.NS, o.Kind, o.Name})
				workloads[string(key)] = Asset{"Workload", string(key), o.NS + "/" + o.Kind + "/" + o.Name}
			}
		}
		for _, c := range o.Containers {
			addImage(images, c)
		}
	}
	if len(workloads) == 1 {
		for _, a := range workloads {
			return a
		}
	}
	addImage(images, r.Finding.Container)
	if digest.MatchString(full) {
		images[full] = Asset{"Image", full, full}
	}
	if len(images) == 1 {
		for _, a := range images {
			return a
		}
	}
	if full != "" {
		typ := "Resource"
		if r.Asset.Type == "google.compute.Instance" {
			typ = "VM"
		}
		name := r.Asset.DisplayName
		if name == "" {
			name = full
		}
		return Asset{typ, full, name}
	}
	return Asset{Type: "Unknown", Name: "Unknown"}
}
func addImage(images map[string]Asset, c container) {
	key := c.ImageID
	if !digest.MatchString(key) {
		key = c.URI
	}
	if digest.MatchString(key) {
		name := c.URI
		if name == "" {
			name = key
		}
		images[key] = Asset{"Image", key, name}
	}
}
