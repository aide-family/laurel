package inspection

import (
	"context"
	"fmt"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// Color codes for terminal output - optimized for dark themes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[91m" // Bright red
	ColorGreen  = "\033[92m" // Bright green
	ColorYellow = "\033[93m" // Bright yellow
	ColorBlue   = "\033[94m" // Bright blue
	ColorPurple = "\033[95m" // Bright magenta
	ColorCyan   = "\033[96m" // Bright cyan
	ColorWhite  = "\033[97m" // Bright white
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
	ColorGray   = "\033[90m" // Bright black/gray
)

// Icons for different status levels
const (
	IconCPU     = ""
	IconCore    = ""
	IconUsage   = ""
	IconWarning = ""
	IconSuccess = ""
	IconInfo    = ""
	IconBar     = ""
	IconEmpty   = ""
	IconMemory  = ""
	IconRAM     = ""
	IconSwap    = ""
	IconStorage = ""
)

// CPUInspection performs a comprehensive CPU inspection with professional formatting
func CPUInspection() {
	// Get CPU information
	cpuInfo, err := cpu.Info()
	if err != nil {
		fmt.Printf("%s%s Error getting CPU info: %v%s\n", ColorRed, IconWarning, err, ColorReset)
		return
	}

	// Display CPU basic information
	displayCPUInfo(cpuInfo)

	// Get CPU usage statistics
	displayCPUUsage()

	// Get per-core usage
	displayPerCoreUsage()
	fmt.Println()
}

// displayCPUInfo shows basic CPU information
func displayCPUInfo(cpuInfo []cpu.InfoStat) {
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
func displayCPUUsage() {
	// Get CPU usage with 1 second interval
	percent, err := cpu.PercentWithContext(context.Background(), 1*time.Second, false)
	if err != nil {
		fmt.Printf("%s%s Error getting CPU usage: %v%s\n", ColorRed, IconWarning, err, ColorReset)
		return
	}

	if len(percent) == 0 {
		fmt.Printf("%s%s No CPU usage data available%s\n", ColorYellow, IconWarning, ColorReset)
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
			getColoredUsage(usage),
			getColoredStatus(usage),
		})
	}

	fmt.Println(t.Render())
}

// displayPerCoreUsage shows per-core CPU usage
func displayPerCoreUsage() {
	// Get per-core usage
	percent, err := cpu.PercentWithContext(context.Background(), 1*time.Second, true)
	if err != nil {
		fmt.Printf("%s%s Error getting per-core usage: %v%s\n", ColorRed, IconWarning, err, ColorReset)
		return
	}

	if len(percent) == 0 {
		fmt.Printf("%s%s No per-core usage data available%s\n", ColorYellow, IconWarning, ColorReset)
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
			getColoredUsage(usage),
			getColoredStatus(usage),
		})
	}

	fmt.Println(t.Render())
}

// getUsageStatus returns status text based on usage percentage
func getUsageStatus(usage float64) string {
	switch {
	case usage >= 90:
		return "Critical"
	case usage >= 70:
		return "High"
	case usage >= 50:
		return "Moderate"
	default:
		return "Normal"
	}
}

// getUsageColor returns color code based on usage percentage
func getUsageColor(usage float64) string {
	switch {
	case usage >= 95:
		return "\033[91m" // Bright red for super high
	case usage >= 80:
		return "\033[91m" // Bright red for high
	case usage >= 60:
		return "\033[93m" // Bright yellow for warning
	default:
		return "\033[92m" // Bright green for normal
	}
}

// getColoredUsage returns colored usage percentage
func getColoredUsage(usage float64) string {
	color := getUsageColor(usage)
	return fmt.Sprintf("%s%.1f%%%s", color, usage, ColorReset)
}

// getColoredStatus returns colored status text
func getColoredStatus(usage float64) string {
	color := getUsageColor(usage)
	status := getUsageStatus(usage)
	return fmt.Sprintf("%s%s%s", color, status, ColorReset)
}

// MemoryInspection performs a comprehensive memory inspection with professional formatting
func MemoryInspection() {
	// Get memory information
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("%s%s Error getting memory info: %v%s\n", ColorRed, IconWarning, err, ColorReset)
		return
	}

	// Display memory basic information
	displayMemoryInfo(memInfo)

	// Display memory usage statistics
	displayMemoryUsage(memInfo)

	// Display swap information
	displaySwapInfo()
	fmt.Println()
}

// displayMemoryInfo shows basic memory information
func displayMemoryInfo(memInfo *mem.VirtualMemoryStat) {
	// Create memory information table
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnMagentaWhite)
	t.SetTitle("Memory Information")
	t.AppendHeader(table.Row{"Total Memory", "Available Memory", "Used Memory", "Free Memory", "Active Memory", "Inactive Memory", "Wired Memory"})
	t.AppendRow(table.Row{
		formatBytes(memInfo.Total),
		formatBytes(memInfo.Available),
		formatBytes(memInfo.Used),
		formatBytes(memInfo.Free),
		formatBytes(memInfo.Active),
		formatBytes(memInfo.Inactive),
		formatBytes(memInfo.Wired),
	})

	fmt.Println(t.Render())
}

// displayMemoryUsage shows memory usage statistics
func displayMemoryUsage(memInfo *mem.VirtualMemoryStat) {
	usagePercent := memInfo.UsedPercent

	// Create memory usage table
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnBlueWhite)
	t.SetTitle("Memory Usage")
	t.AppendHeader(table.Row{"Usage Percentage", "Status"})

	t.AppendRow(table.Row{
		getColoredUsage(usagePercent),
		getColoredStatus(usagePercent),
	})

	fmt.Println(t.Render())
}

// displaySwapInfo shows swap memory information
func displaySwapInfo() {
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnYellowWhite)
	swapInfo, err := mem.SwapMemory()
	if err != nil {
		// Create swap information table for unavailable swap

		t.SetTitle("Swap Information")
		t.AppendHeader(table.Row{"Status"})

		t.AppendRow(table.Row{IconWarning + " Swap not available or disabled"})
		fmt.Println(t.Render())
		return
	}

	t.SetTitle("Swap Information")
	t.AppendHeader(table.Row{"Total Swap", "Used Swap", "Free Swap", "Usage Percentage"})

	t.AppendRow(table.Row{
		formatBytes(swapInfo.Total),
		formatBytes(swapInfo.Used),
		formatBytes(swapInfo.Free),
		func() string {
			if swapInfo.Total == 0 {
				return getColoredUsage(0.0)
			}
			swapPercent := (float64(swapInfo.Used) / float64(swapInfo.Total)) * 100
			return getColoredUsage(swapPercent)
		}(),
	})

	fmt.Println(t.Render())
}

// formatBytes converts bytes to human readable format
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// getMemoryPressureText returns text for memory pressure indicator
func getMemoryPressureText(usage float64) string {
	switch {
	case usage >= 90:
		return "Critical"
	case usage >= 80:
		return "High"
	case usage >= 60:
		return "Moderate"
	default:
		return "Low"
	}
}
