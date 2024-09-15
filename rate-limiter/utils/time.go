package utils

import "time"

type TimeSource interface {
	Now() time.Time
}

type RealTimeSource struct{}

func (receiver *RealTimeSource) Now() time.Time {
	return time.Now()
}

type FakeTimeSource struct {
	FixedTime time.Time
}

func (f *FakeTimeSource) Now() time.Time {
	return f.FixedTime
}

func (f *FakeTimeSource) SetTime(t time.Time) {
	f.FixedTime = t
}
