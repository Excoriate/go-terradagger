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
	ForwardUnixSockets(container *dagger.Container) *dagger.Container
	AddEnvVars(envVars map[string]string, container *dagger.Container) *dagger.Container
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

	if r.container.IsPrivateGitSupportEnabled() {
		base = r.ForwardUnixSockets(base)
	}

	if !r.container.IsKeepEntryPoint() {
		base = base.WithoutEntrypoint()
	}

	if r.container.GetWorkdir() != "" {
		base = base.WithWorkdir(r.container.GetWorkdir())
	}

	if r.container.IsCacheInvalidated() {
		cacheBuster := r.container.GetCacheBusterEnvVar()
		base = base.WithEnvVariable(cacheBuster.Name, cacheBuster.Value)
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

func (r *runtime) ForwardUnixSockets(container *dagger.Container) *dagger.Container {
	unixSocketPath := r.td.Engine.GetEngine().Host().UnixSocket(r.container.GetSSHAuthSockEnvVar().Value)

	return container.WithEnvVariable(r.container.GetGitSSHEnvVar().Name, r.container.GetGitSSHEnvVar().Value).
		WithEnvVariable(r.container.GetSSHAuthSockEnvVar().Name, r.container.GetSSHAuthSockEnvVar().Value).
		WithUnixSocket(r.container.GetSSHAuthSockEnvVar().Value, unixSocketPath)
}

func (r *runtime) AddEnvVars(envVars map[string]string, container *dagger.Container) *dagger.Container {
	for k, v := range envVars {
		container = container.WithEnvVariable(k, v)
	}

	return container
}
