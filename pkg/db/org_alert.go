package db

import (
	"context"
	"strings"
	"time"

	"github.com/ca-risken/core/pkg/model"
	orgalert "github.com/ca-risken/core/proto/org_alert"
	"gorm.io/gorm"
)

type OrgAlertRepository interface {
	// OrganizationNotification
	ListOrgNotification(ctx context.Context, organizationID uint32) ([]*model.OrganizationNotification, error)
	GetOrgNotification(ctx context.Context, organizationID, notificationID uint32) (*model.OrganizationNotification, error)
	UpsertOrgNotification(ctx context.Context, data *model.OrganizationNotification) (*model.OrganizationNotification, error)
	DeleteOrgNotification(ctx context.Context, organizationID, notificationID uint32) error
	ListOrgNotificationByProjectID(ctx context.Context, projectID uint32) ([]*model.OrganizationNotification, error)
	ListOrgAlertCondNotification(ctx context.Context, organizationID, projectID, alertConditionID, notificationID uint32) ([]*orgalert.OrgAlertCondNotification, error)
	GetOrgAlertCondNotification(ctx context.Context, organizationID, projectID, alertConditionID, notificationID uint32) (*orgalert.OrgAlertCondNotification, error)
	InsertOrgAlertCondNotification(ctx context.Context, organizationID, projectID, alertConditionID, notificationID uint32) error
	DeleteOrgAlertCondNotification(ctx context.Context, organizationID, projectID, alertConditionID, notificationID uint32) error
	UpdateOrgAlertCondNotificationCache(ctx context.Context, organizationID, projectID, alertConditionID, notificationID, cacheSecond uint32) (*orgalert.OrgAlertCondNotification, error)
	ExistsOrgAlertConditionMembership(ctx context.Context, organizationID, projectID, alertConditionID uint32) (bool, error)
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
	err := c.Master.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("organization_id = ? AND notification_id = ?", data.OrganizationID, data.NotificationID).Assign(data).FirstOrCreate(&retData).Error; err != nil {
			return err
		}
		return tx.Exec(insertOrgAlertCondNotificationByOrgNotification, retData.OrganizationID, retData.NotificationID).Error
	})
	if err != nil {
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
	return c.Master.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(deleteOrgAlertCondNotificationByOrgNotification, organizationID, notificationID).Error; err != nil {
			return err
		}
		return tx.Exec(deleteOrgNotification, organizationID, notificationID).Error
	})
}

const (
	insertOrgAlertCondNotificationByOrgNotification = `
		insert ignore into organization_alert_cond_notification (
			organization_id, project_id, alert_condition_id, notification_id
		)
		select orgn.organization_id, op.project_id, ac.alert_condition_id, orgn.notification_id
		from organization_notification orgn
		inner join organization_project op on op.organization_id = orgn.organization_id
		inner join alert_condition ac on ac.project_id = op.project_id
		where orgn.organization_id = ? and orgn.notification_id = ?
	`
	deleteOrgAlertCondNotificationByOrgNotification = `
		delete from organization_alert_cond_notification
		where organization_id = ? and notification_id = ?
	`
)

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

type organizationAlertCondNotification struct {
	OrganizationID   uint32
	ProjectID        uint32
	AlertConditionID uint32
	NotificationID   uint32
	CacheSecond      uint32
	NotifiedAt       time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

const (
	listOrgAlertCondNotification = `
		select *
		from organization_alert_cond_notification
		where organization_id = ?
	`
	selectOrgAlertCondNotification = `
		select *
		from organization_alert_cond_notification
		where organization_id = ?
		and project_id = ?
		and alert_condition_id = ?
		and notification_id = ?
	`
	insertOrgAlertCondNotification = `
		insert ignore into organization_alert_cond_notification (
			organization_id,
			project_id,
			alert_condition_id,
			notification_id
		) values (?, ?, ?, ?)
	`
	deleteOrgAlertCondNotification = `
		delete from organization_alert_cond_notification
		where organization_id = ?
		and project_id = ?
		and alert_condition_id = ?
		and notification_id = ?
	`
	updateOrgAlertCondNotificationCache = `
		update organization_alert_cond_notification
		set cache_second = ?
		where organization_id = ?
		and project_id = ?
		and alert_condition_id = ?
		and notification_id = ?
	`
	selectOrgAlertConditionMembership = `
		select exists (
			select 1
			from organization_project op
			inner join alert_condition ac on ac.project_id = op.project_id
			where op.organization_id = ?
			and op.project_id = ?
			and ac.alert_condition_id = ?
		)
	`
)

func (c *Client) ListOrgAlertCondNotification(ctx context.Context, organizationID, projectID, alertConditionID, notificationID uint32) ([]*orgalert.OrgAlertCondNotification, error) {
	query := strings.Builder{}
	query.WriteString(listOrgAlertCondNotification)
	args := []interface{}{organizationID}
	filters := []struct {
		column string
		value  uint32
	}{
		{column: "project_id", value: projectID},
		{column: "alert_condition_id", value: alertConditionID},
		{column: "notification_id", value: notificationID},
	}
	for _, filter := range filters {
		if filter.value == 0 {
			continue
		}
		query.WriteString(" and ")
		query.WriteString(filter.column)
		query.WriteString(" = ?")
		args = append(args, filter.value)
	}
	query.WriteString(" order by project_id, alert_condition_id, notification_id limit 1000")
	var rows []*organizationAlertCondNotification
	if err := c.Slave.WithContext(ctx).Raw(query.String(), args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	data := make([]*orgalert.OrgAlertCondNotification, 0, len(rows))
	for _, row := range rows {
		data = append(data, convertOrgAlertCondNotification(row))
	}
	return data, nil
}

func (c *Client) GetOrgAlertCondNotification(ctx context.Context, organizationID, projectID, alertConditionID, notificationID uint32) (*orgalert.OrgAlertCondNotification, error) {
	return c.getOrgAlertCondNotification(ctx, c.Slave, organizationID, projectID, alertConditionID, notificationID)
}

func (c *Client) InsertOrgAlertCondNotification(ctx context.Context, organizationID, projectID, alertConditionID, notificationID uint32) error {
	return c.Master.WithContext(ctx).Exec(insertOrgAlertCondNotification, organizationID, projectID, alertConditionID, notificationID).Error
}

func (c *Client) DeleteOrgAlertCondNotification(ctx context.Context, organizationID, projectID, alertConditionID, notificationID uint32) error {
	return c.Master.WithContext(ctx).Exec(deleteOrgAlertCondNotification, organizationID, projectID, alertConditionID, notificationID).Error
}

func (c *Client) UpdateOrgAlertCondNotificationCache(ctx context.Context, organizationID, projectID, alertConditionID, notificationID, cacheSecond uint32) (*orgalert.OrgAlertCondNotification, error) {
	if err := c.Master.WithContext(ctx).Exec(updateOrgAlertCondNotificationCache, cacheSecond, organizationID, projectID, alertConditionID, notificationID).Error; err != nil {
		return nil, err
	}
	return c.getOrgAlertCondNotification(ctx, c.Master, organizationID, projectID, alertConditionID, notificationID)
}

func (c *Client) ExistsOrgAlertConditionMembership(ctx context.Context, organizationID, projectID, alertConditionID uint32) (bool, error) {
	var exists bool
	if err := c.Slave.WithContext(ctx).Raw(selectOrgAlertConditionMembership, organizationID, projectID, alertConditionID).Scan(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

func (c *Client) getOrgAlertCondNotification(ctx context.Context, database *gorm.DB, organizationID, projectID, alertConditionID, notificationID uint32) (*orgalert.OrgAlertCondNotification, error) {
	var row organizationAlertCondNotification
	if err := database.WithContext(ctx).Raw(selectOrgAlertCondNotification, organizationID, projectID, alertConditionID, notificationID).First(&row).Error; err != nil {
		return nil, err
	}
	return convertOrgAlertCondNotification(&row), nil
}

func convertOrgAlertCondNotification(row *organizationAlertCondNotification) *orgalert.OrgAlertCondNotification {
	if row == nil {
		return &orgalert.OrgAlertCondNotification{}
	}
	notifiedAt := row.NotifiedAt.Unix()
	if row.NotifiedAt.IsZero() || (row.NotifiedAt.Year() == 1970 && row.NotifiedAt.YearDay() == 1 && row.NotifiedAt.Hour() == 0 && row.NotifiedAt.Minute() == 0 && row.NotifiedAt.Second() == 0) {
		notifiedAt = 0
	}
	return &orgalert.OrgAlertCondNotification{
		OrganizationId:   row.OrganizationID,
		ProjectId:        row.ProjectID,
		AlertConditionId: row.AlertConditionID,
		NotificationId:   row.NotificationID,
		CacheSecond:      row.CacheSecond,
		NotifiedAt:       notifiedAt,
		CreatedAt:        row.CreatedAt.Unix(),
		UpdatedAt:        row.UpdatedAt.Unix(),
	}
}

type OrgAlertNotificationTarget struct {
	OrganizationID   uint32
	ProjectID        uint32
	AlertConditionID uint32
	NotificationID   uint32
	CacheSecond      uint32
	NotifiedAt       time.Time
	Type             string
	NotifySetting    string
}

const (
	listOrgAlertNotificationTarget = `
		select oacn.organization_id, oacn.project_id, oacn.alert_condition_id,
			oacn.notification_id, oacn.cache_second, oacn.notified_at,
			orgn.type, orgn.notify_setting
		from organization_alert_cond_notification oacn
		inner join organization_notification orgn
			on orgn.organization_id = oacn.organization_id
			and orgn.notification_id = oacn.notification_id
		inner join organization_project op
			on op.organization_id = oacn.organization_id
			and op.project_id = oacn.project_id
		where oacn.project_id = ? and oacn.alert_condition_id = ?
		order by oacn.organization_id, oacn.notification_id
	`
	existsOrgAlertNotificationTarget = `
		select exists (
			select 1
			from organization_project op
			inner join organization_alert_cond_notification oacn
				on oacn.organization_id = op.organization_id
				and oacn.project_id = op.project_id
			where oacn.organization_id = ?
			and oacn.project_id = ?
			and oacn.alert_condition_id = ?
			and oacn.notification_id = ?
		)
	`
	updateOrgAlertCondNotificationNotifiedAt = `
		update organization_alert_cond_notification
		set notified_at = ?
		where organization_id = ?
		and project_id = ?
		and alert_condition_id = ?
		and notification_id = ?
	`
)

func (c *Client) ListOrgAlertNotificationTarget(ctx context.Context, projectID, alertConditionID uint32) ([]*OrgAlertNotificationTarget, error) {
	var targets []*OrgAlertNotificationTarget
	if err := c.Slave.WithContext(ctx).Raw(listOrgAlertNotificationTarget, projectID, alertConditionID).Scan(&targets).Error; err != nil {
		return nil, err
	}
	return targets, nil
}

func (c *Client) ExistsOrgAlertNotificationTarget(ctx context.Context, organizationID, projectID, alertConditionID, notificationID uint32) (bool, error) {
	var exists bool
	if err := c.Master.WithContext(ctx).Raw(existsOrgAlertNotificationTarget, organizationID, projectID, alertConditionID, notificationID).Scan(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

func (c *Client) UpdateOrgAlertCondNotificationNotifiedAt(ctx context.Context, organizationID, projectID, alertConditionID, notificationID uint32, notifiedAt time.Time) error {
	return c.Master.WithContext(ctx).Exec(updateOrgAlertCondNotificationNotifiedAt, notifiedAt, organizationID, projectID, alertConditionID, notificationID).Error
}
