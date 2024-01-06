package terradagger

import (
	"fmt"
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/utils"
)

type ResolverPathOptions struct {
	BasePath      string
	FileOrDirName string
}

func ResolveTerraDaggerPath(options *ResolverPathOptions) (string, error) {
	if options == nil {
		return "", nil
	}

	if options.BasePath == "" {
		return "", fmt.Errorf("the base path is empty")
	}

	if options.FileOrDirName == "" {
		return options.BasePath, nil
	}

	return filepath.Join(options.BasePath, options.FileOrDirName), nil
}

func ResolveTerraDaggerDirE(options *ResolverPathOptions) (string, error) {
	if options == nil {
		return "", nil
	}

	resolvedPath, err := ResolveTerraDaggerPath(options)
	if err != nil {
		return "", err
	}

	dirUtils := utils.NewDirUtils()

	if err := dirUtils.IsValidDirE(resolvedPath); err != nil {
		return "", err
	}

	return resolvedPath, nil
}

func ResolveTerraDaggerFileE(options *ResolverPathOptions) (string, error) {
	if options == nil {
		return "", nil
	}

	resolvedPath, err := ResolveTerraDaggerPath(options)
	if err != nil {
		return "", err
	}

	if err := utils.FileExistE(resolvedPath); err != nil {
		return "", err
	}

	return resolvedPath, nil
}

type ResolveTerraDaggerConfigDirPathOptions struct {
	WorkspaceSRCPath   string
	ID                 string
	TerraDaggerDirName string
}

func ResolveTerraDaggerConfigDirPath(options *ResolveTerraDaggerConfigDirPathOptions) (string, error) {
	if options == nil {
		return "", fmt.Errorf("failed to resolve the terra dagger dir path, the options are nil")
	}

	if options.ID == "" {
		return "", fmt.Errorf("failed to resolve the terra dagger dir path, the terra dagger id is empty")
	}

	if options.TerraDaggerDirName == "" {
		return "", fmt.Errorf("failed to resolve the terra dagger dir path, the terra dagger dir name is empty")
	}

	if options.WorkspaceSRCPath == "" {
		return "", fmt.Errorf("failed to resolve the terra dagger dir path, the workspace src path is empty")
	}

	return filepath.Join(options.WorkspaceSRCPath, options.TerraDaggerDirName, options.ID), nil
}
