package option

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "laurel",
	Short: "Moon's custom exporter",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(metricCmd)
	RootCmd.AddCommand(systemCmd)
}
