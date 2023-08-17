package helpers

import (
	"os"
	"path/filepath"
	"strings"
)

func LookupEnvString(key string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return ""
}

func WriteFile(content, filePath string) error {
	// Open the file for writing, creating it if it doesn't exist
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	// Convert the string to a byte slice
	contentBytes := []byte(content)
	// Write the content to the file
	_, err = file.Write(contentBytes)
	return err
}

func TrimFileExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func IsStringInList(value string, list []string) bool {
	for _, str := range list {
		if strings.Contains(value, str) {
			return true
		}
	}
	return false
}

func ReadFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	stat, _ := file.Stat()
	content := make([]byte, stat.Size())
	_, err = file.Read(content)
	return string(content), err
}
