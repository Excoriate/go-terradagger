package terradagger

import "strings"

const (
	BashEntrypoint = "bash"
	ShEntrypoint   = "sh"
)

type BuildTerraformCommandOptions struct {
	Binary      string
	Command     string
	CommandArgs []string // Terraform-specific arguments
}

type BuildTerragruntCommandOptions struct {
	Binary         string   // The Terragrunt binary to use
	Subcommand     string   // e.g., "run-all"
	TerragruntArgs []string // Global options that apply to Terragrunt
	Command        string   // The Terraform command to run, e.g., "init", "apply"
	CommandArgs    []string // Arguments for the Terraform command
}

func BuildTerraformCommand(args BuildTerraformCommandOptions) string {
	var sb strings.Builder

	if args.Binary == "" {
		args.Binary = "terraform"
	}

	sb.WriteString(args.Binary)
	if args.Command != "" {
		sb.WriteString(" " + args.Command)
	}

	for _, arg := range args.CommandArgs {
		sb.WriteString(" " + arg)
	}

	return sb.String()
}

func BuildTerragruntCommand(args BuildTerragruntCommandOptions) string {
	var sb strings.Builder

	if args.Binary == "" {
		args.Binary = "terragrunt"
	}

	// Add "terragrunt" as the binary
	sb.WriteString("terragrunt")

	// If there's a subcommand (like "run-all"), add it
	if args.Subcommand != "" {
		sb.WriteString(" " + args.Subcommand)
	}

	// Append global Terragrunt arguments
	for _, arg := range args.TerragruntArgs {
		sb.WriteString(" " + arg)
	}

	// Add the Terraform command if present
	if args.Command != "" {
		sb.WriteString(" " + args.Command)
	}

	// Append the Terraform-specific arguments
	for _, arg := range args.CommandArgs {
		sb.WriteString(" " + arg)
	}

	return sb.String()
}

func BuildCMDWithBash(command string) []string {
	return []string{BashEntrypoint, "-c", command}
}

func BuildCMDWithSH(command string) []string {
	return []string{ShEntrypoint, "-c", command}
}
