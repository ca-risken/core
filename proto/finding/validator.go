package finding

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	// PutFindingBatchMaxLength is the max number of `finding` data per request.
	PutFindingBatchMaxLength = 50
	// PutResourceBatchMaxLength is the max number of `resource` data per request.
	PutResourceBatchMaxLength = 50
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
		validation.Field(&l.Sort, validation.In(
			"finding_id", "description", "data_source", "resource_name", "score", "updated_at")),
		validation.Field(&l.Direction, validation.In("asc", "desc")),
		validation.Field(&l.Offset, validation.Min(0)),
		validation.Field(&l.Limit, validation.Min(0), validation.Max(200)),
	)
}

// Validate BatchListFindingRequest
func (b *BatchListFindingRequest) Validate() error {
	return validation.ValidateStruct(b,
		validation.Field(&b.ProjectId, validation.Required),
		validation.Field(&b.DataSource, validation.Each(validation.Length(0, 64))),
		validation.Field(&b.ResourceName, validation.Each(validation.Length(0, 255))),
		validation.Field(&b.FromScore, validation.Min(0.0), validation.Max(1.0)),
		validation.Field(&b.ToScore, validation.Min(0.0), validation.Max(1.0)),
		validation.Field(&b.FromAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&b.ToAt, validation.Min(0), validation.Max(253402268399)),   //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&b.Tag, validation.Each(validation.Length(0, 64))),
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

// Validate PutFindingBatchRequest
func (p *PutFindingBatchRequest) Validate() error {
	if validation.IsEmpty(p.Finding) {
		return errors.New("Required finding parameter")
	}
	if validation.IsEmpty(p.ProjectId) {
		return errors.New("Required project_id parameter")
	}
	if err := validation.ValidateStruct(p,
		validation.Field(&p.Finding, validation.Length(1, PutFindingBatchMaxLength))); err != nil {
		return err
	}
	for _, f := range p.Finding {
		pfr := &PutFindingRequest{
			ProjectId: p.ProjectId,
			Finding:   f.Finding,
		}
		if err := pfr.Validate(); err != nil {
			return err
		}
		if f.Recommend != nil {
			if err := f.Recommend.Validate(); err != nil {
				return err
			}
		}
		for _, t := range f.Tag {
			if err := t.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
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
		validation.Field(&l.Sort, validation.In(
			"finding_tag_id", "tag", "updated_at")),
		validation.Field(&l.Direction, validation.In("asc", "desc")),
		validation.Field(&l.Offset, validation.Min(0)),
		validation.Field(&l.Limit, validation.Min(0), validation.Max(200)),
	)
}

// Validate ListFindingTagNameRequest
func (l *ListFindingTagNameRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ProjectId, validation.Required),
		validation.Field(&l.FromAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&l.ToAt, validation.Min(0), validation.Max(253402268399)),   //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&l.Sort, validation.In(
			"finding_tag_id", "tag", "updated_at")),
		validation.Field(&l.Direction, validation.In("asc", "desc")),
		validation.Field(&l.Offset, validation.Min(0)),
		validation.Field(&l.Limit, validation.Min(0), validation.Max(200)),
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

// Validate ClearScoreRequest
func (c *ClearScoreRequest) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.DataSource, validation.Required, validation.Length(0, 64)),
		validation.Field(&c.Tag, validation.Each(validation.Length(0, 64))),
		validation.Field(&c.BeforeAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}

// Validate ListResourceRequest
func (l *ListResourceRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ProjectId, validation.Required),
		validation.Field(&l.ResourceName, validation.Each(validation.Length(0, 255))),
		validation.Field(&l.FromAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&l.ToAt, validation.Min(0), validation.Max(253402268399)),   //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&l.Tag, validation.Each(validation.Length(0, 64))),
		validation.Field(&l.Sort, validation.In(
			"resource_id", "resource_name", "updated_at")),
		validation.Field(&l.Direction, validation.In("asc", "desc")),
		validation.Field(&l.Offset, validation.Min(0)),
		validation.Field(&l.Limit, validation.Min(0), validation.Max(200)),
		validation.Field(&l.Namespace, validation.Length(0, 64)),
		validation.Field(&l.ResourceType, validation.Length(0, 64)),
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

// Validate PutResourceBatchRequest
func (p *PutResourceBatchRequest) Validate() error {
	if validation.IsEmpty(p.Resource) {
		return errors.New("Required reosurce parameter")
	}
	if validation.IsEmpty(p.ProjectId) {
		return errors.New("Required project_id parameter")
	}
	if err := validation.ValidateStruct(p,
		validation.Field(&p.Resource, validation.Length(1, PutResourceBatchMaxLength))); err != nil {
		return err
	}
	for _, r := range p.Resource {
		prr := &PutResourceRequest{
			ProjectId: p.ProjectId,
			Resource:  r.Resource,
		}
		if err := prr.Validate(); err != nil {
			return err
		}
		for _, t := range r.Tag {
			if err := t.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
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
		validation.Field(&l.Sort, validation.In(
			"resource_tag_id", "tag", "updated_at")),
		validation.Field(&l.Direction, validation.In("asc", "desc")),
		validation.Field(&l.Offset, validation.Min(0)),
		validation.Field(&l.Limit, validation.Min(0), validation.Max(200)),
	)
}

// Validate ListResourceTagNameRequest
func (l *ListResourceTagNameRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ProjectId, validation.Required),
		validation.Field(&l.FromAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&l.ToAt, validation.Min(0), validation.Max(253402268399)),   //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&l.Sort, validation.In(
			"resource_tag_id", "tag", "updated_at")),
		validation.Field(&l.Direction, validation.In("asc", "desc")),
		validation.Field(&l.Offset, validation.Min(0)),
		validation.Field(&l.Limit, validation.Min(0), validation.Max(200)),
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

// Validate for ListFindingSettingRequest
func (l *ListFindingSettingRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ProjectId, validation.Required),
	)
}

// Validate for GetFindingSettingRequest
func (g *GetFindingSettingRequest) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.ProjectId, validation.Required),
		validation.Field(&g.FindingSettingId, validation.Required),
	)
}

// Validate for PutFindingSettingRequest
func (p *PutFindingSettingRequest) Validate() error {
	if validation.IsEmpty(p.FindingSetting) {
		return errors.New("Required finding_setting parameter")
	}
	if err := validation.ValidateStruct(p,
		validation.Field(&p.ProjectId, validation.In(p.FindingSetting.ProjectId)),
	); err != nil {
		return err
	}
	return p.FindingSetting.Validate()
}

// Validate for DeleteFindingSettingRequest
func (d *DeleteFindingSettingRequest) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ProjectId, validation.Required),
		validation.Field(&d.FindingSettingId, validation.Required),
	)
}

// Validate for GetRecommendRequest
func (g *GetRecommendRequest) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.ProjectId, validation.Required),
		validation.Field(&g.FindingId, validation.Required),
	)
}

// Validate for PutRecommendRequest
func (p *PutRecommendRequest) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.ProjectId, validation.Required),
		validation.Field(&p.FindingId, validation.Required),
		validation.Field(&p.DataSource, validation.Required, validation.Length(0, 64)),
		validation.Field(&p.Type, validation.Required, validation.Length(0, 128)),
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
		validation.Field(&f.OriginalScore, validation.Min(0.0), validation.Max(f.OriginalMaxScore)),
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
func (p *PendFindingForUpsert) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.FindingId, validation.Required),
		validation.Field(&p.ProjectId, validation.Required),
		validation.Field(&p.Note, validation.Length(0, 128)),
	)
}

// Validate for FindingSettingForUpsert
func (f *FindingSettingForUpsert) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.ProjectId, validation.Required),
		validation.Field(&f.ResourceName, validation.Required),
		validation.Field(&f.Setting, validation.Required, is.JSON),
	)
}

// Validate for RecommendForBatch
func (r *RecommendForBatch) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Type, validation.Required, validation.Length(0, 128)),
	)
}

// Validate for FindingTagForBatch
func (f *FindingTagForBatch) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Tag, validation.Required, validation.Length(0, 64)),
	)
}

// Validate for ResourceTagForBatch
func (r *ResourceTagForBatch) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Tag, validation.Required, validation.Length(0, 64)),
	)
}
