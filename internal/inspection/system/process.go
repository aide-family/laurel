package system

import (
	"fmt"
	"sort"

	"github.com/aide-family/laurel/internal/inspection"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
)

func NewProcessInspector() inspection.Inspector {
	return &processInspector{}
}

type processInspector struct {
}

// processInfo holds process information for sorting
type processInfo struct {
	PID         int32
	Name        string
	CPUPercent  float64
	MemoryMB    float64
	Status      string
	CreateTime  int64
	CommandLine string
}

// Show implements inspection.Inspector.
func (p *processInspector) Show() {
	processes, err := process.Processes()
	if err != nil {
		fmt.Printf("%s%s Error getting process info: %v\n", inspection.ColorRed, inspection.ColorReset, err)
		return
	}

	// Display top 5 processes by CPU usage
	p.displayTopProcessesByCPU(processes)

	// Display top 5 processes by memory usage
	p.displayTopProcessesByMemory(processes)

	fmt.Println()
}

// displayTopProcessesByCPU shows top processes by CPU usage
func (p *processInspector) displayTopProcessesByCPU(processes []*process.Process) {
	var processInfos []processInfo

	// Limit the number of processes to check to avoid high CPU usage
	maxProcesses := 100
	if len(processes) > maxProcesses {
		processes = processes[:maxProcesses]
	}

	// Get CPU usage for processes
	for _, proc := range processes {
		cpuPercent, err := proc.CPUPercent()
		if err != nil {
			continue
		}

		name, err := proc.Name()
		if err != nil {
			name = "Unknown"
		}

		status, err := proc.Status()
		statusStr := "Unknown"
		if err == nil && len(status) > 0 {
			statusStr = status[0]
		}

		createTime, err := proc.CreateTime()
		if err != nil {
			createTime = 0
		}

		cmdline, err := proc.Cmdline()
		if err != nil {
			cmdline = ""
		}

		processInfos = append(processInfos, processInfo{
			PID:         proc.Pid,
			Name:        name,
			CPUPercent:  cpuPercent,
			Status:      statusStr,
			CreateTime:  createTime,
			CommandLine: cmdline,
		})
	}

	// Sort by CPU usage (descending)
	sort.Slice(processInfos, func(i, j int) bool {
		return processInfos[i].CPUPercent > processInfos[j].CPUPercent
	})

	// Create top processes by CPU table
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnBlueWhite)
	t.SetTitle("Top 5 Processes by CPU Usage")
	t.AppendHeader(table.Row{"PID", "Name", "CPU %", "Status", "Command"})

	// Show top 5 processes
	for i, proc := range processInfos {
		if i >= 5 {
			break
		}

		// Truncate command line if too long
		cmd := proc.CommandLine
		if len(cmd) > 50 {
			cmd = cmd[:47] + "..."
		}

		t.AppendRow(table.Row{
			fmt.Sprintf("%d", proc.PID),
			proc.Name,
			inspection.GetColoredUsage(proc.CPUPercent),
			func() string {
				switch proc.Status {
				case "R", "running":
					return fmt.Sprintf("%s%s%s", inspection.ColorGreen, proc.Status, inspection.ColorReset)
				case "S", "sleeping":
					return fmt.Sprintf("%s%s%s", inspection.ColorBlue, proc.Status, inspection.ColorReset)
				case "Z", "zombie":
					return fmt.Sprintf("%s%s%s", inspection.ColorRed, proc.Status, inspection.ColorReset)
				default:
					return proc.Status
				}
			}(),
			cmd,
		})
	}

	fmt.Println(t.Render())
}

// displayTopProcessesByMemory shows top processes by memory usage
func (p *processInspector) displayTopProcessesByMemory(processes []*process.Process) {
	var processInfos []processInfo

	// Limit the number of processes to check to avoid high CPU usage
	maxProcesses := 100
	if len(processes) > maxProcesses {
		processes = processes[:maxProcesses]
	}

	// Get memory usage for processes
	for _, proc := range processes {
		memInfo, err := proc.MemoryInfo()
		if err != nil {
			continue
		}

		// Convert bytes to MB
		memoryMB := float64(memInfo.RSS) / 1024 / 1024

		name, err := proc.Name()
		if err != nil {
			name = "Unknown"
		}

		status, err := proc.Status()
		statusStr := "Unknown"
		if err == nil && len(status) > 0 {
			statusStr = status[0]
		}

		createTime, err := proc.CreateTime()
		if err != nil {
			createTime = 0
		}

		cmdline, err := proc.Cmdline()
		if err != nil {
			cmdline = ""
		}

		processInfos = append(processInfos, processInfo{
			PID:         proc.Pid,
			Name:        name,
			MemoryMB:    memoryMB,
			Status:      statusStr,
			CreateTime:  createTime,
			CommandLine: cmdline,
		})
	}

	// Sort by memory usage (descending)
	sort.Slice(processInfos, func(i, j int) bool {
		return processInfos[i].MemoryMB > processInfos[j].MemoryMB
	})

	// Get total system memory for percentage calculation
	vmem, err := mem.VirtualMemory()
	if err != nil {
		vmem = &mem.VirtualMemoryStat{Total: 1} // Avoid division by zero
	}

	// Create top processes by memory table
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnYellowWhite)
	t.SetTitle("Top 5 Processes by Memory Usage")
	t.AppendHeader(table.Row{"PID", "Name", "Memory (MB)", "Memory %", "Status", "Command"})

	// Show top 5 processes
	for i, proc := range processInfos {
		if i >= 5 {
			break
		}

		// Calculate memory percentage
		memPercent := (float64(proc.MemoryMB*1024*1024) / float64(vmem.Total)) * 100

		// Truncate command line if too long
		cmd := proc.CommandLine
		if len(cmd) > 40 {
			cmd = cmd[:37] + "..."
		}

		t.AppendRow(table.Row{
			fmt.Sprintf("%d", proc.PID),
			proc.Name,
			fmt.Sprintf("%.1f", proc.MemoryMB),
			inspection.GetColoredUsage(memPercent),
			func() string {
				switch proc.Status {
				case "R", "running":
					return fmt.Sprintf("%s%s%s", inspection.ColorGreen, proc.Status, inspection.ColorReset)
				case "S", "sleeping":
					return fmt.Sprintf("%s%s%s", inspection.ColorBlue, proc.Status, inspection.ColorReset)
				case "Z", "zombie":
					return fmt.Sprintf("%s%s%s", inspection.ColorRed, proc.Status, inspection.ColorReset)
				default:
					return proc.Status
				}
			}(),
			cmd,
		})
	}

	fmt.Println(t.Render())
}
