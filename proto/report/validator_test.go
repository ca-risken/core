package report

import (
	"testing"
)

func TestValidateGetReportRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetReportFindingRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetReportFindingRequest{ProjectId: 111, FromDate: "2000-01-01", ToDate: "9999-12-31", Score: 0.5},
			wantErr: false,
		},
		{
			name:    "NG (required ProjectId)",
			input:   &GetReportFindingRequest{},
			wantErr: true,
		},
		{
			name:    "NG (FromDate format)",
			input:   &GetReportFindingRequest{ProjectId: 111, FromDate: "hogehoge", ToDate: "9999-12-31", Score: 0.5},
			wantErr: true,
		},
		{
			name:    "NG (ToDate format)",
			input:   &GetReportFindingRequest{ProjectId: 111, FromDate: "2000-01-01", ToDate: "hogehoge", Score: 0.5},
			wantErr: true,
		},
		{
			name:    "NG (Score < 0.0)",
			input:   &GetReportFindingRequest{ProjectId: 111, FromDate: "2000-01-01", ToDate: "9999-12-31", Score: -0.1},
			wantErr: true,
		},
		{
			name:    "NG (Score > 1.0)",
			input:   &GetReportFindingRequest{ProjectId: 111, FromDate: "2000-01-01", ToDate: "9999-12-31", Score: 1.1},
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

func TestValidateGetReportAllRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetReportFindingAllRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetReportFindingAllRequest{FromDate: "2000-01-01", ToDate: "9999-12-31", Score: 0.5},
			wantErr: false,
		},
		{
			name:    "NG (FromDate format)",
			input:   &GetReportFindingAllRequest{FromDate: "hogehoge", ToDate: "9999-12-31", Score: 0.5},
			wantErr: true,
		},
		{
			name:    "NG (ToDate format)",
			input:   &GetReportFindingAllRequest{FromDate: "2000-01-01", ToDate: "hogehoge", Score: 0.5},
			wantErr: true,
		},
		{
			name:    "NG (Score < 0.0)",
			input:   &GetReportFindingAllRequest{FromDate: "2000-01-01", ToDate: "9999-12-31", Score: -0.1},
			wantErr: true,
		},
		{
			name:    "NG (Score > 1.0)",
			input:   &GetReportFindingAllRequest{FromDate: "2000-01-01", ToDate: "9999-12-31", Score: 1.1},
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
