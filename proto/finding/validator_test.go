package finding

import (
	"testing"
	"time"
)

func TestValidate_ListFindingRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *ListFindingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListFindingRequest{ProjectId: 111, DataSource: []string{"ds1", "ds2"}, ResourceName: []string{"rn1", "rn2"}, FromScore: 0.0, ToScore: 1.0, FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListFindingRequest{DataSource: []string{"ds1", "ds2"}, ResourceName: []string{"rn1", "rn2"}, FromScore: 0.0, ToScore: 1.0, FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too long resource_name",
			input:   &ListFindingRequest{ProjectId: 111, ResourceName: []string{"123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=12345678901234567890123456789012345678901234567890123456"}},
			wantErr: true,
		},
		{
			name:    "NG too long data_source",
			input:   &ListFindingRequest{ProjectId: 111, DataSource: []string{"12345678901234567890123456789012345678901234567890123456789012345"}},
			wantErr: true,
		},
		{
			name:    "NG small from_score",
			input:   &ListFindingRequest{ProjectId: 111, FromScore: -0.1},
			wantErr: true,
		},
		{
			name:    "NG big from_score",
			input:   &ListFindingRequest{ProjectId: 111, FromScore: 1.1},
			wantErr: true,
		},
		{
			name:    "NG small to_score",
			input:   &ListFindingRequest{ProjectId: 111, ToScore: -0.1},
			wantErr: true,
		},
		{
			name:    "NG big to_score",
			input:   &ListFindingRequest{ProjectId: 111, ToScore: 1.1},
			wantErr: true,
		},
		{
			name:    "NG small from_at",
			input:   &ListFindingRequest{ProjectId: 111, FromAt: -1},
			wantErr: true,
		},
		{
			name:    "NG big from_at",
			input:   &ListFindingRequest{ProjectId: 111, FromAt: 253402268400},
			wantErr: true,
		},
		{
			name:    "NG small to_at",
			input:   &ListFindingRequest{ProjectId: 111, ToAt: -1},
			wantErr: true,
		},
		{
			name:    "NG big to_at",
			input:   &ListFindingRequest{ProjectId: 111, ToAt: 253402268400},
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

func TestValidate_GetFindingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetFindingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetFindingRequest{ProjectId: 1, FindingId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetFindingRequest{FindingId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(finding_id)",
			input:   &GetFindingRequest{ProjectId: 1},
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

func TestValidate_PutFindingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutFindingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutFindingRequest{ProjectId: 1, Finding: &FindingForUpsert{DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1, OriginalScore: 1.0, OriginalMaxScore: 1.0}},
			wantErr: false,
		},
		{
			name:    "NG Required(finding)",
			input:   &PutFindingRequest{ProjectId: 999},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != tag.project_id)",
			input:   &PutFindingRequest{ProjectId: 999, Finding: &FindingForUpsert{DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1, OriginalScore: 1.0, OriginalMaxScore: 1.0}},
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

func TestValidate_DeleteFindingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteFindingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteFindingRequest{ProjectId: 1, FindingId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteFindingRequest{FindingId: 1001},
			wantErr: true,
		},
		{
			name:    "NG required(finding_id)",
			input:   &DeleteFindingRequest{ProjectId: 1},
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

func TestValidate_ListFindingTagRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListFindingTagRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListFindingTagRequest{ProjectId: 1, FindingId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListFindingTagRequest{FindingId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Required(finding_id)",
			input:   &ListFindingTagRequest{ProjectId: 1},
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

func TestValidate_TagFindingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *TagFindingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &TagFindingRequest{ProjectId: 1, Tag: &FindingTagForUpsert{FindingId: 1001, ProjectId: 1, TagKey: "k", TagValue: "v"}},
			wantErr: false,
		},
		{
			name:    "NG Required(tag)",
			input:   &TagFindingRequest{ProjectId: 999},
			wantErr: true,
		},
		{
			name:    "NG Required(project_id)",
			input:   &TagFindingRequest{Tag: &FindingTagForUpsert{FindingId: 1001, ProjectId: 1, TagKey: "k", TagValue: "v"}},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != tag.project_id)",
			input:   &TagFindingRequest{ProjectId: 999, Tag: &FindingTagForUpsert{FindingId: 1001, ProjectId: 1, TagKey: "k", TagValue: "v"}},
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

func TestValidate_UntagFindingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *UntagFindingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &UntagFindingRequest{ProjectId: 1, FindingTagId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &UntagFindingRequest{FindingTagId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Required(finding_tag_id)",
			input:   &UntagFindingRequest{ProjectId: 1},
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

func TestValidate_ListResourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListResourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListResourceRequest{ProjectId: 1},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListResourceRequest{},
			wantErr: true,
		},
		{
			name:    "NG Length(esource_name)",
			input:   &ListResourceRequest{ProjectId: 1, ResourceName: []string{"123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=12345678901234567890123456789012345678901234567890123456"}},
			wantErr: true,
		},
		{
			name:    "NG too small from_sum_score",
			input:   &ListResourceRequest{ProjectId: 1, FromSumScore: -0.1},
			wantErr: true,
		},
		{
			name:    "NG too small to_sum_score",
			input:   &ListResourceRequest{ProjectId: 1, ToSumScore: -0.1},
			wantErr: true,
		},
		{
			name:    "NG small from_at",
			input:   &ListResourceRequest{ProjectId: 1, FromAt: -1},
			wantErr: true,
		},
		{
			name:    "NG big from_at",
			input:   &ListResourceRequest{ProjectId: 1, FromAt: 253402268400},
			wantErr: true,
		},
		{
			name:    "NG small to_at",
			input:   &ListResourceRequest{ProjectId: 1, ToAt: -1},
			wantErr: true,
		},
		{
			name:    "NG big to_at",
			input:   &ListResourceRequest{ProjectId: 1, ToAt: 253402268400},
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

func TestValidate_GetResourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetResourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetResourceRequest{ProjectId: 1, ResourceId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetResourceRequest{ResourceId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Required(resource_id)",
			input:   &GetResourceRequest{ProjectId: 1},
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

func TestValidate_PutResourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutResourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutResourceRequest{ProjectId: 1, Resource: &ResourceForUpsert{ResourceName: "rn", ProjectId: 1}},
			wantErr: false,
		},
		{
			name:    "NG Required(resource)",
			input:   &PutResourceRequest{ProjectId: 999},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != tag.project_id)",
			input:   &PutResourceRequest{ProjectId: 999, Resource: &ResourceForUpsert{ResourceName: "rn", ProjectId: 1}},
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

func TestValidate_DeleteResourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteResourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteResourceRequest{ProjectId: 1, ResourceId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteResourceRequest{ResourceId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Required(resource_id)",
			input:   &DeleteResourceRequest{ProjectId: 1},
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

func TestValidate_ListResourceTagRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListResourceTagRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListResourceTagRequest{ProjectId: 1, ResourceId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListResourceTagRequest{ResourceId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Required(resource_id)",
			input:   &ListResourceTagRequest{ProjectId: 1},
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

func TestValidate_TagResourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *TagResourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &TagResourceRequest{ProjectId: 1, Tag: &ResourceTagForUpsert{ResourceId: 1001, ProjectId: 1, TagKey: "k", TagValue: "v"}},
			wantErr: false,
		},
		{
			name:    "NG Required(tag)",
			input:   &TagResourceRequest{ProjectId: 999},
			wantErr: true,
		},
		{
			name:    "NG Required(project_id)",
			input:   &TagResourceRequest{Tag: &ResourceTagForUpsert{ResourceId: 1001, ProjectId: 1, TagKey: "k", TagValue: "v"}},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != tag.project_id)",
			input:   &TagResourceRequest{ProjectId: 999, Tag: &ResourceTagForUpsert{ResourceId: 1001, ProjectId: 1, TagKey: "k", TagValue: "v"}},
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

func TestValidate_UntagResourceRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *UntagResourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &UntagResourceRequest{ProjectId: 1, ResourceTagId: 1001},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &UntagResourceRequest{ResourceTagId: 1001},
			wantErr: true,
		},
		{
			name:    "NG Required(resource_tag_id)",
			input:   &UntagResourceRequest{ProjectId: 1},
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

func TestValidate_FindingForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *FindingForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: false,
		},
		{
			name:    "NG too long Description",
			input:   &FindingForUpsert{Description: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=1", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG required DataSource",
			input:   &FindingForUpsert{Description: "desc", DataSource: "", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG too long DataSource",
			input:   &FindingForUpsert{Description: "desc", DataSource: "12345678901234567890123456789012345678901234567890123456789012345", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG required DataSourceId",
			input:   &FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG too long DataSourceId",
			input:   &FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=01234567890123456789012345678901234567890123456789123456", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG required resource name",
			input:   &FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG too long resource name",
			input:   &FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=12345678901234567890123456789012345678901234567890123456", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG nil OriginalScore",
			input:   &FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG too small OriginalScore",
			input:   &FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: -0.1, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG OriginalScore bigger than OriginalMaxScore",
			input:   &FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 100.01, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG nil OriginalMaxScore",
			input:   &FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG too small OriginalMaxScore",
			input:   &FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: -0.01, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG too big OriginalMaxScore",
			input:   &FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 999.991, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG invalid json Data",
			input:   &FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"`},
			wantErr: true,
		},
		{
			name:    "NG invalid json Data2",
			input:   &FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{key: value}`},
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

func TestValidate_FindingTagForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *FindingTagForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &FindingTagForUpsert{FindingId: 1001, TagKey: "key", TagValue: "value"},
			wantErr: false,
		},
		{
			name:    "NG required FindingId",
			input:   &FindingTagForUpsert{FindingId: 0, TagKey: "key", TagValue: "value"},
			wantErr: true,
		},
		{
			name:    "NG required TagKey",
			input:   &FindingTagForUpsert{FindingId: 1001, TagKey: "", TagValue: "value"},
			wantErr: true,
		},
		{
			name:    "NG too long TagKey",
			input:   &FindingTagForUpsert{FindingId: 1001, TagKey: "12345678901234567890123456789012345678901234567890123456789012345", TagValue: "value"},
			wantErr: true,
		},
		{
			name:    "NG too long TagValue",
			input:   &FindingTagForUpsert{FindingId: 1001, TagKey: "key", TagValue: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=1"},
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

func TestValidate_ResourceForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *ResourceForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ResourceForUpsert{ResourceName: "rn", ProjectId: 1001},
			wantErr: false,
		},
		{
			name:    "NG required ResourceName",
			input:   &ResourceForUpsert{ResourceName: "", ProjectId: 1001},
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

func TestValidate_ResourceTagForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *ResourceTagForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ResourceTagForUpsert{ResourceId: 1001, TagKey: "key", TagValue: "value"},
			wantErr: false,
		},
		{
			name:    "NG required FindingId",
			input:   &ResourceTagForUpsert{ResourceId: 0, TagKey: "key", TagValue: "value"},
			wantErr: true,
		},
		{
			name:    "NG required TagKey",
			input:   &ResourceTagForUpsert{ResourceId: 1001, TagKey: "", TagValue: "value"},
			wantErr: true,
		},
		{
			name:    "NG too long TagKey",
			input:   &ResourceTagForUpsert{ResourceId: 1001, TagKey: "12345678901234567890123456789012345678901234567890123456789012345", TagValue: "value"},
			wantErr: true,
		},
		{
			name:    "NG too long TagValue",
			input:   &ResourceTagForUpsert{ResourceId: 1001, TagKey: "key", TagValue: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=1"},
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
