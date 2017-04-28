package uppercut

import (
	"context"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

type Uppercut struct {
	RequestHandler     fasthttp.RequestHandler
	Counters           []Counter
	BeforeCounters     []Counter
	AfterCounters      []Counter
	Hooks              []Hook
	BeforeHooks        []Hook
	AfterHooks         []Hook
}

func NewUppercut(h fasthttp.RequestHandler) *Uppercut {
	return &Uppercut{RequestHandler: h}
}

func (m *Uppercut) AddCounters(c Counter) {
	m.Counters = append(m.Counters, c)
}

func (m *Uppercut) AddBeforeCounters(c Counter) {
	m.BeforeCounters = append(m.BeforeCounters, c)
}

func (m *Uppercut) AddAfterCounters(c Counter) {
	m.AfterCounters = append(m.AfterCounters, c)
}

func (m *Uppercut) AddHooks(c Counter) {
	m.Hooks = append(m.Hooks, c)
}

func (m *Uppercut) AddBeforeHooks(c Counter) {
	m.BeforeHooks = append(m.BeforeHooks, c)
}

func (m *Uppercut) AddAfterHooks(c Counter) {
	m.AfterHooks = append(m.AfterHooks, c)
}

func (m Uppercut) Handler(requestCtx *fasthttp.RequestCtx) {
	beforeC := append(m.Counters, m.BeforeCounters...)
	afterC := append(m.Counters, m.AfterCounters...)
	beforeS := append(m.Hooks, m.BeforeHooks...)
	afterS := append(m.Hooks, m.AfterHooks...)

	upperCut(beforeC, requestCtx)
	hook(beforeS, requestCtx)

	m.RequestHandler(requestCtx)

	upperCut(afterC, requestCtx)
	hook(afterS, requestCtx)
}

func hook(hooks []Hook, requestCtx *fasthttp.RequestCtx) {
	for _, m := range hooks {
		m.Call(requestCtx)
	}
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
