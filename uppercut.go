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
	SyncCounters       []Counter
	BeforeSyncCounters []Counter
	AfterSyncCounters  []Counter
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

func (m *Uppercut) AddSyncCounters(c Counter) {
	m.SyncCounters = append(m.SyncCounters, c)
}

func (m *Uppercut) AddBeforeSyncCounters(c Counter) {
	m.BeforeSyncCounters = append(m.BeforeSyncCounters, c)
}

func (m *Uppercut) AddAfterSyncCounters(c Counter) {
	m.AfterSyncCounters = append(m.AfterSyncCounters, c)
}

func (m Uppercut) Handler(requestCtx *fasthttp.RequestCtx) {
	beforeC := append(m.Counters, m.BeforeCounters...)
	afterC := append(m.Counters, m.AfterCounters...)
	beforeS := append(m.SyncCounters, m.BeforeSyncCounters...)
	afterS := append(m.SyncCounters, m.AfterSyncCounters...)

	upperCut(beforeC, requestCtx)
	sycCut(beforeS, requestCtx)

	m.RequestHandler(requestCtx)

	upperCut(afterC, requestCtx)
	sycCut(afterS, requestCtx)
}

func sycCut(counters []Counter, requestCtx *fasthttp.RequestCtx) {
	for _, m := range counters {
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
