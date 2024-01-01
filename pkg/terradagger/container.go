package terradagger

import (
	"fmt"
	"strings"

	"github.com/Excoriate/go-terradagger/pkg/o11y"

	"github.com/Excoriate/go-terradagger/pkg/commands"
	"github.com/Excoriate/go-terradagger/pkg/config"
	"github.com/Excoriate/go-terradagger/pkg/erroer"
	"github.com/Excoriate/go-terradagger/pkg/utils"

	"dagger.io/dagger"
)

const defaultImageVersion = "latest"

type CreateNewContainerOptions struct {
	Image   string
	Version string
}

type ContainerValidator interface {
	validate(options *CreateNewContainerOptions) error
}

type ContainerValidatorImpl struct {
	logger o11y.LoggerInterface
}

func (cv *ContainerValidatorImpl) validate(options *CreateNewContainerOptions) error {
	if options == nil {
		return fmt.Errorf("options cannot be nil")
	}

	if options.Image == "" {
		return fmt.Errorf("the image while creating a new runtime cannot be nil or empty")
	}

	if strings.Contains(options.Image, ":") {
		return fmt.Errorf("the image while creating a new runtime cannot include the version")
	}

	if strings.Contains(options.Image, ":") {
		return fmt.Errorf("the image while creating a new runtime cannot include the version")
	}

	return nil
}

func newContainerValidator(logger o11y.LoggerInterface) ContainerValidator {
	if logger == nil {
		logger = o11y.DefaultLogger()
	}
	return &ContainerValidatorImpl{
		logger: logger,
	}
}

type ContainerFactory interface {
	createNewContainer(options *CreateNewContainerOptions) (*Container, error)
	buildImageName(image, version string) string
	withEnvVars(container *dagger.Container, envVars map[string]string) *dagger.Container
	withDirs(container *dagger.Container, mountDir *dagger.Directory,
		workDirPath string, excludeDirsExtra []string) *dagger.Container
	withCommands(container *dagger.Container, cmds commands.DaggerEngineCMDs) *dagger.Container
}

type ContainerClient struct {
	td *TD
}

type Container struct {
	Image            string
	Version          string
	ImageWithVersion string
	DaggerContainer  *dagger.Container
}

func newContainerClient(td *TD) ContainerFactory {
	return &ContainerClient{
		td: td,
	}
}

func (cc *ContainerClient) buildImageName(image, version string) string {
	if version == "" {
		version = defaultImageVersion
	}

	return fmt.Sprintf("%s:%s", image, version)
}

func (cc *ContainerClient) createNewContainer(options *CreateNewContainerOptions) (*Container, error) {
	if options == nil {
		return nil, &erroer.ErrTerraDaggerInvalidArgumentError{
			Details: "options cannot be nil",
		}
	}

	if options.Image == "" {
		return nil, &erroer.ErrTerraDaggerInvalidArgumentError{
			Details: "the image while creating a new runtime cannot be nil or empty",
		}
	}

	imageWithVersion := cc.buildImageName(options.Image, options.Version)

	tdContainer := &Container{
		Image:            options.Image,
		Version:          options.Version,
		ImageWithVersion: imageWithVersion,
		DaggerContainer:  cc.td.DaggerBackend.Container().From(imageWithVersion),
	}

	tdContainer.DaggerContainer = cc.td.DaggerBackend.Container().From(imageWithVersion)

	cc.td.Logger.Info(fmt.Sprintf("Created a new runtime with image %s", imageWithVersion))
	return tdContainer, nil
}

func (cc *ContainerClient) withEnvVars(container *dagger.Container, envVars map[string]string) *dagger.Container {
	for key, value := range envVars {
		container = container.WithEnvVariable(key, value)
	}

	return container
}

func (cc *ContainerClient) withDirs(container *dagger.Container, mountDir *dagger.Directory,
	workDirPath string, excludeDirsExtra []string) *dagger.Container {
	excludeDirsDefault := cc.td.Config.Dagger.Excluded.ExcludedDirs
	excludeFilesDefault := cc.td.Config.Dagger.Excluded.ExcludedFiles

	if len(excludeDirsExtra) > 0 {
		cc.td.Logger.Info(fmt.Sprintf("Excluding extra dirs: %s", excludeDirsExtra))
	}

	container = container.WithDirectory(config.MountPathPrefixInDagger, mountDir,
		dagger.ContainerWithDirectoryOpts{
			Exclude: utils.MisSlices(excludeDirsDefault, excludeFilesDefault, excludeDirsExtra),
		})

	container = container.WithWorkdir(workDirPath)
	return container
}

func (cc *ContainerClient) withCommands(container *dagger.Container, cmds commands.DaggerEngineCMDs) *dagger.Container {
	for _, cmds := range cmds {
		for _, cmd := range cmds {
			container = container.WithExec(cmd)
		}
	}

	return container
}
