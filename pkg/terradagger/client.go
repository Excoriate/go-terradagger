package terradagger

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Excoriate/go-terradagger/pkg/commands"
	"github.com/Excoriate/go-terradagger/pkg/env"
	"github.com/Excoriate/go-terradagger/pkg/erroer"
	"github.com/Excoriate/go-terradagger/pkg/o11y"
	"github.com/Excoriate/go-terradagger/pkg/utils"

	"dagger.io/dagger"
)

var defaultLogOutput io.Writer = os.Stdout // Default log output is stdout
var defaultRootDir = "."                   // Default root directory is current directory

type Client struct {
	// Implementation details, and internals.
	Logger          o11y.LoggerInterface
	ID              string
	Ctx             context.Context
	CurrentDir      string
	RootDirRelative string
	HomeDir         string
	HostEnvVars     map[string]string
	// Client           *dagger.Client
	MountDir     string
	DaggerClient *dagger.Client
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
}

type Core interface {
	// Configure the terradagger, which includes:
	// 1. The dagger client is connected, and properly configured.
	// 2. The image is pulled, and the container is configured.
	// 3. The container is mounted, and the workdir is set.
	Configure(options *ClientConfigOptions) (*dagger.Container, error)

	// Run the terradagger.
	Run(container *dagger.Container) error

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
		CurrentDir:  utils.GetCurrentDir(),
		HomeDir:     utils.GetHomeDir(),
		HostEnvVars: env.GetAllFromHost(),
	}

	if options == nil {
		logger.Warn("No options were passed to the terradagger client. Using default options.")
		options = &ClientOptions{RootDir: defaultRootDir}
	}

	client.RootDirRelative = options.RootDir
	mountDirPath, err := resolveMountDirPath(options.RootDir)
	if err != nil {
		return nil, fmt.Errorf("terraDagger initialization error: failed to resolve mount directory with root directory %s: %w", options.RootDir, err)
	}

	client.MountDir = mountDirPath
	logger.Info("Mount directory resolved", "mountDir", client.MountDir)

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

func (p *Client) Configure(options *ClientConfigOptions) (*dagger.Container, error) {
	dirs := getDirs(p.DaggerClient, options.MountDir, options.Workdir)

	tdContainer := NewContainer(p)
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
		p.Logger.Info(fmt.Sprintf("These directories were passed explicitly to be excluded: %v", options.ExcludedDirs))
	}

	container = tdContainer.withDirs(container, dirs.MountDir, dirs.WorkDirPathInContainer,
		options.ExcludedDirs)

	container = tdContainer.withCommands(container, options.TerraDaggerCMDs)

	if len(options.EnvVars) > 0 {
		container = tdContainer.withEnvVars(container, options.EnvVars)
	}

	return container, nil
}

func (p *Client) Run(container *dagger.Container) error {
	_, err := container.Stdout(p.Ctx)
	if err != nil {
		return err
	}

	return nil
}

// RunAndReturnOutput TODO: Implement this method.
// func (p *Client) RunAndReturnOutput(container *dagger.Container) (string, error) {
// 	return "", nil
// }
