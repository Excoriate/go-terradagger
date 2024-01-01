package terradagger

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	ctx := context.Background()
	currentDir, _ := os.Getwd()

	tests := []struct {
		name             string
		options          *ClientOptions
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:        "With valid rootDir on its format, but it does not exist.",
			options:     &ClientOptions{RootDir: "valid/relative/path"},
			expectError: true,
		},
		{
			name:        "With nil options",
			options:     nil,
			expectError: false,
		},
		{
			name:             "With invalid options",
			options:          &ClientOptions{RootDir: "../invalid/root/path"},
			expectError:      true,
			expectedErrorMsg: "TerraDagger initialization error: failed to resolve root directory ../invalid/root/pat",
		},
		{
			name:             "The rootDir is not relative",
			options:          &ClientOptions{RootDir: currentDir},
			expectError:      true,
			expectedErrorMsg: "is not a relative path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(ctx, tt.options)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.NotEmpty(t, client.ID)
				assert.Equal(t, ctx, client.Ctx)

				if tt.options == nil {
					assert.Equal(t, config.defaultSRCPath, client.ClientOptions.RootDir)
				} else {
					assert.Equal(t, tt.options.RootDir, client.ClientOptions.RootDir)
				}

				assert.NotEmpty(t, client.Paths.CurrentDir)
				assert.NotEmpty(t, client.Paths.HomeDir)
				assert.NotEmpty(t, client.Paths.RootDirRelative)
				assert.NotEmpty(t, client.Paths.RootDirAbsolute)
				assert.NotEmpty(t, client.Paths.MountDirPath)
				assert.NotEmpty(t, client.Paths.TerraDagger)
				assert.NotEmpty(t, client.HostEnvVars)
				assert.NotEmpty(t, client.Dirs.TerraDaggerDir)
				assert.NotEmpty(t, client.Dirs.TerraDaggerExportDir)
			}
		})
	}
}

func TestCreateTerraDaggerDirs(t *testing.T) {
	ctx := context.Background()
	td, err := New(ctx, &ClientOptions{RootDir: "."})

	terraDaggerExportPath := resolveTerraDaggerExportPath(td.Paths.TerraDagger, td.ID)

	assert.NoError(t, err)

	tests := []struct {
		name            string
		client          *TD
		failIfDirExists bool
		expectError     bool
	}{
		{
			name:            "Happy path, all correct and the .terradagger dir do not exist.",
			client:          td,
			failIfDirExists: false,
			expectError:     false,
		},
		// The failIfDirExists is set to true, and the dir already exist.
		{
			name:            "The failIfDirExists is set to true, and the dir already exist.",
			client:          td,
			failIfDirExists: true,
			expectError:     true,
		},
		// The failIfDirExists is set to true, and the export path already exist.
		{
			name:            "The failIfDirExists is set to false, and the export path already exist.",
			client:          td,
			failIfDirExists: false,
			expectError:     false,
		},
		{
			// The terraDaggerDir already exist, so it won't do anything.
			name:            "The terraDaggerDir already exist, so it won't do anything.",
			client:          td,
			failIfDirExists: false,
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.failIfDirExists && tt.name == "The failIfDirExists is set to true, "+
				"and the dir already exist." {
				_ = os.MkdirAll(td.Paths.TerraDagger, os.ModePerm)
			}

			if tt.name == "The failIfDirExists is set to false, and the export path already exist." {
				_ = os.MkdirAll(terraDaggerExportPath, os.ModePerm)
			}

			err := tt.client.CreateTerraDaggerDirs(tt.failIfDirExists)

			if tt.expectError {
				assert.Error(t, err)
				if tt.name == "The failIfDirExists is set to true, and the dir already exist." {
					assert.Contains(t, err.Error(), "already exist")
				}
			} else {
				assert.NoError(t, err)
				assert.DirExists(t, tt.client.Paths.TerraDagger)
				assert.DirExists(t, terraDaggerExportPath)
			}
		})
	}

	// Clean up
	_ = os.RemoveAll(td.Paths.TerraDagger)
}
