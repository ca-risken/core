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
	ListOrgAlertCondNotification(ctx context.Context, organizationID, projectID, alertConditionID, notificationID, pageSize, pageOffset uint32) ([]*orgalert.OrgAlertCondNotification, error)
	GetOrgAlertCondNotification(ctx context.Context, organizationID, projectID, alertConditionID, notificationID uint32) (*orgalert.OrgAlertCondNotification, error)
	UpdateOrgAlertCondNotificationCache(ctx context.Context, organizationID, projectID, alertConditionID, notificationID, cacheSecond uint32) (*orgalert.OrgAlertCondNotification, error)
	UpdateOrgAlertProjectNotificationEnabled(ctx context.Context, organizationID, projectID, notificationID uint32, enabled bool) ([]*orgalert.OrgAlertCondNotification, error)
	UpdateOrgAlertProjectNotificationCache(ctx context.Context, organizationID, projectID, notificationID, cacheSecond uint32) ([]*orgalert.OrgAlertCondNotification, error)
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
			organization_id, project_id, alert_condition_id, notification_id, cache_second
		)
		select orgn.organization_id, op.project_id, ac.alert_condition_id, orgn.notification_id, 2592000
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
	Enabled          bool
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
	selectOrgAlertCondNotificationForUpdate = `
		select oacn.*
		from organization_alert_cond_notification oacn
		inner join organization_project op
			on op.organization_id = oacn.organization_id
			and op.project_id = oacn.project_id
		inner join alert_condition ac
			on ac.project_id = oacn.project_id
			and ac.alert_condition_id = oacn.alert_condition_id
		inner join organization_notification orgn
			on orgn.organization_id = oacn.organization_id
			and orgn.notification_id = oacn.notification_id
		where oacn.organization_id = ?
		and oacn.project_id = ?
		and oacn.alert_condition_id = ?
		and oacn.notification_id = ?
		for update
	`
	updateOrgAlertCondNotificationCache = `
		update organization_alert_cond_notification
		set cache_second = ?
		where organization_id = ?
		and project_id = ?
		and alert_condition_id = ?
		and notification_id = ?
	`
	selectOrganizationProjectForUpdate = `
		select *
		from organization_project
		where organization_id = ?
		and project_id = ?
		for update
	`
	selectOrgAlertProjectNotificationForUpdate = `
		select oacn.*
		from organization_alert_cond_notification oacn
		inner join organization_project op
			on op.organization_id = oacn.organization_id
			and op.project_id = oacn.project_id
		inner join organization_notification orgn
			on orgn.organization_id = oacn.organization_id
			and orgn.notification_id = oacn.notification_id
		where oacn.organization_id = ?
		and oacn.project_id = ?
		and oacn.notification_id = ?
		for update
	`
	updateOrgAlertProjectNotificationEnabled = `
		update organization_alert_cond_notification
		set enabled = ?
		where organization_id = ?
		and project_id = ?
		and notification_id = ?
	`
	updateOrgAlertProjectNotificationCache = `
		update organization_alert_cond_notification
		set cache_second = ?
		where organization_id = ?
		and project_id = ?
		and notification_id = ?
	`
)

func (c *Client) ListOrgAlertCondNotification(ctx context.Context, organizationID, projectID, alertConditionID, notificationID, pageSize, pageOffset uint32) ([]*orgalert.OrgAlertCondNotification, error) {
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
	query.WriteString(" order by project_id, alert_condition_id, notification_id limit ? offset ?")
	args = append(args, pageSize, pageOffset)
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

func (c *Client) UpdateOrgAlertCondNotificationCache(ctx context.Context, organizationID, projectID, alertConditionID, notificationID, cacheSecond uint32) (*orgalert.OrgAlertCondNotification, error) {
	var data *orgalert.OrgAlertCondNotification
	err := c.Master.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row organizationAlertCondNotification
		if err := tx.Raw(selectOrgAlertCondNotificationForUpdate, organizationID, projectID, alertConditionID, notificationID).First(&row).Error; err != nil {
			return err
		}
		if err := tx.Exec(updateOrgAlertCondNotificationCache, cacheSecond, organizationID, projectID, alertConditionID, notificationID).Error; err != nil {
			return err
		}
		updated, err := c.getOrgAlertCondNotification(ctx, tx, organizationID, projectID, alertConditionID, notificationID)
		if err != nil {
			return err
		}
		data = updated
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) UpdateOrgAlertProjectNotificationEnabled(ctx context.Context, organizationID, projectID, notificationID uint32, enabled bool) ([]*orgalert.OrgAlertCondNotification, error) {
	return c.updateOrgAlertProjectNotification(ctx, organizationID, projectID, notificationID, func(tx *gorm.DB) error {
		return tx.Exec(updateOrgAlertProjectNotificationEnabled, enabled, organizationID, projectID, notificationID).Error
	})
}

func (c *Client) UpdateOrgAlertProjectNotificationCache(ctx context.Context, organizationID, projectID, notificationID, cacheSecond uint32) ([]*orgalert.OrgAlertCondNotification, error) {
	return c.updateOrgAlertProjectNotification(ctx, organizationID, projectID, notificationID, func(tx *gorm.DB) error {
		return tx.Exec(updateOrgAlertProjectNotificationCache, cacheSecond, organizationID, projectID, notificationID).Error
	})
}

func (c *Client) updateOrgAlertProjectNotification(ctx context.Context, organizationID, projectID, notificationID uint32, update func(*gorm.DB) error) ([]*orgalert.OrgAlertCondNotification, error) {
	var data []*orgalert.OrgAlertCondNotification
	err := c.Master.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var organizationProject model.OrganizationProject
		if err := tx.Raw(selectOrganizationProjectForUpdate, organizationID, projectID).First(&organizationProject).Error; err != nil {
			return err
		}
		var rows []*organizationAlertCondNotification
		if err := tx.Raw(selectOrgAlertProjectNotificationForUpdate, organizationID, projectID, notificationID).Scan(&rows).Error; err != nil {
			return err
		}
		if len(rows) == 0 {
			return gorm.ErrRecordNotFound
		}
		if err := update(tx); err != nil {
			return err
		}
		if err := tx.Raw(listOrgAlertCondNotification+" and project_id = ? and notification_id = ? order by alert_condition_id", organizationID, projectID, notificationID).Scan(&rows).Error; err != nil {
			return err
		}
		data = make([]*orgalert.OrgAlertCondNotification, 0, len(rows))
		for _, row := range rows {
			data = append(data, convertOrgAlertCondNotification(row))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
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
	epoch := time.Date(1970, time.January, 1, 0, 0, 0, 0, row.NotifiedAt.Location())
	if row.NotifiedAt.IsZero() || row.NotifiedAt.Equal(epoch) {
		notifiedAt = 0
	}
	return &orgalert.OrgAlertCondNotification{
		OrganizationId:   row.OrganizationID,
		ProjectId:        row.ProjectID,
		AlertConditionId: row.AlertConditionID,
		NotificationId:   row.NotificationID,
		Enabled:          row.Enabled,
		CacheSecond:      row.CacheSecond,
		NotifiedAt:       notifiedAt,
		CreatedAt:        row.CreatedAt.Unix(),
		UpdatedAt:        row.UpdatedAt.Unix(),
	}
}

type OrgAlertNotificationTarget struct {
	OrganizationID   uint32
	OrganizationName string
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
		select oacn.organization_id, org.name as organization_name,
			oacn.project_id, oacn.alert_condition_id,
			oacn.notification_id, oacn.cache_second, oacn.notified_at,
			orgn.type, orgn.notify_setting
		from organization_alert_cond_notification oacn
		inner join organization org
			on org.organization_id = oacn.organization_id
		inner join organization_notification orgn
			on orgn.organization_id = oacn.organization_id
			and orgn.notification_id = oacn.notification_id
		inner join organization_project op
			on op.organization_id = oacn.organization_id
			and op.project_id = oacn.project_id
		where oacn.project_id = ? and oacn.alert_condition_id = ?
		and oacn.enabled = true
		order by oacn.organization_id, oacn.notification_id
	`
	getOrgAlertNotificationTarget = `
		select oacn.organization_id, org.name as organization_name,
			oacn.project_id, oacn.alert_condition_id,
			oacn.notification_id, oacn.cache_second, oacn.notified_at,
			orgn.type, orgn.notify_setting
		from organization_alert_cond_notification oacn
		inner join organization org
			on org.organization_id = oacn.organization_id
		inner join organization_project op
			on op.organization_id = oacn.organization_id
			and op.project_id = oacn.project_id
		inner join organization_notification orgn
			on orgn.organization_id = oacn.organization_id
			and orgn.notification_id = oacn.notification_id
		where oacn.organization_id = ?
		and oacn.project_id = ?
		and oacn.alert_condition_id = ?
		and oacn.notification_id = ?
		and oacn.enabled = true
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

func (c *Client) GetOrgAlertNotificationTarget(ctx context.Context, organizationID, projectID, alertConditionID, notificationID uint32) (*OrgAlertNotificationTarget, error) {
	var target OrgAlertNotificationTarget
	if err := c.Master.WithContext(ctx).Raw(getOrgAlertNotificationTarget, organizationID, projectID, alertConditionID, notificationID).First(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

func (c *Client) UpdateOrgAlertCondNotificationNotifiedAt(ctx context.Context, organizationID, projectID, alertConditionID, notificationID uint32, notifiedAt time.Time) error {
	return c.Master.WithContext(ctx).Exec(updateOrgAlertCondNotificationNotifiedAt, notifiedAt, organizationID, projectID, alertConditionID, notificationID).Error
}
