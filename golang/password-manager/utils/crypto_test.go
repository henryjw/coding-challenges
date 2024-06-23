package utils_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"password-manager/utils"
	"testing"
)

func TestHashIsDeterministic(t *testing.T) {
	value := "testvalue"
	salt := "testsalt"

	hash1 := utils.Hash(value, salt)
	hash2 := utils.Hash(value, salt)

	assert.Equal(t, hash1, hash2, "The generated hashes should be the same")
}

func TestEncryptionDeterministic(t *testing.T) {
	value := "test"
	// String must be 32 characters long
	key := "123e4567e89b12d3a456426614174000"

	cipherValue1, encryptError1 := utils.Encrypt(value, key)
	assert.Nil(t, encryptError1)

	cipherValue2, encryptError2 := utils.Encrypt(value, key)
	assert.Nil(t, encryptError2)

	assert.Equal(t, cipherValue1, cipherValue2)
}

func TestSymmetricEncryption(t *testing.T) {
	value := "test"
	// String must be 32 characters long
	key := "123e4567e89b12d3a456426614174000"

	cipherValue, encryptError := utils.Encrypt(value, key)

	if !assert.Nil(t, encryptError) {
		fmt.Printf("%v\n", encryptError)
		return
	}

	decryptedValue, decryptError := utils.Decrypt(cipherValue, key)

	if !assert.Nil(t, decryptError) {
		return
	}

	assert.Equal(t, value, decryptedValue)
}

func TestEncryptTruncatesLongKey(t *testing.T) {
	value := "test"
	longKey := "123e4567e89b12d3a456426614174000aaaaaaaaaaaaaaaaaaaaaaaaa"
	correctKey := "123e4567e89b12d3a456426614174000"

	cipherValueLongKey, errLongKey := utils.Encrypt(value, longKey)
	assert.Nil(t, errLongKey)

	cipherValueCorrectKey, errCorrectKey := utils.Encrypt(value, correctKey)
	assert.Nil(t, errCorrectKey)

	assert.Equal(t, cipherValueLongKey, cipherValueCorrectKey)
}

func TestDecryptTruncatesLongKey(t *testing.T) {
	value := "test"
	longKey := "123e4567e89b12d3a456426614174000aaaaaaaaaaaaaaaaaaaaaaaaa"
	correctKey := "123e4567e89b12d3a456426614174000"

	cipherValueLongKey, errLongKey := utils.Encrypt(value, longKey)
	assert.Nil(t, errLongKey)

	cipherValueCorrectKey, errCorrectKey := utils.Encrypt(value, correctKey)
	assert.Nil(t, errCorrectKey)

	decryptedValueLongKey, errDecryptLongKey := utils.Decrypt(cipherValueLongKey, longKey)
	assert.Nil(t, errDecryptLongKey)

	decryptedValueCorrectKey, errDecryptCorrectKey := utils.Decrypt(cipherValueCorrectKey, correctKey)
	assert.Nil(t, errDecryptCorrectKey)

	assert.Equal(t, decryptedValueLongKey, decryptedValueCorrectKey)

}
