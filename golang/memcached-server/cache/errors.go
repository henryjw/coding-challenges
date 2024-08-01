package cache

import "fmt"

type KeyNotFoundError struct {
	Key string
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("key not found: %s", e.Key)
}

type EmptyKeyError struct{}

func (e *EmptyKeyError) Error() string {
	return "Key must have a length greater than 0"
}
