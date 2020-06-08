package finding

import (
	"testing"
	"time"
)

func TestValidate_ListFindingRequest(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name    string
		input   ListFindingRequest
		wantErr bool
	}{
		{
			name:    "OK empty",
			input:   ListFindingRequest{},
			wantErr: false,
		},
		{
			name:    "OK full parameters",
			input:   ListFindingRequest{UserId: 1001, ProjectId: []uint32{111, 222}, DataSource: []string{"ds1", "ds2"}, ResourceName: []string{"rn1", "rn2"}, FromScore: 0.0, ToScore: 1.0, FromAt: now.Unix(), ToAt: now.Unix()},
			wantErr: false,
		},
		{
			name:    "NG too long resource_name",
			input:   ListFindingRequest{ResourceName: []string{"123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=12345678901234567890123456789012345678901234567890123456"}},
			wantErr: true,
		},
		{
			name:    "NG too long data_source",
			input:   ListFindingRequest{DataSource: []string{"12345678901234567890123456789012345678901234567890123456789012345"}},
			wantErr: true,
		},
		{
			name:    "NG small from_score",
			input:   ListFindingRequest{FromScore: -0.1},
			wantErr: true,
		},
		{
			name:    "NG big from_score",
			input:   ListFindingRequest{FromScore: 1.1},
			wantErr: true,
		},
		{
			name:    "NG small to_score",
			input:   ListFindingRequest{ToScore: -0.1},
			wantErr: true,
		},
		{
			name:    "NG big to_score",
			input:   ListFindingRequest{ToScore: 1.1},
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
		input   GetFindingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   GetFindingRequest{FindingId: 1001},
			wantErr: false,
		},
		{
			name:    "NG required",
			input:   GetFindingRequest{},
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
		input   DeleteFindingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   DeleteFindingRequest{FindingId: 1001},
			wantErr: false,
		},
		{
			name:    "NG required",
			input:   DeleteFindingRequest{},
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
		input   ListFindingTagRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   ListFindingTagRequest{FindingId: 1001},
			wantErr: false,
		},
		{
			name:    "NG required",
			input:   ListFindingTagRequest{},
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
		input   UntagFindingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   UntagFindingRequest{FindingTagId: 1001},
			wantErr: false,
		},
		{
			name:    "NG required",
			input:   UntagFindingRequest{},
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
		input   ListResourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   ListResourceRequest{},
			wantErr: false,
		},
		{
			name:    "NG too long resource_name",
			input:   ListResourceRequest{ResourceName: []string{"123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=12345678901234567890123456789012345678901234567890123456"}},
			wantErr: true,
		},
		{
			name:    "NG too small from_sum_score",
			input:   ListResourceRequest{FromSumScore: -0.1},
			wantErr: true,
		},
		{
			name:    "NG too small to_sum_score",
			input:   ListResourceRequest{ToSumScore: -0.1},
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
		input   GetResourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   GetResourceRequest{ResourceId: 1001},
			wantErr: false,
		},
		{
			name:    "NG required",
			input:   GetResourceRequest{},
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
		input   DeleteResourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   DeleteResourceRequest{ResourceId: 1001},
			wantErr: false,
		},
		{
			name:    "NG required",
			input:   DeleteResourceRequest{},
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
		input   ListResourceTagRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   ListResourceTagRequest{ResourceId: 1001},
			wantErr: false,
		},
		{
			name:    "NG required",
			input:   ListResourceTagRequest{},
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
		input   UntagResourceRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   UntagResourceRequest{ResourceTagId: 1001},
			wantErr: false,
		},
		{
			name:    "NG required",
			input:   UntagResourceRequest{},
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
		input   FindingForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   FindingForUpsert{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: false,
		},
		{
			name:    "NG too long Description",
			input:   FindingForUpsert{FindingId: 1001, Description: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=1", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG required DataSource",
			input:   FindingForUpsert{FindingId: 1001, Description: "desc", DataSource: "", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG too long DataSource",
			input:   FindingForUpsert{FindingId: 1001, Description: "desc", DataSource: "12345678901234567890123456789012345678901234567890123456789012345", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG required DataSourceId",
			input:   FindingForUpsert{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG too long DataSourceId",
			input:   FindingForUpsert{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=01234567890123456789012345678901234567890123456789123456", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG required resource name",
			input:   FindingForUpsert{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG too long resource name",
			input:   FindingForUpsert{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=12345678901234567890123456789012345678901234567890123456", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG nil OriginalScore",
			input:   FindingForUpsert{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG too small OriginalScore",
			input:   FindingForUpsert{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: -0.1, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG OriginalScore bigger than OriginalMaxScore",
			input:   FindingForUpsert{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 100.01, OriginalMaxScore: 100.0, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG nil OriginalMaxScore",
			input:   FindingForUpsert{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG too small OriginalMaxScore",
			input:   FindingForUpsert{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: -0.01, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG too big OriginalMaxScore",
			input:   FindingForUpsert{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 999.991, Data: `{"key": "value"}`},
			wantErr: true,
		},
		{
			name:    "NG invalid json Data",
			input:   FindingForUpsert{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{"key": "value"`},
			wantErr: true,
		},
		{
			name:    "NG invalid json Data2",
			input:   FindingForUpsert{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1001, OriginalScore: 50.5, OriginalMaxScore: 100.0, Data: `{key: value}`},
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
		input   FindingTagForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   FindingTagForUpsert{FindingTagId: 1001001, FindingId: 1001, TagKey: "key", TagValue: "value"},
			wantErr: false,
		},
		{
			name:    "NG required FindingId",
			input:   FindingTagForUpsert{FindingTagId: 1001001, FindingId: 0, TagKey: "key", TagValue: "value"},
			wantErr: true,
		},
		{
			name:    "NG required TagKey",
			input:   FindingTagForUpsert{FindingTagId: 1001001, FindingId: 1001, TagKey: "", TagValue: "value"},
			wantErr: true,
		},
		{
			name:    "NG too long TagKey",
			input:   FindingTagForUpsert{FindingTagId: 1001001, FindingId: 1001, TagKey: "12345678901234567890123456789012345678901234567890123456789012345", TagValue: "value"},
			wantErr: true,
		},
		{
			name:    "NG too long TagValue",
			input:   FindingTagForUpsert{FindingTagId: 1001001, FindingId: 1001, TagKey: "key", TagValue: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=1"},
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
		input   ResourceForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   ResourceForUpsert{ResourceId: 1001, ResourceName: "rn", ProjectId: 1001},
			wantErr: false,
		},
		{
			name:    "NG required ResourceName",
			input:   ResourceForUpsert{ResourceId: 1001, ResourceName: "", ProjectId: 1001},
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
		input   ResourceTagForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   ResourceTagForUpsert{ResourceTagId: 1001001, ResourceId: 1001, TagKey: "key", TagValue: "value"},
			wantErr: false,
		},
		{
			name:    "NG required FindingId",
			input:   ResourceTagForUpsert{ResourceTagId: 1001001, ResourceId: 0, TagKey: "key", TagValue: "value"},
			wantErr: true,
		},
		{
			name:    "NG required TagKey",
			input:   ResourceTagForUpsert{ResourceTagId: 1001001, ResourceId: 1001, TagKey: "", TagValue: "value"},
			wantErr: true,
		},
		{
			name:    "NG too long TagKey",
			input:   ResourceTagForUpsert{ResourceTagId: 1001001, ResourceId: 1001, TagKey: "12345678901234567890123456789012345678901234567890123456789012345", TagValue: "value"},
			wantErr: true,
		},
		{
			name:    "NG too long TagValue",
			input:   ResourceTagForUpsert{ResourceTagId: 1001001, ResourceId: 1001, TagKey: "key", TagValue: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=1"},
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
