package cache

import "fmt"

type KeyNotFoundError struct {
	Key string
}

type KeyAlreadyExistsError struct {
	Key string
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("key not found: %s", e.Key)
}

func (e *KeyAlreadyExistsError) Error() string {
	return fmt.Sprintf("key already exists: %s", e.Key)
}

type EmptyKeyError struct{}

func (e *EmptyKeyError) Error() string {
	return "key must have a length greater than 0"
}
