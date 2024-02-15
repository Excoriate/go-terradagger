package container

import (
	"testing"

	"github.com/Excoriate/go-terradagger/pkg/config"
)

func TestNewImageConfig(t *testing.T) {
	image := "custom/image"
	version := "1.0.0"
	imgConfig := NewImageConfig(image, version).(*ImageConfig)

	if imgConfig.image != image {
		t.Errorf("Expected image %s, got %s", image, imgConfig.image)
	}
	if imgConfig.version != version {
		t.Errorf("Expected version %s, got %s", version, imgConfig.version)
	}
}

func TestImageConfig_GetImageTerraform(t *testing.T) {
	customImage := "custom/terraform"
	imgConfig := NewImageConfig(customImage, "").(*ImageConfig)

	if got := imgConfig.GetImageTerraform(); got != customImage {
		t.Errorf("GetImageTerraform() = %v, want %v", got, customImage)
	}

	defaultImgConfig := NewImageConfig("", "").(*ImageConfig)
	if got := defaultImgConfig.GetImageTerraform(); got != config.TerraformDefaultImage {
		t.Errorf("GetImageTerraform() = %v, want %v", got, config.TerraformDefaultImage)
	}
}

func TestImageConfig_GetImageTerragrunt(t *testing.T) {
	customImage := "custom/terragrunt"
	imgConfig := NewImageConfig(customImage, "").(*ImageConfig)

	if got := imgConfig.GetImageTerragrunt(); got != customImage {
		t.Errorf("GetImageTerragrunt() = %v, want %v", got, customImage)
	}

	defaultImgConfig := NewImageConfig("", "").(*ImageConfig)
	if got := defaultImgConfig.GetImageTerragrunt(); got != config.TerragruntDefaultImage {
		t.Errorf("GetImageTerragrunt() = %v, want %v", got, config.TerragruntDefaultImage)
	}
}

func TestImageConfig_GetVersion(t *testing.T) {
	customVersion := "0.15.0"
	imgConfig := NewImageConfig("", customVersion).(*ImageConfig)

	if got := imgConfig.GetVersion(); got != customVersion {
		t.Errorf("GetVersion() = %v, want %v", got, customVersion)
	}

	defaultImgConfig := NewImageConfig("", "").(*ImageConfig)
	if got := defaultImgConfig.GetVersion(); got != config.DefaultImageVersion {
		t.Errorf("GetVersion() = %v, want %v", got, config.DefaultImageVersion)
	}
}

func TestImageConfig_GetTerraformContainerImage(t *testing.T) {
	customImage := "custom/terraform"
	customVersion := "0.15.0"
	expected := customImage + ":" + customVersion
	imgConfig := NewImageConfig(customImage, customVersion).(*ImageConfig)

	if got := imgConfig.GetTerraformContainerImage(); got != expected {
		t.Errorf("GetTerraformContainerImage() = %v, want %v", got, expected)
	}
}

func TestImageConfig_GetTerragruntContainerImage(t *testing.T) {
	customImage := "custom/terragrunt"
	customVersion := "0.15.0"
	expected := customImage + ":" + customVersion
	imgConfig := NewImageConfig(customImage, customVersion).(*ImageConfig)

	if got := imgConfig.GetTerragruntContainerImage(); got != expected {
		t.Errorf("GetTerragruntContainerImage() = %v, want %v", got, expected)
	}
}
