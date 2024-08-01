package cache

import (
	"errors"
	"reflect"
	"testing"
)

func TestGetEmptyCache(t *testing.T) {
	cache := New(1)

	_, err := cache.Get("value")

	if err == nil {
		t.Fatal("Expected error")
	}

	target := &KeyNotFoundError{}
	if !errors.As(err, &target) {
		t.Fatalf("Unexpected error type. Expected %v, got %v\n", reflect.TypeOf(target), reflect.TypeOf(err))
	}
}

func TestGetEmptyKey(t *testing.T) {
	cache := New(1)

	_, err := cache.Get("")

	if err == nil {
		t.Fatal("Expected error")
	}

	target := &EmptyKeyError{}
	if !errors.As(err, &target) {
		t.Fatalf("Unexpected error type. Expected %v, got %v\n", reflect.TypeOf(target), reflect.TypeOf(err))
	}
}

func TestSetEmptyKey(t *testing.T) {
	cache := New(1)

	err := cache.Set("", "v")

	if err == nil {
		t.Fatal("Expected error")
	}

	target := &EmptyKeyError{}
	if !errors.As(err, &target) {
		t.Fatalf("Unexpected error type. Expected %v, got %v\n", reflect.TypeOf(target), reflect.TypeOf(err))
	}
}

func TestSet(t *testing.T) {
	cache := New(10)
	cache.Set("test", "val")

	if cache.Size() != 1 {
		t.Fatalf("Incorrect cache size. Expected: %d, got: %d'\n", 1, cache.Size())
	}
}

func TestSetAndGet(t *testing.T) {
	cache := New(10)
	cache.Set("test", "val")
	value, err := cache.Get("test")

	if err != nil {
		t.Fatal(err)
	}

	if value != "val" {
		t.Fatalf("Incorrect value. Expected '%s', got '%s'\n", "test", "val")
	}
}

func TestBasicDelete(t *testing.T) {
	cache := New(10)
	cache.Set("test", "val")

	err := cache.Delete("test")

	if err != nil {
		t.Fatal(err)
	}

	if cache.Size() != 0 {
		t.Fatalf("Incorrect cache size. Expected: %d, got: %d\n", 0, cache.Size())
	}
}

func TestEviction(t *testing.T) {
	cache := New(2)
	cache.Set("key1", "val1")
	cache.Set("key2", "val2")
	cache.Set("key3", "val3") // key1 should be evicted after this operation
	cache.Set("key4", "val4")

	if cache.Size() != 2 {
		t.Fatalf("Incorrect cache size. Expected: %d, got: %d\n", 2, cache.Size())
	}

	_, err1 := cache.Get("key1")
	_, err2 := cache.Get("key2")
	val3, _ := cache.Get("key3")
	val4, _ := cache.Get("key4")

	expectedErr := &KeyNotFoundError{}

	if !errors.As(err1, &expectedErr) {
		t.Errorf("Unexpected error type. Expected %v, got %v\n", reflect.TypeOf(expectedErr), reflect.TypeOf(err1))
	}

	if !errors.As(err2, &expectedErr) {
		t.Errorf("Unexpected error type. Expected %v, got %v\n", reflect.TypeOf(expectedErr), reflect.TypeOf(err1))
	}

	if val3 != "val3" {
		t.Errorf("Unexpected value. Expected '%s', got '%s'\n", "val3", val3)
	}

	if val4 != "val4" {
		t.Errorf("Unexpected value. Expected '%s', got '%s'\n", "val4", val4)
	}

}
