package iam

import (
	"testing"
)

func TestValidate_ListUserRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListUserRequest
		wantErr bool
	}{
		{
			name:    "OK multi",
			input:   &ListUserRequest{ProjectId: 111, Name: "nm"},
			wantErr: false,
		},
		{
			name:    "OK single(project_id)",
			input:   &ListUserRequest{ProjectId: 111},
			wantErr: false,
		},
		{
			name:    "NG length",
			input:   &ListUserRequest{ProjectId: 111, Name: "12345678901234567890123456789012345678901234567890123456789012345"},
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

func TestValidate_PutUserRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutUserRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutUserRequest{User: &UserForUpsert{Sub: "sub", Name: "nm", Activated: true}},
			wantErr: false,
		},
		{
			name:    "NG empty(user)",
			input:   &PutUserRequest{},
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

func TestValidate_UserForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *UserForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &UserForUpsert{Sub: "sub", Name: "nm", Activated: true},
			wantErr: false,
		},
		{
			name:    "NG required(sub)",
			input:   &UserForUpsert{Name: "nm", Activated: true},
			wantErr: true,
		},
		{
			name:    "NG required(name)",
			input:   &UserForUpsert{Sub: "sub", Activated: true},
			wantErr: true,
		},
		{
			name:    "NG length(sub)",
			input:   &UserForUpsert{Sub: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789=12345678901234567890123456789012345678901234567890123456", Name: "nm", Activated: true},
			wantErr: true,
		},
		{
			name:    "NG length(sub)",
			input:   &UserForUpsert{Sub: "sub", Name: "12345678901234567890123456789012345678901234567890123456789012345", Activated: true},
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

func TestValidate_ListRoleRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListRoleRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListRoleRequest{ProjectId: 123, Name: "nm"},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListRoleRequest{Name: "nm"},
			wantErr: true,
		},
		{
			name:    "NG Length(name)",
			input:   &ListRoleRequest{ProjectId: 123, Name: "12345678901234567890123456789012345678901234567890123456789012345"},
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

func TestValidate_GetRoleRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetRoleRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetRoleRequest{ProjectId: 123, RoleId: 123},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetRoleRequest{RoleId: 123},
			wantErr: true,
		},
		{
			name:    "NG Required(role_id)",
			input:   &GetRoleRequest{ProjectId: 123},
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

func TestValidate_PutRoleRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutRoleRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutRoleRequest{ProjectId: 123, Role: &RoleForUpsert{ProjectId: 123, Name: "nm"}},
			wantErr: false,
		},
		{
			name:    "NG Empty(role)",
			input:   &PutRoleRequest{ProjectId: 123},
			wantErr: true,
		},
		{
			name:    "NG NotMatch(project_id)",
			input:   &PutRoleRequest{ProjectId: 123, Role: &RoleForUpsert{ProjectId: 999, Name: "nm"}},
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

func TestValidate_RoleForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *RoleForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &RoleForUpsert{ProjectId: 123, Name: "nm"},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &RoleForUpsert{Name: "nm"},
			wantErr: true,
		},
		{
			name:    "NG Required(name)",
			input:   &RoleForUpsert{ProjectId: 123},
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

func TestValidate_DeleteRoleRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteRoleRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeleteRoleRequest{ProjectId: 123, RoleId: 123},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteRoleRequest{RoleId: 123},
			wantErr: true,
		},
		{
			name:    "NG Required(role_id)",
			input:   &DeleteRoleRequest{ProjectId: 123},
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

func TestValidate_AttachRoleRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *AttachRoleRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &AttachRoleRequest{ProjectId: 123, UserId: 1, RoleId: 1},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &AttachRoleRequest{UserId: 1, RoleId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(user_id)",
			input:   &AttachRoleRequest{ProjectId: 123, RoleId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(role_id)",
			input:   &AttachRoleRequest{ProjectId: 123, UserId: 1},
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

func TestValidate_DetachRoleRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DetachRoleRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DetachRoleRequest{ProjectId: 123, UserId: 1, RoleId: 1},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DetachRoleRequest{UserId: 1, RoleId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(user_id)",
			input:   &DetachRoleRequest{ProjectId: 123, RoleId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(role_id)",
			input:   &DetachRoleRequest{ProjectId: 123, UserId: 1},
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

func TestValidate_ListPolicyRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListPolicyRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &ListPolicyRequest{ProjectId: 123, Name: "nm"},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &ListPolicyRequest{Name: "nm"},
			wantErr: true,
		},
		{
			name:    "NG length(name)",
			input:   &ListPolicyRequest{ProjectId: 123, Name: "12345678901234567890123456789012345678901234567890123456789012345"},
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

func TestValidate_GetPolicyRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *GetPolicyRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &GetPolicyRequest{ProjectId: 123, PolicyId: 1},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &GetPolicyRequest{PolicyId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(policy_id)",
			input:   &GetPolicyRequest{ProjectId: 123},
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

func TestValidate_PutPolicyRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *PutPolicyRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PutPolicyRequest{ProjectId: 123, Policy: &PolicyForUpsert{ProjectId: 123, Name: "nm", ActionPtn: ".*", ResourcePtn: ".*"}},
			wantErr: false,
		},
		{
			name:    "NG Empty(policy)",
			input:   &PutPolicyRequest{ProjectId: 123},
			wantErr: true,
		},
		{
			name:    "NG NotMatch(project_id)",
			input:   &PutPolicyRequest{ProjectId: 123, Policy: &PolicyForUpsert{ProjectId: 999, Name: "nm", ActionPtn: ".*", ResourcePtn: ".*"}},
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

func TestValidate_PolicyForUpsert(t *testing.T) {
	cases := []struct {
		name    string
		input   *PolicyForUpsert
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &PolicyForUpsert{ProjectId: 123, Name: "nm", ActionPtn: ".*", ResourcePtn: ".*"},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &PolicyForUpsert{Name: "nm", ActionPtn: ".*", ResourcePtn: ".*"},
			wantErr: true,
		},
		{
			name:    "NG Required(name)",
			input:   &PolicyForUpsert{ProjectId: 123, ActionPtn: ".*", ResourcePtn: ".*"},
			wantErr: true,
		},
		{
			name:    "NG Length(name)",
			input:   &PolicyForUpsert{ProjectId: 123, Name: "12345678901234567890123456789012345678901234567890123456789012345", ActionPtn: ".*", ResourcePtn: ".*"},
			wantErr: true,
		},
		{
			name:    "NG Required(action_ptn)",
			input:   &PolicyForUpsert{ProjectId: 123, Name: "nm", ResourcePtn: ".*"},
			wantErr: true,
		},
		{
			name:    "NG Not regexp pattern(action_ptn)",
			input:   &PolicyForUpsert{ProjectId: 123, Name: "nm", ActionPtn: "*", ResourcePtn: ".*"},
			wantErr: true,
		},
		{
			name:    "NG Required(resource_ptn)",
			input:   &PolicyForUpsert{ProjectId: 123, Name: "nm", ActionPtn: ".*"},
			wantErr: true,
		},
		{
			name:    "NG Not regexp pattern(resource_ptn)",
			input:   &PolicyForUpsert{ProjectId: 123, Name: "nm", ActionPtn: ".*", ResourcePtn: "*xxx"},
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

func TestValidate_DeletePolicyRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeletePolicyRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DeletePolicyRequest{ProjectId: 123, PolicyId: 1},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeletePolicyRequest{PolicyId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(policy_id)",
			input:   &DeletePolicyRequest{ProjectId: 123},
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

func TestValidate_AttachPolicyRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *AttachPolicyRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &AttachPolicyRequest{ProjectId: 123, RoleId: 1, PolicyId: 1},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &AttachPolicyRequest{RoleId: 1, PolicyId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(role_id)",
			input:   &AttachPolicyRequest{ProjectId: 123, PolicyId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(policy_id)",
			input:   &AttachPolicyRequest{ProjectId: 123, RoleId: 1},
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

func TestValidate_DetachPolicyRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DetachPolicyRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &DetachPolicyRequest{ProjectId: 123, RoleId: 1, PolicyId: 1},
			wantErr: false,
		},
		{
			name:    "NG Required(project_id)",
			input:   &DetachPolicyRequest{RoleId: 1, PolicyId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(role_id)",
			input:   &DetachPolicyRequest{ProjectId: 123, PolicyId: 1},
			wantErr: true,
		},
		{
			name:    "NG Required(policy_id)",
			input:   &DetachPolicyRequest{ProjectId: 123, RoleId: 1},
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

func TestValidate_IsAdminRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *IsAdminRequest
		wantErr bool
	}{
		{
			name:    "OK",
			input:   &IsAdminRequest{UserId: 1},
			wantErr: false,
		},
		{
			name:    "NG Required(userID)",
			input:   &IsAdminRequest{UserId: 0},
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
