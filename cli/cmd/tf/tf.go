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
			Workspace: viper.GetString("workspace"),
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
				ModulePath:                   viper.GetString("module"),
				EnableSSHPrivateGit:          true,
				TerraformVersion:             viper.GetString("tf-version"),
				EnvVarsToInjectByKeyFromHost: []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_SESSION_TOKEN"},
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
