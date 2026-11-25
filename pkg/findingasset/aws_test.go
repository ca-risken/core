package findingasset

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestAWSClassification(t *testing.T) {
	const user = "arn:aws:iam::123456789012:user/team/alice"
	const bucket = "arn:aws:s3:::example-bucket"
	const instance = "arn:aws:ec2:ap-northeast-1:123456789012:instance/i-123abc"
	cases := []struct{ name, source, data, typ, arn string }{
		{"access user", "aws:access-analyzer", `{"Resource":"` + user + `"}`, "IAMUser", user},
		{"access bucket", "aws:access-analyzer", `{"Resource":"` + bucket + `","Principal":{"AWS":"*"}}`, "S3Bucket", bucket},
		{"admin user", "aws:admin-checker", `{"user_arn":"` + user + `"}`, "IAMUser", user},
		{"admin role", "aws:admin-checker", `{"role_arn":"arn:aws:iam::123456789012:role/team/worker"}`, "IAMRole", "arn:aws:iam::123456789012:role/team/worker"},
		{"admin ambiguous", "aws:admin-checker", `{"user_arn":"` + user + `","role_arn":"other"}`, "Unknown", ""},
		{"cloud instance", "aws:cloudsploit", `{"data":{"Resource":"` + instance + `"}}`, "EC2Instance", instance},
		{"cloud setting", "aws:cloudsploit", `{"data":{"Resource":"arn:aws:ec2:ap-northeast-1:123456789012:security-group/sg-123"}}`, "EC2SecurityGroup", "arn:aws:ec2:ap-northeast-1:123456789012:security-group/sg-123"},
		{"cloud global resource", "aws:cloudsploit", `{"data":{"Resource":"arn:aws:cloudfront::123456789012:distribution/EXAMPLE"}}`, "Resource", "arn:aws:cloudfront::123456789012:distribution/EXAMPLE"},
		{"missing regional scope", "aws:cloudsploit", `{"data":{"Resource":"arn:aws:ec2:::instance/i-123abc"}}`, "Unknown", ""},
		{"empty resource identity", "aws:cloudsploit", `{"data":{"Resource":"arn:aws:ec2:ap-northeast-1:123456789012:instance/"}}`, "Unknown", ""},
		{"admin wrong kind", "aws:admin-checker", `{"user_arn":"arn:aws:iam::123456789012:role/team/worker"}`, "Unknown", ""},
		{"invalid partition", "aws:access-analyzer", `{"Resource":"arn:awsINVALID:iam::123456789012:user/alice"}`, "Unknown", ""},

		{"cloud unknown", "aws:cloudsploit", `{"data":{"Resource":"N/A"}}`, "Unknown", ""},
		{"cloud no fallback", "aws:cloudsploit", `{}`, "Unknown", ""},
		{"guard instance", "aws:guard-duty", `{"Arn":"arn:aws:guardduty:ap-northeast-1:123456789012:detector/d/finding/f","AccountId":"123456789012","Region":"ap-northeast-1","Resource":{"ResourceType":"Instance","InstanceDetails":{"InstanceId":"i-123abc"}}}`, "EC2Instance", instance},
		{"guard bucket", "aws:guard-duty", `{"Resource":{"ResourceType":"S3Bucket","S3BucketDetails":[{"Arn":"` + bucket + `"},{"Arn":"` + bucket + `"}]}}`, "S3Bucket", bucket},
		{"guard multiple", "aws:guard-duty", `{"Resource":{"ResourceType":"S3Bucket","S3BucketDetails":[{"Arn":"` + bucket + `"},{"Arn":"arn:aws:s3:::other"}]}}`, "Unknown", ""},
		{"guard incomplete instance", "aws:guard-duty", `{"Resource":{"ResourceType":"Instance","InstanceDetails":{"InstanceId":"i-123abc"}}}`, "Unknown", ""},
		{"guard username not arn", "aws:guard-duty", `{"Resource":{"ResourceType":"AccessKey","AccessKeyDetails":{"UserName":"alice","UserType":"IAMUser"}}}`, "Unknown", ""},
		{"invalid json", "aws:access-analyzer", `{`, "Unknown", ""},
		{"null", "aws:access-analyzer", `null`, "Unknown", ""},
		{"not target arn", "aws:access-analyzer", `{"Resource":"arn:aws:guardduty:ap-northeast-1:123456789012:detector/d/finding/f"}`, "Unknown", ""},
		{"portscan excluded", "aws:portscan", `{"data":{"Resource":"` + instance + `"}}`, "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := Classify(c.source, "arn:aws:s3:::must-not-fallback", c.data)
			expected := Asset{Type: c.typ}
			if c.arn != "" {
				sum := sha256.Sum256([]byte(c.arn))
				expected.Key = base64.RawURLEncoding.EncodeToString(sum[:])
				if a.Key == "" || len(a.Key) != 43 {
					t.Fatal(a)
				}
			}
			if a.Type != c.typ || a.Key != expected.Key {
				t.Fatalf("got %+v want type=%s key=%s", a, c.typ, expected.Key)
			}
			if c.arn != "" && a.Name != c.arn {
				t.Fatal(a)
			}
		})
	}
}
func quoteARN(s string) string { b, _ := json.Marshal(s); return string(b) }

func TestAWSIdentityScope(t *testing.T) {
	cases := []struct {
		name, left, right string
		same              bool
	}{
		{"same ARN", "arn:aws:iam::123456789012:user/team/alice", "arn:aws:iam::123456789012:user/team/alice", true},
		{"path", "arn:aws:iam::123456789012:user/team/alice", "arn:aws:iam::123456789012:user/alice", false},
		{"case", "arn:aws:iam::123456789012:user/Alice", "arn:aws:iam::123456789012:user/alice", false},
		{"account", "arn:aws:iam::123456789012:user/alice", "arn:aws:iam::223456789012:user/alice", false},
		{"partition", "arn:aws:iam::123456789012:user/alice", "arn:aws-us-gov:iam::123456789012:user/alice", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := Classify("aws:access-analyzer", "", `{"Resource":`+quoteARN(c.left)+`}`)
			b := Classify("aws:admin-checker", "", `{"user_arn":`+quoteARN(c.right)+`}`)
			if (a.Key == b.Key) != c.same {
				t.Fatalf("%+v %+v", a, b)
			}
		})
	}
}
