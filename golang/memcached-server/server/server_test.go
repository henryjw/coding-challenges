package server

import (
	"memcached-server/cache"
	"memcached-server/utils"
	"testing"
)

func TestProcessSetCommand(t *testing.T) {
	server := New(cache.New(-1))
	result, err := server.processCommand(utils.Command{
		Name:      "set",
		Key:       "test_key",
		Noreply:   false,
		ByteCount: 5,
		ExpiresIn: 0,
		Flags:     uint16(5),
	}, "hello")

	if err != nil {
		t.Fatal(err)
	}

	if result != "STORED" {
		t.Errorf("Unexpected result: %s\n", result)
	}
}

func TestProcessAddCommand(t *testing.T) {
	server := New(cache.New(-1))
	result, err := server.processCommand(utils.Command{
		Name:      "add",
		Key:       "test_key",
		Noreply:   false,
		ByteCount: 5,
		ExpiresIn: 0,
		Flags:     uint16(5),
	}, "hello")

	if err != nil {
		t.Fatal(err)
	}

	if result != "STORED" {
		t.Errorf("Unexpected result: %s\n", result)
	}
}

func TestProcessAddCommand_KeyAlreadyExists(t *testing.T) {
	key := "test_key"
	c := cache.New(-1)
	server := New(c)

	err := c.Set(key, cache.Data{})

	if err != nil {
		t.Fatalf(err.Error())
	}

	result, err := server.processCommand(utils.Command{
		Name:      "add",
		Key:       key,
		Noreply:   false,
		ByteCount: 5,
		ExpiresIn: 0,
		Flags:     uint16(5),
	}, "hello")

	if err != nil {
		t.Fatal(err)
	}

	if result != "NOT_STORED" {
		t.Errorf("Unexpected result: %s\n", result)
	}
}

func TestProcessGetCommand(t *testing.T) {
	server := New(cache.New(-1))
	_, err := server.processCommand(utils.Command{
		Name:      "set",
		Key:       "test_key",
		Noreply:   false,
		ByteCount: 5,
		ExpiresIn: 0,
		Flags:     uint16(8),
	}, "hello")

	if err != nil {
		t.Fatal(err)
	}

	result, err := server.processCommand(utils.Command{
		Name:      "get",
		Key:       "test_key",
		Noreply:   false,
		ByteCount: 0,
		ExpiresIn: 0,
		Flags:     uint16(5),
	}, "")

	if result != "VALUE hello 8 5" {
		t.Errorf("Unexpected result: '%s'. Expected: '%s'\n", result, "hello")
	}
}

func TestGetCommandKeyDoesntExist(t *testing.T) {
	server := New(cache.New(-1))
	result, err := server.processCommand(utils.Command{
		Name:      "get",
		Key:       "test_key",
		Noreply:   false,
		ByteCount: 5,
		ExpiresIn: 0,
		Flags:     uint16(5),
	}, "hello")

	if err != nil {
		t.Fatalf("Unexpected error: %v\n", err)
	}

	if result != "END" {
		t.Fatalf("Expected result to be empty string. Got '%s'\n", result)
	}
}
