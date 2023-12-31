package terradagger

import (
	"path/filepath"
	"testing"

	"github.com/Excoriate/go-terradagger/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestResolveTerraDaggerPath(t *testing.T) {
	tests := []struct {
		name          string
		rootDir       string
		expectedError bool
	}{
		{
			name:          "Empty root directory",
			rootDir:       "",
			expectedError: false,
		},
		{
			name:          "Valid relative path",
			rootDir:       "some/valid/relative/path",
			expectedError: false,
		},
		{
			name:          "Invalid relative path",
			rootDir:       "../invalid/path",
			expectedError: false,
		},
		{
			name:          "Should resolve to the current directory",
			rootDir:       ".",
			expectedError: false,
		},
		{
			name:          "Invalid, it's an absolute path",
			rootDir:       "/some/absolute/path",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := resolveTerraDaggerPath(tt.rootDir)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.rootDir == "" {
					tt.rootDir = config.defaultSRCPath
				}
				expectedPath, _ := filepath.Abs(tt.rootDir)
				expectedPath = filepath.Join(expectedPath, config.terraDaggerDir)
				assert.Equal(t, expectedPath, path)

				// The generated terraDaggerPath should have the terraDaggerDir
				assert.Contains(t, path, config.terraDaggerDir)
			}
		})
	}
}

func TestResolveMountDirPath(t *testing.T) {
	dirUtils := utils.DirUtils{}
	currentDir := dirUtils.GetCurrentDir()

	tests := []struct {
		name          string
		mountDirPath  string
		expectedPath  string
		expectedError bool
	}{
		{
			name:          "Empty mountDirPath",
			mountDirPath:  "",
			expectedPath:  filepath.Join(currentDir, "."),
			expectedError: false,
		},
		{
			name:          "Current directory mountDirPath",
			mountDirPath:  ".",
			expectedPath:  currentDir,
			expectedError: false,
		},
		{
			name:          "Valid relative mountDirPath but it does not exist",
			mountDirPath:  "valid/path",
			expectedPath:  filepath.Join(currentDir, "valid/path"),
			expectedError: true,
		},
		{
			name:          "Invalid relative mountDirPath",
			mountDirPath:  "..///asd",
			expectedPath:  filepath.Join(currentDir, "..///asd"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := resolveMountDirPath(tt.mountDirPath)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPath, path)
			}
		})
	}
}

func TestResolveRootDir(t *testing.T) {
	tests := []struct {
		name          string
		rootDir       string
		expectedError bool
	}{
		{
			name:          "Empty root directory should resolve to the default '.'",
			rootDir:       "",
			expectedError: false,
		},
		{
			name:          "Should resolve to the current directory",
			rootDir:       ".",
			expectedError: false,
		},
		{
			name:          "Invalid, it's an absolute path",
			rootDir:       "/some/absolute/path",
			expectedError: true,
		},
		{
			name:          "The filepath.Abs return an error, it's not possible to parse this path as absolute",
			rootDir:       "some/invalid/path",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootDirConfig, err := resolveSRC(tt.rootDir)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.rootDir == "" {
					tt.rootDir = config.defaultSRCPath
				}
				expectedPath, _ := filepath.Abs(tt.rootDir)
				assert.Equal(t, expectedPath, rootDirConfig.RootDirAbsolute)
				assert.Equal(t, tt.rootDir, rootDirConfig.RootDirRelative)
			}
		})
	}
}

func TestResolveTerraDaggerExportPath(t *testing.T) {
	tests := []struct {
		name            string
		terraDaggerPath string
		terraDaggerID   string
		expectedError   bool
	}{
		{
			name:            "TerraDagger path is 'TerraDagger' and the ID is valid",
			terraDaggerPath: "asda",
			terraDaggerID:   "some-valid-id",
			expectedError:   false,
		},
		{
			name:            "TerraDagger path is 'TerraDagger' and the ID is valid",
			terraDaggerPath: ".terradagger",
			terraDaggerID:   "some-valid-id-again",
			expectedError:   false,
		},
	}

	for _, tt := range tests {
		path := resolveTerraDaggerExportPath(tt.terraDaggerPath, tt.terraDaggerID)
		assert.NotEmptyf(t, path, "expected path to be non-empty")
		assert.Contains(t, path, tt.terraDaggerID)
		assert.Contains(t, path, config.terraDaggerExportDir)
		assert.Equal(t, path[len(path)-len(config.terraDaggerExportDir):], config.terraDaggerExportDir)
	}
}
