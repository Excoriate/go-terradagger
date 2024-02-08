package tf

import (
	"context"

	"github.com/Excoriate/go-terradagger/pkg/terragrunt"

	"github.com/Excoriate/go-terradagger/cli/internal/tui"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
	"github.com/Excoriate/go-terradagger/pkg/terraformcore"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	terraformDir        string
	terraformVars       map[string]string
	terraformTFVarFiles []string
	terraformVersion    string
	all                 bool
)

var Cmd = &cobra.Command{
	Use:   "tf",
	Short: "Execute terraform CI Jobs using Dagger",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		// Build the UX.
		ux := &struct {
			Msg   tui.MessageWriter
			Title tui.TitleWriter
		}{
			Msg:   tui.NewMessageWriter(),
			Title: tui.NewTitleWriter(),
		}

		ux.Title.ShowTitle("TerraDagger CLI")

		td := terradagger.New(ctx, &terradagger.Options{
			Workspace: "../",
		})

		// Start the engine (and the Dagger backend)
		if err := td.StartEngine(); err != nil {
			ux.Msg.ShowError(tui.MessageOptions{
				Message: "Unable to start the Dagger engine",
				Error:   err,
			})
		}

		defer td.Engine.GetEngine().Close()

		// -------------------------------
		// terraform
		// -------------------------------
		//terraformOptions := terraformcore.WithOptions(td, &terraformcore.TfOptions{
		//  ModulePath: "test/terraform/root-module-1",
		//})

		//_, initErr := terraform.InitE(td, terraformOptions, terraform.InitOptions{})
		//// Run terraform init
		//if initErr != nil {
		//	ux.Msg.ShowError(tui.MessageOptions{
		//		Message: "Error initializing terraform",
		//		Error:   initErr,
		//	})
		//}

		// -------------------------------
		// Terragrunt
		// -------------------------------
		terragruntOptions :=
			terraformcore.WithOptions(td, &terraformcore.TfOptions{
				ModulePath:                   "test/terraform/terragrunt-2",
				EnableSSHPrivateGit:          true,
				EnvVarsToInjectByKeyFromHost: []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"},
			})

		_, initTgErr := terragrunt.InitE(td, terragruntOptions, terragrunt.InitOptions{}, terragrunt.GlobalOptions{})
		if initTgErr != nil {
			ux.Msg.ShowError(tui.MessageOptions{
				Message: "Error initializing terragrunt",
				Error:   initTgErr,
			})
		}

		_, planTgErr := terragrunt.PlanE(td, terragruntOptions, terragrunt.PlanOptions{}, terragrunt.GlobalOptions{})
		if planTgErr != nil {
			ux.Msg.ShowError(tui.MessageOptions{
				Message: "Error planning terragrunt",
				Error:   planTgErr,
			})
		}
	},
}

func AddFlags() {
	Cmd.PersistentFlags().BoolVarP(&all, "all", "", false, "Execute all recipes in the 'examples' folder.")
	Cmd.PersistentFlags().StringVarP(&terraformDir, "terraform-dir", "", "",
		"The directory where the terraform code resides. "+
			"It is also the directory that'll be mounted into Dagger's container.")
	Cmd.PersistentFlags().StringToStringVar(&terraformVars, "terraform-vars", map[string]string{},
		"Variables to pass to terraform.")
	Cmd.PersistentFlags().StringSliceVar(&terraformTFVarFiles, "terraform-tfvar-files", []string{},
		"TFVar files to pass to terraform.")
	Cmd.PersistentFlags().StringVarP(&terraformVersion, "terraform-version", "", "",
		"The version of terraform to use.")

	_ = viper.BindPFlag("all", Cmd.PersistentFlags().Lookup("all"))
	_ = viper.BindPFlag("terraform-dir", Cmd.PersistentFlags().Lookup("terraform-dir"))
	_ = viper.BindPFlag("terraform-vars", Cmd.PersistentFlags().Lookup("terraform-vars"))
	_ = viper.BindPFlag("terraform-tfvar-files", Cmd.PersistentFlags().Lookup("terraform-tfvar-files"))
	_ = viper.BindPFlag("terraform-version", Cmd.PersistentFlags().Lookup("terraform-version"))
}

func init() {
	AddFlags()
}
