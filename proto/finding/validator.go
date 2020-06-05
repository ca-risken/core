package finding

import (
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Validate ListFindingRequest
func (l *ListFindingRequest) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(l.DataSource, validation.Each(validation.Length(0, 64))),
		validation.Field(l.ResourceName, validation.Each(validation.Length(0, 200))),
		validation.Field(l.FromScore, validation.Min(0), validation.Max(1)),
		validation.Field(l.ToScore, validation.Min(0), validation.Max(1)),
	)
}

// Validate GetFinding
func (g *GetFindingRequest) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(g.FindingId, validation.Required),
	)
}

// Validate PutFindingRequest
func (p *PutFindingRequest) Validate() error {
	return p.Finding.Validate()
}

// Validate DeleteFindingRequest
func (d *DeleteFindingRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(d.FindingId, validation.Required),
	)
}

// Validate ListFindingTagRequest
func (l *ListFindingTagRequest) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(l.FindingId, validation.Required),
	)
}

// Validate TagFindingRequest
func (t *TagFindingRequest) Validate() error {
	return t.Tag.Validate()
}

// Validate UntagFindingRequest
func (u *UntagFindingRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(u.FindingTagId, validation.Required),
	)
}

// Validate ListResourceRequest
func (l *ListResourceRequest) Validate() error {
	return nil
}

// Validate GetResourceRequest
func (g *GetResourceRequest) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(g.ResourceId, validation.Required),
	)
}

// Validate PutResourceRequest
func (p *PutResourceRequest) Validate() error {
	return p.Resource.Validate()
}

// Validate DeleteResourceRequest
func (d *DeleteResourceRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(d.ResourceId, validation.Required),
	)
}

// Validate ListResourceTagRequest
func (l *ListResourceTagRequest) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(l.ResourceId, validation.Required),
	)
}

// Validate TagResourceRequest
func (t *TagResourceRequest) Validate() error {
	return t.Tag.Validate()
}

// Validate UntagResourceRequest
func (u *UntagResourceRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(u.ResourceTagId, validation.Required),
	)
}

/*
 * entities
**/

// Validate FindingForUpsert
func (f *FindingForUpsert) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(f.Description, validation.Length(0, 200)),
		validation.Field(f.DataSource, validation.Required, validation.Length(0, 64)),
		validation.Field(f.DataSourceId, validation.Required, validation.Length(0, 255)),
		validation.Field(f.ResourceName, validation.Required),
		validation.Field(f.OriginalScore, validation.Required, validation.Length(0, 100)),
		validation.Field(f.OriginalMaxScore, validation.Required, validation.Length(0, 100)),
		validation.Field(f.Data, is.JSON),
	)
}

// Validate FindingTagForUpsert
func (f *FindingTagForUpsert) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(f.FindingId, validation.Required),
		validation.Field(f.TagKey, validation.Required, validation.Length(0, 64)),
		validation.Field(f.TagKey, validation.Length(0, 200)),
	)
}

// Validate ResourceForUpsert
func (r *ResourceForUpsert) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(r.ResourceName, validation.Required),
	)
}

// Validate ResourceTagForUpsert
func (r *ResourceTagForUpsert) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(r.ResourceId, validation.Required),
		validation.Field(r.TagKey, validation.Required, validation.Length(0, 64)),
		validation.Field(r.TagValue, validation.Length(0, 200)),
	)
}
