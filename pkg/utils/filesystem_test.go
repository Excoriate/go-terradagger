package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestDirExistAndHasContent(t *testing.T) {
// 	dirUtils := DirUtils{}
// 	t.Execute("empty directory path", func(t *testing.T) {
// 		err := dirUtils.DirExistAndHasContent("")
// 		assert.Error(t, err, "Expected error due to empty directory path not found")
// 		assert.EqualError(t, err, "directory path cannot be empty")
// 	})
//
// 	t.Execute("non existent directory", func(t *testing.T) {
// 		err := dirUtils.DirExistAndHasContent("/nonexistent_directory")
// 		assert.Error(t, err, "Expected error due to non existent directory")
// 		assert.Contains(t, err.Error(), "does not exist in current directory")
// 	})
//
// 	t.Execute("existent directory", func(t *testing.T) {
// 		testDir, _ := ioutil.TempDir("", "existent_directory")
// 		defer os.RemoveAll(testDir)
//
// 		err := dirUtils.DirExistAndHasContent(testDir)
// 		assert.NoError(t, err, "Unexpected error for existent directory")
// 	})
// }

func TestFindGitRepoDir(t *testing.T) {
	t.Run("Find git repo in parent directories", func(t *testing.T) {
		repoDir, err := FindGitRepoDir(7)
		assert.NoError(t, err, "expected no error")
		assert.NotEmptyf(t, repoDir, "expected repoDir to be non-empty")
	})

	t.Run("Do not find git repo if levels are not enough", func(t *testing.T) {
		_, err := FindGitRepoDir(1)
		assert.Error(t, err, "expected error")
		assert.Contains(t, err.Error(), "no Git repository found",
			"expected error message to state a git repository was not found")
	})
}

func TestIsRelativePath(t *testing.T) {
	t.Run("relative path", func(t *testing.T) {
		assert.NoError(t, IsRelativeE("relative/path"))
	})

	t.Run("absolute path", func(t *testing.T) {
		assert.Error(t, IsRelativeE("/absolute/path"))
	})
}
