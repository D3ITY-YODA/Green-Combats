package clock

import "time"

type Clock interface {
	Now() time.Time
}

type systemClock struct{}

func New() Clock {
	return systemClock{}
}

func (systemClock) Now() time.Time {
	return time.Now()
}

type Fixed struct {
	T time.Time
}

func NewFixed(t time.Time) *Fixed {
	return &Fixed{T: t}
}

func (f *Fixed) Now() time.Time {
	return f.T
}

func (f *Fixed) Advance(d time.Duration) {
	f.T = f.T.Add(d)
}
