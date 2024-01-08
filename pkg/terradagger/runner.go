package terradagger

import (
	"fmt"

	"dagger.io/dagger"
)

type Runner interface {
	RunOnly(c *ClientInstance) error
	getContainer(c *ClientInstance) *dagger.Container
	RunWithExport(c *ClientInstance, options *RuntWithExportOptions) error
}

type RunnerImpl struct {
	td *TD
}

func NewRunner(td *TD) Runner {
	return &RunnerImpl{
		td: td,
	}
}

type RunnerError struct {
	ErrWrapped error
	Details    string
}

const RunnerErrPrefix = "the runner failed when it tried to run the terradagger client"

func (e *RunnerError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", RunnerErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", RunnerErrPrefix, e.Details, e.ErrWrapped.Error())
}

func (r *RunnerImpl) getContainer(c *ClientInstance) *dagger.Container {
	if c == nil {
		return nil
	}

	return c.runtimeContainer.DaggerContainer
}

func (r *RunnerImpl) RunOnly(c *ClientInstance) error {
	if c == nil {
		return &RunnerError{
			Details: "the client instance is nil",
		}
	}

	container := r.getContainer(c)

	_, err := container.Stdout(c.td.Ctx)
	if err != nil {
		return &RunnerError{
			ErrWrapped: err,
			Details:    "the runner failed to get the stdout of the container",
		}
	}

	return nil
}

type RuntWithExportOptions struct {
	// TODO: Complete this implementation.
	// CustomDestinationPathInHost is the custom path in the host where the file or directory will be exported.
	// If this is not set, the default path will be used.
	CustomDestinationPathInHost string
}

func (r *RunnerImpl) RunWithExport(c *ClientInstance, options *RuntWithExportOptions) error {
	if c == nil {
		return &RunnerError{
			Details: "the client instance is nil",
		}
	}

	if options == nil {
		options = &RuntWithExportOptions{}
		r.td.Logger.Info("the options passed to the RunWithExport API are nil, using default ones...")
	}

	transferCfg := c.Config.runtime.containerHostInterop
	filesToExport := transferCfg.transferToHost.Files
	dirsToExport := transferCfg.transferToHost.Dirs
	container := r.getContainer(c)

	exp := NewExporter(r.td)

	if exp.IsAdvancedExportEnabled(c) {
		return exp.ExportAdvance(c, &ExportAdvanceOptions{
			WorkDirPathInDagger:   transferCfg.transferToHost.WorkDirPath,
			TerraDaggerCachePath:  c.Config.Paths.CachePathAbs,
			TerraDaggerExportPath: c.Config.Paths.ExportPathAbs,
			FilesToExport:         transferCfg.transferToHost,
			DirsToExport:          transferCfg.transferToHost,
		})
	}

	// Export a single file or directory to the host. Simple scenario.
	if len(filesToExport) > 0 {
		file := filesToExport[0]

		if _, err := container.File(file.SourcePathInContainer).
			Export(c.td.Ctx, file.DestinationPathInHostAbs); err != nil {
			return &RunnerError{
				ErrWrapped: err,
				Details:    "the runner failed to export the file from the container",
			}
		}
	}

	if len(dirsToExport) > 0 {
		dir := dirsToExport[0]

		if _, err := container.Directory(dir.SourcePathInContainer).
			Export(c.td.Ctx, dir.DestinationPathInHostAbs); err != nil {
			return &RunnerError{
				ErrWrapped: err,
				Details:    "the runner failed to export the directory from the container",
			}
		}
	}

	return nil
}
