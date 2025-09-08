package option

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aide-family/laurel/internal/inspection"
	"github.com/aide-family/laurel/internal/inspection/system"
)

var systemInspectionOptions []string

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "System information",
	Run: func(cmd *cobra.Command, args []string) {
		options := strings.Join(systemInspectionOptions, ",")
		inspectors := make([]inspection.Inspector, 0, 10)
		for _, option := range strings.Split(options, ",") {
			switch strings.ToLower(option) {
			case "cpu":
				inspectors = append(inspectors, system.NewCPUInspector())
			case "memory":
				inspectors = append(inspectors, system.NewMemoryInspector())
			// case "disk":
			// 	inspection.DiskInspection()
			// case "network":
			// 	inspection.NetworkInspection()
			// case "process":
			// 	inspection.ProcessInspection()
			case "all", "a":
				inspectors = append(inspectors, system.NewCPUInspector())
				inspectors = append(inspectors, system.NewMemoryInspector())
				continue
			default:
				fmt.Println("Invalid inspection option: ", option)
			}
		}
		inspection.Show(inspectors...)
	},
}

func init() {
	systemCmd.Flags().StringSliceVarP(&systemInspectionOptions, "inspection", "i", []string{"all"}, "Inspection options")
}
