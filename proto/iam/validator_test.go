package iam

import (
	"testing"
)

func TestValidate_IsAuthorizedRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *IsAuthorizedRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/hoge-bucket"},
			wantErr: false,
		},
		{
			name:    "NG Required(userID)",
			input:   &IsAuthorizedRequest{ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/hoge-bucket"},
			wantErr: true,
		},
		{
			name:    "NG Required(projectID)",
			input:   &IsAuthorizedRequest{UserId: 111, ActionName: "finding/PutFinding", ResourceName: "aws:guardduty/hoge-bucket"},
			wantErr: true,
		},
		{
			name:    "NG Required(ActionName)",
			input:   &IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "", ResourceName: "aws:guardduty/hoge-bucket"},
			wantErr: true,
		},
		{
			name:    "NG Invalid format(ActionName)",
			input:   &IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "findingPutFinding", ResourceName: "aws:guardduty/hoge-bucket"},
			wantErr: true,
		},
		{
			name:    "NG Invalid format(ActionName)",
			input:   &IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding/", ResourceName: "aws:guardduty/hoge-bucket"},
			wantErr: true,
		},
		{
			name:    "NG Required(ResourceName)",
			input:   &IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: ""},
			wantErr: true,
		},
		{
			name:    "NG Invalid format(ResourceName)",
			input:   &IsAuthorizedRequest{UserId: 111, ProjectId: 1001, ActionName: "finding/PutFinding", ResourceName: "/hoge-bucket"},
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

func TestValidate_GetUserRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetUserRequest
		wantErr bool
	}{
		{
			name:    "OK multi",
			input:   &GetUserRequest{UserId: 111, Sub: "1001"},
			wantErr: false,
		},
		{
			name:    "OK single(user_id)",
			input:   &GetUserRequest{UserId: 111},
			wantErr: false,
		},
		{
			name:    "OK single(sub)",
			input:   &GetUserRequest{Sub: "1001"},
			wantErr: false,
		},
		{
			name:    "NG requred",
			input:   &GetUserRequest{UserId: 0, Sub: ""},
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
