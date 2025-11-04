package cmd

import (
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
)

const flagName = "log-level"

func setLogLevelFlag(cmd *cobra.Command) {
	cmd.Flags().String(flagName, "info", "Log level: debug, info, warn, error")
}

func getFlagName() string {
	return flagName
}

func parseLogLevel(logStr string) slog.Level {
	var level slog.Level
	switch strings.ToLower(logStr) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		slog.Warn("Unknown log level, fallback to info.", "given", logStr)
		level = slog.LevelInfo
	}
	return level
}
