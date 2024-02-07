package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func DirExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func DirExistE(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}

	return nil
}

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
