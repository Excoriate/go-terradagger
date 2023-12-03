package errors

import "fmt"

type ErrTerraformBackendFileIsNotFound struct {
	ErrWrapped      error
	BackendFilePath string
}

func (e *ErrTerraformBackendFileIsNotFound) Error() string {
	return fmt.Sprintf("The backend file %s is invalid, or it does not exist: %s",
		e.BackendFilePath, e.ErrWrapped)
}

type ErrTerraformInitFailedToStart struct {
	ErrWrapped error
	Details    string
}

func (e *ErrTerraformInitFailedToStart) Error() string {
	return fmt.Sprintf("Failed to start the terraform init command: %s: %s",
		e.Details, e.ErrWrapped)
}

type ErrTerraformOptionsAreInvalid struct {
	ErrWrapped error
	Details    string
}

func (e *ErrTerraformOptionsAreInvalid) Error() string {
	return fmt.Sprintf("The terraform options are invalid: %s: %s",
		e.Details, e.ErrWrapped)
}

type ErrTerraformPlanFailedToStart struct {
	ErrWrapped error
	Details    string
}

func (e *ErrTerraformPlanFailedToStart) Error() string {
	return fmt.Sprintf("Failed to start the terraform plan command: %s: %s",
		e.Details, e.ErrWrapped)
}

type ErrTerraformPlanFilePathIsInvalid struct {
	ErrWrapped   error
	PlanFilePath string
	TerraformDir string
}

func (e *ErrTerraformPlanFilePathIsInvalid) Error() string {
	return fmt.Sprintf("The plan file path %s is invalid, or it does not exist in the terraform dir %s: %s",
		e.PlanFilePath, e.TerraformDir, e.ErrWrapped)
}

type ErrTerraformVarFileIsInvalid struct {
	ErrWrapped   error
	VarFilePath  string
	TerraformDir string
}

func (e *ErrTerraformVarFileIsInvalid) Error() string {
	return fmt.Sprintf("The var file path %s is invalid, or it does not exist in the terraform dir %s: %s",
		e.VarFilePath, e.TerraformDir, e.ErrWrapped)
}

type ErrTerraformApplyFailedToStart struct {
	ErrWrapped error
	Details    string
}

func (e *ErrTerraformApplyFailedToStart) Error() string {
	return fmt.Sprintf("Failed to start the terraform apply command: %s: %s",
		e.Details, e.ErrWrapped)
}

type ErrTerraformDestroyFailedToStart struct {
	ErrWrapped error
	Details    string
}

func (e *ErrTerraformDestroyFailedToStart) Error() string {
	return fmt.Sprintf("Failed to start the terraform destroy command: %s: %s",
		e.Details, e.ErrWrapped)
}
