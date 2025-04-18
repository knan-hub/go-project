package cmd

import (
	"errors"
	"go-project/cmd/api"
	"go-project/core/sdk/pkg"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "go-project",
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(pkg.Red("requires at least one arg"))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(api.StartCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
