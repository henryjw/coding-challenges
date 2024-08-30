package rateLimiter

import (
	"errors"
	"fmt"
	"rate-limiter/m/v2/utils"
)

var InvalidRateLimitAlgorithmError = errors.New("invalid rate limit algorithm")

const AlgorithmTokenBucket = "token_bucket"
const AlgorithmFixedWindowCounter = "fixed_window_counter"

type RateLimiter struct {
	algorithm                   string
	maxAllowedRequestsPerSecond uint
}

type RequestInfo struct {
	IPAddress string
	Endpoint  string
}

type IRateLimiter interface {
	AllowRequest(requestInfo RequestInfo) (bool, error)
}

type Config struct {
	Algorithm                   string
	MaxAllowedRequestsPerMinute uint
}

func New(config Config) (IRateLimiter, error) {
	switch config.Algorithm {
	case AlgorithmTokenBucket:
		return NewTokenBucketRateLimiter(config.MaxAllowedRequestsPerMinute, &utils.RealTimeSource{}), nil
	case AlgorithmFixedWindowCounter:
		return NewFixedWindowCounterRateLimiter(config.MaxAllowedRequestsPerMinute, &utils.RealTimeSource{}), nil
	default:
		return nil, InvalidRateLimitAlgorithmError
	}
}

func generateRequestKey(info RequestInfo) string {
	return fmt.Sprintf("%s:%s", info.Endpoint, info.IPAddress)
}
