package terradagger

import (
	"fmt"

	"github.com/Excoriate/go-terradagger/pkg/utils"
)

type BackupFilesError struct {
	ErrWrapped error
	Details    string
}

const BackupFilesErrPrefix = "the backer failed when tried to backup the files"

func (e *BackupFilesError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", BackupFilesErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", BackupFilesErrPrefix, e.Details, e.ErrWrapped.Error())
}

type Backer interface {
	BackupFiles(options *BackupContentOptions) error
	BackupDirs(options *BackupContentOptions) error
	BackupManaged(options *BackupOptions) error
}

type BackerImpl struct {
	td *TD
}

func NewBacker(td *TD) Backer {
	return &BackerImpl{
		td: td,
	}
}

type BackupContentOptions struct {
	SourcePathAbs      string
	DestinationPathAbs string
}

func (b *BackerImpl) BackupFiles(options *BackupContentOptions) error {
	return utils.CopyFile(options.SourcePathAbs, options.DestinationPathAbs)
}

type BackupDirsOptions struct {
	Dirs []string
}

func (b *BackerImpl) BackupDirs(options *BackupContentOptions) error {
	return utils.CopyDir(options.SourcePathAbs, options.DestinationPathAbs)
}

type BackupOptions struct {
	Files []*DataBackupContent
	Dirs  []*DataBackupContent
}

func (b *BackerImpl) BackupManaged(options *BackupOptions) error {
	if options == nil {
		return &BackupFilesError{
			Details: "the backup options are nil",
		}
	}

	files := options.Files
	for _, file := range files {
		if err := b.BackupFiles(&BackupContentOptions{
			SourcePathAbs:      file.SourcePathAbs,
			DestinationPathAbs: file.DestinationPathAbs,
		}); err != nil {
			return &BackupFilesError{
				ErrWrapped: err,
				Details:    fmt.Sprintf("failed to backup the file %s", file.SourcePathAbs),
			}
		}
	}

	dirs := options.Dirs
	for _, dir := range dirs {
		if err := b.BackupDirs(&BackupContentOptions{
			SourcePathAbs:      dir.SourcePathAbs,
			DestinationPathAbs: dir.DestinationPathAbs,
		}); err != nil {
			return &BackupFilesError{
				ErrWrapped: err,
				Details:    fmt.Sprintf("failed to backup the dir %s", dir.SourcePathAbs),
			}
		}
	}

	return nil
}
