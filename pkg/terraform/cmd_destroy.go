package terraform

import (
	"path/filepath"

	"github.com/Excoriate/go-terradagger/pkg/commands"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"

	"github.com/Excoriate/go-terradagger/pkg/erroer"
	"github.com/Excoriate/go-terradagger/pkg/utils"
)

type DestroyOptions struct {
	// VarFiles is a list of terraform var files to use when running terraform destroy
	VarFiles []string

	// Vars is a map of terraform vars to use when running terraform destroy
	Vars map[string]interface{}
}

func (o *DestroyOptions) validateCMDOptions(terraformDir string) error {
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

		if err := utils.FileExistE(varFilePath); err != nil {
			return &erroer.ErrTerraformVarFileIsInvalid{
				ErrWrapped:   err,
				VarFilePath:  varFile,
				TerraformDir: terraformDir,
			}
		}

		// If the file doesn't have the .json or tfvars extension, fail.
		if filepath.Ext(varFilePath) != ".json" && filepath.Ext(varFilePath) != ".tfvars" {
			return &erroer.ErrTerraformVarFileIsInvalid{
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

func Destroy(td *terradagger.TD, options *Options, destroyOptions *DestroyOptions) error {
	setDefaultOptions(td, options)

	if destroyOptions == nil {
		destroyOptions = &DestroyOptions{}
	}

	if err := options.validate(); err != nil {
		return &erroer.ErrTerraformDestroyFailedToStart{
			ErrWrapped: err,
			Details:    "the options passed to the terraform destroy command are invalid",
		}
	}

	if err := destroyOptions.validateCMDOptions(options.TerraformModulePath); err != nil {
		return &erroer.ErrTerraformDestroyFailedToStart{
			ErrWrapped: err,
			Details:    "the destroy options passed to the terraform destroy command are invalid",
		}
	}

	td.Logger.Info("All the options are valid, and the terraform destroy command can be started.")

	// Setting the required terraform init args.
	tfInitCMD := initCMDDefault()

	// Setting the required terraform destroy args, based on options.
	tfDestroyArgs := &commands.CmdArgs{}

	if len(destroyOptions.VarFiles) > 0 {
		for _, varFile := range destroyOptions.VarFiles {
			tfDestroyArgs.AddNew(commands.CommandArgument{
				ArgName:  "-var-file",
				ArgValue: varFile,
				ArgType:  commands.ArgTypeValue,
			})
		}
	}

	// In case of -vars, we need to convert the map to a slice of CommandArgument.
	if len(destroyOptions.Vars) > 0 {
		varsArgsConverted, err := convertInputVars(destroyOptions.Vars)
		if err != nil {
			return &erroer.ErrTerraformDestroyFailedToStart{
				ErrWrapped: err,
				Details:    "the destroy options passed to the terraform destroy command are invalid",
			}
			// Handle the error accordingly
		}

		// Append each converted var argument
		for _, arg := range varsArgsConverted {
			tfDestroyArgs.AddNew(arg)
		}
	}

	// Add the -auto-approve flag to the destroy command as a default
	// TODO: I'm not sure about this. Perhaps it should be a flag in the options?
	tfDestroyArgs.AddNew(commands.CommandArgument{
		ArgName: "-auto-approve",
		ArgType: commands.ArgTypeFlag,
	})

	tfDestroyCDM := commands.NewTerraDaggerCMD("terraform", "destroy", tfDestroyArgs.FormatArguments())
	tfDestroyCDM.OmitBinaryNameInCommand = true

	// Add the necessary commands to the list of commands to run.
	cmds := []commands.TerraDaggerCMD{
		*tfInitCMD,
		*tfDestroyCDM,
	}

	tfCMDDagger := commands.ConvertCommandsToDaggerFormat(cmds)

	tfImage := resolveTerraformImage(options)
	tfVersion := resolveTerraformVersion(options)

	// Configuring the options.
	tdOptions := &terradagger.ClientOptions{
		ContainerOptions: &terradagger.InstanceContainerOptions{
			Image:   tfImage,
			Version: tfVersion,
		},
		WorkDirPath:     options.TerraformModulePath,
		TerraDaggerCMDs: tfCMDDagger,
	}

	i := terradagger.NewInstance(td)
	err := i.Validate(tdOptions)

	if err != nil {
		return &erroer.ErrTerraformInitFailedToStart{
			ErrWrapped: err,
			Details:    "the terraform destroy command could not be validated",
		}
	}

	cfg, err := i.Configure(tdOptions)
	if err != nil {
		return err
	}

	clientInstance, err := i.PrepareInstance(cfg)
	if err != nil {
		return err
	}

	if err := td.Run(clientInstance, nil); err != nil {
		return &erroer.ErrTerraformInitFailedToStart{
			ErrWrapped: err,
			Details:    "the terraform destroy command failed to run",
		}
	}

	return nil
}
