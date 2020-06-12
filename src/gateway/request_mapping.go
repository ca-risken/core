package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/go-chi/chi"
)

func mappingListFindingRequest(r *http.Request) *finding.ListFindingRequest {
	req := finding.ListFindingRequest{}
	if param := r.URL.Query().Get("project_id"); param != "" {
		req.ProjectId = commaSeparatorID(param)
	}
	if param := r.URL.Query().Get("data_source"); param != "" {
		req.DataSource = commaSeparator(param)
	}
	if param := r.URL.Query().Get("resource_name"); param != "" {
		req.ResourceName = commaSeparator(param)
	}
	if param := r.URL.Query().Get("from_score"); param != "" {
		req.FromScore = parseScore(param)
	}
	if param := r.URL.Query().Get("to_score"); param != "" {
		req.ToScore = parseScore(param)
	}
	if param := r.URL.Query().Get("from_at"); param != "" {
		req.FromAt = parseAt(param)
	}
	if param := r.URL.Query().Get("to_at"); param != "" {
		req.ToAt = parseAt(param)
	}
	return &req
}

func mappingGetFindingRequest(r *http.Request) *finding.GetFindingRequest {
	return &finding.GetFindingRequest{
		FindingId: parseUint64(chi.URLParam(r, "finding_id")),
	}
}

func mappingPutFindingRequest(r *http.Request) *finding.PutFindingRequest {
	param := finding.FindingForUpsert{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		appLogger.Warnf("Invalid parameter in PutFindingRequest, err: %+v", err)
		return &finding.PutFindingRequest{Finding: &param}
	}
	return &finding.PutFindingRequest{
		Finding: &finding.FindingForUpsert{
			Description:      param.Description,
			DataSource:       param.DataSource,
			DataSourceId:     param.DataSourceId,
			ResourceName:     param.ResourceName,
			ProjectId:        param.ProjectId,
			OriginalScore:    param.OriginalScore,
			OriginalMaxScore: param.OriginalMaxScore,
			Data:             param.Data,
		},
	}
}

func commaSeparatorID(param string) []uint32 {
	separated := []uint32{}
	for _, p := range strings.Split(param, ",") {
		if i, err := strconv.Atoi(p); err == nil {
			separated = append(separated, uint32(i))
		}
	}
	return separated
}

func commaSeparator(param string) []string {
	separated := []string{}
	for _, p := range strings.Split(param, ",") {
		if p != "" {
			separated = append(separated, p)
		}
	}
	return separated
}

func parseScore(score string) float32 {
	f, err := strconv.ParseFloat(score, 32)
	if err != nil {
		return 0.0
	}
	return float32(f)
}

func parseAt(at string) int64 {
	i, err := strconv.ParseInt(at, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func parseUint64(str string) uint64 {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}
	return i
}
