package cmd

import (
  "context"
  "github.com/Excoriate/go-terradagger/cli/cmd/tf"
  "log"
  "os"

  "github.com/spf13/viper"

  "github.com/spf13/cobra"
)

var awsRegion string
var awsAccessKeyID string
var awsSecretAccessKey string

var tfWorkDir string

var rootCmd = &cobra.Command{
  Use:   "terradagger",
  Short: "Pipelines as code using Dagger.io",
  Run: func(cmd *cobra.Command, args []string) {
    _ = cmd.Help()
  },
}

func init() {
  rootCmd.PersistentFlags().StringVarP(&awsRegion, "aws-region", "r", "us-east-1",
    "AWS region to use this terradagger, e.g. us-east-1")

  rootCmd.PersistentFlags().StringVarP(&awsAccessKeyID, "aws-access-key-id", "a", "",
    "AWS access key id to use this terradagger, e.g. AKIA...")

  rootCmd.PersistentFlags().StringVarP(&awsSecretAccessKey, "aws-secret-access-key", "s", "",
    "AWS secret access key to use this terradagger, e.g. 1a2b3c...")

  rootCmd.PersistentFlags().StringVarP(&tfWorkDir, "tf-work-dir", "w", "",
    "Path to the Terraform working directory, e.g. ./modules/<my_module>. ")

  _ = viper.BindPFlag("aws-region", rootCmd.PersistentFlags().Lookup("aws-region"))
  _ = viper.BindPFlag("aws-access-key-id", rootCmd.PersistentFlags().Lookup("aws-access-key-id"))
  _ = viper.BindPFlag("aws-secret-access-key", rootCmd.PersistentFlags().Lookup("aws-secret-access-key"))
  _ = viper.BindPFlag("tf-work-dir", rootCmd.PersistentFlags().Lookup("tf-work-dir"))

  if err := rootCmd.MarkPersistentFlagRequired("aws-region"); err != nil {
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
