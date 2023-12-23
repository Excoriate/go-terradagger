package terraform

import (
	"fmt"

	"github.com/Excoriate/go-terradagger/pkg/commands"
	"github.com/Excoriate/go-terradagger/pkg/erroer"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
)

type PlanOptions struct {
	// VarFiles is a list of terraform var files to use when running terraform plan
	VarFiles []string

	// PlanOutFilePath is the path to save the plan file
	PlanOutFilePath string

	// Vars is a map of terraform vars to use when running terraform plan
	Vars map[string]interface{}

	RefreshOnly bool
	Refresh     bool
}

func (o *PlanOptions) validateCMDOptions(options *Options) error {
	if o.VarFiles == nil {
		o.VarFiles = []string{}
	}

	if o.Vars == nil {
		o.Vars = map[string]interface{}{}
	}

	// TODO: Add support for output the plan file.
	// TOOD: Add option if the Vars are passed, a .hcl variables.tf should be parsed.

	// Check each of the *.tfvars passed.
	var varFilesNormalised []string
	for _, varFile := range o.VarFiles {
		varFilePath, err := getTerraformFileRelativePath(options.TerraformRootDir,
			options.TerraformDir, varFile)

		if err != nil {
			return &erroer.ErrTerraformPlanFailedToStart{
				ErrWrapped: err,
				Details:    fmt.Sprintf("the var file %s is invalid", varFile),
			}
		}

		// Add it to the list of var files.
		varFilesNormalised = append(varFilesNormalised, varFilePath)
	}

	o.VarFiles = varFilesNormalised

	return nil
}

func Plan(td *terradagger.Client, options *Options, planOptions *PlanOptions) error {
	setDefaultOptions(td, options)

	if planOptions == nil {
		planOptions = &PlanOptions{}
	}

	if err := options.validate(); err != nil {
		return &erroer.ErrTerraformPlanFailedToStart{
			ErrWrapped: err,
			Details:    "the options passed to the terraform command are invalid",
		}
	}

	if err := planOptions.validateCMDOptions(options); err != nil {
		return &erroer.ErrTerraformPlanFailedToStart{
			ErrWrapped: err,
			Details:    "the plan options passed to the terraform command are invalid",
		}
	}

	td.Logger.Info("All the options are valid, and the terraform plan command can be started.")

	tfInitCMD := initCMDDefault()

	// Setting the required terraform plan args, based on options.
	tfPlanArgs := &commands.CmdArgs{}
	if planOptions.PlanOutFilePath != "" {
		tfPlanArgs.AddNew(commands.CommandArgument{
			ArgName:  "-out",
			ArgValue: planOptions.PlanOutFilePath,
			ArgType:  commands.ArgTypeValue,
		})
	}

	if len(planOptions.VarFiles) > 0 {
		for _, varFile := range planOptions.VarFiles {
			tfPlanArgs.AddNew(commands.CommandArgument{
				ArgName:  "-var-file",
				ArgValue: varFile,
				ArgType:  commands.ArgTypeValue,
			})
		}
	}

	// In case of -vars, we need to convert the map to a slice of CommandArgument.
	if len(planOptions.Vars) > 0 {
		varsArgsConverted, err := convertInputVars(planOptions.Vars)
		if err != nil {
			return &erroer.ErrTerraformPlanFailedToStart{
				ErrWrapped: err,
				Details:    "the plan options passed to the terraform command are invalid",
			}
			// Handle the error accordingly
		}

		// Append each converted var argument
		for _, arg := range varsArgsConverted {
			tfPlanArgs.AddNew(arg)
		}
	}

	tfPlanCMD := commands.NewTerraDaggerCMD("terraform", "plan", tfPlanArgs.FormatArguments())
	tfPlanCMD.OmitBinaryNameInCommand = true

	// Add the necessary commands to the list of commands to run.
	cmds := []commands.TerraDaggerCMD{
		*tfInitCMD,
		*tfPlanCMD,
	}

	tfCMDDagger := commands.ConvertCommandsToDaggerFormat(cmds)

	tfImage := resolveTerraformImage(options)
	tfVersion := resolveTerraformVersion(options)

	// Configuring the options.
	tdOptions := &terradagger.ClientConfigOptions{
		Image:           tfImage,
		Version:         tfVersion,
		Workdir:         options.TerraformDir,
		MountDir:        td.MountDir,
		TerraDaggerCMDs: tfCMDDagger,
	}

	tdOptions.EnvVars = resolveEnvVarsByOptions(options)

	c, err := td.Configure(tdOptions)

	if err != nil {
		return &erroer.ErrTerraformPlanFailedToStart{
			ErrWrapped: err,
			Details:    "the terraform plan command failed to start",
		}
	}

	// Run the container.
	return td.Run(c)
}
