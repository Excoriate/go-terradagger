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

var TgCMD = &cobra.Command{
	Use:   "tg",
	Short: "Execute terraform commands using go-terradagger",
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

		tgOptions :=
			terraformcore.WithOptions(td, &terraformcore.TfOptions{
				ModulePath:                   viper.GetString("module"),
				EnableSSHPrivateGit:          true,
				TerraformVersion:             viper.GetString("terraform-version"),
				EnvVarsToInjectByKeyFromHost: []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_SESSION_TOKEN"},
			})

		_, tgInitErr := terragrunt.InitE(td, tgOptions, terragrunt.InitOptions{}, nil)
		if tgInitErr != nil {
			ux.Msg.ShowError(tui.MessageOptions{
				Message: tgInitErr.Error(),
				Error:   tgInitErr,
			})
		}

		_, tgPlanErr := terragrunt.PlanE(td, tgOptions, terragrunt.PlanOptions{}, nil)

		if tgPlanErr != nil {
			ux.Msg.ShowError(tui.MessageOptions{
				Message: tgInitErr.Error(),
				Error:   tgPlanErr,
			})
		}
	},
}
