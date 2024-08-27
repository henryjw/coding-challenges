package rateLimiter

import (
	"reflect"
	"testing"
)

func TestNewTokenBucketRateLimiter(t *testing.T) {
	instance, err := New(Config{
		Algorithm:                   RateLimitTokenBucket,
		MaxAllowedRequestsPerMinute: uint(5),
	})

	if err != nil {
		t.Fatal(err)
	}

	instanceType := reflect.TypeOf(instance)
	expectedInstanceType := reflect.TypeOf(&TokenBucketRateLimiter{})

	if instanceType != expectedInstanceType {
		t.Errorf("Unexpected instance type. Expected %v, got %v\n", expectedInstanceType, instanceType)
	}
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
