package option

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/aide-family/laurel/internal/inspection"
	"github.com/aide-family/laurel/internal/inspection/system"
)

var systemInspectionOptions []string

var systemInspectionOptionKeys = []string{
	"cpu",
	"memory",
	"disk",
	"network",
	"process",
}

var systemInspectionOptionValues = map[string]inspection.Inspector{
	"cpu":     system.NewCPUInspector(),
	"memory":  system.NewMemoryInspector(),
	"disk":    system.NewDiskInspector(),
	"network": system.NewNetworkInspector(),
	"process": system.NewProcessInspector(),
}

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "System information",
	Run: func(cmd *cobra.Command, args []string) {
		options := strings.Join(systemInspectionOptions, ",")
		inspectors := make([]inspection.Inspector, 0, 10)
		for _, option := range strings.Split(options, ",") {
			if option == "all" || option == "a" {
				inspectors = make([]inspection.Inspector, 0, 10)
				for _, option := range systemInspectionOptionKeys {
					inspectors = append(inspectors, systemInspectionOptionValues[option])
				}
				continue
			}
			if inspector, ok := systemInspectionOptionValues[option]; ok {
				inspectors = append(inspectors, inspector)
				continue
			}
		}
		inspection.Show(inspectors...)
	},
}

func init() {
	systemCmd.Flags().StringSliceVarP(&systemInspectionOptions, "inspection", "i", []string{"all"}, "Inspection options")
}
