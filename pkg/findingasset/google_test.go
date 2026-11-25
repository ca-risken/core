package findingasset

import "testing"

func TestGoogleAssets(t *testing.T) {
	cases := []struct{ name, source, data, typ string }{
		{"asset bucket", "google:asset", `{"asset":{"name":"//storage.googleapis.com/demo"}}`, "Resource"},
		{"cloudsploit bucket", "google:cloudsploit", `{"resource":"//storage.googleapis.com/demo"}`, "Resource"},
		{"instance", "google:cloudsploit", `{"resource":"//compute.googleapis.com/projects/p/zones/z/instances/v"}`, "VM"},
		{"ambiguous name", "google:cloudsploit", `{"resource":"demo"}`, "Unknown"},
		{"account check", "google:cloudsploit", `{"resource":"N/A"}`, "Unknown"},
		{"malformed", "google:asset", `{`, "Unknown"},
		{"excluded", "google:portscan", `{"resource":"//storage.googleapis.com/demo"}`, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := Classify(c.source, "", c.data)
			if a.Type != c.typ {
				t.Fatalf("got %+v", a)
			}
		})
	}
	a := Classify("google:asset", "", `{"asset":{"name":"//storage.googleapis.com/demo"}}`)
	b := Classify("google:cloudsploit", "", `{"resource":"//storage.googleapis.com/demo"}`)
	if a.Key != b.Key || a.Type != b.Type {
		t.Fatal("same full name must group together")
	}
}
