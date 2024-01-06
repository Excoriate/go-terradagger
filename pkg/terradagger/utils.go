package terradagger

import (
	"fmt"
	"path/filepath"
)

func buildWorkDirPath(mountDirPath, workDirPath string) string {
	mountDirPath = filepath.Clean(mountDirPath)
	workDirPath = filepath.Clean(workDirPath)

	return filepath.Join(mountDirPath, workDirPath)
}

type GetTerraDaggerPathOptions struct {
	MountDirPath  string
	WorkDirPath   string
	FileOrDirName string
}

func GetTerraDaggerPath(options *GetTerraDaggerPathOptions) (string, error) {
	if options == nil {
		return "", fmt.Errorf("failed to get the terra dagger path, the options are nil")
	}

	if options.MountDirPath == "" {
		return "", fmt.Errorf("failed to get the terra dagger path, the mount dir path is empty")
	}

	if options.WorkDirPath == "" {
		return filepath.Join(options.MountDirPath, options.FileOrDirName), nil
	}

	return filepath.Join(options.MountDirPath, options.WorkDirPath, options.FileOrDirName), nil
}

func GetTerraDaggerWorkDirPath(mountDirPath, workDirPath string) string {
	return buildWorkDirPath(mountDirPath, workDirPath)
}

type GetTerraDaggerConfigPathOptions struct {
	TerraDaggerID string
	ConfigDir     string
	FileOrDirName string
}

func GetTerraDaggerConfigPath(options *GetTerraDaggerConfigPathOptions) (string, error) {
	if options == nil {
		return "", fmt.Errorf("failed to get the terra dagger config path, the options are nil")
	}

	if options.TerraDaggerID == "" {
		return "", fmt.Errorf("failed to get the terra dagger config path, the terra dagger id is empty")
	}

	if options.ConfigDir == "" {
		return "", fmt.Errorf("failed to get the terra dagger config path, the config dir is empty")
	}

	if options.FileOrDirName == "" {
		return "", fmt.Errorf("failed to get the terra dagger config path, the file or dir name is empty")
	}

	return filepath.Join(".terradagger", options.TerraDaggerID, options.FileOrDirName), nil
}
