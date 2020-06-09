####################################################
## grpcurl example
####################################################
.PHONY: grpcurl-list
grpcurl-list:
	grpcurl -plaintext localhost:8081 list core.finding.FindingService

.PHONY: grpcurl-fiding-list
grpcurl-fiding-list:
	grpcurl \
		-plaintext \
		-d '{"project_id":["1001","1002"], "from_score":0.0, "to_score":1.0, "from_at": 1560000000, "to_at": 1660000000}' \
		localhost:8081 core.finding.FindingService.ListFinding

.PHONY: grpcurl-finding-put
grpcurl-finding-put:
	grpcurl \
		-plaintext \
		-d '{"finding":{"description":"desc", "data_source":"ds", "data_source_id":"ds-001", "resource_name":"rn", "project_id":1001, "original_score":55.51, "original_max_score":100.0, "data":"{\"key\":\"value\"}"}}' \
		localhost:8081 core.finding.FindingService.PutFinding

.PHONY: grpcurl-finding-delete
grpcurl-finding-delete:
	grpcurl \
		-plaintext \
		-d '{"finding_id": 1004}' \
		localhost:8081 core.finding.FindingService.DeleteFinding

.PHONY: grpcurl-list-finding-tag
grpcurl-list-finding-tag:
	grpcurl \
		-plaintext \
		-d '{"finding_id": 1001}' \
		localhost:8081 core.finding.FindingService.ListFindingTag

