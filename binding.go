package prenn

import "encoding/json"

func BindBody[T any](ctx *Context) (T, error) {
	var v T
	if err := json.Unmarshal(ctx.Body(), &v); err != nil {
		return v, err
	}
	return v, ctx.validate.Struct(v)
}

func MustBindBody[T any](ctx *Context) T {
	v, err := BindBody[T](ctx)
	if err != nil {
		panic(ExceptionBadRequest.WithError(err))
	}
	return v
}
