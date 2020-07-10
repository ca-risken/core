package iam

import (
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/vikyd/zero"
)

// Validate GetUserRequest
func (g *GetUserRequest) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.UserId, validation.When(zero.IsZeroVal(g.Sub), validation.Required.Error("UserId or Sub is required."))),
		validation.Field(&g.Sub, validation.When(zero.IsZeroVal(g.UserId), validation.Required.Error("UserId or Sub is required."))),
	)
}

// Validate ListFindingRequest
func (i *IsAuthorizedRequest) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.UserId, validation.Required),
		validation.Field(&i.ProjectId, validation.Required),
		// validation.Field(&i.ActionName, validation.Required, validation.By(compilableRegexp(i.ActionName))),
		// validation.Field(&i.ResourceName, validation.Required, validation.By(compilableRegexp(i.ResourceName))),
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
