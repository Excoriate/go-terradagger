package terraform

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

func Init(td *terradagger.TD, tfOpts terraformcore.TfGlobalOptions, options InitOptions) (*dagger.Container, container.Runtime, error) {
	tfRun := terraformcore.NewTerraformRunner(td, tfOpts)

	return tfRun.RunInit(config.IacToolTerraform, &terraformcore.InitArgsOptions{
		NoColor:           options.NoColor,
		BackendConfigFile: options.BackendConfigFile,
		Upgrade:           options.Upgrade,
		TfGlobalOptions:   tfOpts,
	})
}

func InitE(td *terradagger.TD, tfOpts terraformcore.TfGlobalOptions, options InitOptions) (string, error) {
	tfRun := terraformcore.NewTerraformRunner(td, tfOpts)

	return tfRun.RunInitE(config.IacToolTerraform, &terraformcore.InitArgsOptions{
		NoColor:           options.NoColor,
		BackendConfigFile: options.BackendConfigFile,
		Upgrade:           options.Upgrade,
		TfGlobalOptions:   tfOpts,
	})
}
