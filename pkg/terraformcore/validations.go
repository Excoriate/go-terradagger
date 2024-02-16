package terraformcore

import (
	"fmt"

	"github.com/Excoriate/go-terradagger/pkg/erroer"
	"github.com/Excoriate/go-terradagger/pkg/utils"
)

type TfGlobalValidator interface {
	IsModulePathValid() error
	ModulePathHasTerraformCode() error
	ModulePathHasTerragruntHCL() error
}

func (o *tfOptions) IsModulePathValid() error {
	if o.GetModulePath() == "" {
		return erroer.NewErrTerraDaggerInvalidArgumentError("the module path is empty", nil)
	}

	modulePathFull := o.GetModulePathFull()

	if err := utils.IsValidDirE(modulePathFull); err != nil {
		return erroer.NewErrTerraDaggerInvalidArgumentError(fmt.Sprintf("the module path %s is not valid", modulePathFull), err)
	}

	return nil
}

// ModulePathHasTerraformCode checks if the module path has terraformed code
// It checks for the existence of a file with the extension.tf or.hcl
func (o *tfOptions) ModulePathHasTerraformCode() error {
	modulePathFull := o.GetModulePathFull()

	return utils.DirHasContentWithCertainExtension(modulePathFull, []string{".tf"})
}

func (o *tfOptions) ModulePathHasTerragruntHCL() error {
	modulePathFull := o.GetModulePathFull()

	return utils.DirHasContentWithCertainExtension(modulePathFull, []string{".hcl"})
}
