package alert

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	projectproto "github.com/ca-risken/core/proto/project"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vikyd/zero"
	"golang.org/x/sync/semaphore"
	"gorm.io/gorm"
)

func (a *AlertService) AnalyzeAlert(ctx context.Context, req *alert.AnalyzeAlertRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	projects, err := a.projectClient.ListProject(ctx, &projectproto.ListProjectRequest{ProjectId: req.ProjectId})
	if err != nil {
		return nil, err
	}
	if len(projects.Project) == 0 {
		return nil, fmt.Errorf("not found project: ProjectID:=%v", req.ProjectId)
	}
	project := projects.Project[0]

	requestID := a.getRequestID(ctx, req.ProjectId)
	// 有効なalertConditionの取得
	a.logger.Infof(ctx, "start ListEnabledAlertCondition: RequestID=%s", requestID)
	alertConditions, err := a.repository.ListEnabledAlertCondition(ctx, req.ProjectId, req.AlertConditionId)
	a.logger.Infof(ctx, "finish ListEnabledAlertCondition: RequestID=%s", requestID)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		a.logger.Error(ctx, err)
		return nil, err
	}

	// マッチング
	a.logger.Infof(ctx, "start matching: RequestID=%s", requestID)
	for _, alertCondition := range *alertConditions {
		err := a.AnalyzeAlertByCondition(ctx, &alertCondition, project)
		if err != nil {
			a.logger.Error(ctx, err)
			return nil, err
		}
	}
	a.logger.Infof(ctx, "finish matching: RequestID=%s", requestID)

	// 無効のalertConditionの取得
	a.logger.Infof(ctx, "start ListDisabledAlertCondition: RequestID=%s", requestID)
	disabledAlertConditions, err := a.repository.ListDisabledAlertCondition(ctx, req.ProjectId, req.AlertConditionId)
	a.logger.Infof(ctx, "finish ListDisabledAlertCondition: RequestID=%s", requestID)
	noRecord = errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		a.logger.Error(ctx, err)
		return nil, err
	}
	// 無効なalertConditionに紐づくAlertを削除
	a.logger.Infof(ctx, "start DeleteAlertByAnalyze: RequestID=%s", requestID)
	for _, alertCondition := range *disabledAlertConditions {
		err := a.DeleteAlertByAnalyze(ctx, &alertCondition)
		if err != nil {
			a.logger.Error(ctx, err)
			return nil, err
		}
	}
	a.logger.Infof(ctx, "finish DeleteAlertByAnalyze: RequestID=%s", requestID)
	return &empty.Empty{}, nil
}

func (a *AlertService) AnalyzeAlertByCondition(ctx context.Context, alertCondition *model.AlertCondition, project *projectproto.Project) error {
	// AlertRuleの取得
	alertRules, err := a.repository.ListAlertRuleByAlertConditionID(ctx, alertCondition.ProjectID, alertCondition.AlertConditionID)
	if err != nil {
		a.logger.Errorf(ctx, "Failed list AlertRule by AlertConditionID. alertConditionID: %v, err: %v", alertCondition.AlertConditionID, err)
		return err
	}
	var matchFindingIDs []uint64
	isFirst := true
	a.logger.Info(ctx, "start matching per rule")
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
	a.logger.Info(ctx, "finish matching per rule")
	if len(matchFindingIDs) > 0 {
		registAlert, existsNewFindings, err := a.RegistAlertByAnalyze(ctx, alertCondition, matchFindingIDs)
		if err != nil {
			return err
		}
		// AlertがACTIVE、かつMatchしている場合はAlert通知を行う
		if registAlert.Status == alert.Status_ACTIVE.String() {
			err = a.NotificationAlert(ctx, alertCondition, registAlert, alertRules, project, &matchFindingIDs, existsNewFindings)
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

// RegistAlertByAnalyze
// アラートを登録・更新します。既に登録済みのアラートがある場合は処理をスキップします。
// Return:
//  1. *model.Alert (登録したアラート)
//  2. bool (新規Findingがあったかどうか)
//  3. error (エラー)
func (a *AlertService) RegistAlertByAnalyze(ctx context.Context, alertCondition *model.AlertCondition, findingIDs []uint64) (*model.Alert, bool, error) {
	// AlertConditionに該当するAlertが既に存在しているか確認
	savedData, err := a.repository.GetAlertByAlertConditionIDStatus(ctx, alertCondition.ProjectID, alertCondition.AlertConditionID, []string{"ACTIVE", "PENDING"})
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		return nil, false, err
	}

	// 既に登録済みの場合は登録済みのAlertと内容が一致しているか確認
	// 一致している場合処理を終了する
	// 登録済みかつStatusがPENDINGの場合、PENDING
	var status string
	var alertID uint32
	compareLatestAlertFinding := &compareLatestAlertFindingResult{}
	if !noRecord {
		compareLatestAlertFinding, err = a.compareLatestAlertFinding(ctx, savedData, alertCondition, findingIDs)
		if err != nil {
			return nil, false, err
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
	registeredData, err := a.repository.UpsertAlert(ctx, data)
	if err != nil {
		a.logger.Errorf(ctx, "Error occurred when upsert alert. alertConditionID: %v, err: %v", alertCondition.AlertConditionID, err)
		return nil, false, err
	}

	// 過去のアラートと状態が同じなら以下の処理はスキップ
	if compareLatestAlertFinding.isMatchAlertFindings {
		return registeredData, compareLatestAlertFinding.existsNewFindings, nil
	}

	// AlertHistoryに登録するための現在のRelAlertFindingを整形
	findingHistory, err := a.makeFindingIDs(ctx, findingIDs)
	if err != nil {
		return nil, false, err
	}

	historyType := getHistoryType(alertID)
	// AlertHistoryの登録
	dataAlertHistory := &model.AlertHistory{
		HistoryType:    historyType,
		AlertID:        registeredData.AlertID,
		Description:    registeredData.Description,
		FindingHistory: findingHistory,
		Severity:       registeredData.Severity,
		ProjectID:      registeredData.ProjectID,
	}
	_, err = a.repository.UpsertAlertHistory(ctx, dataAlertHistory)
	if err != nil {
		a.logger.Errorf(ctx, "Error occurred when upsert AlertHistory. err: %v", err)
		return nil, false, err
	}

	//RelAlertFindingの更新 (削除して再登録)
	err = a.deleteRelAlertFindingByAlertID(ctx, registeredData.ProjectID, registeredData.AlertID)
	if err != nil {
		return nil, false, err
	}
	for _, findingID := range findingIDs {
		data := &model.RelAlertFinding{
			AlertID:   registeredData.AlertID,
			FindingID: findingID,
			ProjectID: registeredData.ProjectID,
		}
		_, err := a.repository.UpsertRelAlertFinding(ctx, data)
		if err != nil {
			return nil, false, err
		}
	}

	return registeredData, compareLatestAlertFinding.existsNewFindings, nil
}

func (a *AlertService) DeleteAlertByAnalyze(ctx context.Context, alertCondition *model.AlertCondition) error {
	// Alertの削除
	savedData, err := a.repository.GetAlertByAlertConditionIDStatus(ctx, alertCondition.ProjectID, alertCondition.AlertConditionID, []string{"ACTIVE", "PENDING"})
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		a.logger.Errorf(ctx, "Failed get alert by alertConditionIDStatus, err: %v", err)
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
	err = a.repository.DeactivateAlert(ctx, data)
	if err != nil {
		return err
	}

	// AlertHistoryに登録するための現在のRelAlertFindingを整形
	findingHistory, err := a.makeFindingIDs(ctx, []uint64{})
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
	_, errAlertHistory := a.repository.UpsertAlertHistory(ctx, dataAlertHistory)
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

func (a *AlertService) deleteRelAlertFindingByAlertID(ctx context.Context, projectID, alertID uint32) error {
	listRelAlertFinding, err := a.repository.ListRelAlertFinding(ctx, projectID, alertID, uint64(0), 0, time.Now().Unix())
	if err != nil {
		return err
	}
	for _, relAlertFinding := range *listRelAlertFinding {
		err := a.repository.DeleteRelAlertFinding(ctx, relAlertFinding.ProjectID, relAlertFinding.AlertID, relAlertFinding.FindingID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AlertService) NotificationAlert(
	ctx context.Context,
	alertCondition *model.AlertCondition,
	alert *model.Alert,
	rules *[]model.AlertRule,
	project *projectproto.Project,
	findingIDs *[]uint64,
	existsNewFindings bool,
) error {
	alertCondNotifications, err := a.repository.ListAlertCondNotification(ctx, alertCondition.ProjectID, alertCondition.AlertConditionID, 0, 0, time.Now().Unix())
	if err != nil {
		return err
	}
	findings, err := a.getFindingDetailsForNotification(ctx, project.ProjectId, findingIDs)
	if err != nil {
		return err
	}
	for _, alertCondNotification := range *alertCondNotifications {
		// 連続通知を防ぐ
		if !existsNewFindings && time.Now().Unix() < alertCondNotification.NotifiedAt.Unix()+int64(alertCondNotification.CacheSecond) {
			continue
		}

		notification, err := a.repository.GetNotification(ctx, alertCondition.ProjectID, alertCondNotification.NotificationID)
		if err != nil {
			return err
		}
		switch notification.Type {
		case "slack":
			err = sendSlackNotification(ctx, a.baseURL, notification.NotifySetting, alert, project, rules, findings, a.defaultLocale)
			if err != nil {
				return fmt.Errorf("notify error: notification_id=%d, err=%w", notification.NotificationID, err)
			}
		default:
			a.logger.Warn(ctx, "This notification_type is unimprement.", notification.Type)
		}
		// 通知時刻を更新する
		dataAlertCondNotification := &model.AlertCondNotification{
			AlertConditionID: alertCondNotification.AlertConditionID,
			NotificationID:   alertCondNotification.NotificationID,
			CacheSecond:      alertCondNotification.CacheSecond,
			NotifiedAt:       time.Now(),
			ProjectID:        alertCondNotification.ProjectID,
		}
		_, err = a.repository.UpsertAlertCondNotification(ctx, dataAlertCondNotification)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AlertService) analyzeAlertByRule(ctx context.Context, alertRule *model.AlertRule) (bool, *[]uint64, error) {
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
	resp, err := a.findingClient.BatchListFinding(ctx, param)
	if err != nil {
		a.logger.Errorf(ctx, "Failed to BatchListFinding, request=%+v, err=%+v", param, err)
		return false, &[]uint64{}, err
	}
	a.logger.Infof(ctx, "Got BatchListFinding, request=%+v, count=%d", param, resp.Count) // Debug
	return resp.Count >= alertRule.FindingCnt, &resp.FindingId, nil
}

type compareLatestAlertFindingResult struct {
	isMatchAlertFindings bool
	existsNewFindings    bool
}

func (a *AlertService) compareLatestAlertFinding(
	ctx context.Context, savedAlert *model.Alert, alertCondition *model.AlertCondition, findingIDs []uint64,
) (
	*compareLatestAlertFindingResult, error,
) {
	result := compareLatestAlertFindingResult{
		isMatchAlertFindings: false,
		existsNewFindings:    false,
	}
	if savedAlert.Description != alertCondition.Description {
		return &result, nil
	}
	if savedAlert.Severity != alertCondition.Severity {
		return &result, nil
	}
	now := time.Now().Unix()
	relAlertFindings, err := a.repository.ListRelAlertFinding(ctx, savedAlert.ProjectID, savedAlert.AlertID, 0, 0, now)
	if err != nil {
		return nil, err
	}
	alertFindingMap := map[uint64]bool{}
	for _, relAlertFinding := range *relAlertFindings {
		alertFindingMap[relAlertFinding.FindingID] = true
	}
	for _, findingID := range findingIDs {
		if _, exists := alertFindingMap[findingID]; !exists {
			result.existsNewFindings = true
			break
		}
	}
	if len(*relAlertFindings) == len(findingIDs) && !result.existsNewFindings {
		result.isMatchAlertFindings = true
	}
	return &result, nil
}

func (a *AlertService) makeFindingIDs(ctx context.Context, findingIDs []uint64) (string, error) {
	mapFindingIDs := map[string][]uint64{"finding_id": findingIDs}
	bytes, err := json.Marshal(mapFindingIDs)
	if err != nil {
		a.logger.Error(ctx, "JSON marshal error when making FindingIDs ", err)
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

// for Logging
func makeRandomStr(digit uint32) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 乱数を生成
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to read random, err=%w", err)
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}

func (a *AlertService) getRequestID(ctx context.Context, projectID uint32) string {
	rand, err := makeRandomStr(10)
	if err != nil {
		a.logger.Warnf(ctx, "Failed to make random string, err=%+v", err)
	}
	return fmt.Sprintf("%d-%s", projectID, rand)
}

func (a *AlertService) AnalyzeAlertAll(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	a.logger.Info(ctx, "start AnalyzeAlertAll")
	list, err := a.projectClient.ListProject(ctx, &projectproto.ListProjectRequest{})
	if err != nil {
		a.logger.Errorf(ctx, "Failed to list project API, err=%+v", err)
		return nil, err
	}
	projectNum := len(list.Project)
	if projectNum < 1 {
		a.logger.Warn(ctx, "There are no project")
		return &empty.Empty{}, nil
	}
	allResource := a.maxAnalyzeAPICall // max groutine resources
	analyzeAlertResource := int64(1)
	sem := semaphore.NewWeighted(allResource)
	var wg sync.WaitGroup
	wg.Add(projectNum)
	for i := range list.Project {
		if err := sem.Acquire(ctx, analyzeAlertResource); err != nil {
			a.logger.Errorf(ctx, "Failed to acquire resource, err=%+v", err)
			return nil, err
		}
		// launch goroutine
		go func(p *projectproto.Project) {
			defer wg.Done()
			active, err := a.projectClient.IsActive(ctx, &projectproto.IsActiveRequest{ProjectId: p.ProjectId})
			if err != nil {
				// TODO AnalyzeAlertの呼び出しを非同期化したら、通常のログ出力に変更してAnalyzeAlertAllがerrorを返すよう修正
				a.logger.Notifyf(ctx, logging.ErrorLevel, "Failed to call API (project.IsActive), err=%+v", err)
			}
			if active != nil && active.Active {
				if _, err := a.AnalyzeAlert(ctx, &alert.AnalyzeAlertRequest{ProjectId: p.ProjectId}); err != nil {
					// TODO AnalyzeAlertの呼び出しを非同期化したら、通常のログ出力に変更してAnalyzeAlertAllがerrorを返すよう修正
					a.logger.Notifyf(ctx, logging.ErrorLevel, "Failed to AnalyzeAlert, project_id=%d, err=%+v", p.ProjectId, err)
				}
			}
			time.Sleep(1 * time.Second)
			sem.Release(analyzeAlertResource)
		}(list.Project[i])
	}
	wg.Wait()
	a.logger.Info(ctx, "end AnalyzeAlertAll")
	return &empty.Empty{}, nil
}
