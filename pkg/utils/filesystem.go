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

func FileExistE(path string) error {
	cleanPath := filepath.Clean(path)

	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", cleanPath)
		}
		return fmt.Errorf("error checking the path %s: %v", cleanPath, err)
	}

	if info.IsDir() {
		return fmt.Errorf("%s is a directory", cleanPath)
	}

	return nil
}

func FileExist(path string) bool {
	cleanPath := filepath.Clean(path)

	info, err := os.Stat(cleanPath)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func DeleteFileE(path string) error {
	cleanPath := filepath.Clean(path)

	if err := os.Remove(cleanPath); err != nil {
		return fmt.Errorf("error deleting file %s: %v", cleanPath, err)
	}

	return nil
}

func DeleteDirE(path string) error {
	cleanPath := filepath.Clean(path)

	if err := os.RemoveAll(cleanPath); err != nil {
		return fmt.Errorf("error deleting directory %s: %v", cleanPath, err)
	}

	return nil
}

func CreateFileE(path string) error {
	cleanPath := filepath.Clean(path)

	if err := os.MkdirAll(filepath.Dir(cleanPath), 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %v", filepath.Dir(cleanPath), err)
	}

	if _, err := os.Create(cleanPath); err != nil {
		return fmt.Errorf("error creating file %s: %v", cleanPath, err)
	}

	return nil
}

func IsAFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func IsAFileE(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error checking the path %s: %v", path, err)
	}

	if info.IsDir() {
		return fmt.Errorf("%s is a directory", path)
	}

	return nil
}

func IsRelativeE(path string) error {
	if filepath.IsAbs(path) {
		return fmt.Errorf("path %s is not relative", path)
	}
	return nil
}

func IsRelative(path string) bool {
	return filepath.IsAbs(path)
}

func IsAbsolute(path string) bool {
	return filepath.IsAbs(path)
}

func IsAbsoluteE(path string) error {
	if !filepath.IsAbs(path) {
		return fmt.Errorf("path %s is not absolute", path)
	}
	return nil
}

func GetFilesByExtension(path, extension string) ([]string, error) {
	var files []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == extension {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error getting files by extension %s in path %s: %v", extension, path, err)
	}

	return files, nil
}
