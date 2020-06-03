package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/CyberAgent/mimosa-core/proto/finding"
	"google.golang.org/grpc"
)

type gatewayService struct {
	findingSvcAddr string
	findingSvcConn *grpc.ClientConn
}

func newGatewayService(ctx context.Context, conf gatewayConf) (gatewayService, error) {
	svc := gatewayService{
		findingSvcAddr: conf.FindingSvcAddr,
	}
	if err := mustConnGRPC(ctx, &svc.findingSvcConn, svc.findingSvcAddr); err != nil {
		return svc, err
	}
	return svc, nil
}

func mustConnGRPC(ctx context.Context, conn **grpc.ClientConn, addr string) error {
	var err error
	*conn, err = grpc.DialContext(ctx, addr,
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Second*3),
	)
	if err != nil {
		return err
	}
	return nil
}

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

func writeResponse(w http.ResponseWriter, status int, body map[string]interface{}) {
	if body == nil {
		w.WriteHeader(status)
		return
	}
	buf, err := json.Marshal(body)
	if err != nil {
		appLogger.Errorf("Response body JSON marshal error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	w.Write(buf)
}
