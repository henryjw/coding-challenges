package cache

import (
	"log"
	"strconv"
	"testing"
	"time"
)

func TestActiveExpirationCleanup(t *testing.T) {
	cache := New(-1)

	defer cache.stopCleanupBackgroundTask()

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

	defer cache.stopCleanupBackgroundTask()

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

func TestCleanupStressTest(t *testing.T) {
	cache := New(-1)
	numEntries := 10_000_000

	log.Printf("Generating %d test cases...\n", numEntries)
	for i := range numEntries {
		cache.Set(strconv.Itoa(i), Data{
			ExpiresAt: time.Now(),
		})
	}
	log.Println("Done generating test cases")

	log.Println("Clearing expired entries...")
	startTime := time.Now()
	cache.clearExpiredData()
	runTimeMs := time.Now().UnixMilli() - startTime.UnixMilli()
	log.Printf("Done clearing expired entries in %dms\n", runTimeMs)

	// NOTE: The current runtime is between 5-7 seconds. See if you can get it under 1 second
	// Also, note that it takes significantly less time (2-3 seconds) to run a sweep when no records are deleted.
	// There might be some inefficiencies in the deletion of keys
	expectedMaxRuntimeMs := int64(10_000)

	if runTimeMs > expectedMaxRuntimeMs {
		t.Errorf("Expected cache clear to take less than %dms, took %dms\n", expectedMaxRuntimeMs, runTimeMs)
	}

	if cache.Size() != 0 {
		t.Errorf("Expected cache to be empty. Got size = %d\n", cache.Size())
	}
}
