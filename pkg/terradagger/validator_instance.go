package terradagger

import (
	"fmt"

	"github.com/Excoriate/go-terradagger/pkg/env"
)

type clientInstanceValidator interface {
	IsEnvVarOptionsValid(options *EnvVarOptions) error
	IsContainerOptionsValid(options *InstanceContainerOptions) error
}

type clientInstanceValidatorImpl struct {
	clientInstance *InstanceImpl
}

func newClientInstanceValidator(clientInstance *InstanceImpl) clientInstanceValidator {
	return &clientInstanceValidatorImpl{
		clientInstance: clientInstance,
	}
}

type EnvVarOptionsError struct {
	ErrWrapped error
	Details    string
}

const envVarOptionsErrPrefix = "the environment variable options passed to the terradagger client" +
	" instance are invalid"

func (e *EnvVarOptionsError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", envVarOptionsErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", envVarOptionsErrPrefix, e.Details, e.ErrWrapped.Error())
}

func (v *clientInstanceValidatorImpl) IsEnvVarOptionsValid(options *EnvVarOptions) error {
	if options == nil {
		return nil
	}

	if len(options.EnvVars) == 0 {
		return &EnvVarOptionsError{
			Details: "the environment variables map is empty",
		}
	}

	if len(options.CopyEnvVarsFromHostByKeys) > 0 && options.MirrorEnvVarsFromHost {
		return &EnvVarOptionsError{
			Details: "the environment variables cannot be copied from the host and mirrored at the same time",
		}
	}

	if len(options.CopyEnvVarsFromHostByKeys) > 0 {
		for _, key := range options.CopyEnvVarsFromHostByKeys {
			if _, err := env.GetEnvVarByKey(key, true); err != nil {
				return &EnvVarOptionsError{
					ErrWrapped: err,
					Details:    fmt.Sprintf("the environment variable %s cannot be copied from the host", key),
				}
			}
		}
	}

	return nil
}

type ContainerOptionsError struct {
	ErrWrapped error
	Details    string
}

const containerOptionsErrPrefix = "the container options passed to the terradagger client" +
	" instance are invalid"

func (e *ContainerOptionsError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", containerOptionsErrPrefix, e.Details)
	}

	return fmt.Sprintf("%s: %s: %s", containerOptionsErrPrefix, e.Details, e.ErrWrapped.Error())
}

func (v *clientInstanceValidatorImpl) IsContainerOptionsValid(options *InstanceContainerOptions) error {
	containerValidator := newContainerValidator(v.clientInstance.td.Logger)
	if err := containerValidator.validate(&CreateNewContainerOptions{
		Image:   options.Image,
		Version: options.Version,
	}); err != nil {
		return &ContainerOptionsError{
			ErrWrapped: err,
			Details:    "the container options are invalid",
		}
	}

	return nil
}
