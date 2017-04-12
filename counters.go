package uppercut

import (
	"context"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

type Uppercut struct {
	RequestHandler fasthttp.RequestHandler
	Counters       []Counter
	BeforeCounters []Counter
	AfterCounters  []Counter
}

func NewUppercut(h fasthttp.RequestHandler) *Uppercut {
	return &Uppercut{RequestHandler: h}
}

func (m Uppercut) AddCounters(c Counter) {
	m.Counters = append(m.Counters, c)
}

func (m Uppercut) AddBeforeCounters(c Counter) {
	m.BeforeCounters = append(m.BeforeCounters, c)
}

func (m Uppercut) AddAfterCounters(c Counter) {
	m.AfterCounters = append(m.AfterCounters, c)
}

func (m Uppercut) Handler(requestCtx *fasthttp.RequestCtx) {
	beforeM := append(m.Counters, m.BeforeCounters...)
	afterM := append(m.Counters, m.AfterCounters...)

	upperCut(beforeM, requestCtx)
	m.RequestHandler(requestCtx)
	upperCut(afterM, requestCtx)
}

type Counter interface {
	Call(ctx *fasthttp.RequestCtx)
}

type CounterFunc func(ctx *fasthttp.RequestCtx)

func (m CounterFunc) Call(ctx *fasthttp.RequestCtx) {
	m(ctx)
}

func upperCut(counters []Counter, requestCtx *fasthttp.RequestCtx) {
	wg := &sync.WaitGroup{}
	queue := make(chan Counter)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for _, m := range counters {
		wg.Add(1)
		go deQueue(ctx, wg, requestCtx, queue)
		enqueue(queue, m)
	}
	wg.Wait()
}

func deQueue(ctx context.Context, wg *sync.WaitGroup, requestCtx *fasthttp.RequestCtx, queue chan Counter) {
BREAK:
	for {
		select {
		case <-ctx.Done():
			break BREAK
		case excuter := <-queue:
			excuter.Call(requestCtx)
			wg.Done()
		}
	}
}

func enqueue(queue chan Counter, job Counter) {
	queue <- job
}
