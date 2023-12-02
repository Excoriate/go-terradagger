package terraform

import (
  "github.com/Excoriate/go-terradagger/pkg/commands"
  "github.com/Excoriate/go-terradagger/pkg/errors"
  "github.com/Excoriate/go-terradagger/pkg/terradagger"
  "github.com/Excoriate/go-terradagger/pkg/utils"
  "path/filepath"
)

type PlanOptions struct {
  // VarFiles is a list of terraform var files to use when running terraform plan
  VarFiles []string

  // PlanOutFilePath is the path to save the plan file
  PlanOutFilePath string

  // Vars is a map of terraform vars to use when running terraform plan
  Vars map[string]interface{}
}

func (o *PlanOptions) validateCMDOptions(terraformDir string) error {
  if o.VarFiles == nil {
    o.VarFiles = []string{}
  }

  if o.Vars == nil {
    o.Vars = map[string]interface{}{}
  }

  if o.PlanOutFilePath != "" {
    planOutFilePath := filepath.Join(terraformDir, o.PlanOutFilePath)

    if err := utils.FileExist(planOutFilePath); err != nil {
      return &errors.ErrTerraformPlanFilePathIsInvalid{
        ErrWrapped:   err,
        PlanFilePath: o.PlanOutFilePath,
        TerraformDir: terraformDir,
      }
    }

    o.PlanOutFilePath = planOutFilePath
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

func Plan(td *terradagger.Client, options *Options, planOptions *PlanOptions) error {
  if options == nil {
    options = &Options{}
  }

  if planOptions == nil {
    planOptions = &PlanOptions{}
  }

  if err := options.validate(); err != nil {
    return &errors.ErrTerraformPlanFailedToStart{
      ErrWrapped: err,
      Details:    "the options passed to the terraform command are invalid",
    }
  }

  if err := planOptions.validateCMDOptions(options.TerraformDir); err != nil {
    return &errors.ErrTerraformPlanFailedToStart{
      ErrWrapped: err,
      Details:    "the plan options passed to the terraform command are invalid",
    }
  }

  td.Logger.Info("All the options are valid, and the terraform plan command can be started.")

  var varsAsArgs []commands.Args
  if len(planOptions.Vars) > 0 {
    args, err := commands.ConvertMapIntoTerraformVarsOption(planOptions.Vars)
    if err != nil {
      return &errors.ErrTerraformPlanFailedToStart{
        ErrWrapped: err,
        Details:    "the plan options passed to the terraform command are invalid",
      }
    }

    varsAsArgs = args
  }

  tfInitCMD := commands.GetTerraformCommand("init", []commands.Args{
    {
      ArgName: "-backend=false",
    }, // This can be dynamically passed. An optional backend configuration can be used
  })
  tfPlanCMD := commands.GetTerraformCommand("plan", varsAsArgs)
  tfInitCMD.OmitBinaryNameInCommand = true
  tfPlanCMD.OmitBinaryNameInCommand = true

  // Convert to a terraDagger format, in this case, there are more than
  // one command to run.
  cmds := []commands.Command{
    tfInitCMD,
    tfPlanCMD,
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
    return &errors.ErrTerraformInitFailedToStart{
      ErrWrapped: err,
      Details:    "the terraform init command could not be configured",
    }
  }

  // Run the container.
  return td.Run(c)
}
