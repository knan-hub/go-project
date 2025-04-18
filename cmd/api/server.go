package api

import (
	"github.com/spf13/cobra"
)

var (
	configYml string
	apiCheck  bool
	StartCmd  = &cobra.Command{
		Use:          "server",
		Example:      "go-project server -c config/settings.yml",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
	}
)

func setup() {
}
