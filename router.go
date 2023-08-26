package prenn

import (
	"encoding/json"
	"errors"
	"net/http"
)

func New() *Router {
	r := &Router{}
	r.Group = Group{
		router: r,
	}
	return r
}

type Router struct {
	Group
	routes []route
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, params, err := router.findRoute(r.Method, r.URL.Path)
	switch err {
	case errNotFound:
		http.NotFound(w, r)
	case errMethodNotAllowed:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	default:
		data := route.handle(&Context{
			Context:   r.Context(),
			Writer:    w,
			Request:   r,
			Params:    params,
			RoutePath: route.path,
		})

		switch data := data.(type) {
		case nil:
			w.WriteHeader(http.StatusNoContent)
		case error:
			http.Error(w, data.Error(), http.StatusInternalServerError)
		default:
			v, _ := json.Marshal(data)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(v)
		}
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
