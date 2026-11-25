package findingremediation

import (
	"fmt"
	"testing"
)

func TestSCCWorkMethods(t *testing.T) {
	cases := []struct{ name, typ, nameJSON, files, offending, want string }{
		{"OS before app path", "OS", `"openssl"`, `[{"path":"/app/lib/a"}]`, `[]`, "os_update"},
		{"OS without name", "OS", `""`, `[]`, `[]`, "os_update"},
		{"absolute app", "NPM", `"next"`, `[{"path":"/app/node_modules/a/package.json"}]`, `[]`, "app_dependency_update"},
		{"relative app", "PYPI", `"flask"`, `[{"path":"app/venv/lib/a"}]`, `[]`, "app_dependency_update"},
		{"disk relative app", "GO", `"stdlib"`, `[{"disk_path":{"relative_path":"app/bin/server"}}]`, `[]`, "app_dependency_update"},
		{"offending package files", "GO", `"stdlib"`, `[]`, `[{"path":"/app/bin/server"}]`, "app_dependency_update"},
		{"merge both arrays", "NPM", `"a"`, `[{"path":"/app/a"}]`, `[{"path":"app/b"}]`, "app_dependency_update"},
		{"direct path before disk", "GO", `"a"`, `[{"path":"/usr/bin/tool","disk_path":{"relative_path":"app/tool"}}]`, `[]`, "library_update"},
		{"disk fallback", "GO", `"a"`, `[{"path":"","disk_path":{"relative_path":"app/tool"}}]`, `[]`, "app_dependency_update"},
		{"VM binary", "GO", `"stdlib"`, `[{"disk_path":{"relative_path":"home/kubernetes/bin/kubelet"}}]`, `[]`, "library_update"},
		{"workload relative binary", "GO", `"stdlib"`, `[{"path":"usr/local/bin/gosu"}]`, `[]`, "library_update"},
		{"mixed app and system", "GO", `"a"`, `[{"path":"/app/a"},{"path":"/usr/bin/a"}]`, `[]`, "library_update"},
		{"multiple app files", "NPM", `"a"`, `[{"path":"/app/a"},{"path":"/app/sub/b"}]`, `[]`, "app_dependency_update"},
		{"null root plus app", "NPM", `"a"`, `null`, `[{"path":"/app/a"}]`, "library_update"},
		{"empty arrays", "NPM", `"a"`, `[]`, `[]`, "library_update"},
		{"missing path", "NPM", `"a"`, `[{}]`, `[]`, "library_update"},
		{"missing package", "NPM", `""`, `[{"path":"/app/a"}]`, `[]`, "review"},
		{"parent traversal", "GO", `"a"`, `[{"path":"/app/../bin/a"}]`, `[]`, "library_update"},
		{"current directory", "GO", `"a"`, `[{"path":"/app/./a"}]`, `[]`, "library_update"},
		{"trailing parent", "GO", `"a"`, `[{"path":"/app/a/.."}]`, `[]`, "library_update"},
		{"directory", "GO", `"a"`, `[{"path":"/app/a/"}]`, `[]`, "library_update"},
		{"app root is not file", "GO", `"a"`, `[{"path":"/app"}]`, `[]`, "library_update"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data := fmt.Sprintf(`{"finding":{"files":%s,"vulnerability":{"offending_package":{"package_type":%q,"package_name":%s,"files":%s}}}}`, c.files, c.typ, c.nameJSON, c.offending)
			if got := Classify("google:scc", "per-finding-id", "SOFTWARE_VULNERABILITY", data); got != c.want {
				t.Fatalf("got %s want %s", got, c.want)
			}
		})
	}
	for _, c := range []struct{ name, data string }{{"invalid JSON", "{"}, {"empty object", "{}"}, {"no package", `{"finding":{"files":[{"path":"app/a"}]}}`}} {
		t.Run(c.name, func(t *testing.T) {
			if got := Classify("google:scc", "id", "SOFTWARE_VULNERABILITY", c.data); got != "review" {
				t.Fatalf("got %s", got)
			}
		})
	}
}
