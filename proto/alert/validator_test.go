package alert

import (
	"testing"
	"time"
)

func TestValidate_ListAlertRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *ListAlertRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListAlertRequest{ProjectId: 111, Severity: []string{"high"}, Description: "test_list_alert", FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListAlertRequest{Severity: []string{"high"}, Description: "test_list_alert", FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:  "NG wrong severity",
			input: &ListAlertRequest{ProjectId: 111, Severity: []string{"error"}, Description: "test_list_alert", FromAt: now.Unix(), ToAt: now.Unix()},

			wantErr: true,
		},
		{
			name:  "NG too long description",
			input: &ListAlertRequest{ProjectId: 111, Severity: []string{"high"}, Description: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901", FromAt: now.Unix(), ToAt: now.Unix()},

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

func TestValidate_GetAlertRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetAlertRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetAlertRequest{ProjectId: 1, AlertId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetAlertRequest{AlertId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(alert_id)",
			input:   &GetAlertRequest{ProjectId: 1},
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

func TestValidate_PutAlertRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutAlertRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutAlertRequest{ProjectId: 1001, Alert: &AlertForUpsert{AlertConditionId: 1001, Description: "test_alert", Severity: "high", Activated: true, ProjectId: 1001}},
			wantErr: false,
		},
		{
			name:    "NG Required(alert)",
			input:   &PutAlertRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:  "NG Not Equal(project_id != tag.project_id)",
			input: &PutAlertRequest{ProjectId: 1000, Alert: &AlertForUpsert{AlertConditionId: 1001, Description: "test_alert", Severity: "high", Activated: true, ProjectId: 1001}},

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

func TestValidate_DeleteAlertRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteAlertRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteAlertRequest{ProjectId: 1, AlertId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteAlertRequest{AlertId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(alert_id)",
			input:   &DeleteAlertRequest{ProjectId: 1},
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

func TestValidate_ListAlertHistoryRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *ListAlertHistoryRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListAlertHistoryRequest{ProjectId: 1001, AlertId: 1001, HistoryType: []string{"created"}, Severity: []string{"high"}, FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListAlertHistoryRequest{AlertId: 1001, HistoryType: []string{"created"}, Severity: []string{"high"}, FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG wrong history_type",
			input:   &ListAlertHistoryRequest{ProjectId: 1001, AlertId: 1001, HistoryType: []string{"error"}, Severity: []string{"high"}, FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG wrong severity",
			input:   &ListAlertHistoryRequest{ProjectId: 1001, AlertId: 1001, HistoryType: []string{"updated"}, Severity: []string{"info"}, FromAt: now.Unix(), ToAt: now.Unix()},
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

func TestValidate_GetAlertHistoryRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetAlertHistoryRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetAlertHistoryRequest{ProjectId: 1, AlertHistoryId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetAlertHistoryRequest{AlertHistoryId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(alert_history_id)",
			input:   &GetAlertHistoryRequest{ProjectId: 1},
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

func TestValidate_PutAlertHistoryRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutAlertHistoryRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutAlertHistoryRequest{ProjectId: 1001, AlertHistory: &AlertHistoryForUpsert{AlertId: 1001, HistoryType: "created", Description: "test_alert", Severity: "high", ProjectId: 1001}},
			wantErr: false,
		},
		{
			name:    "NG Required(alert_history)",
			input:   &PutAlertHistoryRequest{ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != tag.project_id)",
			input:   &PutAlertHistoryRequest{ProjectId: 1000, AlertHistory: &AlertHistoryForUpsert{AlertId: 1001, HistoryType: "created", Description: "test_alert", Severity: "high", ProjectId: 1001}},
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

func TestValidate_DeleteAlertHistoryRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteAlertHistoryRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteAlertHistoryRequest{ProjectId: 1, AlertHistoryId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteAlertHistoryRequest{AlertHistoryId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(alert_history_id)",
			input:   &DeleteAlertHistoryRequest{ProjectId: 1},
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

func TestValidate_ListRelAlertFindingRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *ListRelAlertFindingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListRelAlertFindingRequest{ProjectId: 1001, AlertId: 1001, FindingId: 1001, FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListRelAlertFindingRequest{AlertId: 1001, FindingId: 1001, FromAt: now.Unix(), ToAt: now.Unix()},
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

func TestValidate_GetRelAlertFindingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetRelAlertFindingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetRelAlertFindingRequest{ProjectId: 1001, AlertId: 1001, FindingId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetRelAlertFindingRequest{AlertId: 1001, FindingId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(alert_id)",
			input:   &GetRelAlertFindingRequest{ProjectId: 1001, FindingId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(finding_id)",
			input:   &GetRelAlertFindingRequest{ProjectId: 1001, AlertId: 1001},
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

func TestValidate_PutRelAlertFindingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutRelAlertFindingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutRelAlertFindingRequest{ProjectId: 1001, RelAlertFinding: &RelAlertFindingForUpsert{AlertId: 1001, FindingId: 1001, ProjectId: 1001}},
			wantErr: false,
		},
		{
			name:    "NG Required(rel_alert_finding)",
			input:   &PutRelAlertFindingRequest{ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != tag.project_id)",
			input:   &PutRelAlertFindingRequest{ProjectId: 1000, RelAlertFinding: &RelAlertFindingForUpsert{AlertId: 1001, FindingId: 1001, ProjectId: 1001}},
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

func TestValidate_DeleteRelAlertFindingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteRelAlertFindingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteRelAlertFindingRequest{ProjectId: 1, AlertId: 1001, FindingId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteRelAlertFindingRequest{AlertId: 1001, FindingId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(alert_id)",
			input:   &DeleteRelAlertFindingRequest{ProjectId: 1, FindingId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(finding_id)",
			input:   &DeleteRelAlertFindingRequest{ProjectId: 1, AlertId: 1001},
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

func TestValidate_ListAlertConditionRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *ListAlertConditionRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListAlertConditionRequest{ProjectId: 1001, Severity: []string{"high", "medium"}, Enabled: true, FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListAlertConditionRequest{Severity: []string{"high"}, Enabled: true, FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG wrong severity",
			input:   &ListAlertConditionRequest{ProjectId: 1001, Severity: []string{"error"}, Enabled: true, FromAt: now.Unix(), ToAt: now.Unix()},
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

func TestValidate_GetAlertConditionRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetAlertConditionRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetAlertConditionRequest{ProjectId: 1001, AlertConditionId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetAlertConditionRequest{AlertConditionId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(alert_condition_id)",
			input:   &GetAlertConditionRequest{ProjectId: 1001},
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

func TestValidate_PutAlertConditionRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutAlertConditionRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutAlertConditionRequest{ProjectId: 1001, AlertCondition: &AlertConditionForUpsert{Description: "test_alert_condition", Severity: "high", ProjectId: 1001, AndOr: "and", Enabled: true}},
			wantErr: false,
		},
		{
			name:    "NG Required(alert_condition)",
			input:   &PutAlertConditionRequest{ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != tag.project_id)",
			input:   &PutAlertConditionRequest{ProjectId: 1000, AlertCondition: &AlertConditionForUpsert{Description: "test_alert_condition", Severity: "high", ProjectId: 1001, AndOr: "and", Enabled: true}},
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

func TestValidate_DeleteAlertConditionRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteAlertConditionRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteAlertConditionRequest{ProjectId: 1001, AlertConditionId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteAlertConditionRequest{AlertConditionId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(alert_condition_id)",
			input:   &DeleteAlertConditionRequest{ProjectId: 1001},
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

func TestValidate_ListAlertRuleRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *ListAlertRuleRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListAlertRuleRequest{ProjectId: 1001, FromScore: 0.1, ToScore: 1.0, FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListAlertRuleRequest{FromScore: 0.1, ToScore: 1.0, FromAt: now.Unix(), ToAt: now.Unix()},
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

func TestValidate_GetAlertRuleRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetAlertRuleRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetAlertRuleRequest{ProjectId: 1001, AlertRuleId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetAlertRuleRequest{AlertRuleId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(alert_condition_id)",
			input:   &GetAlertRuleRequest{ProjectId: 1001},
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

func TestValidate_PutAlertRuleRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutAlertRuleRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutAlertRuleRequest{ProjectId: 1001, AlertRule: &AlertRuleForUpsert{Name: "test_alert_rule", Score: 1.0, ProjectId: 1001, ResourceName: "test_resource", Tag: "test_tag", FindingCnt: 1}},
			wantErr: false,
		},
		{
			name:    "NG Required(alert_condition)",
			input:   &PutAlertRuleRequest{ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != tag.project_id)",
			input:   &PutAlertRuleRequest{ProjectId: 1000, AlertRule: &AlertRuleForUpsert{Name: "test_alert_rule", Score: 1.0, ProjectId: 1001, ResourceName: "test_resource", Tag: "test_tag", FindingCnt: 1}},
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

func TestValidate_DeleteAlertRuleRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteAlertRuleRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteAlertRuleRequest{ProjectId: 1001, AlertRuleId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteAlertRuleRequest{AlertRuleId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(alert_id)",
			input:   &DeleteAlertRuleRequest{ProjectId: 1001},
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

func TestValidate_ListAlertCondRuleRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *ListAlertCondRuleRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListAlertCondRuleRequest{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001, FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListAlertCondRuleRequest{AlertConditionId: 1001, AlertRuleId: 1001, FromAt: now.Unix(), ToAt: now.Unix()},
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

func TestValidate_GetAlertCondRuleRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetAlertCondRuleRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetAlertCondRuleRequest{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetAlertCondRuleRequest{AlertConditionId: 1001, AlertRuleId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(alert_condition_id)",
			input:   &GetAlertCondRuleRequest{ProjectId: 1001, AlertRuleId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(alert_rule_id)",
			input:   &GetAlertCondRuleRequest{ProjectId: 1001, AlertConditionId: 1001},
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

func TestValidate_PutAlertCondRuleRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutAlertCondRuleRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutAlertCondRuleRequest{ProjectId: 1001, AlertCondRule: &AlertCondRuleForUpsert{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001}},
			wantErr: false,
		},
		{
			name:    "NG Required(alert_condition)",
			input:   &PutAlertCondRuleRequest{ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != tag.project_id)",
			input:   &PutAlertCondRuleRequest{ProjectId: 1000, AlertCondRule: &AlertCondRuleForUpsert{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001}},
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

func TestValidate_DeleteAlertCondRuleRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteAlertCondRuleRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteAlertCondRuleRequest{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteAlertCondRuleRequest{AlertConditionId: 1001, AlertRuleId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(alert_condition_id)",
			input:   &DeleteAlertCondRuleRequest{ProjectId: 1001, AlertRuleId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(alert_rule_id)",
			input:   &DeleteAlertCondRuleRequest{ProjectId: 1001, AlertConditionId: 1001},
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

func TestValidate_ListNotificationRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *ListNotificationRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListNotificationRequest{ProjectId: 1001, Type: "test_notification", FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListNotificationRequest{Type: "test_notification", FromAt: now.Unix(), ToAt: now.Unix()},
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

func TestValidate_GetNotificationRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetNotificationRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetNotificationRequest{ProjectId: 1001, NotificationId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetNotificationRequest{NotificationId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(notification_id)",
			input:   &GetNotificationRequest{NotificationId: 1001},
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

func TestValidate_PutNotificationRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutNotificationRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutNotificationRequest{ProjectId: 1001, Notification: &NotificationForUpsert{ProjectId: 1001, Name: "test_name", Type: "test_type", NotifySetting: `{"test_setting": "test_val"}`}},
			wantErr: false,
		},
		{
			name:    "NG Required(alert_condition)",
			input:   &PutNotificationRequest{ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != tag.project_id)",
			input:   &PutNotificationRequest{ProjectId: 1000, Notification: &NotificationForUpsert{ProjectId: 1001, Name: "test_name", Type: "test_type", NotifySetting: `{"test_setting": "test_val"}`}},
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

func TestValidate_DeleteNotificationRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteNotificationRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteNotificationRequest{ProjectId: 1001, NotificationId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteNotificationRequest{NotificationId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(notification_id)",
			input:   &DeleteNotificationRequest{ProjectId: 1001},
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

func TestValidate_ListAlertCondNotificationRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListAlertCondNotificationRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListAlertCondNotificationRequest{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListAlertCondNotificationRequest{AlertConditionId: 1001, NotificationId: 1001},
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

func TestValidate_GetAlertCondNotificationRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetAlertCondNotificationRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetAlertCondNotificationRequest{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetAlertCondNotificationRequest{AlertConditionId: 1001, NotificationId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(alert_condition_id)",
			input:   &GetAlertCondNotificationRequest{ProjectId: 1001, AlertConditionId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(notification_id)",
			input:   &GetAlertCondNotificationRequest{ProjectId: 1001, NotificationId: 1001},
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

func TestValidate_PutAlertCondNotificationRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutAlertCondNotificationRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutAlertCondNotificationRequest{ProjectId: 1001, AlertCondNotification: &AlertCondNotificationForUpsert{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001}},
			wantErr: false,
		},
		{
			name:    "NG Required(alert_condition)",
			input:   &PutAlertCondNotificationRequest{ProjectId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != tag.project_id)",
			input:   &PutAlertCondNotificationRequest{ProjectId: 1000, AlertCondNotification: &AlertCondNotificationForUpsert{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001}},
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

func TestValidate_DeleteAlertCondNotificationRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteAlertCondNotificationRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteAlertCondNotificationRequest{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteAlertCondNotificationRequest{AlertConditionId: 1001, NotificationId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(alert_condition_id)",
			input:   &DeleteAlertCondNotificationRequest{ProjectId: 1001, NotificationId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(notification_id)",
			input:   &DeleteAlertCondNotificationRequest{ProjectId: 1001, AlertConditionId: 1001},
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

func TestValidate_AnalyzeAlertRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *AnalyzeAlertRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &AnalyzeAlertRequest{ProjectId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &AnalyzeAlertRequest{},
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
