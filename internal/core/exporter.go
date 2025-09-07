// Package core provides the core functionality for the exporter.
package core

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/aide-family/laurel/internal/collectors/system"
	"github.com/aide-family/laurel/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Exporter struct {
	registry *prometheus.Registry
	config   *config.Config
	server   *http.Server
}

func NewExporter(registry *prometheus.Registry, config *config.Config) *Exporter {
	return &Exporter{registry: registry, config: config}
}

func (e *Exporter) Start(ctx context.Context) error {
	systemCollector := system.NewSystemCollector(&e.config.SystemCollectorConfig)
	e.registry.MustRegister(systemCollector...)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(e.registry, promhttp.HandlerOpts{}))
	e.server = &http.Server{
		Addr:    e.config.Server.Address,
		Handler: mux,
	}
	go func() {
		if err := e.server.ListenAndServe(); err != nil {
			slog.Error("failed to start server", "error", err)
		}
	}()
	return nil
}

func (e *Exporter) Stop(ctx context.Context) error {
	if err := e.server.Shutdown(ctx); err != nil {
		slog.Error("failed to stop server", "error", err)
		if err := e.server.Close(); err != nil {
			slog.Error("failed to close server", "error", err)
		}
		return err
	}
	return nil
}
