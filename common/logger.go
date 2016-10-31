package common

import (
	"log"
	"os"
)

// Logger type to allow for a single point of logging configuration
type Logger struct {
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
}

// NewLogger creates a new logger based on user config.
func NewLogger() *Logger {
	info := log.New(os.Stdout, "[INFO]", 0)
	warn := log.New(os.Stdout, "[WARN]", 0)
	err := log.New(os.Stderr, "[ERROR]", 0)

	return &Logger{
		Info:  info,
		Warn:  warn,
		Error: err,
	}
}
