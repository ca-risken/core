// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: organization/entity.proto

package organization

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on Organization with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Organization) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Organization with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in OrganizationMultiError, or
// nil if none found.
func (m *Organization) ValidateAll() error {
	return m.validate(true)
}

func (m *Organization) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for OrganizationId

	// no validation rules for Name

	// no validation rules for Description

	// no validation rules for CreatedAt

	// no validation rules for UpdatedAt

	if len(errors) > 0 {
		return OrganizationMultiError(errors)
	}

	return nil
}

// OrganizationMultiError is an error wrapping multiple validation errors
// returned by Organization.ValidateAll() if the designated constraints aren't met.
type OrganizationMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrganizationMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrganizationMultiError) AllErrors() []error { return m }

// OrganizationValidationError is the validation error returned by
// Organization.Validate if the designated constraints aren't met.
type OrganizationValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrganizationValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrganizationValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrganizationValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrganizationValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrganizationValidationError) ErrorName() string { return "OrganizationValidationError" }

// Error satisfies the builtin error interface
func (e OrganizationValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrganization.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrganizationValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrganizationValidationError{}

// Validate checks the field values on OrganizationProject with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *OrganizationProject) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OrganizationProject with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// OrganizationProjectMultiError, or nil if none found.
func (m *OrganizationProject) ValidateAll() error {
	return m.validate(true)
}

func (m *OrganizationProject) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for OrganizationId

	// no validation rules for ProjectId

	// no validation rules for CreatedAt

	// no validation rules for UpdatedAt

	if len(errors) > 0 {
		return OrganizationProjectMultiError(errors)
	}

	return nil
}

// OrganizationProjectMultiError is an error wrapping multiple validation
// errors returned by OrganizationProject.ValidateAll() if the designated
// constraints aren't met.
type OrganizationProjectMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrganizationProjectMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrganizationProjectMultiError) AllErrors() []error { return m }

// OrganizationProjectValidationError is the validation error returned by
// OrganizationProject.Validate if the designated constraints aren't met.
type OrganizationProjectValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrganizationProjectValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrganizationProjectValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrganizationProjectValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrganizationProjectValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrganizationProjectValidationError) ErrorName() string {
	return "OrganizationProjectValidationError"
}

// Error satisfies the builtin error interface
func (e OrganizationProjectValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrganizationProject.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrganizationProjectValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrganizationProjectValidationError{}

// Validate checks the field values on OrganizationInvitation with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *OrganizationInvitation) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OrganizationInvitation with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// OrganizationInvitationMultiError, or nil if none found.
func (m *OrganizationInvitation) ValidateAll() error {
	return m.validate(true)
}

func (m *OrganizationInvitation) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for OrganizationId

	// no validation rules for ProjectId

	// no validation rules for Status

	// no validation rules for CreatedAt

	// no validation rules for UpdatedAt

	if len(errors) > 0 {
		return OrganizationInvitationMultiError(errors)
	}

	return nil
}

// OrganizationInvitationMultiError is an error wrapping multiple validation
// errors returned by OrganizationInvitation.ValidateAll() if the designated
// constraints aren't met.
type OrganizationInvitationMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrganizationInvitationMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrganizationInvitationMultiError) AllErrors() []error { return m }

// OrganizationInvitationValidationError is the validation error returned by
// OrganizationInvitation.Validate if the designated constraints aren't met.
type OrganizationInvitationValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrganizationInvitationValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrganizationInvitationValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrganizationInvitationValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrganizationInvitationValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrganizationInvitationValidationError) ErrorName() string {
	return "OrganizationInvitationValidationError"
}

// Error satisfies the builtin error interface
func (e OrganizationInvitationValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrganizationInvitation.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrganizationInvitationValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrganizationInvitationValidationError{}
