package utils

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"time"
)

// Logger represents a simple logger with colored output
type Logger struct {
	Verbose bool
	OutFile *os.File
}

// NewLogger creates a new logger instance
func NewLogger(verbose bool, outFile *os.File) *Logger {
	return &Logger{
		Verbose: verbose,
		OutFile: outFile,
	}
}

// Info logs informational messages
func (l *Logger) Info(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", color.BlueString("[INFO]"), message)
	l.writeToFile("INFO", message)
}

// Success logs success messages
func (l *Logger) Success(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", color.GreenString("[SUCCESS]"), message)
	l.writeToFile("SUCCESS", message)
}

// Warning logs warning messages
func (l *Logger) Warning(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", color.YellowString("[WARNING]"), message)
	l.writeToFile("WARNING", message)
}

// Error logs error messages
func (l *Logger) Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", color.RedString("[ERROR]"), message)
	l.writeToFile("ERROR", message)
}

// Debug logs debug messages (only when verbose mode is enabled)
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.Verbose {
		message := fmt.Sprintf(format, args...)
		fmt.Printf("%s %s\n", color.CyanString("[DEBUG]"), message)
		l.writeToFile("DEBUG", message)
	}
}

// writeToFile writes log messages to the output file if specified
func (l *Logger) writeToFile(level, message string) {
	if l.OutFile != nil {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		_, err := fmt.Fprintf(l.OutFile, "[%s] [%s] %s\n", timestamp, level, message)
		if err != nil {
			fmt.Printf("%s Failed to write to log file: %v\n", color.RedString("[ERROR]"), err)
		}
	}
}

// Result logs a subdomain discovery result
func (l *Logger) Result(subdomain string, found bool, ip string) {
	if found {
		message := fmt.Sprintf("Found: %s [%s]", subdomain, ip)
		fmt.Printf("%s %s\n", color.GreenString("[FOUND]"), message)
		l.writeToFile("FOUND", message)
	} else if l.Verbose {
		message := fmt.Sprintf("Not Found: %s", subdomain)
		fmt.Printf("%s %s\n", color.RedString("[NOT FOUND]"), message)
		l.writeToFile("NOT FOUND", message)
	}
}

// ServiceResult logs a service discovery result
func (l *Logger) ServiceResult(subdomain string, port int, service string, info string) {
	message := fmt.Sprintf("%s:%d - %s - %s", subdomain, port, service, info)
	fmt.Printf("%s %s\n", color.CyanString("[SERVICE]"), message)
	l.writeToFile("SERVICE", message)
}

// FileResult logs a file extraction result
func (l *Logger) FileResult(subdomain string, filePath string, success bool, size int64) {
	if success {
		message := fmt.Sprintf("Extracted: %s - %s (Size: %d bytes)", subdomain, filePath, size)
		fmt.Printf("%s %s\n", color.MagentaString("[FILE]"), message)
		l.writeToFile("FILE", message)
	} else if l.Verbose {
		message := fmt.Sprintf("Failed to extract: %s - %s", subdomain, filePath)
		fmt.Printf("%s %s\n", color.YellowString("[FILE]"), message)
		l.writeToFile("FILE", message)
	}
}