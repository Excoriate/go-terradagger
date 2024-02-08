package cmd

import (
	"context"
	"os"

	"github.com/Excoriate/go-terradagger/cli/cmd/tf"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var (
	workspace string
	module    string
	tfVersion string
)
var rootCmd = &cobra.Command{
	Use:   "terradagger",
	Short: "A portable way to run your infrastructure-as-code in Containers, powered by Dagger.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&workspace, "workspace", "w", "", "The workspace to run the Dagger engine in")
	rootCmd.PersistentFlags().StringVarP(&module, "module", "m", "", "The module to run the Dagger engine in")
	rootCmd.PersistentFlags().StringVarP(&tfVersion, "tf-version", "v", "", "The terraform version to use")
	_ = viper.BindPFlags(rootCmd.PersistentFlags())

	rootCmd.AddCommand(tf.Cmd)
}

func Execute() {
	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		os.Exit(1)
	}
}
