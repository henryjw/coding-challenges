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
	mutex                  *sync.Mutex
}

type TokenBucketRateLimiter struct {
	maxAllowedRequestsPerMinute uint
	buckets                     map[string]*BucketData
	timeSource                  utils.TimeSource
	mutex                       *sync.RWMutex
}

func NewTokenBucketRateLimiter(maxAllowedRequestsPerMinute uint, timeSource utils.TimeSource) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		maxAllowedRequestsPerMinute: maxAllowedRequestsPerMinute,
		buckets:                     make(map[string]*BucketData),
		timeSource:                  timeSource,
		mutex:                       &sync.RWMutex{},
	}
}

func (receiver *TokenBucketRateLimiter) AllowRequest(requestInfo RequestInfo) (bool, error) {
	receiver.mutex.Lock()

	bucketKey := generateRequestKey(requestInfo)
	bucketData, exists := receiver.buckets[bucketKey]

	if !exists {
		bucketData = &BucketData{
			numTokens:              receiver.maxAllowedRequestsPerMinute,
			previousTokenGrantTime: receiver.timeSource.Now(),
			mutex:                  &sync.Mutex{},
		}
	}

	receiver.buckets[bucketKey] = bucketData

	receiver.mutex.Unlock()

	bucketData.mutex.Lock()
	defer bucketData.mutex.Unlock()

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

	return true, nil
}

func (receiver *TokenBucketRateLimiter) getNumberOfTokensRemaining(info RequestInfo) uint {
	bucketData, ok := receiver.buckets[generateRequestKey(info)]

	if !ok {
		return 0
	}

	return bucketData.numTokens
}
