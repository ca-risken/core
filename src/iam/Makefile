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
	grpcurl -plaintext localhost:8082 list core.iam.IAMService

.PHONY: is-authorized
is-authorized:
	grpcurl \
		-plaintext \
		-d '{"user_id":1001, "project_id":1001, "action_name":"finding/GetFinding", "resource_name":"aws:guardduty/s3-bucket-name"}' \
		localhost:8082 core.iam.IAMService.IsAuthorized
