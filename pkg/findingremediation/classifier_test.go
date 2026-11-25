package findingremediation

import "testing"

func TestClassify(t *testing.T) {
	cases := []struct{ name, source, typ, risk, want string }{
		{"SCC vulnerability", "google:scc", "finding-one", "SOFTWARE_VULNERABILITY\nfirst", "review"},
		{"SCC other finding same method", "google:scc", "finding-two", "OS_VULNERABILITY\nsecond", "review"},
		{"SCC explicit public access", "google:scc", "id", "PUBLIC_BUCKET_ACL", "public_access"},
		{"no prose guessing", "google:scc", "id", "UNKNOWN\nSOFTWARE_VULNERABILITY", "review"},
		{"Google bucket", "google:asset", "storage.googleapis.com/Bucket", "", "public_access"},
		{"Google service account", "google:asset", "iam.googleapis.com/ServiceAccount", "", "admin_permissions"},
		{"Google cloudsploit", "google:cloudsploit", "Storage/bucketAllUsersPolicy", "", "public_access"},
		{"unmapped plugin", "google:cloudsploit", "Storage/newPlugin", "", "review"},
		{"admin", "aws:admin-checker", "admin", "", "admin_permissions"},
		{"unused", "aws:admin-checker", "access-report", "", "unused_permissions"},
		{"access analyzer", "aws:access-analyzer", "s3", "", "access_policy"},
		{"guardduty", "aws:guard-duty", "UnauthorizedAccess:EC2/SSHBruteForce", "", "threat"},
		{"portscan excluded", "google:portscan", "storage.googleapis.com/Bucket", "", "review"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Classify(c.source, c.typ, c.risk, ""); got != c.want {
				t.Fatalf("got %s want %s", got, c.want)
			}
		})
	}
}
