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
	// mirrorAllEnvVarsFromHost is a flag to scan the environment variables and inject them into the terraform code
	// The variables that'll be injected are the ones that start with TF_VAR_
	mirrorAllEnvVarsFromHost bool
	// autoDetectTFVarsFromHost is a flag to scan the environment variables and inject them into the terraform code
	autoDetectTFVarsFromHost bool
	// autoDetectAWSKeysFromHost is a flag to scan the environment variables and inject them into the terraform code
	autoDetectAWSKeysFromHost bool
	// customContainerImage is the custom image to use for the terraform container
	customContainerImage string
	// enableSSHPrivateGit is a flag to use SSH for the modules
	enableSSHPrivateGit bool
	// invalidateCache is a flag to invalidate the cache
	invalidateCache bool
	// envVarsToInjectByKeyFromHost is a slice of environment variables to inject into the container
	envVarsToInjectByKeyFromHost []string
}

type TfOptions struct {
	// ModulePath is the directory of the terraform code
	ModulePath string
	// TerraformVersion is the version of terraform to use
	TerraformVersion string
	// MirrorAllEnvVarsFromHost is a flag to scan the environment variables and inject them into the terraform code
	// The variables that'll be injected are the ones that start with TF_VAR_
	MirrorAllEnvVarsFromHost bool
	// AutoDetectTFVarsFromHost is a flag to scan the environment variables and inject them into the terraform code
	AutoDetectTFVarsFromHost bool
	// AutoDetectAWSKeysFromHost is a flag to scan the environment variables and inject them into the terraform code
	AutoDetectAWSKeysFromHost bool
	// CustomContainerImage is the custom image to use for the terraform container
	CustomContainerImage string
	// EnableSSHPrivateGit is a flag to use SSH for the modules
	EnableSSHPrivateGit bool
	// InvalidateCache is a flag to invalidate the cache
	InvalidateCache bool
	// EnvVarsToInjectByKeyFromHost is a slice of environment variables to inject into the container
	EnvVarsToInjectByKeyFromHost []string
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
	IsAutoDetectTFVarsFromHost() bool
	IsAutoDetectAWSKeysFromHost() bool
	IsMirrorAllEnvVarsFromHost() bool
	GetEnvVarsToInjectByKeyFromHost() []string
}

func WithOptions(td *terradagger.TD, o *TfOptions) TfGlobalOptions {
	return &tfOptions{
		td:                           td,
		terraformVersion:             o.TerraformVersion,
		mirrorAllEnvVarsFromHost:     o.MirrorAllEnvVarsFromHost,
		autoDetectTFVarsFromHost:     o.AutoDetectTFVarsFromHost,
		autoDetectAWSKeysFromHost:    o.AutoDetectAWSKeysFromHost,
		enableSSHPrivateGit:          o.EnableSSHPrivateGit,
		modulePath:                   o.ModulePath,
		customContainerImage:         o.CustomContainerImage,
		invalidateCache:              o.InvalidateCache,
		envVarsToInjectByKeyFromHost: o.EnvVarsToInjectByKeyFromHost,
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

func (o *tfOptions) IsAutoDetectTFVarsFromHost() bool {
	return o.autoDetectTFVarsFromHost
}

func (o *tfOptions) IsAutoDetectAWSKeysFromHost() bool {
	return o.autoDetectAWSKeysFromHost
}

func (o *tfOptions) IsMirrorAllEnvVarsFromHost() bool {
	return o.mirrorAllEnvVarsFromHost
}

func (o *tfOptions) GetEnvVarsToInjectByKeyFromHost() []string {
	return o.envVarsToInjectByKeyFromHost
}
