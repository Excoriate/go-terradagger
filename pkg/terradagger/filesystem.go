package terradagger

import (
	"fmt"
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/config"
	"github.com/Excoriate/go-terradagger/pkg/utils"

	"dagger.io/dagger"
)

type fsResolverClient struct {
	td *TD
}

type InstancePathsResolved struct {
	TerraDagger         string
	ExportInHost        string // To host
	Cache               string
	ImportFromContainer string // From ContainerClient.
}

type fsResolver interface {
	IsMountPathValid(mountPath string) error
	IsMountAndWorkDirPathValid(mountPathRelative, workDirPathRelative string) error
	ResolveMountPathIfEmpty(mountPath string) string
	ResolveWorkDirIfEmpty(workDirPath string) string
	ResolveAuxPaths(TerraDaggerDirPathResolved string) ResolveAuxPathsResult
}

func newDirResolverClient(td *TD) fsResolver {
	return &fsResolverClient{
		td: td,
	}
}

func (fs *fsResolverClient) IsMountPathValid(mountPath string) error {
	if err := config.IsAValidTerraDaggerDirRelative(mountPath); err != nil {
		return fmt.Errorf("failed to Validate the mount path, the mount path is invalid: %w", err)
	}

	return nil
}

func (fs *fsResolverClient) IsMountAndWorkDirPathValid(mountPathRelative, workDirPathRelative string) error {
	fullPath := filepath.Join(mountPathRelative, workDirPathRelative)
	dirUtils := utils.DirUtils{}
	if err := dirUtils.IsValidDirE(fullPath); err != nil {
		return fmt.Errorf("failed to Validate the mount and work dir path, the path %s is invalid: %w", fullPath, err)
	}

	return nil
}

func (fs *fsResolverClient) ResolveMountPathIfEmpty(mountPath string) string {
	if mountPath == "" {
		mountPath = fs.td.Config.TerraDagger.Paths.Host.CurrentDir
		fs.td.Logger.Debug("Automatically resolved the mount path to the default '.' path")
	}

	return mountPath
}

func (fs *fsResolverClient) ResolveWorkDirIfEmpty(workDirPath string) string {
	if workDirPath == "" {
		workDirPath = fs.td.Config.TerraDagger.Paths.Host.CurrentDir
		fs.td.Logger.Debug("Automatically resolved the work dir path to the default '.' path")
	}

	return workDirPath
}

type ResolveAuxPathsResult struct {
	ExportPath    string
	ExportPathAbs string
	CachePath     string
	CachePathAbs  string
	ImportPath    string
	ImportPathAbs string
}

func (fs *fsResolverClient) ResolveAuxPaths(terraDaggerConfigDirPath string) ResolveAuxPathsResult {
	absPath, _ := utils.ConvertToAbsolute(terraDaggerConfigDirPath)

	return ResolveAuxPathsResult{
		ExportPath:    filepath.Join(terraDaggerConfigDirPath, fs.td.Config.TerraDagger.Dirs.TerraDaggerExportDir),
		ExportPathAbs: filepath.Join(absPath, fs.td.Config.TerraDagger.Dirs.TerraDaggerExportDir),
		CachePath:     filepath.Join(terraDaggerConfigDirPath, fs.td.Config.TerraDagger.Dirs.TerraDaggerCacheDir),
		CachePathAbs:  filepath.Join(absPath, fs.td.Config.TerraDagger.Dirs.TerraDaggerCacheDir),
		ImportPath:    filepath.Join(terraDaggerConfigDirPath, fs.td.Config.TerraDagger.Dirs.TerraDaggerImportDir),
		ImportPathAbs: filepath.Join(absPath, fs.td.Config.TerraDagger.Dirs.TerraDaggerImportDir),
	}
}

type ResolveWorkDirResult struct {
	WorkDirPath    string
	WorkDirPathAbs string
}

type resolveWorkDirOptions struct {
	WorkDirPathRelative string
	MountPathAbsolute   string
	MountPrefix         string
}

type DirConfig struct {
	MountDir               *dagger.Directory
	WorkDirPath            string
	WorkDirPathInContainer string
}
