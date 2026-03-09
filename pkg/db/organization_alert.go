package db

import (
	"context"

	"github.com/ca-risken/core/pkg/model"
)

type OrganizationAlertRepository interface {
	// OrganizationNotification
	ListOrganizationNotification(ctx context.Context, organizationID uint32) ([]*model.OrganizationNotification, error)
	GetOrganizationNotification(ctx context.Context, organizationID, notificationID uint32) (*model.OrganizationNotification, error)
	UpsertOrganizationNotification(ctx context.Context, data *model.OrganizationNotification) (*model.OrganizationNotification, error)
	DeleteOrganizationNotification(ctx context.Context, organizationID, notificationID uint32) error
}

var _ OrganizationAlertRepository = (*Client)(nil)

func (c *Client) ListOrganizationNotification(ctx context.Context, organizationID uint32) ([]*model.OrganizationNotification, error) {
	query := `select * from organization_notification where organization_id = ? order by notification_id`
	var data []*model.OrganizationNotification
	if err := c.Slave.WithContext(ctx).Raw(query, organizationID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

const selectGetOrganizationNotification = `
	select *
	from organization_notification
	where organization_id = ?
	and notification_id = ?
`

func (c *Client) GetOrganizationNotification(ctx context.Context, organizationID, notificationID uint32) (*model.OrganizationNotification, error) {
	var data model.OrganizationNotification
	if err := c.Slave.WithContext(ctx).Raw(selectGetOrganizationNotification, organizationID, notificationID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertOrganizationNotification(ctx context.Context, data *model.OrganizationNotification) (*model.OrganizationNotification, error) {
	var retData model.OrganizationNotification
	if err := c.Master.WithContext(ctx).Where("organization_id = ? AND notification_id = ?", data.OrganizationID, data.NotificationID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

const deleteOrganizationNotification = `
	delete from organization_notification
	where organization_id = ?
	and notification_id = ?
`

func (c *Client) DeleteOrganizationNotification(ctx context.Context, organizationID, notificationID uint32) error {
	return c.Master.WithContext(ctx).Exec(deleteOrganizationNotification, organizationID, notificationID).Error
}

