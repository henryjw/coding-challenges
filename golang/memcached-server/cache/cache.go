package cache

import (
	"container/list"
	"math"
)

type keyValue struct {
	Key   string
	Value string
}

// Cache simple in-memory cache
type Cache struct {
	// Least frequently used elements are in the fronts
	accessList  *list.List
	lookupTable map[string]*list.Element
	Capacity    int
}

// New Creates new Cache instance with a given capacity. Capacity will be unbounded if `capacity <= 0`
func New(capacity int) *Cache {
	if capacity <= 0 {
		capacity = math.MaxInt
	}

	return &Cache{
		accessList:  list.New(),
		lookupTable: make(map[string]*list.Element),
		Capacity:    capacity,
	}
}

// Size Returns the number of keys currently in the cache
func (receiver *Cache) Size() int {
	return len(receiver.lookupTable)
}

// Set stores key with given value in the cache. Returns error if key is invalid (e.g., empty string)
func (receiver *Cache) Set(key string, value string) error {
	if len(key) < 1 {
		return &EmptyKeyError{}
	}

	if element, exists := receiver.lookupTable[key]; exists {
		element.Value = value
		receiver.accessList.MoveToBack(element)

		return nil
	}

	if receiver.Size() == receiver.Capacity {
		leastRecentlyUsedElement := receiver.accessList.Front()
		prevKey := leastRecentlyUsedElement.Value.(*keyValue).Key

		delete(receiver.lookupTable, prevKey)

		leastRecentlyUsedElement.Value.(*keyValue).Key = key
		leastRecentlyUsedElement.Value.(*keyValue).Value = value
		receiver.accessList.MoveToBack(leastRecentlyUsedElement)

		receiver.lookupTable[key] = leastRecentlyUsedElement
		return nil
	}

	element := receiver.accessList.PushBack(&keyValue{Key: key, Value: value})
	receiver.lookupTable[key] = element

	return nil
}

// Get retrieves value from the cache by key. Returns error if key is not found or if the key is invalid (e.g., empty string)
func (receiver *Cache) Get(key string) (string, error) {
	if len(key) < 1 {
		return "", &EmptyKeyError{}
	}

	if _, exists := receiver.lookupTable[key]; !exists {
		return "", &KeyNotFoundError{}
	}

	element := receiver.lookupTable[key]
	value := element.Value.(*keyValue).Value

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
