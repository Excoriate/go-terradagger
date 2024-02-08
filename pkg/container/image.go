package container

import (
	"fmt"

	"github.com/Excoriate/go-terradagger/pkg/config"
)

type Image interface {
	GetImageTerraform() string
	GetImageTerragrunt() string
	GetImageDefaultTerraform() string
	GetImageDefaultTerragrunt() string
	GetVersion() string
	GetVersionDefault() string
	GetTerraformContainerImage() string
	GetTerragruntContainerImage() string
}

type ImageConfig struct {
	image   string
	version string
}

func NewImageConfig(image, version string) Image {
	return &ImageConfig{
		image:   image,
		version: version,
	}
}

func (o *ImageConfig) GetImageTerraform() string {
	if o.image == "" {
		return o.GetImageDefaultTerraform()
	}

	return o.image
}

func (o *ImageConfig) GetImageTerragrunt() string {
	if o.image == "" {
		return o.GetImageDefaultTerragrunt()
	}

	return o.image
}

func (o *ImageConfig) GetImageDefaultTerraform() string {
	return config.TerraformDefaultImage
}

func (o *ImageConfig) GetImageDefaultTerragrunt() string {
	return config.TerragruntDefaultImage
}

func (o *ImageConfig) GetVersionDefault() string {
	return config.DefaultImageVersion
}

func (o *ImageConfig) GetVersion() string {
	if o.version == "" {
		return o.GetVersionDefault()
	}

	return o.version
}

func (o *ImageConfig) GetTerraformContainerImage() string {
	return fmt.Sprintf("%s:%s", o.GetImageTerraform(), o.GetVersion())
}

func (o *ImageConfig) GetTerragruntContainerImage() string {
	return fmt.Sprintf("%s:%s", o.GetImageTerragrunt(), o.GetVersion())
}
