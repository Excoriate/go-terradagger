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
	IsValidDirRelative(path string) error
	IsValidDirAbsolute(path string) error
	isValidDirCommon(path string) error
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

func (du *DirUtils) IsValidDirRelative(path string) error {
	if filepath.IsAbs(path) {
		return fmt.Errorf("expected relative path, but got absolute path: %s", path)
	}

	return du.isValidDirCommon(path)
}

func (du *DirUtils) IsValidDirAbsolute(path string) error {
	if !filepath.IsAbs(path) {
		return fmt.Errorf("expected absolute path, but got relative path: %s", path)
	}

	return du.isValidDirCommon(path)
}

// isValidDirCommon is a helper function containing the shared logic for path validation.
func (du *DirUtils) isValidDirCommon(path string) error {
	cleanPath := filepath.Clean(path)

	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s, error: %w", cleanPath, err)
		}
		return fmt.Errorf("failed to stat the path: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", cleanPath)
	}

	return nil
}
