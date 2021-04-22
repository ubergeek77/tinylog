package tinylog

// functions.go
// Contains the exported, user-accessible functions

import (
	"fmt"
	"os"
)

// DefaultLogger
// Get a new logger with a default config
func DefaultLogger() *Logger {
	return NewLogger(NewConfig())
}

// NewLogger
// Generate a new logger from a given configuration
func NewLogger(cfg Config) *Logger {
	logger := Logger(cfg)
	logger.ApplyConfig(Config(logger))
	return &logger
}

// NewTaggedLogger
// Conveniently create a new logger with a given tag
func NewTaggedLogger(logTag string, logColor string) *Logger {
	cfg := NewConfig()
	cfg.LogPrefix = GenerateTag(logTag, logColor, cfg)
	return NewLogger(cfg)
}

// NewColor
// Generate an ANSI color sequence from an escape code
// All supported ANSI sequences are supported here, even 256 bit colors
func NewColor(code string) string {
	return "\033[" + code + "m"
}

// GenerateTag
// Given some text, and a desired color, generate a string to use as a log level tag
func GenerateTag(tagText string, colorSequence string, cfg Config) string {
	// Use fmt.Sprintf to create the formatted tag with the desired inner formatting and color, resetting the color after
	formattedTag := fmt.Sprintf("[%s%s%s]", colorSequence, fmt.Sprintf(cfg.LevelTextInnerFormat, tagText), cfg.ResetColor)

	// Format the formatted tag with the desired outer format
	formattedTag = fmt.Sprintf(cfg.LevelTextOuterFormat, formattedTag)

	// Calculate the functional display length by disregarding ANSI sequences
	displayLength := len(removeAnsi(formattedTag))

	// Calculate a padding length by combining the desired padding amount to the amount of characters added by the ANSI sequence
	// This allows the tag text to be padded appropriately even with long 256-bit ANSI sequences
	padLength := cfg.LevelTextPadding + (len(formattedTag) - displayLength)

	// Negate the padding length if left justification is enabled
	if cfg.LevelTextLeftJustify {
		padLength = padLength * -1
	}

	// Finally create the padding string with the calculated padding length
	padString := "%" + fmt.Sprintf("%v", padLength) + "s"

	// Return the formatted tag with the final padding applied
	return fmt.Sprintf(padString, formattedTag)
}

// NewConfig
// Generate a new Config containing the preferred defaults
func NewConfig() Config {
	cfg := Config{}

	cfg.LogToOutput = true
	cfg.LogToFile = false

	cfg.DebugFile = ""
	cfg.InfoFile = ""
	cfg.ErrorFile = ""

	cfg.PrintTime = true
	cfg.TimePattern = "[Jan 02 2006 @ 15:04:05.000] "
	cfg.TimeColor = ColorGray

	cfg.PrintLevel = true

	cfg.LogLevel = TraceLevel

	cfg.LevelTextInnerFormat = "%7s"
	cfg.LevelTextOuterFormat = "%s "
	cfg.LevelTextPadding = 10
	cfg.LevelTextLeftJustify = false

	cfg.LogPrefix = ""
	cfg.LogSuffix = "\n"

	cfg.DisableColors = false

	cfg.TraceColor = ColorWhite
	cfg.DebugColor = ColorGreen
	cfg.InfoColor = ColorCyan
	cfg.WarningColor = ColorYellow
	cfg.ErrorColor = ColorMagenta
	cfg.FatalColor = ColorRed
	cfg.PanicColor = ColorRed
	cfg.ResetColor = ColorReset

	cfg.LevelText.TRACE = GenerateTag("TRACE", cfg.TraceColor, cfg)
	cfg.LevelText.DEBUG = GenerateTag("DEBUG", cfg.DebugColor, cfg)
	cfg.LevelText.INFO = GenerateTag("INFO", cfg.InfoColor, cfg)
	cfg.LevelText.WARNING = GenerateTag("WARNING", cfg.WarningColor, cfg)
	cfg.LevelText.ERROR = GenerateTag("ERROR", cfg.ErrorColor, cfg)
	cfg.LevelText.FATAL = GenerateTag("FATAL", cfg.FatalColor, cfg)
	cfg.LevelText.PANIC = GenerateTag("PANIC", cfg.PanicColor, cfg)

	cfg.DebugWriter = getDebugWriter(cfg)
	cfg.InfoWriter = getInfoWriter(cfg)
	cfg.ErrorWriter = getErrorWriter(cfg)

	return cfg
}

// ApplyConfig
// Apply a configuration to an existing logger
func (l *Logger) ApplyConfig(cfg Config) {
	if cfg.DisableColors {
		cfg.TraceColor = ""
		cfg.DebugColor = ""
		cfg.InfoColor = ""
		cfg.WarningColor = ""
		cfg.ErrorColor = ""
		cfg.FatalColor = ""
		cfg.PanicColor = ""
		cfg.ResetColor = ""
	}

	// Only regenerate a LevelText if the user has changed the color of it
	// We require a config object here so that the user can keep track of their own settings
	// Otherwise, ApplyConfig will overwrite the Logger's settings directly, which might lose a custom color if
	// DisableColors is enabled

	defaults := NewConfig()
	formattingChanged := cfg.LevelTextInnerFormat != defaults.LevelTextInnerFormat
	formattingChanged = formattingChanged || cfg.LevelTextOuterFormat != defaults.LevelTextOuterFormat
	formattingChanged = formattingChanged || cfg.LevelTextLeftJustify != defaults.LevelTextLeftJustify
	formattingChanged = formattingChanged || cfg.LevelTextPadding != defaults.LevelTextPadding
	if cfg.LevelText.TRACE == defaults.LevelText.TRACE && (cfg.TraceColor != defaults.TraceColor || formattingChanged) {
		cfg.LevelText.TRACE = GenerateTag("TRACE", cfg.TraceColor, cfg)
	}
	if cfg.LevelText.DEBUG == defaults.LevelText.DEBUG && (cfg.DebugColor != defaults.DebugColor || formattingChanged) {
		cfg.LevelText.DEBUG = GenerateTag("DEBUG", cfg.DebugColor, cfg)
	}
	if cfg.LevelText.INFO == defaults.LevelText.INFO && (cfg.InfoColor != defaults.InfoColor || formattingChanged) {
		cfg.LevelText.INFO = GenerateTag("INFO", cfg.InfoColor, cfg)
	}
	if cfg.LevelText.WARNING == defaults.LevelText.WARNING && (cfg.WarningColor != defaults.WarningColor || formattingChanged) {
		cfg.LevelText.WARNING = GenerateTag("WARNING", cfg.WarningColor, cfg)
	}
	if cfg.LevelText.ERROR == defaults.LevelText.ERROR && (cfg.ErrorColor != defaults.ErrorColor || formattingChanged) {
		cfg.LevelText.ERROR = GenerateTag("ERROR", cfg.ErrorColor, cfg)
	}
	if cfg.LevelText.FATAL == defaults.LevelText.FATAL && (cfg.FatalColor != defaults.FatalColor || formattingChanged) {
		cfg.LevelText.FATAL = GenerateTag("FATAL", cfg.FatalColor, cfg)
	}
	if cfg.LevelText.PANIC == defaults.LevelText.PANIC && (cfg.PanicColor != defaults.PanicColor || formattingChanged) {
		cfg.LevelText.PANIC = GenerateTag("PANIC", cfg.PanicColor, cfg)
	}

	cfg.DebugWriter = getDebugWriter(cfg)
	cfg.InfoWriter = getInfoWriter(cfg)
	cfg.ErrorWriter = getErrorWriter(cfg)

	*l = Logger(cfg)
}

// All of the logging functions below

func (l *Logger) Trace(msg string) {
	l.doLog(l.DebugWriter, TraceLevel, l.LevelText.TRACE, msg)
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	l.doLogf(l.DebugWriter, TraceLevel, l.LevelText.TRACE, format, args...)
}

func (l *Logger) Debug(msg string) {
	l.doLog(l.DebugWriter, DebugLevel, l.LevelText.DEBUG, msg)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.doLogf(l.DebugWriter, DebugLevel, l.LevelText.DEBUG, format, args...)
}

func (l *Logger) Info(msg string) {
	l.doLog(l.InfoWriter, InfoLevel, l.LevelText.INFO, msg)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.doLogf(l.InfoWriter, InfoLevel, l.LevelText.INFO, format, args...)
}

func (l *Logger) Warning(msg string) {
	l.doLog(l.InfoWriter, WarningLevel, l.LevelText.WARNING, msg)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.doLogf(l.InfoWriter, WarningLevel, l.LevelText.WARNING, format, args...)
}

func (l *Logger) Warn(msg string) {
	l.Warning(msg)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warningf(format, args...)
}

func (l *Logger) Error(msg string) {
	l.doLog(l.ErrorWriter, ErrorLevel, l.LevelText.ERROR, msg)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.doLogf(l.ErrorWriter, ErrorLevel, l.LevelText.ERROR, format, args...)
}

func (l *Logger) Fatal(msg string) {
	l.doLog(l.ErrorWriter, FatalLevel, l.LevelText.FATAL, msg)
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.doLogf(l.ErrorWriter, FatalLevel, l.LevelText.FATAL, format, args...)
	os.Exit(1)
}

func (l *Logger) Panic(msg string) {
	l.doLog(l.ErrorWriter, PanicLevel, l.LevelText.PANIC, msg)
	panic(msg)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.doLogf(l.ErrorWriter, PanicLevel, l.LevelText.PANIC, format, args...)
	panic(fmt.Sprintf(format, args...))
}
