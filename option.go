package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/spf13/cobra"

	"github.com/aide-family/laurel/internal/config"
	"github.com/aide-family/laurel/internal/core"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "laurel",
	Short: "Moon's custom exporter",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var metricCmd = &cobra.Command{
	Use:   "metric",
	Short: "Metric commands",
	Run: func(cmd *cobra.Command, args []string) {
		registry := prometheus.NewRegistry()
		config, err := config.Load(configFile)
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

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cpu.Info())
		fmt.Println(cpu.Times(true))
		fmt.Println(cpu.CountsWithContext(context.Background(), true))
		fmt.Println(cpu.PercentWithContext(context.Background(), 1*time.Second, true))
		fmt.Println(cpu.PercentWithContext(context.Background(), 1*time.Second, false))
	},
}

func init() {
	rootCmd.AddCommand(metricCmd)
	metricCmd.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "config file")

	rootCmd.AddCommand(testCmd)
}
