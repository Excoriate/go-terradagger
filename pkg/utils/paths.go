package utils

import "path/filepath"

func ConvertToAbsFomRelative(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	return filepath.Abs(path)
}
