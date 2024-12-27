package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/Casagrande-Lucas/dnd/config"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

// Logger struct encapsulates the log.Logger and the log file
type Logger struct {
	mu          sync.Mutex
	logger      *log.Logger
	file        *os.File
	level       LogLevel
	currentDate string
}

// NewLogger initializes the Logger, creating the logs directory and log file
func NewLogger(cfg *config.Config) (*Logger, error) {
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("could not create logs directory: %v", err)
	}

	currentDate := time.Now().Format("2006-01-02")
	logFileName := currentDate + ".log"
	logFilePath := filepath.Join("logs", logFileName)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("could not open log file: %v", err)
	}

	logger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)

	l := &Logger{
		logger:      logger,
		file:        logFile,
		currentDate: currentDate,
	}

	switch cfg.APP.ENV {
	case "debug":
		l.level = DEBUG
	case "dev":
		l.level = INFO
	case "release":
		l.level = ERROR
	default:
		l.level = ERROR
	}

	return l, nil
}

// Close closes the log file
func (l *Logger) Close() {
	err := l.file.Close()
	if err != nil {
		log.Fatalf("could not close log file: %v", err)
	}
}

// log writes a log entry with a specific level
func (l *Logger) log(level LogLevel, levelString string, msg string, v ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	currentDate := time.Now().Format("2006-01-02")
	if l.currentDate != currentDate {
		err := l.file.Close()
		if err != nil {
			log.Fatalf("could not close log file: %v", err)
		}

		logFileName := currentDate + ".log"
		logFilePath := filepath.Join("logs", logFileName)
		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("could not open log file: %v\n", err)
			return
		}
		l.logger.SetOutput(logFile)
		l.file = logFile
		l.currentDate = currentDate

		l.cleanOldLogs()
	}

	message := fmt.Sprintf(msg, v...)
	err := l.logger.Output(3, fmt.Sprintf("[%s] %s", levelString, message))
	if err != nil {
		log.Fatalf("could not write to log file: %v", err)
	}
}

// Debug logs an debug message
func (l *Logger) Debug(msg string, v ...interface{}) {
	l.log(DEBUG, "DEBUG", msg, v...)
}

// Info logs an informational message
func (l *Logger) Info(msg string, v ...interface{}) {
	l.log(INFO, "INFO", msg, v...)
}

// Warning logs a warning message
func (l *Logger) Warning(msg string, v ...interface{}) {
	l.log(WARNING, "WARNING", msg, v...)
}

// Error logs an error message
func (l *Logger) Error(msg string, v ...interface{}) {
	l.log(ERROR, "ERROR", msg, v...)
}

// cleanOldLogs removes log files older than the most recent 5
func (l *Logger) cleanOldLogs() {
	entries, err := os.ReadDir("logs")
	if err != nil {
		fmt.Printf("could not read logs directory: %v\n", err)
		return
	}

	var logFiles []string

	for _, entry := range entries {
		if !entry.IsDir() && len(entry.Name()) == len("2006-01-02.log") {
			if _, err := time.Parse("2006-01-02.log", entry.Name()); err == nil {
				logFiles = append(logFiles, entry.Name())
			}
		}
	}

	sort.Strings(logFiles)

	if len(logFiles) > 5 {
		filesToDelete := logFiles[:len(logFiles)-5]
		for _, fileName := range filesToDelete {
			err := os.Remove(filepath.Join("logs", fileName))
			if err != nil {
				fmt.Printf("could not remove log file %s: %v\n", fileName, err)
			}
		}
	}
}
