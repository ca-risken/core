package main

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type param struct {
	Param1 uint32    `json:"param1"`
	Param2 string    `json:"param2"`
	Param3 float32   `json:"param3"`
	Param4 []uint32  `json:"param4"`
	Param5 []string  `json:"param5"`
	Param6 []float32 `json:"param6"`
}

func TestBindQuery(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    *param
		wantErr bool
	}{
		{
			name:  "OK",
			input: `param1=123&param2=aaa&param3=1.1&param4=1,2,3&param5=aaa,bbb,ccc&param6=0.1,0.2,0.3`,
			want:  &param{Param1: 123, Param2: "aaa", Param3: 1.1, Param4: []uint32{1, 2, 3}, Param5: []string{"aaa", "bbb", "ccc"}, Param6: []float32{0.1, 0.2, 0.3}},
		},
		{
			name:  "OK Zero value",
			input: `param4=&param5=&param6=`,
			want:  &param{Param1: 0, Param2: "", Param3: 0, Param4: []uint32{0}, Param5: []string{}, Param6: []float32{0}},
		},
		{
			name:    "NG Parse error",
			input:   `param1=string`,
			want:    &param{},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/test?"+c.input, nil)
			got := param{}
			err := bindQuery(&got, req)
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
	cases := []struct {
		name    string
		input   string
		want    *param
		wantErr bool
	}{
		{
			name:  "OK",
			input: `{"param1":123, "param2":"aaa", "param3":11.1}`,
			want:  &param{Param1: 123, Param2: "aaa", Param3: 11.1},
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
			err := bindBodyJSON(&got, req)
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

func TestStringSeparator(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		delimiter rune
		want      []string
	}{
		{
			name:  "single param",
			input: "aaa",
			want:  []string{"aaa"},
		},
		{
			name:      "multiple params",
			input:     "aaa,bbb,ccc",
			delimiter: ',',
			want:      []string{"aaa", "bbb", "ccc"},
		},
		{
			name:      "exsits blank",
			input:     "aaa,,ccc",
			delimiter: ',',
			want:      []string{"aaa", "ccc"},
		},
		{
			name:      "other delimiter",
			input:     "aaa#ccc",
			delimiter: '#',
			want:      []string{"aaa", "ccc"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := stringSeparator(c.input, c.delimiter)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Unexpected result: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
