package system

import (
	"fmt"

	"github.com/aide-family/laurel/internal/inspection"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shirou/gopsutil/v4/disk"
	"golang.org/x/sync/errgroup"
)

func NewDiskInspector() inspection.Inspector {
	return &diskInspector{}
}

type diskInspector struct {
}

// Show implements inspection.Inspector.
func (d *diskInspector) Show() {
	var eg errgroup.Group
	var partitions []disk.PartitionStat
	var usage *disk.UsageStat
	var err error

	eg.Go(func() error {
		partitions, err = disk.Partitions(false)
		return err
	})

	eg.Go(func() error {
		usage, err = disk.Usage("/")
		return err
	})

	if err := eg.Wait(); err != nil {
		fmt.Printf("%s%s Error getting disk info: %v\n", inspection.ColorRed, inspection.ColorReset, err)
		return
	}

	// Display disk partitions information
	d.displayDiskPartitions(partitions)

	// Display disk usage information
	d.displayDiskUsage(usage)

	fmt.Println()
}

// displayDiskPartitions shows disk partitions information
func (d *diskInspector) displayDiskPartitions(partitions []disk.PartitionStat) {
	if len(partitions) == 0 {
		fmt.Printf("%s%s No disk partitions available\n", inspection.ColorYellow, inspection.ColorReset)
		return
	}

	// Create disk partitions table
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnMagentaWhite)
	t.SetTitle("Disk Partitions")
	t.AppendHeader(table.Row{"Device", "Mountpoint", "Fstype", "Total", "Used", "Free", "Usage %", "Status"})

	// Track seen mountpoints to avoid duplicates
	seenMountpoints := make(map[string]bool)

	for _, partition := range partitions {
		// Skip if we've already seen this mountpoint
		if seenMountpoints[partition.Mountpoint] {
			continue
		}
		seenMountpoints[partition.Mountpoint] = true

		// Get usage information for each partition
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			// If we can't get usage info, show partition info without usage
			t.AppendRow(table.Row{
				partition.Device,
				partition.Mountpoint,
				partition.Fstype,
				"N/A",
				"N/A",
				"N/A",
				"N/A",
				"N/A",
			})
			continue
		}

		// Handle zero total size to avoid division by zero
		var usagePercent float64
		if usage.Total > 0 {
			usagePercent = (float64(usage.Used) / float64(usage.Total)) * 100
		} else {
			usagePercent = 0
		}

		t.AppendRow(table.Row{
			partition.Device,
			partition.Mountpoint,
			partition.Fstype,
			inspection.FormatBytes(usage.Total),
			inspection.FormatBytes(usage.Used),
			inspection.FormatBytes(usage.Free),
			inspection.GetColoredUsage(usagePercent),
			inspection.GetColoredStatus(usagePercent),
		})
	}

	fmt.Println(t.Render())
}

// displayDiskUsage shows disk usage information
func (d *diskInspector) displayDiskUsage(usage *disk.UsageStat) {
	if usage == nil {
		fmt.Printf("%s%s No disk usage data available\n", inspection.ColorYellow, inspection.ColorReset)
		return
	}

	// Create disk usage table
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnBlueWhite)
	t.SetTitle("Disk Usage")
	t.AppendHeader(table.Row{"Path", "Total", "Used", "Free", "Usage %", "Status"})

	usagePercent := (float64(usage.Used) / float64(usage.Total)) * 100
	t.AppendRow(table.Row{
		usage.Path,
		inspection.FormatBytes(usage.Total),
		inspection.FormatBytes(usage.Used),
		inspection.FormatBytes(usage.Free),
		inspection.GetColoredUsage(usagePercent),
		inspection.GetColoredStatus(usagePercent),
	})

	fmt.Println(t.Render())
}
