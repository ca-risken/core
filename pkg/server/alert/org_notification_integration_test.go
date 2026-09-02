//go:build integration

package alert

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/pkg/model"
	projectproto "github.com/ca-risken/core/proto/project"
)

type deleteMembershipBeforeOrgRecheck struct {
	db.AlertRepository
	client         *db.Client
	organizationID uint32
	projectID      uint32
	once           sync.Once
	err            error
}

func (r *deleteMembershipBeforeOrgRecheck) GetOrgAlertNotificationTarget(ctx context.Context, organizationID, projectID, alertConditionID, notificationID uint32) (*db.OrgAlertNotificationTarget, error) {
	if organizationID == r.organizationID && projectID == r.projectID {
		r.once.Do(func() {
			r.err = r.client.RemoveProjectsInOrganization(ctx, r.organizationID, r.projectID)
		})
		if r.err != nil {
			return nil, r.err
		}
	}
	return r.AlertRepository.GetOrgAlertNotificationTarget(ctx, organizationID, projectID, alertConditionID, notificationID)
}

func TestOrgNotificationIntegration(t *testing.T) {
	client := newIntegrationDBClient(t)
	ctx := context.Background()

	t.Run("transactional lifecycle", func(t *testing.T) {
		project, err := client.CreateProject(ctx, "t010-lifecycle-project")
		if err != nil {
			t.Fatal(err)
		}
		organization, err := client.CreateOrganization(ctx, "t010-lifecycle-org", "integration")
		if err != nil {
			t.Fatal(err)
		}
		condition, err := client.UpsertAlertCondition(ctx, &model.AlertCondition{ProjectID: project.ProjectID, Description: "lifecycle", Severity: "low", AndOr: "and", Enabled: true})
		if err != nil {
			t.Fatal(err)
		}
		notification, err := client.UpsertOrgNotification(ctx, &model.OrganizationNotification{OrganizationID: organization.OrganizationID, Name: "lifecycle", Type: "slack", NotifySetting: `{"webhook_url":"http://127.0.0.1/unused"}`})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := client.PutOrganizationProject(ctx, organization.OrganizationID, project.ProjectID); err != nil {
			t.Fatal(err)
		}
		assertOrgRelation(t, client, organization.OrganizationID, project.ProjectID, condition.AlertConditionID, notification.NotificationID, 2592000, 0, true)

		preservedAt := time.Now().Add(-time.Minute).Truncate(time.Second)
		if err := client.Master.WithContext(ctx).Exec(`update organization_alert_cond_notification set cache_second = ?, notified_at = ? where organization_id = ? and project_id = ? and alert_condition_id = ? and notification_id = ?`, 73, preservedAt, organization.OrganizationID, project.ProjectID, condition.AlertConditionID, notification.NotificationID).Error; err != nil {
			t.Fatal(err)
		}
		if _, err := client.PutOrganizationProject(ctx, organization.OrganizationID, project.ProjectID); err != nil {
			t.Fatal(err)
		}
		assertOrgRelation(t, client, organization.OrganizationID, project.ProjectID, condition.AlertConditionID, notification.NotificationID, 73, preservedAt.Unix(), true)

		if err := client.RemoveProjectsInOrganization(ctx, organization.OrganizationID, project.ProjectID); err != nil {
			t.Fatal(err)
		}
		assertOrgRelation(t, client, organization.OrganizationID, project.ProjectID, condition.AlertConditionID, notification.NotificationID, 0, 0, false)
		if _, err := client.PutOrganizationProject(ctx, organization.OrganizationID, project.ProjectID); err != nil {
			t.Fatal(err)
		}
		assertOrgRelation(t, client, organization.OrganizationID, project.ProjectID, condition.AlertConditionID, notification.NotificationID, 2592000, 0, true)

		if err := client.DeleteAlertCondition(ctx, project.ProjectID, condition.AlertConditionID); err != nil {
			t.Fatal(err)
		}
		assertOrgRelation(t, client, organization.OrganizationID, project.ProjectID, condition.AlertConditionID, notification.NotificationID, 0, 0, false)
		if _, err := client.UpsertAlertCondition(ctx, condition); err != nil {
			t.Fatal(err)
		}
		assertOrgRelation(t, client, organization.OrganizationID, project.ProjectID, condition.AlertConditionID, notification.NotificationID, 2592000, 0, true)

		if err := client.DeleteOrgNotification(ctx, organization.OrganizationID, notification.NotificationID); err != nil {
			t.Fatal(err)
		}
		assertOrgRelation(t, client, organization.OrganizationID, project.ProjectID, condition.AlertConditionID, notification.NotificationID, 0, 0, false)
		if _, err := client.UpsertOrgNotification(ctx, notification); err != nil {
			t.Fatal(err)
		}
		assertOrgRelation(t, client, organization.OrganizationID, project.ProjectID, condition.AlertConditionID, notification.NotificationID, 2592000, 0, true)

		deleteCases := []struct {
			name   string
			remove func() error
		}{
			{name: "organization", remove: func() error { return client.DeleteOrganization(ctx, organization.OrganizationID) }},
		}
		for _, tc := range deleteCases {
			t.Run("owner delete "+tc.name, func(t *testing.T) {
				if err := tc.remove(); err != nil {
					t.Fatal(err)
				}
				assertOrgRelation(t, client, organization.OrganizationID, project.ProjectID, condition.AlertConditionID, notification.NotificationID, 0, 0, false)
			})
		}

		project2, err := client.CreateProject(ctx, "t010-deleted-project")
		if err != nil {
			t.Fatal(err)
		}
		organization2, err := client.CreateOrganization(ctx, "t010-project-owner-org", "integration")
		if err != nil {
			t.Fatal(err)
		}
		condition2, err := client.UpsertAlertCondition(ctx, &model.AlertCondition{ProjectID: project2.ProjectID, Description: "owner delete", Severity: "low", AndOr: "and", Enabled: true})
		if err != nil {
			t.Fatal(err)
		}
		notification2, err := client.UpsertOrgNotification(ctx, &model.OrganizationNotification{OrganizationID: organization2.OrganizationID, Name: "owner delete", Type: "slack", NotifySetting: `{"webhook_url":"http://127.0.0.1/unused"}`})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := client.PutOrganizationProject(ctx, organization2.OrganizationID, project2.ProjectID); err != nil {
			t.Fatal(err)
		}
		if err := client.DeleteProject(ctx, project2.ProjectID); err != nil {
			t.Fatal(err)
		}
		assertOrgRelation(t, client, organization2.OrganizationID, project2.ProjectID, condition2.AlertConditionID, notification2.NotificationID, 0, 0, false)
	})

	t.Run("independent project and organization webhooks", func(t *testing.T) {
		var projectCalls atomic.Int32
		var organizationCalls atomic.Int32
		projectWebhook := webhookServer(t, &projectCalls)
		defer projectWebhook.Close()
		organizationWebhook := webhookServer(t, &organizationCalls)
		defer organizationWebhook.Close()

		project, err := client.CreateProject(ctx, "t010-notification-project")
		if err != nil {
			t.Fatal(err)
		}
		organization, err := client.CreateOrganization(ctx, "t010-notification-org", "integration")
		if err != nil {
			t.Fatal(err)
		}
		condition, err := client.UpsertAlertCondition(ctx, &model.AlertCondition{ProjectID: project.ProjectID, Description: "notify", Severity: "low", AndOr: "and", Enabled: true})
		if err != nil {
			t.Fatal(err)
		}
		orgNotification, err := client.UpsertOrgNotification(ctx, &model.OrganizationNotification{OrganizationID: organization.OrganizationID, Name: "org", Type: "slack", NotifySetting: fmt.Sprintf(`{"webhook_url":%q}`, organizationWebhook.URL)})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := client.PutOrganizationProject(ctx, organization.OrganizationID, project.ProjectID); err != nil {
			t.Fatal(err)
		}
		projectNotification, err := client.UpsertNotification(ctx, &model.Notification{ProjectID: project.ProjectID, Name: "project", Type: "slack", NotifySetting: fmt.Sprintf(`{"webhook_url":%q}`, projectWebhook.URL)})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := client.UpsertAlertCondNotification(ctx, &model.AlertCondNotification{ProjectID: project.ProjectID, AlertConditionID: condition.AlertConditionID, NotificationID: projectNotification.NotificationID, CacheSecond: 3600, NotifiedAt: time.Unix(0, 0)}); err != nil {
			t.Fatal(err)
		}
		t.Run("DATETIME(0) rounding boundary", func(t *testing.T) {
			upperBound := time.Now().Truncate(time.Second)
			roundedUpdatedAt := upperBound.Add(time.Second)
			if err := client.Master.WithContext(ctx).Exec(`update alert_cond_notification set updated_at = ? where project_id = ? and alert_condition_id = ? and notification_id = ?`, roundedUpdatedAt, project.ProjectID, condition.AlertConditionID, projectNotification.NotificationID).Error; err != nil {
				t.Fatal(err)
			}
			publicRelations, err := client.ListAlertCondNotification(ctx, project.ProjectID, condition.AlertConditionID, projectNotification.NotificationID, 0, upperBound.Unix())
			if err != nil {
				t.Fatal(err)
			}
			if len(*publicRelations) != 0 {
				t.Fatalf("period search unexpectedly included rounded future row: %+v", *publicRelations)
			}
			currentRelations, err := client.ListAlertCondNotificationForNotification(ctx, project.ProjectID, condition.AlertConditionID)
			if err != nil {
				t.Fatal(err)
			}
			if len(*currentRelations) != 1 || (*currentRelations)[0].NotificationID != projectNotification.NotificationID {
				t.Fatalf("notification lookup missed rounded row: %+v", *currentRelations)
			}
		})
		if _, err := client.UpdateOrgAlertCondNotificationCache(ctx, organization.OrganizationID, project.ProjectID, condition.AlertConditionID, orgNotification.NotificationID, 3600); err != nil {
			t.Fatal(err)
		}

		service := AlertService{repository: client, logger: logging.NewLogger(), baseURL: "http://risken.example", defaultLocale: LocaleEn}
		alertData := &model.Alert{ProjectID: project.ProjectID, AlertConditionID: condition.AlertConditionID, Description: "integration", Severity: "low", Status: "open"}
		projectData := &projectproto.Project{ProjectId: project.ProjectID, Name: project.Name}
		rules := &[]model.AlertRule{}
		findings := &[]uint64{}
		cases := []struct {
			name                  string
			resetProject          bool
			resetOrganization     bool
			wantProjectCalls      int32
			wantOrganizationCalls int32
		}{
			{name: "both", wantProjectCalls: 1, wantOrganizationCalls: 1},
			{name: "neither suppressed", wantProjectCalls: 1, wantOrganizationCalls: 1},
			{name: "project only", resetProject: true, wantProjectCalls: 2, wantOrganizationCalls: 1},
			{name: "organization only", resetOrganization: true, wantProjectCalls: 2, wantOrganizationCalls: 2},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if tc.resetProject {
					setProjectNotifiedAt(t, client, project.ProjectID, condition.AlertConditionID, projectNotification.NotificationID, time.Unix(0, 0))
				}
				if tc.resetOrganization {
					setOrgNotifiedAt(t, client, organization.OrganizationID, project.ProjectID, condition.AlertConditionID, orgNotification.NotificationID, time.Unix(0, 0))
				}
				if err := service.NotificationAlert(ctx, condition, alertData, rules, projectData, findings, false); err != nil {
					t.Fatal(err)
				}
				if got := projectCalls.Load(); got != tc.wantProjectCalls {
					t.Fatalf("project webhook calls: want=%d got=%d", tc.wantProjectCalls, got)
				}
				if got := organizationCalls.Load(); got != tc.wantOrganizationCalls {
					t.Fatalf("organization webhook calls: want=%d got=%d", tc.wantOrganizationCalls, got)
				}
			})
		}
		assertNotifiedAfterEpoch(t, client, "alert_cond_notification", "project_id = ? and alert_condition_id = ? and notification_id = ?", project.ProjectID, condition.AlertConditionID, projectNotification.NotificationID)
		assertNotifiedAfterEpoch(t, client, "organization_alert_cond_notification", "organization_id = ? and project_id = ? and alert_condition_id = ? and notification_id = ?", organization.OrganizationID, project.ProjectID, condition.AlertConditionID, orgNotification.NotificationID)

		deletedOrganization, err := client.CreateOrganization(ctx, "t010-send-time-delete-org", "integration")
		if err != nil {
			t.Fatal(err)
		}
		deletedNotification, err := client.UpsertOrgNotification(ctx, &model.OrganizationNotification{OrganizationID: deletedOrganization.OrganizationID, Name: "deleted", Type: "slack", NotifySetting: fmt.Sprintf(`{"webhook_url":%q}`, organizationWebhook.URL)})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := client.PutOrganizationProject(ctx, deletedOrganization.OrganizationID, project.ProjectID); err != nil {
			t.Fatal(err)
		}
		setProjectNotifiedAt(t, client, project.ProjectID, condition.AlertConditionID, projectNotification.NotificationID, time.Now())
		setOrgNotifiedAt(t, client, organization.OrganizationID, project.ProjectID, condition.AlertConditionID, orgNotification.NotificationID, time.Now())
		wrapped := &deleteMembershipBeforeOrgRecheck{AlertRepository: client, client: client, organizationID: deletedOrganization.OrganizationID, projectID: project.ProjectID}
		service.repository = wrapped
		before := organizationCalls.Load()
		if err := service.NotificationAlert(ctx, condition, alertData, rules, projectData, findings, false); err != nil {
			t.Fatal(err)
		}
		if got := organizationCalls.Load(); got != before {
			t.Fatalf("send-time deleted target was notified: before=%d after=%d", before, got)
		}
		assertOrgRelation(t, client, deletedOrganization.OrganizationID, project.ProjectID, condition.AlertConditionID, deletedNotification.NotificationID, 0, 0, false)
		t.Logf("webhook evidence: project=%d organization=%d", projectCalls.Load(), organizationCalls.Load())
	})
}

func newIntegrationDBClient(t *testing.T) *db.Client {
	t.Helper()
	host := os.Getenv("T010_DB_HOST")
	database := os.Getenv("T010_DB_NAME")
	if host == "" || database == "" {
		t.Skip("T010_DB_HOST and T010_DB_NAME are required")
	}
	port, err := strconv.Atoi(envOrDefault("T010_DB_PORT", "3306"))
	if err != nil {
		t.Fatal(err)
	}
	client, err := db.NewClient(&db.Config{
		MasterHost: host, SlaveHost: host,
		MasterUser: envOrDefault("T010_DB_USER", "root"), SlaveUser: envOrDefault("T010_DB_USER", "root"),
		MasterPassword: os.Getenv("T010_DB_PASSWORD"), SlavePassword: os.Getenv("T010_DB_PASSWORD"),
		Schema: database, Port: port, MaxConnection: 4,
	}, logging.NewLogger())
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func webhookServer(t *testing.T, count *atomic.Int32) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		count.Add(1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
}

func assertOrgRelation(t *testing.T, client *db.Client, organizationID, projectID, alertConditionID, notificationID, wantCache uint32, wantNotifiedAt int64, wantExists bool) {
	t.Helper()
	relation, err := client.GetOrgAlertCondNotification(context.Background(), organizationID, projectID, alertConditionID, notificationID)
	if !wantExists {
		if err == nil {
			t.Fatalf("unexpected organization relation: %+v", relation)
		}
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	if relation.CacheSecond != wantCache || relation.NotifiedAt != wantNotifiedAt {
		t.Fatalf("organization relation state: want cache=%d notified_at=%d, got=%+v", wantCache, wantNotifiedAt, relation)
	}
}

func setProjectNotifiedAt(t *testing.T, client *db.Client, projectID, alertConditionID, notificationID uint32, notifiedAt time.Time) {
	t.Helper()
	if err := client.Master.Exec(`update alert_cond_notification set notified_at = ? where project_id = ? and alert_condition_id = ? and notification_id = ?`, notifiedAt, projectID, alertConditionID, notificationID).Error; err != nil {
		t.Fatal(err)
	}
}

func setOrgNotifiedAt(t *testing.T, client *db.Client, organizationID, projectID, alertConditionID, notificationID uint32, notifiedAt time.Time) {
	t.Helper()
	if err := client.Master.Exec(`update organization_alert_cond_notification set notified_at = ? where organization_id = ? and project_id = ? and alert_condition_id = ? and notification_id = ?`, notifiedAt, organizationID, projectID, alertConditionID, notificationID).Error; err != nil {
		t.Fatal(err)
	}
}

func assertNotifiedAfterEpoch(t *testing.T, client *db.Client, table, predicate string, args ...interface{}) {
	t.Helper()
	var notifiedAt time.Time
	if err := client.Master.Raw("select notified_at from "+table+" where "+predicate, args...).Scan(&notifiedAt).Error; err != nil {
		t.Fatal(err)
	}
	if !notifiedAt.After(time.Unix(0, 0)) {
		t.Fatalf("%s notified_at was not updated: %s", table, notifiedAt)
	}
}
