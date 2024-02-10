package terraformcore

import (
	"fmt"
	"path/filepath"

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

	srcAbsolute := o.td.Config.GetWorkspace()
	modulePathFull := filepath.Join(srcAbsolute, o.GetModulePath())

	if err := utils.IsValidDirE(modulePathFull); err != nil {
		return erroer.NewErrTerraDaggerInvalidArgumentError(fmt.Sprintf("the module path %s is not valid", modulePathFull), err)
	}

	return nil
}

// ModulePathHasTerraformCode checks if the module path has terraform code
// It checks for the existence of a file with the extension.tf or.hcl
func (o *tfOptions) ModulePathHasTerraformCode() error {
	srcAbsolute := o.td.Config.GetWorkspace()
	modulePathFull := filepath.Join(srcAbsolute, o.GetModulePath())

	return utils.DirHasContentWithCertainExtension(modulePathFull, []string{".tf"})
}

func (o *tfOptions) ModulePathHasTerragruntHCL() error {
	srcAbsolute := o.td.Config.GetWorkspace()
	modulePathFull := filepath.Join(srcAbsolute, o.GetModulePath())

	return utils.DirHasContentWithCertainExtension(modulePathFull, []string{".hcl"})
}
