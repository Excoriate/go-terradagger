package terradagger

import (
	"fmt"
	"path/filepath"

	"dagger.io/dagger"

	"github.com/Excoriate/go-terradagger/pkg/utils"

	"github.com/Excoriate/go-terradagger/pkg/commands"
)

// ClientOptions are the options for the terradagger instance.
type ClientOptions struct {
	// ContainerOptions are the options for the runtime that will be used to run the terradagger.
	ContainerOptions *InstanceContainerOptions
	// ExcludeOptions are the options for the files and directories that should be excluded from the terradagger.
	ExcludeOptions *ExcludeOptions
	// WorkDirPath is the path where the terradagger will be executed in the runtime.
	WorkDirPath string
	// TerraDaggerCMDs     [][]string
	TerraDaggerCMDs commands.DaggerEngineCMDs
	// EnvVarOptions are the options for the env vars that should be passed to the terradagger.
	EnvVarOptions *EnvVarOptions
	// ExportFromContainer are the options for the files and directories that should be required from the terradagger.
	ExportFromContainer *ExportFromContainerOptions
	// ImportToContainer are the options for the files and directories that should be required from the terradagger.
	ImportToContainer *ImportToContainerOptions
	// PreRequisites are the options for the files and directories that should be required from the terradagger.
	PreRequisites *PreRequisites
}

type PreRequisites struct {
	WorkDir  *Requisites
	MountDir *Requisites
}

type Requisites struct {
	RequiredFiles          []string
	RequiredDirs           []string
	RequiredFileExtensions []string
}

type ImportToContainerOptions struct {
	LookupFromWorkDir bool // It'll look for the files in the work dir, and not in the import dir.
	FileNames         []string
	DirNames          []string
}

type ExportFromContainerOptions struct {
	FileNames                  []string
	DirNames                   []string
	FailIfNotExistInContainer  bool
	OverrideIfExistInHost      bool
	OverrideIfExistInContainer bool
}

type EnvVarOptions struct {
	EnvVars                   map[string]string
	CopyEnvVarsFromHostByKeys []string
	MirrorEnvVarsFromHost     bool
}

type ImportExportOptions struct {
	ImportFileNames []string
	ImportDirNames  []string
	ExportFileNames []string
	ExportDirNames  []string
}

// InstanceConfig is the configuration of the terradagger instance.
type InstanceConfig struct {
	ID            string
	Paths         *InstancePaths
	ClientOptions *ClientOptions
	runtime       *runtimeConfig
}

type runtimeConfig struct {
	image                string
	version              string
	mountDir             *dagger.Directory
	containerHostInterop *containerHostInteropConfig
	envVars              map[string]string
	excludeDirs          []string
	excludeFiles         []string
	commands             commands.DaggerEngineCMDs
}

// InstancePaths is the Paths of the terradagger instance.
type InstancePaths struct {
	// TerraDagger is the path of the terradagger directory .terradagger
	TerraDagger string
	// ExportPath is the path of the export directory in the host,
	// formatted as: .terradagger/ID/export
	// TerraDaggerAbs is the absolute path of the terradagger directory .terradagger
	TerraDaggerAbs string
	ExportPath     string
	// ImportPath is the path of the import directory in the host,
	// formatted as: .terradagger/ID/import
	// ExportPathAbs is the absolute path of the export directory in the host,
	ExportPathAbs string
	// ImportPath is the path of the import directory in the host,
	ImportPath string
	// ImportPathAbs is the absolute path of the import directory in the host,
	ImportPathAbs string
	// CachePath is the path of the CachePath directory in the host,
	// formatted as: .terradagger/ID/.terradagger-CachePath
	CachePath string
	// MountPath is the path of the mount directory in the runtime,
	// formatted as: /mnt/MountPath
	// CachePathAbs is the absolute path of the CachePath directory in the host,
	CachePathAbs string
	// MountPath is the path of the mount directory in the runtime,
	MountPath string
	// MountPathAbs is the absolute path of the mount directory in the runtime,
	MountPathAbs string
	// WorkDirPath is the path of the work directory in the runtime,
	// formatted as: /mnt/MountPath/WorkDirPath
	WorkDirPath string
	// WorkDirPathDagger is the path of the work directory in the runtime, which includes
	// the dagger mount prefix, formatted as: /mnt/MountPath/WorkDirPath
	WorkDirPathDagger string
	// WorkDirPathAbs is the absolute path of the work directory in the runtime,
	WorkDirPathAbs string
	// mountPrefix is the prefix of the mount directory in the runtime,
	mountPrefix string
}

type InstanceContainerOptions struct {
	Image   string
	Version string
}

type ExcludeOptions struct {
	ExcludedDirs []string
	ExcludeFiles []string
}

type containerHostInteropConfig struct {
	copyFilesToHost       []string // Equivalent to export.
	copyDirsToHost        []string // Equivalent to export.
	copyToContainerConfig []*DataTransfer
}

type Instance interface {
	Validate(options *ClientOptions) error
	Configure(options *ClientOptions) (*InstanceConfig, error)
	PrepareInstance(cfg *InstanceConfig) (*ClientInstance, error)
}

type ClientInstance struct {
	ID               string
	Config           *InstanceConfig
	td               *TD
	runtimeContainer *Container
}

type InstanceImpl struct {
	td *TD
}

func NewInstance(td *TD) Instance {
	return &InstanceImpl{
		td: td,
	}
}

type ClientValidationError struct {
	ErrWrapped error
	Details    string
}

type ClientConfigurationError struct {
	ErrWrapped error
	Details    string
}

const clientInstanceValidationErrPrefix = "the terradagger client (instance) failed to be validated"
const clientInstanceConfigurationErrPrefix = "the terradagger client (instance) failed to be configured"

func (e *ClientValidationError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", clientInstanceValidationErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", clientInstanceValidationErrPrefix, e.Details, e.ErrWrapped.Error())
}

func (e *ClientConfigurationError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", clientInstanceConfigurationErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", clientInstanceConfigurationErrPrefix, e.Details, e.ErrWrapped.Error())
}

func (i *InstanceImpl) Validate(options *ClientOptions) error {
	if options == nil {
		return fmt.Errorf("failed to Validate the terradagger instance, the options are nil")
	}

	instanceValidator := newClientInstanceValidator(i)
	mountPath := i.td.Config.TerraDagger.Paths.Workspace.SRC
	mountPathAbs := i.td.Config.TerraDagger.Paths.Workspace.SRCAbsolute

	if err := instanceValidator.IsEnvVarOptionsValid(options.EnvVarOptions); err != nil {
		return &ClientValidationError{
			ErrWrapped: err,
			Details:    "the env var options are invalid",
		}
	}

	if err := instanceValidator.IsContainerOptionsValid(options.ContainerOptions); err != nil {
		return &ClientValidationError{
			ErrWrapped: err,
			Details:    "the container options are invalid",
		}
	}

	if len(options.TerraDaggerCMDs) == 0 {
		return fmt.Errorf("failed to Validate the terradagger instance, the terradagger commands are empty")
	}

	mountDirValidator := NewMountDirValidator(i.td)

	if err := mountDirValidator.IsWorkDirValid(mountPath, options.WorkDirPath); err != nil {
		return &ClientValidationError{
			ErrWrapped: err,
			Details:    "the work dir path is invalid",
		}
	}

	if err := instanceValidator.IsExcludedOptionsValid(mountPathAbs, options.ExcludeOptions); err != nil {
		return &ClientValidationError{
			ErrWrapped: err,
			Details:    "the excluded options are invalid",
		}
	}

	if err := instanceValidator.IsImportToContainerOptionsValid(mountPath, options.WorkDirPath, options.ImportToContainer); err != nil {
		return &ClientValidationError{
			ErrWrapped: err,
			Details:    "the import to container options are invalid",
		}
	}

	if err := instanceValidator.IsPreRequisitesValid(mountPath, options.WorkDirPath,
		options.PreRequisites.MountDir); err != nil {
		return &ClientValidationError{
			ErrWrapped: err,
			Details:    "the pre-requisites for the mountDir are invalid",
		}
	}

	if err := instanceValidator.IsPreRequisitesValid(mountPath, options.WorkDirPath,
		options.PreRequisites.WorkDir); err != nil {
		return &ClientValidationError{
			ErrWrapped: err,
			Details:    "the pre-requisites for the workDir are invalid",
		}
	}

	return nil
}

func getEmptyInstanceConfig() *InstanceConfig {
	return &InstanceConfig{
		Paths: &InstancePaths{},
		runtime: &runtimeConfig{
			containerHostInterop: &containerHostInteropConfig{},
			envVars:              map[string]string{},
			excludeDirs:          []string{},
			excludeFiles:         []string{},
			mountDir:             &dagger.Directory{},
		},
		ClientOptions: &ClientOptions{},
		ID:            "",
	}
}

func setDefaultOptionsIfNotSet(options *ClientOptions) *ClientOptions {
	if options.ExcludeOptions == nil {
		options.ExcludeOptions = &ExcludeOptions{
			ExcludedDirs: []string{},
			ExcludeFiles: []string{},
		}
	}

	return options
}

func (i *InstanceImpl) Configure(options *ClientOptions) (*InstanceConfig, error) {
	instanceCfg := getEmptyInstanceConfig()
	setDefaultOptionsIfNotSet(options)

	// 1. Generating the unique identifier.
	instanceCfg.ID = i.td.ID
	i.td.Logger.Info(fmt.Sprintf("The terradagger instance ID is: %s", instanceCfg.ID))

	mountPath := i.td.Config.TerraDagger.Paths.Workspace.SRC
	mountPathAbs := i.td.Config.TerraDagger.Paths.Workspace.SRCAbsolute

	// 2. Resolve the terraDaggerPath, now that the ID is known.
	tdConfigPath, err := ResolveTerraDaggerConfigDirPath(&ResolveTerraDaggerConfigDirPathOptions{
		WorkspaceSRCPath:   mountPath,
		ID:                 instanceCfg.ID,
		TerraDaggerDirName: i.td.Config.TerraDagger.Dirs.TerraDaggerDir,
	})

	i.td.Logger.Info(fmt.Sprintf("The terradagger config path is: %s", tdConfigPath))

	if err != nil {
		return nil, &ClientConfigurationError{
			ErrWrapped: err,
			Details:    fmt.Sprintf("failed to resolve the terradagger configuration path (.terradagger) with path: %s", tdConfigPath),
		}
	}

	instanceCfg.Paths.TerraDagger = tdConfigPath
	instanceCfg.Paths.TerraDaggerAbs, _ = utils.ConvertToAbsolute(tdConfigPath)

	// 3. Resolve the export path
	auxPaths := i.td.fsResolverClient.ResolveAuxPaths(instanceCfg.Paths.TerraDagger)
	instanceCfg.Paths.ExportPath = auxPaths.ExportPath
	instanceCfg.Paths.ExportPathAbs = auxPaths.ExportPathAbs
	// instanceCfg.Paths.ImportPath = auxPaths.ImportPath
	// When importing, it's picking up the files and directories from the directories
	// that the target files/dirs where originally copied to, so it's not necessary to
	// resolve the import path.
	instanceCfg.Paths.ImportPath = auxPaths.ExportPath
	instanceCfg.Paths.ImportPathAbs = auxPaths.ExportPathAbs
	instanceCfg.Paths.CachePath = auxPaths.CachePath
	instanceCfg.Paths.CachePathAbs = auxPaths.CachePathAbs

	// 4. Resolve the mount path. The mount path is inherited from the
	// TerraDagger td, since it's equivalent to the Workspace dir.
	instanceCfg.Paths.MountPath = mountPath
	instanceCfg.Paths.MountPathAbs = mountPathAbs
	instanceCfg.Paths.WorkDirPathDagger = filepath.Join(i.td.Config.Dagger.Paths.MountPathPrefix, options.WorkDirPath)

	// 5. Resolve the workDir path, the relative, and absolute.
	instanceCfg.Paths.WorkDirPath = options.WorkDirPath
	instanceCfg.Paths.WorkDirPathAbs = filepath.Join(mountPathAbs, options.WorkDirPath)
	instanceCfg.Paths.mountPrefix = i.td.Config.Dagger.Paths.MountPathPrefix

	// 6. Resolve the runtime configuration.
	mountDirInDaggerFormat, err := i.td.daggerConfigClient.
		GetMountPathAsDaggerDir(mountPath, i.td.DaggerBackend)

	if err != nil {
		return nil, fmt.Errorf("failed to Configure the terradagger instance, the dagger dir is invalid: %w", err)
	}

	instanceCfg.runtime = &runtimeConfig{
		image:                options.ContainerOptions.Image,
		version:              options.ContainerOptions.Version,
		mountDir:             mountDirInDaggerFormat,
		commands:             options.TerraDaggerCMDs,
		containerHostInterop: &containerHostInteropConfig{},
	}

	// 7. Add the excluded files, and dirs mixed with the default, and the passed options (if any).
	instanceCfg.runtime.excludeDirs = utils.MixSlices(i.td.Config.Dagger.Excluded.ExcludedDirs, options.ExcludeOptions.ExcludedDirs)
	instanceCfg.runtime.excludeFiles = utils.MixSlices(i.td.Config.Dagger.Excluded.ExcludedFiles,
		options.ExcludeOptions.ExcludeFiles)

	// 8. Create directories required for this instance.
	_, err = i.td.dirManagerClient.CreateTerraDaggerDir(&CreateTerraDaggerDirOptions{
		TerraDaggerPathResolved: instanceCfg.Paths.TerraDagger,
		SkipCreationIfExist:     true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to Configure the terradagger instance, the creation of the terra dagger dir failed: %w", err)
	}

	if err := i.td.dirManagerClient.CreateAuxDirsInTDDir([]CreateAuxDirsInTDDirOptions{
		{
			TerraDaggerPathResolved: instanceCfg.Paths.TerraDagger,
			SkipCreationIfExist:     true,
			AuxDirName:              i.td.Config.TerraDagger.Dirs.TerraDaggerExportDir,
		},
		{
			TerraDaggerPathResolved: instanceCfg.Paths.TerraDagger,
			SkipCreationIfExist:     true,
			AuxDirName:              i.td.Config.TerraDagger.Dirs.TerraDaggerImportDir,
		},
		{
			TerraDaggerPathResolved: instanceCfg.Paths.TerraDagger,
			SkipCreationIfExist:     true,
			AuxDirName:              i.td.Config.TerraDagger.Dirs.TerraDaggerCacheDir,
		},
	}); err != nil {
		return nil, fmt.Errorf("failed to Configure the terradagger instance, the creation of the aux dirs in the terra dagger dir failed: %w", err)
	}

	// 9 Configure the host/container interop.
	clientConfigurator := NewClientConfigurator(i)
	exportCfg, err := clientConfigurator.ConfigureExportFromContainer(&ConfigureExportFromContainerOptions{
		ParamOptions:      options.ExportFromContainer,
		WorkDirPathDagger: instanceCfg.Paths.WorkDirPathDagger,
		ExportPathInHost:  instanceCfg.Paths.ExportPath,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to Configure the terradagger instance, the configuration of the export from container failed: %w", err)
	}

	if exportCfg != nil {
		instanceCfg.runtime.containerHostInterop.copyFilesToHost = exportCfg.PathsCopyFilesToHost
		instanceCfg.runtime.containerHostInterop.copyDirsToHost = exportCfg.PathsCopyDirsToHost
	}

	// 10. Import from container to host (files and dirs)
	importCfg, err := clientConfigurator.ConfigureImportToContainer(&ConfigureImportToContainerOptions{
		ParamOptions:           options.ImportToContainer,
		WorkDirPathInContainer: instanceCfg.Paths.WorkDirPathDagger,
		WorkDirPathInHost:      instanceCfg.Paths.WorkDirPathAbs,
		ClientImportPathInHost: instanceCfg.Paths.ImportPath,
		SourceImportPathInHost: instanceCfg.Paths.ExportPathAbs,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to Configure the terradagger instance, the configuration of the import to container failed: %w", err)
	}

	instanceCfg.runtime.containerHostInterop.copyToContainerConfig = importCfg

	// if options.ImportToContainer != nil {
	// 	if len(options.ImportToContainer.FileNames) > 0 {
	// 		for _, file := range options.ImportToContainer.FileNames {
	// 			var fileInHostPath string
	// 			if options.ImportToContainer.LookupFromWorkDir {
	// 				fileInHostPath = filepath.Join(instanceCfg.Paths.WorkDirPathAbs, file)
	// 			} else {
	// 				fileInHostPath = filepath.Join(instanceCfg.Paths.ImportPathAbs, file)
	// 			}
	//
	// 			instanceCfg.runtime.containerHostInterop.copyFilesToContainer = append(instanceCfg.runtime.containerHostInterop.copyFilesToContainer, fileInHostPath)
	// 		}
	// 	}
	//
	// 	if len(options.ImportToContainer.DirNames) > 0 {
	// 		for _, dir := range options.ImportToContainer.DirNames {
	// 			var dirInHostPath string
	// 			if options.ImportToContainer.LookupFromWorkDir {
	// 				dirInHostPath = filepath.Join(instanceCfg.Paths.WorkDirPathAbs, dir)
	// 			} else {
	// 				dirInHostPath = filepath.Join(instanceCfg.Paths.ImportPathAbs, dir)
	// 			}
	//
	// 			// Check if the dirInHostPath resolved exist, and it's valid.
	// 			if err := config.IsAValidTerraDaggerDirAbsolute(dirInHostPath); err != nil {
	// 				return nil, fmt.Errorf("failed to Configure the terradagger instance, the dir %s is invalid: %w", dir, err)
	// 			}
	//
	// 			instanceCfg.runtime.containerHostInterop.copyDirsToContainer = append(instanceCfg.runtime.containerHostInterop.copyDirsToContainer, dirInHostPath)
	// 		}
	// 	}
	// }
	//

	// 11. Env vars
	clientEnvVars, err := clientConfigurator.ConfigureEnvVars(&ConfigureEnvVarsOptions{
		ParamOptions: options.EnvVarOptions,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to Configure the terradagger instance, the configuration of the env vars failed: %w", err)
	}

	instanceCfg.runtime.envVars = clientEnvVars.EnvVars

	// 12. Pass client options
	instanceCfg.ClientOptions = options

	return instanceCfg, nil
}

func (i *InstanceImpl) PrepareInstance(cfg *InstanceConfig) (*ClientInstance, error) {
	if cfg == nil {
		return nil, fmt.Errorf("failed to PrepareInstance, the instance config is nil")
	}

	instanceID := cfg.ID
	instance := &ClientInstance{
		ID:               instanceID,
		Config:           cfg,
		td:               i.td,
		runtimeContainer: &Container{},
	}

	i.td.Logger.Info(fmt.Sprintf("Preparing the terradagger instance with ID: %s", instanceID))

	// 1. Creating the TerraDagger runtime container.
	tdContainer, err := i.td.
		CreateTerraDaggerContainer(&CreateTerraDaggerContainerOptions{
			Image:   cfg.runtime.image,
			Version: cfg.runtime.version,
		})

	if err != nil {
		return nil, fmt.Errorf("failed to PrepareInstance, the creation of the terradagger runtime failed: %w", err)
	}

	instance.runtimeContainer = tdContainer

	// 2. Configure the TerraDagger runtime container
	runtimeContainerConfigured, err := i.td.
		ConfigureTerraDaggerContainer(&ConfigureTerraDaggerContainerOptions{
			Container:    tdContainer,
			EnvVars:      cfg.runtime.envVars,
			MountDir:     cfg.runtime.mountDir,
			WorkDirPath:  cfg.Paths.WorkDirPathDagger,
			ExcludedDirs: cfg.runtime.excludeDirs,
			ExcludeFiles: cfg.runtime.excludeFiles,
			Commands:     cfg.runtime.commands,
		})

	if err != nil {
		return nil, fmt.Errorf("failed to PrepareInstance, the configuration of the terradagger runtime failed: %w", err)
	}

	instance.runtimeContainer = runtimeContainerConfigured
	return instance, nil
}
