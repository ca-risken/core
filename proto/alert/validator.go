package alert

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	is "github.com/go-ozzo/ozzo-validation/v4/is"
)

// Validate ListAlertRequest
func (r *ListAlertRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.Severity, validation.Each(validation.In("high", "medium", "low"))),
		validation.Field(&r.Description, validation.Length(0, 200)),
		validation.Field(&r.FromAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&r.ToAt, validation.Min(0), validation.Max(253402268399)),   //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}

// Validate GetAlertRequest
func (r *GetAlertRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.AlertId, validation.Required),
	)
}

// Validate PutAlertRequest
func (r *PutAlertRequest) Validate() error {
	if validation.IsEmpty(r.Alert) {
		return errors.New("Required alert parameter")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.In(r.Alert.ProjectId)),
	); err != nil {
		return err
	}
	return r.Alert.Validate()
}

// Validate DeleteAlertRequest
func (r *DeleteAlertRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.AlertId, validation.Required),
	)
}

// Validate ListAlertHistoryRequest
func (r *ListAlertHistoryRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

// Validate GetAlertHistoryRequest
func (r *GetAlertHistoryRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.AlertHistoryId, validation.Required),
	)
}

// Validate PutAlertHistoryRequest
func (r *PutAlertHistoryRequest) Validate() error {
	if validation.IsEmpty(r.AlertHistory) {
		return errors.New("Required alert_history parameter")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.In(r.AlertHistory.ProjectId)),
	); err != nil {
		return err
	}
	return r.AlertHistory.Validate()
}

// Validate DeleteAlertHistoryRequest
func (r *DeleteAlertHistoryRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.AlertHistoryId, validation.Required),
	)
}

// Validate ListRelAlertFindingRequest
func (r *ListRelAlertFindingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.FromAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&r.ToAt, validation.Min(0), validation.Max(253402268399)),   //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}

// Validate GetRelAlertFindingRequest
func (r *GetRelAlertFindingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.AlertId, validation.Required),
		validation.Field(&r.FindingId, validation.Required),
	)
}

// Validate PutRelAlertFindingRequest
func (r *PutRelAlertFindingRequest) Validate() error {
	if validation.IsEmpty(r.RelAlertFinding) {
		return errors.New("Required rel_alert_finding parameter")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.In(r.RelAlertFinding.ProjectId)),
	); err != nil {
		return err
	}
	return r.RelAlertFinding.Validate()
}

// Validate DeleteRelAlertFindingRequest
func (r *DeleteRelAlertFindingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.AlertId, validation.Required),
		validation.Field(&r.FindingId, validation.Required),
	)
}

// Validate ListAlertConditionRequest
func (r *ListAlertConditionRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.Severity, validation.Each(validation.In("high", "medium", "low"))),
		validation.Field(&r.FromAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&r.ToAt, validation.Min(0), validation.Max(253402268399)),   //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}

// Validate GetAlertConditionRequest
func (r *GetAlertConditionRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.AlertConditionId, validation.Required),
	)
}

// Validate PutAlertConditionRequest
func (r *PutAlertConditionRequest) Validate() error {
	if validation.IsEmpty(r.AlertCondition) {
		return errors.New("Required rel_alert_finding parameter")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.In(r.AlertCondition.ProjectId)),
	); err != nil {
		return err
	}
	return r.AlertCondition.Validate()
}

// Validate DeleteAlertConditionRequest
func (r *DeleteAlertConditionRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.AlertConditionId, validation.Required),
	)
}

// Validate ListAlertRuleRequest
func (r *ListAlertRuleRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.FromScore, validation.Min(0.0), validation.Max(1.0)),
		validation.Field(&r.ToScore, validation.Min(0.0), validation.Max(1.0)),
		validation.Field(&r.FromAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&r.ToAt, validation.Min(0), validation.Max(253402268399)),   //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}

// Validate GetAlertRuleRequest
func (r *GetAlertRuleRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.AlertRuleId, validation.Required),
	)
}

// Validate PutAlertRuleRequest
func (r *PutAlertRuleRequest) Validate() error {
	if validation.IsEmpty(r.AlertRule) {
		return errors.New("Required alert_rule parameter")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.In(r.AlertRule.ProjectId)),
	); err != nil {
		return err
	}
	return r.AlertRule.Validate()
}

// Validate DeleteAlertRuleRequest
func (r *DeleteAlertRuleRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.AlertRuleId, validation.Required),
	)
}

// Validate ListAlertCondRuleRequest
func (r *ListAlertCondRuleRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.FromAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&r.ToAt, validation.Min(0), validation.Max(253402268399)),   //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}

// Validate GetAlertCondRuleRequest
func (r *GetAlertCondRuleRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.AlertConditionId, validation.Required),
		validation.Field(&r.AlertRuleId, validation.Required),
	)
}

// Validate PutAlertCondRuleRequest
func (r *PutAlertCondRuleRequest) Validate() error {
	if validation.IsEmpty(r.AlertCondRule) {
		return errors.New("Required alert_condition_rule parameter")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.In(r.AlertCondRule.ProjectId)),
	); err != nil {
		return err
	}
	return r.AlertCondRule.Validate()
}

// Validate DeleteAlertCondRuleRequest
func (r *DeleteAlertCondRuleRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.AlertConditionId, validation.Required),
		validation.Field(&r.AlertRuleId, validation.Required),
	)
}

// Validate ListNotificationRequest
func (r *ListNotificationRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.FromAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&r.ToAt, validation.Min(0), validation.Max(253402268399)),   //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}

// Validate GetNotificationRequest
func (r *GetNotificationRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.NotificationId, validation.Required),
	)
}

// Validate PutNotificationRequest
func (r *PutNotificationRequest) Validate() error {
	if validation.IsEmpty(r.Notification) {
		return errors.New("Required notification parameter")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.In(r.Notification.ProjectId)),
	); err != nil {
		return err
	}
	return r.Notification.Validate()
}

// Validate DeleteNotificationRequest
func (r *DeleteNotificationRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.NotificationId, validation.Required),
	)
}

// Validate TestNotificationRequest
func (r *TestNotificationRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.NotificationId, validation.Required),
	)
}

// Validate ListAlertCondNotificationRequest
func (r *ListAlertCondNotificationRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.FromAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&r.ToAt, validation.Min(0), validation.Max(253402268399)),   //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}

// Validate GetAlertCondNotificationRequest
func (r *GetAlertCondNotificationRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.NotificationId, validation.Required),
		validation.Field(&r.AlertConditionId, validation.Required),
	)
}

// Validate PutAlertCondNotificationRequest
func (r *PutAlertCondNotificationRequest) Validate() error {
	if validation.IsEmpty(r.AlertCondNotification) {
		return errors.New("Required rel_put_alert_condition_notification parameter")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.In(r.AlertCondNotification.ProjectId)),
	); err != nil {
		return err
	}
	return r.AlertCondNotification.Validate()
}

// Validate DeleteAlertCondNotificationRequest
func (r *DeleteAlertCondNotificationRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.NotificationId, validation.Required),
		validation.Field(&r.AlertConditionId, validation.Required),
	)
}

// Validate AnalyzeAlertRequest
func (r *AnalyzeAlertRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

/*
 * entities
**/

// Validate AlertForUpsert
func (e *AlertForUpsert) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.AlertConditionId, validation.Required),
		validation.Field(&e.Description, validation.Required, validation.Length(0, 200)),
		validation.Field(&e.Severity, validation.Required, validation.In("high", "medium", "low")),
		validation.Field(&e.Status, validation.Required, validation.In(Status_ACTIVE, Status_PENDING)),
		validation.Field(&e.Status, validation.Required),
		validation.Field(&e.ProjectId, validation.Required),
	)
}

// Validate AlertHistoryForUpsert
func (e *AlertHistoryForUpsert) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.Description, validation.Required, validation.Length(0, 200)),
		validation.Field(&e.HistoryType, validation.Required, validation.In("created", "updated", "deleted")),
		validation.Field(&e.Severity, validation.Required, validation.In("high", "medium", "low")),
		validation.Field(&e.AlertId, validation.Required),
		validation.Field(&e.ProjectId, validation.Required),
		validation.Field(&e.FindingHistory, is.JSON),
	)
}

// Validate RelAlertFindingForUpsert
func (e *RelAlertFindingForUpsert) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.ProjectId, validation.Required),
		validation.Field(&e.AlertId, validation.Required),
		validation.Field(&e.FindingId, validation.Required),
	)
}

// Validate AlertConditionForUpsert
func (e *AlertConditionForUpsert) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.ProjectId, validation.Required),
		validation.Field(&e.Description, validation.Required, validation.Length(0, 200)),
		validation.Field(&e.AndOr, validation.Required, validation.In("and", "or")),
		validation.Field(&e.Severity, validation.Required, validation.In("high", "medium", "low")),
	)
}

// Validate AlertRuleForUpsert
func (e *AlertRuleForUpsert) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.ProjectId, validation.Required),
		validation.Field(&e.Name, validation.Required, validation.Length(0, 200)),
		validation.Field(&e.Score, validation.Min(0.0), validation.Max(1.0)),
		validation.Field(&e.ResourceName, validation.Length(0, 255)),
		validation.Field(&e.Tag, validation.Length(0, 64)),
		validation.Field(&e.FindingCnt, validation.Min(uint(1))),
	)
}

// Validate AlertCondRuleForUpsert
func (e *AlertCondRuleForUpsert) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.ProjectId, validation.Required),
		validation.Field(&e.AlertConditionId, validation.Required),
		validation.Field(&e.AlertRuleId, validation.Required),
	)
}

// Validate NotificationForUpsert
func (e *NotificationForUpsert) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.ProjectId, validation.Required),
		validation.Field(&e.Type, validation.Required, validation.Length(0, 64)),
		validation.Field(&e.Name, validation.Required, validation.Length(0, 200)),
		validation.Field(&e.NotifySetting, validation.Required, validation.By(validateNotifySetting)),
	)
}

// Validate AlertCondNotificationForUpsert
func (e *AlertCondNotificationForUpsert) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.ProjectId, validation.Required),
		validation.Field(&e.NotificationId, validation.Required),
		validation.Field(&e.AlertConditionId, validation.Required),
		validation.Field(&e.NotifiedAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}

type slackNotifySetting struct {
	WebhookURL string `json:"webhook_url"`
}

func validateNotifySetting(value interface{}) error {
	s, _ := value.(string)
	var setting slackNotifySetting
	if err := json.Unmarshal([]byte(s), &setting); err != nil {
		return fmt.Errorf("invalid json, %w", err)
	}
	if strings.TrimSpace(setting.WebhookURL) == "" {
		return errors.New("required webhook url in json")
	}
	err := validation.Validate(strings.TrimSpace(setting.WebhookURL), validation.Required, is.URL)
	if err != nil {
		return err
	}
	return nil
}
