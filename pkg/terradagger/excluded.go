package terradagger

import (
	"fmt"
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/utils"
)

type ExcludedFiles interface {
	IsExcludeFileValid(mountPath, fileToExclude string) error
}

type ExcludedDirs interface {
	IsExcludeDirValid(mountPath, dirToExclude string) error
}

type ExcludedFilesImpl struct {
	td *TD
}

type ExcludedDirsImpl struct {
	td *TD
}

type ExcludedFilesError struct {
	ErrWrapped error
	Details    string
}

type ExcludedDirsError struct {
	ErrWrapped error
	Details    string
}

const excludedFilesErrPrefix = "the excluded files are invalid"
const excludedDirsErrPrefix = "the excluded dirs are invalid"

func (e *ExcludedFilesError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", excludedFilesErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", excludedFilesErrPrefix, e.Details, e.ErrWrapped.Error())
}

func (e *ExcludedDirsError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", excludedDirsErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", excludedDirsErrPrefix, e.Details, e.ErrWrapped.Error())
}

func NewExcludedFiles(td *TD) ExcludedFiles {
	return &ExcludedFilesImpl{
		td: td,
	}
}

func NewExcludedDirs(td *TD) ExcludedDirs {
	return &ExcludedDirsImpl{
		td: td,
	}
}

func (v *ExcludedFilesImpl) IsExcludeFileValid(mountPath, fileToExclude string) error {
	if mountPath == "" {
		return &ExcludedFilesError{
			Details: "the mount dir path is empty",
		}
	}

	if fileToExclude == "" {
		return &ExcludedFilesError{
			Details: "the file to exclude path is empty",
		}
	}

	filePath := filepath.Join(mountPath, fileToExclude)
	if err := utils.FileExistE(filePath); err != nil {
		return &ExcludedFilesError{
			Details: fmt.Sprintf("the file %s to exclude does not exist", filePath),
		}
	}

	if err := utils.IsAFileE(filePath); err != nil {
		return &ExcludedFilesError{
			Details: fmt.Sprintf("the file %s to exclude is not a file", filePath),
		}
	}

	return nil
}

func (v *ExcludedDirsImpl) IsExcludeDirValid(mountPath, dirToExclude string) error {
	if mountPath == "" {
		return &ExcludedDirsError{
			Details: "the mount dir path is empty",
		}
	}

	if dirToExclude == "" {
		return &ExcludedDirsError{
			Details: "the dir to exclude path is empty",
		}
	}

	dirPath := filepath.Join(mountPath, dirToExclude)
	dirUtils := utils.DirUtils{}
	if err := dirUtils.IsValidDirE(dirPath); err != nil {
		return &ExcludedDirsError{
			Details: fmt.Sprintf("the dir %s to exclude is not a valid dir", dirPath),
		}
	}

	return nil
}
