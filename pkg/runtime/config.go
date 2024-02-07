package runtime

import "dagger.io/dagger"

type runtime struct {
	id        string
	container *dagger.Container
	dependsOn []string // IDs of other runtimes
}

type Options struct {
	Containers []*dagger.Container
}

type Config interface {
	GetID() string
	GetDependsOn(runtimeId string) []string
}
