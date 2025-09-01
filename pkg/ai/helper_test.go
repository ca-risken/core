package ai

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExtractSingleJSONObject(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "Valid JSON object",
			input:   `{"key": "value", "number": 42}`,
			want:    `{"key": "value", "number": 42}`,
			wantErr: false,
		},
		{
			name:    "JSON object with text before",
			input:   `Here is the JSON: {"key": "value", "number": 42}`,
			want:    `{"key": "value", "number": 42}`,
			wantErr: false,
		},
		{
			name:    "JSON object with text after",
			input:   `{"key": "value", "number": 42} and some text after`,
			want:    `{"key": "value", "number": 42}`,
			wantErr: false,
		},
		{
			name:    "JSON object with text before and after",
			input:   `Some text before {"key": "value", "number": 42} and after`,
			want:    `{"key": "value", "number": 42}`,
			wantErr: false,
		},
		{
			name:    "Nested JSON object",
			input:   `{"outer": {"inner": {"key": "value"}}, "array": [1, 2, 3]}`,
			want:    `{"outer": {"inner": {"key": "value"}}, "array": [1, 2, 3]}`,
			wantErr: false,
		},
		{
			name:    "JSON with escaped braces in string",
			input:   `{"message": "This has \\{escaped\\} braces", "valid": true}`,
			want:    `{"message": "This has \\{escaped\\} braces", "valid": true}`,
			wantErr: false,
		},
		{
			name:    "Multiple JSON objects - returns outermost",
			input:   `{"first": {"nested": "value"}} {"second": "object"}`,
			want:    ``,
			wantErr: true,
		},
		{
			name:    "Empty JSON object",
			input:   `{}`,
			want:    `{}`,
			wantErr: false,
		},
		{
			name:    "JSON array",
			input:   `[{"key": "value"}, {"key": "value2"}]`,
			want:    ``,
			wantErr: true,
		},
		{
			name:    "No opening brace",
			input:   `"key": "value"}`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "No closing brace",
			input:   `{"key": "value"`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "No braces at all",
			input:   `just plain text`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Empty string",
			input:   ``,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Only opening brace",
			input:   `{`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Only closing brace",
			input:   `}`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Closing brace before opening brace",
			input:   `} some text {`,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractSingleJSONObject(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractSingleJSONObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractSingleJSONObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertSchema(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	type NestedStruct struct {
		ID   int `json:"id"`
		Data struct {
			Value   string `json:"value"`
			Numbers []int  `json:"numbers"`
			Active  bool   `json:"active"`
		} `json:"data"`
	}

	tests := []struct {
		name    string
		input   string
		schema  interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name:   "Valid JSON with simple struct",
			input:  `{"name": "John Doe", "age": 30, "email": "john@example.com"}`,
			schema: TestStruct{},
			want: &TestStruct{
				Name:  "John Doe",
				Age:   30,
				Email: "john@example.com",
			},
			wantErr: false,
		},
		{
			name:   "Valid JSON with text before and after",
			input:  `Here is the data: {"name": "Jane Doe", "age": 25, "email": "jane@example.com"} end of data`,
			schema: TestStruct{},
			want: &TestStruct{
				Name:  "Jane Doe",
				Age:   25,
				Email: "jane@example.com",
			},
			wantErr: false,
		},
		{
			name:   "Valid JSON with whitespace",
			input:  "   \n\t  {\"name\": \"Bob\", \"age\": 40, \"email\": \"bob@example.com\"}  \n  ",
			schema: TestStruct{},
			want: &TestStruct{
				Name:  "Bob",
				Age:   40,
				Email: "bob@example.com",
			},
			wantErr: false,
		},
		{
			name: "Valid nested JSON",
			input: `{
				"id": 123,
				"data": {
					"value": "test",
					"numbers": [1, 2, 3],
					"active": true
				}
			}`,
			schema: NestedStruct{},
			want: &NestedStruct{
				ID: 123,
				Data: struct {
					Value   string `json:"value"`
					Numbers []int  `json:"numbers"`
					Active  bool   `json:"active"`
				}{
					Value:   "test",
					Numbers: []int{1, 2, 3},
					Active:  true,
				},
			},
			wantErr: false,
		},
		{
			name:   "Partial JSON matching schema",
			input:  `{"name": "Alice", "age": 35}`,
			schema: TestStruct{},
			want: &TestStruct{
				Name:  "Alice",
				Age:   35,
				Email: "", // missing field gets zero value
			},
			wantErr: false,
		},
		{
			name:    "Invalid JSON syntax",
			input:   `{"name": "John", "age": 30, "email": "john@example.com"`,
			schema:  TestStruct{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "No JSON braces found",
			input:   `just plain text without JSON`,
			schema:  TestStruct{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty string",
			input:   ``,
			schema:  TestStruct{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid JSON - wrong data type",
			input:   `{"name": "John", "age": "thirty", "email": "john@example.com"}`,
			schema:  TestStruct{},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "Empty JSON object",
			input:  `{}`,
			schema: TestStruct{},
			want: &TestStruct{
				Name:  "",
				Age:   0,
				Email: "",
			},
			wantErr: false,
		},
		{
			name:    "Malformed JSON with extra comma",
			input:   `{"name": "John", "age": 30,}`,
			schema:  TestStruct{},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got any
			var err error

			switch tt.schema.(type) {
			case TestStruct:
				got, err = ConvertSchema(tt.input, TestStruct{})
			case NestedStruct:
				got, err = ConvertSchema(tt.input, NestedStruct{})
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("ConvertSchema() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
