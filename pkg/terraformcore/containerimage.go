package terraformcore

import (
	"fmt"

	"github.com/Excoriate/go-terradagger/pkg/container"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
)

// getContainerImageCfg resolves the container image to use for the given IAC configuration and Terraform options.
func getContainerImageCfg(td *terradagger.TD, iacConfig IacConfig, tfOpts TfGlobalOptions) container.Image {
	var imageCfg container.Image
	if tfOpts.GetCustomContainerImage() != "" {
		td.Log.Warn(fmt.Sprintf("using custom container image: %s", tfOpts.GetCustomContainerImage()))
		imageCfg = container.NewImageConfig(tfOpts.GetCustomContainerImage(), tfOpts.GetTerraformVersion())
	} else {
		imageCfg = container.NewImageConfig(iacConfig.GetContainerImage(), tfOpts.GetTerraformVersion())
	}

	return imageCfg
}
