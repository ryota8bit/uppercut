package uppercut

import "github.com/valyala/fasthttp"

type Counter interface {
	Call(ctx *fasthttp.RequestCtx)
}

type CounterFunc func(ctx *fasthttp.RequestCtx)

func (m CounterFunc) Call(ctx *fasthttp.RequestCtx) {
	m(ctx)
}
