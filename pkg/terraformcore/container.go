package terraformcore

import (
	"fmt"

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
