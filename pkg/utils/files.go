package utils

import (
	"fmt"
	"os"
	"strings"
)

func IsValidFileE(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", filePath)
		}

		return fmt.Errorf("unexpected error when checking the file %s: %v", filePath, err)
	}

	return nil
}

func FileHasExtension(filePath, extension string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	if extension == "" {
		return fmt.Errorf("extension cannot be empty")
	}

	if !strings.HasSuffix(filePath, extension) {
		return fmt.Errorf("the file %s does not have the extension %s", filePath, extension)
	}

	return nil
}
