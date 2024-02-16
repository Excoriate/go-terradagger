package terraformcore

import (
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/erroer"
	"github.com/Excoriate/go-terradagger/pkg/utils"
)

type InitArgsOptions struct {
	// NoColor is a flag to disable colors in terraform output
	NoColor bool
	// BackendConfigFile is the path to the backend config file
	BackendConfigFile string
	// Upgrade is a flag to upgrade the modules and plugins
	Upgrade bool

	// TfGlobalOptions is a struct that contains the global options for the terraform binary
	// It implements the TfGlobalOptions interface
	TfGlobalOptions TfGlobalOptions
}

type InitArgs interface {
	GetArgNoColor() []string
	GetArgNoColourValue() bool
	GetArgBackendConfigFile() []string
	GetArgBackendConfigFileValue() string
	GetArgUpgrade() []string
	GetArgUpgradeValue() bool

	// InitArgsValidator is an interface for validating the init args,
	// And also inherits from the TfArgs interface
	InitArgsValidator
}

type InitArgsValidator interface {
	BackendFileIsValid() error
	TfArgs
}

func (ti *InitArgsOptions) GetArgNoColor() []string {
	arg := []string{"-no-color"}
	if ti.NoColor {
		return arg
	}

	return []string{}
}

func (ti *InitArgsOptions) GetArgNoColourValue() bool {
	return ti.NoColor
}

func (ti *InitArgsOptions) GetArgBackendConfigFile() []string {
	arg := []string{"-backend-config", ti.BackendConfigFile}
	if ti.BackendConfigFile != "" {
		return arg
	}

	return []string{}
}

func (ti *InitArgsOptions) GetArgBackendConfigFileValue() string {
	return ti.BackendConfigFile
}

func (ti *InitArgsOptions) GetArgUpgrade() []string {
	arg := []string{"-upgrade"}
	if ti.Upgrade {
		return arg
	}

	return []string{}
}

func (ti *InitArgsOptions) GetArgUpgradeValue() bool {
	return ti.Upgrade
}

func (ti *InitArgsOptions) BackendFileIsValid() error {
	beCfgFile := ti.GetArgBackendConfigFileValue()

	if beCfgFile == "" {
		return nil
	}

	beCfgFilePath := filepath.Join(ti.TfGlobalOptions.GetModulePathFull(), beCfgFile)

	if err := utils.IsValidFileE(beCfgFilePath); err != nil {
		return erroer.NewErrTerraformCoreInvalidArgumentError("the backend file is not valid", err)
	}

	return nil
}

func (ti *InitArgsOptions) AreValid() error {
	if err := ti.BackendFileIsValid(); err != nil {
		return erroer.NewErrTerraformCoreInvalidArgumentError("the backend file is not valid", err)
	}

	return nil
}
