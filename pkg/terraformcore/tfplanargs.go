package terraformcore

import (
	"fmt"
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/erroer"

	"github.com/Excoriate/go-terradagger/pkg/utils"
)

type PlanArgsOptions struct {
	// RefreshOnly is a flag to only refresh the state. Equivalent to
	// terraform plan - refresh-only
	RefreshOnly bool
	// TerraformVarFiles is a list of terraform var files to use
	TerraformVarFiles []string
	// Vars is a list of terraform vars to use
	Vars []TFInputVariable

	// TfGlobalOptions is a struct that contains the global options for the terraform binary
	// It implements the TfGlobalOptions interface
	TfGlobalOptions TfGlobalOptions
}

type PlanArgs interface {
	GetArgRefreshOnly() []string
	GetArgTerraformVarFiles() []string
	GetArgVars() []string

	// PlanArgsValidator is an interface for validating the plan args,
	// And also inherits from the TfArgs interface
	PlanArgsValidator
}

type PlanArgsValidator interface {
	VarFilesAreValid() error
	TfArgs
}

func (po *PlanArgsOptions) GetArgRefreshOnly() []string {
	if po.RefreshOnly {
		return []string{"-refresh-only"}
	}
	return []string{}
}

func (po *PlanArgsOptions) GetArgTerraformVarFiles() []string {
	var args []string
	for _, file := range po.TerraformVarFiles {
		args = append(args, fmt.Sprintf("-var-file=%s", file))
	}
	return args
}

func (po *PlanArgsOptions) GetArgVars() []string {
	var args []string
	for _, v := range po.Vars {
		args = append(args, fmt.Sprintf("-var '%s=%s'", v.Name, utils.EscapeValues(v.Value)))
	}
	return args
}

func (po *PlanArgsOptions) VarFilesAreValid() error {
	varFiles := po.GetArgTerraformVarFiles()

	for _, file := range varFiles {
		tfVarFilePath := filepath.Join(po.TfGlobalOptions.GetModulePathFull(), file)
		if err := TfVarFilesExistAndAreValid([]string{tfVarFilePath}); err != nil {
			return err
		}
	}

	return nil
}

func (po *PlanArgsOptions) AreValid() error {
	if err := po.VarFilesAreValid(); err != nil {
		return erroer.NewErrTerraformCoreInvalidArgumentError("the var files are not valid", err)
	}

	return nil
}
