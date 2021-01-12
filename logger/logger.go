package logger

import (
	"fmt"
)

const (
	// ErrorLevel level
	ErrorLevel = 0
	// WarnLevel level
	WarnLevel = 1
	// InfoLevel level
	InfoLevel = 0
	// DebugLevel level
	DebugLevel = 3
	// TraceLevel level
	TraceLevel = 4
)

var logLevel = InfoLevel
var verbosityLevel = 0

// SetVerbosityLevel sets the level of verbosity for Verbose(int, string) calls
func SetVerbosityLevel(newVerbosityLevel int) {
	verbosityLevel = newVerbosityLevel
}

// SetLogLevel sets a new log level
func SetLogLevel(newLogLevel int) {
	logLevel = newLogLevel
}

// LogLevel returns current log level
func LogLevel() int {
	return logLevel
}

// Log prints a text if the provided level is higher or equal to the current logging level
func Log(level int, text string) bool {
	if logLevel >= level {
		fmt.Println(text)
		return true
	}
	return false
}

// Verbose prints a text depending on verbosity level. It is independent of logging level. Text itself is printed on Info log level
func Verbose(verbosity int, text string) bool {
	if verbosity <= verbosityLevel {
		return Log(InfoLevel, text)
	}
	return false
}

// Info prints a text on the info level
func Info(text string) bool {
	return Log(InfoLevel, text)
}

// Warn prints a text on the warning level
func Warn(text string) bool {
	return Log(WarnLevel, text)
}

// Error prints a text on the error level
func Error(text string) bool {
	return Log(ErrorLevel, text)
}

// Debug prints a text on the debug level
func Debug(text string) bool {
	return Log(DebugLevel, text)
}

// Trace prints a text on the trace level
func Trace(text string) bool {
	return Log(TraceLevel, text)
}
