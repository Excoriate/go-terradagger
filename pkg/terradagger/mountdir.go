package terradagger

import (
	"fmt"

	"github.com/Excoriate/go-terradagger/pkg/utils"
)

type MountDirValidator interface {
	IsMountDirValid(mountPath string) error
	IsWorkDirValid(mountPath, workDirPath string) error
}

type MountDirValidatorImpl struct {
	td *TD
}

type MountDirError struct {
	ErrWrapped error
	Details    string
}

const mountDirErrPrefix = "the mount dir passed to the terradagger client instance is invalid"

func (e *MountDirError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", mountDirErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", mountDirErrPrefix, e.Details, e.ErrWrapped.Error())
}

func NewMountDirValidator(td *TD) MountDirValidator {
	return &MountDirValidatorImpl{
		td: td,
	}
}

func (v *MountDirValidatorImpl) IsMountDirValid(mountPath string) error {
	if mountPath == "" {
		return &MountDirError{
			Details: "the mount dir path is empty",
		}
	}

	du := utils.NewDirUtils()
	if err := du.IsValidDirE(mountPath); err != nil {
		return &MountDirError{
			Details: "the mount dir path is invalid",
		}
	}

	return nil
}

func (v *MountDirValidatorImpl) IsWorkDirValid(mountPath, workDirPath string) error {
	if err := v.IsMountDirValid(mountPath); err != nil {
		return err
	}

	if workDirPath == "" {
		return &MountDirError{
			Details: "the work dir path is empty",
		}
	}

	if mountPath == workDirPath {
		return &MountDirError{
			Details: "the work dir path cannot be the same as the mount dir path",
		}
	}

	workDirFullPath := buildWorkDirPath(mountPath, workDirPath)
	du := utils.NewDirUtils()

	if err := du.IsValidDirE(workDirFullPath); err != nil {
		return &MountDirError{
			Details: "the work dir path is invalid",
		}
	}

	return nil
}
