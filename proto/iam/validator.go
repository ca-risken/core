package iam

import (
	"errors"
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Validate ListUserRequest
func (l *ListUserRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Name, validation.Length(0, 64)),
	)
}

// Validate GetUserRequest
func (g *GetUserRequest) Validate() error {
	errMsg := "UserId or Sub or UserIdpKey is required."
	return validation.ValidateStruct(g,
		validation.Field(&g.UserId, validation.When(g.Sub == "" && g.UserIdpKey == "", validation.Required.Error(errMsg))),
		validation.Field(&g.Sub, validation.When(g.UserId == 0 && g.UserIdpKey == "", validation.Required.Error(errMsg))),
		validation.Field(&g.UserIdpKey, validation.When(g.UserId == 0 && g.Sub == "", validation.Required.Error(errMsg))),
	)
}

// Validate PutUserRequest
func (p *PutUserRequest) Validate() error {
	if validation.IsEmpty(p.User) {
		return errors.New("Required user parameter")
	}
	return p.User.Validate()
}

// Validate UserForUpsert
func (u *UserForUpsert) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Sub, validation.Required, validation.Length(0, 255)),
		validation.Field(&u.Name, validation.Required, validation.Length(0, 64)),
	)
}

// Validate ListRoleRequest
func (l *ListRoleRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ProjectId, validation.Required),
		validation.Field(&l.Name, validation.Length(0, 64)),
	)
}

// ValidateForAdmin ListRoleRequest
func (l *ListRoleRequest) ValidateForAdmin() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Name, validation.Length(0, 64)),
	)
}

// Validate GetRoleRequest
func (g *GetRoleRequest) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.ProjectId, validation.Required),
		validation.Field(&g.RoleId, validation.Required),
	)
}

// ValidateForAdmin GetRoleRequest
func (g *GetRoleRequest) ValidateForAdmin() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.RoleId, validation.Required),
	)
}

// Validate PutRoleRequest
func (p *PutRoleRequest) Validate() error {
	if validation.IsEmpty(p.Role) {
		return errors.New("Required role parameter")
	}
	if err := validation.ValidateStruct(p,
		validation.Field(&p.ProjectId, validation.In(p.Role.ProjectId)),
	); err != nil {
		return err
	}
	return p.Role.Validate()
}

// Validate RoleForUpsert
func (r *RoleForUpsert) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.Name, validation.Required, validation.Length(0, 64)),
	)
}

// Validate DeleteRoleRequest
func (d *DeleteRoleRequest) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ProjectId, validation.Required),
		validation.Field(&d.RoleId, validation.Required),
	)
}

// Validate AttachRoleRequest
func (a *AttachRoleRequest) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.ProjectId, validation.Required),
		validation.Field(&a.UserId, validation.Required),
		validation.Field(&a.RoleId, validation.Required),
	)
}

// ValidateForAdmin AttachRoleRequest
func (a *AttachRoleRequest) ValidateForAdmin() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.UserId, validation.Required),
		validation.Field(&a.RoleId, validation.Required),
	)
}

// Validate DetachRoleRequest
func (d *DetachRoleRequest) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ProjectId, validation.Required),
		validation.Field(&d.UserId, validation.Required),
		validation.Field(&d.RoleId, validation.Required),
	)
}

// ValidateForAdmin DetachRoleRequest
func (d *DetachRoleRequest) ValidateForAdmin() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.UserId, validation.Required),
		validation.Field(&d.RoleId, validation.Required),
	)
}

// Validate ListPolicyRequest
func (l *ListPolicyRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ProjectId, validation.Required),
		validation.Field(&l.Name, validation.Length(0, 64)),
	)
}

// Validate GetPolicyRequest
func (g *GetPolicyRequest) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.ProjectId, validation.Required),
		validation.Field(&g.PolicyId, validation.Required),
	)
}

// Validate PutPolicyRequest
func (p *PutPolicyRequest) Validate() error {
	if validation.IsEmpty(p.Policy) {
		return errors.New("Required policy parameter")
	}
	if err := validation.ValidateStruct(p,
		validation.Field(&p.ProjectId, validation.In(p.Policy.ProjectId)),
	); err != nil {
		return err
	}
	return p.Policy.Validate()
}

// Validate PolicyForUpsert
func (r *PolicyForUpsert) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.Name, validation.Required, validation.Length(0, 64)),
		validation.Field(&r.ActionPtn, validation.Required, validation.By(compilableRegexp(r.ActionPtn))),
		validation.Field(&r.ResourcePtn, validation.Required, validation.By(compilableRegexp(r.ResourcePtn))),
	)
}

// Validate DeletePolicyRequest
func (d *DeletePolicyRequest) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ProjectId, validation.Required),
		validation.Field(&d.PolicyId, validation.Required),
	)
}

// Validate AttachPolicyRequest
func (a *AttachPolicyRequest) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.ProjectId, validation.Required),
		validation.Field(&a.RoleId, validation.Required),
		validation.Field(&a.PolicyId, validation.Required),
	)
}

// Validate DetachPolicyRequest
func (d *DetachPolicyRequest) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ProjectId, validation.Required),
		validation.Field(&d.RoleId, validation.Required),
		validation.Field(&d.PolicyId, validation.Required),
	)
}

// Validate ListAccessTokenRequest
func (l *ListAccessTokenRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ProjectId, validation.Required),
		validation.Field(&l.Name, validation.Length(0, 64)),
	)
}

// Validate AuthenticateAccessTokenRequest
func (a *AuthenticateAccessTokenRequest) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.ProjectId, validation.Required),
		validation.Field(&a.AccessTokenId, validation.Required),
		validation.Field(&a.PlainTextToken, validation.Required),
	)
}

// Validate PutAccessTokenRequest
func (p *PutAccessTokenRequest) Validate() error {
	if validation.IsEmpty(p.AccessToken) {
		return errors.New("Required access_token parameter")
	}
	if err := validation.ValidateStruct(p,
		validation.Field(&p.ProjectId, validation.In(p.AccessToken.ProjectId)),
	); err != nil {
		return err
	}
	return p.AccessToken.Validate()
}

// Validate AccessTokenForUpsert
func (a *AccessTokenForUpsert) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.ProjectId, validation.Required),
		validation.Field(&a.Name, validation.Required, validation.Length(0, 64)),
		validation.Field(&a.Description, validation.Length(0, 255)),
		validation.Field(&a.ExpiredAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&a.LastUpdatedUserId, validation.Required),
	)
}

// Validate DeleteAccessTokenRequest
func (d *DeleteAccessTokenRequest) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ProjectId, validation.Required),
		validation.Field(&d.AccessTokenId, validation.Required),
	)
}

// Validate AttachAccessTokenRoleRequest
func (a *AttachAccessTokenRoleRequest) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.ProjectId, validation.Required),
		validation.Field(&a.RoleId, validation.Required),
		validation.Field(&a.AccessTokenId, validation.Required),
	)
}

// Validate DetachAccessTokenRoleRequest
func (d *DetachAccessTokenRoleRequest) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ProjectId, validation.Required),
		validation.Field(&d.RoleId, validation.Required),
		validation.Field(&d.AccessTokenId, validation.Required),
	)
}

// Validate ListUserReservedRequest
func (r *ListUserReservedRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.UserIdpKey, validation.Length(0, 255)),
	)
}

// Validate PutUserReservedRequest
func (r *PutUserReservedRequest) Validate() error {
	if validation.IsEmpty(r.UserReserved) {
		return errors.New("Required UserReserved parameter")
	}
	return r.UserReserved.Validate()
}

// Validate DeleteUserReservedRequest
func (r *DeleteUserReservedRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.ReservedId, validation.Required),
	)
}

// Validate UserReservedForUpsert
func (r *UserReservedForUpsert) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.UserIdpKey, validation.Required, validation.Length(0, 255)),
		validation.Field(&r.RoleId, validation.Required),
	)
}

// Validate IsAuthorizedRequest
func (i *IsAuthorizedRequest) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.UserId, validation.Required),
		validation.Field(&i.ProjectId, validation.Required),
		// must format: "<service-name>/<action-name>"
		validation.Field(&i.ActionName, validation.Required, validation.Match(regexp.MustCompile(`^(\w|-)+/(\w|-)+$`))),
		// must format: "<prefix>/<prefix>/.../<resource-name>"
		validation.Field(&i.ResourceName, validation.Required, validation.Match(regexp.MustCompile(`^(\w|-|:|/)+/.+$`))),
	)
}

// Validate IsAuthorizedAdminRequest
func (i *IsAuthorizedAdminRequest) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.UserId, validation.Required),
		// must format: "<service-name>/<action-name>"
		validation.Field(&i.ActionName, validation.Required, validation.Match(regexp.MustCompile(`^(\w|-)+/(\w|-)+$`))),
		// must format: "<prefix>/<prefix>/.../<resource-name>"
		validation.Field(&i.ResourceName, validation.Required, validation.Match(regexp.MustCompile(`^(\w|-|:|/)+/.+$`))),
	)
}

// Validate IsAuthorizedTokenRequest
func (i *IsAuthorizedTokenRequest) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.AccessTokenId, validation.Required),
		validation.Field(&i.ProjectId, validation.Required),
		// must format: "<service-name>/<action-name>"
		validation.Field(&i.ActionName, validation.Required, validation.Match(regexp.MustCompile(`^(\w|-)+/(\w|-)+$`))),
		// must format: "<prefix>/<prefix>/.../<resource-name>"
		validation.Field(&i.ResourceName, validation.Required, validation.Match(regexp.MustCompile(`^(\w|-|:|/)+/.+$`))),
	)
}

// Check the `ptn`(string) that is compilable regexp pattern
func compilableRegexp(ptn string) validation.RuleFunc {
	return func(value interface{}) error {
		s, _ := value.(string)
		if s != ptn {
			return fmt.Errorf("Unexpected string, got: %+v", ptn)
		}
		if _, err := regexp.Compile(ptn); err != nil {
			return fmt.Errorf("Could not regexp complie, pattern=%s, err=%+v", ptn, err)
		}
		return nil
	}
}

// Validate IsAdmin
func (i *IsAdminRequest) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.UserId, validation.Required),
	)
}
