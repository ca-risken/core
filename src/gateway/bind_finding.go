package main

import (
	"net/http"

	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/go-chi/chi"
)

func bindListFindingRequest(r *http.Request) *finding.ListFindingRequest {
	req := finding.ListFindingRequest{}
	if err := bindQuery(&req, r); err != nil {
		appLogger.Warnf("Invalid parmeter. err=%+v", err)
	}
	if len(req.ProjectId) > 0 {
		req.ProjectId = ignoreZeroValue(req.ProjectId)
	}
	if len(req.DataSource) > 0 {
		req.DataSource = commaSeparator(req.DataSource[0])
	}
	if len(req.ResourceName) > 0 {
		req.ResourceName = commaSeparator(req.ResourceName[0])
	}
	return &req
}

func bindGetFindingRequest(r *http.Request) *finding.GetFindingRequest {
	return &finding.GetFindingRequest{
		FindingId: parseUint64(chi.URLParam(r, "finding_id")),
	}
}

func bindPutFindingRequest(r *http.Request) *finding.PutFindingRequest {
	req := finding.FindingForUpsert{}
	if err := bindBodyJSON(&req, r); err != nil {
		appLogger.Warnf("Invalid PutFindingRequest. err=%+v", err)
	}
	return &finding.PutFindingRequest{Finding: &req}
}

func bindDeleteFindingRequest(r *http.Request) *finding.DeleteFindingRequest {
	req := finding.DeleteFindingRequest{}
	if err := bindBodyJSON(&req, r); err != nil {
		appLogger.Warnf("Invalid DeleteFindingRequest. err=%+v", err)
	}
	return &req
}

func bindListFindingTagRequest(r *http.Request) *finding.ListFindingTagRequest {
	return &finding.ListFindingTagRequest{
		FindingId: parseUint64(chi.URLParam(r, "finding_id")),
	}
}

func bindTagFindingRequest(r *http.Request) *finding.TagFindingRequest {
	req := finding.FindingTagForUpsert{}
	if err := bindBodyJSON(&req, r); err != nil {
		appLogger.Warnf("Invalid FindingTagForUpsert. err=%+v", err)
	}
	return &finding.TagFindingRequest{Tag: &req}
}

func bindUntagFindingRequest(r *http.Request) *finding.UntagFindingRequest {
	req := finding.UntagFindingRequest{}
	if err := bindBodyJSON(&req, r); err != nil {
		appLogger.Warnf("Invalid UntagFindingRequest. err=%+v", err)
	}
	return &req
}

func bindListResourceRequest(r *http.Request) *finding.ListResourceRequest {
	req := finding.ListResourceRequest{}
	if err := bindQuery(&req, r); err != nil {
		appLogger.Warnf("Invalid parmeter. err=%+v", err)
	}
	if len(req.ProjectId) > 0 {
		req.ProjectId = ignoreZeroValue(req.ProjectId)
	}
	if len(req.ResourceName) > 0 {
		req.ResourceName = commaSeparator(req.ResourceName[0])
	}

	return &req
}

func bindGetResourceRequest(r *http.Request) *finding.GetResourceRequest {
	return &finding.GetResourceRequest{
		ResourceId: parseUint64(chi.URLParam(r, "resource_id")),
	}
}

func bindPutResourceRequest(r *http.Request) *finding.PutResourceRequest {
	req := finding.ResourceForUpsert{}
	if err := bindBodyJSON(&req, r); err != nil {
		appLogger.Warnf("Invalid ResourceForUpsert. err=%+v", err)
	}
	return &finding.PutResourceRequest{Resource: &req}
}

func bindDeleteResourceRequest(r *http.Request) *finding.DeleteResourceRequest {
	req := finding.DeleteResourceRequest{}
	if err := bindBodyJSON(&req, r); err != nil {
		appLogger.Warnf("Invalid DeleteResourceRequest. err=%+v", err)
	}
	return &req
}

func bindListResourceTagRequest(r *http.Request) *finding.ListResourceTagRequest {
	return &finding.ListResourceTagRequest{
		ResourceId: parseUint64(chi.URLParam(r, "resource_id")),
	}
}

func bindTagResourceRequest(r *http.Request) *finding.TagResourceRequest {
	req := finding.ResourceTagForUpsert{}
	if err := bindBodyJSON(&req, r); err != nil {
		appLogger.Warnf("Invalid ResourceTagForUpsert. err=%+v", err)
	}
	return &finding.TagResourceRequest{Tag: &req}
}

func bindUntagResourceRequest(r *http.Request) *finding.UntagResourceRequest {
	req := finding.UntagResourceRequest{}
	if err := bindBodyJSON(&req, r); err != nil {
		appLogger.Warnf("Invalid UntagResourceRequest. err=%+v", err)
	}
	return &req
}
