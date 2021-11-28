package gtu

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

const GET = "GET"
const HEAD = "HEAD"
const POST = "POST"
const PUT = "PUT"
const PATCH = "PATCH"
const DELETE = "DELETE"

type RequestOption func(r *http.Request)

func HeaderOption(header map[string]string) RequestOption {
	return func(r *http.Request) {
		for headerKey, headerValue := range header {
			r.Header.Add(headerKey, headerValue)
		}
	}
}

func QueryOption(query map[string]string) RequestOption {
	return func(r *http.Request) {
		reqQuery := r.URL.Query()
		for headerKey, headerValue := range query {
			reqQuery.Add(headerKey, headerValue)
		}
		r.URL.RawQuery = reqQuery.Encode()
	}
}

type ValidationFunc func(t *testing.T, resp *httptest.ResponseRecorder)

func prepareRequest(method string, url string, body io.Reader, options ...RequestOption) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for _, option := range options {
		option(req)
	}
	return req, nil
}

func Simple(
	t *testing.T,
	engine *gin.Engine,
	method string, url string, body io.Reader,
	validate ValidationFunc,
	options ...RequestOption,
) {
	req, err := prepareRequest(method, url, body, options...)
	if err != nil {
		t.Error(err.Error())
	}
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	validate(t, w)
}

type ValidationFuncForJSONAPI func(t *testing.T, expectedResponse interface{})

func JSONAPI(
	t *testing.T,
	engine *gin.Engine,
	method string, url string, body io.Reader,
	expectedResponse interface{},
	validate ValidationFuncForJSONAPI,
	options ...RequestOption,
) {
	req, err := prepareRequest(method, url, body, options...)
	if err != nil {
		t.Error(err.Error())
	}
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if err := json.Unmarshal(w.Body.Bytes(), expectedResponse); err != nil {
		t.Error(err.Error())
	}

	validate(t, expectedResponse)
}
