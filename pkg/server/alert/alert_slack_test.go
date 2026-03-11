package alert

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/project"
	"github.com/jarcoal/httpmock"
	"github.com/slack-go/slack"
)

func TestSendSlackNotification(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://hogehoge.com", httpmock.NewStringResponder(200, "mocked"))
	httpmock.RegisterResponder("POST", "http://fugafuga.com", httpmock.NewErrorResponder(errors.New("Something Wrong")))
	testFindings := &findingDetail{}
	a := AlertService{}
	cases := []struct {
		name          string
		notifySetting string
		alert         *model.Alert
		project       *project.Project
		wantErr       bool
	}{
		{
			name:          "OK",
			notifySetting: `{"webhook_url":"http://hogehoge.com"}`,
			alert:         &model.Alert{},
			project:       &project.Project{},
			wantErr:       false,
		},
		{
			name:          "NG Json.Marshal Error",
			notifySetting: `{"webhook_url":http://hogehoge.com"}`,
			alert:         &model.Alert{},
			project:       &project.Project{},
			wantErr:       true,
		},
		{
			name:          "Warn webhook_url not set",
			notifySetting: `{}`,
			alert:         &model.Alert{},
			project:       &project.Project{},
			wantErr:       false,
		},
		{
			name:          "HTTP Error",
			notifySetting: `{"webhook_url":"http://fugafuga.com"}`,
			alert:         &model.Alert{},
			project:       &project.Project{},
			wantErr:       true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := a.sendSlackNotification(context.Background(), "unused", c.notifySetting, c.alert, c.project, &[]model.AlertRule{}, testFindings, LocaleEn)
			if (got != nil && !c.wantErr) || (got == nil && c.wantErr) {
				t.Fatalf("Unexpected error: %+v", got)
			}
		})
	}
}

func TestGenerateRuleList(t *testing.T) {
	cases := []struct {
		name  string
		input *[]model.AlertRule
		want  string
	}{
		{
			name: "1 line",
			input: &[]model.AlertRule{
				{AlertRuleID: 1, Name: "aaa"},
			},
			want: "- aaa",
		},
		{
			name: "Multi lines",
			input: &[]model.AlertRule{
				{AlertRuleID: 1, Name: "aaa"},
				{AlertRuleID: 2, Name: "bbb"},
				{AlertRuleID: 3, Name: "ccc"},
			},
			want: "- aaa\n- bbb\n- ccc",
		},
		{
			name:  "Nil input",
			input: nil,
			want:  "",
		},
		{
			name: "Too many rules(max=3)",
			input: &[]model.AlertRule{
				{AlertRuleID: 1, Name: "aaa"},
				{AlertRuleID: 2, Name: "bbb"},
				{AlertRuleID: 3, Name: "ccc"},
				{AlertRuleID: 4, Name: "ddd"},
				{AlertRuleID: 5, Name: "eee"},
				{AlertRuleID: 6, Name: "fff"},
				{AlertRuleID: 7, Name: "ggg"},
				{AlertRuleID: 8, Name: "hhh"},
				{AlertRuleID: 9, Name: "iii"},
				{AlertRuleID: 10, Name: "jjj"},
			},
			want: "- aaa\n- bbb\n- ccc\n- ...",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := generateRuleList(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected result: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetFindingAttachment(t *testing.T) {
	cases := []struct {
		name     string
		input    *findingDetail
		wantNum  int
		wantLast slack.AttachmentField
	}{
		{
			name: "without ai summary",
			input: &findingDetail{
				FindingCount: 1,
				Exampls: []*findingExample{{
					FindingID:    1,
					Description:  "desc",
					ResourceName: "resource",
					DataSource:   "ds",
					Score:        0.9,
					Tags:         []string{"tag1"},
				}},
			},
			wantNum: 4,
			wantLast: slack.AttachmentField{
				Title: "Tags",
				Value: "`tag1`",
			},
		},
		{
			name: "with ai summary",
			input: &findingDetail{
				FindingCount: 1,
				Exampls: []*findingExample{{
					FindingID:    1,
					Description:  "desc",
					ResourceName: "resource",
					DataSource:   "ds",
					Score:        0.9,
					Tags:         []string{"tag1"},
					AISummary:    "summary",
				}},
			},
			wantNum: 5,
			wantLast: slack.AttachmentField{
				Title: "AI Summary",
				Value: "summary",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getFindingAttachment("https://example.com", 1, c.input, LocaleEn)
			if len(got) != 1 {
				t.Fatalf("Unexpected attachment count: got=%d", len(got))
			}
			if len(got[0].Fields) != c.wantNum {
				t.Fatalf("Unexpected field count: got=%d want=%d", len(got[0].Fields), c.wantNum)
			}
			lastField := got[0].Fields[len(got[0].Fields)-1]
			if !reflect.DeepEqual(lastField, c.wantLast) {
				t.Fatalf("Unexpected last field: got=%+v want=%+v", lastField, c.wantLast)
			}
		})
	}
}
