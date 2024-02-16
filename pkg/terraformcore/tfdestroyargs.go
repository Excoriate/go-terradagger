package terraformcore

import (
	"fmt"
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/erroer"

	"github.com/Excoriate/go-terradagger/pkg/utils"
)

type DestroyArgsOptions struct {
	// RefreshOnly is a flag to only refresh the state. Equivalent to
	// terraform plan -refresh-only
	RefreshOnly bool
	// TerraformVarFiles is a list of terraform var files to use
	TerraformVarFiles []string
	// Vars is a list of terraform vars to use
	Vars []TFInputVariable
	// AutoApprove is a flag to auto approve the plan
	AutoApprove bool

	// TfGlobalOptions is a struct that contains the global options for the terraform binary
	// It implements the TfGlobalOptions interface
	TfGlobalOptions TfGlobalOptions
}

type DestroyArgs interface {
	GetArgRefreshOnly() []string
	GetArgRefreshOnlyValue() bool
	GetArgTerraformVarFiles() []string
	GetArgTerraformVarFilesValue() []string
	GetArgVars() []string
	GetArgVarsValue() []TFInputVariable
	GetArgAutoApprove() []string
	GetArgAutoApproveValue() bool

	// DestroyArgsValidator is an interface for validating the destroy args,
	// And also inherits from the TfArgs interface
	DestroyArgsValidator
}

type DestroyArgsValidator interface {
	VarFilesAreValid() error
	TfArgs
}

func (po *DestroyArgsOptions) GetArgRefreshOnly() []string {
	if po.RefreshOnly {
		return []string{"-refresh-only"}
	}
	return []string{}
}

func (po *DestroyArgsOptions) GetArgRefreshOnlyValue() bool {
	return po.RefreshOnly
}

func (po *DestroyArgsOptions) GetArgTerraformVarFiles() []string {
	var args []string
	for _, file := range po.TerraformVarFiles {
		args = append(args, fmt.Sprintf("-var-file=%s", file))
	}
	return args
}

func (po *DestroyArgsOptions) GetArgTerraformVarFilesValue() []string {
	return po.TerraformVarFiles
}

func (po *DestroyArgsOptions) GetArgVars() []string {
	var args []string
	for _, v := range po.Vars {
		args = append(args, fmt.Sprintf("-var '%s=%s'", v.Name, utils.EscapeValues(v.Value)))
	}
	return args
}

func (po *DestroyArgsOptions) GetArgVarsValue() []TFInputVariable {
	return po.Vars
}

func (po *DestroyArgsOptions) GetArgAutoApprove() []string {
	if po.AutoApprove {
		return []string{"-auto-approve"}
	}
	return []string{}
}

func (po *DestroyArgsOptions) GetArgAutoApproveValue() bool {
	return po.AutoApprove
}

func (po *DestroyArgsOptions) VarFilesAreValid() error {
	varFiles := po.GetArgTerraformVarFilesValue()

	for _, file := range varFiles {
		tfVarFilePath := filepath.Join(po.TfGlobalOptions.GetModulePathFull(), file)
		if err := TfVarFilesExistAndAreValid([]string{tfVarFilePath}); err != nil {
			return err
		}
	}

	return nil
}

func (po *DestroyArgsOptions) AreValid() error {
	if err := po.VarFilesAreValid(); err != nil {
		return erroer.NewErrTerraformCoreInvalidArgumentError("the var files are not valid", err)
	}

	return nil
}
