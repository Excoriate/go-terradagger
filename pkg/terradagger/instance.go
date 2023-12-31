package terradagger

import (
	"fmt"
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/env"

	"github.com/Excoriate/go-terradagger/pkg/config"

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
	// ImportExportOptions are the options for the files and directories that should be imported/exported from the terradagger.
	ImportExportOptions *ImportExportOptions
	// EnvVarOptions are the options for the env vars that should be passed to the terradagger.
	EnvVarOptions *EnvVarOptions
	// WorkDirPreRequisites are the options for the files and directories that should be required from the terradagger.
	WorkDirPreRequisites *Requisites
	// MountDirPreRequisites are the options for the files and directories that should be required from the terradagger.
	MountDirPreRequisites *Requisites
}

type Requisites struct {
	RequiredFiles          []string
	RequiredDirs           []string
	RequiredFileExtensions []string
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
	containerHostInterop *ContainerHostInteropConfig
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
	ExportPath string
	// ImportPath is the path of the import directory in the host,
	// formatted as: .terradagger/ID/import
	ImportPath string
	// CachePath is the path of the CachePath directory in the host,
	// formatted as: .terradagger/ID/.terradagger-CachePath
	CachePath string
	// MountPath is the path of the mount directory in the runtime,
	// formatted as: /mnt/MountPath
	MountPath string
	// MountPathAbsolute is the absolute path of the mount directory in the runtime,
	MountPathAbsolute string
	// WorkDirPath is the path of the work directory in the runtime,
	// formatted as: /mnt/MountPath/WorkDirPath
	WorkDirPath string
	// WorkDirPathAbsolute is the absolute path of the work directory in the runtime,
	// WorkDirPathDagger is the path of the work directory in the runtime, which includes
	// the dagger mount prefix, formatted as: /mnt/MountPath/WorkDirPath
	WorkDirPathDagger   string
	WorkDirPathAbsolute string
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

type ContainerHostInteropConfig struct {
	CopyFromHostToContainerBasePath string
	CopyFromContainerToHostBasePath string
	copyFromContainerToHostFiles    []string
	copyFromHostToContainerFiles    []string
	copyFromContainerToHostDirs     []string
	copyFromHostToContainerDirs     []string
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

func (i *InstanceImpl) Validate(options *ClientOptions) error {
	if options == nil {
		return fmt.Errorf("failed to Validate the terradagger instance, the options are nil")
	}

	cv := newContainerValidator(i.td.Logger)

	if err := cv.validate(&CreateNewContainerOptions{
		Image:   options.ContainerOptions.Image,
		Version: options.ContainerOptions.Version,
	}); err != nil {
		return fmt.Errorf("failed to Validate the terradagger instance, the runtime options are nil")
	}

	if options.TerraDaggerCMDs == nil {
		return fmt.Errorf("failed to Validate the terradagger instance, the terradagger commands are nil")
	}

	if len(options.TerraDaggerCMDs) == 0 {
		return fmt.Errorf("failed to Validate the terradagger instance, the terradagger commands are empty")
	}

	if options.EnvVarOptions != nil {
		if len(options.EnvVarOptions.EnvVars) == 0 {
			i.td.Logger.Info("The env vars are empty, so this 'terraDagger' instance will not have custom env vars")
		}

		if len(options.EnvVarOptions.CopyEnvVarsFromHostByKeys) > 0 && options.EnvVarOptions.MirrorEnvVarsFromHost {
			return fmt.Errorf("failed to Validate the terradagger instance, the env vars options are invalid, " +
				"the 'CopyEnvVarsFromHostByKeys' and 'MirrorEnvVarsFromHost' options cannot be used together")
		}

		if len(options.EnvVarOptions.CopyEnvVarsFromHostByKeys) > 0 {
			for _, key := range options.EnvVarOptions.CopyEnvVarsFromHostByKeys {
				_, err := env.GetEnvVarByKey(key, true)
				if err != nil {
					return fmt.Errorf("failed to Validate the terradagger instance, the env vars options are invalid, "+
						"the env var with key %s does not exist", key)
				}
			}
		}
	}

	mountPath := i.td.Config.TerraDagger.Paths.Workspace.SRC

	if err := i.td.fsResolverClient.IsMountAndWorkDirPathValid(mountPath, options.WorkDirPath); err != nil {
		return fmt.Errorf("failed to Validate the terradagger instance, the mount and work dir path is invalid: %w", err)
	}

	if options.ExcludeOptions == nil {
		i.td.Logger.Warn("The exclude options are empty, " +
			"so this 'terraDagger' instance will exclude only the default directories and files")
	} else {
		if err := config.AreFilesToExcludeValid(mountPath, options.ExcludeOptions.ExcludeFiles); err != nil {
			return fmt.Errorf("failed to Validate the terradagger instance, the files to exclude are invalid: %w", err)
		}

		if err := config.AreDirsToExcludeValid(mountPath, options.ExcludeOptions.ExcludedDirs); err != nil {
			return fmt.Errorf("failed to Validate the terradagger instance, the dirs to exclude are invalid: %w", err)
		}
	}

	mountPathAbs := i.td.Config.TerraDagger.Paths.Workspace.SRCAbsolute

	if options.ImportExportOptions != nil {
		filesToCopyToContainer := options.ImportExportOptions.ExportFileNames
		dirsToCopyToContainer := options.ImportExportOptions.ExportDirNames

		if len(filesToCopyToContainer) > 0 {
			for _, file := range filesToCopyToContainer {
				fileHostPath := filepath.Join(mountPathAbs, file)
				if err := i.td.fsResolverClient.IsFileValidInHost(fileHostPath); err != nil {
					return fmt.Errorf("failed to Configure the terradagger instance, the file %s is invalid: %w", file, err)
				}
			}
		}

		if len(dirsToCopyToContainer) > 0 {
			for _, dir := range dirsToCopyToContainer {
				dirHostPath := filepath.Join(mountPathAbs, dir)
				if err := config.IsAValidTerraDaggerDirAbsolute(dirHostPath); err != nil {
					return fmt.Errorf("failed to Configure the terradagger instance, the dir %s is invalid: %w", dir, err)
				}
			}
		}
	}

	if options.WorkDirPreRequisites != nil {
		if len(options.WorkDirPreRequisites.RequiredDirs) > 0 {
		}
	}

	return nil
}

func getEmptyInstanceConfig() *InstanceConfig {
	return &InstanceConfig{
		Paths: &InstancePaths{},
		runtime: &runtimeConfig{
			containerHostInterop: &ContainerHostInteropConfig{},
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
	instanceCfg.ID = utils.GetUUID()
	i.td.Logger.Info(fmt.Sprintf("The terradagger instance ID is: %s", instanceCfg.ID))

	mountPath := i.td.Config.TerraDagger.Paths.Workspace.SRC
	mountPathAbs := i.td.Config.TerraDagger.Paths.Workspace.SRCAbsolute

	// 2. Resolve the terraDaggerPath, now that the ID is known.
	terraDaggerDirPath, err := i.td.fsResolverClient.
		ResolveTerraDaggerDirPath(&ResolveTerraDaggerDirPathOptions{
			WorkspaceSRCPath: mountPath,
			ID:               instanceCfg.ID,
		})

	if err != nil {
		return nil, fmt.Errorf("failed to Configure the terradagger instance, the terra dagger dir path is invalid: %w", err)
	}

	instanceCfg.Paths.TerraDagger = terraDaggerDirPath

	// 3. Resolve the export path
	auxPaths := i.td.fsResolverClient.ResolveAuxPaths(instanceCfg.Paths.TerraDagger)
	instanceCfg.Paths.ExportPath = auxPaths.ExportPath
	instanceCfg.Paths.ImportPath = auxPaths.ImportPath
	instanceCfg.Paths.CachePath = auxPaths.CachePath

	// 4. Resolve the mount path. The mount path is inherited from the
	// TerraDagger td, since it's equivalent to the Workspace dir.
	instanceCfg.Paths.MountPath = mountPath
	instanceCfg.Paths.MountPathAbsolute = mountPathAbs
	instanceCfg.Paths.WorkDirPathDagger = filepath.Join(i.td.Config.Dagger.Paths.MountPathPrefix, options.WorkDirPath)

	// 5. Resolve the workDir path, the relative, and absolute.
	instanceCfg.Paths.WorkDirPath = options.WorkDirPath
	instanceCfg.Paths.WorkDirPathAbsolute = filepath.Join(mountPathAbs, options.WorkDirPath)
	instanceCfg.Paths.mountPrefix = i.td.Config.Dagger.Paths.MountPathPrefix

	// 6. Resolve the runtime configuration.
	mountDirInDaggerFormat, err := i.td.daggerConfigClient.
		GetMountPathAsDaggerDir(mountPath, i.td.DaggerBackend)

	if err != nil {
		return nil, fmt.Errorf("failed to Configure the terradagger instance, the dagger dir is invalid: %w", err)
	}

	instanceCfg.runtime = &runtimeConfig{
		image:    options.ContainerOptions.Image,
		version:  options.ContainerOptions.Version,
		mountDir: mountDirInDaggerFormat,
		commands: options.TerraDaggerCMDs,
	}

	// 7. Add the excluded files, and dirs mixed with the default, and the passed options (if any).
	instanceCfg.runtime.excludeDirs = utils.MixSlices(i.td.Config.Dagger.Excluded.ExcludedDirs, options.ExcludeOptions.ExcludedDirs)
	instanceCfg.runtime.excludeFiles = utils.MixSlices(i.td.Config.Dagger.Excluded.ExcludedFiles,
		options.ExcludeOptions.ExcludeFiles)

	// 8 Configure the base paths for the runtime/host interoperation.
	interopConfig := &ContainerHostInteropConfig{
		CopyFromHostToContainerBasePath: instanceCfg.Paths.ExportPath,
		CopyFromContainerToHostBasePath: instanceCfg.Paths.WorkDirPath,
	}

	instanceCfg.runtime.containerHostInterop = interopConfig

	// 9. Env vars
	if options.EnvVarOptions != nil {
		customEnvVars := options.EnvVarOptions.EnvVars
		var varsScannedFromHostByKey map[string]string

		if len(options.EnvVarOptions.CopyEnvVarsFromHostByKeys) > 0 {
			varsScannedFromHostByKey = map[string]string{}
			for _, key := range options.EnvVarOptions.CopyEnvVarsFromHostByKeys {
				envVar, envKeyErr := env.GetEnvVarByKey(key, true)
				if envKeyErr != nil {
					return nil, fmt.Errorf("failed to Configure the terradagger instance, the env vars options are invalid, "+
						"the env var with key %s does not exist", key)
				}

				varsScannedFromHostByKey[key] = envVar
			}
		}

		var hostEnvVars map[string]string
		if options.EnvVarOptions.MirrorEnvVarsFromHost {
			hostEnvVars = env.GetAllFromHost()
		}

		instanceCfg.runtime.envVars = utils.MergeMaps(customEnvVars, varsScannedFromHostByKey, hostEnvVars)
	}

	// 10. Create directories required for this instance.
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

	// 11. If files, and directories were passed, validate their path.
	if options.ImportExportOptions != nil {
		filesToCopyToContainer := options.ImportExportOptions.ExportFileNames
		dirsToCopyToContainer := options.ImportExportOptions.ExportDirNames

		if len(filesToCopyToContainer) > 0 {
			for _, file := range filesToCopyToContainer {
				fileHostPath := filepath.Join(mountPathAbs, file)
				if err := i.td.fsResolverClient.IsFileValidInHost(fileHostPath); err != nil {
					return nil, fmt.Errorf("failed to Configure the terradagger instance, the file %s is invalid: %w", file, err)
				}

				instanceCfg.runtime.containerHostInterop.copyFromHostToContainerFiles = append(instanceCfg.runtime.containerHostInterop.copyFromHostToContainerFiles, fileHostPath)
			}
		}

		if len(dirsToCopyToContainer) > 0 {
			for _, dir := range dirsToCopyToContainer {
				dirHostPath := filepath.Join(mountPathAbs, dir)
				if err := config.IsAValidTerraDaggerDirAbsolute(dirHostPath); err != nil {
					return nil, fmt.Errorf("failed to Configure the terradagger instance, the dir %s is invalid: %w", dir, err)
				}

				instanceCfg.runtime.containerHostInterop.copyFromHostToContainerDirs = append(instanceCfg.runtime.containerHostInterop.copyFromHostToContainerDirs, dirHostPath)
			}
		}
	}

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
