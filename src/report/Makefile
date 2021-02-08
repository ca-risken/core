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
	grpcurl -plaintext localhost:8005 list core.report.ReportService

.PHONY: get-report
get-report:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001}' \
		localhost:8005 core.report.ReportService.GetReportFinding

.PHONY: get-report-all
get-report-all:
	grpcurl \
		-plaintext \
		-d '{}' \
		localhost:8005 core.report.ReportService.GetReportFindingAll

.PHONY: collect-report
collect-report:
	grpcurl \
		-plaintext \
		-d '{}' \
		localhost:8005 core.report.ReportService.CollectReportFinding
