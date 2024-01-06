package terradagger

import (
	"fmt"

	"github.com/Excoriate/go-terradagger/pkg/utils"

	"github.com/Excoriate/go-terradagger/pkg/env"
)

type clientInstanceValidator interface {
	IsEnvVarOptionsValid(options *EnvVarOptions) error
	IsContainerOptionsValid(options *InstanceContainerOptions) error
	IsExcludedOptionsValid(mountPath string, options *ExcludeOptions) error
	IsImportToContainerOptionsValid(mountPath, workDirPath string, options *ImportToContainerOptions) error
	IsPreRequisitesValid(mountPath, workDirPath string, requisites *Requisites) error
}

type clientInstanceValidatorImpl struct {
	clientInstance *InstanceImpl
}

func newClientInstanceValidator(clientInstance *InstanceImpl) clientInstanceValidator {
	return &clientInstanceValidatorImpl{
		clientInstance: clientInstance,
	}
}

type EnvVarOptionsError struct {
	ErrWrapped error
	Details    string
}

const envVarOptionsErrPrefix = "the environment variable options passed to the terradagger client" +
	" instance are invalid"

func (e *EnvVarOptionsError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", envVarOptionsErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", envVarOptionsErrPrefix, e.Details, e.ErrWrapped.Error())
}

func (v *clientInstanceValidatorImpl) IsEnvVarOptionsValid(options *EnvVarOptions) error {
	if options == nil {
		return nil
	}

	if len(options.EnvVars) == 0 {
		return &EnvVarOptionsError{
			Details: "the environment variables map is empty",
		}
	}

	if len(options.CopyEnvVarsFromHostByKeys) > 0 && options.MirrorEnvVarsFromHost {
		return &EnvVarOptionsError{
			Details: "the environment variables cannot be copied from the host and mirrored at the same time",
		}
	}

	if len(options.CopyEnvVarsFromHostByKeys) > 0 {
		for _, key := range options.CopyEnvVarsFromHostByKeys {
			if _, err := env.GetEnvVarByKey(key, true); err != nil {
				return &EnvVarOptionsError{
					ErrWrapped: err,
					Details:    fmt.Sprintf("the environment variable %s cannot be copied from the host", key),
				}
			}
		}
	}

	return nil
}

type ContainerOptionsError struct {
	ErrWrapped error
	Details    string
}

const containerOptionsErrPrefix = "the container options passed to the terradagger client" +
	" instance are invalid"

func (e *ContainerOptionsError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", containerOptionsErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", containerOptionsErrPrefix, e.Details, e.ErrWrapped.Error())
}

func (v *clientInstanceValidatorImpl) IsContainerOptionsValid(options *InstanceContainerOptions) error {
	containerValidator := NewContainerValidator(v.clientInstance.td.Logger)
	if err := containerValidator.ValidateContainerImage(&CreateNewContainerOptions{
		Image:   options.Image,
		Version: options.Version,
	}); err != nil {
		return &ContainerOptionsError{
			ErrWrapped: err,
			Details:    "the container options are invalid",
		}
	}

	return nil
}

type ExcludeOptionsError struct {
	ErrWrapped error
	Details    string
}

const excludeOptionsErrPrefix = "the exclude options passed to the terradagger client" +
	" instance are invalid"

func (e *ExcludeOptionsError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", excludeOptionsErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", excludeOptionsErrPrefix, e.Details, e.ErrWrapped.Error())
}

func (v *clientInstanceValidatorImpl) IsExcludedOptionsValid(mountPath string, options *ExcludeOptions) error {
	if options == nil {
		return nil
	}

	if mountPath == "" {
		return &ExcludeOptionsError{
			Details: "the mount dir path is empty",
		}
	}

	exFilesClient := NewExcludedFiles(v.clientInstance.td)
	exDirsClient := NewExcludedDirs(v.clientInstance.td)

	for _, fileToExclude := range options.ExcludeFiles {
		if err := exFilesClient.IsExcludeFileValid(mountPath, fileToExclude); err != nil {
			return &ExcludeOptionsError{
				ErrWrapped: err,
				Details:    "the exclude files options are invalid",
			}
		}
	}

	for _, dirToExclude := range options.ExcludedDirs {
		if err := exDirsClient.IsExcludeDirValid(mountPath, dirToExclude); err != nil {
			return &ExcludeOptionsError{
				ErrWrapped: err,
				Details:    "the exclude dirs options are invalid",
			}
		}
	}

	return nil
}

type ImportToContainerOptionsError struct {
	ErrWrapped error
	Details    string
}

const importToContainerOptionsErrPrefix = "the import to container options passed to the terradagger client" +
	" instance are invalid"

func (e *ImportToContainerOptionsError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", importToContainerOptionsErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", importToContainerOptionsErrPrefix, e.Details, e.ErrWrapped.Error())
}

func (v *clientInstanceValidatorImpl) IsImportToContainerOptionsValid(mountPath, workDirPath string, options *ImportToContainerOptions) error {
	if options == nil {
		return nil
	}

	if mountPath == "" {
		return &ImportToContainerOptionsError{
			Details: "the mount dir path is empty",
		}
	}

	files := options.FileNames
	dirs := options.DirNames

	for _, file := range files {
		if file == "" {
			return &ImportToContainerOptionsError{
				Details: "the file name is empty",
			}
		}

		filePath, _ := GetTerraDaggerPath(&GetTerraDaggerPathOptions{
			WorkDirPath:   workDirPath,
			MountDirPath:  mountPath,
			FileOrDirName: file,
		})

		if err := utils.FileExistE(filePath); err != nil {
			return &ImportToContainerOptionsError{
				ErrWrapped: err,
				Details:    fmt.Sprintf("the file %s to import to the container does not exist", filePath),
			}
		}
	}

	du := utils.NewDirUtils()
	for _, dir := range dirs {
		if dir == "" {
			return &ImportToContainerOptionsError{
				Details: "the dir name is empty",
			}
		}

		dirPath, _ := GetTerraDaggerPath(&GetTerraDaggerPathOptions{
			WorkDirPath:   workDirPath,
			MountDirPath:  mountPath,
			FileOrDirName: dir,
		})

		if err := du.DirExistE(dirPath); err != nil {
			return &ImportToContainerOptionsError{
				ErrWrapped: err,
				Details:    fmt.Sprintf("the dir %s to import to the container does not exist", dirPath),
			}
		}
	}

	return nil
}

type ExportFromContainerOptionsError struct {
	ErrWrapped error
	Details    string
}

const exportFromContainerOptionsErrPrefix = "the export from container options passed to the terradagger client" +
	" instance are invalid"

func (e *ExportFromContainerOptionsError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", exportFromContainerOptionsErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", exportFromContainerOptionsErrPrefix, e.Details, e.ErrWrapped.Error())
}

type RequisitesError struct {
	ErrWrapped error
	Details    string
}

const requisitesErrPrefix = "the requisites passed to the terradagger client" +
	" instance are invalid"

func (e *RequisitesError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", requisitesErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", requisitesErrPrefix, e.Details, e.ErrWrapped.Error())
}

func (v *clientInstanceValidatorImpl) IsPreRequisitesValid(mountPath, workDirPath string, requisites *Requisites) error {
	if requisites == nil {
		return nil
	}

	if mountPath == "" {
		return &RequisitesError{
			Details: "the mount dir path is empty",
		}
	}

	if workDirPath == "" {
		return &RequisitesError{
			Details: "the work dir path is empty",
		}
	}

	workDirPath = GetTerraDaggerWorkDirPath(mountPath, workDirPath)

	for _, ext := range requisites.RequiredFileExtensions {
		files, err := utils.GetFilesByExtension(workDirPath, ext)
		if err != nil {
			return &RequisitesError{
				ErrWrapped: err,
				Details:    fmt.Sprintf("the files with extension %s could not be found", ext),
			}
		}

		if len(files) == 0 {
			return &RequisitesError{
				Details: fmt.Sprintf("the files with extension %s could not be found", ext),
			}
		}

		for _, file := range files {
			v.clientInstance.td.Logger.Info(fmt.Sprintf("the file %s with extension %s was found", file, ext))
		}
	}

	for _, file := range requisites.RequiredFiles {
		filePath, _ := GetTerraDaggerPath(&GetTerraDaggerPathOptions{
			WorkDirPath:   workDirPath,
			MountDirPath:  mountPath,
			FileOrDirName: file,
		})

		if err := utils.FileExistE(filePath); err != nil {
			return &RequisitesError{
				ErrWrapped: err,
				Details:    fmt.Sprintf("the file %s could not be found", filePath),
			}
		}
	}

	for _, dir := range requisites.RequiredDirs {
		dirPath, _ := GetTerraDaggerPath(&GetTerraDaggerPathOptions{
			WorkDirPath:   workDirPath,
			MountDirPath:  mountPath,
			FileOrDirName: dir,
		})

		dirUtils := utils.NewDirUtils()

		if err := dirUtils.DirExistE(dirPath); err != nil {
			return &RequisitesError{
				ErrWrapped: err,
				Details:    fmt.Sprintf("the dir %s could not be found", dirPath),
			}
		}
	}

	return nil
}
