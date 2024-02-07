package daggerx

import (
	"context"
	"os"

	"github.com/Excoriate/go-terradagger/pkg/logger"

	"dagger.io/dagger"
)

type Engine interface {
	Start(ctx context.Context, options ...dagger.ClientOpt) (*dagger.Client, error)
	GetEngine() *dagger.Client
}

type DaggerEngine struct {
	l logger.Log
	c *dagger.Client
}

func New(l logger.Log) Engine {
	return &DaggerEngine{
		l: l,
	}
}

func (b *DaggerEngine) Start(ctx context.Context, options ...dagger.ClientOpt) (*dagger.
	Client, error) {
	var c context.Context
	if ctx == nil {
		c = context.Background()
	} else {
		c = ctx
	}

	var daggerOptions []dagger.ClientOpt

	var daggerClient *dagger.Client

	if len(options) == 0 {
		daggerOptions = append(daggerOptions, dagger.WithLogOutput(os.Stdout))
	} else {
		daggerOptions = append(daggerOptions, options...)
	}

	daggerClient, err := dagger.Connect(c, daggerOptions...)
	if err != nil {
		return nil, err
	}

	b.c = daggerClient
	return daggerClient, nil
}

func (b *DaggerEngine) GetEngine() *dagger.Client {
	return b.c
}
