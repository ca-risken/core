package report

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Validate GetReportRequest
func (r *GetReportFindingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

/*
 * entities
**/
