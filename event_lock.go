package lock

import (
    "context"
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
    ctx     context.Context
    cancels []context.CancelFunc
}

func (e *EventLock) Wait(wait ...time.Time) <-chan struct{} {
    if len(wait) < 1 {
        ctx, cancel := context.WithCancel(e.ctx)
        e.cancels = append(e.cancels, cancel)
        return ctx.Done()
    }
    ctx, cancel := context.WithDeadline(e.ctx, wait[0])
    e.cancels = append(e.cancels, cancel)
    return ctx.Done()
}

func (e *EventLock) Notify() {
    for _, cancel := range e.cancels {
        cancel()
    }
}
