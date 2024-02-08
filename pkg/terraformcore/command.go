package terraformcore

import "github.com/Excoriate/go-terradagger/pkg/config"

var (
	tfEntryPoint      = "terraform"
	tgEntryPoint      = "terragrunt"
	tfInitCommand     = "init"
	tfPlanCommand     = "plan"
	tfApplyCommand    = "apply"
	tfDestroyCommand  = "destroy"
	tfValidateCommand = "validate"
)

type TfCmd struct{}

type CommandConfig interface {
	GetEntryPoint(iaacTool string) string
	GetInitCommand() string
	GetPlanCommand() string
	GetApplyCommand() string
	GetDestroyCommand() string
	GetValidateCommand() string
}

func NewTerraformCommandConfig() CommandConfig {
	return &TfCmd{}
}

func (t *TfCmd) GetEntryPoint(iaacTool string) string {
	if iaacTool == "" {
		return tfEntryPoint
	}

	if iaacTool == config.IacToolTerragrunt {
		return tgEntryPoint
	}

	return tfEntryPoint
}

func (t *TfCmd) GetInitCommand() string {
	return tfInitCommand
}

func (t *TfCmd) GetPlanCommand() string {
	return tfPlanCommand
}

func (t *TfCmd) GetApplyCommand() string {
	return tfApplyCommand
}

func (t *TfCmd) GetDestroyCommand() string {
	return tfDestroyCommand
}

func (t *TfCmd) GetValidateCommand() string {
	return tfValidateCommand
}
