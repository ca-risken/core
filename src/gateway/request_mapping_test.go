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
			name:  "full set",
			input: "project_id=hoge&since=20201231",
			want:  &finding.ListFindingRequest{ProjectId: "hoge", Since: "20201231"},
		},
		{
			name:  "no set values",
			input: "project_id=&since=",
			want:  &finding.ListFindingRequest{ProjectId: "", Since: ""},
		},
		{
			name:  "no query param",
			input: "",
			want:  &finding.ListFindingRequest{ProjectId: "", Since: ""},
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
