package terradagger

import (
	"fmt"

	"github.com/Excoriate/go-terradagger/pkg/utils"

	"github.com/Excoriate/go-terradagger/pkg/env"
)

type ClientConfigurator interface {
	ConfigureExportFromContainer(options *ConfigureExportFromContainerOptions) (*DataTransferToHost, error)
	ConfigureImportToContainer(options *ConfigureImportToContainerOptions) (*DataTransferToContainer, error)
	ConfigureEnvVars(options *ConfigureEnvVarsOptions) (*ConfigureEnvVarsResult, error)
}

type ClientConfiguratorImpl struct {
	clientInstance *InstanceImpl
}

type ClientConfiguratorError struct {
	ErrWrapped error
	Details    string
}

const ClientConfiguratorErrPrefix = "the client configurator failed to configure the terradagger client"

func (e *ClientConfiguratorError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", ClientConfiguratorErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", ClientConfiguratorErrPrefix, e.Details, e.ErrWrapped.Error())
}

func NewClientConfigurator(clientInstance *InstanceImpl) ClientConfigurator {
	return &ClientConfiguratorImpl{
		clientInstance: clientInstance,
	}
}

type ConfigureExportFromContainerOptions struct {
	ParamOptions      *ExportFromContainerOptions
	WorkDirPathDagger string
	ExportPathInHost  string
}

type ConfigureExportFromContainerResult struct {
	PathsCopyDirsToHost  []string
	PathsCopyFilesToHost []string
}

func (c *ClientConfiguratorImpl) ConfigureExportFromContainer(
	options *ConfigureExportFromContainerOptions) (*DataTransferToHost, error) {
	if options == nil || options.ParamOptions == nil {
		c.clientInstance.td.Logger.Info("no files or dirs to export from the container")
		return nil, nil
	}

	if options.WorkDirPathDagger == "" {
		return nil, &ClientConfiguratorError{
			Details: "the work dir path in dagger is empty",
		}
	}

	if options.ExportPathInHost == "" {
		return nil, &ClientConfiguratorError{
			Details: "the export path in host is empty",
		}
	}

	paramOptions := options.ParamOptions

	if len(paramOptions.FileNames) == 0 && len(paramOptions.DirNames) == 0 {
		return nil, &ClientConfiguratorError{
			Details: "the export options are empty",
		}
	}

	result := &DataTransferToHost{}
	result.WorkDirPath = options.WorkDirPathDagger

	for _, fileName := range paramOptions.FileNames {
		exportFilePath := fmt.Sprintf("%s/%s", options.WorkDirPathDagger, fileName)
		exportPathInHost := fmt.Sprintf("%s/%s", options.ExportPathInHost, fileName)
		exportPathInHostAbs, _ := utils.ConvertToAbsolute(exportPathInHost)

		result.Files = append(result.Files, TransferToHost{
			SourcePathInContainer:    exportFilePath,
			DestinationPathInHostAbs: exportPathInHostAbs,
		})

		if paramOptions.OverrideIfExistInHost {
			if err := c.clientInstance.td.OverrideTerraDaggerFile(&OverrideTerraDaggerFileOptions{
				FileName:             fileName,
				TerraDaggerID:        c.clientInstance.td.ID,
				TerraDaggerConfigDir: "export",
			}); err != nil {
				return nil, &ClientConfiguratorError{
					ErrWrapped: err,
					Details:    fmt.Sprintf("failed to override the file %s in the host", fileName),
				}
			}
		}

		c.clientInstance.td.Logger.Info(fmt.Sprintf("The file %s will be exported from the container to %s", fileName, options.ExportPathInHost))
	}

	for _, dirName := range paramOptions.DirNames {
		exportDirPath := fmt.Sprintf("%s/%s", options.WorkDirPathDagger, dirName)
		exportPathInHost := fmt.Sprintf("%s/%s", options.ExportPathInHost, dirName)
		exportPathInHostAbs, _ := utils.ConvertToAbsolute(exportPathInHost)

		result.Dirs = append(result.Dirs, TransferToHost{
			SourcePathInContainer:    exportDirPath,
			DestinationPathInHostAbs: exportPathInHostAbs,
		})

		if paramOptions.OverrideIfExistInHost {
			if err := c.clientInstance.td.OverrideTerraDaggerDir(&OverrideTerraDaggerDirOptions{
				DirName:              dirName,
				TerraDaggerID:        c.clientInstance.td.ID,
				TerraDaggerConfigDir: "export",
			}); err != nil {
				return nil, &ClientConfiguratorError{
					ErrWrapped: err,
					Details:    fmt.Sprintf("failed to override the dir %s in the host", dirName),
				}
			}
		}

		c.clientInstance.td.Logger.Info(fmt.Sprintf("The dir %s will be exported from the container to %s", dirName, options.ExportPathInHost))
	}

	return result, nil
}

type ConfigureEnvVarsOptions struct {
	ParamOptions *EnvVarOptions
}

type ConfigureEnvVarsResult struct {
	EnvVars map[string]string
}

func (c *ClientConfiguratorImpl) ConfigureEnvVars(options *ConfigureEnvVarsOptions) (*ConfigureEnvVarsResult, error) {
	if options == nil || options.ParamOptions == nil {
		c.clientInstance.td.Logger.Info("no env vars to configure")
		return &ConfigureEnvVarsResult{}, nil
	}

	paramOptions := options.ParamOptions
	envVarsResult := &ConfigureEnvVarsResult{}

	// Copy environment variables from host by keys
	if len(paramOptions.CopyEnvVarsFromHostByKeys) > 0 {
		envVarsResult.EnvVars = make(map[string]string)
		for _, key := range paramOptions.CopyEnvVarsFromHostByKeys {
			envVar, envKeyErr := env.GetEnvVarByKey(key, true)
			if envKeyErr != nil {
				return nil, fmt.Errorf("failed to configure the terradagger instance, the env vars options are invalid, "+
					"the env var with key %s does not exist: %w", key, envKeyErr)
			}
			envVarsResult.EnvVars[key] = envVar
			c.clientInstance.td.Logger.Info("configured env var from host", "key", key, "value", envVar)
		}
	}

	// Mirror environment variables from host
	if paramOptions.MirrorEnvVarsFromHost {
		hostEnvVars := env.GetAllFromHost()
		if envVarsResult.EnvVars == nil {
			envVarsResult.EnvVars = hostEnvVars
		} else {
			for key, value := range hostEnvVars {
				envVarsResult.EnvVars[key] = value
			}
		}
		c.clientInstance.td.Logger.Info("mirrored env vars from host")
	}

	// Merge custom environment variables
	for key, value := range paramOptions.EnvVars {
		if envVarsResult.EnvVars == nil {
			envVarsResult.EnvVars = make(map[string]string)
		}
		envVarsResult.EnvVars[key] = value
		c.clientInstance.td.Logger.Info("configured custom env var", "key", key, "value", value)
	}

	return envVarsResult, nil
}

type ConfigureImportToContainerOptions struct {
	ParamOptions           *ImportToContainerOptions
	WorkDirPathInContainer string
	ClientImportPathInHost string
	WorkDirPathInHost      string
	SourceImportPathInHost string
}

func (c *ClientConfiguratorImpl) ConfigureImportToContainer(
	options *ConfigureImportToContainerOptions) (*DataTransferToContainer, error) {
	if options == nil || options.ParamOptions == nil {
		c.clientInstance.td.Logger.Info("no files or dirs to import to the container")
		return nil, nil
	}

	if options.WorkDirPathInContainer == "" {
		return nil, &ClientConfiguratorError{
			Details: "the work dir path in container is empty",
		}
	}

	if options.SourceImportPathInHost == "" {
		return nil, &ClientConfiguratorError{
			Details: "the source import path in host is empty",
		}
	}

	paramOptions := options.ParamOptions

	if len(paramOptions.FileNames) == 0 && len(paramOptions.DirNames) == 0 {
		return nil, &ClientConfiguratorError{
			Details: "the import options are empty",
		}
	}

	if paramOptions.LookupFromWorkDir && options.WorkDirPathInHost == "" {
		return nil, &ClientConfiguratorError{
			Details: "if the lookup from work dir is enabled, the work dir path in host must be specified",
		}
	}

	if options.ClientImportPathInHost == "" && !paramOptions.LookupFromWorkDir {
		return nil, &ClientConfiguratorError{
			Details: "the client import path (resolved) is required if the lookup from work dir is disabled or not set",
		}
	}

	result := &DataTransferToContainer{}
	result.WorkDirPath = options.WorkDirPathInContainer

	for _, fileName := range paramOptions.FileNames {
		var sourcePathInHost string
		if paramOptions.LookupFromWorkDir {
			sourcePathInHost = fmt.Sprintf("%s/%s", options.WorkDirPathInHost, fileName)
		} else {
			sourcePathInHost = fmt.Sprintf("%s/%s", options.ClientImportPathInHost, fileName)
		}

		if err := utils.FileExistE(sourcePathInHost); err != nil {
			return nil, &ClientConfiguratorError{
				ErrWrapped: err,
				Details:    fmt.Sprintf("the file %s to import does not exist", sourcePathInHost),
			}
		}

		if err := utils.IsAFileE(sourcePathInHost); err != nil {
			return nil, &ClientConfiguratorError{
				ErrWrapped: err,
				Details:    fmt.Sprintf("the file %s to import is not a file", sourcePathInHost),
			}
		}

		sourcePathInHostAbs, _ := utils.ConvertToAbsolute(sourcePathInHost)

		result.Files = append(result.Files, TransferToContainer{
			SourcePathInHostAbs:        sourcePathInHostAbs,
			DestinationPathInContainer: fmt.Sprintf("%s/%s", options.WorkDirPathInContainer, fileName),
		})

	}

	for _, dirName := range paramOptions.DirNames {
		var sourcePathInHost string
		if paramOptions.LookupFromWorkDir {
			sourcePathInHost = fmt.Sprintf("%s/%s", options.WorkDirPathInHost, dirName)
		} else {
			sourcePathInHost = fmt.Sprintf("%s/%s", options.ClientImportPathInHost, dirName)
		}

		dirUtils := utils.DirUtils{}
		if err := dirUtils.IsValidDirE(sourcePathInHost); err != nil {
			return nil, &ClientConfiguratorError{
				ErrWrapped: err,
				Details:    fmt.Sprintf("the dir %s to import is not a valid dir", sourcePathInHost),
			}
		}

		sourcePathAbs, _ := utils.ConvertToAbsolute(sourcePathInHost)

		result.Dirs = append(result.Dirs, TransferToContainer{
			SourcePathInHostAbs:        sourcePathAbs,
			DestinationPathInContainer: fmt.Sprintf("%s/%s", options.WorkDirPathInContainer, dirName),
		})

	}

	return result, nil
}
