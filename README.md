# UpperCut!!
middleware chain for [fasthttprouter](https://github.com/buaazp/fasthttprouter), middleware runs in goroutine

# usage
```golang
counter := uppercut.NewUppercut(router.Handler)
counter.AddCounters(loggerMiddleware)
counter.AddBeforeCounters(panicHandler)
counter.AddAfterCounters(recoverHandler)
fasthttp.ListenAndServe(":8080", counter.Handler)
```

# Counter interface sample
```golang
func (l *LoggingHandler) Call(ctx *fasthttp.RequestCtx) {
	l.logger.Log("request: ", fmt.Sprintf("Host: %s, Path: %s, Method: %s", ctx.Host(), ctx.Path(), ctx.Method()))
}
```
