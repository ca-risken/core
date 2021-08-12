####################################################
## grpcurl example
####################################################
.PHONY: all
all: help

.PHONY: help
help:
	@echo "Usage: make <sub-command>"
	@echo "\n---------------- sub-command list ----------------"
	@cat Makefile | grep -e "^.PHONY:" | grep -v "all" | cut -f2 -d' '

.PHONY: list-service
list-service:
	grpcurl -plaintext iam.core.svc.cluster.local:8002 list core.iam.IAMService

.PHONY: is-authorized
is-authorized:
	grpcurl \
		-plaintext \
		-d '{"user_id":1001, "project_id":1001, "action_name":"finding/GetFinding", "resource_name":"aws:guardduty/s3-bucket-name"}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.IsAuthorized

.PHONY: is-authorized-token
is-authorized-token:
	grpcurl \
		-plaintext \
		-d '{"access_token_id":1001, "project_id":1001, "action_name":"finding/GetFinding", "resource_name":"aws:guardduty/s3-bucket-name"}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.IsAuthorizedToken

.PHONY: is-admin
is-admin:
	grpcurl \
		-plaintext \
		-d '{"user_id":1001}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.IsAdmin

.PHONY: list-user
list-user:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "activated":"true", "name":"john"}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.ListUser

.PHONY: get-user
get-user:
	grpcurl \
		-plaintext \
		-d '{"user_id":1001, "sub":"alice"}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.GetUser

.PHONY: put-user
put-user:
	grpcurl \
		-plaintext \
		-d '{"user": {"sub":"alice", "sub":"sub", "name":"name"}}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.PutUser

.PHONY: list-role
list-role:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "name":"admin-role"}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.ListRole

.PHONY: get-role
get-role:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "role_id":1001}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.GetRole

.PHONY: put-role
put-role:
	grpcurl \
		-plaintext \
		-d '{"role":{"name":"test-role", "project_id":1001}}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.PutRole

.PHONY: delete-role
delete-role:
	grpcurl \
		-plaintext \
		-d '{"role_id":1004, "project_id":1001}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.DeleteRole

.PHONY: attach-role
attach-role:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "role_id":1005, "user_id":1001}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.AttachRole

.PHONY: detach-role
detach-role:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "role_id":1005, "user_id":1001}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.DetachRole

.PHONY: list-policy
list-policy:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "name":"admin-policy"}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.ListPolicy

.PHONY: get-policy
get-policy:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "policy_id":1001}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.GetPolicy

.PHONY: put-policy
put-policy:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "policy":{"name":"test-policy", "project_id":1001, "action_ptn":".*", "resource_ptn":".*"}}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.PutPolicy

.PHONY: delete-policy
delete-policy:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "policy_id":1004}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.DeletePolicy

.PHONY: attach-policy
attach-policy:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "role_id":1001, "policy_id":1005}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.AttachPolicy
 
.PHONY: detach-policy
detach-policy:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "role_id":1001, "policy_id":1005}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.DetachPolicy

.PHONY: list-access-token
list-access-token:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.ListAccessToken

.PHONY: authenticate-access-token
authenticate-access-token:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "access_token_id": 1001,"plain_text_token":"test-token"}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.AuthenticateAccessToken

.PHONY: put-access-token
put-access-token:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "access_token":{"plain_text_token":"test-token", "name":"test", "project_id":1001, "description":"description", "expired_at":2628676885, "last_updated_uesr_id":1001}}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.PutAccessToken

.PHONY: delete-access-token
delete-access-token:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "access_token_id": 1001}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.DeleteAccessToken

.PHONY: attach-access-token-role
attach-access-token-role:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "role_id":1002, "access_token_id": 1001}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.AttachAccessTokenRole

.PHONY: detach-access-token-role
detach-access-token-role:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "role_id":1002, "access_token_id": 1001}' \
		iam.core.svc.cluster.local:8002 core.iam.IAMService.DetachAccessTokenRole

