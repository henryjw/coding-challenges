package cache

import (
	"errors"
	"math"
	"reflect"
	"testing"
	"time"
)

// TODO: figure out how to generate exhaustive test cases to have higher confidence that this works as expected

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

	err := cache.Set("", Data{
		Flags:     uint16(32),
		Value:     "hello",
		ByteCount: 5,
	})

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
	cache.Set("test", Data{
		Flags:     uint16(32),
		Value:     "hello",
		ByteCount: 5,
	})

	if cache.Size() != 1 {
		t.Fatalf("Incorrect cache size. Expected: %d, got: %d'\n", 1, cache.Size())
	}
}

func TestSetOverrideExistingValue(t *testing.T) {
	data1 := Data{
		Flags:     uint16(32),
		Value:     "hello",
		ByteCount: 5,
	}

	data2 := Data{
		Flags:     uint16(13),
		Value:     "hi",
		ByteCount: 2,
	}

	cache := New(10)
	cache.Set("test", data1)
	cache.Set("test", data2)

	val, _ := cache.Get("test")

	if !reflect.DeepEqual(val, data2) {
		t.Errorf("Unexpected value. Expected '%v', got '%v'\n", data2, val)
	}
}

func TestSetAndGet(t *testing.T) {
	data := Data{
		Flags:     uint16(32),
		Value:     "hello",
		ByteCount: 5,
	}
	cache := New(10)
	cache.Set("test", data)
	value, err := cache.Get("test")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(value, data) {
		t.Fatalf("Incorrect value. Expected '%s', got '%s'\n", "test", "val")
	}
}

func TestBasicDelete(t *testing.T) {
	cache := New(10)
	cache.Set("test", Data{
		Flags:     uint16(32),
		Value:     "hello",
		ByteCount: 5,
	})

	err := cache.Delete("test")

	if err != nil {
		t.Fatal(err)
	}

	if cache.Size() != 0 {
		t.Fatalf("Incorrect cache size. Expected: %d, got: %d\n", 0, cache.Size())
	}
}

func TestDeleteNoMatch(t *testing.T) {
	cache := New(10)
	cache.Set("key", Data{
		Flags:     uint16(32),
		Value:     "hello",
		ByteCount: 5,
	})

	deleteErr := cache.Delete("key1")

	if deleteErr != nil {
		t.Fatal(deleteErr)
	}

	if cache.Size() != 1 {
		t.Fatalf("Incorrect cache size. Expected: %d, got: %d\n", 1, cache.Size())
	}

}

func TestCapacity(t *testing.T) {
	cache := New(10)

	if cache.Capacity != 10 {
		t.Fatalf("Expected capacity to be %d, got capacity = %d\n", 10, cache.Capacity)
	}
}

func TestUnboundedCapacity(t *testing.T) {
	cache := New(0)
	if cache.Capacity != math.MaxInt {
		t.Fatalf("Expected cache capacity to be unbounded (%d), got capacity = %d\n", math.MaxInt, cache.Capacity)
	}

	cache = New(-1)
	if cache.Capacity != math.MaxInt {
		t.Fatalf("Expected cache capacity to be unbounded (%d), got capacity = %d\n", math.MaxInt, cache.Capacity)
	}

	cache = New(-1300304)
	if cache.Capacity != math.MaxInt {
		t.Fatalf("Expected cache capacity to be unbounded (%d), got capacity = %d\n", math.MaxInt, cache.Capacity)
	}
}

func TestEviction(t *testing.T) {
	data1 := Data{
		Flags:     uint16(32),
		Value:     "hello",
		ByteCount: 5,
	}
	data2 := Data{
		Flags:     uint16(13),
		Value:     "hi",
		ByteCount: 2,
	}
	data3 := Data{
		Flags:     uint16(9),
		Value:     "hey",
		ByteCount: 3,
	}
	data4 := Data{
		Flags:     uint16(999),
		Value:     "yes",
		ByteCount: 3,
	}

	cache := New(2)
	cache.Set("key1", data1)
	cache.Set("key2", data2)
	cache.Set("key3", data3) // key1 should be evicted after this operation
	cache.Set("key4", data4)

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

	if !reflect.DeepEqual(val3, data3) {
		t.Errorf("Unexpected value. Expected '%v', got '%v'\n", data3, val3)
	}

	if !reflect.DeepEqual(val4, data4) {
		t.Errorf("Unexpected value. Expected '%v', got '%v'\n", data4, val4)
	}
}

func TestAdd(t *testing.T) {
	cache := New(-1)
	data := Data{
		Value:     "hello",
		ByteCount: 5,
		Flags:     uint16(1),
	}

	err := cache.Add("test", data)

	if err != nil {
		t.Fatal(err)
	}

	cachedData, err := cache.Get("test")

	if err != nil {
		t.Fatalf("Error getting stored value: %v\n", err)
	}

	if !reflect.DeepEqual(data, cachedData) {
		t.Errorf("Unexpected value. Expected '%v', got '%v'\n", data, cachedData)
	}
}

func TestAdd_KeyAlreadyExists(t *testing.T) {
	cache := New(-1)

	cache.Set("test", Data{})
	err := cache.Add("test", Data{})

	expectedErr := &KeyAlreadyExistsError{}

	if !errors.As(err, &expectedErr) {
		t.Fatalf("Unexpected error type. Expected %v, got %v\n", reflect.TypeOf(expectedErr), reflect.TypeOf(err))
	}
}

func TestReplace(t *testing.T) {
	cache := New(-1)
	data := Data{
		Value:     "hello",
		ByteCount: 5,
		Flags:     uint16(1),
	}

	cache.Set("test", Data{})
	err := cache.Replace("test", data)

	if err != nil {
		t.Fatal(err)
	}

	cachedData, err := cache.Get("test")

	if err != nil {
		t.Fatalf("Error getting stored value: %v\n", err)
	}

	if !reflect.DeepEqual(data, cachedData) {
		t.Errorf("Unexpected value. Expected '%v', got '%v'\n", data, cachedData)
	}
}

func TestReplace_KeyDoesntExist(t *testing.T) {
	cache := New(-1)

	err := cache.Replace("test", Data{})

	expectedErr := &KeyNotFoundError{}

	if !errors.As(err, &expectedErr) {
		t.Fatalf("Unexpected error type. Expected %v, got %v\n", reflect.TypeOf(expectedErr), reflect.TypeOf(err))
	}
}

func TestAppend(t *testing.T) {
	cache := New(-1)

	cache.Set("test", Data{
		Value:     "hello",
		ByteCount: 5,
		Flags:     uint16(1),
	})

	err := cache.Append("test", Data{
		Value:     ", world!",
		ByteCount: 8,
		Flags:     uint16(2),
	})

	if err != nil {
		t.Fatal(err)
	}

	cachedData, err := cache.Get("test")

	if err != nil {
		t.Fatalf("Error getting stored value: %v\n", err)
	}

	expected := Data{
		Value:     "hello, world!",
		ByteCount: 13,
		Flags:     uint16(1),
	}

	if !reflect.DeepEqual(expected, cachedData) {
		t.Errorf("Unexpected value. Expected '%v', got '%v'\n", expected, cachedData)
	}
}

func TestAppend_KeyDoesntExist(t *testing.T) {
	cache := New(-1)

	err := cache.Append("test", Data{})

	expectedErr := &KeyNotFoundError{}

	if !errors.As(err, &expectedErr) {
		t.Fatalf("Unexpected error type. Expected %v, got %v\n", reflect.TypeOf(expectedErr), reflect.TypeOf(err))
	}
}

func TestPrepend(t *testing.T) {
	cache := New(-1)

	cache.Set("test", Data{
		Value:     "hi!",
		ByteCount: 3,
		Flags:     uint16(1),
	})

	err := cache.Prepend("test", Data{
		Value:     "hello ",
		ByteCount: 6,
		Flags:     uint16(2),
	})

	if err != nil {
		t.Fatal(err)
	}

	cachedData, err := cache.Get("test")

	if err != nil {
		t.Fatalf("Error getting stored value: %v\n", err)
	}

	expected := Data{
		Value:     "hello hi!",
		ByteCount: 9,
		Flags:     uint16(1),
	}

	if !reflect.DeepEqual(expected, cachedData) {
		t.Errorf("Unexpected value. Expected '%v', got '%v'\n", expected, cachedData)
	}
}

func TestPrepend_KeyDoesntExist(t *testing.T) {
	cache := New(-1)

	err := cache.Prepend("test", Data{})

	expectedErr := &KeyNotFoundError{}

	if !errors.As(err, &expectedErr) {
		t.Fatalf("Unexpected error type. Expected %v, got %v\n", reflect.TypeOf(expectedErr), reflect.TypeOf(err))
	}
}

func TestKeyExpiration_Get(t *testing.T) {
	cache := New(-1)
	key := "test"
	data := Data{
		Value:     "hello",
		ByteCount: 5,
		Flags:     uint16(0),
		ExpiresAt: time.UnixMilli(1),
	}

	err := cache.Set(key, data)

	if err != nil {
		t.Fatal(err)
	}

	_, err = cache.Get(key)

	if err == nil {
		t.Fatal("Expected error")
	}

	expectedErr := &KeyNotFoundError{}

	if !errors.As(err, &expectedErr) {
		t.Fatalf("Unexpected error type. Expected %v, got %v\n", reflect.TypeOf(expectedErr), reflect.TypeOf(err))
	}
}

func TestKeyExpiration_Add(t *testing.T) {
	cache := New(-1)

	cache.Set("test", Data{
		Value:     "hello",
		ByteCount: 5,
		Flags:     uint16(1),
		ExpiresAt: time.UnixMilli(1),
	})

	data := Data{
		Value:     "hi",
		ByteCount: 2,
		Flags:     uint16(2),
	}

	err := cache.Add("test", data)

	if err != nil {
		t.Fatal(err)
	}

	cachedData, err := cache.Get("test")

	if err != nil {
		t.Fatalf("Error getting stored value: %v\n", err)
	}

	if !reflect.DeepEqual(data, cachedData) {
		t.Errorf("Unexpected value. Expected '%v', got '%v'\n", data, cachedData)
	}
}
func TestKeyExpiration_Append(t *testing.T) {
	cache := New(-1)

	cache.Set("test", Data{
		Value:     "hello",
		ByteCount: 5,
		Flags:     uint16(1),
		ExpiresAt: time.UnixMilli(1),
	})

	err := cache.Append("test", Data{
		Value:     ", world!",
		ByteCount: 8,
		Flags:     uint16(2),
	})

	if err == nil {
		t.Fatal("Expected error")
	}

	expectedErr := &KeyNotFoundError{}

	if !errors.As(err, &expectedErr) {
		t.Fatalf("Unexpected error type. Expected %v, got %v\n", reflect.TypeOf(expectedErr), reflect.TypeOf(err))
	}
}

func TestKeyExpiration_Prepend(t *testing.T) {
	cache := New(-1)

	cache.Set("test", Data{
		Value:     "hello",
		ByteCount: 5,
		Flags:     uint16(1),
		ExpiresAt: time.UnixMilli(1),
	})

	err := cache.Prepend("test", Data{
		Value:     ", world!",
		ByteCount: 8,
		Flags:     uint16(2),
	})

	if err == nil {
		t.Fatal("Expected error")
	}

	expectedErr := &KeyNotFoundError{}

	if !errors.As(err, &expectedErr) {
		t.Fatalf("Unexpected error type. Expected %v, got %v\n", reflect.TypeOf(expectedErr), reflect.TypeOf(err))
	}
}

func TestKeyExpiration_Replace(t *testing.T) {
	cache := New(-1)
	data := Data{
		Value:     "hello",
		ByteCount: 5,
		Flags:     uint16(1),
		ExpiresAt: time.UnixMilli(1),
	}

	cache.Set("test", data)
	err := cache.Replace("test", data)

	if err == nil {
		t.Fatal("Expected error")
	}

	expectedErr := &KeyNotFoundError{}

	if !errors.As(err, &expectedErr) {
		t.Fatalf("Unexpected error type. Expected %v, got %v\n", reflect.TypeOf(expectedErr), reflect.TypeOf(err))
	}
}

func TestKeyExpiration_Delete(t *testing.T) {
	cache := New(10)
	cache.Set("test", Data{
		Flags:     uint16(32),
		Value:     "hello",
		ByteCount: 5,
		ExpiresAt: time.UnixMilli(1),
	})

	err := cache.Delete("test")

	if err != nil {
		t.Fatal(err)
	}

	if cache.Size() != 0 {
		t.Fatalf("Incorrect cache size. Expected: %d, got: %d\n", 0, cache.Size())
	}
}
