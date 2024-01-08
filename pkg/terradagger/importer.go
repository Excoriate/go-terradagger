package terradagger

import (
	"fmt"

	"dagger.io/dagger"
)

type ContainerImporter interface {
	AddDataToImportInContainer(container *dagger.Container,
		options *DataTransferToContainer) (*dagger.Container, error)
}

type ContainerImporterImpl struct {
	td *TD
}

type ContainerImporterError struct {
	ErrWrapped error
	Details    string
}

const ContainerImporterErrPrefix = "the container importer failed to manage the import configuration"

func (e *ContainerImporterError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", ContainerImporterErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", ContainerImporterErrPrefix, e.Details, e.ErrWrapped.Error())
}

func NewContainerImporter(td *TD) ContainerImporter {
	return &ContainerImporterImpl{
		td: td,
	}
}

func validateDataTransferOptions(options *DataTransferToContainer) error {
	if options == nil {
		return &ContainerImporterError{
			Details: "the data transfer options are nil",
		}
	}

	if len(options.Files) == 0 && len(options.Dirs) == 0 {
		return &ContainerImporterError{
			Details: "there are no files or dirs to import",
		}
	}

	for _, file := range options.Files {
		if file.DestinationPathInContainer == "" {
			return &ContainerImporterError{
				Details: "one of the files to import is empty",
			}
		}

		if file.SourcePathInHostAbs == "" {
			return &ContainerImporterError{
				Details: "one of the files to import is empty",
			}
		}
	}

	for _, dir := range options.Dirs {
		if dir.DestinationPathInContainer == "" {
			return &ContainerImporterError{
				Details: "one of the dirs to import is empty",
			}
		}

		if dir.SourcePathInHostAbs == "" {
			return &ContainerImporterError{
				Details: "one of the dirs to import is empty",
			}
		}
	}

	return nil
}

func (c *ContainerImporterImpl) AddDataToImportInContainer(container *dagger.Container, options *DataTransferToContainer) (*dagger.Container, error) {
	if container == nil {
		return nil, &ContainerImporterError{
			Details: "the container is nil",
		}
	}

	if err := validateDataTransferOptions(options); err != nil {
		return nil, &ContainerImporterError{
			ErrWrapped: err,
			Details:    "the data transfer options are invalid",
		}
	}

	for _, file := range options.Files {
		daggerFile := c.td.DaggerBackend.Host().File(file.SourcePathInHostAbs)
		container = container.WithFile(file.DestinationPathInContainer, daggerFile)
	}

	for _, dir := range options.Dirs {
		daggerDir := c.td.DaggerBackend.Host().Directory(dir.SourcePathInHostAbs)
		container = container.WithDirectory(dir.DestinationPathInContainer, daggerDir)
	}

	return container, nil
}
