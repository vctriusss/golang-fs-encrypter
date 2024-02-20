package crypto


import (
	"crypto/cipher"
	"crypto/rand"
)


func EncryptBytes(bytes []byte, gcm *cipher.AEAD) ([]byte, error) {
	nonce := make([]byte, (*gcm).NonceSize())
    _, err := rand.Read(nonce)
    if err != nil {
        return nil, err
    }

    return (*gcm).Seal(nonce, nonce, bytes, nil), nil
}


func DecryptBytes(bytesEncrypted []byte, gcm *cipher.AEAD) ([]byte, error) {
    nonceSize := (*gcm).NonceSize()
    nonce, fileBytesEncrypted := bytesEncrypted[:nonceSize], bytesEncrypted[nonceSize:]

    fileBytes, err := (*gcm).Open(nil, []byte(nonce), []byte(fileBytesEncrypted), nil)
    if err != nil {
        return nil, err
    }

	return fileBytes, nil
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