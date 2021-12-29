package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Encrypt encrypts data using a password. The password must be between 14 - 32 characters.
func Encrypt(data, password []byte) ([]byte, error) {
	if len(password) > 32 {
		return nil, errors.New("password must be between 14 - 32 chars")
	}

	if len(password) < 32 {
		password = append(password, make([]byte, 32-len(password))...)
	}

	c, err := aes.NewCipher(password)
	if err != nil {
		return nil, err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// Decrypt decrypts data using a password.
func Decrypt(data, password []byte) ([]byte, error) {
	if len(password) > 32 {
		return nil, errors.New("failed to decrypt")
	}

	if len(password) < 32 {
		password = append(password, make([]byte, 32-len(password))...)
	}

	c, err := aes.NewCipher(password)
	if err != nil {
		return nil, err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, data, nil)
}
