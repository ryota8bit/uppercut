package uppercut

import "github.com/valyala/fasthttp"

type Hook interface {
	Call(ctx *fasthttp.RequestCtx)
}

type HookFunc func(ctx *fasthttp.RequestCtx)

func (m HookFunc) Call(ctx *fasthttp.RequestCtx) {
	m(ctx)
}
