package lock

import (
    "context"
    "sync"
    "time"
)

func NewSignalLock(ctx context.Context) *SignalLock {
    if ctx == nil {
        ctx = context.Background()
    }
    return &SignalLock{
        ctx:     ctx,
        signals: map[int][]context.CancelFunc{},
    }
}

type SignalLock struct {
    mx      sync.Mutex
    ctx     context.Context
    signals map[int][]context.CancelFunc
}

func (s *SignalLock) Wait(signal int) context.Context {
    s.mx.Lock()
    defer s.mx.Unlock()
    ctx, cancel := context.WithCancel(s.ctx)
    cancels, ok := s.signals[signal]
    if !ok {
        cancels = []context.CancelFunc{}
    }
    cancels = append(cancels, cancel)
    s.signals[signal] = cancels
    return ctx
}

func (s *SignalLock) WaitTime(signal int, wait time.Time) context.Context {
    s.mx.Lock()
    defer s.mx.Unlock()
    ctx, cancel := context.WithDeadline(s.ctx, wait)
    cancels, ok := s.signals[signal]
    if !ok {
        cancels = []context.CancelFunc{}
    }
    cancels = append(cancels, cancel)
    s.signals[signal] = cancels
    return ctx
}

func (s *SignalLock) Notify(signal int) {
    s.mx.Lock()
    defer s.mx.Unlock()
    cancels, ok := s.signals[signal]
    if !ok {
        return
    }
    for _, cancel := range cancels {
        cancel()
    }
    s.signals[signal] = []context.CancelFunc{}
}
