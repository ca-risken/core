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
	grpcurl -plaintext project.core.svc.cluster.local:8001 list core.project.ProjectService

.PHONY: list-project
list-project:
	grpcurl \
		-plaintext \
		-d '{"user_id":1002, "project_id":1001, "name":"project-a"}' \
		project.core.svc.cluster.local:8003 core.project.ProjectService.ListProject

.PHONY: create-project
create-project:
	grpcurl \
		-plaintext \
		-d '{"user_id":1001, "name":"project-x"}' \
		project.core.svc.cluster.local:8003 core.project.ProjectService.CreateProject

.PHONY: update-project
update-project:
	grpcurl \
		-plaintext \
		-d '{"project_id":1004, "name":"project-xxx"}' \
		project.core.svc.cluster.local:8003 core.project.ProjectService.UpdateProject

.PHONY: delete-project
delete-project:
	grpcurl \
		-plaintext \
		-d '{"project_id":1005}' \
		project.core.svc.cluster.local:8003 core.project.ProjectService.DeleteProject

.PHONY: tag-project
tag-project:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "tag":"key:value", "color":"blue"}' \
		project.core.svc.cluster.local:8003 core.project.ProjectService.TagProject

.PHONY: tag-project2
tag-project2:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "tag":"key:value2", "color":"green"}' \
		project.core.svc.cluster.local:8003 core.project.ProjectService.TagProject

.PHONY: untag-project
untag-project:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "tag":"key:value"}' \
		project.core.svc.cluster.local:8003 core.project.ProjectService.UntagProject

.PHONY: is-active
is-active:
	grpcurl \
		-plaintext \
		-d '{"project_id":1}' \
		project.core.svc.cluster.local:8003 core.project.ProjectService.IsActive
