package main

import (
	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
)

const selectGetPendFinding = `select * from pend_finding where project_id = ? and finding_id = ?`

func (f *findingDB) GetPendFinding(projectID uint32, findingID uint64) (*model.PendFinding, error) {
	var data model.PendFinding
	if err := f.Master.Raw(selectGetPendFinding, projectID, findingID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

const insertPendFinding = `
INSERT INTO pend_finding
  (finding_id, project_id, note)
VALUES
  (?, ?, ?)
ON DUPLICATE KEY UPDATE
  updated_at = CURRENT_TIMESTAMP()
`

func (f *findingDB) UpsertPendFinding(pend *finding.PendFindingForUpsert) (*model.PendFinding, error) {
	if err := f.Master.Exec(insertPendFinding, pend.FindingId, pend.ProjectId, pend.Note).Error; err != nil {
		return nil, err
	}
	return f.GetPendFinding(pend.ProjectId, pend.FindingId)
}

const deletePendFinding = `delete from pend_finding where project_id = ? and finding_id = ?`

func (f *findingDB) DeletePendFinding(projectID uint32, findingID uint64) error {
	return f.Master.Exec(deletePendFinding, projectID, findingID).Error
}
