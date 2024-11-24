package auth

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

const (
	time    = 1         // number of iterations
	memory  = 64 * 1024 // memory in KiB
	threads = 4         // parallelism
	keyLen  = 32        // length of the generated key
)

func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func HashPassword(password string) (string, error) {
	salt, err := generateSalt(16)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
	fullHash := append(salt, hash...)
	return base64.RawStdEncoding.EncodeToString(fullHash), nil
}

func VerifyPassword(password string, fullHash string) bool {
	data, err := base64.RawStdEncoding.DecodeString(fullHash)
	if err != nil {
		return false
	}

	salt := data[:16]
	hash := data[16:]

	newHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	return string(hash) == string(newHash)
}
