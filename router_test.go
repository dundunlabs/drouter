package prenn_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dundunlabs/prenn"
)

var router = prenn.New()

func simpleHandler(ctx *prenn.Context) any {
	return nil
}

func paramsHandler(ctx *prenn.Context) any {
	return ctx.Params
}

func adminMiddleware(next prenn.Handler) prenn.Handler {
	return func(ctx *prenn.Context) any {
		return errors.New("403 Forbidden")
	}
}

func init() {
	router.GET("", simpleHandler)
	router.WithGroup("api/:version", func(g prenn.Group) {
		g.WithGroup("users", func(g prenn.Group) {
			g.POST("", paramsHandler)
			g.PUT(":id", paramsHandler)
			g.PATCH(":id", paramsHandler)
			g.DELETE(":id", paramsHandler)
		})
		g.GET("*any", paramsHandler)
	})
	router.WithGroup("admin/api/:version", func(g prenn.Group) {
		g.Use(adminMiddleware)

		g.POST("customers", simpleHandler)
	})
}

func TestRoot(t *testing.T) {
	res := fetch(http.MethodGet, "/", nil)
	if res.StatusCode != http.StatusNoContent {
		t.Error("should return statusCode 204")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	res := fetch(http.MethodPost, "/", nil)
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Error("should return statusCode 405")
	}
}

func TestNotFound(t *testing.T) {
	res := fetch(http.MethodGet, "/404", nil)
	if res.StatusCode != http.StatusNotFound {
		t.Error("should return statusCode 404")
	}
}

func TestDynamicRoutes(t *testing.T) {
	res := fetch(http.MethodDelete, "/api/v1/users/100", nil)
	body, _ := io.ReadAll(res.Body)
	if got, want := string(body), "{\"id\":\"100\",\"version\":\"v1\"}"; got != want {
		t.Errorf("should return body %q, got %q", want, got)
	}

	res = fetch(http.MethodGet, "/api/v1/foo/bar", nil)
	body, _ = io.ReadAll(res.Body)
	if got, want := string(body), "{\"any\":\"foo/bar\",\"version\":\"v1\"}"; got != want {
		t.Errorf("should return body %q, got %q", want, got)
	}
}

func TestMiddleware(t *testing.T) {
	res := fetch(http.MethodPost, "/admin/api/v1/customers", nil)
	body, _ := io.ReadAll(res.Body)
	if got, want := string(body), "403 Forbidden\n"; got != want {
		t.Errorf("should return body %q, got %q", want, got)
	}
}

func fetch(method string, path string, body io.Reader) *http.Response {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, body)
	router.ServeHTTP(w, r)
	return w.Result()
}
