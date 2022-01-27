package project

import (
	"testing"
)

func TestValidate_ListProjectRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *ListProjectRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &ListProjectRequest{ProjectId: 123, UserId: 1, Name: "name"},
		},
		{
			name:    "NG Length",
			input:   &ListProjectRequest{ProjectId: 123, UserId: 1, Name: "12345678901234567890123456789012345678901234567890123456789012345"},
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

func TestValidate_CreateProjectRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *CreateProjectRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &CreateProjectRequest{UserId: 1, Name: "name"},
		},
		{
			name:    "NG Required(user_id)",
			input:   &CreateProjectRequest{Name: "name"},
			wantErr: true,
		},
		{
			name:    "NG Required(name)",
			input:   &CreateProjectRequest{UserId: 1},
			wantErr: true,
		},
		{
			name:    "NG Length",
			input:   &CreateProjectRequest{UserId: 1, Name: "12345678901234567890123456789012345678901234567890123456789012345"},
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

func TestValidate_UpdateProjectRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *UpdateProjectRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &UpdateProjectRequest{ProjectId: 1, Name: "name"},
		},
		{
			name:    "NG Required(project_id)",
			input:   &UpdateProjectRequest{Name: "name"},
			wantErr: true,
		},
		{
			name:    "NG Required(name)",
			input:   &UpdateProjectRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG Length",
			input:   &UpdateProjectRequest{ProjectId: 1, Name: "12345678901234567890123456789012345678901234567890123456789012345"},
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

func TestValidate_DeleteProjectRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *DeleteProjectRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &DeleteProjectRequest{ProjectId: 1},
		},
		{
			name:    "NG Required(project_id)",
			input:   &DeleteProjectRequest{},
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

func TestValidate_TagProjectRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *TagProjectRequest
		wantErr bool
	}{
		{
			name:  "OK(full)",
			input: &TagProjectRequest{ProjectId: 1, Tag: "tag", Color: "color"},
		},
		{
			name:  "OK(minimum)",
			input: &TagProjectRequest{ProjectId: 1, Tag: "tag"},
		},
		{
			name:  "OK Max Length(tag)",
			input: &TagProjectRequest{ProjectId: 1, Tag: "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012"},
		},
		{
			name:    "NG Required(project_id)",
			input:   &TagProjectRequest{Tag: "tag"},
			wantErr: true,
		},
		{
			name:    "NG Required(tag)",
			input:   &TagProjectRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG Length(tag)",
			input:   &TagProjectRequest{ProjectId: 1, Tag: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123"},
			wantErr: true,
		},
		{
			name:    "NG Length(color)",
			input:   &TagProjectRequest{ProjectId: 1, Tag: "tag", Color: "123456789012345678901234567890123"},
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

func TestValidate_UntagProjectRequest(t *testing.T) {
	cases := []struct {
		name    string
		input   *UntagProjectRequest
		wantErr bool
	}{
		{
			name:  "OK",
			input: &UntagProjectRequest{ProjectId: 1, Tag: "tag"},
		},
		{
			name:  "OK Max Length(tag)",
			input: &UntagProjectRequest{ProjectId: 1, Tag: "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012"},
		},
		{
			name:    "NG Required(project_id)",
			input:   &UntagProjectRequest{Tag: "tag"},
			wantErr: true,
		},
		{
			name:    "NG Required(tag)",
			input:   &UntagProjectRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:    "NG Length(tag)",
			input:   &UntagProjectRequest{ProjectId: 1, Tag: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123"},
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
