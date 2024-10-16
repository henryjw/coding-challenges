package cache

import "testing"

func TestKeyNotFoundError(t *testing.T) {
	err := KeyNotFoundError{Key: "key1"}

	if err.Error() != "key not found: key1" {
		t.Errorf("Unexpected error message: '%s'\n", err.Error())
	}

	err = KeyNotFoundError{Key: "key2"}

	if err.Error() != "key not found: key2" {
		t.Errorf("Unexpected error message: '%s'\n", err.Error())
	}
}

func TestEmptyKeyError(t *testing.T) {
	err := EmptyKeyError{}

	if err.Error() != "key must have a length greater than 0" {
		t.Errorf("Unexpected error message: '%s'\n", err.Error())
	}
}
