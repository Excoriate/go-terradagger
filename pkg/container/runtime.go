package container

import (
	"dagger.io/dagger"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
)

type runtime struct {
	container Container
	td        *terradagger.TD
}

type Command []string

type Runtime interface {
	CreateContainer() *dagger.Container
	OverrideWorkdir(workdir string, container *dagger.Container) *dagger.Container
	AddCommands(commands []Command, container *dagger.Container) *dagger.Container
	RunAndGetStdout(container *dagger.Container) (string, error)
}

func New(container Container, td *terradagger.TD) Runtime {
	return &runtime{
		container: container,
		td:        td,
	}
}

func (r *runtime) CreateContainer() *dagger.Container {
	containerImageCfg := r.container.GetImageConfig()
	containerImage := containerImageCfg.GetTerraformContainerImage()
	mntPathPrefix := r.container.GetMountPathPrefix()
	mountDir := r.container.GetMountDir(r.td.Engine.GetEngine())

	base := r.td.Engine.GetEngine().Container().From(containerImage).
		WithMountedDirectory(mntPathPrefix, mountDir)

	if !r.container.IsKeepEntryPoint() {
		base = base.WithoutEntrypoint()
	}

	if r.container.GetWorkdir() != "" {
		base = base.WithWorkdir(r.container.GetWorkdir())
	}

	return base
}

func (r *runtime) OverrideWorkdir(workdir string, container *dagger.Container) *dagger.Container {
	return container.WithWorkdir(workdir)
}

func (r *runtime) AddCommands(commands []Command, container *dagger.Container) *dagger.Container {
	if len(commands) > 0 {
		for _, cmd := range commands {
			container = container.WithExec(cmd)
		}
	}

	return container
}

func (r *runtime) RunAndGetStdout(container *dagger.Container) (string, error) {
	return container.Stdout(r.td.Ctx)
}
