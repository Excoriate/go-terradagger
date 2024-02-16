package terraformcore

import (
	"fmt"

	"github.com/Excoriate/go-terradagger/pkg/config"

	"github.com/Excoriate/go-terradagger/pkg/utils"
)

// TfVarFilesExistAndAreValid checks if the terraform var files exist and are valid
// It checks if the files exist, have the correct extension and are valid files
func TfVarFilesExistAndAreValid(varFiles []string) error {
	cfg := config.Options{}

	for _, file := range varFiles {
		if err := utils.IsValidFileE(file); err != nil {
			return fmt.Errorf("the terraform var file %s is not valid: %w", file, err)
		}

		// Check if the file has the correct extension
		if err := utils.FileHasExtension(file, cfg.GetTfVarsExtension()); err != nil {
			return fmt.Errorf("the terraform var file %s does not have the correct extension: %w", file, err)
		}
	}

	return nil
}
