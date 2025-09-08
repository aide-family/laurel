package inspection

import (
	"context"
	"fmt"
	"strings"
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
	fmt.Printf("\n%s%s%s%s%s%s CPU System Inspection %s%s%s%s\n",
		ColorGray, strings.Repeat("=", 25), ColorReset,
		ColorBold, ColorCyan, IconCPU, IconCPU, ColorReset,
		ColorGray, strings.Repeat("=", 25))

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
	t.AppendHeader(table.Row{"Model", "Total Cores", "Physical Cores", "Logical Cores", "Frequency", "Cache Size", "Flags"})

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
			strings.Join(info.Flags, ","),
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
			fmt.Sprintf("%.1f%%", usage),
			getUsageStatus(usage),
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
			fmt.Sprintf("%.1f%%", usage),
			getUsageStatus(usage),
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

// MemoryInspection performs a comprehensive memory inspection with professional formatting
func MemoryInspection() {
	fmt.Printf("\n%s%s%s%s%s%s Memory System Inspection %s%s%s%s\n",
		ColorGray, strings.Repeat("=", 25), ColorReset,
		ColorBold, ColorCyan, IconMemory, IconMemory, ColorReset,
		ColorGray, strings.Repeat("=", 25))

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
}

// displayMemoryInfo shows basic memory information
func displayMemoryInfo(memInfo *mem.VirtualMemoryStat) {
	// Create memory information table
	t := table.NewWriter()
	t.SetTitle("Memory Information")
	t.AppendHeader(table.Row{"Property", "Value"})

	t.AppendRow(table.Row{
		IconRAM + " Total Memory",
		formatBytes(memInfo.Total),
	})

	t.AppendRow(table.Row{
		IconRAM + " Available Memory",
		formatBytes(memInfo.Available),
	})

	t.AppendRow(table.Row{
		IconRAM + " Used Memory",
		formatBytes(memInfo.Used),
	})

	t.AppendRow(table.Row{
		IconRAM + " Free Memory",
		formatBytes(memInfo.Free),
	})

	fmt.Println(t.Render())
}

// displayMemoryUsage shows memory usage statistics
func displayMemoryUsage(memInfo *mem.VirtualMemoryStat) {
	usagePercent := memInfo.UsedPercent

	// Create memory usage table
	t := table.NewWriter()
	t.SetTitle("Memory Usage")
	t.AppendHeader(table.Row{"Metric", "Value"})

	t.AppendRow(table.Row{
		IconUsage + " Usage Percentage",
		fmt.Sprintf("%.1f%%", usagePercent),
	})

	t.AppendRow(table.Row{
		IconRAM + " Memory Pressure",
		getMemoryPressureText(usagePercent),
	})

	fmt.Println(t.Render())
}

// displaySwapInfo shows swap memory information
func displaySwapInfo() {
	swapInfo, err := mem.SwapMemory()
	if err != nil {
		// Create swap information table for unavailable swap
		t := table.NewWriter()
		t.SetTitle("Swap Information")
		t.AppendHeader(table.Row{"Status"})

		t.AppendRow(table.Row{IconWarning + " Swap not available or disabled"})
		fmt.Println(t.Render())
		return
	}

	// Create swap information table
	t := table.NewWriter()
	t.SetTitle("Swap Information")
	t.AppendHeader(table.Row{"Property", "Value"})

	t.AppendRow(table.Row{
		IconSwap + " Total Swap",
		formatBytes(swapInfo.Total),
	})

	t.AppendRow(table.Row{
		IconSwap + " Used Swap",
		formatBytes(swapInfo.Used),
	})

	t.AppendRow(table.Row{
		IconSwap + " Free Swap",
		formatBytes(swapInfo.Free),
	})

	// Swap usage percentage
	if swapInfo.Total > 0 {
		swapPercent := (float64(swapInfo.Used) / float64(swapInfo.Total)) * 100

		t.AppendRow(table.Row{
			IconUsage + " Usage Percentage",
			fmt.Sprintf("%.1f%%", swapPercent),
		})
	}

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
