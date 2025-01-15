package datautils

import (
	cRand "crypto/rand"
	mRand "math/rand"
)

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = charset[mRand.Intn(len(charset))]
	}

	return string(result)
}

func SecureRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	if _, err := cRand.Read(result); err != nil {
		return RandomString(length)
	}

	for i := range result {
		result[i] = charset[int(result[i])%len(charset)]
	}

	return string(result)
}
