package alert

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

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
			got := a.sendSlackNotification(context.Background(), "unused", c.notifySetting, "", c.alert, c.project, &[]model.AlertRule{}, testFindings, LocaleEn)
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
		name      string
		input     *findingDetail
		wantNum   int
		wantFirst slack.AttachmentField
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
			wantFirst: slack.AttachmentField{
				Value: "<https://example.com/finding/finding?project_id=1&finding_id=1&from_score=0&status=1&from=slack|View alert details in RISKEN>",
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
					AISummary:    `{"blocks":[{"type":"text","text":"summary"}]}`,
				}},
			},
			wantNum: 5,
			wantFirst: slack.AttachmentField{
				Title: "AI Summary",
				Value: "summary",
			},
		},
		{
			name: "with ai summary markdown link",
			input: &findingDetail{
				FindingCount: 1,
				Exampls: []*findingExample{{
					FindingID:    1,
					Description:  "desc",
					ResourceName: "resource",
					DataSource:   "ds",
					Score:        0.9,
					Tags:         []string{"tag1"},
					AISummary:    `{"blocks":[{"type":"text","text":"確認してください"},{"type":"link","label":"GitHubリンク","url":"https://github.com/ca-risken/security-review-test/blob/34d724422060a79eaa04a42b278cb7dab10b75d7/test/review-code/main.go#L30-L30"}]}`,
				}},
			},
			wantNum: 5,
			wantFirst: slack.AttachmentField{
				Title: "AI Summary",
				Value: "確認してください\n<https://github.com/ca-risken/security-review-test/blob/34d724422060a79eaa04a42b278cb7dab10b75d7/test/review-code/main.go#L30-L30|GitHubリンク>",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := &AlertService{}
			got := a.getFindingAttachment(context.Background(), "https://example.com", 1, c.input, LocaleEn)
			if len(got) != 1 {
				t.Fatalf("Unexpected attachment count: got=%d", len(got))
			}
			if len(got[0].Fields) != c.wantNum {
				t.Fatalf("Unexpected field count: got=%d want=%d", len(got[0].Fields), c.wantNum)
			}
			firstField := got[0].Fields[0]
			if !reflect.DeepEqual(firstField, c.wantFirst) {
				t.Fatalf("Unexpected first field: got=%+v want=%+v", firstField, c.wantFirst)
			}
		})
	}
}

func TestGetFindingAttachmentUsesLocaleAwareRISKENLinkLabel(t *testing.T) {
	a := &AlertService{}
	got := a.getFindingAttachment(context.Background(), "https://example.com", 1, &findingDetail{
		FindingCount: 1,
		Exampls: []*findingExample{{
			FindingID:    1,
			Description:  "desc",
			ResourceName: "resource",
			DataSource:   "ds",
			Score:        0.9,
			Tags:         []string{"tag1"},
			AISummary:    `{"blocks":[{"type":"text","text":"summary"}]}`,
		}},
	}, LocaleJa)

	if len(got) != 1 {
		t.Fatalf("Unexpected attachment count: got=%d", len(got))
	}
	if len(got[0].Fields) < 2 {
		t.Fatalf("Unexpected field count: got=%d", len(got[0].Fields))
	}
	want := "<https://example.com/finding/finding?project_id=1&finding_id=1&from_score=0&status=1&from=slack|アラートの詳細をRISKENで確認>"
	if got[0].Fields[1].Value != want {
		t.Fatalf("Unexpected RISKEN link label: got=%q want=%q", got[0].Fields[1].Value, want)
	}

	gotEn := a.getFindingAttachment(context.Background(), "https://example.com", 1, &findingDetail{
		FindingCount: 1,
		Exampls: []*findingExample{{
			FindingID:    1,
			Description:  "desc",
			ResourceName: "resource",
			DataSource:   "ds",
			Score:        0.9,
			Tags:         []string{"tag1"},
			AISummary:    `{"blocks":[{"type":"text","text":"summary"}]}`,
		}},
	}, LocaleEn)
	wantEn := "<https://example.com/finding/finding?project_id=1&finding_id=1&from_score=0&status=1&from=slack|View alert details in RISKEN>"
	if gotEn[0].Fields[1].Value != wantEn {
		t.Fatalf("Unexpected RISKEN link label for en: got=%q want=%q", gotEn[0].Fields[1].Value, wantEn)
	}

	gotDefault := a.getFindingAttachment(context.Background(), "https://example.com", 1, &findingDetail{
		FindingCount: 1,
		Exampls: []*findingExample{{
			FindingID:    1,
			Description:  "desc",
			ResourceName: "resource",
			DataSource:   "ds",
			Score:        0.9,
			Tags:         []string{"tag1"},
			AISummary:    `{"blocks":[{"type":"text","text":"summary"}]}`,
		}},
	}, "")
	if gotDefault[0].Fields[1].Value != wantEn {
		t.Fatalf("Unexpected RISKEN link label for default locale: got=%q want=%q", gotDefault[0].Fields[1].Value, wantEn)
	}
}

func TestRenderAlertAISummary(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "payload",
			input: `{"blocks":[{"type":"text","text":"summary"},{"type":"link","label":"GitHub","url":"https://example.com"}]}`,
			want:  "summary\n<https://example.com|GitHub>",
		},
		{
			name:  "escape mrkdwn in text and label",
			input: `{"blocks":[{"type":"text","text":"notify <!here> & review <this>"},{"type":"link","label":"Git|Hub > docs","url":"https://example.com"}]}`,
			want:  "notify &lt;!here&gt; &amp; review &lt;this&gt;\n<https://example.com|Git¦Hub &gt; docs>",
		},
		{
			name:  "delimit url before Japanese closing parenthesis",
			input: `{"blocks":[{"type":"text","text":"google.golang.org/grpc（http://google.golang.org/grpc）の既知の脆弱性を検知しました。"}]}`,
			want:  "google.golang.org/grpc（<http://google.golang.org/grpc>）の既知の脆弱性を検知しました。",
		},
		{
			name:  "delimit url before ASCII closing parenthesis",
			input: `{"blocks":[{"type":"text","text":"see (https://example.com/path) now"}]}`,
			want:  "see (<https://example.com/path>) now",
		},
		{
			name:  "keep balanced parentheses in url",
			input: `{"blocks":[{"type":"text","text":"see https://example.com/wiki/Function_(math) now"}]}`,
			want:  "see <https://example.com/wiki/Function_(math)> now",
		},
		{
			name:  "delimit url before trailing punctuation",
			input: `{"blocks":[{"type":"text","text":"see https://example.com/path., next https://example.com/a]; done"}]}`,
			want:  "see <https://example.com/path>., next <https://example.com/a>]; done",
		},
		{
			name:  "escape ampersand in url",
			input: `{"blocks":[{"type":"text","text":"see https://example.com/path?x=1&y=2 now"}]}`,
			want:  "see <https://example.com/path?x=1&amp;y=2> now",
		},
		{
			name:  "delimit punctuation before unbalanced closing parenthesis",
			input: `{"blocks":[{"type":"text","text":"see https://example.com/path.) now"}]}`,
			want:  "see <https://example.com/path>.) now",
		},
		{
			name:  "do not link url without host",
			input: `{"blocks":[{"type":"text","text":"end https://]. now"}]}`,
			want:  "end https://]. now",
		},
		{
			name:  "normalize Japanese path",
			input: `{"blocks":[{"type":"text","text":"see https://example.com/脆弱性。確認してください。"}]}`,
			want:  "see <https://example.com/%E8%84%86%E5%BC%B1%E6%80%A7>。確認してください。",
		},
		{
			name:  "keep IPv6 host brackets",
			input: `{"blocks":[{"type":"text","text":"see (http://[2001:db8::1]) now"}]}`,
			want:  "see (<http://[2001:db8::1]>) now",
		},
		{
			name:  "delimit url before full width separators",
			input: `{"blocks":[{"type":"text","text":"https://example.com/a？ https://example.com/b！ https://example.com/c． https://example.com/d・ https://example.com/e　next"}]}`,
			want:  "<https://example.com/a>？ <https://example.com/b>！ <https://example.com/c>． <https://example.com/d>・ <https://example.com/e>　next",
		},
		{
			name:  "drop unsafe url in link block",
			input: `{"blocks":[{"type":"text","text":"summary"},{"type":"link","label":"malicious","url":"https://example.com/a><!channel>"}]}`,
			want:  "summary",
		},
		{
			name:  "drop url with whitespace in link block",
			input: "{\"blocks\":[{\"type\":\"text\",\"text\":\"summary\"},{\"type\":\"link\",\"label\":\"docs\",\"url\":\"https://example.com/path\\nnext\"}]}",
			want:  "summary",
		},
		{
			name:  "payload wrapped by json code fence",
			input: "```json\n{\"blocks\":[{\"type\":\"text\",\"text\":\"summary\"}]}\n```",
			want:  "summary",
		},
		{
			name:  "payload wrapped by plain code fence",
			input: "```\n{\"blocks\":[{\"type\":\"text\",\"text\":\"summary\"}]}\n```",
			want:  "summary",
		},
		{
			name:  "invalid payload",
			input: "確認: [GitHubリンク](https://example.com/path)",
			want:  "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := renderAlertAISummary(c.input)
			if got != c.want {
				t.Fatalf("Unexpected rendered value: got=%q want=%q", got, c.want)
			}
		})
	}
}

func TestBuildSlackAttachments(t *testing.T) {
	alert := &model.Alert{
		Description: "alert-desc",
		Severity:    "high",
	}
	project := &project.Project{
		ProjectId: 1,
		Name:      "project-name",
	}
	rules := &[]model.AlertRule{
		{Name: "rule-1"},
	}
	findings := &findingDetail{
		FindingCount: 1,
		Exampls: []*findingExample{{
			FindingID:    10,
			Description:  "finding-desc",
			ResourceName: "resource-1",
			DataSource:   "ds-1",
			Score:        0.9,
			Tags:         []string{"tag-1"},
		}},
	}

	cases := []struct {
		name                 string
		organizationName     string
		wantOrganizationName string
		wantOrganizationRow  bool
	}{
		{name: "project notification does not show organization"},
		{
			name:                 "organization notification shows escaped source",
			organizationName:     "org <!channel> & <https://example.com|link>",
			wantOrganizationName: "org &lt;!channel&gt; &amp; &lt;https://example.com|link&gt;",
			wantOrganizationRow:  true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := (&AlertService{}).buildSlackAttachments(context.Background(), "https://example.com", c.organizationName, alert, project, rules, findings, LocaleJa)

			if len(got) != 2 {
				t.Fatalf("Unexpected attachment count: got=%d want=2", len(got))
			}
			if !strings.Contains(got[0].Fields[0].Value, "/finding/finding?project_id=1&finding_id=10") {
				t.Fatalf("First attachment should be finding block: got=%+v", got[0].Fields[0].Value)
			}
			if !strings.Contains(got[1].Fields[0].Value, "alert-desc") {
				t.Fatalf("Last attachment should be alert block: got=%+v", got[1].Fields[0].Value)
			}
			organizationRows := 0
			for _, field := range got[1].Fields {
				if field.Title == "🏢 Organization" {
					organizationRows++
					if field.Value != c.wantOrganizationName {
						t.Fatalf("Unexpected organization name: got=%q want=%q", field.Value, c.wantOrganizationName)
					}
				}
			}
			if (organizationRows == 1) != c.wantOrganizationRow {
				t.Fatalf("Unexpected organization row count: got=%d", organizationRows)
			}
		})
	}
}

func TestGetFindingAttachmentAddsActionButtons(t *testing.T) {
	a := &AlertService{slackActionSigningSecret: "test-secret"}
	got := a.getFindingAttachment(context.Background(), "https://example.com", 1, &findingDetail{
		FindingCount: 1,
		Exampls: []*findingExample{{
			FindingID:    10,
			Description:  "finding-desc",
			ResourceName: "resource-1",
			DataSource:   "ds-1",
			Score:        0.9,
			Tags:         []string{"tag-1"},
		}},
	}, LocaleJa)

	if len(got) != 1 {
		t.Fatalf("Unexpected attachment count: got=%d want=1", len(got))
	}
	if len(got[0].Actions) != 2 {
		t.Fatalf("Unexpected action count: got=%d want=2", len(got[0].Actions))
	}
	if got[0].Actions[0].Name != "pend_finding" {
		t.Fatalf("Unexpected first action: got=%s want=pend_finding", got[0].Actions[0].Name)
	}
	if got[0].Actions[0].Text != "PEND" {
		t.Fatalf("Unexpected first action label: got=%s", got[0].Actions[0].Text)
	}
	if got[0].Actions[0].URL != "" {
		t.Fatalf("Unexpected first action URL: got=%s", got[0].Actions[0].URL)
	}
	var pendPayload slackActionPayload
	if err := json.Unmarshal([]byte(got[0].Actions[0].Value), &pendPayload); err != nil {
		t.Fatalf("Unexpected first action value: %s", got[0].Actions[0].Value)
	}
	if pendPayload.Action != slackActionButtonPend || pendPayload.ProjectID != 1 || pendPayload.FindingID != 10 {
		t.Fatalf("Unexpected first action payload: got=%+v", pendPayload)
	}
	if !validSlackActionPayload(pendPayload, "test-secret", time.Now()) {
		t.Fatalf("Unexpected first action payload signature: got=%+v", pendPayload)
	}
	if got[0].Actions[1].Name != "archive_finding" {
		t.Fatalf("Unexpected second action: got=%s want=archive_finding", got[0].Actions[1].Name)
	}
	if got[0].Actions[1].Text != "Archive" {
		t.Fatalf("Unexpected second action label: got=%s", got[0].Actions[1].Text)
	}
	if got[0].Actions[1].URL != "" {
		t.Fatalf("Unexpected second action URL: got=%s", got[0].Actions[1].URL)
	}
	var archivePayload slackActionPayload
	if err := json.Unmarshal([]byte(got[0].Actions[1].Value), &archivePayload); err != nil {
		t.Fatalf("Unexpected second action value: %s", got[0].Actions[1].Value)
	}
	if archivePayload.Action != slackActionButtonArchive || archivePayload.ProjectID != 1 || archivePayload.FindingID != 10 {
		t.Fatalf("Unexpected second action payload: got=%+v", archivePayload)
	}
	if !validSlackActionPayload(archivePayload, "test-secret", time.Now()) {
		t.Fatalf("Unexpected second action payload signature: got=%+v", archivePayload)
	}
}

func TestGetFindingAttachmentSkipsActionButtonsWithoutSigningSecret(t *testing.T) {
	got := (&AlertService{}).getFindingAttachment(context.Background(), "https://example.com", 1, &findingDetail{
		FindingCount: 1,
		Exampls: []*findingExample{{
			FindingID:    10,
			Description:  "finding-desc",
			ResourceName: "resource-1",
			DataSource:   "ds-1",
			Score:        0.9,
			Tags:         []string{"tag-1"},
		}},
	}, LocaleJa)

	if len(got) != 1 {
		t.Fatalf("Unexpected attachment count: got=%d want=1", len(got))
	}
	if len(got[0].Actions) != 0 {
		t.Fatalf("Unexpected action count: got=%d want=0", len(got[0].Actions))
	}
}

func TestBuildSlackActionPayloadValue(t *testing.T) {
	now := time.Unix(1779721200, 0)
	got, err := buildSlackActionPayloadValue(slackActionButtonPend, 1001, 123456, "secret", now)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	var payload slackActionPayload
	if err := json.Unmarshal([]byte(got), &payload); err != nil {
		t.Fatalf("Unexpected payload json: %v", err)
	}
	if payload.Action != slackActionButtonPend {
		t.Fatalf("Unexpected action: got=%s", payload.Action)
	}
	if payload.ProjectID != 1001 {
		t.Fatalf("Unexpected project_id: got=%d", payload.ProjectID)
	}
	if payload.FindingID != 123456 {
		t.Fatalf("Unexpected finding_id: got=%d", payload.FindingID)
	}
	if payload.IssuedAt != now.Unix() {
		t.Fatalf("Unexpected issued_at: got=%d", payload.IssuedAt)
	}
	if payload.ExpiresAt != now.Add(slackActionPayloadTTL).Unix() {
		t.Fatalf("Unexpected expires_at: got=%d", payload.ExpiresAt)
	}
	if !validSlackActionPayload(payload, "secret", now) {
		t.Fatalf("Unexpected signature: got=%s", payload.Signature)
	}
}

func TestValidSlackActionPayloadRejectsExpiredPayload(t *testing.T) {
	now := time.Unix(1779721200, 0)
	got, err := buildSlackActionPayloadValue(slackActionButtonPend, 1001, 123456, "secret", now)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	var payload slackActionPayload
	if err := json.Unmarshal([]byte(got), &payload); err != nil {
		t.Fatalf("Unexpected payload json: %v", err)
	}
	if validSlackActionPayload(payload, "secret", now.Add(slackActionPayloadTTL+time.Second)) {
		t.Fatal("Expected expired payload to be rejected")
	}
}
