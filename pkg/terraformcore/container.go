package terraformcore

import (
	"fmt"

	"dagger.io/dagger"

	"github.com/Excoriate/go-terradagger/pkg/container"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
)

type TerraformContainerConfig interface {
	GetTfOptions() TfGlobalOptions
	GetIacConfig() IacConfig
}

type TerraformContainerConfigOptions struct {
	tfOptions TfGlobalOptions
	iacConfig IacConfig
}

func (t *TerraformContainerConfigOptions) GetTfOptions() TfGlobalOptions {
	return t.tfOptions
}

func (t *TerraformContainerConfigOptions) GetIacConfig() IacConfig {
	return t.iacConfig
}

type TerraformContainerSetup interface {
	getContainerImageCfg(td *terradagger.TD) container.Image
	getContainerRuntime(td *terradagger.TD, imageCfg container.Image) container.Runtime
	AddEnvVarsToTerraformContainer(td *terradagger.TD, runtime container.Runtime, tfContainer *dagger.Container) *dagger.Container
}

// getContainerImageCfg resolves the container image to use for the given IAC configuration and Terraform options.
func (t *TerraformContainerConfigOptions) getContainerImageCfg(td *terradagger.TD) container.Image {
	var imageCfg container.Image
	if t.tfOptions.GetCustomContainerImage() != "" {
		td.Log.Warn(fmt.Sprintf("using custom container image: %s", t.tfOptions.GetCustomContainerImage()))
		imageCfg = container.NewImageConfig(t.tfOptions.GetCustomContainerImage(), t.tfOptions.GetTerraformVersion())
	} else {
		imageCfg = container.NewImageConfig(t.iacConfig.GetContainerImage(), t.tfOptions.GetTerraformVersion())
	}

	return imageCfg
}

// getContainerCfg resolves the container configuration to use for the given IAC configuration and Terraform options.
func (t *TerraformContainerConfigOptions) getContainerRuntime(td *terradagger.TD, imageCfg container.Image) container.Runtime {
	containerCfg := container.Config{
		MountPathAbs:         td.Config.GetWorkspaceAbs(),
		Workdir:              t.tfOptions.GetModulePath(),
		ContainerImage:       imageCfg,
		KeepEntryPoint:       false,                                // This will override the container's entrypoint with the command we want to run.
		AddPrivateGitSupport: t.tfOptions.GetEnableSSHPrivateGit(), // Add support for private git repos.
	}

	return container.New(&containerCfg, td)
}

// AddEnvVarsToTerraformContainer configures the container with the appropriate environment variables.
func (t *TerraformContainerConfigOptions) AddEnvVarsToTerraformContainer(td *terradagger.TD, runtime container.Runtime, tfContainer *dagger.Container) *dagger.Container {
	tfOpts := t.GetTfOptions()

	// Mirror all host environment variables if specified.
	if tfOpts.IsMirrorAllEnvVarsFromHost() {
		return runtime.AddEnvVars(td.Config.GetHostEnvVars(), tfContainer)
	}

	// Add AWS keys from host if auto-detection is enabled.
	if tfOpts.IsAutoDetectAWSKeysFromHost() {
		tfContainer = runtime.AddEnvVars(td.Config.GetAWSEnvVars(), tfContainer)
	}

	// Inject specified environment variables by keys from the host.
	if len(tfOpts.GetEnvVarsToInjectByKeyFromHost()) > 0 {
		envVarsToInject := td.Config.GetEnvVarsByKeys(tfOpts.GetEnvVarsToInjectByKeyFromHost())
		tfContainer = runtime.AddEnvVars(envVarsToInject, tfContainer)
	}

	// Automatically detect and add TF_VAR_* environment variables from the host.
	if tfOpts.IsAutoDetectTFVarsFromHost() {
		tfContainer = runtime.AddEnvVars(td.Config.GetTerraformEnvVars(), tfContainer)
	}

	return tfContainer
}
