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
	IsFileValidInHost(filePath string) error
	IsMountPathValid(mountPath string) error
	IsMountAndWorkDirPathValid(mountPathRelative, workDirPathRelative string) error
	ResolveMountPathIfEmpty(mountPath string) string
	ResolveWorkDirIfEmpty(workDirPath string) string
	ResolveTerraDaggerDirPath(options *ResolveTerraDaggerDirPathOptions) (string, error)
	ResolveAuxPaths(TerraDaggerDirPathResolved string) ResolveAuxPathsResult
}

func newDirResolverClient(td *TD) fsResolver {
	return &fsResolverClient{
		td: td,
	}
}

func (fs *fsResolverClient) IsFileValidInHost(filePath string) error {
	if err := utils.FileExists(filePath); err != nil {
		return fmt.Errorf("failed to validate the file in the host, the file %s does not exist: %w", filePath, err)
	}

	return nil
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
	if err := dirUtils.IsValidDir(fullPath); err != nil {
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

type ResolveTerraDaggerDirPathOptions struct {
	WorkspaceSRCPath string
	ID               string
}

func (fs *fsResolverClient) ResolveTerraDaggerDirPath(options *ResolveTerraDaggerDirPathOptions) (string, error) {
	if options == nil {
		return "", fmt.Errorf("failed to resolve the terra dagger dir path, the options are nil")
	}

	if options.ID == "" {
		return "", fmt.Errorf("failed to resolve the terra dagger dir path, the terra dagger id is empty")
	}

	if options.WorkspaceSRCPath == "" {
		options.WorkspaceSRCPath = fs.td.Config.TerraDagger.Paths.Workspace.SRC // Automatically resolve the default '.' path
		fs.td.Logger.Warn("Automatically resolved the src path to the default '.' path")
	}

	var srcPathAbs string
	if utils.IsAbsolute(options.WorkspaceSRCPath) {
		if err := config.IsAValidTerraDaggerDirAbsolute(options.WorkspaceSRCPath); err != nil {
			return "", fmt.Errorf("failed to resolve the terra dagger dir path, the src path is invalid: %w", err)
		}

		srcPathAbs = options.WorkspaceSRCPath
	} else {
		if err := config.IsAValidTerraDaggerDirRelative(options.WorkspaceSRCPath); err != nil {
			return "", fmt.Errorf("failed to resolve the terra dagger dir path, the src path is invalid: %w", err)
		}

		srcPathAbs, _ = filepath.Abs(options.WorkspaceSRCPath)
	}

	terraDaggerDirPath := filepath.Join(srcPathAbs, fs.td.Config.TerraDagger.Dirs.TerraDaggerDir, options.ID)
	return terraDaggerDirPath, nil
}

type ResolveAuxPathsResult struct {
	ExportPath string
	CachePath  string
	ImportPath string
}

func (fs *fsResolverClient) ResolveAuxPaths(TerraDaggerDirPathResolved string) ResolveAuxPathsResult {
	return ResolveAuxPathsResult{
		ExportPath: filepath.Join(TerraDaggerDirPathResolved, fs.td.Config.TerraDagger.Dirs.TerraDaggerExportDir),
		CachePath:  filepath.Join(TerraDaggerDirPathResolved, fs.td.Config.TerraDagger.Dirs.TerraDaggerCacheDir),
		ImportPath: filepath.Join(TerraDaggerDirPathResolved, fs.td.Config.TerraDagger.Dirs.TerraDaggerImportDir),
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
