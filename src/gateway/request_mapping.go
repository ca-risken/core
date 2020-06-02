package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/CyberAgent/mimosa-core/proto/finding"
)

func mappingListFindingRequest(r *http.Request) *finding.ListFindingRequest {
	req := &finding.ListFindingRequest{}
	if key := r.URL.Query().Get("project_id"); key != "" {
		req.ProjectId = commaSeparator(r.URL.Query().Get("project_id"))
	}
	if key := r.URL.Query().Get("data_source"); key != "" {
		req.DataSource = commaSeparator(r.URL.Query().Get("data_source"))
	}
	if key := r.URL.Query().Get("resource_name"); key != "" {
		req.ResourceName = commaSeparator(r.URL.Query().Get("resource_name"))
	}
	if key := r.URL.Query().Get("from_score"); key != "" {
		if score, err := parseScore(r.URL.Query().Get("from_score")); err == nil {
			req.FromScore = score
		}
	}
	if key := r.URL.Query().Get("to_score"); key != "" {
		if score, err := parseScore(r.URL.Query().Get("to_score")); err == nil {
			req.ToScore = score
		}
	}
	if key := r.URL.Query().Get("from_at"); key != "" {
		if t, err := parseTimeParam(r.URL.Query().Get("from_at")); err == nil {
			req.FromAt = t
		}
	}
	if key := r.URL.Query().Get("to_at"); key != "" {
		if t, err := parseTimeParam(r.URL.Query().Get("to_at")); err == nil {
			req.ToAt = t
		}
	}

	return req
}

func commaSeparator(param string) []int32 {
	separated := []int32{}
	for _, p := range strings.Split(param, ",") {
		if i, err := strconv.Atoi(p); err == nil {
			append(separated, i)
		}
	}
	return separated
}

func parseScore(score string) (float32, error) {
	f, err := strconv.ParseFloat(score, 32)
	if err != nil {
		return 0.0, err
	}
	return float32(f), nil
}

func parseTimeParam(at string) (int64, error) {
	i, err := strconv.ParseInt(at, 10, 64)
	if err != nil {
		return i, err
	}
	return i, nil
}
