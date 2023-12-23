package erroer

import "fmt"

type ErrTerraDaggerInvalidArgumentError struct {
	ErrWrapped error
	Details    string
}

func (e *ErrTerraDaggerInvalidArgumentError) Error() string {
	return fmt.Sprintf("The argument is either invalid, missing or it's not supported by TerraDagger: %s: %s",
		e.Details, e.ErrWrapped)
}

type ErrTerraDaggerConfigurationError struct {
	ErrWrapped error
	Details    string
}

func (e *ErrTerraDaggerConfigurationError) Error() string {
	return fmt.Sprintf("The TerraDagger configuration is invalid. The job failed to start: %s: %s",
		e.Details, e.ErrWrapped)
}

type ErrTerraDaggerInvalidMountPath struct {
	ErrWrapped error
	MountPath  string
}

func (e *ErrTerraDaggerInvalidMountPath) Error() string {
	return fmt.Sprintf("The mount path is invalid: %s: %s", e.MountPath, e.ErrWrapped)
}

type ErrTerraDaggerInitializationError struct {
	ErrWrapped error
	Details    string
}

func (e *ErrTerraDaggerInitializationError) Error() string {
	return fmt.Sprintf("The TerraDagger initialization failed: %s: %s", e.Details, e.ErrWrapped)
}

type ErrTerraDaggerContainerFailedToCreate struct {
	ErrWrapped error
	Image      string
}

func (e *ErrTerraDaggerContainerFailedToCreate) Error() string {
	return fmt.Sprintf("The TerraDagger container (dagger) with the image '%s' failed to create: %s", e.Image, e.ErrWrapped)
}
