package tinylog

// tinylog.go
// Contains the private core functions

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

// Retrieve a generic writer
func getWriter(logToFile bool, fileOut string, writer io.Writer) io.Writer {
	if logToFile && fileOut != "" {
		tmp, err := os.OpenFile(fileOut, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			NewLogger(Config{}).Warningf("Failed to open log file for writing: %s", fileOut)
		}
		if writer != nil {
			writer = io.MultiWriter(writer, tmp)
		} else {
			writer = tmp
		}
	}
	return writer
}

// Retrieve the DebugWriter based on the specified config
func getDebugWriter(cfg Config) io.Writer {
	var writer io.Writer
	writer = nil
	if cfg.LogToOutput {
		writer = os.Stdout
	}
	debugWriter := getWriter(cfg.LogToFile, cfg.DebugFile, writer)
	if debugWriter == nil {
		debugWriter = ioutil.Discard
	}
	return debugWriter
}

// Retrieve the InfoWriter based on the specified config
func getInfoWriter(cfg Config) io.Writer {
	var writer io.Writer
	writer = nil
	if cfg.LogToOutput {
		writer = os.Stdout
	}
	infoWriter := getWriter(cfg.LogToFile, cfg.InfoFile, writer)
	if infoWriter == nil {
		infoWriter = ioutil.Discard
	}
	return infoWriter
}

// Retrieve the ErrorWriter based on the specified config
func getErrorWriter(cfg Config) io.Writer {
	var writer io.Writer
	writer = nil
	if cfg.LogToOutput {
		writer = os.Stderr
	}
	errorWriter := getWriter(cfg.LogToFile, cfg.ErrorFile, writer)
	if errorWriter == nil {
		errorWriter = ioutil.Discard
	}
	return errorWriter
}

// Remove ANSI color sequences from a string
// Useful for quickly disabling terminal color output
func removeAnsi(str string) string {
	ansiPattern := "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	removeAnsi := regexp.MustCompile(ansiPattern)
	return removeAnsi.ReplaceAllString(str, "")
}

// Create a formatted time string based on the current config
func (l *Logger) generateTimeString() string {
	if l.PrintTime {
		return l.TimeColor + time.Now().Format(l.TimePattern) + l.ResetColor
	}
	return ""
}

// Log a message
func (l *Logger) doLog(writer io.Writer, level int, levelText string, msg string) {
	if level >= l.LogLevel {
		_, _ = writer.Write([]byte(l.generateTimeString() + l.LogPrefix + levelText + msg + l.LogSuffix))
	}
}

// Log a formatted message
func (l *Logger) doLogf(writer io.Writer, level int, levelText string, format string, args ...interface{}) {
	if level >= l.LogLevel {
		l.doLog(writer, level, levelText, fmt.Sprintf(format, args...))
	}
}
