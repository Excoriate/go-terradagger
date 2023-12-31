package terradagger

import "fmt"

type interopClient struct {
	td *TD
}

type interop interface {
	FilesToCopyFromHostAreValid(o *FilesToCopyFromHostAreValidOptions) error
	DirsToCopyFromHostAreValid(o *DirsToCopyFromHostAreValidOptions) error
}

type FilesToCopyFromHostAreValidOptions struct {
	Files             []string
	BaseExportPathAbs string
}

type DirsToCopyFromHostAreValidOptions struct {
	Dirs           []string
	BaseExportPath string
}

func newInteropClient(td *TD) interop {
	return &interopClient{
		td: td,
	}
}

func (c *interopClient) FilesToCopyFromHostAreValid(o *FilesToCopyFromHostAreValidOptions) error {
	if o == nil {
		return fmt.Errorf("failed to validate the files to copy from host, the options are nil")
	}

	if o.Files == nil {
		return fmt.Errorf("failed to validate the files to copy from host, the files are nil")
	}

	if o.BaseExportPathAbs == "" {
		return fmt.Errorf("failed to validate the files to copy from host, the base export path cannot be empty")
	}

	if len(o.Files) > 0 {
		for _, file := range o.Files {
			if file == "" {
				return fmt.Errorf("failed to validate the files to copy from host, the file cannot be empty")
			}
		}
	}

	return nil
}

func (c *interopClient) DirsToCopyFromHostAreValid(o *DirsToCopyFromHostAreValidOptions) error {
	return nil
}
