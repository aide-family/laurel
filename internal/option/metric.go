package option

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"

	"github.com/aide-family/laurel/internal/config"
	"github.com/aide-family/laurel/internal/core"
)

var metricConfigFile string

var metricCmd = &cobra.Command{
	Use:   "metric",
	Short: "Metric commands",
	Run: func(cmd *cobra.Command, args []string) {
		registry := prometheus.NewRegistry()
		config, err := config.Load(metricConfigFile)
		if err != nil {
			slog.Error("failed to load config", "error", err)
			os.Exit(1)
		}

		exporter := core.NewExporter(registry, config)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		exporter.Start(ctx)

		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
		<-signalCh
		if err := exporter.Stop(ctx); err != nil {
			slog.Error("failed to stop exporter", "error", err)
		}
	},
}

func init() {
	metricCmd.Flags().StringVarP(&metricConfigFile, "config", "c", "config.yaml", "config file")
}
