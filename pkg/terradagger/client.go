package terradagger

import (
	"context"

	"github.com/Excoriate/go-terradagger/pkg/utils"

	"github.com/Excoriate/go-terradagger/pkg/config"

	"github.com/Excoriate/go-terradagger/pkg/daggerx"

	"github.com/Excoriate/go-terradagger/pkg/logger"
)

type Options struct {
	Workspace     string
	EnvVars       map[string]string
	ExcludeDirs   []string
	ExcludedFiles []string
}

type Client interface {
	StartEngine() error
}

type TD struct {
	Ctx    context.Context
	Log    logger.Log
	Engine daggerx.Engine
	Config config.Config
	ID     string
}

func New(ctx context.Context, options *Options) *TD {
	td := &TD{
		Log: logger.NewLogger().Logger,
		Ctx: ctx,
		ID:  utils.GetUUID(),
	}

	td.Engine = daggerx.New(td.Log)
	td.Config = config.New(options.Workspace, options.EnvVars, options.ExcludeDirs, options.ExcludedFiles)

	return td
}

// StartEngine starts the Dagger engine.
func (td *TD) StartEngine() error {
	_, err := td.Engine.Start(td.Ctx)
	if err != nil {
		return err
	}

	return nil
}
