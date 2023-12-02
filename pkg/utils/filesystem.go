package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

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

func FindGitRepoDir(levels int) (string, error) {
	// Get the current working directory
	pathname, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(pathname)
	if err != nil {
		return "", fmt.Errorf("error converting path %s to absolute path: %w", pathname, err)
	}

	for i := 0; i < levels; i++ {
		gitPath := filepath.Join(absPath, ".git")
		if stat, err := os.Stat(gitPath); err == nil && stat.IsDir() {
			return absPath, nil
		}
		parentPath := filepath.Dir(absPath)

		// To avoid going beyond the root ("/" or "C:\"), check if we're already at the root
		if parentPath == absPath {
			return "", fmt.Errorf("reached root directory, no Git repository found")
		}

		absPath = parentPath
	}

	return "", fmt.Errorf("no Git repository found in %s or any of its parent directories", pathname)
}

func IsValidDir(path string) error {
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
			return nil
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

func GetHomeDir() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}

func GetCurrentDir() string {
	currentDir, _ := os.Getwd()
	return currentDir
}

type IsSubDirOfOptions struct {
	ParentDir string
	ChildDir  string
}

type IsValidRelativeToBaseOptions struct {
	BaseDir      string
	RelativePath string
}

func FileExist(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", path)
		}
		return fmt.Errorf("error checking the path %s: %v", path, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%s is a directory", path)
	}
	return nil
}
