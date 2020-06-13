####################################################
## curl example
####################################################
.PHONY: all
all: help

.PHONY: help
help:
	@echo "Usage: make <sub-command>"
	@echo "\n---------------- sub-command list ----------------"
	@cat Makefile | grep -e "^.PHONY:" | grep -v "all" | cut -f2 -d' '

.PHONY: list-finding
list-finding:
	curl -is -XGET \
		http://localhost:8080/finding?project_id=1001,1002

.PHONY: get-finding
get-finding:
	curl -is -XGET \
		http://localhost:8080/finding/1001/

.PHONY: put-finding
put-finding:
	curl -is -XPOST -H 'Content-Type: application/json' \
		--data '{"description":"desc", "data_source":"ds", "data_source_id":"ds-004", "resource_name":"rn", "project_id":1001, "original_score":55.51, "original_max_score":100.0, "data":"{\"key\":\"value\"}"}' \
		http://localhost:8080/finding/put/

.PHONY: delete-finding
delete-finding:
	curl -is -XPOST -H 'Content-Type: application/json' \
		--data '{"finding_id":1008}' \
		http://localhost:8080/finding/delete/

.PHONY: list-finding-tag
list-finding-tag:
	curl -is -XGET \
		http://localhost:8080/finding/1001/tag/

.PHONY: tag-finding
tag-finding:
	curl -is -XPOST -H 'Content-Type: application/json' \
		--data '{"finding_id":1001, "tag_key":"test", "tag_value":"true"}' \
		http://localhost:8080/finding/tag/

.PHONY: untag-finding
untag-finding:
	curl -is -XPOST -H 'Content-Type: application/json' \
		--data '{"finding_tag_id":1002}' \
		http://localhost:8080/finding/untag/

.PHONY: list-resource
list-resource:
	curl -is -XGET \
		http://localhost:8080/resource?project_id=1001,1002

.PHONY: get-resource
get-resource:
	curl -is -XGET \
		http://localhost:8080/resource/1001/

.PHONY: put-resource
put-resource:
	curl -is -XPOST -H 'Content-Type: application/json' \
		--data '{"resource_name":"rn", "project_id":1001}' \
		http://localhost:8080/resource/put/

.PHONY: delete-resource
delete-resource:
	curl -is -XPOST -H 'Content-Type: application/json' \
		--data '{"resource_id":1003}' \
		http://localhost:8080/resource/delete/

.PHONY: list-resource-tag
list-resource-tag:
	curl -is -XGET -H 'Content-Type: application/json' \
		http://localhost:8080/resource/1001/tag/

.PHONY: tag-resource
tag-resource:
	curl -is -XPOST -H 'Content-Type: application/json' \
		--data '{"resource_id":1001, "tag_key":"test", "tag_value":"true"}' \
		http://localhost:8080/resource/tag/

.PHONY: untag-resource
untag-resource:
	curl -is -XPOST -H 'Content-Type: application/json' \
		--data '{"resource_tag_id":1005}' \
		http://localhost:8080/resource/untag/
