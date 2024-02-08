package container

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Excoriate/go-terradagger/pkg/utils"

	"dagger.io/dagger"
)

type Config struct {
	Workdir              string
	MountPathAbs         string
	EnvVars              map[string]string
	AddPrivateGitSupport bool
	ContainerImage       Image
	MountPathPrefix      string
	KeepEntryPoint       bool
	InvalidateCache      bool
}

type EnvVar struct {
	Name   string
	Value  string
	Expand bool
}

var (
	cacheBusterEnvVar = EnvVar{
		Name:   "CACHE_BUSTER",
		Value:  fmt.Sprintf("%d", time.Now().Unix()),
		Expand: false,
	}

	gitSSHEnvVar = EnvVar{
		Name:   "GIT_SSH_COMMAND",
		Value:  utils.GetSSHGitSecureConnectCommand(),
		Expand: false,
	}

	sshAuthSockEnvVar = EnvVar{
		Name:   "SSH_AUTH_SOCK",
		Value:  utils.GetSSHAuthSock(),
		Expand: false,
	}
)

type Container interface {
	GetMountDirPath() string
	GetMountDir(client *dagger.Client) *dagger.Directory
	GetDir(dirPathAbs string, client *dagger.Client) *dagger.Directory
	GetMountPathPrefix() string
	GetImageConfig() Image

	GetWorkdir() string

	IsKeepEntryPoint() bool
	GetEnvVars() map[string]string
	IsCacheInvalidated() bool
	IsPrivateGitSupportEnabled() bool
	GetCacheBusterEnvVar() EnvVar
	GetGitSSHEnvVar() EnvVar
	GetSSHAuthSockEnvVar() EnvVar
}

func (o *Config) GetMountDir(client *dagger.Client) *dagger.Directory {
	return client.Host().Directory(o.GetMountDirPath())
}

func (o *Config) GetMountDirPath() string {
	if o.MountPathAbs == "" {
		currentDir, _ := os.Getwd()
		return currentDir
	}

	mountDirAbs, _ := filepath.Abs(o.MountPathAbs)
	return mountDirAbs
}

func (o *Config) GetDir(dirPathAbs string, client *dagger.Client) *dagger.Directory {
	return client.Host().Directory(dirPathAbs)
}

func (o *Config) GetMountPathPrefix() string {
	if o.MountPathPrefix == "" {
		return "/mnt"
	}

	return o.MountPathPrefix
}

func (o *Config) GetImageConfig() Image {
	return o.ContainerImage
}

func (o *Config) GetWorkdir() string {
	return fmt.Sprintf("%s/%s", o.GetMountPathPrefix(), o.Workdir)
}

func (o *Config) IsKeepEntryPoint() bool {
	return o.KeepEntryPoint
}

func (o *Config) GetEnvVars() map[string]string {
	return o.EnvVars
}

func (o *Config) IsCacheInvalidated() bool {
	return o.InvalidateCache
}

func (o *Config) IsPrivateGitSupportEnabled() bool {
	return o.AddPrivateGitSupport
}

func (o *Config) GetCacheBusterEnvVar() EnvVar {
	return cacheBusterEnvVar
}

func (o *Config) GetGitSSHEnvVar() EnvVar {
	return gitSSHEnvVar
}

func (o *Config) GetSSHAuthSockEnvVar() EnvVar {
	return sshAuthSockEnvVar
}
