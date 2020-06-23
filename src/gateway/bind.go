package main

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/schema"
)

var (
	decoder = newDecoder()
)

func newDecoder() *schema.Decoder {
	d := schema.NewDecoder()
	d.IgnoreUnknownKeys(true)
	d.ZeroEmpty(true)
	d.SetAliasTag("json")
	d.RegisterConverter([]string{}, func(input string) reflect.Value {
		return reflect.ValueOf(stringSeparator(input, ','))
	})
	return d
}

// bind bindding request parameter
func bind(out interface{}, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := bindQuery(out, r); err != nil {
			appLogger.Warnf("Could not `bindQuery`, url=%s, err=%+v", r.URL.RequestURI(), err)
		}
		return
	case http.MethodPost, http.MethodPut, http.MethodDelete:
		if err := bindBodyJSON(out, r); err != nil {
			appLogger.Warnf("Could not `bindBodyJSON`, url=%s, err=%+v", r.URL.RequestURI(), err)
		}
		return
	default:
		appLogger.Warnf("Unexpected HTTP Method, method=%s", r.Method)
	}
	return
}

// bindQuery bindding query parameter
func bindQuery(out interface{}, r *http.Request) error {
	return decoder.Decode(out, r.URL.Query())
}

// bindBodyJSON bindding body parameter binding
func bindBodyJSON(out interface{}, r *http.Request) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(out)
}

func stringSeparator(input string, delimiter rune) []string {
	separated := []string{}
	for _, p := range strings.Split(input, string(delimiter)) {
		if p != "" {
			separated = append(separated, p)
		}
	}
	return separated
}
