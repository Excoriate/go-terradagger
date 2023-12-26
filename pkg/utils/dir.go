package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

type DirUtilities interface {
	GetHomeDir() string
	GetCurrentDir() string
	DirExist(path string) bool
	IsValidDir(path string) error
	DirExistAndHasContent(dirPath string) error
}

type DirUtils struct{}

func (du *DirUtils) GetHomeDir() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}

func (du *DirUtils) GetCurrentDir() string {
	currentDir, _ := os.Getwd()
	return currentDir
}

func (du *DirUtils) DirExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (du *DirUtils) IsValidDir(path string) error {
	// Ensure we have an absolute path to avoid confusion with relative paths.
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Use os.Stat to check if the path exists and get file info.
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			// The path does not exist.
			return fmt.Errorf("path does not exist: %s. Failed with error: %w", absPath, err)
		}
		// There was some problem accessing the path.
		return fmt.Errorf("failed to stat the path: %w", err)
	}

	// Check if the path is a directory.
	if !info.IsDir() {
		// The path is not a directory.
		return fmt.Errorf("path is not a directory: %s", absPath)
	}

	// The path is a valid directory.
	return nil
}

func (du *DirUtils) DirExistAndHasContent(dirPath string) error {
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
