package terradagger

import (
	"fmt"
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/utils"
)

type Exporter interface {
	IsAdvancedExportEnabled(c *ClientInstance) bool
	validateExportOptions(options *ExportAdvanceOptions) error
	ExportAdvance(c *ClientInstance, options *ExportAdvanceOptions) error
	FilterContentInCache(cachePathAbs string, files, dirs []string) (*CacheContent, error)
}

type ExporterImpl struct {
	td *TD
}

func NewExporter(td *TD) Exporter {
	return &ExporterImpl{
		td: td,
	}
}

type ExporterError struct {
	ErrWrapped error
	Details    string
}

const ExporterErrPrefix = "the exporter failed when tried to export or manage the export configuration"

func (e *ExporterError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", ExporterErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", ExporterErrPrefix, e.Details, e.ErrWrapped.Error())
}

// IsAdvancedExportEnabled validate if the runtimeContainer interop configuration (data transfer)
// either on the files or dirs has more than one occurrence in either; if so,
// it's considered advanced, otherwise, it's considered simple.
func (ei *ExporterImpl) IsAdvancedExportEnabled(c *ClientInstance) bool {
	if c == nil {
		return false
	}

	transferCfg := c.Config.runtime.containerHostInterop

	totalTransfers := len(transferCfg.transferToHost.Files) + len(transferCfg.transferToHost.Dirs)
	return totalTransfers != 1
}

type ExportAdvanceOptions struct {
	WorkDirPathInDagger   string
	TerraDaggerCachePath  string
	TerraDaggerExportPath string
	FilesToExport         *DataTransferToHost
	DirsToExport          *DataTransferToHost
}

func (ei *ExporterImpl) validateExportOptions(options *ExportAdvanceOptions) error {
	if options == nil {
		return &ExporterError{
			Details: "the export options are nil",
		}
	}

	if options.WorkDirPathInDagger == "" {
		return &ExporterError{
			Details: "the work dir path in dagger is empty",
		}
	}

	if options.TerraDaggerCachePath == "" {
		return &ExporterError{
			Details: "the terra dagger cache path is empty",
		}
	}

	if options.TerraDaggerExportPath == "" {
		return &ExporterError{
			Details: "the terra dagger export path is empty",
		}
	}

	if options.FilesToExport == nil {
		return &ExporterError{
			Details: "the files to export are nil",
		}
	}

	if options.DirsToExport == nil {
		return &ExporterError{
			Details: "the dirs to export are nil",
		}
	}

	return nil
}

func (ei *ExporterImpl) ExportAdvance(c *ClientInstance, options *ExportAdvanceOptions) error {
	if err := ei.validateExportOptions(options); err != nil {
		return &ExporterError{
			ErrWrapped: err,
			Details:    "the export options are invalid",
		}
	}

	container := c.runtimeContainer.DaggerContainer

	// Resolving names, and paths.
	workDirPath := options.WorkDirPathInDagger
	workDirPathName := filepath.Base(workDirPath)
	workDirPathNameInCache := fmt.Sprintf("%s-tmp", workDirPathName)
	workDirPathInCache := filepath.Join(options.TerraDaggerCachePath, workDirPathNameInCache)

	// Exporting the work dir from the container to the cache.
	if _, err := container.Directory(workDirPath).
		Export(c.td.Ctx, workDirPathInCache); err != nil {
		return &ExporterError{
			ErrWrapped: err,
			Details:    "the exporter failed to export the work dir from the container",
		}
	}

	ei.td.Logger.Info(fmt.Sprintf("the work dir %s was exported from the container to the cache", workDirPathInCache))

	// Filtering the cache content, if the export was successful.
	var filesToFoundInCache []string
	var dirsToExportPaths []string

	for _, file := range options.FilesToExport.Files {
		fileName := filepath.Base(file.DestinationPathInHostAbs)
		filesToFoundInCache = append(filesToFoundInCache, fileName)
	}

	for _, dir := range options.DirsToExport.Dirs {
		dirName := filepath.Base(dir.DestinationPathInHostAbs)
		dirsToExportPaths = append(dirsToExportPaths, dirName)
	}

	cacheContent, err := ei.FilterContentInCache(workDirPathInCache, filesToFoundInCache, dirsToExportPaths)
	if err != nil {
		return &ExporterError{
			ErrWrapped: err,
			Details:    "the exporter failed to filter the content in the cache",
		}
	}

	ei.td.Logger.Info(fmt.Sprintf("the work dir %s was filtered in the cache", workDirPathInCache))

	if len(cacheContent.Files) == 0 && len(cacheContent.Dirs) == 0 {
		return &ExporterError{
			Details: "the exporter failed to filter the content in the cache",
		}
	}

	// Copy the filtered files, and dirs from the .cache to the .export path.
	if err := utils.CopyPaths(cacheContent.Files, cacheContent.Dirs, options.TerraDaggerExportPath); err != nil {
		return &ExporterError{
			ErrWrapped: err,
			Details:    "the exporter failed to copy the filtered content from the cache to the export path",
		}
	}

	ei.td.Logger.Info(fmt.Sprintf("the work dir %s was copied from the cache to the export path", workDirPathInCache))

	// delete the old work dir in the cache.
	if err := utils.DeleteDirE(workDirPathInCache); err != nil {
		return &ExporterError{
			ErrWrapped: err,
			Details:    "the exporter failed to delete the old work dir in the cache",
		}
	}

	ei.td.Logger.Info(fmt.Sprintf("the work dir %s was deleted from the cache", workDirPathInCache))

	return nil
}

type CacheContent struct {
	Files []string
	Dirs  []string
}

func (ei *ExporterImpl) FilterContentInCache(cachePathAbs string, files, dirs []string) (*CacheContent, error) {
	foundPaths, err := utils.FilterPathsByNames(cachePathAbs, files, dirs)
	if err != nil {
		return nil, &ExporterError{
			ErrWrapped: err,
			Details:    "the exporter failed to filter the content in the cache",
		}
	}

	return &CacheContent{
		Files: foundPaths.FilesPathFound,
		Dirs:  foundPaths.DirsPathFound,
	}, nil
}
