package terraform

import (
	"fmt"
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/utils"
)

var allowedTerraformFileExtensions = []string{".tf", ".hcl", ".json", ".tfvars"}

func getTerraformFileRelativePath(rootDir, terraformDir, file string) (string, error) {
	// If the file is already an absolute path, return it.
	if filepath.IsAbs(file) {
		return file, nil
	}

	// Otherwise, join the root dir and terraform dir to the file.
	terraformFilePath := filepath.Join(rootDir, terraformDir, file)

	var fileHasAllowedExtension bool
	for _, allowedExtension := range allowedTerraformFileExtensions {
		if filepath.Ext(terraformFilePath) == allowedExtension {
			fileHasAllowedExtension = true
			break
		}
	}

	if !fileHasAllowedExtension {
		return "", fmt.Errorf("the terraform file %s does not have a valid extension", terraformFilePath)
	}

	// Check if the file exists.
	if err := utils.FileExistE(terraformFilePath); err != nil {
		return "", fmt.Errorf("the terraform file %s does not exist", terraformFilePath)
	}

	return terraformFilePath, nil
}
