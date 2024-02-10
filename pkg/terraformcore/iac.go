package terraformcore

import (
	"dagger.io/dagger"
	"github.com/Excoriate/go-terradagger/pkg/config"
	"github.com/Excoriate/go-terradagger/pkg/container"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
)

type IacLifeCycleCommand interface {
	Init(td *terradagger.TD, tfOpts TfGlobalOptions, options InitArgs, extraArgs []string) (*dagger.Container, container.Runtime, error)
	InitE(td *terradagger.TD, tfOpts TfGlobalOptions, options InitArgs, extraArgs []string) (string, error)
	Plan(td *terradagger.TD, tfOpts TfGlobalOptions, options PlanArgs, extraArgs []string) (*dagger.Container, container.Runtime, error)
	PlanE(td *terradagger.TD, tfOpts TfGlobalOptions, options PlanArgs, extraArgs []string) (string, error)
	Apply(td *terradagger.TD, tfOpts TfGlobalOptions, options ApplyArgs, extraArgs []string) (*dagger.Container, container.Runtime, error)
	ApplyE(td *terradagger.TD, tfOpts TfGlobalOptions, options ApplyArgs, extraArgs []string) (string, error)
	Destroy(td *terradagger.TD, tfOpts TfGlobalOptions, options DestroyArgs, extraArgs []string) (*dagger.Container, container.Runtime, error)
	DestroyE(td *terradagger.TD, tfOpts TfGlobalOptions, options DestroyArgs, extraArgs []string) (string, error)
}

type IacConfigOptions struct {
	Binary         string // Either terragrunt, or terraform.
	ContainerImage string
	Version        string
}

type IacConfig interface {
	GetBinary() string
	GetContainerImage() string
	GetVersion() string
}

func (i *IacConfigOptions) GetBinary() string {
	if i.Binary == "" {
		return config.IacToolTerraform
	}

	return i.Binary
}

func (i *IacConfigOptions) GetContainerImage() string {
	if i.ContainerImage != "" {
		return i.ContainerImage
	}

	binary := i.GetBinary()
	if binary == config.IacToolTerragrunt {
		return config.TerragruntDefaultImage
	}

	return config.TerraformDefaultImage
}

func (i *IacConfigOptions) GetVersion() string {
	if i.Version != "" {
		return i.Version
	}

	binary := i.GetBinary()
	if binary == config.IacToolTerragrunt {
		return config.TerragruntDefaultVersion
	}

	return config.TerraformDefaultVersion
}

type IasC struct {
	Config IacConfig
}
