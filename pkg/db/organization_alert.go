package db

import (
	"context"

	"github.com/ca-risken/core/pkg/model"
)

type OrganizationAlertRepository interface {
	// OrganizationNotification
	ListOrganizationNotification(ctx context.Context, organizationID uint32, notifyType string) ([]*model.OrganizationNotification, error)
	GetOrganizationNotification(ctx context.Context, organizationID, notificationID uint32) (*model.OrganizationNotification, error)
	UpsertOrganizationNotification(ctx context.Context, data *model.OrganizationNotification) (*model.OrganizationNotification, error)
	DeleteOrganizationNotification(ctx context.Context, organizationID, notificationID uint32) error
	ListOrganizationNotificationByProjectID(ctx context.Context, projectID uint32) ([]*model.OrganizationNotification, error)
}

var _ OrganizationAlertRepository = (*Client)(nil)

func (c *Client) ListOrganizationNotification(ctx context.Context, organizationID uint32, notifyType string) ([]*model.OrganizationNotification, error) {
	query := `select * from organization_notification where organization_id = ?`
	var params []interface{}
	params = append(params, organizationID)
	if notifyType != "" {
		query += " and type = ?"
		params = append(params, notifyType)
	}
	query += " order by notification_id"
	var data []*model.OrganizationNotification
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
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

const putOrganizationNotification = `
	INSERT INTO organization_notification (
		notification_id,
		name,
		organization_id,
		type,
		notify_setting
	) VALUES (
		?, ?, ?, ?, ?
	) ON DUPLICATE KEY UPDATE
		name = VALUES(name),
		type = VALUES(type),
		notify_setting = VALUES(notify_setting)
`

func (c *Client) UpsertOrganizationNotification(ctx context.Context, data *model.OrganizationNotification) (*model.OrganizationNotification, error) {
	if err := c.Master.WithContext(ctx).Exec(putOrganizationNotification,
		data.NotificationID, data.Name, data.OrganizationID, data.Type, data.NotifySetting,
	).Error; err != nil {
		return nil, err
	}
	return c.getOrganizationNotificationFromMaster(ctx, data.OrganizationID, data.NotificationID)
}

func (c *Client) getOrganizationNotificationFromMaster(ctx context.Context, organizationID, notificationID uint32) (*model.OrganizationNotification, error) {
	var data model.OrganizationNotification
	if err := c.Master.WithContext(ctx).Raw(selectGetOrganizationNotification, organizationID, notificationID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const deleteOrganizationNotification = `
	delete from organization_notification
	where organization_id = ?
	and notification_id = ?
`

func (c *Client) DeleteOrganizationNotification(ctx context.Context, organizationID, notificationID uint32) error {
	return c.Master.WithContext(ctx).Exec(deleteOrganizationNotification, organizationID, notificationID).Error
}

const listOrganizationNotificationByProjectID = `
	select orgn.*
	from organization_notification orgn
	join organization_project op on orgn.organization_id = op.organization_id
	where op.project_id = ?
	order by orgn.notification_id
`

func (c *Client) ListOrganizationNotificationByProjectID(ctx context.Context, projectID uint32) ([]*model.OrganizationNotification, error) {
	var data []*model.OrganizationNotification
	if err := c.Slave.WithContext(ctx).Raw(listOrganizationNotificationByProjectID, projectID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
