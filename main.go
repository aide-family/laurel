package main

import (
	"log/slog"
	"os"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("failed to execute root command", "error", err)
		os.Exit(1)
	}
}
