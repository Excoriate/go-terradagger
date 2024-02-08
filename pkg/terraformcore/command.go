package terraformcore

var (
	tfEntryPoint      = "terraform"
	tfInitCommand     = "init"
	tfPlanCommand     = "plan"
	tfApplyCommand    = "apply"
	tfDestroyCommand  = "destroy"
	tfValidateCommand = "validate"
)

type TfCmd struct{}

type CommandConfig interface {
	GetEntryPoint() string
	GetInitCommand() string
	GetPlanCommand() string
	GetApplyCommand() string
	GetDestroyCommand() string
	GetValidateCommand() string
}

func NewTerraformCommandConfig() CommandConfig {
	return &TfCmd{}
}

func (t *TfCmd) GetEntryPoint() string {
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
