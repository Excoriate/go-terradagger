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

func Init(td *terradagger.TD, tfOpts terraformcore.TfGlobalOptions, options InitOptions, _ terraformcore.TerragruntConfig) (*dagger.Container, container.Runtime, error) {
	tgRun := terraformcore.NewTerragruntRunner(td, tfOpts, nil)

	return tgRun.RunInit(config.IacToolTerragrunt, &terraformcore.InitArgsOptions{
		NoColor:           options.NoColor,
		BackendConfigFile: options.BackendConfigFile,
		Upgrade:           options.Upgrade,
		TfGlobalOptions:   tfOpts,
	})
}

func InitE(td *terradagger.TD, tfOpts terraformcore.TfGlobalOptions, options InitOptions, _ terraformcore.TerragruntConfig) (string, error) {
	tgRun := terraformcore.NewTerragruntRunner(td, tfOpts, nil)

	return tgRun.RunInitE(config.IacToolTerragrunt, &terraformcore.InitArgsOptions{
		NoColor:           options.NoColor,
		BackendConfigFile: options.BackendConfigFile,
		Upgrade:           options.Upgrade,
		TfGlobalOptions:   tfOpts,
	})
}
