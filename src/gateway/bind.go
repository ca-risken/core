package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/schema"
)

var (
	errInvalidContentType = errors.New("request: invalid Content-Type")
	decoder               = newDecoder()
)

func newDecoder() *schema.Decoder {
	d := schema.NewDecoder()
	d.IgnoreUnknownKeys(true)
	d.ZeroEmpty(true)
	d.SetAliasTag("json")
	return d
}

// Bind request parameter binding
func Bind(out interface{}, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := BindQuery(out, r); err != nil {
			appLogger.Warnf("Could not `bindQuery`, url=%s, err=%+v", r.URL.RequestURI(), err)
		}
		return
	case http.MethodPost, http.MethodPut, http.MethodDelete:
		if err := BindBodyJSON(out, r); err != nil {
			appLogger.Warnf("Could not `bindBodyJSON`, url=%s, err=%+v", r.URL.RequestURI(), err)
		}
		return
	default:
		appLogger.Warnf("Unexpected HTTP Method, method=%s", r.Method)
	}
	return
}

// BindQuery query parameter binding
func BindQuery(out interface{}, r *http.Request) error {
	return decoder.Decode(out, r.URL.Query())
}

// BindBodyJSON body parameter binding
func BindBodyJSON(out interface{}, r *http.Request) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(out)
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

func parseUint64(str string) uint64 {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}
	return i
}
