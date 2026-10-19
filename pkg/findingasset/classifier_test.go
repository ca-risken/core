package findingasset

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

func TestClassify(t *testing.T) {
	cases := []struct {
		name, source, typ, resource, category, id string
		objects                                   []object
		container                                 container
		wantType, wantKey                         string
	}{
		{name: "other source", source: "aws:inspector", category: "OS_VULNERABILITY", id: "CVE-2023-1234"},
		{name: "unsupported category", source: "google:scc", category: "PUBLIC_BUCKET", id: "CVE-2023-1234"},
		{name: "missing identifier", source: "google:scc", category: "OS_VULNERABILITY"},
		{name: "workload", source: "google:scc", typ: "google.container.Cluster", resource: "//container.googleapis.com/projects/test/zones/z/clusters/c", category: "OS_VULNERABILITY", id: "CVE-2023-1234", objects: []object{{NS: "default", Kind: "Deployment", Name: "app"}}, wantType: "Workload", wantKey: `["//container.googleapis.com/projects/test/zones/z/clusters/c","default","Deployment","app"]`},
		{name: "duplicate workload", source: "google:scc", typ: "google.container.Cluster", resource: "cluster", category: "SOFTWARE_VULNERABILITY", id: "GHSA-abcd-1234-efgh", objects: []object{{NS: "n", Kind: "StatefulSet", Name: "app"}, {NS: "n", Kind: "StatefulSet", Name: "app"}}, wantType: "Workload", wantKey: `["cluster","n","StatefulSet","app"]`},
		{name: "ambiguous workloads fall back", source: "google:scc", typ: "google.container.Cluster", resource: "cluster", category: "OS_VULNERABILITY", id: "CVE-2023-1234", objects: []object{{NS: "n", Kind: "Deployment", Name: "a"}, {NS: "n", Kind: "Deployment", Name: "b"}}, wantType: "Resource", wantKey: "cluster"},
		{name: "pod not guessed", source: "google:scc", typ: "google.container.Cluster", resource: "cluster", category: "OS_VULNERABILITY", id: "CVE-2023-1234", objects: []object{{NS: "n", Kind: "Pod", Name: "app-123"}}, wantType: "Resource", wantKey: "cluster"},
		{name: "missing namespace", source: "google:scc", typ: "google.container.Cluster", resource: "cluster", category: "OS_VULNERABILITY", id: "CVE-2023-1234", objects: []object{{Kind: "Deployment", Name: "app"}}, wantType: "Resource", wantKey: "cluster"},
		{name: "image digest", source: "google:scc", category: "CONTAINER_IMAGE_VULNERABILITY", id: "CVE-2023-1234", container: container{ImageID: "registry.example/app@sha256:" + strings.Repeat("a", 64)}, wantType: "Image", wantKey: "registry.example/app"},
		{name: "image digest in URI", source: "google:scc", typ: "google.container.Cluster", resource: "cluster", category: "CONTAINER_IMAGE_VULNERABILITY", id: "CVE-2023-1234", container: container{URI: "registry.example/app@sha256:" + strings.Repeat("b", 64)}, wantType: "Image", wantKey: "registry.example/app"},
		{name: "image resource", source: "google:scc", resource: "registry.example/project/app@sha256:" + strings.Repeat("a", 64), category: "CONTAINER_IMAGE_VULNERABILITY", id: "GHSA-vcvr-r3jv-pc5j", wantType: "Image", wantKey: "registry.example/project/app"},
		{name: "image tag and registry port", source: "google:scc", category: "CONTAINER_IMAGE_VULNERABILITY", id: "CVE-2023-1234", container: container{ImageID: "registry.example:5000/project/app:v1@sha256:" + strings.Repeat("a", 64), URI: "display-name:v1"}, wantType: "Image", wantKey: "registry.example:5000/project/app"},
		{name: "same repository candidates", source: "google:scc", resource: "registry.example/app@sha256:" + strings.Repeat("a", 64), category: "CONTAINER_IMAGE_VULNERABILITY", id: "CVE-2023-1234", container: container{ImageID: "registry.example/app@sha256:" + strings.Repeat("b", 64)}, wantType: "Image", wantKey: "registry.example/app"},
		{name: "different repositories remain ambiguous", source: "google:scc", resource: "cluster", category: "CONTAINER_IMAGE_VULNERABILITY", id: "CVE-2023-1234", objects: []object{{Kind: "Pod", Containers: []container{{ImageID: "registry.example/a/app@sha256:" + strings.Repeat("a", 64)}, {ImageID: "registry.example/b/app@sha256:" + strings.Repeat("b", 64)}}}}, wantType: "Resource", wantKey: "cluster"},
		{name: "vm", source: "google:scc", typ: "google.compute.Instance", resource: "//compute.googleapis.com/projects/p/zones/z/instances/v", category: "OS_VULNERABILITY", id: "CVE-2023-1234", wantType: "VM", wantKey: "//compute.googleapis.com/projects/p/zones/z/instances/v"},
		{name: "unknown", source: "google:scc", category: "OS_VULNERABILITY", id: "CVE-2023-1234", wantType: "Unknown"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var r report
			r.Asset.Type = c.typ
			r.Asset.FullName = c.resource
			r.Finding.Category = c.category
			r.Finding.Vulnerability.CVE.ID = c.id
			r.Finding.Kubernetes.Objects = c.objects
			r.Finding.Container = c.container
			b, _ := json.Marshal(r)
			got := Classify(c.source, c.resource, string(b))
			wantKey := c.wantKey
			if wantKey != "" {
				sum := sha256.Sum256([]byte(wantKey))
				wantKey = base64.RawURLEncoding.EncodeToString(sum[:])
			}
			if got.Type != c.wantType || got.Key != wantKey || (c.wantType == "Image" && got.Name != c.wantKey) {
				t.Fatalf("got %+v want %s/%s", got, c.wantType, wantKey)
			}
		})
	}
}

func TestHashIdentity(t *testing.T) {
	cases := []struct{ name, identity, want string }{
		{"short", "vm", "W86Y9z8-0Mg38nKe2VCbOOpmoVbbf2UzVstv43s2boU"},
		{"long", strings.Repeat("x", 9000), "55fir28FwkzdBkeTzmCp0BMC2aOw5d_1uwwEa4fy9mg"},
		{"workload", "[\"cluster\",\"ns\",\"Deployment\",\"app\"]", "Lxxp4l0xeIUsX2i4WZFkXUeuVWGiKO9s-7WA9FDl6U0"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data, _ := json.Marshal(map[string]any{"asset": map[string]any{"type": "google.compute.Instance", "full_name": c.identity}, "finding": map[string]any{"category": "OS_VULNERABILITY", "vulnerability": map[string]any{"cve": map[string]string{"id": "CVE-2023-1234"}}}})
			got := Classify("google:scc", "", string(data))
			if got.Key != c.want || len(got.Key) != 43 {
				t.Fatalf("got %s want %s", got.Key, c.want)
			}
		})
	}
}
