// Package findingremediation classifies recommendations at ingestion.
package findingremediation

import "strings"

// Classify uses source-owned types; SCC also reuses package/path work rules.
// It does not group by prose similarity or individual recommendation IDs.
// Unknown mappings remain review; classification does not imply resolution.
func Classify(source, typ, risk, findingData string) string {
	switch source {
	case "aws:access-analyzer":
		switch typ {
		case "s3", "iam", "kms", "lambda", "sqs", "secretsmanager":
			return "access_policy"
		}
	case "aws:guard-duty":
		if strings.Contains(typ, ":") && strings.Contains(typ, "/") {
			return "threat"
		}
	case "aws:admin-checker":
		switch typ {
		case "admin":
			return "admin_permissions"
		case "access-report":
			return "unused_permissions"
		}
	case "google:asset":
		switch typ {
		case "storage.googleapis.com/Bucket":
			return "public_access"
		case "iam.googleapis.com/ServiceAccount":
			return "admin_permissions"
		}
	case "google:scc":
		if method := sccMethod(findingData); method != "review" {
			return method
		}
		category, _, _ := strings.Cut(risk, "\n")
		switch strings.TrimSpace(category) {
		case "PUBLIC_BUCKET_ACL", "PUBLIC_BUCKET_IAM":
			return "public_access"
		}
	case "google:cloudsploit":
		switch typ {
		case "Storage/bucketAllUsersPolicy", "BigQuery/datasetAllUsersPolicy":
			return "public_access"
		case "IAM/serviceAccountAdmin":
			return "admin_permissions"
		case "IAM/serviceAccountKeyRotation":
			return "key_rotation"
		}
	case "aws:cloudsploit":
		switch typ {
		case "S3/bucketAllUsersPolicy", "S3/bucketAllUsersAcl":
			return "public_access"
		}
	}
	return "review"
}
