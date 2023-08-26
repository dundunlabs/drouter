package prenn_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dundunlabs/prenn"
)

func BenchmarkRouterSimple(b *testing.B) {
	router := prenn.New()

	router.GET("", simpleHandler)
	router.GET("user/dimfeld", simpleHandler)

	req := httptest.NewRequest("GET", "/user/dimfeld", nil)

	benchRequest(b, router, req)
}

func BenchmarkRouterRoot(b *testing.B) {
	router := prenn.New()

	router.GET("", simpleHandler)
	router.GET("user/dimfeld", simpleHandler)

	req := httptest.NewRequest("GET", "/", nil)

	benchRequest(b, router, req)
}

func BenchmarkRouterParam(b *testing.B) {
	router := prenn.New()

	router.GET("", simpleHandler)
	router.GET("user/:name", simpleHandler)

	req := httptest.NewRequest("GET", "/user/dimfeld", nil)

	benchRequest(b, router, req)
}

func BenchmarkRouterLongParams(b *testing.B) {
	router := prenn.New()

	router.GET("", simpleHandler)
	router.GET("user/:name/:resource", simpleHandler)

	req := httptest.NewRequest("GET", "/user/aaaabbbbccccddddeeeeffff/asdfghjkl", nil)

	benchRequest(b, router, req)
}

func BenchmarkRouterFiveColon(b *testing.B) {
	router := prenn.New()

	router.GET("", simpleHandler)
	router.GET(":a/:b/:c/:d/:e", simpleHandler)

	req := httptest.NewRequest("GET", "/test/test/test/test/test", nil)

	benchRequest(b, router, req)
}

func benchRequest(b *testing.B, router http.Handler, r *http.Request) {
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, r)
	}
}
