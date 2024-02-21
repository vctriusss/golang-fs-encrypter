package crypto

import (
	"crypto/rand"
)

type Cipher interface {
	EncryptBytes([]byte) ([]byte, error)
	DecryptBytes([]byte) ([]byte, error)
}

func GenerateBytes(length int) []byte {
	keyBytes := make([]byte, length/2)
	_, err := rand.Read(keyBytes)

	if err != nil {
		for i := 0; i < length/2; i++ {
			keyBytes[i] = 3
		}
	}

	return keyBytes
}
