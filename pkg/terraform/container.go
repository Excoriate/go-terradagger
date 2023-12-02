package terraform

const defaultImage = "hashicorp/terraform"
const defaultVersion = "1.6.1"

func resolveTerraformImage(options *Options) string {
	if options.TerraformCustomContainerImage != "" {
		return options.TerraformCustomContainerImage
	}

	return defaultImage
}

func resolveTerraformVersion(options *Options) string {
	if options.TerraformVersion != "" {
		return options.TerraformVersion
	}

	return defaultVersion
}
