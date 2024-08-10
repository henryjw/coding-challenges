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

func TestCleanupStressTest_AllItemsExpired(t *testing.T) {
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

	// NOTE: The current runtime is between 2-3 seconds (was 5-7 seconds). See if you can get it under 1 second
	// There might still be some inefficiencies in the deletion of keys
	expectedMaxRuntimeMs := int64(3000)

	if runTimeMs > expectedMaxRuntimeMs {
		t.Errorf("Expected cache clear to take less than %dms, took %dms\n", expectedMaxRuntimeMs, runTimeMs)
	}

	if cache.Size() != 0 {
		t.Errorf("Expected cache to be empty. Got size = %d\n", cache.Size())
	}

	if cache.accessList.Len() != 0 {
		t.Errorf("Expected access list to be empty. Got size %d\n", cache.accessList.Len())
	}
}

func TestCleanupStressTest_NoItemsExpired(t *testing.T) {
	cache := New(-1)
	numEntries := 10_000_000

	log.Printf("Generating %d test cases...\n", numEntries)
	for i := range numEntries {
		cache.Set(strconv.Itoa(i), Data{})
	}
	log.Println("Done generating test cases")

	log.Println("Clearing expired entries...")
	startTime := time.Now()
	cache.clearExpiredData()
	runTimeMs := time.Now().UnixMilli() - startTime.UnixMilli()
	log.Printf("Done clearing expired entries in %dms\n", runTimeMs)

	// NOTE: Runtime is currently 50ms-250ms for 10 million items. Pretty good 👍
	expectedMaxRuntimeMs := int64(300)

	if runTimeMs > expectedMaxRuntimeMs {
		t.Errorf("Expected cache clear to take less than %dms, took %dms\n", expectedMaxRuntimeMs, runTimeMs)
	}

	if cache.Size() != numEntries {
		t.Errorf("No entries should have been deleted. %d were deleted", numEntries-cache.Size())
	}

	if cache.accessList.Len() != numEntries {
		t.Errorf("No entries should have been removed from the access list. %d were deleted", numEntries-cache.accessList.Len())
	}
}
