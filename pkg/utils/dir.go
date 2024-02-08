package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func IsValidDirE(path string) error {
	// Clean the path to remove any unnecessary parts.
	cleanPath := filepath.Clean(path)

	// Use os.Stat to check if the path exists and get file info.
	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			// The path does not exist.
			return fmt.Errorf("path does not exist: %s, error: %w", cleanPath, err)
		}
		// There was some problem accessing the path.
		return fmt.Errorf("failed to stat the path: %w", err)
	}

	// Check if the path is a directory.
	if !info.IsDir() {
		// The path is not a directory.
		return fmt.Errorf("path is not a directory: %s", cleanPath)
	}

	// The path is a valid directory.
	return nil
}

func DirExistAndHasContent(dirPath string) error {
	if dirPath == "" {
		return fmt.Errorf("directory path cannot be empty")
	}

	currentDir, _ := os.Getwd()

	_, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory %s does not exist in current directory %s", dirPath, currentDir)
		}

		return fmt.Errorf("unexpected error when checking the directory %s: %v", dirPath, err)
	}

	return nil
}

func DirHasContentWithCertainExtension(dirPath string, extensions []string) error {
	if dirPath == "" {
		return fmt.Errorf("directory path cannot be empty")
	}

	if len(extensions) == 0 {
		return fmt.Errorf("extensions cannot be empty")
	}

	err := DirExistAndHasContent(dirPath)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read the directory %s: %v", dirPath, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		for _, ext := range extensions {
			if filepath.Ext(file.Name()) == ext {
				return nil
			}
		}
	}

	return fmt.Errorf("directory %s does not contain any files with the following extensions: %v", dirPath, extensions)
}
