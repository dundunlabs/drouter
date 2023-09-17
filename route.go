package drouter

type route struct {
	method      string
	path        string
	pathMatcher PathMacher
	handler     Handler
	group       *Group
}

func (r route) handle(ctx *Context) any {
	h := r.handler
	if r.group != nil {
		mdws := r.group.mergedMdws()
		mdwsLen := len(mdws)
		for i, _ := range mdws {
			h = mdws[mdwsLen-i-1](h)
		}
	}
	return h(ctx)
}
