package drouter

type Handler func(ctx *Context) any

type Middleware func(next Handler) Handler
