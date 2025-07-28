package db

import (
	"context"
	"strings"

	"github.com/ca-risken/core/pkg/model"
)

type ReportRepository interface {
	// Report
	GetReportFinding(context.Context, uint32, []string, string, string, float32) (*[]model.ReportFinding, error)
	GetReportFindingAll(context.Context, []string, string, string, float32) (*[]model.ReportFinding, error)
	CollectReportFinding(ctx context.Context) error
	PurgeReportFinding(ctx context.Context) error
	ListReport(ctx context.Context, projectID uint32) (*[]model.Report, error)
	GetReport(ctx context.Context, projectID uint32, reportID uint32) (*model.Report, error)
	PutReport(ctx context.Context, report *model.Report) (*model.Report, error)
}

var _ ReportRepository = (*Client)(nil)

func (c *Client) GetReportFinding(ctx context.Context, projectID uint32, dataSource []string, fromDate, toDate string, score float32) (*[]model.ReportFinding, error) {
	query := `select r.*,p.name as project_name from report_finding as r, project as p where r.project_id = ? and r.project_id = p.project_id and score >= ?`
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
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetReportFindingAll(ctx context.Context, dataSource []string, fromDate, toDate string, score float32) (*[]model.ReportFinding, error) {
	query := `select r.*,p.name as project_name from report_finding as r, project as p where r.project_id = p.project_id and score >= ?`
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
	if err := c.Slave.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) CollectReportFinding(ctx context.Context) error {
	query := `insert into report_finding (report_date, project_id, data_source, score, count) 
select DATE_ADD(CURRENT_DATE, INTERVAL -1 DAY) as report_date, project_id, data_source, score , count(*) as count 
from finding f
where not exists (select pend_finding.finding_id from pend_finding where f.finding_id = pend_finding.finding_id and (pend_finding.expired_at is NULL or pend_finding.expired_at <= NOW())) 
group by f.project_id, data_source, score ON DUPLICATE KEY UPDATE count=values(count)`
	var data []model.ReportFinding
	if err := c.Master.WithContext(ctx).Raw(query).Scan(&data).Error; err != nil {
		return err
	}
	return nil
}

const deleteReportFinding = "delete from report_finding where report_date < DATE_ADD(CURRENT_DATE, INTERVAL -365 DAY)"

func (c *Client) PurgeReportFinding(ctx context.Context) error {
	if err := c.Master.WithContext(ctx).Exec(deleteReportFinding).Error; err != nil {
		return err
	}
	return nil
}

const selectListReport = `select * from report where project_id = ?`

func (c *Client) ListReport(ctx context.Context, projectID uint32) (*[]model.Report, error) {
	var data []model.Report
	if err := c.Slave.WithContext(ctx).Raw(selectListReport, projectID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const selectGetReport = `select * from report where project_id = ? and report_id = ?`

func (c *Client) GetReport(ctx context.Context, projectID uint32, reportID uint32) (*model.Report, error) {
	var data model.Report
	if err := c.Slave.WithContext(ctx).Raw(selectGetReport, projectID, reportID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertPutReport = `
INSERT INTO report
  (report_id, project_id, name, type, status, content)
VALUES
  (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  name=VALUES(name),
  type=VALUES(type),
  status=VALUES(status),
  content=VALUES(content)
`

func (c *Client) PutReport(ctx context.Context, report *model.Report) (*model.Report, error) {
	if err := c.Master.WithContext(ctx).Exec(insertPutReport,
		report.ReportID,
		report.ProjectID,
		report.Name,
		report.Type,
		report.Status,
		report.Content,
	).Error; err != nil {
		return nil, err
	}
	return c.getReportByName(ctx, report.ProjectID, report.Name) // if insert, return new report(new report_id)
}

// selectGetReportByName: Unique Key(project_id, name)
const selectGetReportByName = `select * from report where project_id = ? and name = ?`

func (c *Client) getReportByName(ctx context.Context, projectID uint32, name string) (*model.Report, error) {
	var data model.Report
	if err := c.Master.WithContext(ctx).Raw(selectGetReportByName, projectID, name).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
