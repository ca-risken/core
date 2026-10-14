package findingasset

import (
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
		{name: "image digest", source: "google:scc", category: "CONTAINER_IMAGE_VULNERABILITY", id: "CVE-2023-1234", container: container{ImageID: "registry.example/app@sha256:" + strings.Repeat("a", 64)}, wantType: "Image", wantKey: "registry.example/app@sha256:" + strings.Repeat("a", 64)},
		{name: "image digest in URI", source: "google:scc", typ: "google.container.Cluster", resource: "cluster", category: "CONTAINER_IMAGE_VULNERABILITY", id: "CVE-2023-1234", container: container{URI: "registry.example/app@sha256:" + strings.Repeat("b", 64)}, wantType: "Image", wantKey: "registry.example/app@sha256:" + strings.Repeat("b", 64)},
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
			if got.Type != c.wantType || got.Key != c.wantKey {
				t.Fatalf("got %+v want %s/%s", got, c.wantType, c.wantKey)
			}
		})
	}
}
