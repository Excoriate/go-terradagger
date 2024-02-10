package terraformcore

import (
	"github.com/Excoriate/go-terradagger/pkg/config"

	"github.com/Excoriate/go-terradagger/pkg/terradagger"
)

type tfOptions struct {
	td      *terradagger.TD
	options *TfOptions
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
	GetEnableSSHPrivateGit() bool
	GetCustomContainerImage() string
	GetInvalidateCache() bool
	IsAutoDetectTFVarsFromHost() bool
	IsAutoDetectAWSKeysFromHost() bool
	IsMirrorAllEnvVarsFromHost() bool
	GetEnvVarsToInjectByKeyFromHost() []string
	TfGlobalValidator
}

// WithOptions is a function that returns a TfGlobalOptions
// with the options passed in. It's a constructor for the TfGlobalOptions
func WithOptions(td *terradagger.TD, o *TfOptions) TfGlobalOptions {
	return &tfOptions{
		td:      td,
		options: o, // all the options passed in here are the ones that are used in the terraform code
	}
}

func (o *tfOptions) GetModulePath() string {
	return o.options.ModulePath
}

func (o *tfOptions) GetTerraformVersion() string {
	if o.options.TerraformVersion == "" {
		return config.TerraformDefaultVersion
	}

	return o.options.TerraformVersion
}

func (o *tfOptions) GetEnableSSHPrivateGit() bool {
	return o.options.EnableSSHPrivateGit
}

func (o *tfOptions) GetCustomContainerImage() string {
	return o.options.CustomContainerImage
}

func (o *tfOptions) GetInvalidateCache() bool {
	return o.options.InvalidateCache
}

func (o *tfOptions) IsAutoDetectTFVarsFromHost() bool {
	return o.options.AutoDetectTFVarsFromHost
}

func (o *tfOptions) IsAutoDetectAWSKeysFromHost() bool {
	return o.options.AutoDetectAWSKeysFromHost
}

func (o *tfOptions) IsMirrorAllEnvVarsFromHost() bool {
	return o.options.MirrorAllEnvVarsFromHost
}

func (o *tfOptions) GetEnvVarsToInjectByKeyFromHost() []string {
	return o.options.EnvVarsToInjectByKeyFromHost
}
