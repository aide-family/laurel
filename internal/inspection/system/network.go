package system

import (
	"fmt"
	"time"

	"github.com/aide-family/laurel/internal/inspection"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shirou/gopsutil/v4/net"
)

func NewNetworkInspector() inspection.Inspector {
	return &networkInspector{}
}

type networkInspector struct {
}

// Show implements inspection.Inspector.
func (n *networkInspector) Show() {
	// Get initial network I/O counters
	ioCounters1, err := net.IOCounters(true)
	if err != nil {
		fmt.Printf("%s%s Error getting network info: %v\n", inspection.ColorRed, inspection.ColorReset, err)
		return
	}

	// Wait for 500ms to calculate speed (reduced from 1 second)
	time.Sleep(500 * time.Millisecond)

	// Get network I/O counters after 500ms
	ioCounters2, err := net.IOCounters(true)
	if err != nil {
		fmt.Printf("%s%s Error getting network info: %v\n", inspection.ColorRed, inspection.ColorReset, err)
		return
	}

	// Display network speed
	n.displayNetworkSpeed(ioCounters1, ioCounters2)

	fmt.Println()
}

// displayNetworkSpeed shows current network upload/download speed
func (n *networkInspector) displayNetworkSpeed(ioCounters1, ioCounters2 []net.IOCountersStat) {
	if len(ioCounters1) == 0 || len(ioCounters2) == 0 {
		fmt.Printf("%s%s No network speed data available\n", inspection.ColorYellow, inspection.ColorReset)
		return
	}

	// Create a map for quick lookup of second set of counters
	ioMap := make(map[string]net.IOCountersStat)
	for _, io := range ioCounters2 {
		ioMap[io.Name] = io
	}

	// Create network speed table
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnMagentaWhite)
	t.SetTitle("Network Speed (Current Upload/Download)")
	t.AppendHeader(table.Row{"Interface", "Upload Speed", "Download Speed", "Status"})

	var totalUploadSpeed, totalDownloadSpeed uint64

	for _, io1 := range ioCounters1 {
		io2, exists := ioMap[io1.Name]
		if !exists {
			continue
		}

		// Calculate speed (bytes per second)
		uploadSpeed := io2.BytesSent - io1.BytesSent
		downloadSpeed := io2.BytesRecv - io1.BytesRecv

		// Only show interfaces with activity
		if uploadSpeed > 0 || downloadSpeed > 0 {
			totalUploadSpeed += uploadSpeed
			totalDownloadSpeed += downloadSpeed

			// Determine status based on activity
			status := "Idle"
			if uploadSpeed > 0 || downloadSpeed > 0 {
				status = "Active"
			}

			t.AppendRow(table.Row{
				io1.Name,
				inspection.FormatBytes(uploadSpeed) + "/s",
				inspection.FormatBytes(downloadSpeed) + "/s",
				func() string {
					if status == "Active" {
						return fmt.Sprintf("%s%s%s", inspection.ColorGreen, status, inspection.ColorReset)
					}
					return fmt.Sprintf("%s%s%s", inspection.ColorGray, status, inspection.ColorReset)
				}(),
			})
		}
	}

	// Add total row
	if totalUploadSpeed > 0 || totalDownloadSpeed > 0 {
		t.AppendSeparator()
		t.AppendRow(table.Row{
			"TOTAL",
			fmt.Sprintf("%s%s%s", inspection.ColorBold, inspection.FormatBytes(totalUploadSpeed)+"/s", inspection.ColorReset),
			fmt.Sprintf("%s%s%s", inspection.ColorBold, inspection.FormatBytes(totalDownloadSpeed)+"/s", inspection.ColorReset),
			"Active",
		})
	}

	fmt.Println(t.Render())
}
