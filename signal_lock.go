package lock

import "context"

type SignalLock struct{
    signals map[int][]context.CancelFunc
}

func (s *SignalLock) Wait() {}

func (s *SignalLock) Notify() {}
