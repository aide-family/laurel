package system

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v4/cpu"

	"github.com/aide-family/laurel/internal/config"
)

var _ prometheus.Collector = (*cpuCollector)(nil)

func NewCPUCollector(config *config.Usage) (prometheus.Collector, error) {
	infoStats, err := cpu.InfoWithContext(context.Background())
	if err != nil {
		slog.Error("failed to get CPU info", "error", err)
		return nil, err
	}
	collector := &cpuCollector{
		config:    config,
		infoStats: infoStats,
	}
	if len(infoStats) == 0 {
		return collector, nil
	}
	info := infoStats[0]
	collector.cpuCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "system_cpu_count",
		Help: "System CPU count",
		ConstLabels: prometheus.Labels{
			"cpu":         strconv.Itoa(int(info.CPU)),
			"vendor_id":   info.VendorID,
			"family":      info.Family,
			"model":       info.Model,
			"stepping":    strconv.Itoa(int(info.Stepping)),
			"physical_id": info.PhysicalID,
			"core_id":     info.CoreID,
			"cores":       strconv.Itoa(int(info.Cores)),
			"model_name":  info.ModelName,
			"mhz":         strconv.Itoa(int(info.Mhz)),
			"cache_size":  strconv.Itoa(int(info.CacheSize)),
			"flags":       strings.Join(info.Flags, ","),
			"microcode":   info.Microcode,
		},
	})

	count, err := cpu.CountsWithContext(context.Background(), true)
	if err != nil {
		slog.Error("failed to get CPU count", "error", err)
		return nil, err
	}
	collector.cpuCount.Set(float64(count))
	cpuUsages := make([]prometheus.Gauge, 0, count)
	for i := 0; i < int(count); i++ {
		cpuUsage := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "system_cpu_usage",
			Help: "System CPU usage",
			ConstLabels: prometheus.Labels{
				"cpu": strconv.Itoa(i),
			},
		})
		cpuUsages = append(cpuUsages, cpuUsage)
	}
	collector.cpuUsages = cpuUsages

	cpuTimesUsages := make([]*prometheus.GaugeVec, 0, count*10)
	for i := 0; i < int(count); i++ {
		cpuTimesUsage := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "system_cpu_times_usage",
			Help: "System CPU times usage",
		}, []string{"cpu", "type"})
		cpuTimesUsages = append(cpuTimesUsages, cpuTimesUsage)
	}
	collector.cpuTimesUsages = cpuTimesUsages

	cpuUsageTotal := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "system_cpu_usage_total",
		Help: "System CPU usage total",
	}, []string{"cpu"})
	collector.cpuUsageTotal = append(collector.cpuUsageTotal, cpuUsageTotal)

	return collector, nil
}

type cpuCollector struct {
	config *config.Usage

	infoStats      []cpu.InfoStat
	cpuCount       prometheus.Gauge
	cpuUsages      []prometheus.Gauge
	cpuTimesUsages []*prometheus.GaugeVec
	cpuUsageTotal  []*prometheus.GaugeVec
}

// Collect implements prometheus.Collector.
func (c *cpuCollector) Collect(ch chan<- prometheus.Metric) {
	slog.Info("collecting CPU metrics")
	if !c.config.Enabled {
		slog.Warn("CPU metrics are not enabled")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()
	countValue, err := cpu.CountsWithContext(ctx, true)
	if err != nil {
		slog.Error("failed to get CPU count", "error", err)
		return
	}
	c.cpuCount.Set(float64(countValue))

	infoStats, err := cpu.InfoWithContext(ctx)
	if err != nil {
		slog.Error("failed to get CPU info", "error", err)
		return
	}
	c.infoStats = infoStats
	cpuUsage, err := cpu.PercentWithContext(ctx, 1*time.Second, true)
	if err != nil {
		slog.Error("failed to get CPU usage", "error", err)
		return
	}
	for i, usage := range cpuUsage {
		if i >= len(c.cpuUsages) {
			break
		}
		c.cpuUsages[i].Set(usage)
	}

	cpuTimes, err := cpu.TimesWithContext(ctx, true)
	if err != nil {
		slog.Error("failed to get CPU times", "error", err)
		return
	}
	for i, times := range cpuTimes {
		if i >= len(c.cpuTimesUsages) {
			break
		}
		c.cpuTimesUsages[i].WithLabelValues(times.CPU, "user").Set(times.User)
		c.cpuTimesUsages[i].WithLabelValues(times.CPU, "system").Set(times.System)
		c.cpuTimesUsages[i].WithLabelValues(times.CPU, "idle").Set(times.Idle)
		c.cpuTimesUsages[i].WithLabelValues(times.CPU, "nice").Set(times.Nice)
		c.cpuTimesUsages[i].WithLabelValues(times.CPU, "iowait").Set(times.Iowait)
		c.cpuTimesUsages[i].WithLabelValues(times.CPU, "irq").Set(times.Irq)
		c.cpuTimesUsages[i].WithLabelValues(times.CPU, "softirq").Set(times.Softirq)
		c.cpuTimesUsages[i].WithLabelValues(times.CPU, "steal").Set(times.Steal)
		c.cpuTimesUsages[i].WithLabelValues(times.CPU, "guest").Set(times.Guest)
		c.cpuTimesUsages[i].WithLabelValues(times.CPU, "guestNice").Set(times.GuestNice)
	}

	cpuUsageTotal, err := cpu.PercentWithContext(ctx, 1*time.Second, false)
	if err != nil {
		slog.Error("failed to get CPU usage total", "error", err)
		return
	}
	for i, usage := range cpuUsageTotal {
		c.cpuUsageTotal[i].WithLabelValues(strconv.Itoa(i)).Set(usage)
	}

	c.cpuCount.Collect(ch)
	for _, cpuUsage := range c.cpuUsages {
		cpuUsage.Collect(ch)
	}
	for _, cpuUsageTotal := range c.cpuUsageTotal {
		cpuUsageTotal.Collect(ch)
	}

	for _, cpuTimesUsage := range c.cpuTimesUsages {
		cpuTimesUsage.Collect(ch)
	}
}

// Describe implements prometheus.Collector.
func (c *cpuCollector) Describe(ch chan<- *prometheus.Desc) {
	c.cpuCount.Describe(ch)
	for _, cpuUsage := range c.cpuUsages {
		cpuUsage.Describe(ch)
	}
	for _, cpuTimesUsage := range c.cpuTimesUsages {
		cpuTimesUsage.Describe(ch)
	}
	for _, cpuUsageTotal := range c.cpuUsageTotal {
		cpuUsageTotal.Describe(ch)
	}
}
