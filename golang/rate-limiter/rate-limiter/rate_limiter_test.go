package rateLimiter

import (
	"reflect"
	"testing"
)

func TestNewTokenBucketRateLimiter(t *testing.T) {
	instance, err := New(Config{
		Algorithm:                   AlgorithmTokenBucket,
		MaxAllowedRequestsPerMinute: uint(5),
	})

	if err != nil {
		t.Fatal(err)
	}

	assertTypes(instance, reflect.TypeOf(&TokenBucketRateLimiter{}), t)
}

func TestNewFixedWindowCounterRateLimiter(t *testing.T) {
	instance, err := New(Config{
		Algorithm:                   AlgorithmFixedWindowCounter,
		MaxAllowedRequestsPerMinute: uint(5),
	})

	if err != nil {
		t.Fatal(err)
	}

	assertTypes(instance, reflect.TypeOf(&FixedWindowCounterRateLimiter{}), t)
}

func TestInvalidRateLimitAlgorithm(t *testing.T) {
	instance, err := New(Config{
		Algorithm:                   "",
		MaxAllowedRequestsPerMinute: uint(1),
	})

	if err == nil {
		t.Error("Expected error")
	}

	if instance != nil {
		t.Error("Instance should be nil")
	}
}

func assertTypes(instance any, expectedType reflect.Type, t *testing.T) {
	instanceType := reflect.TypeOf(instance)
	if reflect.TypeOf(instance) != expectedType {
		t.Errorf("Unexpected instance type. Expected %v, got %v\n", expectedType, instanceType)
	}
}
