package alert

import (
	"testing"
	"time"
)

func TestValidateListAlertRequest(t *testing.T) {
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
			name:    "NG wrong severity",
			input:   &ListAlertRequest{ProjectId: 111, Severity: []string{"error"}, Description: "test_list_alert", FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too long description",
			input:   &ListAlertRequest{ProjectId: 111, Severity: []string{"high"}, Description: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901", FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too small from_at",
			input:   &ListAlertRequest{ProjectId: 111, Severity: []string{"high"}, Description: "test_list_alert", FromAt: -1, ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too large from_at",
			input:   &ListAlertRequest{ProjectId: 111, Severity: []string{"high"}, Description: "test_list_alert", FromAt: 253402268400, ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too small to_at",
			input:   &ListAlertRequest{ProjectId: 111, Severity: []string{"high"}, Description: "test_list_alert", FromAt: now.Unix(), ToAt: -1},
			wantErr: true,
		},
		{
			name:    "NG too large to_at",
			input:   &ListAlertRequest{ProjectId: 111, Severity: []string{"high"}, Description: "test_list_alert", FromAt: now.Unix(), ToAt: 253402268400},
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

func TestValidateGetAlertRequest(t *testing.T) {
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

func TestValidatePutAlertRequest(t *testing.T) {
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

func TestValidateDeleteAlertRequest(t *testing.T) {
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

func TestValidateListAlertHistoryRequest(t *testing.T) {
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
		{
			name:    "NG too small from_at",
			input:   &ListAlertHistoryRequest{ProjectId: 1001, AlertId: 1001, HistoryType: []string{"created"}, Severity: []string{"high"}, FromAt: -1, ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too large from_at",
			input:   &ListAlertHistoryRequest{ProjectId: 1001, AlertId: 1001, HistoryType: []string{"created"}, Severity: []string{"high"}, FromAt: 253402268400, ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too small to_at",
			input:   &ListAlertHistoryRequest{ProjectId: 1001, AlertId: 1001, HistoryType: []string{"created"}, Severity: []string{"high"}, FromAt: now.Unix(), ToAt: -1},
			wantErr: true,
		},
		{
			name:    "NG too large to_at",
			input:   &ListAlertHistoryRequest{ProjectId: 1001, AlertId: 1001, HistoryType: []string{"created"}, Severity: []string{"high"}, FromAt: now.Unix(), ToAt: 253402268400},
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

func TestValidateGetAlertHistoryRequest(t *testing.T) {
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

func TestValidatePutAlertHistoryRequest(t *testing.T) {
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

func TestValidateDeleteAlertHistoryRequest(t *testing.T) {
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

func TestValidateListRelAlertFindingRequest(t *testing.T) {
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
		{
			name:    "NG too small from_at",
			input:   &ListRelAlertFindingRequest{ProjectId: 1001, AlertId: 1001, FindingId: 1001, FromAt: -1, ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too large from_at",
			input:   &ListRelAlertFindingRequest{ProjectId: 1001, AlertId: 1001, FindingId: 1001, FromAt: 253402268400, ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too small to_at",
			input:   &ListRelAlertFindingRequest{ProjectId: 1001, AlertId: 1001, FindingId: 1001, FromAt: now.Unix(), ToAt: -1},
			wantErr: true,
		},
		{
			name:    "NG too large to_at",
			input:   &ListRelAlertFindingRequest{ProjectId: 1001, AlertId: 1001, FindingId: 1001, FromAt: now.Unix(), ToAt: 253402268400},
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

func TestValidateGetRelAlertFindingRequest(t *testing.T) {
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

func TestValidatePutRelAlertFindingRequest(t *testing.T) {
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

func TestValidateDeleteRelAlertFindingRequest(t *testing.T) {
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

func TestValidateListAlertConditionRequest(t *testing.T) {
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
		{
			name:    "NG too small from_at",
			input:   &ListAlertConditionRequest{ProjectId: 1001, Severity: []string{"high", "medium"}, Enabled: true, FromAt: -1, ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too large from_at",
			input:   &ListAlertConditionRequest{ProjectId: 1001, Severity: []string{"high", "medium"}, Enabled: true, FromAt: 253402268400, ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too small to_at",
			input:   &ListAlertConditionRequest{ProjectId: 1001, Severity: []string{"high", "medium"}, Enabled: true, FromAt: now.Unix(), ToAt: -1},
			wantErr: true,
		},
		{
			name:    "NG too large to_at",
			input:   &ListAlertConditionRequest{ProjectId: 1001, Severity: []string{"high", "medium"}, Enabled: true, FromAt: now.Unix(), ToAt: 253402268400},
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

func TestValidateGetAlertConditionRequest(t *testing.T) {
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

func TestValidatePutAlertConditionRequest(t *testing.T) {
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

func TestValidateDeleteAlertConditionRequest(t *testing.T) {
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

func TestValidateListAlertRuleRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *ListAlertRuleRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListAlertRuleRequest{ProjectId: 1001, FromScore: 0.0, ToScore: 1.0, FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListAlertRuleRequest{FromScore: 0.0, ToScore: 1.0, FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too small from_at",
			input:   &ListAlertRuleRequest{ProjectId: 1001, FromScore: 0.0, ToScore: 1.0, FromAt: -1, ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too large from_at",
			input:   &ListAlertRuleRequest{ProjectId: 1001, FromScore: 0.0, ToScore: 1.0, FromAt: 253402268400, ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too small to_at",
			input:   &ListAlertRuleRequest{ProjectId: 1001, FromScore: 0.0, ToScore: 1.0, FromAt: now.Unix(), ToAt: -1},
			wantErr: true,
		},
		{
			name:    "NG too large to_at",
			input:   &ListAlertRuleRequest{ProjectId: 1001, FromScore: 0.0, ToScore: 1.0, FromAt: now.Unix(), ToAt: 253402268400},
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

func TestValidateGetAlertRuleRequest(t *testing.T) {
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

func TestValidatePutAlertRuleRequest(t *testing.T) {
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
			name:    "NG Not Equal(project_id != rule.project_id)",
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

func TestValidateDeleteAlertRuleRequest(t *testing.T) {
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

func TestValidateListAlertCondRuleRequest(t *testing.T) {
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
		{
			name:    "NG too small from_at",
			input:   &ListAlertCondRuleRequest{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001, FromAt: -1, ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too large from_at",
			input:   &ListAlertCondRuleRequest{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001, FromAt: 253402268400, ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too small to_at",
			input:   &ListAlertCondRuleRequest{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001, FromAt: now.Unix(), ToAt: -1},
			wantErr: true,
		},
		{
			name:    "NG too large to_at",
			input:   &ListAlertCondRuleRequest{ProjectId: 1001, AlertConditionId: 1001, AlertRuleId: 1001, FromAt: now.Unix(), ToAt: 253402268400},
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

func TestValidateGetAlertCondRuleRequest(t *testing.T) {
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

func TestValidatePutAlertCondRuleRequest(t *testing.T) {
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

func TestValidateDeleteAlertCondRuleRequest(t *testing.T) {
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

func TestValidateListNotificationRequest(t *testing.T) {
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
		{
			name:    "NG too small from_at",
			input:   &ListNotificationRequest{ProjectId: 1001, Type: "test_notification", FromAt: -1, ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too large from_at",
			input:   &ListNotificationRequest{ProjectId: 1001, Type: "test_notification", FromAt: 253402268400, ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too small to_at",
			input:   &ListNotificationRequest{ProjectId: 1001, Type: "test_notification", FromAt: now.Unix(), ToAt: -1},
			wantErr: true,
		},
		{
			name:    "NG too large to_at",
			input:   &ListNotificationRequest{ProjectId: 1001, Type: "test_notification", FromAt: now.Unix(), ToAt: 253402268400},
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

func TestValidateGetNotificationRequest(t *testing.T) {
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

func TestValidatePutNotificationRequest(t *testing.T) {
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

func TestValidateDeleteNotificationRequest(t *testing.T) {
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

func TestValidateListAlertCondNotificationRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *ListAlertCondNotificationRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListAlertCondNotificationRequest{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001, FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListAlertCondNotificationRequest{AlertConditionId: 1001, NotificationId: 1001, FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too small from_at",
			input:   &ListAlertCondNotificationRequest{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001, FromAt: -1, ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too large from_at",
			input:   &ListAlertCondNotificationRequest{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001, FromAt: 253402268400, ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too small to_at",
			input:   &ListAlertCondNotificationRequest{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001, FromAt: now.Unix(), ToAt: -1},
			wantErr: true,
		},
		{
			name:    "NG too large to_at",
			input:   &ListAlertCondNotificationRequest{ProjectId: 1001, AlertConditionId: 1001, NotificationId: 1001, FromAt: now.Unix(), ToAt: 253402268400},
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

func TestValidateGetAlertCondNotificationRequest(t *testing.T) {
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

func TestValidatePutAlertCondNotificationRequest(t *testing.T) {
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

func TestValidateDeleteAlertCondNotificationRequest(t *testing.T) {
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

func TestValidateAnalyzeAlertRequest(t *testing.T) {
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
