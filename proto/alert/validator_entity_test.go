package alert

import (
	"fmt"
	"testing"
)

func TestValidateAlertForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *AlertForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &AlertForUpsert{AlertConditionId: 1001, Description: "test_alert", Severity: "high", Status: Status_ACTIVE, ProjectId: 1001},
			wantErr: false,
		},
		{
			name:  "NG too long Description",
			input: &AlertForUpsert{AlertConditionId: 1001, Description: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901", Severity: "high", Status: Status_ACTIVE, ProjectId: 1001},

			wantErr: true,
		},
		{
			name:    "NG required Description",
			input:   &AlertForUpsert{AlertConditionId: 1001, Description: "", Severity: "high", Status: Status_ACTIVE, ProjectId: 1001},
			wantErr: true,
		},
		{
			name:  "NG wrong severity",
			input: &AlertForUpsert{AlertConditionId: 1001, Description: "test_alert", Severity: "error", Status: Status_ACTIVE, ProjectId: 1001},

			wantErr: true,
		},
		{
			name:    "NG required severity",
			input:   &AlertForUpsert{AlertConditionId: 1001, Description: "test_alert", Status: Status_ACTIVE, ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG invalid status",
			input:   &AlertForUpsert{AlertConditionId: 1001, Description: "test_alert", Severity: "high", Status: Status_DEACTIVE, ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required status",
			input:   &AlertForUpsert{AlertConditionId: 1001, Description: "test_alert", Severity: "high", ProjectId: 1001},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fmt.Printf("Status: %v\n", c.input.Status)
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidateAlertHistoryForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *AlertHistoryForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &AlertHistoryForUpsert{AlertId: 1001, HistoryType: "created", Description: "test_alert", Severity: "high", ProjectId: 1001},
			wantErr: false,
		},
		{
			name:    "NG wrong history_type",
			input:   &AlertHistoryForUpsert{AlertId: 1001, HistoryType: "error", Description: "test_alert", Severity: "high", ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required history_type",
			input:   &AlertHistoryForUpsert{AlertId: 1001, Description: "test_alert", Severity: "high", ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG too long Description",
			input:   &AlertHistoryForUpsert{AlertId: 1001, HistoryType: "created", Description: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901", Severity: "high", ProjectId: 1001},
			wantErr: true,
		},
		{
			name:  "NG required Description",
			input: &AlertHistoryForUpsert{AlertId: 1001, HistoryType: "created", Severity: "high", ProjectId: 1001},

			wantErr: true,
		},
		{
			name:    "NG wrong severity",
			input:   &AlertHistoryForUpsert{AlertId: 1001, HistoryType: "created", Description: "test_alert", Severity: "error", ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required severity",
			input:   &AlertHistoryForUpsert{AlertId: 1001, HistoryType: "created", Description: "test_alert", ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG invalid json finding_history",
			input:   &AlertHistoryForUpsert{AlertId: 1001, HistoryType: "created", Description: "test_alert", Severity: "high", FindingHistory: "hogehoge", ProjectId: 1001},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidateRelAlertFindingForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *RelAlertFindingForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &RelAlertFindingForUpsert{AlertId: 1001, FindingId: 1001, ProjectId: 1001},
			wantErr: false,
		},
		{
			name:    "NG required finding_id",
			input:   &RelAlertFindingForUpsert{AlertId: 1001, ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required alert_id",
			input:   &RelAlertFindingForUpsert{FindingId: 1001, ProjectId: 1001},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidateAlertConditionForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *AlertConditionForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &AlertConditionForUpsert{Description: "test_alert_condition", Severity: "high", ProjectId: 1001, AndOr: "and", Enabled: true},
			wantErr: false,
		},
		{
			name:    "NG too long description",
			input:   &AlertConditionForUpsert{Description: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901", Severity: "high", ProjectId: 1001, AndOr: "and", Enabled: true},
			wantErr: true,
		},
		{
			name:    "NG required description",
			input:   &AlertConditionForUpsert{Severity: "high", ProjectId: 1001, AndOr: "and", Enabled: true},
			wantErr: true,
		},
		{
			name:    "NG wrong severity",
			input:   &AlertConditionForUpsert{Description: "test_alert_condition", Severity: "error", ProjectId: 1001, AndOr: "and", Enabled: true},
			wantErr: true,
		},
		{
			name:    "NG required severity",
			input:   &AlertConditionForUpsert{Description: "test_alert_condition", ProjectId: 1001, AndOr: "and", Enabled: true},
			wantErr: true,
		},
		{
			name:    "NG required project_id",
			input:   &AlertConditionForUpsert{Description: "test_alert_condition", Severity: "high", AndOr: "and", Enabled: true},
			wantErr: true,
		},
		{
			name:    "NG wrong and_or",
			input:   &AlertConditionForUpsert{Description: "test_alert_condition", Severity: "high", ProjectId: 1001, AndOr: "not", Enabled: true},
			wantErr: true,
		},
		{
			name:    "NG required and_or",
			input:   &AlertConditionForUpsert{Description: "test_alert_condition", Severity: "high", ProjectId: 1001, Enabled: true},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidateAlertRuleForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *AlertRuleForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &AlertRuleForUpsert{Name: "test_alert_rule", Score: 1.0, ProjectId: 1001, ResourceName: "test_resource", FindingCnt: 1},
			wantErr: false,
		},
		{
			name:    "NG too long name",
			input:   &AlertRuleForUpsert{Name: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901", Score: 1.0, ProjectId: 1001, ResourceName: "test_resource", Tag: "test_tag", FindingCnt: 1},
			wantErr: true,
		},
		{
			name:    "NG required name",
			input:   &AlertRuleForUpsert{Score: 1.0, ProjectId: 1001, ResourceName: "test_resource", Tag: "test_tag", FindingCnt: 1},
			wantErr: true,
		},
		{
			name:    "NG too small value score",
			input:   &AlertRuleForUpsert{Name: "test_alert_rule", Score: -0.1, ProjectId: 1001, ResourceName: "test_resource", Tag: "test_tag", FindingCnt: 1},
			wantErr: true,
		},
		{
			name:    "NG too large value score",
			input:   &AlertRuleForUpsert{Name: "test_alert_rule", Score: 1.1, ProjectId: 1001, ResourceName: "test_resource", Tag: "test_tag", FindingCnt: 1},
			wantErr: true,
		},
		{
			name:    "NG required project_id",
			input:   &AlertRuleForUpsert{Name: "test_alert_rule", Score: 1.0, ResourceName: "test_resource", Tag: "test_tag", FindingCnt: 1},
			wantErr: true,
		},
		{
			name:    "NG too long resource_name",
			input:   &AlertRuleForUpsert{Name: "test_alert_rule", Score: 1.0, ProjectId: 1001, ResourceName: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456", Tag: "test_tag", FindingCnt: 1},
			wantErr: true,
		},
		{
			name:    "NG too long tag",
			input:   &AlertRuleForUpsert{Name: "test_alert_rule", Score: 1.0, ProjectId: 1001, ResourceName: "test_resource", Tag: "12345678901234567890123456789012345678901234567890123456789012345", FindingCnt: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidateAlertCondRuleForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *AlertCondRuleForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &AlertCondRuleForUpsert{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001},
			wantErr: false,
		},
		{
			name:    "NG required project_id",
			input:   &AlertCondRuleForUpsert{AlertConditionId: 1001, AlertRuleId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required alert_condition_id",
			input:   &AlertCondRuleForUpsert{ProjectId: 1001, AlertRuleId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required alert_rule_id",
			input:   &AlertCondRuleForUpsert{ProjectId: 1001, AlertConditionId: 1001},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidateNotificationForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *NotificationForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &NotificationForUpsert{ProjectId: 1001, Name: "test_name", Type: "test_type", NotifySetting: `{"test_setting": "test_val"}`},
			wantErr: false,
		},
		{
			name:    "NG required project_id",
			input:   &NotificationForUpsert{Name: "test_name", Type: "test_type", NotifySetting: `{"test_setting": "test_val"}`},
			wantErr: true,
		},
		{
			name:    "NG too long name",
			input:   &NotificationForUpsert{ProjectId: 1001, Name: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901", Type: "test_type", NotifySetting: `{"test_setting": "test_val"}`},
			wantErr: true,
		},
		{
			name:    "NG required name",
			input:   &NotificationForUpsert{ProjectId: 1001, Type: "test_type", NotifySetting: `{"test_setting": "test_val"}`},
			wantErr: true,
		},
		{
			name:    "NG too long type",
			input:   &NotificationForUpsert{ProjectId: 1001, Name: "test_name", Type: "12345678901234567890123456789012345678901234567890123456789012345", NotifySetting: `{"test_setting": "test_val"}`},
			wantErr: true,
		},
		{
			name:    "NG required type",
			input:   &NotificationForUpsert{ProjectId: 1001, Name: "test_name", NotifySetting: `{"test_setting": "test_val"}`},
			wantErr: true,
		},
		{
			name:    "NG invalid json notify_setting",
			input:   &NotificationForUpsert{ProjectId: 1001, Name: "test_name", Type: "test_type", NotifySetting: `{"test_setting": "test_val"`},
			wantErr: true,
		},
		{
			name:    "NG required notify_setting",
			input:   &NotificationForUpsert{ProjectId: 1001, Name: "test_name", Type: "test_type"},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidateAlertCondNotificationForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *AlertCondNotificationForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &AlertCondNotificationForUpsert{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001, NotifiedAt: 1602660000},
			wantErr: false,
		},
		{
			name:    "NG required project_id",
			input:   &AlertCondNotificationForUpsert{AlertConditionId: 1001, NotificationId: 1001, NotifiedAt: 1602660000},
			wantErr: true,
		},
		{
			name:    "NG required alert_condition_id",
			input:   &AlertCondNotificationForUpsert{ProjectId: 1001, NotificationId: 1001, NotifiedAt: 1602660000},
			wantErr: true,
		},
		{
			name:    "NG required notification_id",
			input:   &AlertCondNotificationForUpsert{ProjectId: 1001, AlertConditionId: 1001, NotifiedAt: 1602660000},
			wantErr: true,
		},
		{
			name:    "NG too small notified_at",
			input:   &AlertCondNotificationForUpsert{ProjectId: 1001, AlertConditionId: 1001, NotifiedAt: -1},
			wantErr: true,
		},
		{
			name:    "NG required large notified_at",
			input:   &AlertCondNotificationForUpsert{ProjectId: 1001, AlertConditionId: 1001, NotifiedAt: 253402268400},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}
