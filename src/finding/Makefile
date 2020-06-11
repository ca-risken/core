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

.PHONY: service-list
service-list:
	grpcurl -plaintext localhost:8081 list core.finding.FindingService

.PHONY: fiding-list
fiding-list:
	grpcurl \
		-plaintext \
		-d '{"project_id":["1001","1002"], "from_score":0.0, "to_score":1.0, "from_at": 1560000000, "to_at": 1660000000}' \
		localhost:8081 core.finding.FindingService.ListFinding

.PHONY: finding-put
finding-put:
	grpcurl \
		-plaintext \
		-d '{"finding":{"description":"desc", "data_source":"ds", "data_source_id":"ds-001", "resource_name":"rn", "project_id":1001, "original_score":55.51, "original_max_score":100.0, "data":"{\"key\":\"value\"}"}}' \
		localhost:8081 core.finding.FindingService.PutFinding

.PHONY: finding-delete
finding-delete:
	grpcurl \
		-plaintext \
		-d '{"finding_id": 1004}' \
		localhost:8081 core.finding.FindingService.DeleteFinding

.PHONY: list-finding-tag
list-finding-tag:
	grpcurl \
		-plaintext \
		-d '{"finding_id": 1001}' \
		localhost:8081 core.finding.FindingService.ListFindingTag

.PHONY: tag-finding
tag-finding:
	grpcurl \
		-plaintext \
		-d '{"tag": {"finding_tag_id": 1001, "finding_id": 1001, "tag_key": "key", "tag_value": "value"}}' \
		localhost:8081 core.finding.FindingService.TagFinding

.PHONY: untag-finding
untag-finding:
	grpcurl \
		-plaintext \
		-d '{"finding_tag_id": 1001}' \
		localhost:8081 core.finding.FindingService.UntagFinding

.PHONY: list-resource
list-resource:
	grpcurl \
		-plaintext \
		-d '{"project_id":["1001","1002"], "resource_name":["rn"], "from_sum_score":0.0, "to_sum_score":999.9, "from_at": 1560000000, "to_at": 1660000000}' \
		localhost:8081 core.finding.FindingService.ListResource
