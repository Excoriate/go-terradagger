package tf

import (
	"context"
	"os"

	"github.com/Excoriate/go-terradagger/cli/internal/tui"
	"github.com/Excoriate/go-terradagger/pkg/terradagger"
	"github.com/Excoriate/go-terradagger/pkg/terraform"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	recipe   string // maps to the recipe in the /examples/<recipe>
	scenario string // maps to a config/<scenario>.tfvars per each recipe.
	all      bool
)

var Cmd = &cobra.Command{
	Use:   "tf",
	Short: "Run terraform CI Jobs using Dagger",
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

		ux.Title.ShowTitle("Terraform CI/CD")

		td, err := terradagger.New(ctx, &terradagger.ClientOptions{
			RootDir: "../",
		})

		defer td.DaggerClient.Close()

		if err != nil {
			ux.Msg.ShowError(tui.MessageOptions{
				Message: "Unable to create a new terradagger",
				Error:   err,
			})
			os.Exit(1)
		}

		terraformOptions := &terraform.Options{
			TerraformDir: "modules/default",
		}

		_ = terraform.Init(td, terraformOptions, nil)
		_ = terraform.Plan(td, terraformOptions, &terraform.PlanOptions{
			Vars: map[string]interface{}{
				"aws_region": "us-east-1",
				"is_enabled": true,
			},
		})
	},
}

func AddFlags() {
	Cmd.PersistentFlags().StringVarP(&recipe, "recipe", "", "basic",
		"Recipe to run. By default, "+
			"it'll run the 'basic' recipe in the 'examples' folder.")

	Cmd.PersistentFlags().StringVarP(&scenario, "scenario", "", "fixtures",
		"Scenario to run. By default, "+
			"it'll run the 'fixtures' scenario in the 'config' folder.")

	Cmd.PersistentFlags().BoolVarP(&all, "all", "", false, "Run all recipes in the 'examples' folder.")

	_ = viper.BindPFlag("recipe", Cmd.PersistentFlags().Lookup("recipe"))
	_ = viper.BindPFlag("scenario", Cmd.PersistentFlags().Lookup("scenario"))
	_ = viper.BindPFlag("all", Cmd.PersistentFlags().Lookup("all"))

}
func init() {
	AddFlags()
}
