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
	grpcurl -plaintext localhost:8002 list core.iam.IAMService

.PHONY: list-user
list-user:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "activated":"true", "name":"alice"}' \
		localhost:8002 core.iam.IAMService.ListUser

.PHONY: get-user
get-user:
	grpcurl \
		-plaintext \
		-d '{"user_id":1001, "sub":"alice"}' \
		localhost:8002 core.iam.IAMService.GetUser

.PHONY: put-user
put-user:
	grpcurl \
		-plaintext \
		-d '{"user": {"sub":"alice", "sub":"sub", "name":"name"}}' \
		localhost:8002 core.iam.IAMService.PutUser

.PHONY: list-role
list-role:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "name":"admin-role"}' \
		localhost:8002 core.iam.IAMService.ListRole

.PHONY: get-role
get-role:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "role_id":1001}' \
		localhost:8002 core.iam.IAMService.GetRole

.PHONY: put-role
put-role:
	grpcurl \
		-plaintext \
		-d '{"role":{"name":"test-role", "project_id":1001}}' \
		localhost:8002 core.iam.IAMService.PutRole

.PHONY: delete-role
delete-role:
	grpcurl \
		-plaintext \
		-d '{"role_id":1004, "project_id":1001}' \
		localhost:8002 core.iam.IAMService.DeleteRole

.PHONY: attach-role
attach-role:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "role_id":1005, "user_id":1001}' \
		localhost:8002 core.iam.IAMService.AttachRole

.PHONY: detach-role
detach-role:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "role_id":1005, "user_id":1001}' \
		localhost:8002 core.iam.IAMService.DetachRole

.PHONY: is-authorized
is-authorized:
	grpcurl \
		-plaintext \
		-d '{"user_id":1001, "project_id":1001, "action_name":"finding/GetFinding", "resource_name":"aws:guardduty/s3-bucket-name"}' \
		localhost:8002 core.iam.IAMService.IsAuthorized
