package terragrunt

import (
	"dagger.io/dagger"
	"github.com/Excoriate/go-terradagger/pkg/config"
	"github.com/Excoriate/go-terradagger/pkg/container"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
	"github.com/Excoriate/go-terradagger/pkg/terraformcore"
)

type InitOptions struct {
	// NoColor is a flag to disable colors in terraform output
	NoColor bool
	// BackendConfigFile is the path to the backend config file
	BackendConfigFile string
	// Upgrade is a flag to upgrade the modules and plugins
	Upgrade bool
}

func Init(td *terradagger.TD, tfOpts terraformcore.TfGlobalOptions, options InitOptions, _ GlobalOptions) (*dagger.Container, container.Runtime, error) {
	IaacCfg := terraformcore.IacConfigOptions{
		Binary: config.IacToolTerragrunt,
	}

	tfIaac := terraformcore.IasC{
		Config: &IaacCfg,
	}

	return tfIaac.Init(td, tfOpts, &terraformcore.InitOptions{
		NoColor:           options.NoColor,
		BackendConfigFile: options.BackendConfigFile,
		Upgrade:           options.Upgrade,
	}, []string{})
}

func InitE(td *terradagger.TD, tfOpts terraformcore.TfGlobalOptions, options InitOptions, _ GlobalOptions) (string, error) {
	IaacCfg := terraformcore.IacConfigOptions{
		Binary: config.IacToolTerragrunt,
	}

	tfIaac := terraformcore.IasC{
		Config: &IaacCfg,
	}

	return tfIaac.InitE(td, tfOpts, &terraformcore.InitOptions{
		NoColor:           options.NoColor,
		BackendConfigFile: options.BackendConfigFile,
		Upgrade:           options.Upgrade,
	}, []string{})
}
