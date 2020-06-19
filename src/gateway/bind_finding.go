package main

import (
	"net/http"

	"github.com/CyberAgent/mimosa-core/proto/finding"
)

func bindListFindingRequest(r *http.Request) *finding.ListFindingRequest {
	req := finding.ListFindingRequest{}
	if err := bindQuery(&req, r); err != nil {
		appLogger.Warnf("Invalid parmeter. err=%+v", err)
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
	req := finding.GetFindingRequest{}
	if err := bindQuery(&req, r); err != nil {
		appLogger.Warnf("Invalid parmeter. err=%+v", err)
	}
	return &req
}

func bindPutFindingRequest(r *http.Request) *finding.PutFindingRequest {
	req := finding.PutFindingRequest{}
	if err := bindBodyJSON(&req, r); err != nil {
		appLogger.Warnf("Invalid PutFindingRequest. err=%+v", err)
	}
	return &req
}

func bindDeleteFindingRequest(r *http.Request) *finding.DeleteFindingRequest {
	req := finding.DeleteFindingRequest{}
	if err := bindBodyJSON(&req, r); err != nil {
		appLogger.Warnf("Invalid DeleteFindingRequest. err=%+v", err)
	}
	return &req
}

func bindListFindingTagRequest(r *http.Request) *finding.ListFindingTagRequest {
	req := finding.ListFindingTagRequest{}
	if err := bindQuery(&req, r); err != nil {
		appLogger.Warnf("Invalid parmeter. err=%+v", err)
	}
	return &req
}

func bindTagFindingRequest(r *http.Request) *finding.TagFindingRequest {
	req := finding.TagFindingRequest{}
	if err := bindBodyJSON(&req, r); err != nil {
		appLogger.Warnf("Invalid FindingTagForUpsert. err=%+v", err)
	}
	return &req
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
	if len(req.ResourceName) > 0 {
		req.ResourceName = commaSeparator(req.ResourceName[0])
	}
	return &req
}

func bindGetResourceRequest(r *http.Request) *finding.GetResourceRequest {
	req := finding.GetResourceRequest{}
	if err := bindQuery(&req, r); err != nil {
		appLogger.Warnf("Invalid parmeter. err=%+v", err)
	}
	return &req
}

func bindPutResourceRequest(r *http.Request) *finding.PutResourceRequest {
	req := finding.PutResourceRequest{}
	if err := bindBodyJSON(&req, r); err != nil {
		appLogger.Warnf("Invalid ResourceForUpsert. err=%+v", err)
	}
	return &req
}

func bindDeleteResourceRequest(r *http.Request) *finding.DeleteResourceRequest {
	req := finding.DeleteResourceRequest{}
	if err := bindBodyJSON(&req, r); err != nil {
		appLogger.Warnf("Invalid DeleteResourceRequest. err=%+v", err)
	}
	return &req
}

func bindListResourceTagRequest(r *http.Request) *finding.ListResourceTagRequest {
	req := finding.ListResourceTagRequest{}
	if err := bindQuery(&req, r); err != nil {
		appLogger.Warnf("Invalid parmeter. err=%+v", err)
	}
	return &req
}

func bindTagResourceRequest(r *http.Request) *finding.TagResourceRequest {
	req := finding.TagResourceRequest{}
	if err := bindBodyJSON(&req, r); err != nil {
		appLogger.Warnf("Invalid ResourceTagForUpsert. err=%+v", err)
	}
	return &req
}

func bindUntagResourceRequest(r *http.Request) *finding.UntagResourceRequest {
	req := finding.UntagResourceRequest{}
	if err := bindBodyJSON(&req, r); err != nil {
		appLogger.Warnf("Invalid UntagResourceRequest. err=%+v", err)
	}
	return &req
}
