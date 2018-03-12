package ctx_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"

	"github.com/cutedogspark/echo-custom-context"
	"github.com/stretchr/testify/assert"
)

func TestCustomCtx(t *testing.T) {
	tt := []struct {
		name         string
		givenHandler func(c echo.Context) error
		wantJSON     string
	}{
		{
			name: "200",
			givenHandler: func(c echo.Context) error {
				return c.(ctx.CustomCtx).Resp(http.StatusOK).Data("hello world").Do()
			},
			wantJSON: `{"apiVersion": "v1", "data": "helle world"}`,
		},
		{
			name: "400",
			givenHandler: func(c echo.Context) error {

				errCode := 40000001
				errMsg := "Error Title"
				errDate := ctx.NewErrors()
				errDate.Add("Error Message 1")
				errDate.Add("Error Message 2")

				return c.(ctx.CustomCtx).Resp(errCode).Error(fmt.Sprintf("%v", errMsg)).Code(errCode).Errors(errDate.Error()).Do()
			},
			wantJSON: `{"apiVersion":"1.0","error":{"code":40000001,"message":"Error Title","errors":["Error Message 1","Error Message 2"]}}`,
		},
		{
			name: "400 with error and errors",
			givenHandler: func(c echo.Context) error {
				errs := []interface{}{}
				errs = append(errs, struct {
					name string
				}{"peter"})
				return c.(ctx.CustomCtx).Resp(http.StatusOK).Error("this is error message").Errors(errs).Do()
			},
			wantJSON: `{"apiVersion":"v1","error":{"code":0, "message":"this is error message", "errors":[{"name":"peter"}]}}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					return next(ctx.CustomCtx{c})
				}
			}(tc.givenHandler)(c)

			t.Log(rec.Body.String())
			assert.JSONEq(t, tc.wantJSON, rec.Body.String())
		})
	}
}