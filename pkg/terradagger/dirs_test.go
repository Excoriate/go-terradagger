package terradagger

import (
	"path/filepath"
	"testing"

	"github.com/Excoriate/go-terradagger/pkg/utils"

	"github.com/stretchr/testify/assert"
)

// func TestResolveMountDirPath(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()
// 	// Override the global Util variable with our mock
//
// 	mockDirUtils := mocks.NewMockDirUtilities(mockCtrl)
//
// 	cwd := "user/this/is/cwd/result"
//
// 	testCases := []struct {
// 		name          string
// 		mountDirPath  string
// 		expectedPath  string
// 		expectedError error
// 	}{
// 		{
// 			// mountDirPath is empty, so it'll return the current directory
// 			name:          "mountDirPath is empty",
// 			mountDirPath:  "",
// 			expectedPath:  filepath.Join(cwd, "."),
// 			expectedError: nil,
// 		},
// 		// mountDirPath is ".", so it'll return the current directory
// 		{
// 			name:          "mountDirPath is .",
// 			mountDirPath:  ".",
// 			expectedPath:  filepath.Join(cwd, "."),
// 			expectedError: nil,
// 		},
// 	}
//
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			mockDirUtils.EXPECT().GetCurrentDir().Return(cwd).AnyTimes()
// 			mockDirUtils.EXPECT().IsValidDir(gomock.Any()).Return(tc.expectedError).AnyTimes()
//
// 			path, err := resolveMountDirPath(tc.mountDirPath)
//
// 			// assert.Equal(t, tc.expectedPath, path)
// 			// TODO: Refactor the resolveMountDirPath to accept an interface,
// 			//  so it'll use the cwd instead of the actual current directory.
// 			assert.NotEmpty(t, path)
// 			assert.Equal(t, tc.expectedError, err)
// 		})
// 	}
// }

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
					tt.rootDir = defaultRootDirRelative
				}
				expectedPath, _ := filepath.Abs(tt.rootDir)
				expectedPath = filepath.Join(expectedPath, terraDaggerDir)
				assert.Equal(t, expectedPath, path)

				// The generated terraDaggerPath should have the terraDaggerDir
				assert.Contains(t, path, terraDaggerDir)
			}
		})
	}
}

func TestResolveMountDirPath(t *testing.T) {
	dirUtils := utils.DirUtils{}
	currentDir := dirUtils.GetCurrentDir()

	// Create an invalid directory from a valid one

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
			rootDirConfig, err := resolveRootDir(tt.rootDir)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.rootDir == "" {
					tt.rootDir = defaultRootDirRelative
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
			name:            "terraDagger path is 'terraDagger' and the ID is valid",
			terraDaggerPath: "asda",
			terraDaggerID:   "some-valid-id",
			expectedError:   false,
		},
		{
			name:            "terraDagger path is 'terraDagger' and the ID is valid",
			terraDaggerPath: ".terradagger",
			terraDaggerID:   "some-valid-id-again",
			expectedError:   false,
		},
	}

	for _, tt := range tests {
		path := resolveTerraDaggerExportPath(tt.terraDaggerPath, tt.terraDaggerID)
		assert.NotEmptyf(t, path, "expected path to be non-empty")
		assert.Contains(t, path, tt.terraDaggerID)
		assert.Contains(t, path, terraDaggerExportDir)
		assert.Equal(t, path[len(path)-len(terraDaggerExportDir):], terraDaggerExportDir)
	}
}
