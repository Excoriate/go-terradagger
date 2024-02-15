package container

import (
	"os"
	"reflect"
	"testing"
)

func TestConfig_GetMountDirPath(t *testing.T) {
	currentDir, _ := os.Getwd()
	tests := []struct {
		name         string
		MountPathAbs string
		want         string
	}{
		{"default path", "", currentDir},
		{"absolute path", "/tmp", "/tmp"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				MountPathAbs: tt.MountPathAbs,
			}
			if got := c.GetMountDirPath(); got != tt.want {
				t.Errorf("Config.GetMountDirPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetMountPathPrefix(t *testing.T) {
	tests := []struct {
		name            string
		MountPathPrefix string
		want            string
	}{
		{"default prefix", "", "/mnt"},
		{"custom prefix", "/custom", "/custom"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				MountPathPrefix: tt.MountPathPrefix,
			}
			if got := c.GetMountPathPrefix(); got != tt.want {
				t.Errorf("Config.GetMountPathPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetWorkdir(t *testing.T) {
	c := &Config{
		Workdir:         "work",
		MountPathPrefix: "/mnt",
	}
	want := "/mnt/work"
	if got := c.GetWorkdir(); got != want {
		t.Errorf("Config.GetWorkdir() = %v, want %v", got, want)
	}
}

func TestConfig_IsKeepEntryPoint(t *testing.T) {
	c := &Config{KeepEntryPoint: true}
	if got := c.IsKeepEntryPoint(); got != true {
		t.Errorf("Config.IsKeepEntryPoint() = %v, want %v", got, true)
	}
}

func TestConfig_GetEnvVars(t *testing.T) {
	envVars := map[string]string{"VAR": "VALUE"}
	c := &Config{EnvVars: envVars}
	if got := c.GetEnvVars(); !reflect.DeepEqual(got, envVars) {
		t.Errorf("Config.GetEnvVars() = %v, want %v", got, envVars)
	}
}

func TestConfig_IsCacheInvalidated(t *testing.T) {
	c := &Config{InvalidateCache: true}
	if got := c.IsCacheInvalidated(); got != true {
		t.Errorf("Config.IsCacheInvalidated() = %v, want %v", got, true)
	}
}

func TestConfig_IsPrivateGitSupportEnabled(t *testing.T) {
	c := &Config{AddPrivateGitSupport: true}
	if got := c.IsPrivateGitSupportEnabled(); got != true {
		t.Errorf("Config.IsPrivateGitSupportEnabled() = %v, want %v", got, true)
	}
}

func TestConfig_GetCacheBusterEnvVar(t *testing.T) {
	c := &Config{}
	want := cacheBusterEnvVar
	if got := c.GetCacheBusterEnvVar(); !reflect.DeepEqual(got, want) {
		t.Errorf("Config.GetCacheBusterEnvVar() = %v, want %v", got, want)
	}
}

func TestConfig_GetGitSSHEnvVar(t *testing.T) {
	c := &Config{}
	want := gitSSHEnvVar
	if got := c.GetGitSSHEnvVar(); !reflect.DeepEqual(got, want) {
		t.Errorf("Config.GetGitSSHEnvVar() = %v, want %v", got, want)
	}
}

func TestConfig_GetSSHAuthSockEnvVar(t *testing.T) {
	c := &Config{}
	want := sshAuthSockEnvVar
	if got := c.GetSSHAuthSockEnvVar(); !reflect.DeepEqual(got, want) {
		t.Errorf("Config.GetSSHAuthSockEnvVar() = %v, want %v", got, want)
	}
}
