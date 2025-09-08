package option

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/aide-family/laurel/internal/inspection"
)

var systemInspectionOptions []string

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "System information",
	Run: func(cmd *cobra.Command, args []string) {
		options := strings.Join(systemInspectionOptions, ",")
		if options == "" || options == "all" || options == "a" {
			showAllInspectionOptions()

			return
		}
		for _, option := range strings.Split(options, ",") {
			switch option {
			case "cpu":
				inspection.CPUInspection()
			case "memory":
				inspection.MemoryInspection()
				// case "disk":
				// 	inspection.DiskInspection()
				// case "network":
				// 	inspection.NetworkInspection()
				// case "process":
				// 	inspection.ProcessInspection()
			}
		}
	},
}

func init() {
	systemCmd.Flags().StringSliceVarP(&systemInspectionOptions, "inspection", "i", []string{"all"}, "Inspection options")
}

func showAllInspectionOptions() {
	inspection.CPUInspection()
	inspection.MemoryInspection()
}
