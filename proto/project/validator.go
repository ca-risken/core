package project

// import (
// 	validation "github.com/go-ozzo/ozzo-validation/v4"
// )

// // Validate ListProjectRequest
// func (l *ListProjectRequest) Validate() error {
// 	return validation.ValidateStruct(l,
// 		validation.Field(&l.Name, validation.Length(0, 64)),
// 	)
// }

// // Validate CreateProjectRequest
// func (c *CreateProjectRequest) Validate() error {
// 	return validation.ValidateStruct(c,
// 		validation.Field(&c.UserId, validation.Required),
// 		validation.Field(&c.Name, validation.Required, validation.Length(0, 64)),
// 	)
// }

// // Validate UpdateProjectRequest
// func (u *UpdateProjectRequest) Validate() error {
// 	return validation.ValidateStruct(u,
// 		validation.Field(&u.ProjectId, validation.Required),
// 		validation.Field(&u.Name, validation.Required, validation.Length(0, 64)),
// 	)
// }

// // Validate DeleteProjectRequest
// func (d *DeleteProjectRequest) Validate() error {
// 	return validation.ValidateStruct(d,
// 		validation.Field(&d.ProjectId, validation.Required),
// 	)
// }
