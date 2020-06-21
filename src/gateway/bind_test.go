package main

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestBindQuery(t *testing.T) {
	type param struct {
		ProjectID uint32 `json:"project_id"`
	}
	cases := []struct {
		name    string
		input   string
		want    *param
		wantErr bool
	}{
		{
			name:  "OK",
			input: `project_id=123`,
			want:  &param{ProjectID: 123},
		},
		{
			name:    "NG parse error",
			input:   `project_id=string`,
			want:    &param{},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/test?"+c.input, nil)
			got := param{}
			err := BindQuery(&got, req)
			if err == nil && c.wantErr {
				t.Fatalf("Unexpected no error: wantErr=%t", c.wantErr)
			}
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: wantErr=%t, err=%+v", c.wantErr, err)
			}
			if !reflect.DeepEqual(c.want, &got) {
				t.Fatalf("Unexpected bind: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestBindBodyJSON(t *testing.T) {
	type param struct {
		ProjectID uint32 `json:"project_id"`
	}
	cases := []struct {
		name    string
		input   string
		want    *param
		wantErr bool
	}{
		{
			name:  "OK",
			input: `{"project_id":123}`,
			want:  &param{ProjectID: 123},
		},
		{
			name:    "NG parse error",
			input:   `{"project_id":`,
			want:    &param{},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/test", strings.NewReader(c.input))
			got := param{}
			err := BindBodyJSON(&got, req)
			if err == nil && c.wantErr {
				t.Fatalf("Unexpected no error: wantErr=%t", c.wantErr)
			}
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: wantErr=%t, err=%+v", c.wantErr, err)
			}
			if !reflect.DeepEqual(c.want, &got) {
				t.Fatalf("Unexpected bind: want=%+v, got=%+v", c.want, got)
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
		{
			name:  "blank params",
			input: "aaa,,ccc",
			want:  []string{"aaa", "ccc"},
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
