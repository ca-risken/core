package findingasset

import (
	"encoding/json"
	"regexp"
	"strings"
)

var awsAccountID = regexp.MustCompile(`^[0-9]{12}$`)
var awsPartition = regexp.MustCompile(`^aws(?:-[a-z]+)*$`)
var ec2InstanceID = regexp.MustCompile(`^i-[0-9a-f]+$`)

func awsSource(source string) bool {
	switch source {
	case "aws:access-analyzer", "aws:guard-duty", "aws:cloudsploit", "aws:admin-checker":
		return true
	}
	return false
}

type awsARN struct{ partition, service, region, account, resource string }

func parseAWSARN(value string) (awsARN, bool) {
	p := strings.SplitN(value, ":", 6)
	if len(p) != 6 || p[0] != "arn" || !awsPartition.MatchString(p[1]) || p[2] == "" || p[5] == "" || strings.ContainsAny(value, " \t\r\n") {
		return awsARN{}, false
	}
	if p[4] != "" && !awsAccountID.MatchString(p[4]) {
		return awsARN{}, false
	}
	return awsARN{p[1], p[2], p[3], p[4], p[5]}, true
}

func awsAsset(value string) Asset {
	a, ok := parseAWSARN(value)
	if !ok || a.service == "guardduty" || a.service == "access-analyzer" {
		return unknownAWSAsset()
	}
	switch a.service {
	case "ec2", "lambda", "rds", "eks", "kms", "sns", "sqs", "secretsmanager":
		if !awsAccountID.MatchString(a.account) || a.region == "" {
			return unknownAWSAsset()
		}
	case "iam":
		if !awsAccountID.MatchString(a.account) || a.region != "" {
			return unknownAWSAsset()
		}
	}
	if strings.HasSuffix(a.resource, "/") || strings.HasSuffix(a.resource, ":") {
		return unknownAWSAsset()
	}
	typ := "Resource"
	switch {
	case a.service == "ec2" && strings.HasPrefix(a.resource, "instance/"):
		typ = "EC2Instance"
	case a.service == "ec2" && strings.HasPrefix(a.resource, "security-group/"):
		typ = "EC2SecurityGroup"
	case a.service == "iam" && strings.HasPrefix(a.resource, "user/"):
		typ = "IAMUser"
	case a.service == "iam" && strings.HasPrefix(a.resource, "role/"):
		typ = "IAMRole"
	case a.service == "s3" && !strings.ContainsAny(a.resource, "/:"):
		typ = "S3Bucket"
	case a.service == "lambda" && strings.HasPrefix(a.resource, "function:"):
		typ = "LambdaFunction"
	case a.service == "rds" && strings.HasPrefix(a.resource, "db:"):
		typ = "RDSInstance"
	case a.service == "eks" && strings.HasPrefix(a.resource, "cluster/"):
		typ = "EKSCluster"
	case a.service == "kms" && strings.HasPrefix(a.resource, "key/"):
		typ = "KMSKey"
	case a.service == "sns":
		typ = "SNSTopic"
	case a.service == "sqs":
		typ = "SQSQueue"
	case a.service == "secretsmanager" && strings.HasPrefix(a.resource, "secret:"):
		typ = "SecretsManagerSecret"
	}
	return Asset{typ, value, value}
}
func unknownAWSAsset() Asset { return Asset{Type: "Unknown", Name: "Unknown"} }

func classifyAWS(source, data string) Asset {
	var r struct {
		Resource               json.RawMessage
		Arn, AccountId, Region string
		UserARN                string `json:"user_arn"`
		RoleARN                string `json:"role_arn"`
		Data                   struct{ Resource, AccountID, Region string }
	}
	if json.Unmarshal([]byte(data), &r) != nil {
		return unknownAWSAsset()
	}
	switch source {
	case "aws:access-analyzer":
		var target string
		if json.Unmarshal(r.Resource, &target) != nil {
			return unknownAWSAsset()
		}
		return awsAsset(target)
	case "aws:admin-checker":
		if r.UserARN != "" && r.RoleARN != "" {
			return unknownAWSAsset()
		}
		target := r.UserARN
		if target == "" {
			target = r.RoleARN
		}
		a := awsAsset(target)
		if (r.UserARN != "" && a.Type != "IAMUser") || (r.RoleARN != "" && a.Type != "IAMRole") || (a.Type != "IAMUser" && a.Type != "IAMRole") {
			return unknownAWSAsset()
		}
		return a
	case "aws:cloudsploit":
		return awsAsset(r.Data.Resource)
	case "aws:guard-duty":
		var target struct {
			ResourceType      string
			InstanceDetails   *struct{ InstanceId string }
			S3BucketDetails   []struct{ Arn string }
			EksClusterDetails *struct{ Arn string }
		}
		if json.Unmarshal(r.Resource, &target) != nil {
			return unknownAWSAsset()
		}
		switch target.ResourceType {
		case "Instance":
			findingARN, ok := parseAWSARN(r.Arn)
			if !ok || findingARN.service != "guardduty" || findingARN.account != r.AccountId || findingARN.region != r.Region || !awsAccountID.MatchString(r.AccountId) || r.Region == "" || target.InstanceDetails == nil || !ec2InstanceID.MatchString(target.InstanceDetails.InstanceId) {
				return unknownAWSAsset()
			}
			return awsAsset("arn:" + findingARN.partition + ":ec2:" + r.Region + ":" + r.AccountId + ":instance/" + target.InstanceDetails.InstanceId)
		case "S3Bucket":
			candidates := map[string]Asset{}
			for _, b := range target.S3BucketDetails {
				a := awsAsset(b.Arn)
				if a.Type != "S3Bucket" {
					return unknownAWSAsset()
				}
				candidates[a.Key] = a
			}
			if len(candidates) == 1 {
				for _, a := range candidates {
					return a
				}
			}
		case "EKSCluster":
			if target.EksClusterDetails != nil {
				a := awsAsset(target.EksClusterDetails.Arn)
				if a.Type == "EKSCluster" {
					return a
				}
			}
		case "AccessKey":
			// AccessKeyDetails contains names and IDs, not a canonical IAM ARN.
			// Do not infer an IAM path or turn a role session into a role ARN.
		}
	}
	return unknownAWSAsset()
}
