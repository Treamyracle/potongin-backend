package utils

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const shortCodeLength = 7

// GenerateShortCode membuat string acak aman (cryptographically secure)
func GenerateShortCode() (string, error) {
	b := make([]byte, shortCodeLength)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}

// GenerateAPIKey (bisa gunakan fungsi yang sama atau lebih panjang)
func GenerateAPIKey() (string, error) {
	b := make([]byte, 32) // API Key lebih panjang
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return "pk_" + string(b), nil // "pk_" prefix for "potongin key"
}
