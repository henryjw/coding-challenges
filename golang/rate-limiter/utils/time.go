package utils

import "time"

type TimeSource interface {
	Now() time.Time
}

type RealTimeSource struct{}

func (receiver RealTimeSource) Now() time.Time {
	return time.Now()
}
