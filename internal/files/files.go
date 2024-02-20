package files

import (
	"io"
	"os"
)

func ReadFile(path string) ([]byte, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}


func RewriteFile(path string, newBytes []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(newBytes)
	
	return err
}


func WriteFile(path string, bytes []byte) error {
	file, err := os.OpenFile(path, os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(bytes)
	return err
}