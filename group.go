package prenn

import "net/http"

type Group struct {
	parent *Group
	router *Router
	path   string
	mdws   []Middleware
}

func (g *Group) WithGroup(path string, fn func(g Group)) {
	fn(Group{
		parent: g,
		router: g.router,
		path:   path,
	})
}

func (g *Group) Use(mdws ...Middleware) {
	g.mdws = append(g.mdws, mdws...)
}

func (g Group) GET(path string, handler Handler) {
	g.addRoute(http.MethodGet, path, handler)
}

func (g Group) POST(path string, handler Handler) {
	g.addRoute(http.MethodPost, path, handler)
}

func (g Group) PUT(path string, handler Handler) {
	g.addRoute(http.MethodPut, path, handler)
}

func (g Group) PATCH(path string, handler Handler) {
	g.addRoute(http.MethodPatch, path, handler)
}

func (g Group) DELETE(path string, handler Handler) {
	g.addRoute(http.MethodDelete, path, handler)
}

func (g Group) mergedMdws() []Middleware {
	if g.parent == nil {
		return g.mdws
	}
	return append(g.parent.mergedMdws(), g.mdws...)
}

func (g *Group) addRoute(method string, path string, handler Handler) {
	p := g.pathToGroup() + "/" + path
	r := route{
		method:      method,
		path:        p,
		pathMatcher: newPathMatcher(p),
		handler:     handler,
		group:       g,
	}
	g.router.addRoute(r)
}

func (g Group) pathToGroup() string {
	if g.parent == nil {
		return g.path
	}

	return g.parent.pathToGroup() + "/" + g.path
}
