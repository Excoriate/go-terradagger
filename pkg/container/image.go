package container

import "fmt"

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
	return defaultTerraformImage
}

func (o *ImageConfig) GetImageDefaultTerragrunt() string {
	return defaultTerragruntImage
}

func (o *ImageConfig) GetVersionDefault() string {
	return defaultVersion
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
