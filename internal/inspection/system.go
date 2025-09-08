package inspection

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

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
	IconCPU     = "🖥️"
	IconCore    = "⚡"
	IconUsage   = "📊"
	IconWarning = "⚠️"
	IconSuccess = "✅"
	IconInfo    = "ℹ️"
	IconBar     = "█"
	IconEmpty   = "░"
	IconMemory  = "💾"
	IconRAM     = "🧠"
	IconSwap    = "🔄"
	IconStorage = "💿"
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

	info := cpuInfo[0] // Use first CPU for general info

	// Get physical and logical core counts
	physicalCount, err := cpu.Counts(false)
	if err != nil {
		physicalCount = 0
	}
	logicalCount, err := cpu.Counts(true)
	if err != nil {
		logicalCount = 0
	}

	fmt.Printf("\n%s%s%s CPU Information %s\n", ColorBold, ColorBlue, IconInfo, ColorReset)
	fmt.Printf("%s%s%s Model: %s%s%s\n", ColorWhite, IconCPU, ColorReset, ColorCyan, info.ModelName, ColorReset)

	// Display cores information with proper formatting
	fmt.Printf("%s%s%s Cores: %s%d%s (Physical: %s%d%s, Logical: %s%d%s)\n",
		ColorWhite, IconCore, ColorReset,
		ColorYellow, int(info.Cores), ColorReset,
		ColorYellow, physicalCount, ColorReset,
		ColorYellow, logicalCount, ColorReset)

	if info.Mhz > 0 {
		fmt.Printf("%s%s%s Frequency: %s%.2f GHz%s\n",
			ColorWhite, IconCore, ColorReset,
			ColorGreen, float64(info.Mhz)/1000, ColorReset)
	}

	if info.CacheSize > 0 {
		fmt.Printf("%s%s%s Cache: %s%d KB%s\n",
			ColorWhite, IconCore, ColorReset,
			ColorPurple, info.CacheSize, ColorReset)
	}
}

// displayCPUUsage shows overall CPU usage
func displayCPUUsage() {
	fmt.Printf("\n%s%s%s Overall CPU Usage %s\n", ColorBold, ColorBlue, IconUsage, ColorReset)

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

	usage := percent[0]
	color := getUsageColor(usage)
	bar := generateUsageBar(usage, 20)

	fmt.Printf("%s%s%s Usage: %s%.1f%%%s %s%s\n",
		ColorWhite, IconUsage, ColorReset,
		color, usage, ColorReset,
		ColorGray, bar)
}

// displayPerCoreUsage shows per-core CPU usage
func displayPerCoreUsage() {
	fmt.Printf("\n%s%s%s Per-Core CPU Usage %s\n", ColorBold, ColorBlue, IconCore, ColorReset)

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

	// Display cores in a grid format
	coresPerRow := 4
	for i, usage := range percent {
		if i%coresPerRow == 0 {
			fmt.Print("\n  ")
		}

		color := getUsageColor(usage)
		bar := generateUsageBar(usage, 8)

		fmt.Printf("%sCore %2d: %s%5.1f%%%s %s%s  ",
			ColorWhite, i+1,
			color, usage, ColorReset,
			ColorGray, bar)
	}
	fmt.Println()
}

// getUsageColor returns appropriate color based on usage percentage
func getUsageColor(usage float64) string {
	switch {
	case usage >= 90:
		return ColorRed
	case usage >= 70:
		return ColorYellow
	case usage >= 50:
		return ColorBlue
	default:
		return ColorGreen
	}
}

// generateUsageBar creates a visual bar representation of usage
func generateUsageBar(usage float64, width int) string {
	filled := int(math.Round(usage / 100 * float64(width)))
	if filled > width {
		filled = width
	}

	bar := strings.Repeat(IconBar, filled) + strings.Repeat(IconEmpty, width-filled)
	return fmt.Sprintf("[%s]", bar)
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
	fmt.Printf("\n%s%s%s Memory Information %s\n", ColorBold, ColorBlue, IconInfo, ColorReset)

	// Total memory
	fmt.Printf("%s%s%s Total: %s%s%s\n",
		ColorWhite, IconRAM, ColorReset,
		ColorCyan, formatBytes(memInfo.Total), ColorReset)

	// Available memory
	fmt.Printf("%s%s%s Available: %s%s%s\n",
		ColorWhite, IconRAM, ColorReset,
		ColorGreen, formatBytes(memInfo.Available), ColorReset)

	// Used memory
	fmt.Printf("%s%s%s Used: %s%s%s\n",
		ColorWhite, IconRAM, ColorReset,
		ColorYellow, formatBytes(memInfo.Used), ColorReset)

	// Free memory
	fmt.Printf("%s%s%s Free: %s%s%s\n",
		ColorWhite, IconRAM, ColorReset,
		ColorPurple, formatBytes(memInfo.Free), ColorReset)
}

// displayMemoryUsage shows memory usage statistics
func displayMemoryUsage(memInfo *mem.VirtualMemoryStat) {
	fmt.Printf("\n%s%s%s Memory Usage %s\n", ColorBold, ColorBlue, IconUsage, ColorReset)

	usagePercent := memInfo.UsedPercent
	color := getMemoryUsageColor(usagePercent)
	bar := generateUsageBar(usagePercent, 20)

	fmt.Printf("%s%s%s Usage: %s%.1f%%%s %s%s\n",
		ColorWhite, IconUsage, ColorReset,
		color, usagePercent, ColorReset,
		ColorGray, bar)

	// Memory pressure indicator
	fmt.Printf("%s%s%s Pressure: %s%s%s\n",
		ColorWhite, IconRAM, ColorReset,
		getMemoryPressureColor(usagePercent), getMemoryPressureText(usagePercent), ColorReset)
}

// displaySwapInfo shows swap memory information
func displaySwapInfo() {
	swapInfo, err := mem.SwapMemory()
	if err != nil {
		fmt.Printf("\n%s%s%s Swap Information %s\n", ColorBold, ColorBlue, IconSwap, ColorReset)
		fmt.Printf("%s%s Swap not available or disabled%s\n", ColorYellow, IconWarning, ColorReset)
		return
	}

	fmt.Printf("\n%s%s%s Swap Information %s\n", ColorBold, ColorBlue, IconSwap, ColorReset)

	// Total swap
	fmt.Printf("%s%s%s Total: %s%s%s\n",
		ColorWhite, IconSwap, ColorReset,
		ColorCyan, formatBytes(swapInfo.Total), ColorReset)

	// Used swap
	fmt.Printf("%s%s%s Used: %s%s%s\n",
		ColorWhite, IconSwap, ColorReset,
		ColorYellow, formatBytes(swapInfo.Used), ColorReset)

	// Free swap
	fmt.Printf("%s%s%s Free: %s%s%s\n",
		ColorWhite, IconSwap, ColorReset,
		ColorGreen, formatBytes(swapInfo.Free), ColorReset)

	// Swap usage percentage
	if swapInfo.Total > 0 {
		swapPercent := (float64(swapInfo.Used) / float64(swapInfo.Total)) * 100
		color := getMemoryUsageColor(swapPercent)
		bar := generateUsageBar(swapPercent, 20)

		fmt.Printf("%s%s%s Usage: %s%.1f%%%s %s%s\n",
			ColorWhite, IconUsage, ColorReset,
			color, swapPercent, ColorReset,
			ColorGray, bar)
	}
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

// getMemoryUsageColor returns appropriate color based on memory usage percentage
func getMemoryUsageColor(usage float64) string {
	switch {
	case usage >= 90:
		return ColorRed
	case usage >= 80:
		return ColorYellow
	case usage >= 60:
		return ColorBlue
	default:
		return ColorGreen
	}
}

// getMemoryPressureColor returns color for memory pressure indicator
func getMemoryPressureColor(usage float64) string {
	switch {
	case usage >= 90:
		return ColorRed
	case usage >= 80:
		return ColorYellow
	case usage >= 60:
		return ColorBlue
	default:
		return ColorGreen
	}
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
