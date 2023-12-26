package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirUtils(t *testing.T) {
	du := &DirUtils{}
	homeDir := os.Getenv("HOME") // or use os.UserHomeDir() if you want the actual home directory
	currentDir, _ := os.Getwd()

	tests := []struct {
		name        string
		testFunc    func() string
		expectedRes string
	}{
		{
			name:        "GetHomeDir",
			testFunc:    du.GetHomeDir,
			expectedRes: homeDir,
		},
		{
			name:        "GetCurrentDir",
			testFunc:    du.GetCurrentDir,
			expectedRes: currentDir,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.testFunc()
			assert.Equal(t, tt.expectedRes, res)
		})
	}
}

func TestDirExist(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testDirExist")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir) // Clean up

	du := &DirUtils{}

	assert.True(t, du.DirExist(tempDir), "Temporary directory should exist")
	assert.False(t, du.DirExist(filepath.Join(tempDir, "nonexistent")), "Non-existent directory should not exist")
}

func TestIsValidDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testIsValidDir")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir) // Clean up

	du := &DirUtils{}

	assert.NoError(t, du.IsValidDir(tempDir), "Temporary directory should be valid")
}

func TestDirExistAndHasContent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testDirExistAndHasContent")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir) // Clean up

	tempFile, err := os.CreateTemp(tempDir, "testfile")
	assert.NoError(t, err)
	tempFile.Close()

	du := &DirUtils{}

	assert.NoError(t, du.DirExistAndHasContent(tempDir), "Temporary directory with content should be valid")
	assert.Error(t, du.DirExistAndHasContent(filepath.Join(tempDir, "nonexistent")), "Non-existent directory should return error")
}
