package main

import (
	"net/http"

	"github.com/CyberAgent/mimosa-core/proto/finding"
)

func bindListFindingRequest(r *http.Request) *finding.ListFindingRequest {
	req := finding.ListFindingRequest{}
	bind(&req, r)
	if len(req.DataSource) > 0 {
		req.DataSource = commaSeparator(req.DataSource[0])
	}
	if len(req.ResourceName) > 0 {
		req.ResourceName = commaSeparator(req.ResourceName[0])
	}
	return &req
}

func bindGetFindingRequest(r *http.Request) *finding.GetFindingRequest {
	req := finding.GetFindingRequest{}
	bind(&req, r)
	return &req
}

func bindPutFindingRequest(r *http.Request) *finding.PutFindingRequest {
	req := finding.PutFindingRequest{}
	bind(&req, r)
	return &req
}

func bindDeleteFindingRequest(r *http.Request) *finding.DeleteFindingRequest {
	req := finding.DeleteFindingRequest{}
	bind(&req, r)
	return &req
}

func bindListFindingTagRequest(r *http.Request) *finding.ListFindingTagRequest {
	req := finding.ListFindingTagRequest{}
	bind(&req, r)
	return &req
}

func bindTagFindingRequest(r *http.Request) *finding.TagFindingRequest {
	req := finding.TagFindingRequest{}
	bind(&req, r)
	return &req
}

func bindUntagFindingRequest(r *http.Request) *finding.UntagFindingRequest {
	req := finding.UntagFindingRequest{}
	bind(&req, r)
	return &req
}

func bindListResourceRequest(r *http.Request) *finding.ListResourceRequest {
	req := finding.ListResourceRequest{}
	bind(&req, r)
	if len(req.ResourceName) > 0 {
		req.ResourceName = commaSeparator(req.ResourceName[0])
	}
	return &req
}

func bindGetResourceRequest(r *http.Request) *finding.GetResourceRequest {
	req := finding.GetResourceRequest{}
	bind(&req, r)
	return &req
}

func bindPutResourceRequest(r *http.Request) *finding.PutResourceRequest {
	req := finding.PutResourceRequest{}
	bind(&req, r)
	return &req
}

func bindDeleteResourceRequest(r *http.Request) *finding.DeleteResourceRequest {
	req := finding.DeleteResourceRequest{}
	bind(&req, r)
	return &req
}

func bindListResourceTagRequest(r *http.Request) *finding.ListResourceTagRequest {
	req := finding.ListResourceTagRequest{}
	bind(&req, r)
	return &req
}

func bindTagResourceRequest(r *http.Request) *finding.TagResourceRequest {
	req := finding.TagResourceRequest{}
	bind(&req, r)
	return &req
}

func bindUntagResourceRequest(r *http.Request) *finding.UntagResourceRequest {
	req := finding.UntagResourceRequest{}
	bind(&req, r)
	return &req
}
