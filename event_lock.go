package lock

import (
    "context"
    "sync"
    "time"
)

func NewEventLock(ctx context.Context) *EventLock {
    if ctx == nil {
        ctx = context.Background()
    }
    return &EventLock{
        ctx: ctx,
    }
}

type EventLock struct {
    mx      sync.Mutex
    ctx     context.Context
    cancels []context.CancelFunc
}

func (e *EventLock) Wait() <-chan struct{} {
    ctx, cancel := context.WithCancel(e.ctx)
    e.mx.Lock()
    e.cancels = append(e.cancels, cancel)
    e.mx.Unlock()
    return ctx.Done()
}

func (e *EventLock) WaitTime(wait time.Time) <-chan struct{} {
    ctx, cancel := context.WithDeadline(e.ctx, wait)
    e.mx.Lock()
    e.cancels = append(e.cancels, cancel)
    e.mx.Unlock()
    return ctx.Done()
}

func (e *EventLock) Notify() {
    if len(e.cancels) < 1 {
        return
    }
    for _, cancel := range e.cancels {
        cancel()
    }
    e.mx.Lock()
    e.cancels = []context.CancelFunc{}
    e.mx.Unlock()
}
