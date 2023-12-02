package terradagger

import (
	"fmt"
	"github.com/Excoriate/go-terradagger/pkg/config"
	"github.com/Excoriate/go-terradagger/pkg/errors"
	"github.com/Excoriate/go-terradagger/pkg/utils"
	"path/filepath"

	"dagger.io/dagger"
)

type DirConfig struct {
	MountDir               *dagger.Directory
	WorkDirPath            string
	WorkDirPathInContainer string
}

// getDirs returns the mount directory, and the work directory.
// The mount directory is the directory that is mounted in the container.
// The work directory is the directory that is used by the commands passed.
func getDirs(client *dagger.Client, mountDir, workDir string) *DirConfig {
	mountDirDagger := client.Host().Directory(mountDir)
	workDirPathInContainer := fmt.Sprintf("%s/%s", config.MountPathPrefixInDagger, filepath.Clean(workDir))

	return &DirConfig{
		MountDir:               mountDirDagger,
		WorkDirPath:            workDir,
		WorkDirPathInContainer: workDirPathInContainer,
	}
}

// resolveMountDirPath resolves the mount directory path.
// If the mount directory path is empty, the current directory is used.
func resolveMountDirPath(mountDirPath string) (string, error) {
	currentDir := utils.GetCurrentDir()
	if mountDirPath == "" {
		return filepath.Join(currentDir, "/", "."), nil
	}

	mountDirPath = filepath.Join(currentDir, "/", mountDirPath)

	if err := utils.IsValidDir(mountDirPath); err != nil {
		return "", &errors.ErrTerraDaggerInvalidMountPath{
			ErrWrapped: err,
			MountPath:  mountDirPath,
		}
	}

	return mountDirPath, nil
}
