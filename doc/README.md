# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [alert/entity.proto](#alert/entity.proto)
    - [Alert](#core.alert.Alert)
    - [AlertCondNotification](#core.alert.AlertCondNotification)
    - [AlertCondNotificationForUpsert](#core.alert.AlertCondNotificationForUpsert)
    - [AlertCondRule](#core.alert.AlertCondRule)
    - [AlertCondRuleForUpsert](#core.alert.AlertCondRuleForUpsert)
    - [AlertCondition](#core.alert.AlertCondition)
    - [AlertConditionForUpsert](#core.alert.AlertConditionForUpsert)
    - [AlertForUpsert](#core.alert.AlertForUpsert)
    - [AlertHistory](#core.alert.AlertHistory)
    - [AlertHistoryForUpsert](#core.alert.AlertHistoryForUpsert)
    - [AlertRule](#core.alert.AlertRule)
    - [AlertRuleForUpsert](#core.alert.AlertRuleForUpsert)
    - [Notification](#core.alert.Notification)
    - [NotificationForUpsert](#core.alert.NotificationForUpsert)
    - [RelAlertFinding](#core.alert.RelAlertFinding)
    - [RelAlertFindingForUpsert](#core.alert.RelAlertFindingForUpsert)
  
    - [Status](#core.alert.Status)
  
- [alert/service.proto](#alert/service.proto)
    - [AnalyzeAlertRequest](#core.alert.AnalyzeAlertRequest)
    - [DeleteAlertCondNotificationRequest](#core.alert.DeleteAlertCondNotificationRequest)
    - [DeleteAlertCondRuleRequest](#core.alert.DeleteAlertCondRuleRequest)
    - [DeleteAlertConditionRequest](#core.alert.DeleteAlertConditionRequest)
    - [DeleteAlertHistoryRequest](#core.alert.DeleteAlertHistoryRequest)
    - [DeleteAlertRequest](#core.alert.DeleteAlertRequest)
    - [DeleteAlertRuleRequest](#core.alert.DeleteAlertRuleRequest)
    - [DeleteNotificationRequest](#core.alert.DeleteNotificationRequest)
    - [DeleteRelAlertFindingRequest](#core.alert.DeleteRelAlertFindingRequest)
    - [GetAlertCondNotificationRequest](#core.alert.GetAlertCondNotificationRequest)
    - [GetAlertCondNotificationResponse](#core.alert.GetAlertCondNotificationResponse)
    - [GetAlertCondRuleRequest](#core.alert.GetAlertCondRuleRequest)
    - [GetAlertCondRuleResponse](#core.alert.GetAlertCondRuleResponse)
    - [GetAlertConditionRequest](#core.alert.GetAlertConditionRequest)
    - [GetAlertConditionResponse](#core.alert.GetAlertConditionResponse)
    - [GetAlertHistoryRequest](#core.alert.GetAlertHistoryRequest)
    - [GetAlertHistoryResponse](#core.alert.GetAlertHistoryResponse)
    - [GetAlertRequest](#core.alert.GetAlertRequest)
    - [GetAlertResponse](#core.alert.GetAlertResponse)
    - [GetAlertRuleRequest](#core.alert.GetAlertRuleRequest)
    - [GetAlertRuleResponse](#core.alert.GetAlertRuleResponse)
    - [GetNotificationRequest](#core.alert.GetNotificationRequest)
    - [GetNotificationResponse](#core.alert.GetNotificationResponse)
    - [GetRelAlertFindingRequest](#core.alert.GetRelAlertFindingRequest)
    - [GetRelAlertFindingResponse](#core.alert.GetRelAlertFindingResponse)
    - [ListAlertCondNotificationRequest](#core.alert.ListAlertCondNotificationRequest)
    - [ListAlertCondNotificationResponse](#core.alert.ListAlertCondNotificationResponse)
    - [ListAlertCondRuleRequest](#core.alert.ListAlertCondRuleRequest)
    - [ListAlertCondRuleResponse](#core.alert.ListAlertCondRuleResponse)
    - [ListAlertConditionRequest](#core.alert.ListAlertConditionRequest)
    - [ListAlertConditionResponse](#core.alert.ListAlertConditionResponse)
    - [ListAlertHistoryRequest](#core.alert.ListAlertHistoryRequest)
    - [ListAlertHistoryResponse](#core.alert.ListAlertHistoryResponse)
    - [ListAlertRequest](#core.alert.ListAlertRequest)
    - [ListAlertResponse](#core.alert.ListAlertResponse)
    - [ListAlertRuleRequest](#core.alert.ListAlertRuleRequest)
    - [ListAlertRuleResponse](#core.alert.ListAlertRuleResponse)
    - [ListNotificationRequest](#core.alert.ListNotificationRequest)
    - [ListNotificationResponse](#core.alert.ListNotificationResponse)
    - [ListRelAlertFindingRequest](#core.alert.ListRelAlertFindingRequest)
    - [ListRelAlertFindingResponse](#core.alert.ListRelAlertFindingResponse)
    - [PutAlertCondNotificationRequest](#core.alert.PutAlertCondNotificationRequest)
    - [PutAlertCondNotificationResponse](#core.alert.PutAlertCondNotificationResponse)
    - [PutAlertCondRuleRequest](#core.alert.PutAlertCondRuleRequest)
    - [PutAlertCondRuleResponse](#core.alert.PutAlertCondRuleResponse)
    - [PutAlertConditionRequest](#core.alert.PutAlertConditionRequest)
    - [PutAlertConditionResponse](#core.alert.PutAlertConditionResponse)
    - [PutAlertHistoryRequest](#core.alert.PutAlertHistoryRequest)
    - [PutAlertHistoryResponse](#core.alert.PutAlertHistoryResponse)
    - [PutAlertRequest](#core.alert.PutAlertRequest)
    - [PutAlertResponse](#core.alert.PutAlertResponse)
    - [PutAlertRuleRequest](#core.alert.PutAlertRuleRequest)
    - [PutAlertRuleResponse](#core.alert.PutAlertRuleResponse)
    - [PutNotificationRequest](#core.alert.PutNotificationRequest)
    - [PutNotificationResponse](#core.alert.PutNotificationResponse)
    - [PutRelAlertFindingRequest](#core.alert.PutRelAlertFindingRequest)
    - [PutRelAlertFindingResponse](#core.alert.PutRelAlertFindingResponse)
  
    - [AlertService](#core.alert.AlertService)
  
- [finding/entity.proto](#finding/entity.proto)
    - [Finding](#core.finding.Finding)
    - [FindingForUpsert](#core.finding.FindingForUpsert)
    - [FindingTag](#core.finding.FindingTag)
    - [FindingTagForUpsert](#core.finding.FindingTagForUpsert)
    - [PendFinding](#core.finding.PendFinding)
    - [PendFindingForUpsert](#core.finding.PendFindingForUpsert)
    - [Resource](#core.finding.Resource)
    - [ResourceForUpsert](#core.finding.ResourceForUpsert)
    - [ResourceTag](#core.finding.ResourceTag)
    - [ResourceTagForUpsert](#core.finding.ResourceTagForUpsert)
  
- [finding/service.proto](#finding/service.proto)
    - [DeleteFindingRequest](#core.finding.DeleteFindingRequest)
    - [DeletePendFindingRequest](#core.finding.DeletePendFindingRequest)
    - [DeleteResourceRequest](#core.finding.DeleteResourceRequest)
    - [GetFindingRequest](#core.finding.GetFindingRequest)
    - [GetFindingResponse](#core.finding.GetFindingResponse)
    - [GetPendFindingRequest](#core.finding.GetPendFindingRequest)
    - [GetPendFindingResponse](#core.finding.GetPendFindingResponse)
    - [GetResourceRequest](#core.finding.GetResourceRequest)
    - [GetResourceResponse](#core.finding.GetResourceResponse)
    - [ListFindingRequest](#core.finding.ListFindingRequest)
    - [ListFindingResponse](#core.finding.ListFindingResponse)
    - [ListFindingTagNameRequest](#core.finding.ListFindingTagNameRequest)
    - [ListFindingTagNameResponse](#core.finding.ListFindingTagNameResponse)
    - [ListFindingTagRequest](#core.finding.ListFindingTagRequest)
    - [ListFindingTagResponse](#core.finding.ListFindingTagResponse)
    - [ListResourceRequest](#core.finding.ListResourceRequest)
    - [ListResourceResponse](#core.finding.ListResourceResponse)
    - [ListResourceTagNameRequest](#core.finding.ListResourceTagNameRequest)
    - [ListResourceTagNameResponse](#core.finding.ListResourceTagNameResponse)
    - [ListResourceTagRequest](#core.finding.ListResourceTagRequest)
    - [ListResourceTagResponse](#core.finding.ListResourceTagResponse)
    - [PutFindingRequest](#core.finding.PutFindingRequest)
    - [PutFindingResponse](#core.finding.PutFindingResponse)
    - [PutPendFindingRequest](#core.finding.PutPendFindingRequest)
    - [PutPendFindingResponse](#core.finding.PutPendFindingResponse)
    - [PutResourceRequest](#core.finding.PutResourceRequest)
    - [PutResourceResponse](#core.finding.PutResourceResponse)
    - [TagFindingRequest](#core.finding.TagFindingRequest)
    - [TagFindingResponse](#core.finding.TagFindingResponse)
    - [TagResourceRequest](#core.finding.TagResourceRequest)
    - [TagResourceResponse](#core.finding.TagResourceResponse)
    - [UntagFindingRequest](#core.finding.UntagFindingRequest)
    - [UntagResourceRequest](#core.finding.UntagResourceRequest)
  
    - [FindingService](#core.finding.FindingService)
  
- [iam/entity.proto](#iam/entity.proto)
    - [Policy](#core.iam.Policy)
    - [PolicyForUpsert](#core.iam.PolicyForUpsert)
    - [Role](#core.iam.Role)
    - [RoleForUpsert](#core.iam.RoleForUpsert)
    - [RolePolicy](#core.iam.RolePolicy)
    - [User](#core.iam.User)
    - [UserForUpsert](#core.iam.UserForUpsert)
    - [UserRole](#core.iam.UserRole)
  
- [iam/policy.proto](#iam/policy.proto)
    - [AttachPolicyRequest](#core.iam.AttachPolicyRequest)
    - [AttachPolicyResponse](#core.iam.AttachPolicyResponse)
    - [DeletePolicyRequest](#core.iam.DeletePolicyRequest)
    - [DetachPolicyRequest](#core.iam.DetachPolicyRequest)
    - [GetPolicyRequest](#core.iam.GetPolicyRequest)
    - [GetPolicyResponse](#core.iam.GetPolicyResponse)
    - [ListPolicyRequest](#core.iam.ListPolicyRequest)
    - [ListPolicyResponse](#core.iam.ListPolicyResponse)
    - [PutPolicyRequest](#core.iam.PutPolicyRequest)
    - [PutPolicyResponse](#core.iam.PutPolicyResponse)
  
- [iam/role.proto](#iam/role.proto)
    - [AttachRoleRequest](#core.iam.AttachRoleRequest)
    - [AttachRoleResponse](#core.iam.AttachRoleResponse)
    - [DeleteRoleRequest](#core.iam.DeleteRoleRequest)
    - [DetachRoleRequest](#core.iam.DetachRoleRequest)
    - [GetRoleRequest](#core.iam.GetRoleRequest)
    - [GetRoleResponse](#core.iam.GetRoleResponse)
    - [ListRoleRequest](#core.iam.ListRoleRequest)
    - [ListRoleResponse](#core.iam.ListRoleResponse)
    - [PutRoleRequest](#core.iam.PutRoleRequest)
    - [PutRoleResponse](#core.iam.PutRoleResponse)
  
- [iam/service.proto](#iam/service.proto)
    - [IsAdminRequest](#core.iam.IsAdminRequest)
    - [IsAdminResponse](#core.iam.IsAdminResponse)
    - [IsAuthorizedRequest](#core.iam.IsAuthorizedRequest)
    - [IsAuthorizedResponse](#core.iam.IsAuthorizedResponse)
  
    - [IAMService](#core.iam.IAMService)
  
- [iam/user.proto](#iam/user.proto)
    - [GetUserRequest](#core.iam.GetUserRequest)
    - [GetUserResponse](#core.iam.GetUserResponse)
    - [ListUserRequest](#core.iam.ListUserRequest)
    - [ListUserResponse](#core.iam.ListUserResponse)
    - [PutUserRequest](#core.iam.PutUserRequest)
    - [PutUserResponse](#core.iam.PutUserResponse)
  
- [project/entity.proto](#project/entity.proto)
    - [Project](#core.project.Project)
  
- [project/service.proto](#project/service.proto)
    - [CreateProjectRequest](#core.project.CreateProjectRequest)
    - [CreateProjectResponse](#core.project.CreateProjectResponse)
    - [DeleteProjectRequest](#core.project.DeleteProjectRequest)
    - [ListProjectRequest](#core.project.ListProjectRequest)
    - [ListProjectResponse](#core.project.ListProjectResponse)
    - [UpdateProjectRequest](#core.project.UpdateProjectRequest)
    - [UpdateProjectResponse](#core.project.UpdateProjectResponse)
  
    - [ProjectService](#core.project.ProjectService)
  
- [report/entity.proto](#report/entity.proto)
    - [ReportFinding](#core.report.ReportFinding)
  
- [report/service.proto](#report/service.proto)
    - [GetReportFindingAllRequest](#core.report.GetReportFindingAllRequest)
    - [GetReportFindingAllResponse](#core.report.GetReportFindingAllResponse)
    - [GetReportFindingRequest](#core.report.GetReportFindingRequest)
    - [GetReportFindingResponse](#core.report.GetReportFindingResponse)
  
    - [ReportService](#core.report.ReportService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="alert/entity.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## alert/entity.proto



<a name="core.alert.Alert"></a>

### Alert
Alert


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_id | [uint32](#uint32) |  |  |
| alert_condition_id | [uint32](#uint32) |  |  |
| description | [string](#string) |  |  |
| severity | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| status | [Status](#core.alert.Status) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.alert.AlertCondNotification"></a>

### AlertCondNotification
AlertCondNotification


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_condition_id | [uint32](#uint32) |  |  |
| notification_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| cache_second | [uint32](#uint32) |  |  |
| notified_at | [int64](#int64) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.alert.AlertCondNotificationForUpsert"></a>

### AlertCondNotificationForUpsert
AlertCondNotificationForUpsert


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_condition_id | [uint32](#uint32) |  |  |
| notification_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| cache_second | [uint32](#uint32) |  |  |
| notified_at | [int64](#int64) |  |  |






<a name="core.alert.AlertCondRule"></a>

### AlertCondRule
AlertCondRule


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_condition_id | [uint32](#uint32) |  |  |
| alert_rule_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.alert.AlertCondRuleForUpsert"></a>

### AlertCondRuleForUpsert
AlertCondRuleForUpsert


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_condition_id | [uint32](#uint32) |  |  |
| alert_rule_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |






<a name="core.alert.AlertCondition"></a>

### AlertCondition
AlertCondition


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_condition_id | [uint32](#uint32) |  |  |
| description | [string](#string) |  |  |
| severity | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| and_or | [string](#string) |  |  |
| enabled | [bool](#bool) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.alert.AlertConditionForUpsert"></a>

### AlertConditionForUpsert
AlertConditionForUpsert


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_condition_id | [uint32](#uint32) |  |  |
| description | [string](#string) |  |  |
| severity | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| and_or | [string](#string) |  |  |
| enabled | [bool](#bool) |  |  |






<a name="core.alert.AlertForUpsert"></a>

### AlertForUpsert
AlertForUpsert


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_id | [uint32](#uint32) |  |  |
| alert_condition_id | [uint32](#uint32) |  |  |
| description | [string](#string) |  |  |
| severity | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| status | [Status](#core.alert.Status) |  |  |






<a name="core.alert.AlertHistory"></a>

### AlertHistory
AlertHistory


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_history_id | [uint32](#uint32) |  |  |
| history_type | [string](#string) |  |  |
| alert_id | [uint32](#uint32) |  |  |
| description | [string](#string) |  |  |
| severity | [string](#string) |  |  |
| finding_history | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.alert.AlertHistoryForUpsert"></a>

### AlertHistoryForUpsert
AlertHistoryForUpsert


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_history_id | [uint32](#uint32) |  |  |
| history_type | [string](#string) |  |  |
| alert_id | [uint32](#uint32) |  |  |
| description | [string](#string) |  |  |
| severity | [string](#string) |  |  |
| finding_history | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |






<a name="core.alert.AlertRule"></a>

### AlertRule
AlertRule


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_rule_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| score | [float](#float) |  |  |
| resource_name | [string](#string) |  |  |
| tag | [string](#string) |  |  |
| finding_cnt | [uint32](#uint32) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.alert.AlertRuleForUpsert"></a>

### AlertRuleForUpsert
AlertRuleForUpsert


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_rule_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| score | [float](#float) |  |  |
| resource_name | [string](#string) |  |  |
| tag | [string](#string) |  |  |
| finding_cnt | [uint32](#uint32) |  |  |






<a name="core.alert.Notification"></a>

### Notification
Notification


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| notification_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| type | [string](#string) |  |  |
| notify_setting | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.alert.NotificationForUpsert"></a>

### NotificationForUpsert
NotificationForUpsert


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| notification_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| type | [string](#string) |  |  |
| notify_setting | [string](#string) |  |  |






<a name="core.alert.RelAlertFinding"></a>

### RelAlertFinding
RelAlertFinding


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_id | [uint32](#uint32) |  |  |
| finding_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.alert.RelAlertFindingForUpsert"></a>

### RelAlertFindingForUpsert
RelAlertFindingForUpsert


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_id | [uint32](#uint32) |  |  |
| finding_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |





 


<a name="core.alert.Status"></a>

### Status
Status

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| ACTIVE | 1 |  |
| PENDING | 2 |  |
| DEACTIVE | 3 |  |


 

 

 



<a name="alert/service.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## alert/service.proto



<a name="core.alert.AnalyzeAlertRequest"></a>

### AnalyzeAlertRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_condition_id | [uint32](#uint32) | repeated |  |






<a name="core.alert.DeleteAlertCondNotificationRequest"></a>

### DeleteAlertCondNotificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_condition_id | [uint32](#uint32) |  |  |
| notification_id | [uint32](#uint32) |  |  |






<a name="core.alert.DeleteAlertCondRuleRequest"></a>

### DeleteAlertCondRuleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_condition_id | [uint32](#uint32) |  |  |
| alert_rule_id | [uint32](#uint32) |  |  |






<a name="core.alert.DeleteAlertConditionRequest"></a>

### DeleteAlertConditionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_condition_id | [uint32](#uint32) |  |  |






<a name="core.alert.DeleteAlertHistoryRequest"></a>

### DeleteAlertHistoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_history_id | [uint32](#uint32) |  |  |






<a name="core.alert.DeleteAlertRequest"></a>

### DeleteAlertRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_id | [uint32](#uint32) |  |  |






<a name="core.alert.DeleteAlertRuleRequest"></a>

### DeleteAlertRuleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_rule_id | [uint32](#uint32) |  |  |






<a name="core.alert.DeleteNotificationRequest"></a>

### DeleteNotificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| notification_id | [uint32](#uint32) |  |  |






<a name="core.alert.DeleteRelAlertFindingRequest"></a>

### DeleteRelAlertFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_id | [uint32](#uint32) |  |  |
| finding_id | [uint32](#uint32) |  |  |






<a name="core.alert.GetAlertCondNotificationRequest"></a>

### GetAlertCondNotificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_condition_id | [uint32](#uint32) |  |  |
| notification_id | [uint32](#uint32) |  |  |






<a name="core.alert.GetAlertCondNotificationResponse"></a>

### GetAlertCondNotificationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_cond_notification | [AlertCondNotification](#core.alert.AlertCondNotification) |  |  |






<a name="core.alert.GetAlertCondRuleRequest"></a>

### GetAlertCondRuleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_condition_id | [uint32](#uint32) |  |  |
| alert_rule_id | [uint32](#uint32) |  |  |






<a name="core.alert.GetAlertCondRuleResponse"></a>

### GetAlertCondRuleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_cond_rule | [AlertCondRule](#core.alert.AlertCondRule) |  |  |






<a name="core.alert.GetAlertConditionRequest"></a>

### GetAlertConditionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_condition_id | [uint32](#uint32) |  |  |






<a name="core.alert.GetAlertConditionResponse"></a>

### GetAlertConditionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_condition | [AlertCondition](#core.alert.AlertCondition) |  |  |






<a name="core.alert.GetAlertHistoryRequest"></a>

### GetAlertHistoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_history_id | [uint32](#uint32) |  |  |






<a name="core.alert.GetAlertHistoryResponse"></a>

### GetAlertHistoryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_history | [AlertHistory](#core.alert.AlertHistory) |  |  |






<a name="core.alert.GetAlertRequest"></a>

### GetAlertRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_id | [uint32](#uint32) |  |  |






<a name="core.alert.GetAlertResponse"></a>

### GetAlertResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert | [Alert](#core.alert.Alert) |  |  |






<a name="core.alert.GetAlertRuleRequest"></a>

### GetAlertRuleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_rule_id | [uint32](#uint32) |  |  |






<a name="core.alert.GetAlertRuleResponse"></a>

### GetAlertRuleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_rule | [AlertRule](#core.alert.AlertRule) |  |  |






<a name="core.alert.GetNotificationRequest"></a>

### GetNotificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| notification_id | [uint32](#uint32) |  |  |






<a name="core.alert.GetNotificationResponse"></a>

### GetNotificationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| notification | [Notification](#core.alert.Notification) |  |  |






<a name="core.alert.GetRelAlertFindingRequest"></a>

### GetRelAlertFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_id | [uint32](#uint32) |  |  |
| finding_id | [uint32](#uint32) |  |  |






<a name="core.alert.GetRelAlertFindingResponse"></a>

### GetRelAlertFindingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rel_alert_finding | [RelAlertFinding](#core.alert.RelAlertFinding) |  |  |






<a name="core.alert.ListAlertCondNotificationRequest"></a>

### ListAlertCondNotificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_condition_id | [uint32](#uint32) |  |  |
| notification_id | [uint32](#uint32) |  |  |
| from_at | [int64](#int64) |  |  |
| to_at | [int64](#int64) |  |  |






<a name="core.alert.ListAlertCondNotificationResponse"></a>

### ListAlertCondNotificationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_cond_notification | [AlertCondNotification](#core.alert.AlertCondNotification) | repeated |  |






<a name="core.alert.ListAlertCondRuleRequest"></a>

### ListAlertCondRuleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_condition_id | [uint32](#uint32) |  |  |
| alert_rule_id | [uint32](#uint32) |  |  |
| from_at | [int64](#int64) |  |  |
| to_at | [int64](#int64) |  |  |






<a name="core.alert.ListAlertCondRuleResponse"></a>

### ListAlertCondRuleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_cond_rule | [AlertCondRule](#core.alert.AlertCondRule) | repeated |  |






<a name="core.alert.ListAlertConditionRequest"></a>

### ListAlertConditionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| severity | [string](#string) | repeated |  |
| enabled | [bool](#bool) |  |  |
| from_at | [int64](#int64) |  |  |
| to_at | [int64](#int64) |  |  |






<a name="core.alert.ListAlertConditionResponse"></a>

### ListAlertConditionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_condition | [AlertCondition](#core.alert.AlertCondition) | repeated |  |






<a name="core.alert.ListAlertHistoryRequest"></a>

### ListAlertHistoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_id | [uint32](#uint32) |  |  |
| history_type | [string](#string) | repeated |  |
| severity | [string](#string) | repeated |  |
| from_at | [int64](#int64) |  |  |
| to_at | [int64](#int64) |  |  |






<a name="core.alert.ListAlertHistoryResponse"></a>

### ListAlertHistoryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_history | [AlertHistory](#core.alert.AlertHistory) | repeated |  |






<a name="core.alert.ListAlertRequest"></a>

### ListAlertRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| status | [Status](#core.alert.Status) | repeated |  |
| severity | [string](#string) | repeated |  |
| description | [string](#string) |  |  |
| from_at | [int64](#int64) |  |  |
| to_at | [int64](#int64) |  |  |






<a name="core.alert.ListAlertResponse"></a>

### ListAlertResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert | [Alert](#core.alert.Alert) | repeated |  |






<a name="core.alert.ListAlertRuleRequest"></a>

### ListAlertRuleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| from_score | [float](#float) |  |  |
| to_score | [float](#float) |  |  |
| from_at | [int64](#int64) |  |  |
| to_at | [int64](#int64) |  |  |






<a name="core.alert.ListAlertRuleResponse"></a>

### ListAlertRuleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_rule | [AlertRule](#core.alert.AlertRule) | repeated |  |






<a name="core.alert.ListNotificationRequest"></a>

### ListNotificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| type | [string](#string) |  |  |
| from_at | [int64](#int64) |  |  |
| to_at | [int64](#int64) |  |  |






<a name="core.alert.ListNotificationResponse"></a>

### ListNotificationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| notification | [Notification](#core.alert.Notification) | repeated |  |






<a name="core.alert.ListRelAlertFindingRequest"></a>

### ListRelAlertFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_id | [uint32](#uint32) |  |  |
| finding_id | [uint32](#uint32) |  |  |
| from_at | [int64](#int64) |  |  |
| to_at | [int64](#int64) |  |  |






<a name="core.alert.ListRelAlertFindingResponse"></a>

### ListRelAlertFindingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rel_alert_finding | [RelAlertFinding](#core.alert.RelAlertFinding) | repeated |  |






<a name="core.alert.PutAlertCondNotificationRequest"></a>

### PutAlertCondNotificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_cond_notification | [AlertCondNotificationForUpsert](#core.alert.AlertCondNotificationForUpsert) |  |  |






<a name="core.alert.PutAlertCondNotificationResponse"></a>

### PutAlertCondNotificationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_cond_notification | [AlertCondNotification](#core.alert.AlertCondNotification) |  |  |






<a name="core.alert.PutAlertCondRuleRequest"></a>

### PutAlertCondRuleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_cond_rule | [AlertCondRuleForUpsert](#core.alert.AlertCondRuleForUpsert) |  |  |






<a name="core.alert.PutAlertCondRuleResponse"></a>

### PutAlertCondRuleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_cond_rule | [AlertCondRule](#core.alert.AlertCondRule) |  |  |






<a name="core.alert.PutAlertConditionRequest"></a>

### PutAlertConditionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_condition | [AlertConditionForUpsert](#core.alert.AlertConditionForUpsert) |  |  |






<a name="core.alert.PutAlertConditionResponse"></a>

### PutAlertConditionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_condition | [AlertCondition](#core.alert.AlertCondition) |  |  |






<a name="core.alert.PutAlertHistoryRequest"></a>

### PutAlertHistoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_history | [AlertHistoryForUpsert](#core.alert.AlertHistoryForUpsert) |  |  |






<a name="core.alert.PutAlertHistoryResponse"></a>

### PutAlertHistoryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_history | [AlertHistory](#core.alert.AlertHistory) |  |  |






<a name="core.alert.PutAlertRequest"></a>

### PutAlertRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert | [AlertForUpsert](#core.alert.AlertForUpsert) |  |  |






<a name="core.alert.PutAlertResponse"></a>

### PutAlertResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert | [Alert](#core.alert.Alert) |  |  |






<a name="core.alert.PutAlertRuleRequest"></a>

### PutAlertRuleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| alert_rule | [AlertRuleForUpsert](#core.alert.AlertRuleForUpsert) |  |  |






<a name="core.alert.PutAlertRuleResponse"></a>

### PutAlertRuleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| alert_rule | [AlertRule](#core.alert.AlertRule) |  |  |






<a name="core.alert.PutNotificationRequest"></a>

### PutNotificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| notification | [NotificationForUpsert](#core.alert.NotificationForUpsert) |  |  |






<a name="core.alert.PutNotificationResponse"></a>

### PutNotificationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| notification | [Notification](#core.alert.Notification) |  |  |






<a name="core.alert.PutRelAlertFindingRequest"></a>

### PutRelAlertFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| rel_alert_finding | [RelAlertFindingForUpsert](#core.alert.RelAlertFindingForUpsert) |  |  |






<a name="core.alert.PutRelAlertFindingResponse"></a>

### PutRelAlertFindingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rel_alert_finding | [RelAlertFinding](#core.alert.RelAlertFinding) |  |  |





 

 

 


<a name="core.alert.AlertService"></a>

### AlertService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListAlert | [ListAlertRequest](#core.alert.ListAlertRequest) | [ListAlertResponse](#core.alert.ListAlertResponse) | alert |
| GetAlert | [GetAlertRequest](#core.alert.GetAlertRequest) | [GetAlertResponse](#core.alert.GetAlertResponse) |  |
| PutAlert | [PutAlertRequest](#core.alert.PutAlertRequest) | [PutAlertResponse](#core.alert.PutAlertResponse) |  |
| DeleteAlert | [DeleteAlertRequest](#core.alert.DeleteAlertRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListAlertHistory | [ListAlertHistoryRequest](#core.alert.ListAlertHistoryRequest) | [ListAlertHistoryResponse](#core.alert.ListAlertHistoryResponse) | alert_history |
| GetAlertHistory | [GetAlertHistoryRequest](#core.alert.GetAlertHistoryRequest) | [GetAlertHistoryResponse](#core.alert.GetAlertHistoryResponse) |  |
| PutAlertHistory | [PutAlertHistoryRequest](#core.alert.PutAlertHistoryRequest) | [PutAlertHistoryResponse](#core.alert.PutAlertHistoryResponse) |  |
| DeleteAlertHistory | [DeleteAlertHistoryRequest](#core.alert.DeleteAlertHistoryRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListRelAlertFinding | [ListRelAlertFindingRequest](#core.alert.ListRelAlertFindingRequest) | [ListRelAlertFindingResponse](#core.alert.ListRelAlertFindingResponse) | rel_alert_finding |
| GetRelAlertFinding | [GetRelAlertFindingRequest](#core.alert.GetRelAlertFindingRequest) | [GetRelAlertFindingResponse](#core.alert.GetRelAlertFindingResponse) |  |
| PutRelAlertFinding | [PutRelAlertFindingRequest](#core.alert.PutRelAlertFindingRequest) | [PutRelAlertFindingResponse](#core.alert.PutRelAlertFindingResponse) |  |
| DeleteRelAlertFinding | [DeleteRelAlertFindingRequest](#core.alert.DeleteRelAlertFindingRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListAlertCondition | [ListAlertConditionRequest](#core.alert.ListAlertConditionRequest) | [ListAlertConditionResponse](#core.alert.ListAlertConditionResponse) | alert_condition |
| GetAlertCondition | [GetAlertConditionRequest](#core.alert.GetAlertConditionRequest) | [GetAlertConditionResponse](#core.alert.GetAlertConditionResponse) |  |
| PutAlertCondition | [PutAlertConditionRequest](#core.alert.PutAlertConditionRequest) | [PutAlertConditionResponse](#core.alert.PutAlertConditionResponse) |  |
| DeleteAlertCondition | [DeleteAlertConditionRequest](#core.alert.DeleteAlertConditionRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListAlertRule | [ListAlertRuleRequest](#core.alert.ListAlertRuleRequest) | [ListAlertRuleResponse](#core.alert.ListAlertRuleResponse) | alert_rule |
| GetAlertRule | [GetAlertRuleRequest](#core.alert.GetAlertRuleRequest) | [GetAlertRuleResponse](#core.alert.GetAlertRuleResponse) |  |
| PutAlertRule | [PutAlertRuleRequest](#core.alert.PutAlertRuleRequest) | [PutAlertRuleResponse](#core.alert.PutAlertRuleResponse) |  |
| DeleteAlertRule | [DeleteAlertRuleRequest](#core.alert.DeleteAlertRuleRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListAlertCondRule | [ListAlertCondRuleRequest](#core.alert.ListAlertCondRuleRequest) | [ListAlertCondRuleResponse](#core.alert.ListAlertCondRuleResponse) | alert_cond_rule |
| GetAlertCondRule | [GetAlertCondRuleRequest](#core.alert.GetAlertCondRuleRequest) | [GetAlertCondRuleResponse](#core.alert.GetAlertCondRuleResponse) |  |
| PutAlertCondRule | [PutAlertCondRuleRequest](#core.alert.PutAlertCondRuleRequest) | [PutAlertCondRuleResponse](#core.alert.PutAlertCondRuleResponse) |  |
| DeleteAlertCondRule | [DeleteAlertCondRuleRequest](#core.alert.DeleteAlertCondRuleRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListNotification | [ListNotificationRequest](#core.alert.ListNotificationRequest) | [ListNotificationResponse](#core.alert.ListNotificationResponse) | notification |
| GetNotification | [GetNotificationRequest](#core.alert.GetNotificationRequest) | [GetNotificationResponse](#core.alert.GetNotificationResponse) |  |
| PutNotification | [PutNotificationRequest](#core.alert.PutNotificationRequest) | [PutNotificationResponse](#core.alert.PutNotificationResponse) |  |
| DeleteNotification | [DeleteNotificationRequest](#core.alert.DeleteNotificationRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListAlertCondNotification | [ListAlertCondNotificationRequest](#core.alert.ListAlertCondNotificationRequest) | [ListAlertCondNotificationResponse](#core.alert.ListAlertCondNotificationResponse) | alert_cond_notification |
| GetAlertCondNotification | [GetAlertCondNotificationRequest](#core.alert.GetAlertCondNotificationRequest) | [GetAlertCondNotificationResponse](#core.alert.GetAlertCondNotificationResponse) |  |
| PutAlertCondNotification | [PutAlertCondNotificationRequest](#core.alert.PutAlertCondNotificationRequest) | [PutAlertCondNotificationResponse](#core.alert.PutAlertCondNotificationResponse) |  |
| DeleteAlertCondNotification | [DeleteAlertCondNotificationRequest](#core.alert.DeleteAlertCondNotificationRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| AnalyzeAlert | [AnalyzeAlertRequest](#core.alert.AnalyzeAlertRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) | AnalyzeAlert |

 



<a name="finding/entity.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## finding/entity.proto



<a name="core.finding.Finding"></a>

### Finding
Finding


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding_id | [uint64](#uint64) |  |  |
| description | [string](#string) |  |  |
| data_source | [string](#string) |  |  |
| data_source_id | [string](#string) |  |  |
| resource_name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| original_score | [float](#float) |  |  |
| original_max_score | [float](#float) |  |  |
| score | [float](#float) |  |  |
| data | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.finding.FindingForUpsert"></a>

### FindingForUpsert
Finding For Upsert
(Unique keys: project_id, data_source, data_source_id)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| description | [string](#string) |  |  |
| data_source | [string](#string) |  |  |
| data_source_id | [string](#string) |  |  |
| resource_name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| original_score | [float](#float) |  |  |
| original_max_score | [float](#float) |  |  |
| data | [string](#string) |  |  |






<a name="core.finding.FindingTag"></a>

### FindingTag
FindingTag


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding_tag_id | [uint64](#uint64) |  |  |
| finding_id | [uint64](#uint64) |  |  |
| project_id | [uint32](#uint32) |  |  |
| tag | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.finding.FindingTagForUpsert"></a>

### FindingTagForUpsert
FindingTag For Upsert
(Unique keys: finding_id, tag_key)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding_id | [uint64](#uint64) |  |  |
| project_id | [uint32](#uint32) |  |  |
| tag | [string](#string) |  |  |






<a name="core.finding.PendFinding"></a>

### PendFinding
PendFinding


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding_id | [uint64](#uint64) |  |  |
| project_id | [uint32](#uint32) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.finding.PendFindingForUpsert"></a>

### PendFindingForUpsert
PendFinding For upsert


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding_id | [uint64](#uint64) |  |  |
| project_id | [uint32](#uint32) |  |  |






<a name="core.finding.Resource"></a>

### Resource
Resource


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_id | [uint64](#uint64) |  |  |
| resource_name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.finding.ResourceForUpsert"></a>

### ResourceForUpsert
Resource For upsert
(Unique keys: project_id, resource_name)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |






<a name="core.finding.ResourceTag"></a>

### ResourceTag
ResourceTag


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_tag_id | [uint64](#uint64) |  |  |
| resource_id | [uint64](#uint64) |  |  |
| project_id | [uint32](#uint32) |  |  |
| tag | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.finding.ResourceTagForUpsert"></a>

### ResourceTagForUpsert
ResourceTag For upsert
(Unique keys: resource_id, tag_key)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_id | [uint64](#uint64) |  |  |
| project_id | [uint32](#uint32) |  |  |
| tag | [string](#string) |  |  |





 

 

 

 



<a name="finding/service.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## finding/service.proto



<a name="core.finding.DeleteFindingRequest"></a>

### DeleteFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| finding_id | [uint64](#uint64) |  |  |






<a name="core.finding.DeletePendFindingRequest"></a>

### DeletePendFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| finding_id | [uint64](#uint64) |  |  |






<a name="core.finding.DeleteResourceRequest"></a>

### DeleteResourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| resource_id | [uint64](#uint64) |  |  |






<a name="core.finding.GetFindingRequest"></a>

### GetFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| finding_id | [uint64](#uint64) |  |  |






<a name="core.finding.GetFindingResponse"></a>

### GetFindingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding | [Finding](#core.finding.Finding) |  |  |






<a name="core.finding.GetPendFindingRequest"></a>

### GetPendFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| finding_id | [uint64](#uint64) |  |  |






<a name="core.finding.GetPendFindingResponse"></a>

### GetPendFindingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pend_finding | [PendFinding](#core.finding.PendFinding) |  |  |






<a name="core.finding.GetResourceRequest"></a>

### GetResourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| resource_id | [uint64](#uint64) |  |  |






<a name="core.finding.GetResourceResponse"></a>

### GetResourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource | [Resource](#core.finding.Resource) |  |  |






<a name="core.finding.ListFindingRequest"></a>

### ListFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| data_source | [string](#string) | repeated |  |
| resource_name | [string](#string) | repeated |  |
| from_score | [float](#float) |  |  |
| to_score | [float](#float) |  |  |
| from_at | [int64](#int64) |  |  |
| to_at | [int64](#int64) |  |  |
| tag | [string](#string) | repeated |  |
| sort | [string](#string) |  |  |
| direction | [string](#string) |  |  |
| offset | [int32](#int32) |  |  |
| limit | [int32](#int32) |  |  |






<a name="core.finding.ListFindingResponse"></a>

### ListFindingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding_id | [uint64](#uint64) | repeated |  |
| count | [uint32](#uint32) |  |  |
| total | [uint32](#uint32) |  |  |






<a name="core.finding.ListFindingTagNameRequest"></a>

### ListFindingTagNameRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| from_at | [int64](#int64) |  |  |
| to_at | [int64](#int64) |  |  |
| sort | [string](#string) |  |  |
| direction | [string](#string) |  |  |
| offset | [int32](#int32) |  |  |
| limit | [int32](#int32) |  |  |






<a name="core.finding.ListFindingTagNameResponse"></a>

### ListFindingTagNameResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tag | [string](#string) | repeated |  |






<a name="core.finding.ListFindingTagRequest"></a>

### ListFindingTagRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| finding_id | [uint64](#uint64) |  |  |
| sort | [string](#string) |  |  |
| direction | [string](#string) |  |  |
| offset | [int32](#int32) |  |  |
| limit | [int32](#int32) |  |  |






<a name="core.finding.ListFindingTagResponse"></a>

### ListFindingTagResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tag | [FindingTag](#core.finding.FindingTag) | repeated |  |
| count | [uint32](#uint32) |  |  |
| total | [uint32](#uint32) |  |  |






<a name="core.finding.ListResourceRequest"></a>

### ListResourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| resource_name | [string](#string) | repeated |  |
| from_sum_score | [float](#float) |  |  |
| to_sum_score | [float](#float) |  |  |
| from_at | [int64](#int64) |  |  |
| to_at | [int64](#int64) |  |  |
| tag | [string](#string) | repeated |  |
| sort | [string](#string) |  |  |
| direction | [string](#string) |  |  |
| offset | [int32](#int32) |  |  |
| limit | [int32](#int32) |  |  |






<a name="core.finding.ListResourceResponse"></a>

### ListResourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_id | [uint64](#uint64) | repeated |  |
| count | [uint32](#uint32) |  |  |
| total | [uint32](#uint32) |  |  |






<a name="core.finding.ListResourceTagNameRequest"></a>

### ListResourceTagNameRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| from_at | [int64](#int64) |  |  |
| to_at | [int64](#int64) |  |  |
| sort | [string](#string) |  |  |
| direction | [string](#string) |  |  |
| offset | [int32](#int32) |  |  |
| limit | [int32](#int32) |  |  |






<a name="core.finding.ListResourceTagNameResponse"></a>

### ListResourceTagNameResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tag | [string](#string) | repeated |  |
| count | [uint32](#uint32) |  |  |
| total | [uint32](#uint32) |  |  |






<a name="core.finding.ListResourceTagRequest"></a>

### ListResourceTagRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| resource_id | [uint64](#uint64) |  |  |
| sort | [string](#string) |  |  |
| direction | [string](#string) |  |  |
| offset | [int32](#int32) |  |  |
| limit | [int32](#int32) |  |  |






<a name="core.finding.ListResourceTagResponse"></a>

### ListResourceTagResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tag | [ResourceTag](#core.finding.ResourceTag) | repeated |  |
| count | [uint32](#uint32) |  |  |
| total | [uint32](#uint32) |  |  |






<a name="core.finding.PutFindingRequest"></a>

### PutFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| finding | [FindingForUpsert](#core.finding.FindingForUpsert) |  |  |






<a name="core.finding.PutFindingResponse"></a>

### PutFindingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| finding | [Finding](#core.finding.Finding) |  |  |






<a name="core.finding.PutPendFindingRequest"></a>

### PutPendFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| pend_finding | [PendFindingForUpsert](#core.finding.PendFindingForUpsert) |  |  |






<a name="core.finding.PutPendFindingResponse"></a>

### PutPendFindingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pend_finding | [PendFinding](#core.finding.PendFinding) |  |  |






<a name="core.finding.PutResourceRequest"></a>

### PutResourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| resource | [ResourceForUpsert](#core.finding.ResourceForUpsert) |  |  |






<a name="core.finding.PutResourceResponse"></a>

### PutResourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource | [Resource](#core.finding.Resource) |  |  |






<a name="core.finding.TagFindingRequest"></a>

### TagFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| tag | [FindingTagForUpsert](#core.finding.FindingTagForUpsert) |  |  |






<a name="core.finding.TagFindingResponse"></a>

### TagFindingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tag | [FindingTag](#core.finding.FindingTag) |  |  |
| count | [uint32](#uint32) |  |  |
| total | [uint32](#uint32) |  |  |






<a name="core.finding.TagResourceRequest"></a>

### TagResourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| tag | [ResourceTagForUpsert](#core.finding.ResourceTagForUpsert) |  |  |






<a name="core.finding.TagResourceResponse"></a>

### TagResourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tag | [ResourceTag](#core.finding.ResourceTag) |  |  |






<a name="core.finding.UntagFindingRequest"></a>

### UntagFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| finding_tag_id | [uint64](#uint64) |  |  |






<a name="core.finding.UntagResourceRequest"></a>

### UntagResourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| resource_tag_id | [uint64](#uint64) |  |  |





 

 

 


<a name="core.finding.FindingService"></a>

### FindingService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListFinding | [ListFindingRequest](#core.finding.ListFindingRequest) | [ListFindingResponse](#core.finding.ListFindingResponse) | fiding |
| GetFinding | [GetFindingRequest](#core.finding.GetFindingRequest) | [GetFindingResponse](#core.finding.GetFindingResponse) |  |
| PutFinding | [PutFindingRequest](#core.finding.PutFindingRequest) | [PutFindingResponse](#core.finding.PutFindingResponse) |  |
| DeleteFinding | [DeleteFindingRequest](#core.finding.DeleteFindingRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListFindingTag | [ListFindingTagRequest](#core.finding.ListFindingTagRequest) | [ListFindingTagResponse](#core.finding.ListFindingTagResponse) |  |
| ListFindingTagName | [ListFindingTagNameRequest](#core.finding.ListFindingTagNameRequest) | [ListFindingTagNameResponse](#core.finding.ListFindingTagNameResponse) |  |
| TagFinding | [TagFindingRequest](#core.finding.TagFindingRequest) | [TagFindingResponse](#core.finding.TagFindingResponse) |  |
| UntagFinding | [UntagFindingRequest](#core.finding.UntagFindingRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListResource | [ListResourceRequest](#core.finding.ListResourceRequest) | [ListResourceResponse](#core.finding.ListResourceResponse) | resource |
| GetResource | [GetResourceRequest](#core.finding.GetResourceRequest) | [GetResourceResponse](#core.finding.GetResourceResponse) |  |
| PutResource | [PutResourceRequest](#core.finding.PutResourceRequest) | [PutResourceResponse](#core.finding.PutResourceResponse) |  |
| DeleteResource | [DeleteResourceRequest](#core.finding.DeleteResourceRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListResourceTag | [ListResourceTagRequest](#core.finding.ListResourceTagRequest) | [ListResourceTagResponse](#core.finding.ListResourceTagResponse) |  |
| ListResourceTagName | [ListResourceTagNameRequest](#core.finding.ListResourceTagNameRequest) | [ListResourceTagNameResponse](#core.finding.ListResourceTagNameResponse) |  |
| TagResource | [TagResourceRequest](#core.finding.TagResourceRequest) | [TagResourceResponse](#core.finding.TagResourceResponse) |  |
| UntagResource | [UntagResourceRequest](#core.finding.UntagResourceRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| GetPendFinding | [GetPendFindingRequest](#core.finding.GetPendFindingRequest) | [GetPendFindingResponse](#core.finding.GetPendFindingResponse) | pend_finding |
| PutPendFinding | [PutPendFindingRequest](#core.finding.PutPendFindingRequest) | [PutPendFindingResponse](#core.finding.PutPendFindingResponse) |  |
| DeletePendFinding | [DeletePendFindingRequest](#core.finding.DeletePendFindingRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |

 



<a name="iam/entity.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iam/entity.proto



<a name="core.iam.Policy"></a>

### Policy
Policy


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| policy_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| action_ptn | [string](#string) |  |  |
| resource_ptn | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.iam.PolicyForUpsert"></a>

### PolicyForUpsert
PolicyForUpsert


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| action_ptn | [string](#string) |  |  |
| resource_ptn | [string](#string) |  |  |






<a name="core.iam.Role"></a>

### Role
Role


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.iam.RoleForUpsert"></a>

### RoleForUpsert
RoleForUpsert


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |






<a name="core.iam.RolePolicy"></a>

### RolePolicy
RolePolicy


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role_id | [uint32](#uint32) |  |  |
| policy_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.iam.User"></a>

### User
User


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| sub | [string](#string) |  |  |
| name | [string](#string) |  |  |
| activated | [bool](#bool) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="core.iam.UserForUpsert"></a>

### UserForUpsert
UserForUpsert
(Unique keys: sub)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sub | [string](#string) |  |  |
| name | [string](#string) |  |  |
| activated | [bool](#bool) |  |  |






<a name="core.iam.UserRole"></a>

### UserRole
UserRole


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| role_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |





 

 

 

 



<a name="iam/policy.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iam/policy.proto



<a name="core.iam.AttachPolicyRequest"></a>

### AttachPolicyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| role_id | [uint32](#uint32) |  |  |
| policy_id | [uint32](#uint32) |  |  |






<a name="core.iam.AttachPolicyResponse"></a>

### AttachPolicyResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role_policy | [RolePolicy](#core.iam.RolePolicy) |  |  |






<a name="core.iam.DeletePolicyRequest"></a>

### DeletePolicyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| policy_id | [uint32](#uint32) |  |  |






<a name="core.iam.DetachPolicyRequest"></a>

### DetachPolicyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| role_id | [uint32](#uint32) |  |  |
| policy_id | [uint32](#uint32) |  |  |






<a name="core.iam.GetPolicyRequest"></a>

### GetPolicyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| policy_id | [uint32](#uint32) |  |  |






<a name="core.iam.GetPolicyResponse"></a>

### GetPolicyResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| policy | [Policy](#core.iam.Policy) |  |  |






<a name="core.iam.ListPolicyRequest"></a>

### ListPolicyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| role_id | [uint32](#uint32) |  |  |






<a name="core.iam.ListPolicyResponse"></a>

### ListPolicyResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| policy_id | [uint32](#uint32) | repeated |  |






<a name="core.iam.PutPolicyRequest"></a>

### PutPolicyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| policy | [PolicyForUpsert](#core.iam.PolicyForUpsert) |  |  |






<a name="core.iam.PutPolicyResponse"></a>

### PutPolicyResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| policy | [Policy](#core.iam.Policy) |  |  |





 

 

 

 



<a name="iam/role.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iam/role.proto



<a name="core.iam.AttachRoleRequest"></a>

### AttachRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| user_id | [uint32](#uint32) |  |  |
| role_id | [uint32](#uint32) |  |  |






<a name="core.iam.AttachRoleResponse"></a>

### AttachRoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_role | [UserRole](#core.iam.UserRole) |  |  |






<a name="core.iam.DeleteRoleRequest"></a>

### DeleteRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| role_id | [uint32](#uint32) |  |  |






<a name="core.iam.DetachRoleRequest"></a>

### DetachRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| user_id | [uint32](#uint32) |  |  |
| role_id | [uint32](#uint32) |  |  |






<a name="core.iam.GetRoleRequest"></a>

### GetRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| role_id | [uint32](#uint32) |  |  |






<a name="core.iam.GetRoleResponse"></a>

### GetRoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role | [Role](#core.iam.Role) |  |  |






<a name="core.iam.ListRoleRequest"></a>

### ListRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| user_id | [uint32](#uint32) |  |  |






<a name="core.iam.ListRoleResponse"></a>

### ListRoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role_id | [uint32](#uint32) | repeated |  |






<a name="core.iam.PutRoleRequest"></a>

### PutRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| role | [RoleForUpsert](#core.iam.RoleForUpsert) |  |  |






<a name="core.iam.PutRoleResponse"></a>

### PutRoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role | [Role](#core.iam.Role) |  |  |





 

 

 

 



<a name="iam/service.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iam/service.proto



<a name="core.iam.IsAdminRequest"></a>

### IsAdminRequest
IsAdminRequest
特定プロジェクトに依存しない管理者権限を持っているかどうかを返します


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |






<a name="core.iam.IsAdminResponse"></a>

### IsAdminResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  |  |






<a name="core.iam.IsAuthorizedRequest"></a>

### IsAuthorizedRequest
IsAuthorizedRequest
ユーザからのリクエストに対して、アクションやリソースへの認可を行います


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  | UserID,(e.g.)111 |
| project_id | [uint32](#uint32) |  | ProjectID,(e.g.)1001 |
| action_name | [string](#string) |  | Service&amp;API_name(&lt;service_name&gt;/&lt;API&gt;format),(e.g.)`finding/GetFinding` |
| resource_name | [string](#string) |  | System_resource_name(&lt;prefix&gt;/&lt;resouorce_name&gt;format),(e.g.)`aws:accessAnalyzer/samplebucket` |






<a name="core.iam.IsAuthorizedResponse"></a>

### IsAuthorizedResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  |  |





 

 

 


<a name="core.iam.IAMService"></a>

### IAMService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListUser | [ListUserRequest](#core.iam.ListUserRequest) | [ListUserResponse](#core.iam.ListUserResponse) | User |
| GetUser | [GetUserRequest](#core.iam.GetUserRequest) | [GetUserResponse](#core.iam.GetUserResponse) |  |
| PutUser | [PutUserRequest](#core.iam.PutUserRequest) | [PutUserResponse](#core.iam.PutUserResponse) |  |
| ListRole | [ListRoleRequest](#core.iam.ListRoleRequest) | [ListRoleResponse](#core.iam.ListRoleResponse) | Role |
| GetRole | [GetRoleRequest](#core.iam.GetRoleRequest) | [GetRoleResponse](#core.iam.GetRoleResponse) |  |
| PutRole | [PutRoleRequest](#core.iam.PutRoleRequest) | [PutRoleResponse](#core.iam.PutRoleResponse) |  |
| DeleteRole | [DeleteRoleRequest](#core.iam.DeleteRoleRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| AttachRole | [AttachRoleRequest](#core.iam.AttachRoleRequest) | [AttachRoleResponse](#core.iam.AttachRoleResponse) |  |
| DetachRole | [DetachRoleRequest](#core.iam.DetachRoleRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListPolicy | [ListPolicyRequest](#core.iam.ListPolicyRequest) | [ListPolicyResponse](#core.iam.ListPolicyResponse) | Policy |
| GetPolicy | [GetPolicyRequest](#core.iam.GetPolicyRequest) | [GetPolicyResponse](#core.iam.GetPolicyResponse) |  |
| PutPolicy | [PutPolicyRequest](#core.iam.PutPolicyRequest) | [PutPolicyResponse](#core.iam.PutPolicyResponse) |  |
| DeletePolicy | [DeletePolicyRequest](#core.iam.DeletePolicyRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| AttachPolicy | [AttachPolicyRequest](#core.iam.AttachPolicyRequest) | [AttachPolicyResponse](#core.iam.AttachPolicyResponse) |  |
| DetachPolicy | [DetachPolicyRequest](#core.iam.DetachPolicyRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| IsAuthorized | [IsAuthorizedRequest](#core.iam.IsAuthorizedRequest) | [IsAuthorizedResponse](#core.iam.IsAuthorizedResponse) | 認可（ユーザがリクエストしたアクションや、リソースに対しての認可を行います） |
| IsAdmin | [IsAdminRequest](#core.iam.IsAdminRequest) | [IsAdminResponse](#core.iam.IsAdminResponse) | 特定プロジェクトに依存しない管理者権限を持っているかどうかを返します |

 



<a name="iam/user.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iam/user.proto



<a name="core.iam.GetUserRequest"></a>

### GetUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| sub | [string](#string) |  |  |






<a name="core.iam.GetUserResponse"></a>

### GetUserResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#core.iam.User) |  |  |






<a name="core.iam.ListUserRequest"></a>

### ListUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| activated | [bool](#bool) |  |  |






<a name="core.iam.ListUserResponse"></a>

### ListUserResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) | repeated |  |






<a name="core.iam.PutUserRequest"></a>

### PutUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [UserForUpsert](#core.iam.UserForUpsert) |  |  |






<a name="core.iam.PutUserResponse"></a>

### PutUserResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#core.iam.User) |  |  |





 

 

 

 



<a name="project/entity.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## project/entity.proto



<a name="core.project.Project"></a>

### Project
Project


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |





 

 

 

 



<a name="project/service.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## project/service.proto



<a name="core.project.CreateProjectRequest"></a>

### CreateProjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  | project owner |
| name | [string](#string) |  |  |






<a name="core.project.CreateProjectResponse"></a>

### CreateProjectResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [Project](#core.project.Project) |  |  |






<a name="core.project.DeleteProjectRequest"></a>

### DeleteProjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |






<a name="core.project.ListProjectRequest"></a>

### ListProjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |






<a name="core.project.ListProjectResponse"></a>

### ListProjectResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [Project](#core.project.Project) | repeated |  |






<a name="core.project.UpdateProjectRequest"></a>

### UpdateProjectRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |






<a name="core.project.UpdateProjectResponse"></a>

### UpdateProjectResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project | [Project](#core.project.Project) |  |  |





 

 

 


<a name="core.project.ProjectService"></a>

### ProjectService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListProject | [ListProjectRequest](#core.project.ListProjectRequest) | [ListProjectResponse](#core.project.ListProjectResponse) | project |
| CreateProject | [CreateProjectRequest](#core.project.CreateProjectRequest) | [CreateProjectResponse](#core.project.CreateProjectResponse) |  |
| UpdateProject | [UpdateProjectRequest](#core.project.UpdateProjectRequest) | [UpdateProjectResponse](#core.project.UpdateProjectResponse) |  |
| DeleteProject | [DeleteProjectRequest](#core.project.DeleteProjectRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |

 



<a name="report/entity.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## report/entity.proto



<a name="core.report.ReportFinding"></a>

### ReportFinding
Report


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| report_finding_id | [uint32](#uint32) |  |  |
| report_date | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| project_name | [string](#string) |  |  |
| category | [string](#string) |  |  |
| data_source | [string](#string) |  |  |
| score | [float](#float) |  |  |
| count | [uint32](#uint32) |  |  |





 

 

 

 



<a name="report/service.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## report/service.proto



<a name="core.report.GetReportFindingAllRequest"></a>

### GetReportFindingAllRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| from_date | [string](#string) |  |  |
| to_date | [string](#string) |  |  |
| score | [float](#float) |  |  |
| data_source | [string](#string) | repeated |  |






<a name="core.report.GetReportFindingAllResponse"></a>

### GetReportFindingAllResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| report_finding | [ReportFinding](#core.report.ReportFinding) | repeated |  |






<a name="core.report.GetReportFindingRequest"></a>

### GetReportFindingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| from_date | [string](#string) |  |  |
| to_date | [string](#string) |  |  |
| score | [float](#float) |  |  |
| data_source | [string](#string) | repeated |  |






<a name="core.report.GetReportFindingResponse"></a>

### GetReportFindingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| report_finding | [ReportFinding](#core.report.ReportFinding) | repeated |  |





 

 

 


<a name="core.report.ReportService"></a>

### ReportService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetReportFinding | [GetReportFindingRequest](#core.report.GetReportFindingRequest) | [GetReportFindingResponse](#core.report.GetReportFindingResponse) | report |
| GetReportFindingAll | [GetReportFindingAllRequest](#core.report.GetReportFindingAllRequest) | [GetReportFindingAllResponse](#core.report.GetReportFindingAllResponse) |  |
| CollectReportFinding | [.google.protobuf.Empty](#google.protobuf.Empty) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

