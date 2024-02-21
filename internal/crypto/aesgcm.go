package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

type AESGCM struct {
	key []byte
	gcm cipher.AEAD
}

func NewAESGCM(key []byte) (*AESGCM, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	return &AESGCM{key: key, gcm: gcm}, nil
}

func (c *AESGCM) EncryptBytes(bytes []byte) ([]byte, error) {
	nonce := make([]byte, c.gcm.NonceSize())
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	return c.gcm.Seal(nonce, nonce, bytes, nil), nil
}

func (c *AESGCM) DecryptBytes(bytesEncrypted []byte) ([]byte, error) {
	nonceSize := c.gcm.NonceSize()
	nonce, fileBytesEncrypted := bytesEncrypted[:nonceSize], bytesEncrypted[nonceSize:]

	fileBytes, err := c.gcm.Open(nil, []byte(nonce), []byte(fileBytesEncrypted), nil)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}
