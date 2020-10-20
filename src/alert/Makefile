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
	grpcurl -plaintext localhost:8004 list core.alert.AlertService

.PHONY: list-alert
list-alert:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "description": "", "status":["ACTIVE","PENDING"], "severity":["high","medium","low"],"from_at":1560000000, "to_at":1660000000}' \
		localhost:8004 core.alert.AlertService.ListAlert

.PHONY: get-alert
get-alert:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_id":1001}' \
		localhost:8004 core.alert.AlertService.GetAlert

.PHONY: put-alert
put-alert:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert":{"alert_condition_id": 1003, "description":"test_put_alert", "severity": "low", "status": "ACTIVE", "project_id":1001}}' \
		localhost:8004 core.alert.AlertService.PutAlert

.PHONY: delete-alert
delete-alert:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_id":1003}' \
		localhost:8004 core.alert.AlertService.DeleteAlert

.PHONY: list-alert_history
list-alert_history:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001,"alert_id":1001, "history_type": ["created","deleted"], "severity":["high","medium"], "from_at":1560000000, "to_at":1660000000}' \
		localhost:8004 core.alert.AlertService.ListAlertHistory

.PHONY: get-alert_history
get-alert_history:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_history_id":1001}' \
		localhost:8004 core.alert.AlertService.GetAlertHistory

.PHONY: put-alert_history
put-alert_history:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_history":{"alert_id": 1003, "history_type":"created","description":"test_put_alert_history", "severity": "low", "project_id":1001}}' \
		localhost:8004 core.alert.AlertService.PutAlertHistory

.PHONY: delete-alert_history
delete-alert_history:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_history_id":1003}' \
		localhost:8004 core.alert.AlertService.DeleteAlertHistory

.PHONY: list-rel_alert_finding
list-rel_alert_finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001,"alert_id":1001, "finding_id":1001, "from_at":1560000000, "to_at":1660000000}' \
		localhost:8004 core.alert.AlertService.ListRelAlertFinding

.PHONY: get-rel_alert_finding
get-rel_alert_finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001,"alert_id":1001, "finding_id":1001}' \
		localhost:8004 core.alert.AlertService.GetRelAlertFinding

.PHONY: put-rel_alert_finding
put-rel_alert_finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "rel_alert_finding":{"project_id":1001,"alert_id":1003, "finding_id":1001}}' \
		localhost:8004 core.alert.AlertService.PutRelAlertFinding

.PHONY: delete-rel_alert_finding
delete-rel_alert_finding:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001,"alert_id":1003, "finding_id":1001}' \
		localhost:8004 core.alert.AlertService.DeleteRelAlertFinding

.PHONY: list-alert_condition
list-alert_condition:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "severity":["high","medium"], "enabled": true, "from_at":1560000000, "to_at":1660000000}' \
		localhost:8004 core.alert.AlertService.ListAlertCondition

.PHONY: get-alert_condition
get-alert_condition:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1001}' \
		localhost:8004 core.alert.AlertService.GetAlertCondition

.PHONY: put-alert_condition
put-alert_condition:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition":{"enabled": true, "description":"test_put_alert_condition", "severity": "low", "and_or": "or", "project_id":1001}}' \
		localhost:8004 core.alert.AlertService.PutAlertCondition

.PHONY: delete-alert_condition
delete-alert_condition:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1003}' \
		localhost:8004 core.alert.AlertService.DeleteAlertCondition

.PHONY: list-alert_rule
list-alert_rule:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "from_score":0.0, "to_score":1.0, "from_at":1560000000, "to_at":1660000000}' \
		localhost:8004 core.alert.AlertService.ListAlertRule

.PHONY: get-alert_rule
get-alert_rule:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_rule_id":1001}' \
		localhost:8004 core.alert.AlertService.GetAlertRule

.PHONY: put-alert_rule
put-alert_rule:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_rule":{"name": "test_put_alert_rule", "score": 0.1, "project_id":1001}}' \
		localhost:8004 core.alert.AlertService.PutAlertRule

.PHONY: delete-alert_rule
delete-alert_rule:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_rule_id":1003}' \
		localhost:8004 core.alert.AlertService.DeleteAlertRule

.PHONY: list-alert_cond_rule
list-alert_cond_rule:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1001, "alert_rule_id":1001, "from_at":1560000000, "to_at":1660000000}' \
		localhost:8004 core.alert.AlertService.ListAlertCondRule

.PHONY: get-alert_cond_rule
get-alert_cond_rule:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1001, "alert_rule_id":1001}' \
		localhost:8004 core.alert.AlertService.GetAlertCondRule

.PHONY: put-alert_cond_rule
put-alert_cond_rule:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_cond_rule":{"project_id":1001, "alert_condition_id":1003, "alert_rule_id":1001}}' \
		localhost:8004 core.alert.AlertService.PutAlertCondRule

.PHONY: delete-alert_cond_rule
delete-alert_cond_rule:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1003, "alert_rule_id":1001}' \
		localhost:8004 core.alert.AlertService.DeleteAlertCondRule

.PHONY: list-notification
list-notification:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "type":"slack", "from_at":1560000000, "to_at":1660000000}' \
		localhost:8004 core.alert.AlertService.ListNotification

.PHONY: get-notification
get-notification:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "notification_id":1001}' \
		localhost:8004 core.alert.AlertService.GetNotification

.PHONY: put-notification
put-notification:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001,"notification":{"project_id":1001, "name":"test_notification","type":"slack", "notify_setting":"{\"hoge\":\"fuga\"}"}}' \
		localhost:8004 core.alert.AlertService.PutNotification

.PHONY: delete-notification
delete-notification:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "notification_id":1003}' \
		localhost:8004 core.alert.AlertService.DeleteNotification

.PHONY: analyze-alert
analyze-alert:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001}' \
		localhost:8004 core.alert.AlertService.AnalyzeAlert

.PHONY: list-alert_cond_notification
list-alert_cond_notification:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1001, "notification_id":1001, "from_at":1560000000, "to_at":1660000000}' \
		localhost:8004 core.alert.AlertService.ListAlertCondNotification

.PHONY: get-alert_cond_notification
get-alert_cond_notification:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1001, "notification_id":1001}' \
		localhost:8004 core.alert.AlertService.GetAlertCondNotification

.PHONY: put-alert_cond_notification
put-alert_cond_notification:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_cond_notification":{"project_id":1001, "alert_condition_id":1003, "notification_id":1001,"cache_second":3600,"notified_at":1560000000}}' \
		localhost:8004 core.alert.AlertService.PutAlertCondNotification

.PHONY: delete-alert_cond_notification
delete-alert_cond_notification:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "alert_condition_id":1003, "notification_id":1001}' \
		localhost:8004 core.alert.AlertService.DeleteAlertCondNotification
