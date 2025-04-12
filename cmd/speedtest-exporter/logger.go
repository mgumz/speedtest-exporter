package main

import (
	"fmt"
	"log/slog"
)

func logLevel(level string) (*slog.LevelVar, error) {
	ll := &slog.LevelVar{}
	switch level {
	case "debug":
		ll.Set(slog.LevelDebug)
	case "info":
		ll.Set(slog.LevelInfo)
	case "warn":
		ll.Set(slog.LevelWarn)
	case "error":
		ll.Set(slog.LevelError)
	default:
		return nil, fmt.Errorf("unsupported -log-level: %q", level)
	}
	return ll, nil
}
