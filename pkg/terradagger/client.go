package terradagger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/commands"
	"github.com/Excoriate/go-terradagger/pkg/env"
	"github.com/Excoriate/go-terradagger/pkg/erroer"
	"github.com/Excoriate/go-terradagger/pkg/o11y"
	"github.com/Excoriate/go-terradagger/pkg/utils"

	"dagger.io/dagger"
)

var defaultLogOutput io.Writer = os.Stdout // Default log output is stdout
var defaultRootDirRelative = "."           // Default root directory is current directory

type Client struct {
	// Implementation details, and internals.
	Logger      o11y.LoggerInterface
	ID          string
	Ctx         context.Context
	HostEnvVars map[string]string
	// Client           *dagger.Client
	DaggerClient *dagger.Client
	// Dirs
	Paths *PathsCfg
	Dirs  *DirsCfg
}

type DirsCfg struct {
	TerraDaggerDir       string
	TerraDaggerExportDir string
}

type PathsCfg struct {
	TerraDagger     string
	CurrentDir      string
	HomeDir         string
	RootDirRelative string
	RootDirAbsolute string
	MountDirPath    string
	WorkDirPath     string
}

type ClientConfigOptions struct {
	Image        string
	Version      string
	EnvVars      map[string]string
	Workdir      string
	MountDir     string
	ExcludedDirs []string
	// TerraDaggerCMDs     [][]string
	TerraDaggerCMDs commands.DaggerEngineCMDs
	ExportOptions   *ExportOptions
}

type ExportOptions struct {
	ExportToHostPathCustom string // If it's not set, it'll use the default .terradagger directory.
	// These two properties are internally resolved based on the ExportToHostPathCustom.
	useDefaultTerraDaggerExportPath bool
	exportToHostPath                string
	importFromContainerPath         string
}

type ExcludeOptions struct {
	ExcludedDirs []string
	ExcludeFiles []string
}

type OutputFilesConfig struct {
	File                string
	FilePathInHost      string
	FilePathInContainer string
}

type OutputDirsConfig struct {
	Dir                string
	DirPathInHost      string
	DirPathInContainer string
}

type OutputsCfg struct {
	DirsExported       []*OutputDirsConfig
	FilesExportedPaths []*OutputFilesConfig
}

type Core interface {
	// CreateTerraDaggerDirs creates the terradagger directories.
	// The terradagger directories are:
	// 1. The terradagger directory.
	// 2. The export directory.
	CreateTerraDaggerDirs(failIfDirExist bool) error
	// Configure the terradagger, which includes:
	// 1. The dagger client is connected, and properly configured.
	// 2. The image is pulled, and the container is configured.
	// 3. The container is mounted, and the workdir is set.
	Configure(options *ClientConfigOptions) (*dagger.Container, error)

	// Run the terradagger.
	Run(container *dagger.Container) error
	RunWithExport(container *dagger.Container, outputs *OutputsCfg) (string, error)

	// RunAndReturnOutput runs the terradagger and returns the output.
	// RunAndReturnOutput(container *dagger.Container) (string, error)
}

// newDaggerClient creates a new dagger client.
// If no options are passed, the default options are used.
func newDaggerClient(ctx context.Context, options ...dagger.ClientOpt) (*dagger.
	Client, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	var daggerOptions []dagger.ClientOpt

	if len(options) == 0 {
		return dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	}

	// If options are passed, append them to the daggerOptions.
	daggerOptions = append(daggerOptions, options...)

	return dagger.Connect(ctx, daggerOptions...)
}

// ClientOptions are the options for the terradagger client.
// RootDir is the root directory of the terradagger client.
// WithStderrLogInDaggerClient is a flag to enable stderr as log output for the dagger client.
type ClientOptions struct {
	RootDir                     string
	WithStderrLogInDaggerClient bool
}

// New creates a new terradagger client.
// If no options are passed, the default options are used.
func New(ctx context.Context, options *ClientOptions) (*Client, error) {
	logger := o11y.NewLogger(o11y.LoggerOptions{
		EnableJSONHandler: true,
		EnableStdError:    true,
	})

	id := utils.GetUUID()
	logger.Info("Starting terradagger with id", "id", id)

	client := &Client{
		Logger:      logger,
		ID:          id,
		Ctx:         ctx,
		HostEnvVars: env.GetAllFromHost(),
		Paths: &PathsCfg{
			CurrentDir: utils.GetCurrentDir(),
			HomeDir:    utils.GetHomeDir(),
		},
		Dirs: &DirsCfg{
			TerraDaggerDir:       terraDaggerDir,
			TerraDaggerExportDir: terraDaggerExportDir,
		},
	}

	if options == nil {
		logger.Warn("No options were passed to the terradagger client. Using default options.")
		options = &ClientOptions{RootDir: defaultRootDirRelative}
	}

	if err := utils.IsRelative(options.RootDir); err != nil {
		return nil, fmt.Errorf("the root directory %s is not a relative path", options.RootDir)
	}

	rootDirCfg, err := resolveRootDir(options.RootDir)
	if err != nil {
		return nil, fmt.Errorf("terraDagger initialization error: failed to resolve root directory %s: %w", options.RootDir, err)
	}

	client.Paths.RootDirRelative = rootDirCfg.RootDirRelative
	client.Paths.RootDirAbsolute = rootDirCfg.RootDirAbsolute

	// FIXME: Potential abolition of the mount directory - consider refactor this.
	mountDirPath, err := resolveMountDirPath(client.Paths.RootDirRelative)

	if err != nil {
		return nil, fmt.Errorf("terraDagger initialization error: failed to resolve mount directory with root directory %s: %w", options.RootDir, err)
	}

	client.Paths.MountDirPath = mountDirPath

	terraDaggerPath, err := resolveTerraDaggerPath(client.Paths.RootDirRelative)
	if err != nil {
		return nil, fmt.Errorf("terraDagger initialization error: failed to resolve terradagger directory with root directory %s: %w", options.RootDir, err)
	}

	client.Paths.TerraDagger = terraDaggerPath

	logOutput := getLogOutput(options.WithStderrLogInDaggerClient, logger)

	daggerClient, err := newDaggerClient(ctx, dagger.WithLogOutput(logOutput))
	if err != nil {
		return nil, fmt.Errorf("terraDagger initialization error: failed to start dagger client: %w", err)
	}

	client.DaggerClient = daggerClient
	logger.Info("TerraDagger client started successfully.")
	return client, nil
}

// getLogOutput returns the log output for the dagger client.
func getLogOutput(withStderrLogInDaggerClient bool, logger o11y.LoggerInterface) io.Writer {
	if withStderrLogInDaggerClient {
		logger.Debug("Using stderr as log output for dagger client.")
		return os.Stderr
	}
	return defaultLogOutput
}

func (td *Client) Configure(options *ClientConfigOptions) (*dagger.Container, error) {
	dirs := getDirs(td.DaggerClient, options.MountDir, options.Workdir)

	if options.ExportOptions == nil {
		options.ExportOptions = &ExportOptions{}
	}

	tdContainer := NewContainer(td)
	container, err := tdContainer.create(&NewContainerOptions{
		Image:   options.Image,
		Version: options.Version,
	})

	if err != nil {
		return nil, &erroer.ErrTerraDaggerConfigurationError{
			ErrWrapped: err,
			Details:    "the container could not be created",
		}
	}

	if len(options.ExcludedDirs) > 0 {
		td.Logger.Info(fmt.Sprintf("These directories were passed explicitly to be excluded: %v", options.ExcludedDirs))
	}

	// Specific container options.
	container = tdContainer.withDirs(container, dirs.MountDir, dirs.WorkDirPathInContainer,
		options.ExcludedDirs)

	container = tdContainer.withCommands(container, options.TerraDaggerCMDs)

	if len(options.EnvVars) > 0 {
		container = tdContainer.withEnvVars(container, options.EnvVars)
	}

	if tdDirErr := td.CreateTerraDaggerDirs(false); tdDirErr != nil {
		return nil, &erroer.ErrTerraDaggerConfigurationError{
			ErrWrapped: tdDirErr,
			Details:    "the terradagger directories could not be created",
		}
	}

	// This is always set to the <root-dir>/.terradagger/<id>/export.
	exportToHostPathResolved := resolveTerraDaggerExportPath(td.Paths.TerraDagger, td.ID)

	options.ExportOptions.exportToHostPath = exportToHostPathResolved
	options.ExportOptions.importFromContainerPath = dirs.WorkDirPathInContainer
	if options.ExportOptions.ExportToHostPathCustom == "" {
		options.ExportOptions.ExportToHostPathCustom = exportToHostPathResolved
	}

	return container, nil
}

func (td *Client) Run(container *dagger.Container) error {
	_, err := container.Stdout(td.Ctx)
	if err != nil {
		return err
	}

	return nil
}

type RunWithExportOptions struct {
	TargetFilesFromContainer []string
	TargetDirsFromContainer  []string
	FailIfDirNotExist        bool
	FailIfFileNotExist       bool
}

func (td *Client) RunWithExport(container *dagger.Container, exportOptions *RunWithExportOptions, options *ClientConfigOptions) (*OutputsCfg, error) {
	outputs := &OutputsCfg{
		DirsExported:       []*OutputDirsConfig{},
		FilesExportedPaths: []*OutputFilesConfig{},
	}
	// var exportDestinationPathInHost := resolveTerraDaggerExportPath(td.Paths.TerraDagger, td.ID)

	if len(exportOptions.TargetDirsFromContainer) > 0 {
		for _, dir := range exportOptions.TargetDirsFromContainer {
			dirPathImport := filepath.Join(options.ExportOptions.importFromContainerPath, dir)
			dirPathExport := filepath.Join(options.ExportOptions.exportToHostPath, filepath.Base(dir))

			td.Logger.Info("Exporting directory", "directory", dir, "dirPathImport", dirPathImport, "dirPathExport", dirPathExport)

			if _, err := container.Directory(dirPathImport).
				Export(td.Ctx, dirPathExport); err != nil {
				td.Logger.Error("Failed to export directory", "directory", dir, "error", err)
				return nil, fmt.Errorf("failed to export directory %s: %w", dir, err)
			}

			outputs.DirsExported = append(outputs.DirsExported, &OutputDirsConfig{
				Dir:                dir,
				DirPathInContainer: dirPathImport,
				DirPathInHost:      filepath.Join(dirPathExport, filepath.Base(dir)),
			})
		}
	}

	if len(exportOptions.TargetFilesFromContainer) > 0 {
		for _, file := range exportOptions.TargetFilesFromContainer {
			filePathImport := filepath.Join(options.ExportOptions.importFromContainerPath, file)
			filePathExport := filepath.Join(options.ExportOptions.exportToHostPath, filepath.Base(file))

			td.Logger.Info("Exporting file", "file", file, "filePathImport", filePathImport, "filePathExport", filePathExport)

			if _, err := container.File(filePathImport).
				Export(td.Ctx, filePathExport); err != nil {
				td.Logger.Error("Failed to export file", "file", file, "error", err)
				return nil, fmt.Errorf("failed to export file %s: %w", file, err)
			}

			outputs.FilesExportedPaths = append(outputs.FilesExportedPaths, &OutputFilesConfig{
				File:                file,
				FilePathInContainer: filePathImport,
				FilePathInHost:      filepath.Join(filePathExport, filepath.Base(file)),
			})
		}
	}

	return outputs, nil
}

func (td *Client) CreateTerraDaggerDirs(failIfDirExist bool) error {
	td.Logger.Info("Creating terradagger directories.")
	terraDaggerExportPath := resolveTerraDaggerExportPath(td.Paths.TerraDagger, td.ID)

	if failIfDirExist {
		if utils.DirExist(td.Paths.TerraDagger) {
			return fmt.Errorf("terradagger directory already exists: %s", td.Paths.TerraDagger)
		}
		if utils.DirExist(terraDaggerExportPath) {
			return fmt.Errorf("export path already exists: %s", terraDaggerExportPath)
		}
	}

	err := os.MkdirAll(td.Paths.TerraDagger, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create terradagger directory: %w", err)
	}

	err = os.MkdirAll(terraDaggerExportPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create export path: %w", err)
	}

	td.Logger.Info("Terradagger directories created successfully.")
	return nil
}

// RunAndReturnOutput TODO: Implement this method.
// func (p *Client) RunAndReturnOutput(container *dagger.Container) (string, error) {
// 	return "", nil
// }
