package terraformcore

import (
	"fmt"

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
}

type DestroyArgs interface {
	GetArgRefreshOnly() []string
	GetArgTerraformVarFiles() []string
	GetArgVars() []string
	GetArgAutoApprove() []string
}

func (po *DestroyArgsOptions) GetArgRefreshOnly() []string {
	if po.RefreshOnly {
		return []string{"-refresh-only"}
	}
	return []string{}
}

func (po *DestroyArgsOptions) GetArgTerraformVarFiles() []string {
	var args []string
	for _, file := range po.TerraformVarFiles {
		args = append(args, fmt.Sprintf("-var-file=%s", file))
	}
	return args
}

func (po *DestroyArgsOptions) GetArgVars() []string {
	var args []string
	for _, v := range po.Vars {
		args = append(args, fmt.Sprintf("-var '%s=%s'", v.Name, utils.EscapeValues(v.Value)))
	}
	return args
}

func (po *DestroyArgsOptions) GetArgAutoApprove() []string {
	if po.AutoApprove {
		return []string{"-auto-approve"}
	}
	return []string{}
}
