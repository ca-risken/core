package main

import (
	"net/http"
)

func (g *gatewayService) listFindingHandler(w http.ResponseWriter, r *http.Request) {
	req := bindListFindingRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := g.findingClient.ListFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{"data": resp})
}

func (g *gatewayService) getFindingHandler(w http.ResponseWriter, r *http.Request) {
	req := bindGetFindingRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := g.findingClient.GetFinding(r.Context(), req)
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
	req := bindPutFindingRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := g.findingClient.PutFinding(r.Context(), req)
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
	req := bindDeleteFindingRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := g.findingClient.DeleteFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{"data": resp})
}

func (g *gatewayService) listFindingTagHandler(w http.ResponseWriter, r *http.Request) {
	req := bindListFindingTagRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := g.findingClient.ListFindingTag(r.Context(), req)
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
	req := bindTagFindingRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := g.findingClient.TagFinding(r.Context(), req)
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
	req := bindUntagFindingRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := g.findingClient.UntagFinding(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{"data": resp})
}

func (g *gatewayService) listResourceHandler(w http.ResponseWriter, r *http.Request) {
	req := bindListResourceRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := g.findingClient.ListResource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{"data": resp})
}

func (g *gatewayService) getResourceHandler(w http.ResponseWriter, r *http.Request) {
	req := bindGetResourceRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := g.findingClient.GetResource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{"data": resp})
}

func (g *gatewayService) putResourceHandler(w http.ResponseWriter, r *http.Request) {
	if err := validatePostHeader(r); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	req := bindPutResourceRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := g.findingClient.PutResource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{"data": resp})
}

func (g *gatewayService) deleteResourceHandler(w http.ResponseWriter, r *http.Request) {
	if err := validatePostHeader(r); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	req := bindDeleteResourceRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := g.findingClient.DeleteResource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{"data": resp})
}

func (g *gatewayService) listResourceTagHandler(w http.ResponseWriter, r *http.Request) {
	req := bindListResourceTagRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := g.findingClient.ListResourceTag(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{"data": resp})
}

func (g *gatewayService) tagResourceHandler(w http.ResponseWriter, r *http.Request) {
	if err := validatePostHeader(r); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	req := bindTagResourceRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := g.findingClient.TagResource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{"data": resp})
}

func (g *gatewayService) untagResourceHandler(w http.ResponseWriter, r *http.Request) {
	if err := validatePostHeader(r); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	req := bindUntagResourceRequest(r)
	if err := req.Validate(); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}
	resp, err := g.findingClient.UntagResource(r.Context(), req)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}
	writeResponse(w, http.StatusOK, map[string]interface{}{"data": resp})
}
