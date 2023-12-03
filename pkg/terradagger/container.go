package terradagger

import (
	"fmt"

	"github.com/Excoriate/go-terradagger/pkg/commands"
	"github.com/Excoriate/go-terradagger/pkg/config"
	"github.com/Excoriate/go-terradagger/pkg/errors"
	"github.com/Excoriate/go-terradagger/pkg/utils"

	"dagger.io/dagger"
)

const defaultImageVersion = "latest"

type NewContainerOptions struct {
	Image   string
	Version string
}

type ContainerFactory interface {
	create(options *NewContainerOptions) (*dagger.Container, error)
	// withEnvVars(envVars map[string]string) *dagger.Container
	withDirs(container *dagger.Container, mountDir *dagger.Directory,
		workDirPath string) *dagger.Container
	withCommands(container *dagger.Container, commands [][]string) *dagger.Container
	withEnvVars(container *dagger.Container, envVars map[string]string) *dagger.Container
}

type Container struct {
	client *Client
}

func NewContainer(td *Client) *Container {
	return &Container{
		client: td,
	}
}

func buildImageName(image, version string) string {
	if version == "" {
		version = defaultImageVersion
	}

	return fmt.Sprintf("%s:%s", image, version)
}

func (c *Container) create(options *NewContainerOptions) (*dagger.Container, error) {
	if options == nil {
		return nil, &errors.ErrTerraDaggerInvalidArgumentError{
			Details: "options cannot be nil",
		}
	}

	if options.Image == "" {
		return nil, &errors.ErrTerraDaggerInvalidArgumentError{
			Details: "the image while creating a new container cannot be nil or empty",
		}
	}

	imageWithVersion := buildImageName(options.Image, options.Version)
	c.client.Logger.Info(fmt.Sprintf("Creating a new container with image: %s", imageWithVersion))

	return c.client.DaggerClient.Container().From(imageWithVersion), nil
}

func (c *Container) withDirs(container *dagger.Container, mountDir *dagger.Directory,
	workDirPath string, excludeDirsExtra []string) *dagger.Container {
	container = container.WithDirectory(config.MountPathPrefixInDagger, mountDir,
		dagger.ContainerWithDirectoryOpts{
			Exclude: utils.MisSlices(config.ExcludedDirsDefault, config.ExcludedDirsTerraform, excludeDirsExtra),
		})
	container = container.WithWorkdir(workDirPath)

	return container
}

func (c *Container) withCommands(container *dagger.Container, cmds commands.DaggerEngineCMDs) *dagger.Container {
	for _, cmds := range cmds {
		for _, cmd := range cmds {
			container = container.WithExec(cmd)
		}
	}

	return container
}

func (c *Container) withEnvVars(container *dagger.Container, envVars map[string]string) *dagger.Container {
	for key, value := range envVars {
		container = container.WithEnvVariable(key, value)
	}

	return container
}
