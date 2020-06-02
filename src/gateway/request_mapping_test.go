package main

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/CyberAgent/mimosa-core/proto/finding"
)

func TestMappingListFindingRequest(t *testing.T) {
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
			input: "project_id=aaa&data_source=aws:guardduty&resource_name=resouce&from_score=0.5&to_score=1.0&from_at=1590000000&to_at=1600000000",
			want:  &finding.ListFindingRequest{ProjectId: []string{"aaa"}, DataSource: []string{"aws:guardduty"}, ResourceName: []string{"resouce"}, FromScore: 0.5, ToScore: 1.0, FromAt: 1590000000, ToAt: 1600000000},
		},
		{
			name:  "multiple project_id 1",
			input: "project_id=hoge,fuga",
			want:  &finding.ListFindingRequest{ProjectId: []string{"hoge", "fuga"}},
		},
		{
			name:  "multiple project_id 2",
			input: "project_id=hoge,",
			want:  &finding.ListFindingRequest{ProjectId: []string{"hoge", ""}},
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
			r, _ := http.NewRequest("GET", "/?"+c.input, nil)
			mapping := mappingListFindingRequest(r)
			if !reflect.DeepEqual(mapping, c.want) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, mapping)
			}
		})
	}
}

func TestCommaSeparator(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "single param",
			input: "aaa",
			want:  []string{"aaa"},
		},
		{
			name:  "multiple params",
			input: "aaa,bbb,ccc",
			want:  []string{"aaa", "bbb", "ccc"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := commaSeparator(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected result: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestParseScore(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    float32
		wantErr bool
	}{
		{
			name:    "normal",
			input:   "0.05",
			want:    0.05,
			wantErr: false,
		},
		{
			name:    "parse error",
			input:   "parse error",
			want:    0.00,
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parseScore(c.input)
			if c.wantErr && err == nil {
				t.Fatalf("No error occurred: wantErr=%+v, err=%+v", c.wantErr, err)
			}
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occurred: wantErr=%+v, err=%+v", c.wantErr, err)
			}
			if got != c.want {
				t.Fatalf("Unexpected result: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestParseTimeParam(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{
			name:    "normal",
			input:   "1591034681",
			want:    1591034681,
			wantErr: false,
		},
		{
			name:    "parse error",
			input:   "parse error",
			want:    0,
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parseTimeParam(c.input)
			if c.wantErr && err == nil {
				t.Fatalf("No error occurred: wantErr=%+v, err=%+v", c.wantErr, err)
			}
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occurred: wantErr=%+v, err=%+v", c.wantErr, err)
			}
			if got != c.want {
				t.Fatalf("Unexpected result: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
