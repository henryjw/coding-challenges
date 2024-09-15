package rateLimiter

import (
	"log"
	"rate-limiter/m/v2/utils"
	"time"
)

type fixedWindowInfo struct {
	windowStart      time.Time
	numberOfRequests uint
}

type FixedWindowCounterRateLimiter struct {
	maxAllowedRequestsPerMinute uint
	timeSource                  utils.TimeSource
	windows                     map[string]fixedWindowInfo
}

func NewFixedWindowCounterRateLimiter(maxAllowedRequestsPerMinute uint, timeSource utils.TimeSource) *FixedWindowCounterRateLimiter {
	return &FixedWindowCounterRateLimiter{
		maxAllowedRequestsPerMinute: maxAllowedRequestsPerMinute,
		timeSource:                  timeSource,
		windows:                     make(map[string]fixedWindowInfo),
	}
}

func (f *FixedWindowCounterRateLimiter) AllowRequest(requestInfo RequestInfo) (bool, error) {
	key := generateRequestKey(requestInfo)
	windowInfo, ok := f.windows[key]

	if !ok {
		windowInfo = fixedWindowInfo{
			windowStart:      f.timeSource.Now(),
			numberOfRequests: 0,
		}
		log.Printf("Initialized window for %s\n", key)
	}

	shouldResetCounter := f.timeSource.Now().Sub(windowInfo.windowStart).Seconds() >= 60

	if shouldResetCounter {
		log.Printf("Resetting counter for %s\n", key)
		windowInfo.windowStart = f.timeSource.Now()
		windowInfo.numberOfRequests = 0
	}

	if windowInfo.numberOfRequests >= f.maxAllowedRequestsPerMinute {
		return false, nil
	}

	windowInfo.numberOfRequests += 1

	f.windows[key] = windowInfo

	return true, nil
}

// Note that this doesn't account for when the counter should be reset.
// E.g., if the counter should be reset every minute and this function is called after, say, 5 minutes,
// it'll still return the last count instead of 0
func (f *FixedWindowCounterRateLimiter) getNumberOfRequestsInWindow(info RequestInfo) uint {
	key := generateRequestKey(info)

	return f.windows[key].numberOfRequests
}
