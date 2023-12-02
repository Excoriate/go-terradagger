package daggerio

import "dagger.io/dagger"

type DaggerContainer struct {
	Client *dagger.Client
}

type Container interface {
	Create(image string) *dagger.Container
}

func NewContainer(client *dagger.Client) *DaggerContainer {
	return &DaggerContainer{
		Client: client,
	}
}

func (d *DaggerContainer) Create(image string) *dagger.Container {
	return d.Client.Container().From(image)
}
