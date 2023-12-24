package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Adjust the import path according to your project's structure

func TestExcludedDirsDefault(t *testing.T) {
	expectedDirs := []string{".git/", ".gitignore", ".terra-dagger/", "dist/**", "node_modules/**"}
	assert.Equal(t, expectedDirs, ExcludedDirsDefault, "Default excluded directories should match expected values")
}

func TestExcludedDirsTerraform(t *testing.T) {
	expectedDirs := []string{".terraform/**", ".terraform-lock.hcl", ".terraform-cache/**"}
	assert.Equal(t, expectedDirs, ExcludedDirsTerraform, "Terraform excluded directories should match expected values")
}
