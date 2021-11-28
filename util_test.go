package gtu_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ringsaturn/gtu"
)

func TestSimple(t *testing.T) {
	type args struct {
		t        *testing.T
		engine   *gin.Engine
		method   string
		url      string
		body     io.Reader
		validate gtu.ValidationFunc
		options  []gtu.RequestOption
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "basic",
			args: args{
				t: t,
				engine: func() *gin.Engine {
					e := gin.Default()
					e.GET("/ping", func(c *gin.Context) {
						c.String(http.StatusOK, "pong")
					})
					return e
				}(),
				method: gtu.GET,
				url:    "/ping",
				body:   nil,
				validate: func(t *testing.T, resp *httptest.ResponseRecorder) {
					if resp.Body.String() != "pong" {
						t.Error("not pong")
					}
				},
			},
		},
		{
			name: "validate HTTP Header",
			args: args{
				t: t,
				engine: func() *gin.Engine {
					e := gin.Default()
					e.GET("/check_header", func(c *gin.Context) {
						appName := c.Request.Header.Get("x-app-name")
						if appName == "" {
							c.String(http.StatusBadRequest, "not x-app-name")
							return
						}
						c.String(http.StatusOK, "ok")
					})
					return e
				}(),
				method: "GET",
				url:    "/check_header",
				body:   nil,
				validate: func(t *testing.T, resp *httptest.ResponseRecorder) {
					if resp.Result().StatusCode != http.StatusOK {
						t.Error("bad header")
					}
				},
				options: []gtu.RequestOption{
					gtu.HeaderOption(map[string]string{"x-app-name": "hello"}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gtu.Simple(tt.args.t, tt.args.engine, tt.args.method, tt.args.url, tt.args.body, tt.args.validate, tt.args.options...)
		})
	}
}
