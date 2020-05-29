package main

import (
	"net/http"

	"github.com/CyberAgent/mimosa-core/proto/finding"
)

func mappingListFindingRequest(r *http.Request) *finding.ListFindingRequest {
	return &finding.ListFindingRequest{
		ProjectId: r.URL.Query().Get("project_id"),
		Since:     r.URL.Query().Get("since"),
	}
}
