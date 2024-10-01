package rateLimiter

import (
	"log"
	"rate-limiter/m/v2/utils"
	"sync"
	"time"
)

type BucketData struct {
	previousTokenGrantTime time.Time
	numTokens              uint
}

type TokenBucketRateLimiter struct {
	maxAllowedRequestsPerMinute uint
	buckets                     map[string]BucketData
	timeSource                  utils.TimeSource
	mutex                       *sync.RWMutex
}

func NewTokenBucketRateLimiter(maxAllowedRequestsPerMinute uint, timeSource utils.TimeSource) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		maxAllowedRequestsPerMinute: maxAllowedRequestsPerMinute,
		buckets:                     make(map[string]BucketData),
		timeSource:                  timeSource,
		mutex:                       &sync.RWMutex{},
	}
}

func (receiver *TokenBucketRateLimiter) AllowRequest(requestInfo RequestInfo) (bool, error) {
	// FIXME: use a separate mutex for each bucket. Otherwise, the entire system is blocked when checking if a request
	// should be allowed. Instead, only requests for the same bucket should block
	defer receiver.mutex.Unlock()
	receiver.mutex.Lock()

	bucketKey := generateRequestKey(requestInfo)
	bucketData, exists := receiver.buckets[bucketKey]

	if !exists {
		bucketData = BucketData{
			numTokens:              receiver.maxAllowedRequestsPerMinute,
			previousTokenGrantTime: receiver.timeSource.Now(),
		}

		receiver.buckets[bucketKey] = bucketData
	}

	secondsSincePreviousTokenGrant := uint(receiver.timeSource.Now().Sub(bucketData.previousTokenGrantTime).Seconds())
	numTokensToGrant := uint((float64(secondsSincePreviousTokenGrant) / float64(60)) * float64(receiver.maxAllowedRequestsPerMinute))

	if numTokensToGrant > 0 {
		log.Printf("Granting %d tokens\n", numTokensToGrant)
		bucketData.previousTokenGrantTime = receiver.timeSource.Now()
	}

	bucketData.numTokens = min(receiver.maxAllowedRequestsPerMinute, bucketData.numTokens+numTokensToGrant)

	if bucketData.numTokens < 1 {
		return false, nil
	}

	bucketData.numTokens -= 1
	receiver.buckets[bucketKey] = bucketData

	return true, nil
}

func (receiver *TokenBucketRateLimiter) getNumberOfTokensRemaining(info RequestInfo) uint {
	bucketData, ok := receiver.buckets[generateRequestKey(info)]

	if !ok {
		return 0
	}

	return bucketData.numTokens
}
