# UpperCut!!
middleware chain for [fasthttprouter](https://github.com/buaazp/fasthttprouter), middleware runs in goroutine

# usage
```golang
counters := []uppercut.Counter{loggingHandler}
beforeCounters := []uppercut.Counter{panicHandler}
afterCounters := []uppercut.Counter{recoverHandler}
m := uppercut.Counters{RequestHandler: fasthttprouter.Router, Counters: middlewares, , BeforeCounters: beforeCounters, , AfterCounters: afterCounters}
fasthttp.ListenAndServe(":8080", m.Handler)
```

# Counter interface sample
```golang
func (l *LoggingHandler) Call(ctx *fasthttp.RequestCtx) {
	l.logger.Log("request: ", fmt.Sprintf("Host: %s, Path: %s, Method: %s", ctx.Host(), ctx.Path(), ctx.Method()))
}
```
