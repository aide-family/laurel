package main

import (
	"log/slog"
	"os"

	"github.com/aide-family/laurel/internal/option"
)

func main() {
	if err := option.RootCmd.Execute(); err != nil {
		slog.Error("failed to execute root command", "error", err)
		os.Exit(1)
	}
}
