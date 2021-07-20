package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/alert"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/CyberAgent/mimosa-core/proto/project"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

/**
 * AnalyzeAlert
 */

func (a *alertService) AnalyzeAlert(ctx context.Context, req *alert.AnalyzeAlertRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	requestID := getRequestID(req.ProjectId)
	// 有効なalertConditionの取得
	appLogger.Infof("start ListEnabledAlertCondition: RequestID=%s", requestID)
	_, ss := xray.BeginSubsegment(ctx, "ListEnabledAlertCondition")
	alertConditions, err := a.repository.ListEnabledAlertCondition(ctx, req.ProjectId, req.AlertConditionId)
	ss.Close(err)
	appLogger.Infof("finish ListEnabledAlertCondition: RequestID=%s", requestID)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		appLogger.Error(err)
		return nil, err
	}

	// マッチング
	appLogger.Infof("start matching: RequestID=%s", requestID)
	for _, alertCondition := range *alertConditions {
		err := a.AnalyzeAlertByCondition(ctx, &alertCondition)
		if err != nil {
			appLogger.Error(err)
			return nil, err
		}
	}
	appLogger.Infof("finish matching: RequestID=%s", requestID)

	// 無効のalertConditionの取得
	appLogger.Infof("start ListDisabledAlertCondition: RequestID=%s", requestID)
	_, ss = xray.BeginSubsegment(ctx, "ListDisabledAlertCondition")
	disabledAlertConditions, err := a.repository.ListDisabledAlertCondition(ctx, req.ProjectId, req.AlertConditionId)
	ss.Close(err)
	appLogger.Infof("finish ListDisabledAlertCondition: RequestID=%s", requestID)
	noRecord = errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		appLogger.Error(err)
		return nil, err
	}
	// 無効なalertConditionに紐づくAlertを削除
	appLogger.Infof("start DeleteAlertByAnalyze: RequestID=%s", requestID)
	for _, alertCondition := range *disabledAlertConditions {
		err := a.DeleteAlertByAnalyze(ctx, &alertCondition)
		if err != nil {
			appLogger.Error(err)
			return nil, err
		}
	}
	appLogger.Infof("finish DeleteAlertByAnalyze: RequestID=%s", requestID)
	return &empty.Empty{}, nil
}

func (a *alertService) AnalyzeAlertByCondition(ctx context.Context, alertCondition *model.AlertCondition) error {
	// AlertRuleの取得
	_, ss := xray.BeginSubsegment(ctx, "ListAlertRuleByAlertConditionID")
	alertRules, err := a.repository.ListAlertRuleByAlertConditionID(ctx, alertCondition.ProjectID, alertCondition.AlertConditionID)
	ss.Close(err)
	if err != nil {
		appLogger.Errorf("Failed list AlertRule by AlertConditionID. alertConditionID: %v, err: %v", alertCondition.AlertConditionID, err)
		return err
	}
	var matchFindingIDs []uint64
	isFirst := true
	appLogger.Info("start matching per rule")
	for _, alertRule := range *alertRules {
		isMatchRule, matchFindingIDsByAlert, err := a.analyzeAlertByRule(ctx, &alertRule)
		if err != nil {
			return err
		}
		if !isMatchRule {
			if alertCondition.AndOr == "and" {
				matchFindingIDs = []uint64{} // clear
				break
			} else {
				continue
			}
		}
		if alertCondition.AndOr == "or" || isFirst {
			matchFindingIDs = append(matchFindingIDs, *matchFindingIDsByAlert...)
			isFirst = false
			continue
		}
		var andMatchFindingIDs []uint64
		for _, id := range *matchFindingIDsByAlert {
			if isContainsFindings(id, matchFindingIDs) {
				andMatchFindingIDs = append(andMatchFindingIDs, id)
			}
		}
		if len(andMatchFindingIDs) == 0 {
			matchFindingIDs = []uint64{} // clear
			break
		}
		matchFindingIDs = andMatchFindingIDs
	}
	appLogger.Info("finish matching per rule")
	if len(matchFindingIDs) > 0 {
		registAlert, err := a.RegistAlertByAnalyze(ctx, alertCondition, matchFindingIDs)
		if err != nil {
			return err
		}
		// AlertがACTIVE、かつMatchしている場合はAlert通知を行う
		if registAlert.Status == alert.Status_ACTIVE.String() {
			err = a.NotificationAlert(ctx, alertCondition, registAlert, alertRules)
			if err != nil {
				return err
			}
		}
	} else {
		err = a.DeleteAlertByAnalyze(ctx, alertCondition)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *alertService) RegistAlertByAnalyze(ctx context.Context, alertCondition *model.AlertCondition, findingIDs []uint64) (*model.Alert, error) {
	// AlertConditionに該当するAlertが既に存在しているか確認
	_, ss := xray.BeginSubsegment(ctx, "GetAlertByAlertConditionIDStatus")
	savedData, err := a.repository.GetAlertByAlertConditionIDStatus(ctx, alertCondition.ProjectID, alertCondition.AlertConditionID, []string{"ACTIVE", "PENDING"})
	ss.Close(err)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		return nil, err
	}

	// 既に登録済みの場合は登録済みのAlertと内容が一致しているか確認
	// 一致している場合処理を終了する
	// 登録済みかつStatusがPENDINGの場合、PENDING
	var status string
	var alertID uint32
	var isMatchExisting bool
	if !noRecord {
		_, ss = xray.BeginSubsegment(ctx, "isMatchExistingAlert")
		isMatchExisting, err = a.isMatchExistingAlert(ctx, savedData, alertCondition, findingIDs)
		ss.Close(err)
		if err != nil {
			return nil, err
		}

		alertID = savedData.AlertID
		status = savedData.Status
	} else {
		status = "ACTIVE"
	}

	data := &model.Alert{
		AlertID:          alertID,
		AlertConditionID: alertCondition.AlertConditionID,
		Description:      alertCondition.Description,
		Severity:         alertCondition.Severity,
		ProjectID:        alertCondition.ProjectID,
		Status:           status,
	}
	// upsert Alert
	_, ss = xray.BeginSubsegment(ctx, "UpsertAlert")
	registerdData, err := a.repository.UpsertAlert(ctx, data)
	ss.Close(err)
	if err != nil {
		appLogger.Errorf("Error occured when upsert alert. alertConditionID: %v, err: %v", alertCondition.AlertConditionID, err)
		return nil, err
	}

	// 過去のアラートと状態が同じなら以下の処理はスキップ
	if isMatchExisting {
		return registerdData, nil
	}

	// AlertHistoryに登録するための現在のRelAlertFindingを整形
	findingHistory, err := makeFindingIDs(findingIDs)
	if err != nil {
		return nil, err
	}

	historyType := getHistoryType(alertID)
	// AlertHistoryの登録
	dataAlertHistory := &model.AlertHistory{
		HistoryType:    historyType,
		AlertID:        registerdData.AlertID,
		Description:    registerdData.Description,
		FindingHistory: findingHistory,
		Severity:       registerdData.Severity,
		ProjectID:      registerdData.ProjectID,
	}
	_, ss = xray.BeginSubsegment(ctx, "UpsertAlertHistory")
	_, err = a.repository.UpsertAlertHistory(ctx, dataAlertHistory)
	ss.Close(err)
	if err != nil {
		appLogger.Errorf("Error occured when upsert AlertHistory. err: %v", err)
		return nil, err
	}

	//RelAlertFindingの更新 (削除して再登録)
	err = a.deleteRelAlertFindingByAlertID(ctx, registerdData.ProjectID, registerdData.AlertID)
	if err != nil {
		return nil, err
	}
	for _, findingID := range findingIDs {
		data := &model.RelAlertFinding{
			AlertID:   registerdData.AlertID,
			FindingID: uint32(findingID),
			ProjectID: registerdData.ProjectID,
		}
		_, ss = xray.BeginSubsegment(ctx, "UpsertRelAlertFinding")
		_, err := a.repository.UpsertRelAlertFinding(ctx, data)
		ss.Close(err)
		if err != nil {
			return nil, err
		}
	}

	return registerdData, nil
}

func (a *alertService) DeleteAlertByAnalyze(ctx context.Context, alertCondition *model.AlertCondition) error {
	// Alertの削除
	_, ss := xray.BeginSubsegment(ctx, "GetAlertByAlertConditionIDStatus")
	savedData, err := a.repository.GetAlertByAlertConditionIDStatus(ctx, alertCondition.ProjectID, alertCondition.AlertConditionID, []string{"ACTIVE", "PENDING"})
	ss.Close(err)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		appLogger.Errorf("Failed get alert by alertConditionIDStatus, err: %v", err)
		return err
	}
	// レコードが存在しない場合何もしない
	if noRecord {
		return nil
	}

	data := &model.Alert{
		AlertID:          savedData.AlertID,
		AlertConditionID: alertCondition.AlertConditionID,
		Description:      alertCondition.Description,
		Severity:         alertCondition.Severity,
		ProjectID:        alertCondition.ProjectID,
		Status:           alert.Status_DEACTIVE.String(),
	}
	// update Alert
	_, ss = xray.BeginSubsegment(ctx, "DeactivateAlert")
	err = a.repository.DeactivateAlert(ctx, data)
	ss.Close(err)
	if err != nil {
		return err
	}

	// AlertHistoryに登録するための現在のRelAlertFindingを整形
	findingHistory, err := makeFindingIDs([]uint64{})
	if err != nil {
		return err
	}

	// AlertHistoryの登録
	dataAlertHistory := &model.AlertHistory{
		HistoryType:    "deleted",
		AlertID:        savedData.AlertID,
		Description:    savedData.Description,
		Severity:       savedData.Severity,
		FindingHistory: findingHistory,
		ProjectID:      savedData.ProjectID,
	}
	_, ss = xray.BeginSubsegment(ctx, "UpsertAlertHistory")
	_, errAlertHistory := a.repository.UpsertAlertHistory(ctx, dataAlertHistory)
	ss.Close(err)
	if errAlertHistory != nil {
		return errAlertHistory
	}

	//RelAlertFindingの削除
	err = a.deleteRelAlertFindingByAlertID(ctx, savedData.ProjectID, savedData.AlertID)
	if err != nil {
		return err
	}

	return nil
}

func (a *alertService) deleteRelAlertFindingByAlertID(ctx context.Context, projectID, alertID uint32) error {
	_, ss := xray.BeginSubsegment(ctx, "ListRelAlertFinding")
	listRelAlertFinding, err := a.repository.ListRelAlertFinding(ctx, projectID, alertID, uint32(0), 0, time.Now().Unix())
	ss.Close(err)
	if err != nil {
		return err
	}
	for _, relAlertFinding := range *listRelAlertFinding {
		_, ss = xray.BeginSubsegment(ctx, "DeleteRelAlertFinding")
		err := a.repository.DeleteRelAlertFinding(ctx, relAlertFinding.ProjectID, relAlertFinding.AlertID, relAlertFinding.FindingID)
		ss.Close(err)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *alertService) NotificationAlert(ctx context.Context, alertCondition *model.AlertCondition, alert *model.Alert, rules *[]model.AlertRule) error {
	_, ss := xray.BeginSubsegment(ctx, "ListAlertCondNotification")
	alertCondNotifications, err := a.repository.ListAlertCondNotification(ctx, alertCondition.ProjectID, alertCondition.AlertConditionID, 0, 0, time.Now().Unix())
	ss.Close(err)
	if err != nil {
		return err
	}
	for _, alertCondNotification := range *alertCondNotifications {
		// 連続通知を防ぐ
		if time.Now().Unix() < alertCondNotification.NotifiedAt.Unix()+int64(alertCondNotification.CacheSecond) {
			continue
		}

		_, ss = xray.BeginSubsegment(ctx, "GetNotification")
		notification, err := a.repository.GetNotification(ctx, alertCondition.ProjectID, alertCondNotification.NotificationID)
		ss.Close(err)
		if err != nil {
			return err
		}
		switch notification.Type {
		case "slack":
			_, ss = xray.BeginSubsegment(ctx, "GetProject")
			project, err := a.repository.GetProject(ctx, alert.ProjectID)
			ss.Close(err)
			if err != nil {
				return err
			}
			_, ss = xray.BeginSubsegment(ctx, "sendSlackNotification")
			err = sendSlackNotification(notification.NotifySetting, alert, project, rules)
			ss.Close(err)
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
		_, ss = xray.BeginSubsegment(ctx, "UpsertAlertCondNotification")
		_, err = a.repository.UpsertAlertCondNotification(ctx, dataAlertCondNotification)
		ss.Close(err)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *alertService) analyzeAlertByRule(ctx context.Context, alertRule *model.AlertRule) (bool, *[]uint64, error) {
	param := &finding.BatchListFindingRequest{
		ProjectId: alertRule.ProjectID,
		FromScore: alertRule.Score,
		Status:    finding.FindingStatus_FINDING_ACTIVE,
	}
	if alertRule.ResourceName != "" {
		param.ResourceName = []string{alertRule.ResourceName}
	}
	if alertRule.Tag != "" {
		param.Tag = []string{alertRule.Tag}
	}
	_, ss := xray.BeginSubsegment(ctx, "BatchListFinding")
	resp, err := a.findingClient.BatchListFinding(ctx, param)
	ss.Close(err)
	if err != nil {
		appLogger.Errorf("Failed to BatchListFinding, request=%+v, err=%+v", param, err)
		return false, &[]uint64{}, err
	}
	appLogger.Infof("Got BatchListFinding, request=%+v, count=%d", param, resp.Count) // Debug
	return resp.Count >= alertRule.FindingCnt, &resp.FindingId, nil
}

func (a *alertService) isMatchExistingAlert(ctx context.Context, savedAlert *model.Alert, alertCondition *model.AlertCondition, findingIDs []uint64) (bool, error) {
	if savedAlert.Description != alertCondition.Description {
		return false, nil
	}
	if savedAlert.Severity != alertCondition.Severity {
		return false, nil
	}
	now := time.Now().Unix()
	relAlertFindings, err := a.repository.ListRelAlertFinding(ctx, savedAlert.ProjectID, savedAlert.AlertID, 0, 0, now)
	if err != nil {
		return false, err
	}
	if len(*relAlertFindings) != len(findingIDs) {
		return false, nil
	}
	for _, relAlertFinding := range *relAlertFindings {
		if !isContainsFindings(uint64(relAlertFinding.FindingID), findingIDs) {
			return false, nil
		}
	}
	return true, nil
}

func makeFindingIDs(findingIDs []uint64) (string, error) {
	mapFindingIDs := map[string][]uint64{"finding_id": findingIDs}
	bytes, err := json.Marshal(mapFindingIDs)
	if err != nil {
		appLogger.Error("JSON marshal error when making FindingIDs ", err)
		return "", err
	}
	return string(bytes), nil
}

func isContainsFindings(targetID uint64, findingIDs []uint64) bool {
	for _, findingID := range findingIDs {
		if findingID == targetID {
			return true
		}
	}
	return false
}

func getHistoryType(alertID uint32) string {
	if zero.IsZeroVal(alertID) {
		return "created"
	}
	return "updated"
}

func sendSlackNotification(notifySetting string, alert *model.Alert, project *model.Project, rules *[]model.AlertRule) error {
	var setting slackNotifySetting
	if err := json.Unmarshal([]byte(notifySetting), &setting); err != nil {
		return err
	}
	if zero.IsZeroVal(setting.WebhookURL) {
		appLogger.Warn("Unset webhook_url")
		return nil
	}
	channel := ""
	if !zero.IsZeroVal(setting.Data["channel"]) {
		channel = setting.Data["channel"]
	}
	message := ""
	if !zero.IsZeroVal(setting.Data["message"]) {
		message = setting.Data["message"]
	}
	slackConfig, err := newslackWebhookConfig()
	if err != nil {
		return err
	}

	payload, err := slackConfig.GetPayload(channel, message, alert, project, rules)
	if err != nil {
		return err
	}
	// TODO http tracing
	resp, err := http.PostForm(setting.WebhookURL, url.Values{"payload": {string(payload)}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func sendSlackTestNotification(notifySetting string) error {
	var setting slackNotifySetting
	if err := json.Unmarshal([]byte(notifySetting), &setting); err != nil {
		return err
	}
	if zero.IsZeroVal(setting.WebhookURL) {
		appLogger.Warn("Unset webhook_url")
		return nil
	}
	channel := ""
	if !zero.IsZeroVal(setting.Data["channel"]) {
		channel = setting.Data["channel"]
	}
	slackConfig, err := newslackWebhookConfig()
	if err != nil {
		return err
	}

	payload, err := slackConfig.GetTestPayload(channel)
	if err != nil {
		return err
	}
	// TODO http tracing
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

// for Logging
func makeRandomStr(digit uint32) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 乱数を生成
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("unexpected error...")
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}

func getRequestID(projectID uint32) string {
	rand, err := makeRandomStr(10)
	if err != nil {
		appLogger.Warnf("Failed to make random string, err=%+v", err)
	}
	return fmt.Sprintf("%d-%s", projectID, rand)
}

func (a *alertService) AnalyzeAlertAll(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	appLogger.Info("start AnalyzeAlertAll")
	list, err := a.projectClient.ListProject(ctx, &project.ListProjectRequest{})
	if err != nil {
		appLogger.Errorf("Failed to list project API, err=%+v", err)
		return nil, err
	}
	projectNum := len(list.Project)
	if projectNum < 1 {
		appLogger.Warn("There are no project")
		return &empty.Empty{}, nil
	}
	var wg sync.WaitGroup
	wg.Add(projectNum)
	for i := range list.Project {
		// launch goroutine
		go func(p *project.Project) {
			defer wg.Done()
			if _, err := a.AnalyzeAlert(ctx, &alert.AnalyzeAlertRequest{ProjectId: p.ProjectId}); err != nil {
				appLogger.Warnf("Failed to AnalyzeAlert, project_id=%d, err=%+v", p.ProjectId, err)
			}
			time.Sleep(1 * time.Second)
		}(list.Project[i])
	}
	wg.Wait()
	appLogger.Info("end AnalyzeAlertAll")
	return &empty.Empty{}, nil
}
