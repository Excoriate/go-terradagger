package config

import (
	"os"
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/env"
	"github.com/Excoriate/go-terradagger/pkg/utils"
)

const (
	defaultWorkspace     = "."
	terraDaggerDir       = ".terradagger"
	mountPrefix          = "/mnt"
	terraDaggerExportDir = "export"
	awsKeysPrefix        = "AWS_"
	tfVarExtension       = ".tfvars"
)

var (
	excludedDirsDefault = []string{".git", ".terradagger",
		"dist/**", "node_modules/**", ".cache"}

	excludedFilesDefault           = []string{".gitignore"}
	allowedTerraformFileExtensions = []string{".tf", ".hcl", ".json", ".tfvars"}
)

type Config interface {
	GetTerraDaggerDir() string
	GetTerraDaggerExportDir() string
	GetWorkspace() string
	GetWorkspaceAbs() string
	GetWorkspaceDefault() string
	GetWorkspaceDefaultAbs() string
	GetMountPrefix() string
	GetExcludedDirs() []string
	GetExcludedFiles() []string
	GetCurrentDir() string
	GetHomeDir() string
	GetHostEnvVars() map[string]string
	GetTerraformEnvVars() map[string]string
	GetAWSEnvVars() map[string]string
	GetAllEnvVars() map[string]string
	GetAllowedTerraformFileExtensions() []string

	GetEnvVarsByKeys(keys []string) map[string]string
	GetTfVarsExtension() string
}

type Options struct {
	workspace     string
	envVars       map[string]string
	excludedDirs  []string
	excludedFiles []string
}

func New(workspace string, envVars map[string]string, excludeDirs, excludedFiles []string) Config {
	return &Options{
		workspace:     workspace,
		envVars:       envVars,
		excludedDirs:  excludeDirs,
		excludedFiles: excludedFiles,
	}
}

func (o *Options) GetTerraDaggerDir() string {
	return terraDaggerDir
}

func (o *Options) GetTerraDaggerExportDir() string {
	return terraDaggerExportDir
}

func (o *Options) GetWorkspace() string {
	if o.workspace == "" {
		return defaultWorkspace
	}

	return o.workspace
}

func (o *Options) GetWorkspaceDefault() string {
	return defaultWorkspace
}

func (o *Options) GetWorkspaceDefaultAbs() string {
	return o.GetCurrentDir()
}

func (o *Options) GetWorkspaceAbs() string {
	workspaceAbs, _ := filepath.Abs(o.GetWorkspace())
	return workspaceAbs
}

func (o *Options) GetMountPrefix() string {
	return mountPrefix
}

func (o *Options) GetExcludedDirs() []string {
	return append(excludedDirsDefault, o.excludedDirs...)
}

func (o *Options) GetExcludedFiles() []string {
	return append(excludedFilesDefault, o.excludedFiles...)
}

func (o *Options) GetCurrentDir() string {
	currentDir, _ := os.Getwd()
	return currentDir
}

func (o *Options) GetHomeDir() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}

func (o *Options) GetHostEnvVars() map[string]string {
	return env.GetAllFromHost()
}

func (o *Options) GetTerraformEnvVars() map[string]string {
	envVars, _ := env.GetAllEnvVarsWithPrefix("TF_VAR_")
	return utils.MergeMaps(envVars, o.envVars)
}

func (o *Options) GetAllEnvVars() map[string]string {
	return utils.MergeMaps(o.GetHostEnvVars(), o.GetTerraformEnvVars())
}

func (o *Options) GetAllowedTerraformFileExtensions() []string {
	return allowedTerraformFileExtensions
}

func (o *Options) GetAWSEnvVars() map[string]string {
	awsVars, _ := env.GetAllEnvVarsWithPrefix(awsKeysPrefix)
	return awsVars
}

func (o *Options) GetEnvVarsByKeys(keys []string) map[string]string {
	envVars := map[string]string{}
	for _, key := range keys {
		envVar, _ := env.GetEnvVarByKey(key)
		envVars[key] = envVar
	}

	return envVars
}

func (o *Options) GetTfVarsExtension() string {
	return tfVarExtension
}
