package prenn

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func New() *Router {
	r := &Router{
		validate: validator.New(),
	}
	r.Group = Group{
		router: r,
	}
	return r
}

type Router struct {
	Group
	routes   []route
	validate *validator.Validate
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, params, err := router.findRoute(r.Method, r.URL.Path)
	switch err {
	case errNotFound:
		http.NotFound(w, r)
	case errMethodNotAllowed:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	default:
		ctx := &Context{
			Context:   r.Context(),
			Writer:    w,
			Request:   r,
			Params:    params,
			RoutePath: route.path,
			validate:  router.validate,
		}
		defer func() {
			r := recover()
			handleResult(ctx, r)
		}()
		result := route.handle(ctx)
		handleResult(ctx, result)
	}
}

var (
	errNotFound         = errors.New("prenn: route not found")
	errMethodNotAllowed = errors.New("prenn: method not allowed")
)

func (router *Router) findRoute(method string, path string) (*route, Params, error) {
	matched := false
	for _, r := range router.routes {
		if params, ok := r.pathMatcher.Match(path); ok {
			matched = true
			if r.method == method {
				return &r, params, nil
			}
		}
	}
	if matched {
		return nil, nil, errMethodNotAllowed
	}
	return nil, nil, errNotFound
}

func (router *Router) addRoute(r route) {
	router.routes = append(router.routes, r)
}

func handleResult(ctx *Context, result any) {
	switch result := result.(type) {
	case nil:
		ctx.Writer.WriteHeader(http.StatusNoContent)
	case Exception:
		http.Error(ctx.Writer, result.Error(), result.statusCode)
	case error:
		http.Error(ctx.Writer, result.Error(), http.StatusInternalServerError)
	default:
		v, _ := json.Marshal(result)
		ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		ctx.Writer.Write(v)
	}
}
