package terraformcore

import (
	"fmt"
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/config"

	"github.com/Excoriate/go-terradagger/pkg/erroer"
	"github.com/Excoriate/go-terradagger/pkg/utils"

	"github.com/Excoriate/go-terradagger/pkg/terradagger"
)

type tfOptions struct {
	td *terradagger.TD
	// modulePath is the directory of the terraform code
	modulePath string
	// terraformVersion is the version of terraform to use
	terraformVersion string
	// envVarsAutoInjectFromHost is a flag to scan the environment variables and inject them into the terraform code
	// The variables that'll be injected are the ones that start with TF_VAR_
	envVarsAutoInjectFromHost bool
	// envVarsAutoInjectTFVars is a flag to scan the environment variables and inject them into the terraform code
	envVarsAutoInjectTFVars bool
	// customContainerImage is the custom image to use for the terraform container
	customContainerImage string
	// enableSSHPrivateGit is a flag to use SSH for the modules
	enableSSHPrivateGit bool
	// invalidateCache is a flag to invalidate the cache
	invalidateCache bool
}

type TfOptions struct {
	// ModulePath is the directory of the terraform code
	ModulePath string
	// TerraformVersion is the version of terraform to use
	TerraformVersion string
	// EnvVarsAutoInjectFromHost is a flag to scan the environment variables and inject them into the terraform code
	// The variables that'll be injected are the ones that start with TF_VAR_
	EnvVarsAutoInjectFromHost bool
	// EnvVarsAutoInjectTFVars is a flag to scan the environment variables and inject them into the terraform code
	EnvVarsAutoInjectTFVars bool
	// CustomContainerImage is the custom image to use for the terraform container
	CustomContainerImage string
	// EnableSSHPrivateGit is a flag to use SSH for the modules
	EnableSSHPrivateGit bool
	// InvalidateCache is a flag to invalidate the cache
	InvalidateCache bool
}

type TfGlobalOptions interface {
	GetModulePath() string
	GetTerraformVersion() string
	IsModulePathValid() error
	ModulePathHasTerraformCode() error
	ModulePathHasTerragruntHCL() error
	GetEnableSSHPrivateGit() bool

	GetCustomContainerImage() string
	GetInvalidateCache() bool
}

func WithOptions(td *terradagger.TD, o *TfOptions) TfGlobalOptions {
	return &tfOptions{
		td:                        td,
		terraformVersion:          o.TerraformVersion,
		envVarsAutoInjectFromHost: o.EnvVarsAutoInjectFromHost,
		envVarsAutoInjectTFVars:   o.EnvVarsAutoInjectTFVars,
		enableSSHPrivateGit:       o.EnableSSHPrivateGit,
		modulePath:                o.ModulePath,
		customContainerImage:      o.CustomContainerImage,
		invalidateCache:           o.InvalidateCache,
	}
}

func (o *tfOptions) GetModulePath() string {
	return o.modulePath
}

func (o *tfOptions) GetTerraformVersion() string {
	if o.terraformVersion == "" {
		return config.TerraformDefaultVersion
	}

	return o.terraformVersion
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

func (o *tfOptions) GetEnableSSHPrivateGit() bool {
	return o.enableSSHPrivateGit
}

func (o *tfOptions) GetCustomContainerImage() string {
	return o.customContainerImage
}

func (o *tfOptions) GetInvalidateCache() bool {
	return o.invalidateCache
}
