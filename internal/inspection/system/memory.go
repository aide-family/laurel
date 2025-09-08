package system

import (
	"fmt"

	"github.com/aide-family/laurel/internal/inspection"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shirou/gopsutil/v4/mem"
	"golang.org/x/sync/errgroup"
)

func NewMemoryInspector() inspection.Inspector {
	return &memoryInspector{}
}

type memoryInspector struct {
}

// Show implements inspection.Inspector.
func (m *memoryInspector) Show() {
	var eg errgroup.Group
	var memInfo *mem.VirtualMemoryStat
	var swapInfo *mem.SwapMemoryStat
	var err error

	eg.Go(func() error {
		// Get memory information
		memInfo, err = mem.VirtualMemory()
		return err
	})

	eg.Go(func() error {
		swapInfo, err = mem.SwapMemory()
		return err
	})

	if err := eg.Wait(); err != nil {
		fmt.Printf("%s%s Error getting memory info: %v\n", inspection.ColorRed, inspection.ColorReset, err)
		return
	}

	// Display memory basic information
	m.displayMemoryInfo(memInfo)

	// Display memory usage statistics
	m.displayMemoryUsage(memInfo)

	// Display swap information
	m.displaySwapInfo(swapInfo)
	fmt.Println()
}

// displayMemoryInfo shows basic memory information
func (m *memoryInspector) displayMemoryInfo(memInfo *mem.VirtualMemoryStat) {
	// Create memory information table
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnMagentaWhite)
	t.SetTitle("Memory Information")
	t.AppendHeader(table.Row{"Total Memory", "Available Memory", "Used Memory", "Free Memory", "Active Memory", "Inactive Memory", "Wired Memory"})
	t.AppendRow(table.Row{
		inspection.FormatBytes(memInfo.Total),
		inspection.FormatBytes(memInfo.Available),
		inspection.FormatBytes(memInfo.Used),
		inspection.FormatBytes(memInfo.Free),
		inspection.FormatBytes(memInfo.Active),
		inspection.FormatBytes(memInfo.Inactive),
		inspection.FormatBytes(memInfo.Wired),
	})

	fmt.Println(t.Render())
}

// displayMemoryUsage shows memory usage statistics
func (m *memoryInspector) displayMemoryUsage(memInfo *mem.VirtualMemoryStat) {
	usagePercent := memInfo.UsedPercent

	// Create memory usage table
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnBlueWhite)
	t.SetTitle("Memory Usage")
	t.AppendHeader(table.Row{"Usage Percentage", "Status"})

	t.AppendRow(table.Row{
		inspection.GetColoredUsage(usagePercent),
		inspection.GetColoredStatus(usagePercent),
	})

	fmt.Println(t.Render())
}

// displaySwapInfo shows swap memory information
func (m *memoryInspector) displaySwapInfo(swapInfo *mem.SwapMemoryStat) {
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnYellowWhite)

	t.SetTitle("Swap Information")
	t.AppendHeader(table.Row{"Total Swap", "Used Swap", "Free Swap", "Usage Percentage"})

	t.AppendRow(table.Row{
		inspection.FormatBytes(swapInfo.Total),
		inspection.FormatBytes(swapInfo.Used),
		inspection.FormatBytes(swapInfo.Free),
		func() string {
			if swapInfo.Total == 0 {
				return inspection.GetColoredUsage(0.0)
			}
			swapPercent := (float64(swapInfo.Used) / float64(swapInfo.Total)) * 100
			return inspection.GetColoredUsage(swapPercent)
		}(),
	})

	fmt.Println(t.Render())
}
