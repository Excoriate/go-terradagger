package terradagger

import "dagger.io/dagger"

type ContainerCommand []string

type Runner interface {
	Exec(cmds []ContainerCommand) *dagger.Container
}
