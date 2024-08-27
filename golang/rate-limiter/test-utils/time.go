package testUtils

import "time"

type FakeTimeSource struct {
	fixedTime time.Time
}

func (f FakeTimeSource) Now() time.Time {
	return f.fixedTime
}
