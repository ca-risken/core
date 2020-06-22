package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func TestListFindingHandler(t *testing.T) {
	findingMock := &mockFindingClient{}
	svc := gatewayService{
		findingClient: findingMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *finding.ListFindingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK Empty",
			input:      `project_id=1`,
			mockResp:   &finding.ListFindingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "OK Response exists",
			input:      `project_id=1001&data_source_name=rn`,
			mockResp:   &finding.ListFindingResponse{FindingId: []uint64{111, 222}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid request",
			input:      `data_source_name=rn`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Regist mock response
			if c.mockResp != nil || c.mockErr != nil {
				findingMock.On("ListFinding").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/finding/?"+c.input, nil)
			svc.listFindingHandler(rec, req)
			// Check Response
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestGetFindingHandler(t *testing.T) {
	findingMock := &mockFindingClient{}
	svc := gatewayService{
		findingClient: findingMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *finding.GetFindingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK Empty",
			input:      `project_id=1&finding_id=1001`,
			mockResp:   &finding.GetFindingResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "OK Response exists",
			input:      `project_id=1&finding_id=1002`,
			mockResp:   &finding.GetFindingResponse{Finding: &finding.Finding{FindingId: 1002, DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", Score: 0.5}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid request",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&finding_id=1003`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Regist mock response
			if c.mockResp != nil || c.mockErr != nil {
				findingMock.On("GetFinding").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/finding/detail?"+c.input, nil)
			svc.getFindingHandler(rec, req)
			// Check Response
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestPutFindingHandler(t *testing.T) {
	findingMock := &mockFindingClient{}
	svc := gatewayService{
		findingClient: findingMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *finding.PutFindingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1001, "finding":{"description":"desc", "data_source":"ds", "data_source_id":"ds-004", "resource_name":"rn", "project_id":1001, "original_score":55.51, "original_max_score":100.0, "data":"{\"key\":\"value\"}"}}`,
			mockResp:   &finding.PutFindingResponse{Finding: &finding.Finding{FindingId: 1001, Description: "desc", DataSource: "ds", DataSourceId: "ds-001", ResourceName: "rn", ProjectId: 1, Score: 0.5}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid request",
			input:      `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1, "finding":{"description":"desc", "data_source":"ds", "data_source_id":"ds-001", "resource_name":"rn", "project_id":1, "original_score":55.51, "original_max_score":100.0, "data":"{\"key\":\"value\"}"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Regist mock response
			if c.mockResp != nil || c.mockErr != nil {
				findingMock.On("PutFinding").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/finding/put", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putFindingHandler(rec, req)
			// Check Response
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestDeleteFindingHandler(t *testing.T) {
	findingMock := &mockFindingClient{}
	svc := gatewayService{
		findingClient: findingMock,
	}
	cases := []struct {
		name           string
		input          string
		existsMockResp bool
		mockErr        error
		wantStatus     int
	}{
		{
			name:           "OK",
			input:          `{"project_id":1, "finding_id":1005}`,
			existsMockResp: true,
			wantStatus:     http.StatusOK,
		},
		{
			name:       "NG Invalid request",
			input:      `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:           "NG Backend service error",
			input:          `{"project_id":1, "finding_id":9999}`,
			wantStatus:     http.StatusInternalServerError,
			existsMockResp: true,
			mockErr:        errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Regist mock response
			if c.existsMockResp {
				findingMock.On("DeleteFinding").Return(&empty.Empty{}, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/finding/delete", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteFindingHandler(rec, req)
			// Check Response
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestListFindingTagHandler(t *testing.T) {
	findingMock := &mockFindingClient{}
	svc := gatewayService{
		findingClient: findingMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *finding.ListFindingTagResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK Empty",
			input:      `project_id=1&finding_id=1001`,
			mockResp:   &finding.ListFindingTagResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:  "OK Response exists",
			input: `project_id=1&finding_id=1002`,
			mockResp: &finding.ListFindingTagResponse{Tag: []*finding.FindingTag{
				{FindingTagId: 1, FindingId: 1002, ProjectId: 1, TagKey: "k1", TagValue: "v"},
				{FindingTagId: 2, FindingId: 1002, ProjectId: 1, TagKey: "k2", TagValue: "v"},
			}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid request",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&finding_id=1003`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Regist mock response
			if c.mockResp != nil || c.mockErr != nil {
				findingMock.On("ListFindingTag").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/finding/tag?"+c.input, nil)
			svc.listFindingTagHandler(rec, req)
			// Check Response
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestTagFindingHandler(t *testing.T) {
	findingMock := &mockFindingClient{}
	svc := gatewayService{
		findingClient: findingMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *finding.TagFindingResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1001, "tag":{"finding_id":1001, "project_id":1001, "tag_key":"test", "tag_value":"true"}}`,
			mockResp:   &finding.TagFindingResponse{Tag: &finding.FindingTag{FindingTagId: 1001}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid request",
			input:      `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1001, "tag":{"finding_id":1001, "project_id":1001, "tag_key":"test", "tag_value":"true"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Regist mock response
			if c.mockResp != nil || c.mockErr != nil {
				findingMock.On("TagFinding").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/finding/tag/put", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.tagFindingHandler(rec, req)
			// Check Response
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestUntagFindingHandler(t *testing.T) {
	findingMock := &mockFindingClient{}
	svc := gatewayService{
		findingClient: findingMock,
	}
	cases := []struct {
		name           string
		input          string
		existsMockResp bool
		mockErr        error
		wantStatus     int
	}{
		{
			name:           "OK",
			input:          `{"project_id":1001, "finding_tag_id":1002}`,
			existsMockResp: true,
			wantStatus:     http.StatusOK,
		},
		{
			name:       "NG Invalid request",
			input:      `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:           "NG Backend service error",
			input:          `{"project_id":1001, "finding_tag_id":9999}`,
			wantStatus:     http.StatusInternalServerError,
			existsMockResp: true,
			mockErr:        errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Regist mock response
			if c.existsMockResp {
				findingMock.On("UntagFinding").Return(&empty.Empty{}, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/finding/tag/delete", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.untagFindingHandler(rec, req)
			// Check Response
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestListResourceHandler(t *testing.T) {
	findingMock := &mockFindingClient{}
	svc := gatewayService{
		findingClient: findingMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *finding.ListResourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK Empty",
			input:      `project_id=1&resource_name=rn,rn2`,
			mockResp:   &finding.ListResourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "OK Response exists",
			input:      `project_id=2&resource_name=rn,rn2`,
			mockResp:   &finding.ListResourceResponse{ResourceId: []uint64{111, 222}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid request",
			input:      ``,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=2&resource_name=rn,rn2`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Regist mock response
			if c.mockResp != nil || c.mockErr != nil {
				findingMock.On("ListResource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/resource/?"+c.input, nil)
			svc.listResourceHandler(rec, req)
			// Check Response
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestGetResourceHandler(t *testing.T) {
	findingMock := &mockFindingClient{}
	svc := gatewayService{
		findingClient: findingMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *finding.GetResourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK Empty",
			input:      `project_id=1001&resource_id=1001`,
			mockResp:   &finding.GetResourceResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "OK Response exists",
			input:      `project_id=1001&resource_id=1001`,
			mockResp:   &finding.GetResourceResponse{Resource: &finding.Resource{ResourceId: 1002, ResourceName: "rn"}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid request",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1001&resource_id=1001`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Regist mock response
			if c.mockResp != nil || c.mockErr != nil {
				findingMock.On("GetResource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/resource/detail?"+c.input, nil)
			svc.getResourceHandler(rec, req)
			// Check Response
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestPutResourceHandler(t *testing.T) {
	findingMock := &mockFindingClient{}
	svc := gatewayService{
		findingClient: findingMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *finding.PutResourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1001, "resource":{"resource_name":"rn", "project_id":1001}}`,
			mockResp:   &finding.PutResourceResponse{Resource: &finding.Resource{ResourceId: 1001, ResourceName: "rn", ProjectId: 1001}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid request",
			input:      `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1001, "resource":{"resource_name":"rn", "project_id":1001}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Regist mock response
			if c.mockResp != nil || c.mockErr != nil {
				findingMock.On("PutResource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/resource/put", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.putResourceHandler(rec, req)
			// Check Response
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestDeleteResourceHandler(t *testing.T) {
	findingMock := &mockFindingClient{}
	svc := gatewayService{
		findingClient: findingMock,
	}
	cases := []struct {
		name           string
		input          string
		existsMockResp bool
		mockErr        error
		wantStatus     int
	}{
		{
			name:           "OK",
			input:          `{"project_id":1001, "resource_id":1003}`,
			existsMockResp: true,
			wantStatus:     http.StatusOK,
		},
		{
			name:       "NG Invalid request",
			input:      `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:           "NG Backend service error",
			input:          `{"project_id":1001, "resource_id":1003}`,
			wantStatus:     http.StatusInternalServerError,
			existsMockResp: true,
			mockErr:        errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Regist mock response
			if c.existsMockResp {
				findingMock.On("DeleteResource").Return(&empty.Empty{}, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/resource/delete", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.deleteResourceHandler(rec, req)
			// Check Response
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestListResourceTagHandler(t *testing.T) {
	findingMock := &mockFindingClient{}
	svc := gatewayService{
		findingClient: findingMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *finding.ListResourceTagResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK Empty",
			input:      `project_id=1&resource_id=1001`,
			mockResp:   &finding.ListResourceTagResponse{},
			wantStatus: http.StatusOK,
		},
		{
			name:  "OK Response exists",
			input: `project_id=1&resource_id=1001`,
			mockResp: &finding.ListResourceTagResponse{Tag: []*finding.ResourceTag{
				{ResourceTagId: 1, ResourceId: 1001, ProjectId: 1, TagKey: "k1", TagValue: "v"},
				{ResourceTagId: 2, ResourceId: 1001, ProjectId: 1, TagKey: "k2", TagValue: "v"},
			}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid request",
			input:      `project_id=1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `project_id=1&resource_id=1001`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Regist mock response
			if c.mockResp != nil || c.mockErr != nil {
				findingMock.On("ListResourceTag").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/resource/tag?"+c.input, nil)
			svc.listResourceTagHandler(rec, req)
			// Check Response
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestTagResourceHandler(t *testing.T) {
	findingMock := &mockFindingClient{}
	svc := gatewayService{
		findingClient: findingMock,
	}
	cases := []struct {
		name       string
		input      string
		mockResp   *finding.TagResourceResponse
		mockErr    error
		wantStatus int
	}{
		{
			name:       "OK",
			input:      `{"project_id":1001, "tag":{"resource_id":1001, "project_id":1001, "tag_key":"test", "tag_value":"true"}}`,
			mockResp:   &finding.TagResourceResponse{Tag: &finding.ResourceTag{ResourceTagId: 1001}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NG Invalid request",
			input:      `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "NG Backend service error",
			input:      `{"project_id":1001, "tag":{"resource_id":1001, "project_id":1001, "tag_key":"test", "tag_value":"true"}}`,
			wantStatus: http.StatusInternalServerError,
			mockErr:    errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Regist mock response
			if c.mockResp != nil || c.mockErr != nil {
				findingMock.On("TagResource").Return(c.mockResp, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/resource/tag/put", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.tagResourceHandler(rec, req)
			// Check Response
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

func TestUntagResourceHandler(t *testing.T) {
	findingMock := &mockFindingClient{}
	svc := gatewayService{
		findingClient: findingMock,
	}
	cases := []struct {
		name           string
		input          string
		existsMockResp bool
		mockErr        error
		wantStatus     int
	}{
		{
			name:           "OK",
			input:          `{"project_id":1001, "resource_tag_id":1004}`,
			existsMockResp: true,
			wantStatus:     http.StatusOK,
		},
		{
			name:       "NG Invalid request",
			input:      `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:           "NG Backend service error",
			input:          `{"project_id":1001, "resource_tag_id":1004}`,
			wantStatus:     http.StatusInternalServerError,
			existsMockResp: true,
			mockErr:        errors.New("something wrong"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Regist mock response
			if c.existsMockResp {
				findingMock.On("UntagResource").Return(&empty.Empty{}, c.mockErr).Once()
			}
			// Invoke HTTP Request
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/resource/tag/delete", strings.NewReader(c.input))
			req.Header.Add("Content-Type", "application/json")
			svc.untagResourceHandler(rec, req)
			// Check Response
			if c.wantStatus != rec.Code {
				t.Fatalf("Unexpected HTTP status code: want=%+v, got=%+v", c.wantStatus, rec.Code)
			}
			resp := map[string]interface{}{}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("Unexpected json decode error to response body: err=%+v", err)
			}
			jsonKey := successJSONKey
			if c.wantStatus != http.StatusOK {
				jsonKey = errorJSONKey
			}
			if _, ok := resp[jsonKey]; !ok {
				t.Fatalf("Unexpected no response key: want key=%s", jsonKey)
			}
		})
	}
}

/**
 * Mock Client
**/
type mockFindingClient struct {
	mock.Mock
}

func (m *mockFindingClient) ListFinding(context.Context, *finding.ListFindingRequest, ...grpc.CallOption) (*finding.ListFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListFindingResponse), args.Error(1)
}
func (m *mockFindingClient) GetFinding(context.Context, *finding.GetFindingRequest, ...grpc.CallOption) (*finding.GetFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.GetFindingResponse), args.Error(1)
}
func (m *mockFindingClient) PutFinding(context.Context, *finding.PutFindingRequest, ...grpc.CallOption) (*finding.PutFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.PutFindingResponse), args.Error(1)
}
func (m *mockFindingClient) DeleteFinding(context.Context, *finding.DeleteFindingRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockFindingClient) ListFindingTag(context.Context, *finding.ListFindingTagRequest, ...grpc.CallOption) (*finding.ListFindingTagResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListFindingTagResponse), args.Error(1)
}
func (m *mockFindingClient) TagFinding(context.Context, *finding.TagFindingRequest, ...grpc.CallOption) (*finding.TagFindingResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.TagFindingResponse), args.Error(1)
}
func (m *mockFindingClient) UntagFinding(context.Context, *finding.UntagFindingRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockFindingClient) ListResource(context.Context, *finding.ListResourceRequest, ...grpc.CallOption) (*finding.ListResourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListResourceResponse), args.Error(1)
}
func (m *mockFindingClient) GetResource(context.Context, *finding.GetResourceRequest, ...grpc.CallOption) (*finding.GetResourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.GetResourceResponse), args.Error(1)
}
func (m *mockFindingClient) PutResource(context.Context, *finding.PutResourceRequest, ...grpc.CallOption) (*finding.PutResourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.PutResourceResponse), args.Error(1)
}
func (m *mockFindingClient) DeleteResource(context.Context, *finding.DeleteResourceRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
func (m *mockFindingClient) ListResourceTag(context.Context, *finding.ListResourceTagRequest, ...grpc.CallOption) (*finding.ListResourceTagResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.ListResourceTagResponse), args.Error(1)
}
func (m *mockFindingClient) TagResource(context.Context, *finding.TagResourceRequest, ...grpc.CallOption) (*finding.TagResourceResponse, error) {
	args := m.Called()
	return args.Get(0).(*finding.TagResourceResponse), args.Error(1)
}
func (m *mockFindingClient) UntagResource(context.Context, *finding.UntagResourceRequest, ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called()
	return args.Get(0).(*empty.Empty), args.Error(1)
}
