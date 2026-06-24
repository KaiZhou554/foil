// Package log provides a compatibility shim for playground/log used by the vendored
// apksign code. It routes debug/warn/error calls to the Go standard log package.
package log

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "[apksign] ", log.LstdFlags|log.Lmsgprefix)
var debugEnabled = false

func init() {
	if os.Getenv("FOIL_DEBUG_APKSIGN") == "1" {
		debugEnabled = true
	}
}

// Debug logs a debug-level message. Only emitted when FOIL_DEBUG_APKSIGN=1.
func Debug(tag, msg string, args ...interface{}) {
	if debugEnabled {
		if len(args) > 0 {
			logger.Printf("DEBUG %s: %s %v", tag, msg, args)
		} else {
			logger.Printf("DEBUG %s: %s", tag, msg)
		}
	}
}

// Error logs an error-level message.
func Error(tag, msg string, args ...interface{}) {
	if len(args) > 0 {
		logger.Printf("ERROR %s: %s %v", tag, msg, args)
	} else {
		logger.Printf("ERROR %s: %s", tag, msg)
	}
}

// Warn logs a warning-level message.
func Warn(tag, msg string, args ...interface{}) {
	if len(args) > 0 {
		logger.Printf("WARN %s: %s %v", tag, msg, args)
	} else {
		logger.Printf("WARN %s: %s", tag, msg)
	}
}

// Printf is a convenience shorthand.
func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

// Print is a convenience shorthand.
func Print(args ...interface{}) {
	logger.Print(args...)
}

// Fatalf is a fatal-level log then os.Exit(1).
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

// Fatal is a fatal-level log then os.Exit(1).
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// String is a helper that returns a printable representation of the value.
func String(v interface{}) string {
	return fmt.Sprint(v)
}
