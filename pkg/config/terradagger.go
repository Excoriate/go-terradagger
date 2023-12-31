package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/o11y"

	"github.com/Excoriate/go-terradagger/pkg/env"

	"github.com/Excoriate/go-terradagger/pkg/utils"
)

var defaultSRCPath = "." // Default root directory is current directory
var defaultSRCPathAbsolute, _ = os.Getwd()
var terraDaggerDir = ".terradagger"
var terraDaggerCacheDir = ".cache"
var terraDaggerExportDir = "export"
var terraDaggerImportDir = "import"

type TerraDaggerConfig interface {
	ResolveWorkspacePaths(srcPath string) (*Workspace, error)
	Configure(options *APIParams) (*TerraDagger, error)
	GetDirs() *Dirs
	GetHostPaths() *HostPaths
	GetEnvVars(injectedEnvVars map[string]string) *EnvVars
}

type TerraDaggerConfigClient struct {
	l o11y.LoggerInterface
}

type TerraDagger struct {
	Paths     *Paths
	Dirs      *Dirs
	APIParams *APIParams
	EnvVars   *EnvVars
}

type Paths struct {
	Workspace *Workspace
	Host      *HostPaths
}

type Workspace struct {
	SRC         string
	SRCAbsolute string
}

type HostPaths struct {
	CurrentDir string
	HomeDir    string
}

type Dirs struct {
	TerraDaggerDir       string
	TerraDaggerCacheDir  string
	TerraDaggerExportDir string
	TerraDaggerImportDir string
}

type APIParams struct {
	Workspace                    string
	EnableStdErrForDaggerBackend bool
	EnvVars                      map[string]string
}

type EnvVars struct {
	HostEnvVars     map[string]string
	InjectedEnvVars map[string]string
}

func NewTerraDaggerConfig(logger o11y.LoggerInterface) TerraDaggerConfig {
	var logImpl o11y.LoggerInterface
	if logger == nil {
		logImpl = o11y.DefaultLogger()
	}

	return &TerraDaggerConfigClient{
		l: logImpl,
	}
}

func newEmptyTerraDaggerCfg() *TerraDagger {
	return &TerraDagger{
		Paths: &Paths{
			Workspace: &Workspace{},
			Host:      &HostPaths{},
		},
		Dirs:      &Dirs{},
		APIParams: &APIParams{},
		EnvVars:   &EnvVars{},
	}
}

func (c *TerraDaggerConfigClient) Configure(options *APIParams) (*TerraDagger, error) {
	cfg := newEmptyTerraDaggerCfg()

	if options == nil {
		c.l.Warn("No options passed to SetUp. Using default options.")
		options = &APIParams{
			Workspace: defaultSRCPath,
			EnvVars:   map[string]string{},
		}
	}

	// 1. Resolve the Workspace path
	workspaceCfg, err := c.ResolveWorkspacePaths(options.Workspace)
	if err != nil {
		return nil, err
	}

	cfg.Paths.Workspace = workspaceCfg
	// 2. Resolve the host paths
	cfg.Paths.Host = c.GetHostPaths()
	// 3. Resolve the directories
	cfg.Dirs = c.GetDirs()
	// 4. Resolve the env vars
	cfg.EnvVars = c.GetEnvVars(options.EnvVars)
	// 5. Resolve the API params
	cfg.APIParams = options

	return cfg, nil
}

func (c *TerraDaggerConfigClient) ResolveWorkspacePaths(srcPath string) (*Workspace, error) {
	if srcPath == "" {
		return &Workspace{
			SRC:         defaultSRCPath,
			SRCAbsolute: defaultSRCPathAbsolute,
		}, nil
	}

	if err := IsAValidTerraDaggerDirRelative(srcPath); err != nil {
		return nil, fmt.Errorf("failed to resolve the 'src' path, the src path is invalid: %w", err)
	}

	absPath, _ := filepath.Abs(srcPath)

	return &Workspace{
		SRC:         srcPath,
		SRCAbsolute: absPath,
	}, nil
}

func (c *TerraDaggerConfigClient) GetDirs() *Dirs {
	return &Dirs{
		TerraDaggerDir:       terraDaggerDir,
		TerraDaggerExportDir: terraDaggerExportDir,
		TerraDaggerCacheDir:  terraDaggerCacheDir,
		TerraDaggerImportDir: terraDaggerImportDir,
	}
}

func (c *TerraDaggerConfigClient) GetHostPaths() *HostPaths {
	dirUtils := utils.DirUtils{}
	return &HostPaths{
		CurrentDir: dirUtils.GetCurrentDir(),
		HomeDir:    dirUtils.GetHomeDir(),
	}
}

func (c *TerraDaggerConfigClient) GetEnvVars(injectedEnvVars map[string]string) *EnvVars {
	return &EnvVars{
		HostEnvVars:     env.GetAllFromHost(),
		InjectedEnvVars: injectedEnvVars,
	}
}
