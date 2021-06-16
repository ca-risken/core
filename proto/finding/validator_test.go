package finding

import (
	"testing"
	"time"
)

const (
	len65string  string = "12345678901234567890123456789012345678901234567890123456789012345"
	len129string string = "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=12345678901234567890123456789"
	len201string string = "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=1"
	len256string string = "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=12345678901234567890123456789012345678901234567890123456"
	maxLimit     int32  = 200
)

func TestValidate_ListFindingRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   *ListFindingRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &ListFindingRequest{ProjectId: 1, DataSource: []string{"ds1", "ds2"}, ResourceName: []string{"rn1", "rn2"}, FromScore: 0.0, ToScore: 1.0, FromAt: now.Unix(), ToAt: now.Unix(), Sort: "finding_id", Direction: "asc", Offset: 0, Limit: maxLimit},
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListFindingRequest{DataSource: []string{"ds1", "ds2"}, ResourceName: []string{"rn1", "rn2"}, FromScore: 0.0, ToScore: 1.0, FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: true,
		},
		{
			name:    "NG too long resource_name",
			input:   &ListFindingRequest{ProjectId: 1, ResourceName: []string{len256string}},
			wantErr: true,
		},
		{
			name:    "NG too long data_source",
			input:   &ListFindingRequest{ProjectId: 1, DataSource: []string{len65string}},
			wantErr: true,
		},
		{
			name:    "NG small from_score",
			input:   &ListFindingRequest{ProjectId: 1, FromScore: -0.1},
			wantErr: true,
		},
		{
			name:    "NG big from_score",
			input:   &ListFindingRequest{ProjectId: 1, FromScore: 1.1},
			wantErr: true,
		},
		{
			name:    "NG small to_score",
			input:   &ListFindingRequest{ProjectId: 1, ToScore: -0.1},
			wantErr: true,
		},
		{
			name:    "NG big to_score",
			input:   &ListFindingRequest{ProjectId: 1, ToScore: 1.1},
			wantErr: true,
		},
		{
			name:    "NG small from_at",
			input:   &ListFindingRequest{ProjectId: 1, FromAt: -1},
			wantErr: true,
		},
		{
			name:    "NG big from_at",
			input:   &ListFindingRequest{ProjectId: 1, FromAt: 253402268400},
			wantErr: true,
		},
		{
			name:    "NG small to_at",
			input:   &ListFindingRequest{ProjectId: 1, ToAt: -1},
			wantErr: true,
		},
		{
			name:    "NG big to_at",
			input:   &ListFindingRequest{ProjectId: 1, ToAt: 253402268400},
			wantErr: true,
		},
		{
			name:    "NG too long tag",
			input:   &ListFindingRequest{ProjectId: 1, Tag: []string{len65string}},
			wantErr: true,
		},
		{
			name:    "NG Invalid sort",
			input:   &ListFindingRequest{ProjectId: 1, Sort: "unknown_key"},
			wantErr: true,
		},
		{
			name:    "NG Invalid sort direction",
			input:   &ListFindingRequest{ProjectId: 1, Direction: "unknown_direction"},
			wantErr: true,
		},
		{
			name:    "NG Min offset",
			input:   &ListFindingRequest{ProjectId: 1, Offset: -1},
			wantErr: true,
		},
		{
			name:    "NG Min limit",
			input:   &ListFindingRequest{ProjectId: 1, Limit: -1},
			wantErr: true,
		},
		{
			name:    "NG Max limit",
			input:   &ListFindingRequest{ProjectId: 1, Limit: maxLimit + 1},
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
			input:   &ListFindingTagRequest{ProjectId: 1, FindingId: 1, Sort: "finding_tag_id", Direction: "desc", Offset: 0, Limit: maxLimit},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListFindingTagRequest{FindingId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(finding_id)",
			input:   &ListFindingTagRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG Invalid sort",
			input:   &ListFindingTagRequest{ProjectId: 1, Sort: "unknown_key"},
			wantErr: true,
		},
		{
			name:    "NG Invalid sort direction",
			input:   &ListFindingTagRequest{ProjectId: 1, Direction: "unknown_direction"},
			wantErr: true,
		},
		{
			name:    "NG Min offset",
			input:   &ListFindingTagRequest{ProjectId: 1, Offset: -1},
			wantErr: true,
		},
		{
			name:    "NG Min limit",
			input:   &ListFindingTagRequest{ProjectId: 1, Limit: -1},
			wantErr: true,
		},
		{
			name:    "NG Max limit",
			input:   &ListFindingTagRequest{ProjectId: 1, Limit: maxLimit + 1},
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

func TestValidate_ListFindingTagNameRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListFindingTagNameRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &ListFindingTagNameRequest{ProjectId: 1, FromAt: 0, ToAt: 0, Sort: "finding_tag_id", Direction: "desc", Offset: 0, Limit: maxLimit},
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListFindingTagNameRequest{FromAt: 0, ToAt: 0},
			wantErr: true,
		},
		{
			name:    "NG small from_at",
			input:   &ListFindingTagNameRequest{ProjectId: 1, FromAt: -1},
			wantErr: true,
		},
		{
			name:    "NG big from_at",
			input:   &ListFindingTagNameRequest{ProjectId: 1, FromAt: 253402268400},
			wantErr: true,
		},
		{
			name:    "NG small to_at",
			input:   &ListFindingTagNameRequest{ProjectId: 1, ToAt: -1},
			wantErr: true,
		},
		{
			name:    "NG big to_at",
			input:   &ListFindingTagNameRequest{ProjectId: 1, ToAt: 253402268400},
			wantErr: true,
		},
		{
			name:    "NG Invalid sort",
			input:   &ListFindingTagNameRequest{ProjectId: 1, Sort: "unknown_key"},
			wantErr: true,
		},
		{
			name:    "NG Invalid sort direction",
			input:   &ListFindingTagNameRequest{ProjectId: 1, Direction: "unknown_direction"},
			wantErr: true,
		},
		{
			name:    "NG Min offset",
			input:   &ListFindingTagNameRequest{ProjectId: 1, Offset: -1},
			wantErr: true,
		},
		{
			name:    "NG Min limit",
			input:   &ListFindingTagNameRequest{ProjectId: 1, Limit: -1},
			wantErr: true,
		},
		{
			name:    "NG Max limit",
			input:   &ListFindingTagNameRequest{ProjectId: 1, Limit: maxLimit + 1},
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
			input:   &TagFindingRequest{ProjectId: 1, Tag: &FindingTagForUpsert{FindingId: 1001, ProjectId: 1, Tag: "tag"}},
			wantErr: false,
		},
		{
			name:    "NG Required(tag)",
			input:   &TagFindingRequest{ProjectId: 999},
			wantErr: true,
		},
		{
			name:    "NG Required(project_id)",
			input:   &TagFindingRequest{Tag: &FindingTagForUpsert{FindingId: 1001, ProjectId: 1, Tag: "tag"}},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != tag.project_id)",
			input:   &TagFindingRequest{ProjectId: 999, Tag: &FindingTagForUpsert{FindingId: 1001, ProjectId: 1, Tag: "tag"}},
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
			input:   &ListResourceRequest{ProjectId: 1, Sort: "resource_id", Direction: "desc", Offset: 0, Limit: maxLimit},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListResourceRequest{},
			wantErr: true,
		},
		{
			name:    "NG Length(resource_name)",
			input:   &ListResourceRequest{ProjectId: 1, ResourceName: []string{len256string}},
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
		{
			name:    "NG Length(tag)",
			input:   &ListResourceRequest{ProjectId: 1, Tag: []string{len65string}},
			wantErr: true,
		},
		{
			name:    "NG Invalid sort",
			input:   &ListResourceRequest{ProjectId: 1, Sort: "unknown_key"},
			wantErr: true,
		},
		{
			name:    "NG Invalid sort direction",
			input:   &ListResourceRequest{ProjectId: 1, Direction: "unknown_direction"},
			wantErr: true,
		},
		{
			name:    "NG Min offset",
			input:   &ListResourceRequest{ProjectId: 1, Offset: -1},
			wantErr: true,
		},
		{
			name:    "NG Min limit",
			input:   &ListResourceRequest{ProjectId: 1, Limit: -1},
			wantErr: true,
		},
		{
			name:    "NG Max limit",
			input:   &ListResourceRequest{ProjectId: 1, Limit: maxLimit + 1},
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
			input:   &ListResourceTagRequest{ProjectId: 1, ResourceId: 1001, Sort: "resource_tag_id", Direction: "desc", Offset: 0, Limit: maxLimit},
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
		{
			name:    "NG Invalid sort",
			input:   &ListResourceTagRequest{ProjectId: 1, Sort: "unknown_key"},
			wantErr: true,
		},
		{
			name:    "NG Invalid sort direction",
			input:   &ListResourceTagRequest{ProjectId: 1, Direction: "unknown_direction"},
			wantErr: true,
		},
		{
			name:    "NG Min offset",
			input:   &ListResourceTagRequest{ProjectId: 1, Offset: -1},
			wantErr: true,
		},
		{
			name:    "NG Min limit",
			input:   &ListResourceTagRequest{ProjectId: 1, Limit: -1},
			wantErr: true,
		},
		{
			name:    "NG Max limit",
			input:   &ListResourceTagRequest{ProjectId: 1, Limit: maxLimit + 1},
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

func TestValidate_ListResourceTagNameRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListResourceTagNameRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &ListResourceTagNameRequest{ProjectId: 1, FromAt: 0, ToAt: 0, Sort: "resource_tag_id", Direction: "desc", Offset: 0, Limit: maxLimit},
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListResourceTagNameRequest{FromAt: 0, ToAt: 0},
			wantErr: true,
		},
		{
			name:    "NG small from_at",
			input:   &ListResourceTagNameRequest{ProjectId: 1, FromAt: -1},
			wantErr: true,
		},
		{
			name:    "NG big from_at",
			input:   &ListResourceTagNameRequest{ProjectId: 1, FromAt: 253402268400},
			wantErr: true,
		},
		{
			name:    "NG small to_at",
			input:   &ListResourceTagNameRequest{ProjectId: 1, ToAt: -1},
			wantErr: true,
		},
		{
			name:    "NG big to_at",
			input:   &ListResourceTagNameRequest{ProjectId: 1, ToAt: 253402268400},
			wantErr: true,
		},
		{
			name:    "NG Invalid sort",
			input:   &ListResourceTagNameRequest{ProjectId: 1, Sort: "unknown_key"},
			wantErr: true,
		},
		{
			name:    "NG Invalid sort direction",
			input:   &ListResourceTagNameRequest{ProjectId: 1, Direction: "unknown_direction"},
			wantErr: true,
		},
		{
			name:    "NG Min offset",
			input:   &ListResourceTagNameRequest{ProjectId: 1, Offset: -1},
			wantErr: true,
		},
		{
			name:    "NG Min limit",
			input:   &ListResourceTagNameRequest{ProjectId: 1, Limit: -1},
			wantErr: true,
		},
		{
			name:    "NG Max limit",
			input:   &ListResourceTagNameRequest{ProjectId: 1, Limit: maxLimit + 1},
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
			input:   &TagResourceRequest{ProjectId: 1, Tag: &ResourceTagForUpsert{ResourceId: 1001, ProjectId: 1, Tag: "tag"}},
			wantErr: false,
		},
		{
			name:    "NG Required(tag)",
			input:   &TagResourceRequest{ProjectId: 999},
			wantErr: true,
		},
		{
			name:    "NG Required(project_id)",
			input:   &TagResourceRequest{Tag: &ResourceTagForUpsert{ResourceId: 1001, ProjectId: 1, Tag: "tag"}},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id != tag.project_id)",
			input:   &TagResourceRequest{ProjectId: 999, Tag: &ResourceTagForUpsert{ResourceId: 1001, ProjectId: 1, Tag: "tag"}},
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

func TestValidate_GetPendFindingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetPendFindingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetPendFindingRequest{ProjectId: 1, FindingId: 1},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetPendFindingRequest{FindingId: 1},
			wantErr: true,
		},
		{
			name:    "NG required(finding_id)",
			input:   &GetPendFindingRequest{ProjectId: 1},
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

func TestValidate_PutPendFindingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutPendFindingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutPendFindingRequest{ProjectId: 1, PendFinding: &PendFindingForUpsert{FindingId: 1, ProjectId: 1}},
			wantErr: false,
		},
		{
			name:    "NG Required(pend_finding)",
			input:   &PutPendFindingRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id)",
			input:   &PutPendFindingRequest{ProjectId: 999, PendFinding: &PendFindingForUpsert{FindingId: 1, ProjectId: 1}},
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

func TestValidate_DeletePendFindingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeletePendFindingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeletePendFindingRequest{ProjectId: 1, FindingId: 1},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeletePendFindingRequest{FindingId: 1},
			wantErr: true,
		},
		{
			name:    "NG required(finding_id)",
			input:   &DeletePendFindingRequest{ProjectId: 1},
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

func TestValidate_ListFindingSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListFindingSettingRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &ListFindingSettingRequest{ProjectId: 1, Status: FindingSettingStatus_SETTING_ACTIVE},
		},
		{
			name:  "OK status unknown",
			input: &ListFindingSettingRequest{ProjectId: 1, Status: FindingSettingStatus_SETTING_UNKNOWN},
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListFindingSettingRequest{Status: FindingSettingStatus_SETTING_ACTIVE},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_GetFindingSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetFindingSettingRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &GetFindingSettingRequest{ProjectId: 1, FindingSettingId: 1, Status: FindingSettingStatus_SETTING_ACTIVE},
		},
		{
			name:  "OK status unknown",
			input: &GetFindingSettingRequest{ProjectId: 1, FindingSettingId: 1, Status: FindingSettingStatus_SETTING_UNKNOWN},
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetFindingSettingRequest{FindingSettingId: 1, Status: FindingSettingStatus_SETTING_ACTIVE},
			wantErr: true,
		},
		{
			name:    "NG Required(finding_setting_id)",
			input:   &GetFindingSettingRequest{ProjectId: 1, Status: FindingSettingStatus_SETTING_ACTIVE},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
			} else if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}

func TestValidate_PutFindingSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutFindingSettingRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &PutFindingSettingRequest{ProjectId: 1, FindingSetting: &FindingSettingForUpsert{ProjectId: 1, ResourceName: "rn", Setting: `{"k":"v"}`, Status: FindingSettingStatus_SETTING_ACTIVE}},
		},
		{
			name:    "NG Required(finding_setting)",
			input:   &PutFindingSettingRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG Not Equal(project_id)",
			input:   &PutFindingSettingRequest{ProjectId: 999, FindingSetting: &FindingSettingForUpsert{ProjectId: 1, ResourceName: "rn", Setting: `{"k":"v"}`, Status: FindingSettingStatus_SETTING_ACTIVE}},
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

func TestValidate_DeleteFindingSettingRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteFindingSettingRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &DeleteFindingSettingRequest{ProjectId: 1, FindingSettingId: 1},
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteFindingSettingRequest{FindingSettingId: 1},
			wantErr: true,
		},
		{
			name:    "NG required(finding_setting_id)",
			input:   &DeleteFindingSettingRequest{ProjectId: 1},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.Validate()
			if c.wantErr && err == nil {
				t.Fatal("Unexpected no error")
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
			input:   &FindingForUpsert{Description: len201string, DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG required DataSource",
			input:   &FindingForUpsert{Description: "desc", DataSource: "", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG too long DataSource",
			input:   &FindingForUpsert{Description: "desc", DataSource: len65string, DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG required DataSourceId",
			input:   &FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG too long DataSourceId",
			input:   &FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: len256string, ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG required resource name",
			input:   &FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG too long resource name",
			input:   &FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: len256string, ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
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
			input:   &FindingTagForUpsert{FindingId: 1001, Tag: "tag"},
			wantErr: false,
		},
		{
			name:    "NG required FindingId",
			input:   &FindingTagForUpsert{FindingId: 0, Tag: "tag"},
			wantErr: true,
		},
		{
			name:    "NG required Tag",
			input:   &FindingTagForUpsert{FindingId: 1001},
			wantErr: true,
		},
		{
			name:    "NG too long Tag",
			input:   &FindingTagForUpsert{FindingId: 1001, Tag: len65string},
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
			input:   &ResourceTagForUpsert{ResourceId: 1001, Tag: "tag"},
			wantErr: false,
		},
		{
			name:    "NG required FindingId",
			input:   &ResourceTagForUpsert{ResourceId: 0, Tag: "tag"},
			wantErr: true,
		},
		{
			name:    "NG required Tag",
			input:   &ResourceTagForUpsert{ResourceId: 1001},
			wantErr: true,
		},
		{
			name:    "NG too long Tag",
			input:   &ResourceTagForUpsert{ResourceId: 1001, Tag: len65string},
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

func TestValidate_PendFindingForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *PendFindingForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PendFindingForUpsert{FindingId: 1, ProjectId: 1, Note: "note"},
			wantErr: false,
		},
		{
			name:    "NG required FindingId",
			input:   &PendFindingForUpsert{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG required ProjectID",
			input:   &PendFindingForUpsert{FindingId: 1},
			wantErr: true,
		},
		{
			name:    "NG length Note",
			input:   &PendFindingForUpsert{FindingId: 1, ProjectId: 1, Note: len129string},
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

func TestValidate_FindingSettingForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *FindingSettingForUpsert
		wantErr bool
	}{
		{
			name:  "OK",
			input: &FindingSettingForUpsert{ProjectId: 1, ResourceName: "rn", Setting: "{}", Status: FindingSettingStatus_SETTING_ACTIVE},
		},
		{
			name:    "NG required ProjectId",
			input:   &FindingSettingForUpsert{ResourceName: "rn", Setting: "{}", Status: FindingSettingStatus_SETTING_ACTIVE},
			wantErr: true,
		},
		{
			name:    "NG required ProjectId",
			input:   &FindingSettingForUpsert{ProjectId: 1, Setting: "{}", Status: FindingSettingStatus_SETTING_ACTIVE},
			wantErr: true,
		},
		{
			name:    "NG required Setting",
			input:   &FindingSettingForUpsert{ProjectId: 1, ResourceName: "rn", Status: FindingSettingStatus_SETTING_ACTIVE},
			wantErr: true,
		},
		{
			name:    "NG is not JSON Setting",
			input:   &FindingSettingForUpsert{ProjectId: 1, ResourceName: "rn", Setting: "{", Status: FindingSettingStatus_SETTING_ACTIVE},
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
