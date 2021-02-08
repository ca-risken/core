package main

import (
	"strings"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	_ "github.com/go-sql-driver/mysql"
)

func (f *reportDB) GetReportFinding(projectID uint32, dataSource []string, fromDate, toDate string, score float32) (*[]model.ReportFinding, error) {
	query := `select * from report_finding where project_id = ? and score > ?`
	var params []interface{}
	params = append(params, projectID, score)
	if len(dataSource) != 0 {
		query += " and data_source regexp ?"
		params = append(params, strings.Join(dataSource, "|"))
	}
	if len(fromDate) != 0 {
		query += " and report_date >= ?"
		params = append(params, fromDate)
	}
	if len(toDate) != 0 {
		query += " and report_date <= ?"
		params = append(params, toDate)
	}
	var data []model.ReportFinding
	if err := f.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *reportDB) GetReportFindingAll(dataSource []string, fromDate, toDate string, score float32) (*[]model.ReportFinding, error) {
	query := `select * from report_finding where score > ?`
	var params []interface{}
	params = append(params, score)
	if len(dataSource) != 0 {
		query += " and data_source regexp ?"
		params = append(params, strings.Join(dataSource, "|"))
	}
	if len(fromDate) != 0 {
		query += " and report_date >= ?"
		params = append(params, fromDate)
	}
	if len(toDate) != 0 {
		query += " and report_date <= ?"
		params = append(params, toDate)
	}
	var data []model.ReportFinding
	if err := f.Slave.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *reportDB) CollectReportFinding() error {
	query := `insert into report_finding (report_date, project_id, data_source, score, count) 
	select DATE_ADD(CURRENT_DATE, INTERVAL -1 DAY) as report_date, project_id, data_source, score, count(*) as count 
	from finding f
	where f.finding_id not in (select finding_id from pend_finding)
	group by f.project_id, data_source, score ON DUPLICATE KEY UPDATE count=values(count)`
	var data []model.ReportFinding
	if err := f.Slave.Raw(query).Scan(&data).Error; err != nil {
		return err
	}
	return nil
}
