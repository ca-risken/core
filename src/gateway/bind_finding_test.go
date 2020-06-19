package main

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/go-chi/chi"
)

func TestBindListFindingRequest(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  *finding.ListFindingRequest
	}{
		{
			name:  "no param",
			input: "",
			want:  &finding.ListFindingRequest{},
		},
		{
			name:  "full set",
			input: "project_id=123&data_source=aws:guardduty&resource_name=resouce&from_score=0.5&to_score=1.0&from_at=1590000000&to_at=1600000000",
			want:  &finding.ListFindingRequest{ProjectId: []uint32{123}, DataSource: []string{"aws:guardduty"}, ResourceName: []string{"resouce"}, FromScore: 0.5, ToScore: 1.0, FromAt: 1590000000, ToAt: 1600000000},
		},
		{
			name:  "multiple project_id 1",
			input: "project_id=123,456",
			want:  &finding.ListFindingRequest{ProjectId: []uint32{123, 456}},
		},
		{
			name:  "multiple project_id 2",
			input: "project_id=123,",
			want:  &finding.ListFindingRequest{ProjectId: []uint32{123}},
		},
		{
			name:  "multiple data_source 1",
			input: "data_source=aws:guardduty,github:code-scaning",
			want:  &finding.ListFindingRequest{DataSource: []string{"aws:guardduty", "github:code-scaning"}},
		},
		{
			name:  "multiple resource_name 1",
			input: "resource_name=resource-1,resource-2",
			want:  &finding.ListFindingRequest{ResourceName: []string{"resource-1", "resource-2"}},
		},
		{
			name:  "score parse error",
			input: "from_score=parse_error&to_score=parse_error",
			want:  &finding.ListFindingRequest{},
		},
		{
			name:  "time parse error",
			input: "from_at=parse_error&to_at=parse_error",
			want:  &finding.ListFindingRequest{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/finding?"+c.input, nil)
			got := bindListFindingRequest(req)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected bind: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBindGetFindingRequest(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  *finding.GetFindingRequest
	}{
		{
			name:  "No param",
			input: "",
			want:  &finding.GetFindingRequest{},
		},
		{
			name:  "finding_id",
			input: "1001",
			want:  &finding.GetFindingRequest{FindingId: 1001},
		},
		{
			name:  "parse error",
			input: "xxxxxxxx",
			want:  &finding.GetFindingRequest{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/finding/"+c.input, nil)
			// Requestにパスパラメータ{finding_id}を登録
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("finding_id", c.input)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			got := bindGetFindingRequest(req)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected bind: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBindPutFindingRequest(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  *finding.PutFindingRequest
	}{
		{
			name:  "OK",
			input: `{"description":"desc", "data_source":"ds", "data_source_id":"ds-01", "resource_name":"rn", "project_id":1, "original_score":0.1, "original_max_score":1.0, "data":"{}"}`,
			want:  &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{Description: "desc", DataSource: "ds", DataSourceId: "ds-01", ResourceName: "rn", ProjectId: 1, OriginalScore: 0.1, OriginalMaxScore: 1.0, Data: `{}`}},
		},
		{
			name:  "parse error",
			input: `{"description":"desc`,
			want:  &finding.PutFindingRequest{Finding: &finding.FindingForUpsert{}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/finding/put", strings.NewReader(c.input))
			got := bindPutFindingRequest(req)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected bind: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBindDeleteFindingRequest(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  *finding.DeleteFindingRequest
	}{
		{
			name:  "OK",
			input: `{"finding_id":1001}`,
			want:  &finding.DeleteFindingRequest{FindingId: 1001},
		},
		{
			name:  "parse error",
			input: "xxxxxxxx",
			want:  &finding.DeleteFindingRequest{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/finding/delete", strings.NewReader(c.input))
			got := bindDeleteFindingRequest(req)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected bind: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBindListFindingTagRequest(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  *finding.ListFindingTagRequest
	}{
		{
			name:  "OK",
			input: `1001`,
			want:  &finding.ListFindingTagRequest{FindingId: 1001},
		},
		{
			name:  "parse error",
			input: "xxx",
			want:  &finding.ListFindingTagRequest{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/finding/"+c.input+"/tag", strings.NewReader(c.input))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("finding_id", c.input) // Requestにパスパラメータ{finding_id}を登録
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			got := bindListFindingTagRequest(req)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected bind: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBindTagFindingRequest(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  *finding.TagFindingRequest
	}{
		{
			name:  "OK",
			input: `{"finding_id":111, "tag_key":"key", "tag_value":"value"}`,
			want:  &finding.TagFindingRequest{Tag: &finding.FindingTagForUpsert{FindingId: 111, TagKey: "key", TagValue: "value"}},
		},
		{
			name:  "parse error",
			input: "xxxxxxxx",
			want:  &finding.TagFindingRequest{Tag: &finding.FindingTagForUpsert{}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/finding/tag", strings.NewReader(c.input))
			got := bindTagFindingRequest(req)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected bind: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBindUntagFindingRequest(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  *finding.UntagFindingRequest
	}{
		{
			name:  "OK",
			input: `{"finding_tag_id":1001}`,
			want:  &finding.UntagFindingRequest{FindingTagId: 1001},
		},
		{
			name:  "parse error",
			input: "xxxxxxxx",
			want:  &finding.UntagFindingRequest{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/finding/untag", strings.NewReader(c.input))
			got := bindUntagFindingRequest(req)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected bind: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBindListResourceRequest(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  *finding.ListResourceRequest
	}{
		{
			name:  "OK",
			input: "project_id=111,222&resource_name=aaa,bbb&from_sum_score=0.0&to_sum_score=100.0&from_at=&to_at=",
			want:  &finding.ListResourceRequest{ProjectId: []uint32{111, 222}, ResourceName: []string{"aaa", "bbb"}, FromSumScore: 0.0, ToSumScore: 100.0},
		},
		{
			name:  "OK No param",
			input: "",
			want:  &finding.ListResourceRequest{},
		},
		{
			name:  "score parse error",
			input: "from_sum_score=parse_error&to_sum_score=parse_error",
			want:  &finding.ListResourceRequest{},
		},
		{
			name:  "time parse error",
			input: "from_at=parse_error&to_at=parse_error",
			want:  &finding.ListResourceRequest{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/resource?"+c.input, nil)
			got := bindListResourceRequest(req)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected bind: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBindGetResourceRequest(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  *finding.GetResourceRequest
	}{
		{
			name:  "No param",
			input: "",
			want:  &finding.GetResourceRequest{},
		},
		{
			name:  "resource id",
			input: "1001",
			want:  &finding.GetResourceRequest{ResourceId: 1001},
		},
		{
			name:  "parse error",
			input: "xxxxxxxx",
			want:  &finding.GetResourceRequest{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/resource/"+c.input, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("resource_id", c.input) // Requestにパスパラメータ{resource_id}を登録
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			got := bindGetResourceRequest(req)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected bind: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBindPutResourceRequest(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  *finding.PutResourceRequest
	}{
		{
			name:  "OK",
			input: `{"resource_name":"rn", "project_id":1}`,
			want:  &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{ResourceName: "rn", ProjectId: 1}},
		},
		{
			name:  "parse error",
			input: "xxxxxxxx",
			want:  &finding.PutResourceRequest{Resource: &finding.ResourceForUpsert{}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/resource/put", strings.NewReader(c.input))
			got := bindPutResourceRequest(req)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected bind: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBindDeleteResourceRequest(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  *finding.DeleteResourceRequest
	}{
		{
			name:  "OK",
			input: `{"resource_id":1001}`,
			want:  &finding.DeleteResourceRequest{ResourceId: 1001},
		},
		{
			name:  "parse error",
			input: "xxxxxxxx",
			want:  &finding.DeleteResourceRequest{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/resource/delete", strings.NewReader(c.input))
			got := bindDeleteResourceRequest(req)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected bind: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBindListResourceTagReqest(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  *finding.ListResourceTagRequest
	}{
		{
			name:  "OK",
			input: `1001`,
			want:  &finding.ListResourceTagRequest{ResourceId: 1001},
		},
		{
			name:  "parse error",
			input: "xxx",
			want:  &finding.ListResourceTagRequest{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/resource/"+c.input+"/tag", strings.NewReader(c.input))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("resource_id", c.input) // Requestにパスパラメータ{resource_id}を登録
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			got := bindListResourceTagRequest(req)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected bind: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBindTagResourceRequest(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  *finding.TagResourceRequest
	}{
		{
			name:  "OK",
			input: `{"resource_id":111, "tag_key":"key", "tag_value":"value"}`,
			want:  &finding.TagResourceRequest{Tag: &finding.ResourceTagForUpsert{ResourceId: 111, TagKey: "key", TagValue: "value"}},
		},
		{
			name:  "parse error",
			input: "xxxxxxxx",
			want:  &finding.TagResourceRequest{Tag: &finding.ResourceTagForUpsert{}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/resource/tag", strings.NewReader(c.input))
			got := bindTagResourceRequest(req)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected bind: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBindUntagResourceRequest(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  *finding.UntagResourceRequest
	}{
		{
			name:  "OK",
			input: `{"resource_tag_id":111}`,
			want:  &finding.UntagResourceRequest{ResourceTagId: 111},
		},
		{
			name:  "parse error",
			input: "xxxxxxxx",
			want:  &finding.UntagResourceRequest{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/resource/untag", strings.NewReader(c.input))
			got := bindUntagResourceRequest(req)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected bind: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
