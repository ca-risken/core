package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/alert"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"github.com/vikyd/zero"
)

/**
 * AnalyzeAlert
 */

func (f *alertService) AnalyzeAlert(ctx context.Context, req *alert.AnalyzeAlertRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	// alertConditionの取得
	// 存在しなければ終了
	alertConditions, err := f.repository.ListAlertCondition(req.ProjectId, nil, true, 0, time.Now().Unix())
	noRecordAlertCondition := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecordAlertCondition {
		appLogger.Error(err)
		return nil, err
	}

	// findingの取得
	// 存在しなければ終了
	findings, err := f.getFindings(ctx, req.ProjectId)
	noRecordFinding := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecordFinding {
		appLogger.Error(err)
		return nil, err
	}

	// マッチング

	for _, alertCondition := range *alertConditions {
		err := f.AnalyzeAlertByCondition(ctx, &alertCondition, findings)
		if err != nil {
			appLogger.Error(err)
			return &empty.Empty{}, err
		}
	}
	return &empty.Empty{}, nil
}

func (f *alertService) AnalyzeAlertByCondition(ctx context.Context, alertCondition *model.AlertCondition, findings []*finding.Finding) error {
	isMatch := false

	// AlertRuleの取得
	// 存在しなければ終了
	alertRules, err := f.repository.ListAlertRuleByAlertConditionID(alertCondition.ProjectID, alertCondition.AlertConditionID)
	if err != nil {
		appLogger.Error(err)
		return err
	}
	var matchFindingIDs []uint64
	for _, alertRule := range *alertRules {
		appLogger.Info("alertRule:", alertRule)
		isMatchRule, matchFindingIDsByAlert := analyzeAlertByRule(&alertRule, findings)
		appLogger.Info("isMatchRule:", isMatchRule)
		if isMatchRule {
			isMatch = true
			matchFindingIDs = append(matchFindingIDs, matchFindingIDsByAlert...)
		} else {
			if alertCondition.AndOr == "and" {
				isMatch = false
				break
			}
		}
	}
	if isMatch {
		appLogger.Info("isMatch:", isMatch)

		alertID, err := f.RegistAlertByAnalyze(alertCondition, matchFindingIDs)
		if err != nil {
			appLogger.Error(err)
			return err
		}
		// Matchしている場合はAlert通知を行う
		errNoti := f.NotificationAlert(alertCondition, alertID)
		if errNoti != nil {
			return errNoti
		}
	} else {
		appLogger.Info("isMatch:", isMatch)
		err := f.DeleteAlertByAnalyze(alertCondition)
		if err != nil {
			appLogger.Error(err)
			return err
		}
	}
	return nil
}

func (f *alertService) RegistAlertByAnalyze(alertCondition *model.AlertCondition, findingIDs []uint64) (uint32, error) {

	// Alertの登録
	savedData, err := f.repository.GetAlertByAlertConditionIDWithActivated(alertCondition.ProjectID, alertCondition.AlertConditionID, true)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return 0, err
	}

	// 既に登録済みの場合はalertIDを取得
	var alertID uint32
	if !noRecord {
		alertID = savedData.AlertID
	}

	data := &model.Alert{
		AlertID:          alertID,
		AlertConditionID: alertCondition.AlertConditionID,
		Description:      alertCondition.Description,
		Severity:         alertCondition.Severity,
		ProjectID:        alertCondition.ProjectID,
		Activated:        true,
	}
	// insert Alert
	registerdData, err := f.repository.UpsertAlert(data)
	if err != nil {
		return 0, err
	}

	historyType := getHistoryType(alertID)
	// AlertHistoryの登録
	dataAlertHistory := &model.AlertHistory{
		HistoryType: historyType,
		AlertID:     registerdData.AlertID,
		Description: registerdData.Description,
		Severity:    registerdData.Severity,
		ProjectID:   registerdData.ProjectID,
	}
	_, errAlertHistory := f.repository.UpsertAlertHistory(dataAlertHistory)
	if errAlertHistory != nil {
		return 0, errAlertHistory
	}

	//RelAlertFindingの登録
	appLogger.Info(registerdData)

	for _, findingID := range findingIDs {
		data := &model.RelAlertFinding{
			AlertID:   registerdData.AlertID,
			FindingID: uint32(findingID),
			ProjectID: registerdData.ProjectID,
		}
		_, err := f.repository.UpsertRelAlertFinding(data)
		if err != nil {
			return 0, err
		}
	}

	return registerdData.AlertID, nil
}

func (f *alertService) DeleteAlertByAnalyze(alertCondition *model.AlertCondition) error {

	// Alertの削除
	savedData, err := f.repository.GetAlertByAlertConditionIDWithActivated(alertCondition.ProjectID, alertCondition.AlertConditionID, true)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return err
	}

	// レコードが存在しない、もしくはActivated出ない場合は何もしない
	appLogger.Info("savedData:", savedData)
	if noRecord || !savedData.Activated {
		return nil
	}

	data := &model.Alert{
		AlertID:          savedData.AlertID,
		AlertConditionID: alertCondition.AlertConditionID,
		Description:      alertCondition.Description,
		Severity:         alertCondition.Severity,
		ProjectID:        alertCondition.ProjectID,
		Activated:        false,
	}
	// update Alert
	errDeactivate := f.repository.DeactivateAlert(data)
	if errDeactivate != nil {
		return errDeactivate
	}

	// AlertHistoryの登録
	dataAlertHistory := &model.AlertHistory{
		HistoryType: "deleted",
		AlertID:     savedData.AlertID,
		Description: savedData.Description,
		Severity:    savedData.Severity,
		ProjectID:   savedData.ProjectID,
	}
	_, errAlertHistory := f.repository.UpsertAlertHistory(dataAlertHistory)
	if errAlertHistory != nil {
		return errAlertHistory
	}

	//RelAlertFindingの削除
	listRelAlertFinding, err := f.repository.ListRelAlertFinding(alertCondition.ProjectID, savedData.AlertID, uint32(0), 0, time.Now().Unix())
	if err != nil {
		return err
	}
	for _, relAlertFinding := range *listRelAlertFinding {
		err := f.repository.DeleteRelAlertFinding(relAlertFinding.ProjectID, relAlertFinding.AlertID, relAlertFinding.FindingID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *alertService) getFindings(ctx context.Context, projectID uint32) ([]*finding.Finding, error) {
	req := &finding.ListFindingRequest{ProjectId: projectID}
	list, err := f.findingClient.ListFinding(ctx, req)
	if err != nil {
		appLogger.Error(err)
		return nil, err
	}
	var ret []*finding.Finding
	findingIDs := list.FindingId
	for _, findingID := range findingIDs {
		req := &finding.GetFindingRequest{ProjectId: projectID, FindingId: findingID}
		get, err := f.findingClient.GetFinding(ctx, req)
		if err != nil {
			appLogger.Error(err)
			return nil, err
		}
		ret = append(ret, get.Finding)
	}
	return ret, nil
}

func (f *alertService) NotificationAlert(alertCondition *model.AlertCondition, alertID uint32) error {

	alertCondNotifications, err := f.repository.ListAlertCondNotification(alertCondition.ProjectID, alertCondition.AlertConditionID, 0, 0, time.Now().Unix())
	if err != nil {
		return err
	}
	for _, alertCondNotification := range *alertCondNotifications {
		// 連続通知を防ぐ
		if time.Now().Unix() < alertCondNotification.NotifiedAt.Unix()+int64(alertCondNotification.CacheSecond) {
			continue
		}

		notification, err := f.repository.GetNotification(alertCondition.ProjectID, alertCondNotification.NotificationID)
		if err != nil {
			return err
		}
		switch notification.Type {
		case "slack":
			errNoti := sendSlackNotification(notification.NotifySetting, alertID)
			if errNoti != nil {
				return errNoti
			}
		default:
			appLogger.Warn("This notification_type is unimprement.", notification.Type)
			break
		}
	}
	return nil
}

func analyzeAlertByRule(alertRule *model.AlertRule, findings []*finding.Finding) (bool, []uint64) {
	var matchFindingIDs []uint64
	for _, f := range findings {
		if checkMatchAlertRuleFinding(alertRule, f) {
			matchFindingIDs = append(matchFindingIDs, f.FindingId)
		}
	}
	isMatch := (len(matchFindingIDs) >= int(alertRule.FindingCnt))
	return isMatch, matchFindingIDs
}

func checkMatchAlertRuleFinding(alertRule *model.AlertRule, finding *finding.Finding) bool {
	if alertRule.Score > finding.Score {
		return false
	}
	if !zero.IsZeroVal(alertRule.ResourceName) && strings.Index(finding.ResourceName, alertRule.ResourceName) == -1 {
		return false
	}
	//	if !zero.IsZeroVal(alertRule.Tag) and strings.Index(finding.ResourceName, alertRule.ResourceName) == -1 {
	//		return false
	//	}
	return true
}

func getHistoryType(alertID uint32) string {
	if !zero.IsZeroVal(alertID) {
		return "created"
	}
	return "updated"
}

func sendSlackNotification(notifySetting string, alertID uint32) error {
	var setting slackNotifySetting
	if err := json.Unmarshal([]byte(notifySetting), &setting); err != nil {
		panic(err)
	}
	if zero.IsZeroVal(setting.WebhookURL) {
		appLogger.Warn("Unset webhook_url")
		return nil
	}
	var channel string
	if !zero.IsZeroVal(setting.Data["channel"]) {
		channel = setting.Data["channel"]
	}
	slackAlert := slackWebhookSetting{Channel: channel, AlertID: alertID}
	payload, err := slackAlert.GetPayload()
	if err != nil {
		return err
	}
	resp, err := http.PostForm(setting.WebhookURL, url.Values{"payload": {string(payload)}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

type slackNotifySetting struct {
	WebhookURL string            `json:"webhook_url"`
	Data       map[string]string `json:"data"`
}
