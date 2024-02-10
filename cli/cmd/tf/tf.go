package tf

import (
	"context"

	"github.com/Excoriate/go-terradagger/pkg/terraform"

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

		tfOpptions :=
			terraformcore.WithOptions(td, &terraformcore.TfOptions{
				ModulePath:                   viper.GetString("module"),
				EnableSSHPrivateGit:          true,
				TerraformVersion:             viper.GetString("tf-version"),
				EnvVarsToInjectByKeyFromHost: []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_SESSION_TOKEN"},
			})

		_, tfInitErr := terraform.InitE(td, tfOpptions, terraform.InitOptions{})
		if tfInitErr != nil {
			ux.Msg.ShowError(tui.MessageOptions{
				Message: tfInitErr.Error(),
				Error:   tfInitErr,
			})
		}

		_, tfPlanErr := terraform.PlanE(td, tfOpptions, terraform.PlanOptions{
			Vars: []terraformcore.TFInputVariable{
				{
					Name:  "is_enabled",
					Value: "true",
				}},
		})
		if tfPlanErr != nil {
			ux.Msg.ShowError(tui.MessageOptions{
				Message: tfInitErr.Error(),
				Error:   tfPlanErr,
			})
		}

		_, tfApplyErr := terraform.ApplyE(td, tfOpptions, terraform.ApplyOptions{
			AutoApprove: true,
			Vars: []terraformcore.TFInputVariable{
				{
					Name:  "is_enabled",
					Value: "true",
				}},
		})
		if tfApplyErr != nil {
			ux.Msg.ShowError(tui.MessageOptions{
				Message: tfInitErr.Error(),
				Error:   tfApplyErr,
			})
		}

		_, tfDestroyErr := terraform.DestroyE(td, tfOpptions, terraform.DestroyOptions{
			AutoApprove: true,
			Vars: []terraformcore.TFInputVariable{
				{
					Name:  "is_enabled",
					Value: "true",
				}},
		})
		if tfDestroyErr != nil {
			ux.Msg.ShowError(tui.MessageOptions{
				Message: tfInitErr.Error(),
				Error:   tfDestroyErr,
			})
		}
	},
}
