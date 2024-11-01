package cache

import (
	"container/list"
	"log"
	"math"
	"time"
)

type keyValue struct {
	Key   string
	Value Data
}

type Data struct {
	Value     string
	Flags     uint16
	ByteCount int
	// Time at which the data will expire. Defaults to never expire (`time.UnixMilli(0)`)
	ExpiresAt time.Time
}

// Cache simple in-memory cache
type Cache struct {
	// Least frequently used elements are in the fronts
	accessList           *list.List
	lookupTable          map[string]*list.Element
	shouldRunCleanupTask bool
	Capacity             int
}

// New Creates new Cache instance with a given capacity. Capacity will be unbounded if `capacity <= 0`
func New(capacity int) *Cache {
	if capacity <= 0 {
		capacity = math.MaxInt
	}

	return &Cache{
		accessList:           list.New(),
		lookupTable:          make(map[string]*list.Element),
		Capacity:             capacity,
		shouldRunCleanupTask: false,
	}
}

// Size Returns the number of keys currently in the cache
func (receiver *Cache) Size() int {
	return len(receiver.lookupTable)
}

// Set stores key with given value in the cache. Returns error if key is invalid (e.g., empty string)
func (receiver *Cache) Set(key string, data Data) error {
	if len(key) < 1 {
		return &EmptyKeyError{}
	}

	if element, exists := receiver.lookupTable[key]; exists {
		element.Value.(*keyValue).Value = data
		receiver.accessList.MoveToBack(element)

		return nil
	}

	if receiver.Size() == receiver.Capacity {
		leastRecentlyUsedElement := receiver.accessList.Front()
		prevKey := leastRecentlyUsedElement.Value.(*keyValue).Key

		delete(receiver.lookupTable, prevKey)

		leastRecentlyUsedElement.Value.(*keyValue).Key = key
		leastRecentlyUsedElement.Value.(*keyValue).Value = data
		receiver.accessList.MoveToBack(leastRecentlyUsedElement)

		receiver.lookupTable[key] = leastRecentlyUsedElement
		return nil
	}

	element := receiver.accessList.PushBack(&keyValue{Key: key, Value: data})
	receiver.lookupTable[key] = element

	return nil
}

// Get retrieves value from the cache by key. Returns error if key is not found or if the key is invalid (e.g., empty string)
func (receiver *Cache) Get(key string) (Data, error) {
	if len(key) < 1 {
		return Data{}, &EmptyKeyError{}
	}

	if _, exists := receiver.lookupTable[key]; !exists {
		return Data{}, &KeyNotFoundError{key}
	}

	element := receiver.lookupTable[key]
	value := element.Value.(*keyValue).Value

	if isExpired(value) {
		err := receiver.Delete(key)

		if err != nil {
			return Data{}, err
		}

		return Data{}, &KeyNotFoundError{key}
	}

	receiver.accessList.MoveToBack(element)

	return value, nil
}

// Delete key if it exists. Currently there are no errors for this function
func (receiver *Cache) Delete(key string) error {
	if _, exists := receiver.lookupTable[key]; !exists {
		return nil
	}

	element, _ := receiver.lookupTable[key]

	receiver.accessList.Remove(element)
	delete(receiver.lookupTable, key)

	return nil
}

func (receiver *Cache) Add(key string, data Data) error {
	if receiver.hasKey(key) && !receiver.isKeyExpired(key) {
		return &KeyAlreadyExistsError{Key: key}
	}

	return receiver.Set(key, data)
}

func (receiver *Cache) Replace(key string, data Data) error {
	if !receiver.hasKey(key) || receiver.isKeyExpired(key) {
		return &KeyNotFoundError{Key: key}
	}

	return receiver.Set(key, data)
}

// Append data is appended to the data matching the given key, if exists. Returns error if key doesn't exist
func (receiver *Cache) Append(key string, data Data) error {
	cachedData, err := receiver.Get(key)

	if err != nil {
		return err
	}

	return receiver.Set(key, Data{
		Value:     cachedData.Value + data.Value,
		ByteCount: cachedData.ByteCount + data.ByteCount,
		// There's no requirements in the project regarding the handling of these fields, so
		// just leave it as it is
		Flags:     cachedData.Flags,
		ExpiresAt: cachedData.ExpiresAt,
	})
}

// Prepend data is prepended to the data matching the given key, if exists. Returns error if key doesn't exist
func (receiver *Cache) Prepend(key string, data Data) error {
	cachedData, err := receiver.Get(key)

	if err != nil {
		return err
	}

	return receiver.Set(key, Data{
		Value:     data.Value + cachedData.Value,
		ByteCount: cachedData.ByteCount + data.ByteCount,
		// There's no requirements in the project regarding the handling of these fields, so
		// just leave it as it is
		Flags:     cachedData.Flags,
		ExpiresAt: cachedData.ExpiresAt,
	})
}

// RunExpireDataCleanupBackgroundTask starts background task to clean up expired data. No effect if there's already
// a task running for this cache instance. It's recommended to not set the frequency too low (less than 5 seconds) since it will negatively
// impact performance
func (receiver *Cache) RunExpireDataCleanupBackgroundTask(cleanupFrequencyMs int) {
	go receiver.setUpCleanupBackgroundTask(cleanupFrequencyMs)
}

// Close frees up resources and any running background tasks like Cache.RunExpireDataCleanupBackgroundTask
func (receiver *Cache) Close() {
	receiver.stopCleanupBackgroundTask()
}

func (receiver *Cache) setUpCleanupBackgroundTask(frequencyMs int) {
	receiver.shouldRunCleanupTask = true

	for receiver.shouldRunCleanupTask {
		log.Print("Deleting expired records...")
		receiver.clearExpiredData()
		time.Sleep(time.Millisecond * time.Duration(frequencyMs))
	}
}

func (receiver *Cache) clearExpiredData() {
	node := receiver.accessList.Front()
	sizeBefore := receiver.Size()

	// Iterating through the list is much faster (10-20x) than iterating through keys of the lookupTable.
	for node != nil {
		next := node.Next()
		val := node.Value.(*keyValue)
		if isExpired(val.Value) {
			receiver.Delete(val.Key)
		}

		node = next
	}

	numRecordsDeleted := sizeBefore - receiver.Size()

	if numRecordsDeleted > 0 {
		log.Printf("Deleted %d records\n", numRecordsDeleted)
	}
}

func (receiver *Cache) stopCleanupBackgroundTask() {
	receiver.shouldRunCleanupTask = false
	log.Println("Cleanup task stopped")
}

func (receiver *Cache) hasKey(key string) bool {
	_, ok := receiver.lookupTable[key]

	return ok
}

func (receiver *Cache) isKeyExpired(key string) bool {
	element, ok := receiver.lookupTable[key]

	if !ok {
		return false
	}

	return isExpired(element.Value.(*keyValue).Value)
}

func isExpired(data Data) bool {
	return data.ExpiresAt.UnixMilli() > 0 && time.Now().UnixMilli() > data.ExpiresAt.UnixMilli()
}
