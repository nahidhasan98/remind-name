package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logWriter io.Writer
	logger    *log.Logger
	once      sync.Once
)

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
	DEBUG LogLevel = "DEBUG"
	WARN  LogLevel = "WARN"
)

// Init initializes the logger with config
func Init(debugMode bool, logFile string) {
	once.Do(func() {
		if debugMode {
			logWriter = os.Stdout
		} else {
			logWriter = &lumberjack.Logger{
				Filename:   logFile,
				MaxSize:    10, // megabytes
				MaxBackups: 5,
				MaxAge:     30, //days
				Compress:   true,
			}
		}
		logger = log.New(logWriter, "", log.LstdFlags|log.Lshortfile)
	})
}

func logWithLevel(level LogLevel, msg string, v ...interface{}) {
	// fallback to default logger if not initialized
	if logger == nil {
		fmt.Println("Logger not initialized, using default logger")
		Init(true, "logs/app.log") // default to console if not initialized
	}
	logger.Printf("[%s] %s", level, format(msg, v...))
}

func format(msg string, v ...interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf(msg, v...)
	}
	return msg
}

func Info(msg string, v ...interface{})  { logWithLevel(INFO, msg, v...) }
func Error(msg string, v ...interface{}) { logWithLevel(ERROR, msg, v...) }
func Debug(msg string, v ...interface{}) { logWithLevel(DEBUG, msg, v...) }
func Warn(msg string, v ...interface{})  { logWithLevel(WARN, msg, v...) }
