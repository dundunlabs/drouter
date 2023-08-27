package prenn

import (
	"context"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Context struct {
	context.Context

	Writer    http.ResponseWriter
	Request   *http.Request
	Params    Params
	RoutePath string

	body     []byte
	validate *validator.Validate
}

func (ctx *Context) WithValue(key any, value any) *Context {
	ctx.Context = context.WithValue(ctx.Context, key, value)
	return ctx
}

func (ctx *Context) Param(param string) string {
	return ctx.Params[param]
}

func (ctx *Context) Body() []byte {
	if len(ctx.body) == 0 {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			panic(ExceptionBadRequest.WithError(err))
		}
		ctx.body = body
	}
	return ctx.body
}
