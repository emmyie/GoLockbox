package crypto

import (
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

func EncryptBytes(plaintext, key []byte) ([]byte, error) {
	if len(key) != KeySize {
		return nil, ErrInvalidKeyLength
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, chacha20poly1305.NonceSizeX)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// ciphertext = nonce || aead(plaintext)
	ciphertext := aead.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

func DecryptBytes(ciphertext, key []byte) ([]byte, error) {
	if len(key) != KeySize {
		return nil, ErrInvalidKeyLength
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < chacha20poly1305.NonceSizeX {
		return nil, errors.New("ciphertext too short")
	}

	nonce := ciphertext[:chacha20poly1305.NonceSizeX]
	data := ciphertext[chacha20poly1305.NonceSizeX:]

	plaintext, err := aead.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
