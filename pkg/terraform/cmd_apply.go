package terraform

import (
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/commands"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"

	"github.com/Excoriate/go-terradagger/pkg/errors"
	"github.com/Excoriate/go-terradagger/pkg/utils"
)

type ApplyOptions struct {
	// VarFiles is a list of terraform var files to use when running terraform apply
	VarFiles []string

	// Vars is a map of terraform vars to use when running terraform apply
	Vars map[string]interface{}
}

func (o *ApplyOptions) validateCMDOptions(terraformDir string) error {
	if o.VarFiles == nil {
		o.VarFiles = []string{}
	}

	if o.Vars == nil {
		o.Vars = map[string]interface{}{}
	}

	// Check each of the *.tfvars passed.
	var varFilesNormalised []string
	for _, varFile := range o.VarFiles {
		varFilePath := filepath.Join(terraformDir, varFile)

		if err := utils.FileExist(varFilePath); err != nil {
			return &errors.ErrTerraformVarFileIsInvalid{
				ErrWrapped:   err,
				VarFilePath:  varFile,
				TerraformDir: terraformDir,
			}
		}

		// If the file doesn't have the .json or tfvars extension, fail.
		if filepath.Ext(varFilePath) != ".json" && filepath.Ext(varFilePath) != ".tfvars" {
			return &errors.ErrTerraformVarFileIsInvalid{
				ErrWrapped:   nil,
				VarFilePath:  varFile,
				TerraformDir: terraformDir,
			}
		}

		// Add it to the list of var files.
		varFilesNormalised = append(varFilesNormalised, varFilePath)
	}

	o.VarFiles = varFilesNormalised

	return nil
}

func Apply(td *terradagger.Client, options *Options, applyOptions *ApplyOptions) error {
	if options == nil {
		options = &Options{}
	}

	if applyOptions == nil {
		applyOptions = &ApplyOptions{}
	}

	if err := options.validate(); err != nil {
		return &errors.ErrTerraformApplyFailedToStart{
			ErrWrapped: err,
			Details:    "the options passed to the terraform command are invalid",
		}
	}

	if err := applyOptions.validateCMDOptions(options.TerraformDir); err != nil {
		return &errors.ErrTerraformApplyFailedToStart{
			ErrWrapped: err,
			Details:    "the apply options passed to the terraform command are invalid",
		}
	}

	td.Logger.Info("All the options are valid, and the terraform apply command can be started.")

	// Setting the required terraform init args.
	tfInitArgs := &commands.CmdArgs{}
	tfInitArgs.AddNew(commands.CommandArgument{
		ArgName:  "-input",
		ArgValue: "false",
		ArgType:  commands.ArgTypeFlag,
	})

	tfInitArgs.AddNew(commands.CommandArgument{
		ArgName:  "-backend",
		ArgValue: "false",
		ArgType:  commands.ArgTypeFlag,
	})

	tfInitArgs.AddNew(commands.CommandArgument{
		ArgName:  "-upgrade",
		ArgValue: "false",
		ArgType:  commands.ArgTypeFlag,
	})

	tfInitCMD := commands.NewTerraDaggerCMD("terraform", "init", tfInitArgs.FormatArguments())
	tfInitCMD.OmitBinaryNameInCommand = true

	// Setting the required terraform apply args, based on options.
	tfApplyArgs := &commands.CmdArgs{}

	if len(applyOptions.VarFiles) > 0 {
		for _, varFile := range applyOptions.VarFiles {
			tfApplyArgs.AddNew(commands.CommandArgument{
				ArgName:  "-var-file",
				ArgValue: varFile,
				ArgType:  commands.ArgTypeValue,
			})
		}
	}

	// In case of -vars, we need to convert the map to a slice of CommandArgument.
	if len(applyOptions.Vars) > 0 {
		varsArgsConverted, err := convertInputVars(applyOptions.Vars)
		if err != nil {
			return &errors.ErrTerraformApplyFailedToStart{
				ErrWrapped: err,
				Details:    "the apply options passed to the terraform command are invalid",
			}
			// Handle the error accordingly
		}

		// Append each converted var argument
		for _, arg := range varsArgsConverted {
			tfApplyArgs.AddNew(arg)
		}
	}

	// Add the -auto-approve flag to the apply command as a default
	// TODO: I'm not sure about this. Perhaps it should be a flag in the options?
	tfApplyArgs.AddNew(commands.CommandArgument{
		ArgName: "-auto-approve",
		ArgType: commands.ArgTypeFlag,
	})

	tfApplyCMD := commands.NewTerraDaggerCMD("terraform", "apply", tfApplyArgs.FormatArguments())
	tfApplyCMD.OmitBinaryNameInCommand = true

	// Add the necessary commands to the list of commands to run.
	cmds := []commands.TerraDaggerCMD{
		*tfInitCMD,
		*tfApplyCMD,
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
		return &errors.ErrTerraformApplyFailedToStart{
			ErrWrapped: err,
			Details:    "the terraform apply command failed to start",
		}
	}

	// Run the container.
	return td.Run(c)
}
