package main

import (
	"net/http"

	"github.com/CyberAgent/mimosa-core/proto/finding"
)

func (g *gatewayService) listFindingHandler(w http.ResponseWriter, r *http.Request) {
	req := mappingListFindingRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := finding.NewFindingServiceClient(g.findingSvcConn).ListFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{"data": resp})
}

func (g *gatewayService) getFindingHandler(w http.ResponseWriter, r *http.Request) {
	req := mappingGetFindingRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := finding.NewFindingServiceClient(g.findingSvcConn).GetFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{"data": resp})
}

func (g *gatewayService) putFindingHandler(w http.ResponseWriter, r *http.Request) {
	if err := validatePostHeader(r); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	req := mappingPutFindingRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := finding.NewFindingServiceClient(g.findingSvcConn).PutFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{"data": resp})
}

func (g *gatewayService) deleteFindingHandler(w http.ResponseWriter, r *http.Request) {
	if err := validatePostHeader(r); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	req := mappingDeleteFindingRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := finding.NewFindingServiceClient(g.findingSvcConn).DeleteFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{"data": resp})
}

func (g *gatewayService) listFindingTagHandler(w http.ResponseWriter, r *http.Request) {
	req := mappingListFindingTagRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := finding.NewFindingServiceClient(g.findingSvcConn).ListFindingTag(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{"data": resp})
}

func (g *gatewayService) tagFindingHandler(w http.ResponseWriter, r *http.Request) {
	if err := validatePostHeader(r); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	req := mappingTagFindingRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := finding.NewFindingServiceClient(g.findingSvcConn).TagFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{"data": resp})
}

func (g *gatewayService) untagFindingHandler(w http.ResponseWriter, r *http.Request) {
	if err := validatePostHeader(r); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	req := mappingUntagFindingRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := finding.NewFindingServiceClient(g.findingSvcConn).UntagFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{"data": resp})
}

// func (g *gatewayService) listResourceHandler(w http.ResponseWriter, r *http.Request) {
// 	req := mappingListResourceRequest(r)
// }

// func (g *gatewayService) getResourceHandler(w http.ResponseWriter, r *http.Request) {
// 	req := mappingGetResourceRequest(r)
// }

// func (g *gatewayService) putResourceHandler(w http.ResponseWriter, r *http.Request) {
// 	req := mappingPutResourceRequest(r)
// }

// func (g *gatewayService) deleteResourceHandler(w http.ResponseWriter, r *http.Request) {
// 	req := mappingDeleteResourceRequest(r)
// }

// func (g *gatewayService) listResourceTagHandler(w http.ResponseWriter, r *http.Request) {
// 	req := mappingListResourceTagRequest(r)
// }

// func (g *gatewayService) tagResourceHandler(w http.ResponseWriter, r *http.Request) {
// 	req := mappingTagResourceRequest(r)
// }

// func (g *gatewayService) untagResourceHandler(w http.ResponseWriter, r *http.Request) {
// 	req := mappingUntagResourceRequest(r)
// }
