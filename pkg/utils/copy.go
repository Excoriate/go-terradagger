package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyPaths copies a list of files and directories to the specified destination directory.
func CopyPaths(filePaths, dirPaths []string, destDir string) error {
	// Validate and create the destination directory if it doesn't exist.
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Copy files.
	for _, filePath := range filePaths {
		destFilePath := filepath.Join(destDir, filepath.Base(filePath))
		if err := copyFile(filePath, destFilePath); err != nil {
			return err
		}
	}

	// Copy directories.
	for _, dirPath := range dirPaths {
		destDirPath := filepath.Join(destDir, filepath.Base(dirPath))
		if err := copyDirectory(dirPath, destDirPath); err != nil {
			return err
		}
	}

	return nil
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}
	return nil
}

// copyDirectory recursively copies a directory from src to dst.
func copyDirectory(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("failed to resolve relative path: %w", err)
		}
		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		return copyFile(path, destPath)
	})
}
