package rateLimiter

import (
	"rate-limiter/m/v2/utils"
	"testing"
	"time"
)

func TestFixedWindowCounterRateLimiter_AllowRequest(t *testing.T) {
	limiter := NewFixedWindowCounterRateLimiter(1, &utils.FakeTimeSource{})

	request := RequestInfo{
		IPAddress: "test_ip",
		Endpoint:  "/endpoint",
	}

	allowed, err := limiter.AllowRequest(request)

	if err != nil {
		t.Fatal(err)
	}

	if !allowed {
		t.Fatal("Expected request to be allowed")
	}
}

func TestFixedWindowCounterRateLimiter_Reject(t *testing.T) {
	limiter := NewFixedWindowCounterRateLimiter(1, &utils.FakeTimeSource{})

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

func TestFixedWindowCounterRateLimiter_CounterReset(t *testing.T) {
	fakeTimeSource := &utils.FakeTimeSource{}
	numRequests := uint(5)
	limiter := NewFixedWindowCounterRateLimiter(5, fakeTimeSource)

	request := RequestInfo{
		IPAddress: "test_ip",
		Endpoint:  "/test_endpoint",
	}

	for range numRequests {
		limiter.AllowRequest(request)
	}

	if limiter.getNumberOfRequestsInWindow(request) != numRequests {
		t.Fatalf("Expected number of requests to be %d. Got %d\n", numRequests, limiter.getNumberOfRequestsInWindow(request))
	}

	fakeTimeSource.SetTime(fakeTimeSource.FixedTime.Add(1 * time.Minute))

	allow, err := limiter.AllowRequest(request)

	if err != nil {
		t.Fatal(err)
	}

	if !allow {
		t.Fatal("Expected request to be allowed")
	}

	if limiter.getNumberOfRequestsInWindow(request) != 1 {
		t.Fatalf("Expected number of requests in window to be 1; got %d\n", limiter.getNumberOfRequestsInWindow(request))
	}
}
