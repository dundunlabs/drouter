package prenn_test

import (
	"bytes"
	"encoding/json"
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

func errorHandler(ctx *prenn.Context) any {
	return errors.New("any error")
}

func exceptionHandler(ctx *prenn.Context) any {
	return prenn.ExceptionBadRequest
}

func panicHandler(ctx *prenn.Context) any {
	panic(prenn.ExceptionBadRequest.WithError(errors.New("panic")))
}

func bindingHandler(ctx *prenn.Context) any {
	type Body struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"min=18"`
	}
	var body Body
	ctx.MustBindBody(&body)
	return body
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
			g.GET("", errorHandler)
			g.POST("", bindingHandler)
			g.PUT(":id", exceptionHandler)
			g.PATCH(":id", panicHandler)
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

func TestInternalServerError(t *testing.T) {
	res := fetch(http.MethodGet, "/api/v1/users", nil)
	if res.StatusCode != http.StatusInternalServerError {
		t.Error("should return statusCode 500")
	}
	body, _ := io.ReadAll(res.Body)
	if got, want := string(body), "any error\n"; got != want {
		t.Errorf("should return body %q, got %q", want, got)
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

func TestException(t *testing.T) {
	res := fetch(http.MethodPut, "/api/v1/users/1", nil)
	if res.StatusCode != http.StatusBadRequest {
		t.Error("should return statusCode 400")
	}
	body, _ := io.ReadAll(res.Body)
	if got, want := string(body), "Bad Request\n"; got != want {
		t.Errorf("should return body %q, got %q", want, got)
	}
}

func TestPanic(t *testing.T) {
	res := fetch(http.MethodPatch, "/api/v1/users/1", nil)
	if res.StatusCode != http.StatusBadRequest {
		t.Error("should return statusCode 400")
	}
	body, _ := io.ReadAll(res.Body)
	if got, want := string(body), "panic\n"; got != want {
		t.Errorf("should return body %q, got %q", want, got)
	}
}

func TestBinding(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		data := map[string]any{
			"name": "John",
			"age":  20,
		}
		v, _ := json.Marshal(data)
		res := fetch(http.MethodPost, "/api/v1/users", bytes.NewBuffer(v))
		if res.StatusCode != http.StatusOK {
			t.Error("should return statusCode 200")
		}
		body, _ := io.ReadAll(res.Body)
		if got, want := string(body), "{\"name\":\"John\",\"age\":20}"; got != want {
			t.Errorf("should return body %q, got %q", want, got)
		}
	})

	t.Run("fail", func(t *testing.T) {
		data := map[string]any{
			"name": "John",
			"age":  17,
		}
		v, _ := json.Marshal(data)
		res := fetch(http.MethodPost, "/api/v1/users", bytes.NewBuffer(v))
		if res.StatusCode != http.StatusBadRequest {
			t.Error("should return statusCode 400")
		}
		body, _ := io.ReadAll(res.Body)
		if got, want := string(body), "Key: 'Body.Age' Error:Field validation for 'Age' failed on the 'min' tag\n"; got != want {
			t.Errorf("should return body %q, got %q", want, got)
		}
	})
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
