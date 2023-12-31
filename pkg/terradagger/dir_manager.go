package terradagger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/utils"
)

type dirManagerClient struct {
	td *TD
}

type dirManager interface {
	// CreateTerraDaggerDir creates a new TerraDagger directory, based on the instance
	// of the TerraDagger execution that's running. It's a per-instance.
	CreateTerraDaggerDir(options *CreateTerraDaggerDirOptions) (string, error)
	CreateAuxDirsInTDDir(options []CreateAuxDirsInTDDirOptions) error
}

func newDirManagerClient(td *TD) dirManager {
	return &dirManagerClient{
		td: td,
	}
}

type CreateTerraDaggerDirOptions struct {
	TerraDaggerPathResolved string
	SkipCreationIfExist     bool
}

func (c *dirManagerClient) CreateTerraDaggerDir(options *CreateTerraDaggerDirOptions) (string, error) {
	if options == nil {
		return "", fmt.Errorf("failed to createNewContainer the TerraDagger dir, the options are nil")
	}

	dirUtils := utils.DirUtils{}

	if options.SkipCreationIfExist && dirUtils.DirExist(options.TerraDaggerPathResolved) {
		c.td.Logger.Info(fmt.Sprintf("the TerraDagger dir %s already exists, skipping the creation, as requested", options.TerraDaggerPathResolved))
		return options.TerraDaggerPathResolved, nil
	}

	if err := os.MkdirAll(options.TerraDaggerPathResolved, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to createNewContainer the TerraDagger dir: %w", err)
	}

	return options.TerraDaggerPathResolved, nil
}

type CreateAuxDirsInTDDirOptions struct {
	TerraDaggerPathResolved string
	SkipCreationIfExist     bool
	AuxDirName              string
}

func (c *dirManagerClient) CreateAuxDirsInTDDir(options []CreateAuxDirsInTDDirOptions) error {
	if options == nil {
		return fmt.Errorf("failed to createNewContainer the aux dirs in the TerraDagger dir, the options are nil")
	}

	if len(options) == 0 {
		return fmt.Errorf("failed to createNewContainer the aux dirs in the TerraDagger dir, the options are empty")
	}

	dirUtils := utils.DirUtils{}

	for _, option := range options {
		auxDirPath := filepath.Join(option.TerraDaggerPathResolved, option.AuxDirName)

		if option.SkipCreationIfExist && dirUtils.DirExist(auxDirPath) {
			c.td.Logger.Info(fmt.Sprintf("the aux dir %s already exists, skipping the creation", auxDirPath))
			continue
		}

		if err := os.MkdirAll(auxDirPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to createNewContainer the aux dir %s: %w", auxDirPath, err)
		}
	}

	return nil
}
