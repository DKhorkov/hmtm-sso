package logging

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"sync"
)

var (
	instance *slog.Logger
	once     sync.Once

	LogLevels = struct {
		INFO, DEBUG, WARN, ERROR slog.Level
	}{
		INFO:  slog.LevelInfo,
		DEBUG: slog.LevelDebug,
		WARN:  slog.LevelWarn,
		ERROR: slog.LevelError,
	}
)

// GetInstance implemented as singleton pattern to get Logger instance, created once for whole app:.
func GetInstance(logLevel slog.Level) *slog.Logger {
	once.Do(func() {
		instance = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: logLevel,
				},
			),
		)
	})

	return instance
}

// GetLogTraceback return a string with info about filename, function name and line
// https://stackoverflow.com/questions/25927660/how-to-get-the-current-function-name
func GetLogTraceback() string {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Sprintf("%s on line %d: %s", "Unknown", 0, "Unknown")
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return fmt.Sprintf("%s on line %d: %s", file, line, "Unknown")
	}

	return fmt.Sprintf("%s on line %d: %s", file, line, fn.Name())
}
