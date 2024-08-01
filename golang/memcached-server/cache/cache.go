package cache

import (
	"container/list"
)

type KeyValue struct {
	Key   string
	Value string
}

type Cache struct {
	// Least frequently used elements are in the fronts
	accessList  *list.List
	lookupTable map[string]*list.Element
	Capacity    int
}

func New(capacity int) *Cache {
	return &Cache{
		accessList:  list.New(),
		lookupTable: make(map[string]*list.Element),
		Capacity:    capacity,
	}
}

func (receiver *Cache) Size() int {
	return len(receiver.lookupTable)
}

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
		prevKey := leastRecentlyUsedElement.Value.(*KeyValue).Key

		delete(receiver.lookupTable, prevKey)

		leastRecentlyUsedElement.Value.(*KeyValue).Key = key
		leastRecentlyUsedElement.Value.(*KeyValue).Value = value
		receiver.accessList.MoveToBack(leastRecentlyUsedElement)

		receiver.lookupTable[key] = leastRecentlyUsedElement
		return nil
	}

	element := receiver.accessList.PushBack(&KeyValue{Key: key, Value: value})
	receiver.lookupTable[key] = element

	return nil
}

func (receiver *Cache) Get(key string) (string, error) {
	if len(key) < 1 {
		return "", &EmptyKeyError{}
	}

	if _, exists := receiver.lookupTable[key]; !exists {
		return "", &KeyNotFoundError{}
	}

	element := receiver.lookupTable[key]
	value := element.Value.(*KeyValue).Value

	receiver.accessList.MoveToBack(element)

	return value, nil

}

func (receiver *Cache) Delete(key string) error {
	if _, exists := receiver.lookupTable[key]; !exists {
		return nil
	}

	element, _ := receiver.lookupTable[key]

	receiver.accessList.Remove(element)
	delete(receiver.lookupTable, key)

	return nil
}
