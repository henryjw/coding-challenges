package rateLimiter

import (
	testUtils "rate-limiter/m/v2/test-utils"
	"testing"
)

func TestRejection(t *testing.T) {
	limiter := NewTokenBucketRateLimiter(1, testUtils.FakeTimeSource{})

	request := RequestInfo{
		IPAddress: "test_ip",
		Endpoint:  "/test_endpoint",
	}

	allowed, err := limiter.AllowRequest(request)

	if err != nil {
		t.Fatal(err)
	}

	if !allowed {
		t.Fatal("Expected first request to be allowed")
	}

	allowed, err = limiter.AllowRequest(request)

	if err != nil {
		t.Fatal(err)
	}

	if allowed {
		t.Error("Expected second request not to be allowed")
	}
}

func TestConcurrency(t *testing.T) {
	t.Error("not yet implemented")
}

func TestBucketBucketRefill(t *testing.T) {
	t.Error("not yet implemented")
}

func TestTokenBucketRateLimiter_AllowRequest(t *testing.T) {
	t.Error("not yet implemented")
}
