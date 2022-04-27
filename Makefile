TARGETS = alert finding iam project report
MOCK_TARGETS = $(TARGETS:=.mock)
BUILD_OPT=""
IMAGE_TAG=latest
MANIFEST_TAG=latest
IMAGE_NAME=core
IMAGE_REGISTRY=local

.PHONY: all
all: run

.PHONY: install
install:
	go get \
		google.golang.org/grpc \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/envoyproxy/protoc-gen-validate@v0.6.7 \
		github.com/grpc-ecosystem/go-grpc-middleware

.PHONY: clean
clean:
	rm -f proto/*/*.pb.go
	rm -f doc/*.md

.PHONY: fmt
fmt: proto/**/*.proto
	@clang-format -i proto/**/*.proto

.PHONY: doc
doc: fmt
	protoc \
		--proto_path=proto \
		--error_format=gcc \
		-I $(GOPATH)/src/github.com/envoyproxy/protoc-gen-validate \
		--doc_out=markdown,README.md:doc \
		proto/**/*.proto;

# build without protoc-gen-validate
.PHONY: proto-without-validation
proto-without-validate: fmt
	for svc in "alert" "finding" "iam"; do \
		protoc \
			--proto_path=proto \
			--error_format=gcc \
			--go_out=plugins=grpc,paths=source_relative:proto \
			proto/$$svc/*.proto; \
	done

# build with protoc-gen-validate
.PHONY: proto-validate
proto-validate: fmt
	for svc in "project" "report"; do \
		protoc \
			--proto_path=proto \
			--error_format=gcc \
			-I $(GOPATH)/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v0.6.7 \
			--go_out=plugins=grpc,paths=source_relative:proto \
			--validate_out="lang=go,paths=source_relative:proto" \
			proto/$$svc/*.proto; \
	done

.PHONY: proto
proto : proto-without-validate proto-validate proto-mock

PHONY: build
build: test
	IMAGE_TAG=$(IMAGE_TAG) IMAGE_NAME=$(IMAGE_NAME) BUILD_OPT="$(BUILD_OPT)" . hack/docker-build.sh

PHONY: build-ci
build-ci:
	IMAGE_TAG=$(IMAGE_TAG) IMAGE_NAME=$(IMAGE_NAME) BUILD_OPT="$(BUILD_OPT)" . hack/docker-build.sh
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

PHONY: push-image
push-image:
	docker push $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

PHONY: pull-image
pull-image:
	docker pull $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

PHONY: tag-image
tag-image:
	docker tag $(SOURCE_IMAGE_NAME):$(SOURCE_IMAGE_TAG) $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

PHONY: create-manifest
create-manifest:
	docker manifest create $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(MANIFEST_TAG) $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG_BASE)_linux_amd64 $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG_BASE)_linux_arm64
	docker manifest annotate --arch amd64 $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(MANIFEST_TAG) $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG_BASE)_linux_amd64
	docker manifest annotate --arch arm64 $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(MANIFEST_TAG) $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG_BASE)_linux_arm64

PHONY: push-manifest
push-manifest:
	docker manifest push $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(MANIFEST_TAG)
	docker manifest inspect $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(MANIFEST_TAG)

PHONY: test
test:
	go test ./...

.PHONY: lint
lint: FAKE
	golangci-lint run --timeout 5m

.PHONY: generate-mock
generate-mock: proto-mock
proto-mock: $(MOCK_TARGETS)
%.mock: FAKE
	sh hack/generate-mock.sh proto/$(*)

.PHONY: list-project-service
list-project-service:
	grpcurl -plaintext core.core.svc.cluster.local:8080 list core.core.ProjectService

.PHONY: list-project
list-project:
	grpcurl \
		-plaintext \
		-d '{"user_id":1002, "project_id":1001, "name":"project-a"}' \
		core.core.svc.cluster.local:8080 core.core.ProjectService.ListProject

.PHONY: create-project
create-project:
	grpcurl \
		-plaintext \
		-d '{"user_id":1001, "name":"project-x"}' \
		core.core.svc.cluster.local:8080 core.core.ProjectService.CreateProject

.PHONY: update-project
update-project:
	grpcurl \
		-plaintext \
		-d '{"project_id":1004, "name":"project-xxx"}' \
		core.core.svc.cluster.local:8080 core.core.ProjectService.UpdateProject

.PHONY: delete-project
delete-project:
	grpcurl \
		-plaintext \
		-d '{"project_id":1005}' \
		core.core.svc.cluster.local:8080 core.core.ProjectService.DeleteProject

.PHONY: tag-project
tag-project:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "tag":"key:value", "color":"blue"}' \
		core.core.svc.cluster.local:8080 core.core.ProjectService.TagProject

.PHONY: tag-project2
tag-project2:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "tag":"key:value2", "color":"green"}' \
		core.core.svc.cluster.local:8080 core.core.ProjectService.TagProject

.PHONY: untag-project
untag-project:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "tag":"key:value"}' \
		core.core.svc.cluster.local:8080 core.core.ProjectService.UntagProject

.PHONY: is-active
is-active:
	grpcurl \
		-plaintext \
		-d '{"project_id":1}' \
		core.core.svc.cluster.local:8080 core.core.ProjectService.IsActive

.PHONY: list-alert-service
list-alert-service:
	grpcurl -plaintext core.core.svc.cluster.local:8080 list core.core.AlertService

.PHONY: list-alert
list-alert:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "description": "", "status":["ACTIVE","PENDING"], "severity":["high","medium","low"],"from_at":1560000000, "to_at":1660000000}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.ListAlert

.PHONY: get-alert
get-alert:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.GetAlert

.PHONY: put-alert
put-alert:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert":{"alert_condition_id": 1001, "description":"hogehoge", "severity": "low", "status": 1, "project_id":1001}}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.PutAlert

.PHONY: delete-alert
delete-alert:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_id":1003}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.DeleteAlert

.PHONY: list-alert_history
list-alert_history:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001,"alert_id":1001, "history_type": ["created","deleted"], "severity":["high","medium"], "from_at":1560000000, "to_at":1660000000}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.ListAlertHistory

.PHONY: get-alert_history
get-alert_history:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_history_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.GetAlertHistory

.PHONY: put-alert_history
put-alert_history:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_history":{"alert_id": 1003, "history_type":"created","description":"test_put_alert_history","finding_history":"{\"finding_id\":[1001,1002]}", "severity": "low", "project_id":1001}}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.PutAlertHistory

.PHONY: delete-alert_history
delete-alert_history:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_history_id":1003}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.DeleteAlertHistory

.PHONY: list-rel_alert_finding
list-rel_alert_finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001,"alert_id":1001, "finding_id":1001, "from_at":1560000000, "to_at":1660000000}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.ListRelAlertFinding

.PHONY: get-rel_alert_finding
get-rel_alert_finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001,"alert_id":1001, "finding_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.GetRelAlertFinding

.PHONY: put-rel_alert_finding
put-rel_alert_finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "rel_alert_finding":{"project_id":1001,"alert_id":1003, "finding_id":1001}}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.PutRelAlertFinding

.PHONY: delete-rel_alert_finding
delete-rel_alert_finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001,"alert_id":1003, "finding_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.DeleteRelAlertFinding

.PHONY: list-alert_condition
list-alert_condition:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "severity":["high","medium"], "enabled": true, "from_at":1560000000, "to_at":1660000000}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.ListAlertCondition

.PHONY: get-alert_condition
get-alert_condition:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.GetAlertCondition

.PHONY: put-alert_condition
put-alert_condition:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition":{"enabled": true, "description":"test_put_alert_condition", "severity": "low", "and_or": "or", "project_id":1001}}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.PutAlertCondition

.PHONY: delete-alert_condition
delete-alert_condition:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1003}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.DeleteAlertCondition

.PHONY: list-alert_rule
list-alert_rule:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "from_score":0.0, "to_score":1.0, "from_at":1560000000, "to_at":1660000000}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.ListAlertRule

.PHONY: get-alert_rule
get-alert_rule:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_rule_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.GetAlertRule

.PHONY: put-alert_rule
put-alert_rule:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_rule":{"name": "test_put_alert_rule", "score": 0.1, "project_id":1001}}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.PutAlertRule

.PHONY: delete-alert_rule
delete-alert_rule:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_rule_id":1003}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.DeleteAlertRule

.PHONY: list-alert_cond_rule
list-alert_cond_rule:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1001, "alert_rule_id":1001, "from_at":1560000000, "to_at":1660000000}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.ListAlertCondRule

.PHONY: get-alert_cond_rule
get-alert_cond_rule:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1001, "alert_rule_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.GetAlertCondRule

.PHONY: put-alert_cond_rule
put-alert_cond_rule:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_cond_rule":{"project_id":1001, "alert_condition_id":1003, "alert_rule_id":1001}}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.PutAlertCondRule

.PHONY: delete-alert_cond_rule
delete-alert_cond_rule:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1003, "alert_rule_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.DeleteAlertCondRule

.PHONY: list-notification
list-notification:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "type":"slack", "from_at":1560000000, "to_at":1660000000}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.ListNotification

.PHONY: get-notification
get-notification:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "notification_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.GetNotification

.PHONY: put-notification
put-notification:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001,"notification":{"project_id":1001, "name":"test_notification","type":"slack", "notify_setting":"{\"webhook_url\":\"http://hogehoge.com/fuga/piyo\"}"}}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.PutNotification

.PHONY: delete-notification
delete-notification:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "notification_id":1003}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.DeleteNotification

.PHONY: test-notification
test-notification:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "notification_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.TestNotification

.PHONY: analyze-alert
analyze-alert:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.AnalyzeAlert

.PHONY: list-alert_cond_notification
list-alert_cond_notification:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1001, "notification_id":1001, "from_at":1560000000, "to_at":1660000000}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.ListAlertCondNotification

.PHONY: get-alert_cond_notification
get-alert_cond_notification:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1001, "notification_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.GetAlertCondNotification

.PHONY: put-alert_cond_notification
put-alert_cond_notification:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_cond_notification":{"project_id":1001, "alert_condition_id":1003, "notification_id":1001,"cache_second":3600,"notified_at":1560000000}}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.PutAlertCondNotification

.PHONY: delete-alert_cond_notification
delete-alert_cond_notification:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1003, "notification_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.AlertService.DeleteAlertCondNotification

.PHONY: analyze-alert-all
analyze-alert-all:
	grpcurl \
		-plaintext \
		core.core.svc.cluster.local:8080 core.core.AlertService.AnalyzeAlertAll

.PHONY: list-finding-service
list-finding-service:
	grpcurl -plaintext core.core.svc.cluster.local:8080 list core.core.FindingService

.PHONY: list-finding
list-finding:
	grpcurl \
		-plaintext \
		-d '{"finding_id": "1040", "project_id":1001, "status":0, "sort": "finding_id", "direction": "desc", "offset":0, "limit":10}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.ListFinding

.PHONY: get-finding
get-finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "finding_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.GetFinding

.PHONY: put-finding
put-finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "finding":{"description":"desc", "data_source":"ds", "data_source_id":"ds-001", "resource_name":"rn", "project_id":1001, "original_score":55.51, "original_max_score":100.0, "data":"{\"key\":\"value\"}"}}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.PutFinding

.PHONY: delete-finding
delete-finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "finding_id":1004}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.DeleteFinding

.PHONY: list-finding-tag
list-finding-tag:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "finding_id":1001, "sort":"tag", "direction": "desc", "offset":0, "limit":1}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.ListFindingTag

.PHONY: list-finding-tag-name
list-finding-tag-name:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "sort":"tag", "direction": "asc", "offset":0}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.ListFindingTagName

.PHONY: tag-finding
tag-finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "tag":{"finding_id":1001, "project_id":1001, "tag":"tag"}}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.TagFinding

.PHONY: untag-finding
untag-finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "finding_tag_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.UntagFinding

.PHONY: clear-score
clear-score:
	grpcurl \
		-plaintext \
		-d '{"data_source":"aws:guard-duty", "project_id":1001, "tag":["test1", "test2"]}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.ClearScore

.PHONY: list-resource
list-resource:
	grpcurl \
		-plaintext \
		-d '{"resource_id":"1001", "project_id":1001, "from_sum_score":0.0, "to_sum_score":999.9, "from_at":1560000000, "to_at":253402268399, "tag": ["tag1", "tag:key"],"sort": "resource_id", "direction": "desc", "offset":0, "limit":10}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.ListResource

.PHONY: get-resource
get-resource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001,"resource_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.GetResource

.PHONY: put-resource
put-resource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "resource":{"resource_name":"rn-test", "project_id":1001}}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.PutResource

.PHONY: delete-resource
delete-resource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "resource_id":1004}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.DeleteResource

.PHONY: list-resource-tag
list-resource-tag:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "resource_id":1001, "sort": "tag", "direction": "desc", "offset":0, "limit":1}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.ListResourceTag

.PHONY: list-resource-tag-name
list-resource-tag-name:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "sort": "tag", "direction": "desc", "offset":0, "limit":1}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.ListResourceTagName

.PHONY: tag-resource
tag-resource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "tag":{"resource_id":1001, "project_id":1001, "tag":"tag"}}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.TagResource

.PHONY: untag-resource
untag-resource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "resource_tag_id":1003}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.UntagResource

.PHONY: get-pend-finding
get-pend-finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001,"finding_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.GetPendFinding

.PHONY: put-pend-finding
put-pend-finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "pend_finding":{"finding_id":1001, "project_id":1001, "note":"note"}}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.PutPendFinding

.PHONY: delete-pend-finding
delete-pend-finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "finding_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.DeletePendFinding

.PHONY: list-finding-setting
list-finding-setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.ListFindingSetting

.PHONY: get-finding-setting
get-finding-setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001,"finding_setting_id":1003}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.GetFindingSetting

.PHONY: put-finding-setting
put-finding-setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "finding_setting":{"project_id":1001, "resource_name":"rn", "status":1, "setting": "{}"}}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.PutFindingSetting

.PHONY: delete-finding-setting
delete-finding-setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "finding_setting_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.DeleteFindingSetting

.PHONY: get-recommend
get-recommend:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "finding_id":1}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.GetRecommend

.PHONY: put-recommend
put-recommend:
	grpcurl \
		-plaintext \
		-d '{"project_id":1,"finding_id":1, "data_source":"ds", "type":"c", "risk":"critical", "recommendation":"..."}' \
		core.core.svc.cluster.local:8080 core.core.FindingService.PutRecommend

.PHONY: list-report-service
list-report-service:
	grpcurl -plaintext core.core.svc.cluster.local:8080 list core.core.ReportService

.PHONY: get-report
get-report:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.ReportService.GetReportFinding

.PHONY: get-report-all
get-report-all:
	grpcurl \
		-plaintext \
		-d '{}' \
		core.core.svc.cluster.local:8080 core.core.ReportService.GetReportFindingAll

.PHONY: collect-report
collect-report:
	grpcurl \
		-plaintext \
		-d '{}' \
		core.core.svc.cluster.local:8080 core.core.ReportService.CollectReportFinding

.PHONY: list-iam-service
list-iam-service:
	grpcurl -plaintext core.core.svc.cluster.local:8080 list core.core.IAMService

.PHONY: is-authorized
is-authorized:
	grpcurl \
		-plaintext \
		-d '{"user_id":1001, "project_id":1001, "action_name":"finding/GetFinding", "resource_name":"aws:guardduty/s3-bucket-name"}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.IsAuthorized

.PHONY: is-authorized-token
is-authorized-token:
	grpcurl \
		-plaintext \
		-d '{"access_token_id":1001, "project_id":1001, "action_name":"finding/GetFinding", "resource_name":"aws:guardduty/s3-bucket-name"}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.IsAuthorizedToken

.PHONY: is-admin
is-admin:
	grpcurl \
		-plaintext \
		-d '{"user_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.IsAdmin

.PHONY: list-user
list-user:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "activated":"true", "name":"john"}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.ListUser

.PHONY: get-user
get-user:
	grpcurl \
		-plaintext \
		-d '{"user_id":1001, "sub":"alice"}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.GetUser

.PHONY: put-user
put-user:
	grpcurl \
		-plaintext \
		-d '{"user": {"sub":"alice", "sub":"sub", "name":"name"}}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.PutUser

.PHONY: list-role
list-role:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "name":"admin-role"}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.ListRole

.PHONY: get-role
get-role:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "role_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.GetRole

.PHONY: put-role
put-role:
	grpcurl \
		-plaintext \
		-d '{"role":{"name":"test-role", "project_id":1001}}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.PutRole

.PHONY: delete-role
delete-role:
	grpcurl \
		-plaintext \
		-d '{"role_id":1004, "project_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.DeleteRole

.PHONY: attach-role
attach-role:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "role_id":1005, "user_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.AttachRole

.PHONY: detach-role
detach-role:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "role_id":1005, "user_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.DetachRole

.PHONY: list-policy
list-policy:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "name":"admin-policy"}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.ListPolicy

.PHONY: get-policy
get-policy:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "policy_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.GetPolicy

.PHONY: put-policy
put-policy:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "policy":{"name":"test-policy", "project_id":1001, "action_ptn":".*", "resource_ptn":".*"}}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.PutPolicy

.PHONY: delete-policy
delete-policy:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "policy_id":1004}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.DeletePolicy

.PHONY: attach-policy
attach-policy:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "role_id":1001, "policy_id":1005}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.AttachPolicy
 
.PHONY: detach-policy
detach-policy:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "role_id":1001, "policy_id":1005}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.DetachPolicy

.PHONY: list-access-token
list-access-token:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.ListAccessToken

.PHONY: authenticate-access-token
authenticate-access-token:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "access_token_id": 1001,"plain_text_token":"test-token"}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.AuthenticateAccessToken

.PHONY: put-access-token
put-access-token:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "access_token":{"plain_text_token":"test-token", "name":"test", "project_id":1001, "description":"description", "expired_at":2628676885, "last_updated_uesr_id":1001}}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.PutAccessToken

.PHONY: delete-access-token
delete-access-token:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "access_token_id": 1001}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.DeleteAccessToken

.PHONY: attach-access-token-role
attach-access-token-role:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "role_id":1002, "access_token_id": 1001}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.AttachAccessTokenRole

.PHONY: detach-access-token-role
detach-access-token-role:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "role_id":1002, "access_token_id": 1001}' \
		core.core.svc.cluster.local:8080 core.core.IAMService.DetachAccessTokenRole

.PHONY: analyze-access-token-expiration
analyze-access-token-expiration:
	grpcurl \
		-plaintext \
		core.core.svc.cluster.local:8080 core.core.IAMService.AnalyzeTokenExpiration

FAKE:
