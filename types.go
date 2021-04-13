package tinylog

// types.go
// Defines all the types used by the logger

import "io"

// ColorRed
// The ANSI escape sequence for the color Red
const ColorRed = "\033[31m"

// ColorGreen
// The ANSI escape sequence for the color Green
const ColorGreen = "\033[32m"

// ColorYellow
// The ANSI escape sequence for the color Yellow
const ColorYellow = "\033[33m"

// ColorBlue
// The ANSI escape sequence for the color Blue
const ColorBlue = "\033[34m"

// ColorMagenta
// The ANSI escape sequence for the color Magenta
const ColorMagenta = "\033[35m"

// ColorCyan
// The ANSI escape sequence for the color Cyan
const ColorCyan = "\033[36m"

// ColorWhite
// The ANSI escape sequence for the color White
const ColorWhite = "\033[37m"

// ColorGray
// The ANSI escape sequence for the color Gray
const ColorGray = "\033[30;1m"

// ColorReset
// The ANSI escape sequence to reset the terminal color
const ColorReset = "\033[0m"

// TraceLevel
// The integer representation of the Trace log level
const TraceLevel = 0

// DebugLevel
// The integer representation of the Debug log level
const DebugLevel = 1

// InfoLevel
// The integer representation of the Info log level
const InfoLevel = 2

// WarningLevel
// The integer representation of the Warning log level
const WarningLevel = 3

// ErrorLevel
// The integer representation of the Error log level
const ErrorLevel = 4

// FatalLevel
// The integer representation of the Fatal log level
const FatalLevel = 5

// PanicLevel
// The integer representation of the Panic log level
const PanicLevel = 6

// LevelText
// A type that defines the configured strings for each of the supported log levels
type LevelText struct {
	TRACE   string
	DEBUG   string
	INFO    string
	WARNING string
	ERROR   string
	FATAL   string
	PANIC   string
}

// Config
// A type that defines the settings for a Logger object
type Config struct {
	LogToOutput bool // Whether or not to log to stdout/stderr
	LogToFile   bool // Whether or not file logging is enabled

	DebugFile string // If file logging is enabled - The log file to send TRACE and DEBUG messages to
	InfoFile  string // If file logging is enabled - The log file to send INFO and WARNING messages to
	ErrorFile string // If file logging is enabled - The log file to send ERROR, FATAL, and PANIC messages to

	PrintTime   bool   // Whether or not to print the time
	TimePattern string // The formatting pattern to format the timestamp with: https://golang.org/src/time/format.go
	TimeColor   string // The ANSI escape sequence to use to color the time in the console

	PrintLevel bool // Whether or not to print the log level

	LogLevel int // For filtering messages; the minimum allowed log level that should be logged

	LevelTextInnerFormat string // The `fmt` string to format the level text tag inside of the [braces]
	LevelTextOuterFormat string // The `fmt` string to format the level text tag, including [braces]
	LevelTextPadding     int    // The number of characters to pad the log levels to (including trailing spaces)
	LevelTextLeftJustify bool   // Whether or not to left-justify the log levels

	LogPrefix string // The string that will always be printed at the beginning of each line, ***after the time and before the log level***; useful for printing submodules
	LogSuffix string // The string that will always be printed at the end of each line; useful for newlines

	DisableColors bool // Whether or not color output should be disabled

	TraceColor   string // The full ANSI color sequence to use for TRACE messages
	DebugColor   string // The full ANSI color sequence to use for DEBUG messages
	InfoColor    string // The full ANSI color sequence to use for INFO messages
	WarningColor string // The full ANSI color sequence to use for WARNING messages
	ErrorColor   string // The full ANSI color sequence to use for ERROR messages
	FatalColor   string // The full ANSI color sequence to use for FATAL messages
	PanicColor   string // The full ANSI color sequence to use for PANIC messages
	ResetColor   string // The full ANSI color sequence used to clear coloring

	LevelText LevelText // The object that stores the display text for each of the log levels; can be user-overridden

	DebugWriter io.Writer // The io.Writer to send TRACE and DEBUG messages to
	InfoWriter  io.Writer // The io.Writer to send INFO and WARNING messages to
	ErrorWriter io.Writer // The io.Writer to send ERROR, FATAL, and PANIC messages to
}

// Logger
// The definition of a logger, which is an abstraction of a Config
type Logger Config
