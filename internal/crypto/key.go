package crypto

import (
	"crypto/rand"
	"io"
)

const (
	SaltSize = 16
	KeySize  = 32
)

func GenerateRandomKey(size int) ([]byte, error) {
	key := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}

	return key, nil
}

func GenerateSalt() ([]byte, error) {
	return GenerateRandomKey(SaltSize)
}
