package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/schema"
	"github.com/vikyd/zero"
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

func bindQuery(out interface{}, r *http.Request) error {
	return decoder.Decode(out, r.URL.Query())
}

func bindBodyJSON(out proto.Message, r *http.Request) error {
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

func ignoreZeroValue(param []uint32) []uint32 {
	var result []uint32
	for _, v := range param {
		if zero.IsZeroVal(v) {
			continue
		}
		result = append(result, v)
	}
	return result
}

func parseUint64(str string) uint64 {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}
	return i
}
