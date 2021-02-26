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
	grpcurl -plaintext localhost:8001 list core.finding.FindingService

.PHONY: list-finding
list-finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "sort": "finding_id", "direction": "desc", "offset":0, "limit":1}' \
		localhost:8001 core.finding.FindingService.ListFinding

.PHONY: get-finding
get-finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "finding_id":1001}' \
		localhost:8001 core.finding.FindingService.GetFinding

.PHONY: put-finding
put-finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "finding":{"description":"desc", "data_source":"ds", "data_source_id":"ds-001", "resource_name":"rn", "project_id":1001, "original_score":55.51, "original_max_score":100.0, "data":"{\"key\":\"value\"}"}}' \
		localhost:8001 core.finding.FindingService.PutFinding

.PHONY: delete-finding
delete-finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "finding_id":1004}' \
		localhost:8001 core.finding.FindingService.DeleteFinding

.PHONY: list-finding-tag
list-finding-tag:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "finding_id":1001, "sort":"tag", "direction": "desc", "offset":0, "limit":1}' \
		localhost:8001 core.finding.FindingService.ListFindingTag

.PHONY: list-finding-tag-name
list-finding-tag-name:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "sort":"tag", "direction": "asc", "offset":0}' \
		localhost:8001 core.finding.FindingService.ListFindingTagName

.PHONY: tag-finding
tag-finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "tag":{"finding_id":1001, "project_id":1001, "tag":"tag"}}' \
		localhost:8001 core.finding.FindingService.TagFinding

.PHONY: untag-finding
untag-finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "finding_tag_id":1001}' \
		localhost:8001 core.finding.FindingService.UntagFinding

.PHONY: list-resource
list-resource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "from_sum_score":0.0, "to_sum_score":999.9, "from_at":1560000000, "to_at":253402268399, "tag": ["tag1", "tag:key"],"sort": "resource_id", "direction": "desc", "offset":0, "limit":1}' \
		localhost:8001 core.finding.FindingService.ListResource

.PHONY: get-resource
get-resource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001,"resource_id":1001}' \
		localhost:8001 core.finding.FindingService.GetResource

.PHONY: put-resource
put-resource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "resource":{"resource_name":"rn-test", "project_id":1001}}' \
		localhost:8001 core.finding.FindingService.PutResource

.PHONY: delete-resource
delete-resource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "resource_id":1004}' \
		localhost:8001 core.finding.FindingService.DeleteResource

.PHONY: list-resource-tag
list-resource-tag:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "resource_id":1001, "sort": "tag", "direction": "desc", "offset":0, "limit":1}' \
		localhost:8001 core.finding.FindingService.ListResourceTag

.PHONY: list-resource-tag-name
list-resource-tag-name:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "sort": "tag", "direction": "desc", "offset":0, "limit":1}' \
		localhost:8001 core.finding.FindingService.ListResourceTagName

.PHONY: tag-resource
tag-resource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "tag":{"resource_id":1001, "project_id":1001, "tag":"tag"}}' \
		localhost:8001 core.finding.FindingService.TagResource

.PHONY: untag-resource
untag-resource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "resource_tag_id":1003}' \
		localhost:8001 core.finding.FindingService.UntagResource

.PHONY: get-pend-finding
get-pend-finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001,"finding_id":1001}' \
		localhost:8001 core.finding.FindingService.GetPendFinding

.PHONY: put-pend-finding
put-pend-finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "pend_finding":{"finding_id":1001, "project_id":1001}}' \
		localhost:8001 core.finding.FindingService.PutPendFinding

.PHONY: delete-pend-finding
delete-pend-finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "finding_id":1001}' \
		localhost:8001 core.finding.FindingService.DeletePendFinding
