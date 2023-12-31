package config

import (
	"fmt"
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/utils"
)

func IsAValidTerraDaggerDirRelative(dirPath string) error {
	dirUtils := utils.DirUtils{}
	if dirPath == "" {
		return fmt.Errorf("the dir path is empty")
	}

	if err := utils.IsRelativeE(dirPath); err != nil {
		return fmt.Errorf("the dir path %s is not a relative path", dirPath)
	}

	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return fmt.Errorf("failed to get the absolute path of the dir path %s", dirPath)
	}

	if !dirUtils.DirExist(absPath) {
		return fmt.Errorf("the dir path %s do not exist", dirPath)
	}

	return nil
}

func IsAValidTerraDaggerDirAbsolute(dirPath string) error {
	dirUtils := utils.DirUtils{}
	if dirPath == "" {
		return fmt.Errorf("the dir path is empty")
	}

	if err := utils.IsAbsoluteE(dirPath); err != nil {
		return fmt.Errorf("the dir path %s is not an absolute path", dirPath)
	}

	if !dirUtils.DirExist(dirPath) {
		return fmt.Errorf("the dir path %s do not exist", dirPath)
	}

	return nil
}

func AreFilesToExcludeValid(mountPath string, filesToExclude []string) error {
	if len(filesToExclude) == 0 {
		return nil
	}

	for _, fileToExclude := range filesToExclude {
		filePath := filepath.Join(mountPath, fileToExclude)
		if err := utils.FileExistE(filePath); err != nil {
			return fmt.Errorf("the file %s to exclude does not exist", filePath)
		}

		if err := utils.IsAFileE(filePath); err != nil {
			return fmt.Errorf("the file %s to exclude is not a file", filePath)
		}
	}

	return nil
}

func AreDirsToExcludeValid(mountPath string, dirsToExclude []string) error {
	if len(dirsToExclude) == 0 {
		return nil
	}

	dirUtils := utils.DirUtils{}

	for _, dirToExclude := range dirsToExclude {
		dirPath := filepath.Join(mountPath, dirToExclude)

		if err := dirUtils.IsValidDir(dirPath); err != nil {
			return fmt.Errorf("the dir %s to exclude is not a valid dir", dirPath)
		}
	}

	return nil
}
