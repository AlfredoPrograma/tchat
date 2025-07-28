package log

import (
	log "log"
	"os"
)

type LogLevel string

const (
	LOG_LEVEL_INFO  LogLevel = "INFO"
	LOG_LEVEL_ERROR LogLevel = "ERROR"
	LOG_LEVEL_FATAL LogLevel = "FATAL"
)

func Log(level LogLevel, msg string) {
	log.Printf("[%s] %s", level, msg)

	if level == LOG_LEVEL_FATAL {
		os.Exit(1)
	}
}
