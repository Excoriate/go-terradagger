package config

import (
	"fmt"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/Excoriate/go-terradagger/pkg/o11y"
	"github.com/Excoriate/go-terradagger/pkg/utils"
)

type DaggerBackend struct {
	Paths    *DaggerPaths
	Excluded *ExcludedDirAndFiles
}

type ExcludedDirAndFiles struct {
	ExcludedDirs  []string
	ExcludedFiles []string
}

type DaggerPaths struct {
	MountPathPrefix string
}

type DaggerBackendConfigClient struct {
	l o11y.LoggerInterface
}

func getEmptyDaggerConfig() *DaggerBackend {
	return &DaggerBackend{
		Paths:    &DaggerPaths{},
		Excluded: &ExcludedDirAndFiles{},
	}
}

func NewDaggerBackendConfigClient(logger o11y.LoggerInterface) *DaggerBackendConfigClient {
	var logImpl o11y.LoggerInterface
	if logger == nil {
		logImpl = o11y.DefaultLogger()
	}

	return &DaggerBackendConfigClient{
		l: logImpl,
	}
}

type DaggerBackendConfig interface {
	GetMountPrefix() string
	GetExcludedDirs() []string
	GetExcludedFiles() []string
	Configure() (*DaggerBackend, error)
	GetMountPathAsDaggerDir(mountPath string, daggerClient *dagger.Client) (*dagger.Directory, error)
}

func (c *DaggerBackendConfigClient) GetMountPrefix() string {
	return MountPathPrefixInDagger
}

func (c *DaggerBackendConfigClient) Configure() (*DaggerBackend, error) {
	cfg := getEmptyDaggerConfig()
	// 1. Getting the paths configuration
	cfg.Paths.MountPathPrefix = c.GetMountPrefix()

	// 2. Getting the default excluded directories and files.
	cfg.Excluded.ExcludedDirs = c.GetExcludedDirs()
	cfg.Excluded.ExcludedFiles = c.GetExcludedFiles()

	return cfg, nil
}

type GetWorkDirOptions struct {
	MountPathAbsolute string
	WorkDirRelative   string
}

func (c *DaggerBackendConfigClient) GetExcludedDirs() []string {
	return excludedDirsDefault
}

func (c *DaggerBackendConfigClient) GetExcludedFiles() []string {
	return excludedFilesDefault
}

func (c *DaggerBackendConfigClient) GetMountPathAsDaggerDir(mountPath string, daggerClient *dagger.Client) (*dagger.Directory, error) {
	dirUtils := utils.DirUtils{}

	if daggerClient == nil {
		return nil, fmt.Errorf("failed to get the mount path as a dagger directory, the dagger client is nil")
	}

	if mountPath == "" || mountPath == "." {
		mountPathAbs := dirUtils.GetCurrentDir()
		return daggerClient.Host().Directory(mountPathAbs), nil
	}

	if err := IsAValidTerraDaggerDirRelative(mountPath); err != nil {
		return nil, fmt.Errorf("failed to get the mount path as a dagger directory, the mount path is invalid: %w", err)
	}

	mountPathAbs, _ := filepath.Abs(mountPath)
	mountDirAsDaggerFormat := daggerClient.Host().Directory(mountPathAbs)

	// entries, err := mountDirAsDaggerFormat.Entries(context.Background())
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get the mount path as a dagger directory, the mount path is invalid: %w", err)
	// }
	//
	// if len(entries) == 0 {
	// 	return nil, fmt.Errorf("failed to get the mount path as a dagger directory, the mount path is empty")
	// }
	//

	return mountDirAsDaggerFormat, nil
}
