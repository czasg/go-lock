package lock

type SignalLock struct{}

func (s *SignalLock) Wait() {}

func (s *SignalLock) Notify() {}
