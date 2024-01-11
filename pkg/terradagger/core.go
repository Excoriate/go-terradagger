package terradagger

import (
	"context"
	"fmt"

	"github.com/Excoriate/go-terradagger/pkg/utils"

	"github.com/Excoriate/go-terradagger/pkg/commands"

	"github.com/Excoriate/go-terradagger/pkg/config"

	"github.com/Excoriate/go-terradagger/pkg/daggerio"

	"github.com/Excoriate/go-terradagger/pkg/o11y"

	"dagger.io/dagger"
)

type TD struct {
	ID string
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

type OverrideTerraDaggerFileOptions struct {
	FileName             string
	TerraDaggerID        string
	TerraDaggerConfigDir string
}

type OverrideTerraDaggerDirOptions struct {
	DirName              string
	TerraDaggerID        string
	TerraDaggerConfigDir string
}

type Core interface {
	CreateTerraDaggerContainer(options *CreateTerraDaggerContainerOptions) (*Container, error)
	OverrideTerraDaggerFile(options *OverrideTerraDaggerFileOptions) error
	OverrideTerraDaggerDir(options *OverrideTerraDaggerDirOptions) error
	ConfigureTerraDaggerContainer(options *ConfigureTerraDaggerContainerOptions) (*Container, error)
	newDaggerBackendClient(enableStderrLog bool) error
	Execute(instance *ClientInstance, options *RunOptions) error
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

func (td *TD) OverrideTerraDaggerFile(options *OverrideTerraDaggerFileOptions) error {
	if options == nil {
		return fmt.Errorf("failed to override the terradagger file, the options are nil")
	}

	if options.FileName == "" {
		return fmt.Errorf("failed to override the terradagger file, the file name is empty")
	}

	if options.TerraDaggerID == "" {
		return fmt.Errorf("failed to override the terradagger file, the terradagger id is empty")
	}

	if err := config.IsAValidTerraDaggerConfigDir(options.TerraDaggerConfigDir); err != nil {
		return fmt.Errorf("failed to override the terradagger file, the terradagger config dir is invalid: %w", err)
	}

	filePathToOverride, err := GetTerraDaggerConfigPath(&GetTerraDaggerConfigPathOptions{
		TerraDaggerID: options.TerraDaggerID,
		ConfigDir:     options.TerraDaggerConfigDir,
		FileOrDirName: options.FileName,
	})

	if err != nil {
		return fmt.Errorf("failed to override the terradagger file, the file path is invalid: %w", err)
	}

	// If the file do not exist, fail with an error.
	if err := utils.OverrideFileIfExist(filePathToOverride, true); err != nil {
		return fmt.Errorf("failed to override the terradagger file, the file %s do not exist: %w", filePathToOverride, err)
	}

	return nil
}

func (td *TD) OverrideTerraDaggerDir(options *OverrideTerraDaggerDirOptions) error {
	if options == nil {
		return fmt.Errorf("failed to override the terradagger dir, the options are nil")
	}

	if options.DirName == "" {
		return fmt.Errorf("failed to override the terradagger dir, the dir name is empty")
	}

	if options.TerraDaggerID == "" {
		return fmt.Errorf("failed to override the terradagger dir, the terradagger id is empty")
	}

	if err := config.IsAValidTerraDaggerConfigDir(options.TerraDaggerConfigDir); err != nil {
		return fmt.Errorf("failed to override the terradagger dir, the terradagger config dir is invalid: %w", err)
	}

	dirPathToOverride, err := GetTerraDaggerConfigPath(&GetTerraDaggerConfigPathOptions{
		TerraDaggerID: options.TerraDaggerID,
		ConfigDir:     options.TerraDaggerConfigDir,
		FileOrDirName: options.DirName,
	})

	if err != nil {
		return fmt.Errorf("failed to override the terradagger dir, the dir path is invalid: %w", err)
	}

	// If the dir do not exist, fail with an error.
	if err := utils.OverrideDirIfExist(dirPathToOverride, true); err != nil {
		return fmt.Errorf("failed to override the terradagger dir, the dir %s do not exist: %w", dirPathToOverride, err)
	}

	return nil
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
		ID:                      utils.GetUUID(),
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
	TransferToContainer *DataTransferToContainer
	TransferToHost      *DataTransferToHost
}

func (td *TD) Execute(instance *ClientInstance, options *RunOptions) error {
	r := NewRunner(td)

	if options == nil {
		options = &RunOptions{}
	}

	transferCfg := instance.Config.runtime.containerHostInterop

	if !transferCfg.isTransferToContainerEnabled && !transferCfg.isTransferToHostEnabled {
		return r.RunOnly(instance)
	}

	if transferCfg.isTransferToContainerEnabled {
		ic := NewContainerImporter(td)
		updatedContainer, err := ic.AddDataToImportInContainer(instance.runtimeContainer.
			DaggerContainer, transferCfg.transferToContainer)
		if err != nil {
			return err
		}

		instance.runtimeContainer.DaggerContainer = updatedContainer
	}

	if !transferCfg.isTransferToHostEnabled {
		return r.RunOnly(instance)
	}

	err := r.RunWithExport(instance, &RuntWithExportOptions{})
	if err != nil {
		return err
	}

	if !transferCfg.isBackupInHostEnabled {
		return nil
	}

	// Implement the backup manager.
	bc := NewBacker(td)
	return bc.BackupManaged(&BackupOptions{
		Files: transferCfg.backupToHost.Files,
		Dirs:  transferCfg.backupToHost.Dirs,
	})
}
