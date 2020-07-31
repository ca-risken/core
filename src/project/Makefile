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
	grpcurl -plaintext localhost:8001 list core.project.ProjectService

.PHONY: list-project
list-project:
	grpcurl \
		-plaintext \
		-d '{"user_id":1001, "project_id":1001, "name":"project-a"}' \
		localhost:8003 core.project.ProjectService.ListProject

.PHONY: create-project
create-project:
	grpcurl \
		-plaintext \
		-d '{"user_id":1001, "name":"project-x"}' \
		localhost:8003 core.project.ProjectService.CreateProject

.PHONY: update-project
update-project:
	grpcurl \
		-plaintext \
		-d '{"project_id":1004, "name":"project-xxx"}' \
		localhost:8003 core.project.ProjectService.UpdateProject

.PHONY: delete-project
delete-project:
	grpcurl \
		-plaintext \
		-d '{"project_id":1005}' \
		localhost:8003 core.project.ProjectService.DeleteProject
