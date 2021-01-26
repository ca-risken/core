package report

import (
	"testing"
)

func TestValidateGetReportRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetReportRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetReportRequest{ProjectId: 111},
			wantErr: false,
		},
		{
			name:    "NG (required ProjectId)",
			input:   &GetReportRequest{},
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
