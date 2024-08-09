package cache

import (
	"testing"
	"time"
)

func TestActiveExpirationCleanup(t *testing.T) {
	cache := New(-1)

	cache.Set("test", Data{
		Value:     "hello",
		ByteCount: 5,
		ExpiresAt: time.Now(),
	})

	cache.RunExpireDataCleanupBackgroundTask(100)

	// Give the task some time to run
	time.Sleep(time.Duration(time.Millisecond * 500))

	// Avoid calling any public cache methods since they will check for expiration of the data and clear it if it's expired
	_, ok := cache.lookupTable["test"]

	if ok {
		t.Error("Data should have been deleted")
	}
}

func TestActiveExpirationCleanupNoExpiredRecords(t *testing.T) {
	cache := New(-1)

	cache.Set("test", Data{
		Value:     "hello",
		ByteCount: 5,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(1)),
	})

	cache.RunExpireDataCleanupBackgroundTask(100)

	// Give the task some time to run
	time.Sleep(time.Duration(time.Millisecond * 500))

	// Avoid calling any public cache methods since they will check for expiration of the data and clear it if it's expired
	_, ok := cache.lookupTable["test"]

	if !ok {
		t.Error("Should not have been deleted")
	}
}

func TestActiveExpirationCleanupStops(t *testing.T) {
	cache := New(-1)

	cache.RunExpireDataCleanupBackgroundTask(100)

	time.Sleep(time.Duration(time.Millisecond * 500))

	cache.stopCleanupBackgroundTask()

	cache.Set("test", Data{
		Value:     "hello",
		ByteCount: 5,
		ExpiresAt: time.Now(),
	})

	// Avoid calling any public cache methods since they will check for expiration of the data and clear it if it's expired
	_, ok := cache.lookupTable["test"]

	if !ok {
		t.Error("Should not have been deleted")
	}
}
