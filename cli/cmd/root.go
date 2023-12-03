package cmd

import (
	"context"
	"log"
	"os"

	"github.com/Excoriate/go-terradagger/cli/cmd/tf"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var srcDir string
var rootCmd = &cobra.Command{
	Use:   "terradagger",
	Short: "A portable way to run your infrastructure-as-code in Containers, powered by Dagger.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&srcDir, "src-dir", "s", "", "The directory containing your Daggerfile and Terraform code.")
	_ = viper.BindPFlag("src-dir", rootCmd.PersistentFlags().Lookup("src-dir"))

	if err := rootCmd.MarkPersistentFlagRequired("src-dir"); err != nil {
		log.Fatal(err)
	}

	rootCmd.AddCommand(tf.Cmd)
}

func Execute() {
	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		os.Exit(1)
	}
}
