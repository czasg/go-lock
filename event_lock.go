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

func (e *EventLock) Wait() context.Context {
    ctx, cancel := context.WithCancel(e.ctx)
    e.mx.Lock()
    e.cancels = append(e.cancels, cancel)
    e.mx.Unlock()
    return ctx
}

func (e *EventLock) WaitTime(wait time.Time) context.Context {
    ctx, cancel := context.WithDeadline(e.ctx, wait)
    e.mx.Lock()
    e.cancels = append(e.cancels, cancel)
    e.mx.Unlock()
    return ctx
}

func (e *EventLock) Notify() {
    if len(e.cancels) < 1 {
        return
    }
    e.mx.Lock()
    defer e.mx.Unlock()
    for _, cancel := range e.cancels {
        cancel()
    }
    e.cancels = []context.CancelFunc{}
}
