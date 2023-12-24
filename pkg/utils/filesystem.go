package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

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

func IsRelative(path string) error {
	if filepath.IsAbs(path) {
		return fmt.Errorf("path %s is not relative", path)
	}
	return nil
}
