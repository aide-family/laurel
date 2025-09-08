package system

import (
	"context"
	"fmt"
	"time"

	"github.com/aide-family/laurel/internal/inspection"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shirou/gopsutil/v4/cpu"
	"golang.org/x/sync/errgroup"
)

func NewCPUInspector() inspection.Inspector {
	return &cpuInspection{}
}

type cpuInspection struct {
}

// Show implements inspection.Inspector.
func (c *cpuInspection) Show() {
	var eg errgroup.Group
	var cpuInfo []cpu.InfoStat
	var percentTotal []float64
	var percentPerCore []float64
	var err error

	eg.Go(func() error {
		cpuInfo, err = cpu.Info()
		return err
	})

	eg.Go(func() error {
		percentTotal, err = cpu.PercentWithContext(context.Background(), 1*time.Second, false)
		return err
	})

	eg.Go(func() error {
		percentPerCore, err = cpu.PercentWithContext(context.Background(), 1*time.Second, true)
		return err
	})

	if err := eg.Wait(); err != nil {
		fmt.Printf("%s%s Error getting CPU info: %v\n", inspection.ColorRed, inspection.ColorReset, err)
		return
	}

	// Display CPU basic information
	c.displayCPUInfo(cpuInfo)

	// Get CPU usage statistics
	c.displayCPUUsage(percentTotal)

	// Get per-core usage
	c.displayPerCoreUsage(percentPerCore)

	fmt.Println()
}

// displayCPUInfo shows basic CPU information
func (c *cpuInspection) displayCPUInfo(cpuInfo []cpu.InfoStat) {
	if len(cpuInfo) == 0 {
		return
	}
	// Create CPU information table
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnMagentaWhite)

	t.SetTitle("CPU Information")
	t.AppendHeader(table.Row{"Model", "Total Cores", "Physical Cores", "Logical Cores", "Frequency", "Cache Size"})

	for _, info := range cpuInfo {

		// Get physical and logical core counts
		physicalCount, err := cpu.Counts(false)
		if err != nil {
			physicalCount = 0
		}
		logicalCount, err := cpu.Counts(true)
		if err != nil {
			logicalCount = 0
		}

		t.AppendRow(table.Row{
			info.ModelName,
			fmt.Sprintf("%d", int(info.Cores)),
			fmt.Sprintf("%d", physicalCount),
			fmt.Sprintf("%d", logicalCount),
			fmt.Sprintf("%.2f GHz", float64(info.Mhz)/1000),
			fmt.Sprintf("%d KB", info.CacheSize),
		})
	}

	fmt.Println(t.Render())
}

// displayCPUUsage shows overall CPU usage
func (c *cpuInspection) displayCPUUsage(percent []float64) {
	if len(percent) == 0 {
		fmt.Printf("%s%s No CPU usage data available\n", inspection.ColorYellow, inspection.ColorReset)
		return
	}

	// Create CPU usage table
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnBlueWhite)
	t.SetTitle("CPU Usage")
	t.AppendHeader(table.Row{"Core", "Overall Usage", "Status"})
	for i, usage := range percent {
		t.AppendRow(table.Row{
			fmt.Sprintf("Core %d", i+1),
			inspection.GetColoredUsage(usage),
			inspection.GetColoredStatus(usage),
		})
	}

	fmt.Println(t.Render())
}

// displayPerCoreUsage shows per-core CPU usage
func (c *cpuInspection) displayPerCoreUsage(percent []float64) {
	if len(percent) == 0 {
		fmt.Printf("%s%s No per-core usage data available\n", inspection.ColorYellow, inspection.ColorReset)
		return
	}

	// Create per-core usage table
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnYellowWhite)
	t.SetTitle("Per-Core CPU Usage")
	t.AppendHeader(table.Row{"Core", "Usage", "Status"})

	for i, usage := range percent {
		t.AppendRow(table.Row{
			fmt.Sprintf("Core %d", i+1),
			inspection.GetColoredUsage(usage),
			inspection.GetColoredStatus(usage),
		})
	}

	fmt.Println(t.Render())
}
