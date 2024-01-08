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

type FoundPaths struct {
	FilesPathFound []string
	DirsPathFound  []string
}

func FilterPathsByNames(rootPath string, fileNames, dirNames []string) (*FoundPaths, error) {
	foundPaths := &FoundPaths{}

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the current path is a file and matches any of the provided file names.
		if info.Mode().IsRegular() {
			for _, name := range fileNames {
				if filepath.Base(path) == name {
					foundPaths.FilesPathFound = append(foundPaths.FilesPathFound, path)
					break
				}
			}
		}

		// Check if the current path is a directory and matches any of the provided directory names.
		if info.IsDir() {
			for _, name := range dirNames {
				if filepath.Base(path) == name {
					foundPaths.DirsPathFound = append(foundPaths.DirsPathFound, path)
					break
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return foundPaths, nil
}
