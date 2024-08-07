package server

import (
	"memcached-server/cache"
	"memcached-server/utils"
	"testing"
)

func TestProcessSetCommand(t *testing.T) {
	server := New(cache.New(-1))
	// Is validation required to ensure that the length of the data matches the byte_count field?
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
