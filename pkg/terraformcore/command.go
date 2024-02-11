package terraformcore

import (
	"github.com/Excoriate/go-terradagger/pkg/config"
	"github.com/Excoriate/go-terradagger/pkg/erroer"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
)

var (
	tfEntryPoint      = "terraform"
	tgEntryPoint      = "terragrunt"
	tfInitCommand     = "init"
	tfPlanCommand     = "plan"
	tfApplyCommand    = "apply"
	tfDestroyCommand  = "destroy"
	tfValidateCommand = "validate"
)

type TfLifecycleCMD struct{}

type LifecycleCommandConfig interface {
	GetEntryPoint(iaacTool string) string
	GetInitCommand() string
	GetPlanCommand() string
	GetApplyCommand() string
	GetDestroyCommand() string
	GetValidateCommand() string
}

func (t *TfLifecycleCMD) GetEntryPoint(iaacTool string) string {
	if iaacTool == "" {
		return tfEntryPoint
	}

	if iaacTool == config.IacToolTerragrunt {
		return tgEntryPoint
	}

	return tfEntryPoint
}

func (t *TfLifecycleCMD) GetInitCommand() string {
	return tfInitCommand
}

func (t *TfLifecycleCMD) GetPlanCommand() string {
	return tfPlanCommand
}

func (t *TfLifecycleCMD) GetApplyCommand() string {
	return tfApplyCommand
}

func (t *TfLifecycleCMD) GetDestroyCommand() string {
	return tfDestroyCommand
}

func (t *TfLifecycleCMD) GetValidateCommand() string {
	return tfValidateCommand
}

type GetTerraformLifecycleCMDStringOptions struct {
	iacConfig        IacConfig
	lifecycleCommand string
	args             []string
}

type GenerateTFInitCMDStrOptions struct {
	iacConfig IacConfig
	initArgs  []string
}
type TfLifecycleCMDResolver interface {
	GetTerraformLifecycleCMDString(options *GetTerraformLifecycleCMDStringOptions) (string, error)
	GenerateTFInitCommandStr(options *GenerateTFInitCMDStrOptions) (string, error)
}

func (t *TfLifecycleCMD) GetTerraformLifecycleCMDString(options *GetTerraformLifecycleCMDStringOptions) (string, error) {
	if options == nil {
		return "", erroer.NewErrTerraformCoreInvalidConfigurationError("options cannot be nil", nil)
	}

	var cmdStr string
	cmdBinary := options.iacConfig.GetBinary()

	if cmdBinary == config.IacToolTerragrunt {
		cmdStr = terradagger.BuildTerragruntCommand(terradagger.BuildTerragruntCommandOptions{
			Binary:      cmdBinary,
			Command:     options.lifecycleCommand,
			CommandArgs: options.args,
		})
	} else {
		cmdStr = terradagger.BuildTerraformCommand(terradagger.BuildTerraformCommandOptions{
			Binary:      cmdBinary,
			Command:     options.lifecycleCommand,
			CommandArgs: options.args,
		})
	}
	return cmdStr, nil
}

func (t *TfLifecycleCMD) GenerateTFInitCommandStr(options *GenerateTFInitCMDStrOptions) (string, error) {
	if options == nil {
		return "", erroer.NewErrTerraformCoreInvalidConfigurationError("options cannot be nil", nil)
	}

	cmdBinary := options.iacConfig.GetBinary()
	var initStr string

	if cmdBinary == config.IacToolTerragrunt {
		initStr = terradagger.BuildTerragruntCommand(terradagger.BuildTerragruntCommandOptions{
			Binary:      cmdBinary,
			Command:     t.GetInitCommand(),
			CommandArgs: options.initArgs,
		})
	} else {
		initStr = terradagger.BuildTerraformCommand(terradagger.BuildTerraformCommandOptions{
			Binary:      cmdBinary,
			Command:     t.GetInitCommand(),
			CommandArgs: options.initArgs,
		})
	}

	return initStr, nil
}
