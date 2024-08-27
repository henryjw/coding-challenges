package rateLimiter

import (
	"fmt"
	"log"
	"rate-limiter/m/v2/utils"
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
}

func NewTokenBucketRateLimiter(maxAllowedRequestsPerMinute uint, timeSource utils.TimeSource) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		maxAllowedRequestsPerMinute: maxAllowedRequestsPerMinute,
		buckets:                     make(map[string]BucketData),
		timeSource:                  timeSource,
	}
}

func (receiver *TokenBucketRateLimiter) AllowRequest(requestInfo RequestInfo) (bool, error) {
	bucketKey := generateBucketKey(requestInfo)
	bucketData, exists := receiver.buckets[bucketKey]

	if !exists {
		bucketData = BucketData{
			numTokens:              receiver.maxAllowedRequestsPerMinute,
			previousTokenGrantTime: receiver.timeSource.Now(),
		}

		receiver.buckets[bucketKey] = bucketData
	}

	secondsSincePreviousTokenGrant := uint(receiver.timeSource.Now().Sub(bucketData.previousTokenGrantTime).Seconds())
	numTokensToGrant := (secondsSincePreviousTokenGrant / 60) * receiver.maxAllowedRequestsPerMinute

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

func generateBucketKey(info RequestInfo) string {
	return fmt.Sprintf("%s:%s", info.Endpoint, info.IPAddress)
}
