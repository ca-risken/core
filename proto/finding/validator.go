package finding

import (
	"errors"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Validate ListFindingRequest
func (l *ListFindingRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ProjectId, validation.Required),
		validation.Field(&l.DataSource, validation.Each(validation.Length(0, 64))),
		validation.Field(&l.ResourceName, validation.Each(validation.Length(0, 255))),
		validation.Field(&l.FromScore, validation.Min(0.0), validation.Max(1.0)),
		validation.Field(&l.ToScore, validation.Min(0.0), validation.Max(1.0)),
		validation.Field(&l.FromAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&l.ToAt, validation.Min(0), validation.Max(253402268399)),   //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&l.Tag, validation.Each(validation.Length(0, 64))),
	)
}

// Validate GetFinding
func (g *GetFindingRequest) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.ProjectId, validation.Required),
		validation.Field(&g.FindingId, validation.Required),
	)
}

// Validate PutFindingRequest
func (p *PutFindingRequest) Validate() error {
	if validation.IsEmpty(p.Finding) {
		return errors.New("Required finding parameter")
	}
	if err := validation.ValidateStruct(p,
		validation.Field(&p.ProjectId, validation.In(p.Finding.ProjectId)),
	); err != nil {
		return err
	}
	return p.Finding.Validate()
}

// Validate DeleteFindingRequest
func (d *DeleteFindingRequest) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ProjectId, validation.Required),
		validation.Field(&d.FindingId, validation.Required),
	)
}

// Validate ListFindingTagRequest
func (l *ListFindingTagRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ProjectId, validation.Required),
		validation.Field(&l.FindingId, validation.Required),
	)
}

// Validate ListFindingTagNameRequest
func (l *ListFindingTagNameRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ProjectId, validation.Required),
		validation.Field(&l.FromAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&l.ToAt, validation.Min(0), validation.Max(253402268399)),   //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}

// Validate TagFindingRequest
func (t *TagFindingRequest) Validate() error {
	if validation.IsEmpty(t.Tag) {
		return errors.New("Required tag parameter")
	}
	if err := validation.ValidateStruct(t,
		validation.Field(&t.ProjectId, validation.Required, validation.In(t.Tag.ProjectId)),
	); err != nil {
		return err
	}
	return t.Tag.Validate()
}

// Validate UntagFindingRequest
func (u *UntagFindingRequest) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.ProjectId, validation.Required),
		validation.Field(&u.FindingTagId, validation.Required),
	)
}

// Validate ListResourceRequest
func (l *ListResourceRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ProjectId, validation.Required),
		validation.Field(&l.ResourceName, validation.Each(validation.Length(0, 255))),
		validation.Field(&l.FromSumScore, validation.Min(0.0)),
		validation.Field(&l.ToSumScore, validation.Min(0.0)),
		validation.Field(&l.FromAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&l.ToAt, validation.Min(0), validation.Max(253402268399)),   //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&l.Tag, validation.Each(validation.Length(0, 64))),
	)
}

// Validate GetResourceRequest
func (g *GetResourceRequest) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.ProjectId, validation.Required),
		validation.Field(&g.ResourceId, validation.Required),
	)
}

// Validate PutResourceRequest
func (p *PutResourceRequest) Validate() error {
	if validation.IsEmpty(p.Resource) {
		return errors.New("Required resource parameter")
	}
	if err := validation.ValidateStruct(p,
		validation.Field(&p.ProjectId, validation.In(p.Resource.ProjectId)),
	); err != nil {
		return err
	}
	return p.Resource.Validate()
}

// Validate DeleteResourceRequest
func (d *DeleteResourceRequest) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ProjectId, validation.Required),
		validation.Field(&d.ResourceId, validation.Required),
	)
}

// Validate ListResourceTagRequest
func (l *ListResourceTagRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ProjectId, validation.Required),
		validation.Field(&l.ResourceId, validation.Required),
	)
}

// Validate ListResourceTagNameRequest
func (l *ListResourceTagNameRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ProjectId, validation.Required),
		validation.Field(&l.FromAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&l.ToAt, validation.Min(0), validation.Max(253402268399)),   //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}

// Validate TagResourceRequest
func (t *TagResourceRequest) Validate() error {
	if validation.IsEmpty(t.Tag) {
		return errors.New("Required tag parameter")
	}
	if err := validation.ValidateStruct(t,
		validation.Field(&t.ProjectId, validation.Required, validation.In(t.Tag.ProjectId)),
	); err != nil {
		return err
	}
	return t.Tag.Validate()
}

// Validate UntagResourceRequest
func (u *UntagResourceRequest) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.ProjectId, validation.Required),
		validation.Field(&u.ResourceTagId, validation.Required),
	)
}

// Validate for GetPendFindingRequest
func (g *GetPendFindingRequest) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.ProjectId, validation.Required),
		validation.Field(&g.FindingId, validation.Required),
	)
}

// Validate for PutPendFindingRequest
func (p *PutPendFindingRequest) Validate() error {
	if validation.IsEmpty(p.PendFinding) {
		return errors.New("Required pend_finding parameter")
	}
	if err := validation.ValidateStruct(p,
		validation.Field(&p.ProjectId, validation.In(p.PendFinding.ProjectId)),
	); err != nil {
		return err
	}
	return p.PendFinding.Validate()
}

// Validate for DeletePendFindingRequest
func (d *DeletePendFindingRequest) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ProjectId, validation.Required),
		validation.Field(&d.FindingId, validation.Required),
	)
}

/*
 * entities
**/

// Validate FindingForUpsert
func (f *FindingForUpsert) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Description, validation.Length(0, 200)),
		validation.Field(&f.DataSource, validation.Required, validation.Length(0, 64)),
		validation.Field(&f.DataSourceId, validation.Required, validation.Length(0, 255)),
		validation.Field(&f.ResourceName, validation.Required, validation.Length(0, 255)),
		validation.Field(&f.OriginalScore, validation.NilOrNotEmpty, validation.Min(0.0), validation.Max(f.OriginalMaxScore)),
		validation.Field(&f.OriginalMaxScore, validation.NilOrNotEmpty, validation.Min(0.0), validation.Max(999.99)),
		validation.Field(&f.Data, is.JSON),
	)
}

// Validate FindingTagForUpsert
func (f *FindingTagForUpsert) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.FindingId, validation.Required),
		validation.Field(&f.Tag, validation.Required, validation.Length(0, 64)),
	)
}

// Validate ResourceForUpsert
func (r *ResourceForUpsert) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ResourceName, validation.Required, validation.Length(0, 255)),
	)
}

// Validate ResourceTagForUpsert
func (r *ResourceTagForUpsert) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ResourceId, validation.Required),
		validation.Field(&r.Tag, validation.Required, validation.Length(0, 64)),
	)
}

// Validate for PendFindingForUpsert
func (r *PendFindingForUpsert) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.FindingId, validation.Required),
		validation.Field(&r.ProjectId, validation.Required),
	)
}
