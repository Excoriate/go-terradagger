package terraform

var (
	tfDefaultCacheDir            = ".terraform"
	tfDefaultStateFileName       = "terraform.tfstate"
	tfDefaultStateBackupFileName = "terraform.tfstate.backup"
	tfDefaultLockFileName        = ".terraform.lock.hcl"
	tfEntryPoint                 = "terraform"
	tfInitCommand                = "init"
	tfPlanCommand                = "plan"
	tfApplyCommand               = "apply"
	tfDestroyCommand             = "destroy"
	tfValidateCommand            = "validate"
)

type TfConfig interface {
	GetCacheDir() string
	GetStateFileName() string
	GetStateBackupFileName() string
	GetLockFileName() string
}

type TfCommandsConfig interface {
	GetEntryPoint() string
	GetInitCommand() string
	GetPlanCommand() string
	GetApplyCommand() string
	GetDestroyCommand() string
	GetValidateCommand() string
}

type tfCfg struct{}

func NewTFConfig() TfConfig {
	return &tfCfg{}
}

func (t *tfCfg) GetCacheDir() string {
	return tfDefaultCacheDir
}

func (t *tfCfg) GetStateFileName() string {
	return tfDefaultStateFileName
}

func (t *tfCfg) GetStateBackupFileName() string {
	return tfDefaultStateBackupFileName
}

func (t *tfCfg) GetLockFileName() string {
	return tfDefaultLockFileName
}

func (t *tfCfg) GetEntryPoint() string {
	return tfEntryPoint
}

func (t *tfCfg) GetInitCommand() string {
	return tfInitCommand
}

func (t *tfCfg) GetPlanCommand() string {
	return tfPlanCommand
}

func (t *tfCfg) GetApplyCommand() string {
	return tfApplyCommand
}

func (t *tfCfg) GetDestroyCommand() string {
	return tfDestroyCommand
}

func (t *tfCfg) GetValidateCommand() string {
	return tfValidateCommand
}

type TfCmd struct{}

func NewTfCmd() TfCommandsConfig {
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
