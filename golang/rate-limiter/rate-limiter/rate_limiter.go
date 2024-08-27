package rateLimiter

import (
	"errors"
	"rate-limiter/m/v2/utils"
)

var InvalidRateLimitAlgorithmError = errors.New("invalid rate limit algorithm")

const RateLimitTokenBucket = "token_bucket"

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
	case RateLimitTokenBucket:
		return NewTokenBucketRateLimiter(config.MaxAllowedRequestsPerMinute, &utils.RealTimeSource{}), nil
	default:
		return nil, InvalidRateLimitAlgorithmError
	}
}
