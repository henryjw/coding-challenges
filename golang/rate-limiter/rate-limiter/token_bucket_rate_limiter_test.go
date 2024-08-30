package rateLimiter

import (
	"rate-limiter/m/v2/utils"
	"testing"
	"time"
)

func TestAllowRequest(t *testing.T) {
	limiter := NewTokenBucketRateLimiter(1, &utils.FakeTimeSource{})

	request := RequestInfo{
		IPAddress: "test_ip",
		Endpoint:  "/test_endpoint",
	}

	allowed, err := limiter.AllowRequest(request)

	if err != nil {
		t.Fatal(err)
	}

	if !allowed {
		t.Fatal("Expected request to be allowed")
	}
}

func TestRejection(t *testing.T) {
	limiter := NewTokenBucketRateLimiter(1, &utils.FakeTimeSource{})

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

func TestBucketTokenRefill_Basic(t *testing.T) {
	fakeTimer := &utils.FakeTimeSource{
		FixedTime: time.Unix(0, 0),
	}

	maxAllowedRequestsPerMinute := uint(60)
	numTokensPerSecond := maxAllowedRequestsPerMinute / 60
	limiter := NewTokenBucketRateLimiter(maxAllowedRequestsPerMinute, fakeTimer)

	request := RequestInfo{
		IPAddress: "test_ip",
		Endpoint:  "/test_endpoint",
	}

	for range maxAllowedRequestsPerMinute {
		limiter.AllowRequest(request)
	}

	if numTokens := limiter.getNumberOfTokensRemaining(request); numTokens != 0 {
		t.Fatalf("Expected bucket to be empty. Got %d tokens\n", numTokens)
	}

	fakeTimer.SetTime(fakeTimer.FixedTime.Add(5 * time.Second))

	// After this, the limiter should have generated 5 tokens and consumed 1.
	// So, there should be 4 tokens remaining
	if _, err := limiter.AllowRequest(request); err != nil {
		t.Error(err)
	}

	expectedNumTokens := (numTokensPerSecond * 5) - 1

	if numTokens := limiter.getNumberOfTokensRemaining(request); numTokens != expectedNumTokens {
		t.Errorf("Expected bucket to have %d tokens. Got %d\n", expectedNumTokens, numTokens)
	}
}

func TestBucketTokenRefill_Overflow(t *testing.T) {
	fakeTimer := &utils.FakeTimeSource{
		FixedTime: time.Unix(0, 0),
	}

	maxAllowedRequestsPerMinute := uint(60)
	limiter := NewTokenBucketRateLimiter(maxAllowedRequestsPerMinute, fakeTimer)

	request := RequestInfo{
		IPAddress: "test_ip",
		Endpoint:  "/test_endpoint",
	}

	for range maxAllowedRequestsPerMinute {
		limiter.AllowRequest(request)
	}

	if numTokens := limiter.getNumberOfTokensRemaining(request); numTokens != 0 {
		t.Fatalf("Expected bucket to be empty. Got %d tokens\n", numTokens)
	}

	fakeTimer.SetTime(fakeTimer.FixedTime.Add(9001 * time.Second))

	// After this, the limiter should have generated 60 tokens and consumed 1.
	if _, err := limiter.AllowRequest(request); err != nil {
		t.Error(err)
	}

	expectedNumTokens := maxAllowedRequestsPerMinute - 1

	if numTokens := limiter.getNumberOfTokensRemaining(request); numTokens != expectedNumTokens {
		t.Errorf("Expected bucket to have %d tokens. Got %d\n", expectedNumTokens, numTokens)
	}
}
