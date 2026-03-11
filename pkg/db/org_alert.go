package db

import (
	"context"

	"github.com/ca-risken/core/pkg/model"
)

type OrgAlertRepository interface {
	// OrganizationNotification
	ListOrgNotification(ctx context.Context, organizationID uint32) ([]*model.OrganizationNotification, error)
	GetOrgNotification(ctx context.Context, organizationID, notificationID uint32) (*model.OrganizationNotification, error)
	UpsertOrgNotification(ctx context.Context, data *model.OrganizationNotification) (*model.OrganizationNotification, error)
	DeleteOrgNotification(ctx context.Context, organizationID, notificationID uint32) error
	ListOrgNotificationByProjectID(ctx context.Context, projectID uint32) ([]*model.OrganizationNotification, error)
}

var _ OrgAlertRepository = (*Client)(nil)

func (c *Client) ListOrgNotification(ctx context.Context, organizationID uint32) ([]*model.OrganizationNotification, error) {
	query := `select * from organization_notification where organization_id = ? order by notification_id`
	var data []*model.OrganizationNotification
	if err := c.Slave.WithContext(ctx).Raw(query, organizationID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

const selectGetOrgNotification = `
	select *
	from organization_notification
	where organization_id = ?
	and notification_id = ?
`

func (c *Client) GetOrgNotification(ctx context.Context, organizationID, notificationID uint32) (*model.OrganizationNotification, error) {
	var data model.OrganizationNotification
	if err := c.Slave.WithContext(ctx).Raw(selectGetOrgNotification, organizationID, notificationID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) UpsertOrgNotification(ctx context.Context, data *model.OrganizationNotification) (*model.OrganizationNotification, error) {
	var retData model.OrganizationNotification
	if err := c.Master.WithContext(ctx).Where("organization_id = ? AND notification_id = ?", data.OrganizationID, data.NotificationID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
		return nil, err
	}
	return &retData, nil
}

const deleteOrgNotification = `
	delete from organization_notification
	where organization_id = ?
	and notification_id = ?
`

func (c *Client) DeleteOrgNotification(ctx context.Context, organizationID, notificationID uint32) error {
	return c.Master.WithContext(ctx).Exec(deleteOrgNotification, organizationID, notificationID).Error
}

const selectOrgNotificationByProjectID = `
	select orgn.*
	from organization_notification orgn
	inner join organization_project op on orgn.organization_id = op.organization_id
	where op.project_id = ?
	order by orgn.notification_id
`

func (c *Client) ListOrgNotificationByProjectID(ctx context.Context, projectID uint32) ([]*model.OrganizationNotification, error) {
	var data []*model.OrganizationNotification
	if err := c.Slave.WithContext(ctx).Raw(selectOrgNotificationByProjectID, projectID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

