// Package system provides the system collector.
package system

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/aide-family/laurel/internal/config"
)

func NewSystemCollector(config *config.SystemCollectorConfig) []prometheus.Collector {
	systemCollector := &SystemCollector{}
	systemCollector.AppendCollector(NewCPUCollector, &config.CPUUsage)
	return systemCollector.collectors
}

type SystemCollector struct {
	collectors []prometheus.Collector
}

type CollectorFunc func(config *config.Usage) (prometheus.Collector, error)

func (s *SystemCollector) AppendCollector(f CollectorFunc, config *config.Usage) *SystemCollector {
	collector, err := f(config)
	if err != nil {
		slog.Warn("failed to create collector", "error", err)
		return s
	}
	s.collectors = append(s.collectors, collector)
	return s
}
