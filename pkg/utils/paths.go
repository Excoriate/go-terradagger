package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func IsPathValidAndExist(path string) error {
	cleanPath := filepath.Clean(path)

	_, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path %s does not exist", cleanPath)
		}
		return fmt.Errorf("error checking the path %s: %v", cleanPath, err)
	}

	// At this point, the path exists. We can return nil without distinguishing between a file and a directory.
	return nil
}
