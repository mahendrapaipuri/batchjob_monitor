//go:build cgo
// +build cgo

package frontend

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	querierv1 "github.com/grafana/pyroscope/api/gen/proto/go/querier/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestParseTimeParam(t *testing.T) {
	type resultType struct {
		asTime  time.Time
		asError func() error
	}

	ts, err := parseTime("1582468023986")
	require.NoError(t, err)

	tests := []struct {
		paramName    string
		paramValue   string
		defaultValue time.Time
		result       resultType
	}{
		{ // When data is valid.
			paramName:    "start",
			paramValue:   "1582468023986",
			defaultValue: MinTime,
			result: resultType{
				asTime:  ts,
				asError: nil,
			},
		},
		{ // When data is empty string.
			paramName:    "end",
			paramValue:   "",
			defaultValue: MaxTime,
			result: resultType{
				asTime:  MaxTime,
				asError: nil,
			},
		},
		{ // When data is not valid.
			paramName:    "foo",
			paramValue:   "baz",
			defaultValue: MaxTime,
			result: resultType{
				asTime: time.Time{},
				asError: func() error {
					_, err := parseTime("baz")

					return fmt.Errorf("invalid time value for '%s': %w", "foo", err)
				},
			},
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest( //nolint:noctx
			http.MethodGet,
			"localhost:42/foo?"+test.paramName+"="+test.paramValue,
			nil,
		)
		require.NoError(t, err)

		result := test.result
		if asTime, err := parseTimeParam(req, test.paramName, test.defaultValue); err != nil {
			assert.Equal(t, err.Error(), result.asError().Error())
		} else {
			assert.Equal(t, result.asTime, asTime)
		}
	}
}

func TestParseTime(t *testing.T) {
	ts, err := time.Parse(time.RFC3339Nano, "2015-06-03T13:21:58.555Z")
	require.NoError(t, err)

	tests := []struct {
		input  string
		fail   bool
		result time.Time
	}{
		{
			input: "",
			fail:  true,
		},
		{
			input: "abc",
			fail:  true,
		},
		{
			input: "30s",
			fail:  true,
		},
		{
			input:  "123",
			result: time.Unix(123, 0),
		},
		{
			input:  "123.123",
			result: time.Unix(123, 123000000),
		},
		{
			input:  "2015-06-03T13:21:58.555Z",
			result: ts,
		},
		{
			input:  "2015-06-03T14:21:58.555+01:00",
			result: ts,
		},
		{
			// Test float rounding.
			input:  "1543578564.705",
			result: time.Unix(1543578564, 705*1e6),
		},
		{
			input:  MinTime.Format(time.RFC3339Nano),
			result: MinTime,
		},
		{
			input:  MaxTime.Format(time.RFC3339Nano),
			result: MaxTime,
		},
	}

	for _, test := range tests {
		ts, err := parseTime(test.input)
		if !test.fail {
			require.NoError(t, err)
			// assert.Equal(t, test.result, ts)
			if !ts.Equal(test.result) {
				t.Errorf("%s: expected %s, got %s", test.input, test.result, ts)
			}

			continue
		}

		assert.Error(t, err)
	}
}

func TestParseTSDBQueryParams(t *testing.T) {
	tests := []struct {
		query  string
		uuids  []string
		rmID   string
		rmIDs  []string
		method string
	}{
		{
			query:  "foo{uuid=~\"123|456\",gpuuuid=\"GPU-0123\",ceems_id=\"rm-0\"}",
			uuids:  []string{"123", "456"},
			rmID:   "rm-0",
			rmIDs:  []string{"rm-0", "rm-1"},
			method: "GET",
		},
		{
			query:  "foo{uuid=~\"abc-123|456\",ceems_id=\"rm-0|rm-1\"}",
			uuids:  []string{"abc-123", "456"},
			rmID:   "rm-1",
			rmIDs:  []string{"rm-0", "rm-1"},
			method: "POST",
		},
		{
			query:  "foo{uuid=\"456\",gpuuuid=\"GPU-0123\",ceems_id=\"rm-0\"}",
			uuids:  []string{"456"},
			rmID:   "rm-0",
			rmIDs:  []string{"rm-0"},
			method: "POST",
		},
		{
			query:  "foo{uuid=~\"abc_123|456\"}",
			method: "POST",
		},
	}

	for _, test := range tests {
		var body *strings.Reader

		// Query params
		data := url.Values{}
		data.Set("query", test.query)

		if test.method == "POST" {
			body = strings.NewReader(data.Encode())
		} else {
			body = strings.NewReader("hello")
		}

		req, err := http.NewRequest(test.method, "http://localhost:9090", body) //nolint:noctx
		require.NoError(t, err)

		// For GET request add query to URL
		if test.method == "GET" {
			req.URL.RawQuery = data.Encode()
		} else {
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		}

		p := &ReqParams{}
		err = parseTSDBRequest(p, req)
		require.NoError(t, err)

		assert.Equal(t, test.uuids, p.uuids)
		assert.Equal(t, test.rmID, p.clusterID)

		// Set parameters to request's context
		newReq := setQueryParams(req, p)

		if test.method == "POST" {
			// Check the new request body can still be parsed
			require.NoError(t, newReq.ParseForm())

			// Check if form value can be retrieved
			require.NotEmpty(t, newReq.FormValue("query"))
		}
	}
}

func TestParsePyroQueryParams(t *testing.T) {
	tests := []struct {
		message *querierv1.SelectMergeStacktracesRequest
		uuids   []string
		start   int64
		rmIDs   string
	}{
		{
			message: &querierv1.SelectMergeStacktracesRequest{
				LabelSelector: `{service_name="123"}`,
				Start:         1735209190,
			},
			uuids: []string{"123"},
			start: 1735209190000,
		},
		{
			message: &querierv1.SelectMergeStacktracesRequest{
				LabelSelector: `{service_name="123", ceems_id="default"}`,
				Start:         1735209190,
			},
			uuids: []string{"123"},
			rmIDs: "default",
			start: 1735209190000,
		},
	}

	for _, test := range tests {
		// Query params
		data, err := proto.Marshal(test.message)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "http://localhost:9090", bytes.NewBuffer(data)) //nolint:noctx
		require.NoError(t, err)

		p := &ReqParams{}
		err = parsePyroRequest(p, req)
		require.NoError(t, err)

		assert.Equal(t, test.uuids, p.uuids)
		assert.Equal(t, test.rmIDs, p.clusterID)
		assert.Equal(t, test.start, p.time)
	}
}
