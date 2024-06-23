package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

func Hash(value string, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(value + salt))

	return hex.EncodeToString(hash.Sum(nil))
}

func Encrypt(value string, key string) (string, error) {
	key = key[0:32]
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize()) // constant nonce
	cipherText := gcm.Seal(nonce, nonce, []byte(value), nil)
	return hex.EncodeToString(cipherText), nil
}

func Decrypt(cipherText string, key string) (string, error) {
	key = key[0:32]
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	data, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, cipherText := data[:nonceSize], string(data[nonceSize:])
	plainText, err := gcm.Open(nil, nonce, []byte(cipherText), nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
