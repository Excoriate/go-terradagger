package terradagger

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/commands"

	"github.com/Excoriate/go-terradagger/pkg/config"

	"github.com/Excoriate/go-terradagger/pkg/daggerio"

	"github.com/Excoriate/go-terradagger/pkg/o11y"

	"dagger.io/dagger"
)

type TD struct {
	// Configuration interfaces
	daggerConfigClient      config.DaggerBackendConfig
	terraDaggerConfigClient config.TerraDaggerConfig
	fsResolverClient        fsResolver
	dirManagerClient        dirManager
	interopClient           interop
	containerClient         ContainerFactory
	// Implementation details, and internals.
	Logger o11y.LoggerInterface
	Ctx    context.Context
	// DaggerBackend is the backend of the dagger td.
	DaggerBackend *dagger.Client

	// Main configuration object.
	Config *Config
}

type Config struct {
	TerraDagger *config.TerraDagger
	Dagger      *config.DaggerBackend
}

type Options struct {
	Workspace                    string
	EnableStdErrForDaggerBackend bool
	EnvVars                      map[string]string
}

type Client interface {
	CreateTerraDaggerContainer(options *CreateTerraDaggerContainerOptions) (*Container, error)
	ConfigureTerraDaggerContainer(options *ConfigureTerraDaggerContainerOptions) (*Container, error)
	ResolveRunOptions(instance *ClientInstance) *RunOptions
	newDaggerBackendClient(enableStderrLog bool) error
	Run(instance *ClientInstance, options *RunOptions) error
}

type CreateTerraDaggerContainerOptions struct {
	Image   string
	Version string
}

func (td *TD) CreateTerraDaggerContainer(options *CreateTerraDaggerContainerOptions) (*Container, error) {
	if options == nil {
		return nil, fmt.Errorf("failed to createNewContainer the terradagger runtime, the options are nil")
	}

	newContainer, err := td.containerClient.createNewContainer(&CreateNewContainerOptions{
		Image:   options.Image,
		Version: options.Version,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create a terraDagger runtime")
	}

	return newContainer, nil
}

type ConfigureTerraDaggerContainerOptions struct {
	Container    *Container
	EnvVars      map[string]string
	MountDir     *dagger.Directory
	WorkDirPath  string
	ExcludedDirs []string
	ExcludeFiles []string
	Commands     commands.DaggerEngineCMDs
}

func (td *TD) ConfigureTerraDaggerContainer(options *ConfigureTerraDaggerContainerOptions) (*Container, error) {
	if options == nil {
		return nil, fmt.Errorf("failed to configure the terradagger runtime, the options are nil")
	}

	if options.Container == nil {
		return nil, fmt.Errorf("failed to configure the terradagger runtime, the runtime is nil")
	}

	if options.MountDir == nil {
		return nil, fmt.Errorf("failed to configure the terradagger runtime, the mount dir is nil")
	}

	if options.WorkDirPath == "" {
		return nil, fmt.Errorf("failed to configure the terradagger runtime, the work dir is empty")
	}

	if options.Commands == nil {
		return nil, fmt.Errorf("failed to configure the terradagger runtime, the commands are nil")
	}

	configuredContainer := options.Container

	if len(options.EnvVars) > 0 {
		configuredContainer.DaggerContainer = td.containerClient.
			withEnvVars(options.Container.DaggerContainer, options.EnvVars)
	}

	configuredContainer.DaggerContainer = td.containerClient.
		withDirs(options.Container.DaggerContainer, options.MountDir, options.WorkDirPath, options.ExcludedDirs)

	configuredContainer.DaggerContainer = td.containerClient.
		withCommands(options.Container.DaggerContainer, options.Commands)

	return options.Container, nil
}

func getDefaultOptionsIfEmpty(options *Options) {
	if options == nil {
		options = &Options{}
	}
}

func (td *TD) newDaggerBackendClient(enableStderrLog bool) error {
	if td.DaggerBackend != nil {
		return nil
	}

	bc := &daggerio.Backend{}
	daggerLogCfg := bc.ResolveDaggerLogConfig(enableStderrLog)

	c, err := bc.CreateDaggerBackend(td.Ctx, dagger.WithLogOutput(daggerLogCfg))
	if err != nil {
		return fmt.Errorf("failed to createNewContainer the dagger backend: %w", err)
	}

	td.DaggerBackend = c
	return nil
}

// New creates a new terradagger td.
// If no options are passed, the default options are used.
func New(ctx context.Context, options *Options) (*TD, error) {
	logger := o11y.NewLogger(o11y.LoggerOptions{
		EnableJSONHandler: true,
		EnableStdError:    true,
	})

	td := &TD{
		Logger:                  logger,
		Ctx:                     ctx,
		daggerConfigClient:      config.NewDaggerBackendConfigClient(logger),
		terraDaggerConfigClient: config.NewTerraDaggerConfig(logger),
		Config: &Config{
			TerraDagger: &config.TerraDagger{},
			Dagger:      &config.DaggerBackend{},
		},
	}

	// Resolving the main configuration.
	getDefaultOptionsIfEmpty(options)

	terraDaggerCfg, err := td.terraDaggerConfigClient.Configure(&config.APIParams{
		Workspace:                    options.Workspace,
		EnableStdErrForDaggerBackend: options.EnableStdErrForDaggerBackend,
		EnvVars:                      options.EnvVars,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to configure the terradagger td: %w", err)
	}

	td.Config.TerraDagger = terraDaggerCfg

	// Resolving the Dagger backend configuration.
	daggerBackendCfg, err := td.daggerConfigClient.Configure()
	if err != nil {
		return nil, fmt.Errorf("failed to configure the dagger backend: %w", err)
	}

	td.Config.Dagger = daggerBackendCfg

	// Adding the required clients.
	td.fsResolverClient = newDirResolverClient(td)
	td.dirManagerClient = newDirManagerClient(td)
	td.interopClient = newInteropClient(td)
	td.containerClient = newContainerClient(td)

	// Injecting the dagger backend.
	if err := td.newDaggerBackendClient(options.EnableStdErrForDaggerBackend); err != nil {
		return nil, fmt.Errorf("failed to inject the dagger backend: %w", err)
	}

	logger.Info("TerraDagger td started successfully.")
	return td, nil
}

type RunOptions struct {
	CopyFilesToContainer []string
	CopyDirsToContainer  []string
	CopyFilesToHost      []string
	CopyDirsToHost       []string
}

func (td *TD) ResolveRunOptions(instance *ClientInstance) *RunOptions {
	if instance == nil {
		return nil
	}

	return &RunOptions{
		CopyFilesToContainer: instance.Config.runtime.containerHostInterop.copyFilesToContainer,
		CopyDirsToContainer:  instance.Config.runtime.containerHostInterop.copyDirsToContainer,
		CopyFilesToHost:      instance.Config.runtime.containerHostInterop.copyFilesToHost,
		CopyDirsToHost:       instance.Config.runtime.containerHostInterop.copyDirsToHost,
	}
}

func (td *TD) Run(instance *ClientInstance, options *RunOptions) error {
	if instance == nil {
		return fmt.Errorf("failed to run the terradagger instance, the instance is nil")
	}

	// workDir, err := instance.runtimeContainer.DaggerContainer.Workdir(td.Ctx)
	// if err != nil {
	// 	return err
	// }
	//
	// entriesInMountPath, err := instance.runtimeContainer.DaggerContainer.Directory("/mnt/test-data/terraform/root-module-1").Entries(td.Ctx)
	// if err != nil {
	// 	return err
	// }

	// td.Logger.Info("Entries in mount path", "entries", entriesInMountPath)

	// td.Logger.Info("Workdir", "workdir", workDir)

	if options == nil {
		_, runErr := instance.runtimeContainer.DaggerContainer.Stdout(td.Ctx)
		if runErr != nil {
			return runErr
		}

		return nil
	}

	if len(options.CopyFilesToHost) > 0 {
		for _, file := range options.CopyFilesToHost {
			fileName := filepath.Base(file)
			fileInHostPath := filepath.Join(instance.Config.Paths.ExportPath, fileName)

			if _, err := instance.runtimeContainer.DaggerContainer.File(file).
				Export(td.Ctx, fileInHostPath); err != nil {
				td.Logger.Error("Failed to export file", "file", file, "error", err)
			}
		}
	}

	if len(options.CopyDirsToHost) > 0 {
		for _, dir := range options.CopyDirsToHost {
			dirName := filepath.Base(dir)
			dirInHostPath := filepath.Join(instance.Config.Paths.ExportPath, dirName)

			if _, err := instance.runtimeContainer.DaggerContainer.Directory(dir).
				Export(td.Ctx, dirInHostPath); err != nil {
				td.Logger.Error("Failed to export directory", "directory", dir, "error", err)
			}
		}
	}

	if len(options.CopyDirsToContainer) > 0 {
		for _, dir := range options.CopyDirsToContainer {
			dirName := filepath.Base(dir)
			dirInContainerPath := filepath.Join(instance.Config.Paths.WorkDirPathDagger, dirName)
			dirAsDaggerFormat := td.DaggerBackend.Host().Directory(dir)

			_, err := instance.runtimeContainer.DaggerContainer.WithDirectory(dirInContainerPath, dirAsDaggerFormat).Stdout(td.Ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
