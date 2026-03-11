package crypto

import (
	"errors"

	"golang.org/x/crypto/argon2"
)

func DeriveKey(password string, salt []byte) ([]byte, error) {
	if len(salt) != SaltSize {
		return nil, errors.New("invalid salt size")
	}

	key := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, KeySize)

	return key, nil
}

func DeriveKeyWithSalt(password string, salt []byte) ([]byte, error) {
	return DeriveKey(password, salt)
}
