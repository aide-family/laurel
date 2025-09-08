package inspection

import (
	"fmt"
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

// getUsageStatus returns status text based on usage percentage
func GetUsageStatus(usage float64) string {
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
func GetUsageColor(usage float64) string {
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

// GetColoredUsage returns colored usage percentage
func GetColoredUsage(usage float64) string {
	color := GetUsageColor(usage)
	return fmt.Sprintf("%s%.1f%%%s", color, usage, ColorReset)
}

// GetColoredStatus returns colored status text
func GetColoredStatus(usage float64) string {
	color := GetUsageColor(usage)
	status := GetUsageStatus(usage)
	return fmt.Sprintf("%s%s%s", color, status, ColorReset)
}

// FormatBytes converts bytes to human readable format
func FormatBytes(bytes uint64) string {
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
