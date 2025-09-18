TARGETS = ai alert finding iam project report organization_iam organization
MOCK_TARGETS = $(TARGETS:=.mock)
BUILD_OPT=""
IMAGE_TAG=latest
MANIFEST_TAG=latest
IMAGE_NAME=core
IMAGE_REGISTRY=local
GRPCURL=kubectl run grpcurl --image=fullstorydev/grpcurl -n core --restart=Never --rm -it --
CORE_API_ADDR=core.core.svc.cluster.local:8080

.PHONY: all
all: run

.PHONY: install
install:
	go install \
		google.golang.org/grpc
	go install \
		github.com/golang/protobuf/protoc-gen-go
	go install \
		github.com/envoyproxy/protoc-gen-validate@v0.6.7
	go install \
		github.com/grpc-ecosystem/go-grpc-middleware

.PHONY: clean
clean:
	rm -f proto/*/*.pb.go
	rm -f doc/*.md

.PHONY: fmt
fmt: proto/**/*.proto
	@clang-format -i proto/**/*.proto

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
	for svc in "project" "report" "ai" "organization" "organization_iam"; do \
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
	GO111MODULE=on go test ./...

PHONY: go-test
go-test:
	GO111MODULE=on go test ./...

.PHONY: lint
lint: FAKE
	GO111MODULE=on GOFLAGS=-buildvcs=false golangci-lint run --timeout 5m

.PHONY: generate-mock
generate-mock: proto-mock repository-mock ai-mock

.PHONY: proto-mock
proto-mock: $(MOCK_TARGETS)
%.mock: FAKE
	sh hack/generate-mock.sh proto/$(*)

.PHONY: repository-mock
repository-mock: FAKE
	sh hack/generate-mock.sh pkg/db

.PHONY: ai-mock
ai-mock: FAKE
	sh hack/generate-mock.sh pkg/ai

.PHONY: list-project-service
list-project-service:
	$(GRPCURL) -plaintext $(CORE_API_ADDR) list core.project.ProjectService

.PHONY: list-project
list-project:
	$(GRPCURL) \
		-plaintext \
		-d '{"user_id":1002, "project_id":1001, "name":"project-a", "organization_id":100}' \
		$(CORE_API_ADDR) core.project.ProjectService.ListProject

.PHONY: create-project
create-project:
	$(GRPCURL) \
		-plaintext \
		-d '{"user_id":1001, "name":"project-x"}' \
		$(CORE_API_ADDR) core.project.ProjectService.CreateProject

.PHONY: update-project
update-project:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1004, "name":"project-xxx"}' \
		$(CORE_API_ADDR) core.project.ProjectService.UpdateProject

.PHONY: delete-project
delete-project:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1005}' \
		$(CORE_API_ADDR) core.project.ProjectService.DeleteProject

.PHONY: clean-project
clean-project:
	$(GRPCURL) \
		-plaintext \
		-d '' \
		$(CORE_API_ADDR) core.project.ProjectService.CleanProject

.PHONY: tag-project
tag-project:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "tag":"key:value", "color":"blue"}' \
		$(CORE_API_ADDR) core.project.ProjectService.TagProject

.PHONY: tag-project2
tag-project2:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "tag":"key:value2", "color":"green"}' \
		$(CORE_API_ADDR) core.project.ProjectService.TagProject

.PHONY: untag-project
untag-project:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "tag":"key:value"}' \
		$(CORE_API_ADDR) core.project.ProjectService.UntagProject

.PHONY: is-active
is-active:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1}' \
		$(CORE_API_ADDR) core.project.ProjectService.IsActive

.PHONY: list-alert-service
list-alert-service:
	$(GRPCURL) -plaintext $(CORE_API_ADDR) list core.alert.AlertService

.PHONY: list-alert
list-alert:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "description": "", "status":["ACTIVE","PENDING"], "severity":["high","medium","low"],"from_at":1560000000, "to_at":1660000000}' \
		$(CORE_API_ADDR) core.alert.AlertService.ListAlert

.PHONY: get-alert
get-alert:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_id":1001}' \
		$(CORE_API_ADDR) core.alert.AlertService.GetAlert

.PHONY: put-alert
put-alert:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert":{"alert_condition_id": 1001, "description":"hogehoge", "severity": "low", "status": 1, "project_id":1001}}' \
		$(CORE_API_ADDR) core.alert.AlertService.PutAlert

.PHONY: put-alert-first-viewed-at
put-alert-first-viewed-at:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_id":1001}' \
		$(CORE_API_ADDR) core.alert.AlertService.PutAlertFirstViewedAt

.PHONY: delete-alert
delete-alert:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_id":1003}' \
		$(CORE_API_ADDR) core.alert.AlertService.DeleteAlert

.PHONY: list-alert_history
list-alert_history:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001,"alert_id":1001, "history_type": ["created","deleted"], "severity":["high","medium"], "from_at":1560000000, "to_at":1660000000}' \
		$(CORE_API_ADDR) core.alert.AlertService.ListAlertHistory

.PHONY: get-alert_history
get-alert_history:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_history_id":1001}' \
		$(CORE_API_ADDR) core.alert.AlertService.GetAlertHistory

.PHONY: put-alert_history
put-alert_history:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_history":{"alert_id": 1003, "history_type":"created","description":"test_put_alert_history","finding_history":"{\"finding_id\":[1001,1002]}", "severity": "low", "project_id":1001}}' \
		$(CORE_API_ADDR) core.alert.AlertService.PutAlertHistory

.PHONY: delete-alert_history
delete-alert_history:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_history_id":1003}' \
		$(CORE_API_ADDR) core.alert.AlertService.DeleteAlertHistory

.PHONY: list-rel_alert_finding
list-rel_alert_finding:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001,"alert_id":1001, "finding_id":1001, "from_at":1560000000, "to_at":1660000000}' \
		$(CORE_API_ADDR) core.alert.AlertService.ListRelAlertFinding

.PHONY: get-rel_alert_finding
get-rel_alert_finding:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001,"alert_id":1001, "finding_id":1001}' \
		$(CORE_API_ADDR) core.alert.AlertService.GetRelAlertFinding

.PHONY: put-rel_alert_finding
put-rel_alert_finding:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "rel_alert_finding":{"project_id":1001,"alert_id":1003, "finding_id":1001}}' \
		$(CORE_API_ADDR) core.alert.AlertService.PutRelAlertFinding

.PHONY: delete-rel_alert_finding
delete-rel_alert_finding:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001,"alert_id":1003, "finding_id":1001}' \
		$(CORE_API_ADDR) core.alert.AlertService.DeleteRelAlertFinding

.PHONY: list-alert_condition
list-alert_condition:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "severity":["high","medium"], "enabled": true, "from_at":1560000000, "to_at":1660000000}' \
		$(CORE_API_ADDR) core.alert.AlertService.ListAlertCondition

.PHONY: get-alert_condition
get-alert_condition:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1001}' \
		$(CORE_API_ADDR) core.alert.AlertService.GetAlertCondition

.PHONY: put-alert_condition
put-alert_condition:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_condition":{"enabled": true, "description":"test_put_alert_condition", "severity": "low", "and_or": "or", "project_id":1001}}' \
		$(CORE_API_ADDR) core.alert.AlertService.PutAlertCondition

.PHONY: delete-alert_condition
delete-alert_condition:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1003}' \
		$(CORE_API_ADDR) core.alert.AlertService.DeleteAlertCondition

.PHONY: list-alert_rule
list-alert_rule:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "from_score":0.0, "to_score":1.0, "from_at":1560000000, "to_at":1660000000}' \
		$(CORE_API_ADDR) core.alert.AlertService.ListAlertRule

.PHONY: get-alert_rule
get-alert_rule:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_rule_id":1001}' \
		$(CORE_API_ADDR) core.alert.AlertService.GetAlertRule

.PHONY: put-alert_rule
put-alert_rule:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_rule":{"name": "test_put_alert_rule", "score": 0.1, "project_id":1001}}' \
		$(CORE_API_ADDR) core.alert.AlertService.PutAlertRule

.PHONY: delete-alert_rule
delete-alert_rule:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_rule_id":1003}' \
		$(CORE_API_ADDR) core.alert.AlertService.DeleteAlertRule

.PHONY: list-alert_cond_rule
list-alert_cond_rule:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1001, "alert_rule_id":1001, "from_at":1560000000, "to_at":1660000000}' \
		$(CORE_API_ADDR) core.alert.AlertService.ListAlertCondRule

.PHONY: get-alert_cond_rule
get-alert_cond_rule:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1001, "alert_rule_id":1001}' \
		$(CORE_API_ADDR) core.alert.AlertService.GetAlertCondRule

.PHONY: put-alert_cond_rule
put-alert_cond_rule:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_cond_rule":{"project_id":1001, "alert_condition_id":1003, "alert_rule_id":1001}}' \
		$(CORE_API_ADDR) core.alert.AlertService.PutAlertCondRule

.PHONY: delete-alert_cond_rule
delete-alert_cond_rule:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1003, "alert_rule_id":1001}' \
		$(CORE_API_ADDR) core.alert.AlertService.DeleteAlertCondRule

.PHONY: list-notification
list-notification:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "type":"slack"}' \
		$(CORE_API_ADDR) core.alert.AlertService.ListNotification

.PHONY: list-notification-for-internal
list-notification-for-internal:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "type":"slack"}' \
		$(CORE_API_ADDR) core.alert.AlertService.ListNotificationForInternal

.PHONY: get-notification
get-notification:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "notification_id":1001}' \
		$(CORE_API_ADDR) core.alert.AlertService.GetNotification

.PHONY: put-notification
put-notification:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001,"notification":{"project_id":1001, "name":"test_notification","type":"slack", "notify_setting":"{\"webhook_url\":\"http://hogehoge.com/fuga/piyo\"}"}}' \
		$(CORE_API_ADDR) core.alert.AlertService.PutNotification

.PHONY: put-notification2
put-notification2:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001,"notification":{"project_id":1001, "name":"test_notification","type":"slack", "notify_setting":"{\"channel_id\":\"C023QE39Q0J\"}"}}' \
		$(CORE_API_ADDR) core.alert.AlertService.PutNotification

.PHONY: delete-notification
delete-notification:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "notification_id":1003}' \
		$(CORE_API_ADDR) core.alert.AlertService.DeleteNotification

.PHONY: test-notification
test-notification:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "notification_id":1001}' \
		$(CORE_API_ADDR) core.alert.AlertService.TestNotification

.PHONY: request-project-role-notification
request-project-role-notification:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001,  "user_id": 1001}' \
		$(CORE_API_ADDR) core.alert.AlertService.RequestProjectRoleNotification

.PHONY: analyze-alert
analyze-alert:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001}' \
		$(CORE_API_ADDR) core.alert.AlertService.AnalyzeAlert

.PHONY: list-alert_cond_notification
list-alert_cond_notification:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1001, "notification_id":1001, "from_at":1560000000, "to_at":1660000000}' \
		$(CORE_API_ADDR) core.alert.AlertService.ListAlertCondNotification

.PHONY: get-alert_cond_notification
get-alert_cond_notification:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1001, "notification_id":1001}' \
		$(CORE_API_ADDR) core.alert.AlertService.GetAlertCondNotification

.PHONY: put-alert_cond_notification
put-alert_cond_notification:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_cond_notification":{"project_id":1001, "alert_condition_id":1003, "notification_id":1001,"cache_second":3600,"notified_at":1560000000}}' \
		$(CORE_API_ADDR) core.alert.AlertService.PutAlertCondNotification

.PHONY: delete-alert_cond_notification
delete-alert_cond_notification:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1003, "notification_id":1001}' \
		$(CORE_API_ADDR) core.alert.AlertService.DeleteAlertCondNotification

.PHONY: analyze-alert-all
analyze-alert-all:
	$(GRPCURL) \
		-plaintext \
		$(CORE_API_ADDR) core.alert.AlertService.AnalyzeAlertAll

.PHONY: list-finding-service
list-finding-service:
	$(GRPCURL) -plaintext $(CORE_API_ADDR) list core.finding.FindingService

.PHONY: list-finding
list-finding:
	$(GRPCURL) \
		-plaintext \
		-d '{"alert_id":2, "project_id":1, "status":0, "sort":"finding_id", "direction":"desc", "offset":0, "limit":10}' \
		$(CORE_API_ADDR) core.finding.FindingService.ListFinding

.PHONY: get-finding
get-finding:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "finding_id":1001}' \
		$(CORE_API_ADDR) core.finding.FindingService.GetFinding

.PHONY: put-finding
put-finding:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "finding":{"description":"desc", "data_source":"ds", "data_source_id":"ds-001", "resource_name":"rn", "project_id":1001, "original_score":55.51, "original_max_score":100.0, "data":"{\"key\":\"value\"}"}}' \
		$(CORE_API_ADDR) core.finding.FindingService.PutFinding

.PHONY: delete-finding
delete-finding:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "finding_id":1004}' \
		$(CORE_API_ADDR) core.finding.FindingService.DeleteFinding

.PHONY: list-finding-tag
list-finding-tag:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "finding_id":1001, "sort":"tag", "direction": "desc", "offset":0, "limit":1}' \
		$(CORE_API_ADDR) core.finding.FindingService.ListFindingTag

.PHONY: list-finding-tag-name
list-finding-tag-name:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "sort":"tag", "direction": "asc", "offset":0}' \
		$(CORE_API_ADDR) core.finding.FindingService.ListFindingTagName

.PHONY: tag-finding
tag-finding:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "tag":{"finding_id":1001, "project_id":1001, "tag":"tag"}}' \
		$(CORE_API_ADDR) core.finding.FindingService.TagFinding

.PHONY: untag-finding
untag-finding:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "finding_tag_id":1001}' \
		$(CORE_API_ADDR) core.finding.FindingService.UntagFinding

.PHONY: clear-score
clear-score:
	$(GRPCURL) \
		-plaintext \
		-d '{"data_source":"aws:guard-duty", "project_id":1001, "tag":["aws", "guardduty"], "before_at":1675071350 }' \
		$(CORE_API_ADDR) core.finding.FindingService.ClearScore

.PHONY: list-resource
list-resource:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "from_at":1560000000, "to_at":253402268399, "sort": "resource_id", "direction": "desc", "offset":0, "limit":30, "namespace":"aws", "resource_type":"s3", "tag":["aws", "access-analyzer"]}' \
		$(CORE_API_ADDR) core.finding.FindingService.ListResource

.PHONY: get-resource
get-resource:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001,"resource_id":1001}' \
		$(CORE_API_ADDR) core.finding.FindingService.GetResource

.PHONY: put-resource
put-resource:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "resource":{"resource_name":"rn-test", "project_id":1001}}' \
		$(CORE_API_ADDR) core.finding.FindingService.PutResource

.PHONY: delete-resource
delete-resource:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "resource_id":1004}' \
		$(CORE_API_ADDR) core.finding.FindingService.DeleteResource

.PHONY: list-resource-tag
list-resource-tag:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "resource_id":1001, "sort": "tag", "direction": "desc", "offset":0, "limit":1}' \
		$(CORE_API_ADDR) core.finding.FindingService.ListResourceTag

.PHONY: list-resource-tag-name
list-resource-tag-name:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "sort": "tag", "direction": "desc", "offset":0, "limit":1}' \
		$(CORE_API_ADDR) core.finding.FindingService.ListResourceTagName

.PHONY: tag-resource
tag-resource:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "tag":{"resource_id":1001, "project_id":1001, "tag":"tag"}}' \
		$(CORE_API_ADDR) core.finding.FindingService.TagResource

.PHONY: untag-resource
untag-resource:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "resource_tag_id":1003}' \
		$(CORE_API_ADDR) core.finding.FindingService.UntagResource

.PHONY: untag-by-resource-name
untag-by-resource-name:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "resource_name":"arn:aws:s3:::test-bucket", "tag":"aws"}' \
		$(CORE_API_ADDR) core.finding.FindingService.UntagByResourceName

.PHONY: get-pend-finding
get-pend-finding:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001,"finding_id":1001}' \
		$(CORE_API_ADDR) core.finding.FindingService.GetPendFinding

.PHONY: put-pend-finding
put-pend-finding:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "pend_finding":{"finding_id":1001, "project_id":1001, "pend_user_id":1001, "note":"note"}}' \
		$(CORE_API_ADDR) core.finding.FindingService.PutPendFinding

.PHONY: put-pend-finding-expired
put-pend-finding-expired:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "pend_finding":{"finding_id":1001, "project_id":1001, "note":"note", "expired_at":1675868400}}' \
		$(CORE_API_ADDR) core.finding.FindingService.PutPendFinding

.PHONY: put-pend-finding-active
put-pend-finding-active:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "pend_finding":{"finding_id":1001, "project_id":1001, "note":"note", "expired_at":253402182000}}' \
		$(CORE_API_ADDR) core.finding.FindingService.PutPendFinding

.PHONY: delete-pend-finding
delete-pend-finding:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "finding_id":1001}' \
		$(CORE_API_ADDR) core.finding.FindingService.DeletePendFinding

.PHONY: list-finding-setting
list-finding-setting:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001}' \
		$(CORE_API_ADDR) core.finding.FindingService.ListFindingSetting

.PHONY: get-finding-setting
get-finding-setting:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001,"finding_setting_id":1003}' \
		$(CORE_API_ADDR) core.finding.FindingService.GetFindingSetting

.PHONY: put-finding-setting
put-finding-setting:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "finding_setting":{"project_id":1001, "resource_name":"rn", "status":1, "setting": "{}"}}' \
		$(CORE_API_ADDR) core.finding.FindingService.PutFindingSetting

.PHONY: delete-finding-setting
delete-finding-setting:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "finding_setting_id":1001}' \
		$(CORE_API_ADDR) core.finding.FindingService.DeleteFindingSetting

.PHONY: get-recommend
get-recommend:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1, "finding_id":1}' \
		$(CORE_API_ADDR) core.finding.FindingService.GetRecommend

.PHONY: put-recommend
put-recommend:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1,"finding_id":1, "data_source":"ds", "type":"c", "risk":"critical", "recommendation":"..."}' \
		$(CORE_API_ADDR) core.finding.FindingService.PutRecommend

.PHONY: get-ai-summary
get-ai-summary:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "finding_id": 1001, "lang":"ja"}' \
		$(CORE_API_ADDR) core.finding.FindingService.GetAISummary

.PHONY: get-ai-summary-stream
get-ai-summary-stream:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "finding_id": 1001, "lang":"ja"}' \
		$(CORE_API_ADDR) core.finding.FindingService.GetAISummaryStream

.PHONY: clean-old-resource
clean-old-resource:
	$(GRPCURL) \
		-plaintext \
		-d '{}' \
		$(CORE_API_ADDR) core.finding.FindingService.CleanOldResource

.PHONY: list-report-service
list-report-service:
	$(GRPCURL) -plaintext $(CORE_API_ADDR) list core.report.ReportService

.PHONY: get-report-finding
get-report-finding:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001}' \
		$(CORE_API_ADDR) core.report.ReportService.GetReportFinding

.PHONY: get-report-finding-all
get-report-finding-all:
	$(GRPCURL) \
		-plaintext \
		-d '{}' \
		$(CORE_API_ADDR) core.report.ReportService.GetReportFindingAll

.PHONY: collect-report-finding
collect-report-finding:
	$(GRPCURL) \
		-plaintext \
		-d '{}' \
		$(CORE_API_ADDR) core.report.ReportService.CollectReportFinding

.PHONY: purge-report-finding
purge-report-finding:
	$(GRPCURL) \
		-plaintext \
		-d '{}' \
		$(CORE_API_ADDR) core.report.ReportService.PurgeReportFinding

.PHONY: get-report
get-report:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "report_id":1001}' \
		$(CORE_API_ADDR) core.report.ReportService.GetReport

.PHONY: list-report
list-report:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001}' \
		$(CORE_API_ADDR) core.report.ReportService.ListReport

.PHONY: put-report
put-report:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "name":"report-name", "type":"Markdown", "status":"OK", "content":"# title"}' \
		$(CORE_API_ADDR) core.report.ReportService.PutReport

.PHONY: list-iam-service
list-iam-service:
	$(GRPCURL) -plaintext $(CORE_API_ADDR) list core.iam.IAMService

.PHONY: is-authorized
is-authorized:
	$(GRPCURL) \
		-plaintext \
		-d '{"user_id":1001, "project_id":1001, "action_name":"finding/GetFinding", "resource_name":"aws:guardduty/s3-bucket-name"}' \
		$(CORE_API_ADDR) core.iam.IAMService.IsAuthorized

.PHONY: is-authorized-token
is-authorized-token:
	$(GRPCURL) \
		-plaintext \
		-d '{"access_token_id":1001, "project_id":1001, "action_name":"finding/GetFinding", "resource_name":"aws:guardduty/s3-bucket-name"}' \
		$(CORE_API_ADDR) core.iam.IAMService.IsAuthorizedToken

.PHONY: is-admin
is-admin:
	$(GRPCURL) \
		-plaintext \
		-d '{"user_id":1001}' \
		$(CORE_API_ADDR) core.iam.IAMService.IsAdmin

.PHONY: list-user
list-user:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "activated":"true", "name":"john"}' \
		$(CORE_API_ADDR) core.iam.IAMService.ListUser

.PHONY: get-user
get-user:
	$(GRPCURL) \
		-plaintext \
		-d '{"user_id":1001, "sub":"alice"}' \
		$(CORE_API_ADDR) core.iam.IAMService.GetUser

.PHONY: put-user
put-user:
	$(GRPCURL) \
		-plaintext \
		-d '{"user": {"sub":"sub", "name":"name", "user_idp_key":"user_idp_key"}}' \
		$(CORE_API_ADDR) core.iam.IAMService.PutUser

.PHONY: list-role
list-role:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "name":"admin-role"}' \
		$(CORE_API_ADDR) core.iam.IAMService.ListRole

.PHONY: get-role
get-role:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "role_id":1001}' \
		$(CORE_API_ADDR) core.iam.IAMService.GetRole

.PHONY: put-role
put-role:
	$(GRPCURL) \
		-plaintext \
		-d '{"role":{"name":"test-role", "project_id":1001}}' \
		$(CORE_API_ADDR) core.iam.IAMService.PutRole

.PHONY: delete-role
delete-role:
	$(GRPCURL) \
		-plaintext \
		-d '{"role_id":1004, "project_id":1001}' \
		$(CORE_API_ADDR) core.iam.IAMService.DeleteRole

.PHONY: attach-role
attach-role:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "role_id":1005, "user_id":1001}' \
		$(CORE_API_ADDR) core.iam.IAMService.AttachRole

.PHONY: detach-role
detach-role:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "role_id":1005, "user_id":1001}' \
		$(CORE_API_ADDR) core.iam.IAMService.DetachRole

.PHONY: list-policy
list-policy:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "name":"admin-policy"}' \
		$(CORE_API_ADDR) core.iam.IAMService.ListPolicy

.PHONY: get-policy
get-policy:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "policy_id":1001}' \
		$(CORE_API_ADDR) core.iam.IAMService.GetPolicy

.PHONY: put-policy
put-policy:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "policy":{"name":"test-policy", "project_id":1001, "action_ptn":".*", "resource_ptn":".*"}}' \
		$(CORE_API_ADDR) core.iam.IAMService.PutPolicy

.PHONY: delete-policy
delete-policy:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "policy_id":1004}' \
		$(CORE_API_ADDR) core.iam.IAMService.DeletePolicy

.PHONY: attach-policy
attach-policy:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "role_id":1001, "policy_id":1005}' \
		$(CORE_API_ADDR) core.iam.IAMService.AttachPolicy
 
.PHONY: detach-policy
detach-policy:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "role_id":1001, "policy_id":1005}' \
		$(CORE_API_ADDR) core.iam.IAMService.DetachPolicy

.PHONY: list-access-token
list-access-token:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001}' \
		$(CORE_API_ADDR) core.iam.IAMService.ListAccessToken

.PHONY: authenticate-access-token
authenticate-access-token:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "access_token_id": 1001,"plain_text_token":"test-token"}' \
		$(CORE_API_ADDR) core.iam.IAMService.AuthenticateAccessToken

.PHONY: put-access-token
put-access-token:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "access_token":{"plain_text_token":"test-token", "name":"test", "project_id":1001, "description":"description", "expired_at":2628676885, "last_updated_uesr_id":1001}}' \
		$(CORE_API_ADDR) core.iam.IAMService.PutAccessToken

.PHONY: delete-access-token
delete-access-token:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "access_token_id": 1001}' \
		$(CORE_API_ADDR) core.iam.IAMService.DeleteAccessToken

.PHONY: attach-access-token-role
attach-access-token-role:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "role_id":1002, "access_token_id": 1001}' \
		$(CORE_API_ADDR) core.iam.IAMService.AttachAccessTokenRole

.PHONY: detach-access-token-role
detach-access-token-role:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "role_id":1002, "access_token_id": 1001}' \
		$(CORE_API_ADDR) core.iam.IAMService.DetachAccessTokenRole

.PHONY: analyze-access-token-expiration
analyze-access-token-expiration:
	$(GRPCURL) \
		-plaintext \
		$(CORE_API_ADDR) core.iam.IAMService.AnalyzeTokenExpiration

.PHONY: list-user-reserved
list-user-reserved:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001}' \
		$(CORE_API_ADDR) core.iam.IAMService.ListUserReserved

.PHONY: put-user-reserved
put-user-reserved:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "user_reserved": {"role_id": 1001, "user_idp_key": "reserved_user"}}' \
		$(CORE_API_ADDR) core.iam.IAMService.PutUserReserved

.PHONY: delete-user-reserved
delete-user-reserved:
	$(GRPCURL) \
		-plaintext \
		-d '{"project_id":1001, "reserved_id": 1006}' \
		$(CORE_API_ADDR) core.iam.IAMService.DeleteUserReserved

.PHONY: chat-ai
chat-ai:
	$(GRPCURL) \
		-plaintext \
		-d '{"question":"What mountain is the highest in the world?", "chat_history": [{"role":1, "content":"hello!"}, {"role":2, "content":"Hi, I am a chatbot."}]}' \
		$(CORE_API_ADDR) core.ai.AIService.ChatAI

.PHONY: generate-report
generate-report:
	$(GRPCURL) \
		-plaintext \
		-d '{"prompt":"AWSのFindingレポートを作成してください。データソースごとの解析もお願いします。", "project_id":1001, "name":"report-name"}' \
		$(CORE_API_ADDR) core.ai.AIService.GenerateReport

.PHONY: generate-report2
generate-report2:
	$(GRPCURL) \
		-plaintext \
		-d '{"prompt":"google:sccのFindingを分析して", "project_id":1001}' \
		$(CORE_API_ADDR) core.ai.AIService.GenerateReport

FAKE:
