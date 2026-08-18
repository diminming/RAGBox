package logging

import (
	"log/slog"
	"os"
)

func init() {
	// 1. Create handler options and set the log level
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug, // Display Debug level and above (Debug, Info, Warn, Error)
	}

	// 2. Pass options when initializing the handler (JSON or Text)
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)

	// 3. Set as default logger (optional)
	slog.SetDefault(logger)

}
