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
	// 有効なalertConditionの取得
	alertConditions, err := f.repository.ListAlertCondition(req.ProjectId, nil, true, 0, time.Now().Unix())
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		appLogger.Error(err)
		return nil, err
	}

	// findingの取得
	findings, err := f.repository.ListFinding(req.ProjectId)
	noRecord = gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		appLogger.Error(err)
		return nil, err
	}

	// マッチング
	for _, alertCondition := range *alertConditions {
		err := f.AnalyzeAlertByCondition(ctx, &alertCondition, findings)
		if err != nil {
			appLogger.Error(err)
			return nil, err
		}
	}

	// 無効のalertConditionの取得
	disabledAlertConditions, err := f.repository.ListDisabledAlertCondition(req.ProjectId)
	noRecord = gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		appLogger.Error(err)
		return nil, err
	}
	// 無効なalertConditionに紐づくAlertを削除
	for _, alertCondition := range *disabledAlertConditions {
		err := f.DeleteAlertByAnalyze(&alertCondition)
		if err != nil {
			appLogger.Error(err)
			return nil, err
		}
	}
	return &empty.Empty{}, nil
}

func (f *alertService) AnalyzeAlertByCondition(ctx context.Context, alertCondition *model.AlertCondition, findings *[]model.Finding) error {
	isMatch := false

	// AlertRuleの取得
	alertRules, err := f.repository.ListAlertRuleByAlertConditionID(alertCondition.ProjectID, alertCondition.AlertConditionID)
	if err != nil {
		return err
	}
	var matchFindingIDs []uint64
	for _, alertRule := range *alertRules {
		isMatchRule, matchFindingIDsByAlert, err := f.analyzeAlertByRule(ctx, &alertRule, findings)
		if err != nil {
			return err
		}
		if isMatchRule {
			isMatch = true
			matchFindingIDs = append(matchFindingIDs, *matchFindingIDsByAlert...)
		} else {
			if alertCondition.AndOr == "and" {
				isMatch = false
				break
			}
		}
	}
	if isMatch {
		alertID, err := f.RegistAlertByAnalyze(alertCondition, matchFindingIDs)
		if err != nil {
			return err
		}
		// Matchしている場合はAlert通知を行う
		err = f.NotificationAlert(alertCondition, alertID)
		if err != nil {
			return err
		}
	} else {
		err = f.DeleteAlertByAnalyze(alertCondition)
		if err != nil {
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
	_, err = f.repository.UpsertAlertHistory(dataAlertHistory)
	if err != nil {
		return 0, err
	}

	//RelAlertFindingの更新 (削除して再登録)
	err = f.deleteRelAlertFindingByAlertID(registerdData.ProjectID, registerdData.AlertID)
	if err != nil {
		return 0, err
	}
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
	err = f.deleteRelAlertFindingByAlertID(savedData.ProjectID, savedData.AlertID)
	if err != nil {
		return err
	}

	return nil
}

func (f *alertService) getFindingTags(ctx context.Context, projectID uint32, findingID uint64) (*[]model.FindingTag, error) {
	list, err := f.repository.ListFindingTag(projectID, findingID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (f *alertService) deleteRelAlertFindingByAlertID(projectID, alertID uint32) error {
	listRelAlertFinding, err := f.repository.ListRelAlertFinding(projectID, alertID, uint32(0), 0, time.Now().Unix())
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
			err = sendSlackNotification(notification.NotifySetting, alertID)
			if err != nil {
				return err
			}
		default:
			appLogger.Warn("This notification_type is unimprement.", notification.Type)
			break
		}
		// 通知時刻を更新する
		dataAlertCondNotification := &model.AlertCondNotification{
			AlertConditionID: alertCondNotification.AlertConditionID,
			NotificationID:   alertCondNotification.NotificationID,
			CacheSecond:      alertCondNotification.CacheSecond,
			NotifiedAt:       time.Now(),
			ProjectID:        alertCondNotification.ProjectID,
		}
		_, err = f.repository.UpsertAlertCondNotification(dataAlertCondNotification)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *alertService) analyzeAlertByRule(ctx context.Context, alertRule *model.AlertRule, findings *[]model.Finding) (bool, *[]uint64, error) {
	matchFindingIDs := []uint64{}
	for _, finding := range *findings {
		isMatch, err := f.checkMatchAlertRuleFinding(ctx, alertRule, &finding)
		if err != nil {
			return false, &matchFindingIDs, err
		}
		if isMatch {
			matchFindingIDs = append(matchFindingIDs, finding.FindingID)
		}
	}
	isMatch := (len(matchFindingIDs) >= int(alertRule.FindingCnt))
	return isMatch, &matchFindingIDs, nil
}

func (f *alertService) checkMatchAlertRuleFinding(ctx context.Context, alertRule *model.AlertRule, finding *model.Finding) (bool, error) {
	if alertRule.Score > finding.Score {
		return false, nil
	}
	if !zero.IsZeroVal(alertRule.ResourceName) && strings.Index(finding.ResourceName, alertRule.ResourceName) == -1 {
		return false, nil
	}
	if !zero.IsZeroVal(alertRule.Tag) {
		//findingIDがマッチするTagの収集
		findingTags, err := f.getFindingTags(ctx, alertRule.ProjectID, finding.FindingID)
		if err != nil {
			return false, err
		}
		isMatchTag := false
		for _, findingTag := range *findingTags {
			if findingTag.Tag == alertRule.Tag {
				isMatchTag = true
				break
			}
		}
		// findingとalertRuleのTagがマッチしなければfalseを返す
		if !isMatchTag {
			return false, nil
		}
	}
	return true, nil
}

func getHistoryType(alertID uint32) string {
	if zero.IsZeroVal(alertID) {
		return "created"
	}
	return "updated"
}

func sendSlackNotification(notifySetting string, alertID uint32) error {
	var setting slackNotifySetting
	if err := json.Unmarshal([]byte(notifySetting), &setting); err != nil {
		return err
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
