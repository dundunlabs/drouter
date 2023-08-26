package prenn

import (
	"context"
	"net/http"
)

type Context struct {
	context.Context

	Writer    http.ResponseWriter
	Request   *http.Request
	Params    Params
	RoutePath string
}

func (ctx *Context) WithValue(key any, value any) *Context {
	ctx.Context = context.WithValue(ctx.Context, key, value)
	return ctx
}

func (ctx *Context) Param(param string) string {
	return ctx.Params[param]
}
