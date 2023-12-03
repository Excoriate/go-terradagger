package terraform

import "github.com/Excoriate/go-terradagger/pkg/terradagger"

type LifeCycleOptions struct {
	// InitOptions is the options for the terraform init command
	InitOptions InitOptions

	// PlanOptions is the options for the terraform plan command
	PlanOptions PlanOptions

	// ApplyOptions is the options for the terraform apply command
	ApplyOptions ApplyOptions

	// DestroyOptions is the options for the terraform destroy command
	DestroyOptions DestroyOptions
}

// CI runs the terraform lifecycle commands in CI mode
// TODO: Pending to implement this.
func CI(td *terradagger.Client, options *Options, lifecycleOptions *LifeCycleOptions) error {
	return nil
}
