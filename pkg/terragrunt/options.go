package terragrunt

type GlobalOptions struct {
	// The path to the Terragrunt config file. Default is terragrunt.hcl.
	Config string `json:"terragrunt-config,omitempty"`

	// Write terragrunt-debug.tfvars to working folder to help root-cause issues.
	Debug bool `json:"terragrunt-debug,omitempty"`

	// When this flag is set Terragrunt will not update the remote state bucket.
	DisableBucketUpdate bool `json:"terragrunt-disable-bucket-update,omitempty"`

	// When this flag is set, Terragrunt will not validate the terraform command.
	DisableCommandValidation bool `json:"terragrunt-disable-command-validation,omitempty"`

	// The path where to download Terraform code. Default is .terragrunt-cache in the working directory.
	DownloadDir string `json:"terragrunt-download-dir,omitempty"`

	// Unix-style glob of directories to exclude when running *-all commands.
	ExcludeDir string `json:"terragrunt-exclude-dir,omitempty"`

	// When this flag is set Terragrunt will fail if the remote state bucket needs to be created.
	FailOnStateBucketCreation bool `json:"terragrunt-fail-on-state-bucket-creation,omitempty"`

	// The option fetchs dependency output directly from the state file instead of init dependencies and running terraform on them.
	FetchDependencyOutputFromState bool `json:"terragrunt-fetch-dependency-output-from-state,omitempty"`

	// Session duration for IAM Assume Role session.
	IamAssumeRoleDuration string `json:"terragrunt-iam-assume-role-duration,omitempty"`

	// Name for the IAM Assummed Role session.
	IamAssumeRoleSessionName string `json:"terragrunt-iam-assume-role-session-name,omitempty"`

	// Assume the specified IAM role before executing Terraform.
	IamRole string `json:"terragrunt-iam-role,omitempty"`

	// *-all commands continue processing components even if a dependency fails.
	IgnoreDependencyErrors bool `json:"terragrunt-ignore-dependency-errors,omitempty"`

	// *-all commands will be run disregarding the dependencies.
	IgnoreDependencyOrder bool `json:"terragrunt-ignore-dependency-order,omitempty"`

	// *-all commands will not attempt to include external dependencies.
	IgnoreExternalDependencies bool `json:"terragrunt-ignore-external-dependencies,omitempty"`

	// Unix-style glob of directories to include when running *-all commands.
	IncludeDir string `json:"terragrunt-include-dir,omitempty"`

	// *-all commands will include external dependencies.
	IncludeExternalDependencies bool `json:"terragrunt-include-external-dependencies,omitempty"`

	// When this flag is set output from Terraform sub-commands is prefixed with module path.
	IncludeModulePrefix bool `json:"terragrunt-include-module-prefix,omitempty"`

	// If specified, Terragrunt will output its logs in JSON format.
	JSONLog bool `json:"terragrunt-json-log,omitempty"`

	// Sets the logging level for Terragrunt. Supported levels: panic, fatal, error, warn, info, debug, trace.
	LogLevel string `json:"terragrunt-log-level,omitempty"`

	// If flag is set, 'run-all' will only run the command against Terragrunt modules that include the specified file.
	ModulesThatInclude string `json:"terragrunt-modules-that-include,omitempty"`

	// Don't automatically append -auto-approve to the underlying Terraform commands run with 'run-all'.
	NoAutoApprove bool `json:"terragrunt-no-auto-approve,omitempty"`

	// Don't automatically run 'terraform init' during other terragrunt commands. You must run 'terragrunt init' manually.
	NoAutoInit bool `json:"terragrunt-no-auto-init,omitempty"`

	// Don't automatically re-run command in case of transient errors.
	NoAutoRetry bool `json:"terragrunt-no-auto-retry,omitempty"`

	// If specified, Terragrunt output won't contain any color.
	NoColor bool `json:"terragrunt-no-color,omitempty"`

	// Assume "yes" for all prompts.
	NonInteractive bool `json:"terragrunt-non-interactive,omitempty"`

	// *-all commands parallelism set to at most N modules.
	Parallelism int `json:"terragrunt-parallelism,omitempty"`

	// Download Terraform configurations from the specified source into a temporary folder, and run Terraform in that temporary folder.
	Source string `json:"terragrunt-source,omitempty"`

	// Replace any source URL (including the source URL of a config pulled in with dependency blocks) that has root source with dest.
	SourceMap string `json:"terragrunt-source-map,omitempty"`

	// Delete the contents of the temporary folder to clear out any old, cached source code before downloading new source code into it.
	SourceUpdate bool `json:"terragrunt-source-update,omitempty"`

	// If flag is set, only modules under the directories passed in with '--terragrunt-include-dir' will be included.
	StrictInclude bool `json:"terragrunt-strict-include,omitempty"`

	// If specified, Terragrunt will wrap Terraform stdout and stderr in JSON.
	TfLogsToJSON bool `json:"terragrunt-tf-logs-to-json,omitempty"`

	// Path to the Terraform binary. Default is terraform (on PATH).
	TfPath string `json:"terragrunt-tfpath,omitempty"`

	// Enables caching of includes during partial parsing operations. Will also be used for the --terragrunt-iam-role option if provided.
	UsePartialParseConfigCache bool `json:"terragrunt-use-partial-parse-config-cache,omitempty"`
}
