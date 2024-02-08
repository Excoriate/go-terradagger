package terraformcore

import (
	"dagger.io/dagger"
	"github.com/Excoriate/go-terradagger/pkg/config"
	"github.com/Excoriate/go-terradagger/pkg/container"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
)

type IacCoreCommand interface {
	Init(td *terradagger.TD, tfOpts TfGlobalOptions, options InitArgs) (*dagger.Container, container.Runtime, error)
	InitE(td *terradagger.TD, tfOpts TfGlobalOptions, options InitArgs) (string, error)
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
