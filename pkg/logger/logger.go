package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	// Logger is the global logger instance
	Logger *log.Logger
	// LogFile is the file handle for the log file
	file *os.File
)

// Init initializes the logger with file output
func Init(logDir string) error {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	logFile := filepath.Join(logDir, fmt.Sprintf("api-%s.log", timestamp))

	var err error
	file, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// Initialize logger with file output
	Logger = log.New(file, "[API] ", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}

// Close closes the log file
func Close() {
	if file != nil {
		file.Close()
	}
}
