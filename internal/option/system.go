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
			case "disk":
				inspectors = append(inspectors, system.NewDiskInspector())
			case "network":
				inspectors = append(inspectors, system.NewNetworkInspector())
			case "process":
				inspectors = append(inspectors, system.NewProcessInspector())
			case "all", "a":
				inspectors = append(inspectors, system.NewCPUInspector())
				inspectors = append(inspectors, system.NewMemoryInspector())
				inspectors = append(inspectors, system.NewDiskInspector())
				inspectors = append(inspectors, system.NewNetworkInspector())
				inspectors = append(inspectors, system.NewProcessInspector())
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
