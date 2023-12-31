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

		td, err := terradagger.New(ctx, &terradagger.Options{
			Workspace: "../",
		})

		defer td.DaggerBackend.Close()

		if err != nil {
			ux.Msg.ShowError(tui.MessageOptions{
				Message: "Unable to create a new terradagger client",
				Error:   err,
			})
			os.Exit(1)
		}

		terraformOptions := &terraform.Options{
			TerraformModulePath: "test-data/terraform/root-module-1",
		}

		initErr := terraform.Init(td, terraformOptions, nil)
		if initErr != nil {
			ux.Msg.ShowError(tui.MessageOptions{
				Message: "Unable to run terraform init",
				Error:   initErr,
			})
		}
		planErr := terraform.Plan(td, terraformOptions, &terraform.PlanOptions{
			Vars: map[string]interface{}{
				"is_enabled": true,
			},
		})

		if planErr != nil {
			ux.Msg.ShowError(tui.MessageOptions{
				Message: "Unable to run terraform plan",
				Error:   planErr,
			})
		}

		// _ = terraform.Apply(td, terraformOptions, &terraform.ApplyOptions{
		// 	Vars: map[string]interface{}{
		// 		"is_enabled": true,
		// 	},
		// 	PreserveTFState: true, // Preserve the state file.
		// })
		// _ = terraform.Destroy(td, terraformOptions, &terraform.DestroyOptions{
		// 	Vars: map[string]interface{}{
		// 		"is_enabled": true,
		// 	},
		// })
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
