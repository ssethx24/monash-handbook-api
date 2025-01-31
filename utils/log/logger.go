package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorPurple  = "\033[35m"
	colorCyan    = "\033[36m"
	colorOrange  = "\033[38;5;208m" // True orange color
	colorDarkRed = "\033[0;31m"
	colorWhite   = "\033[37m"
)

type LogLevel int

const (
	LOG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
	SUCCESS
)

var (
	mu sync.Mutex
)

func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2) // 2 levels up the stack
	if !ok {
		return "unknown:0"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

func getTime() string {
	return time.Now().Format("15:04:05")
}

func log(level LogLevel, message string) {
	mu.Lock()
	defer mu.Unlock()

	caller := getCallerInfo()
	timestamp := getTime()

	var levelStr, color string

	switch level {
	case INFO:
		levelStr = "INFO"
		color = colorYellow
	case WARN:
		levelStr = "WARN"
		color = colorOrange
	case ERROR:
		levelStr = "ERROR"
		color = colorRed
	case FATAL:
		levelStr = "FATAL"
		color = colorDarkRed
	case SUCCESS:
		levelStr = "SUCCESS"
		color = colorGreen
	default:
		levelStr = "LOG"
		color = colorWhite
	}

	// Format: [time] [file:line] [LEVEL] message
	logMessage := fmt.Sprintf("%s[%s] [%s] [%s] %s%s\n",
		color,
		timestamp,
		caller,
		levelStr,
		message,
		colorReset)

	fmt.Print(logMessage)

	if level == FATAL {
		os.Exit(1)
	}
}

// Public logging functions
func Info(message string) {
	log(INFO, message)
}

func Warn(message string) {
	log(WARN, message)
}

func Error(message string) {
	log(ERROR, message)
}

func Fatal(message string) {
	log(FATAL, message)
}

func Log(message string) {
	log(LOG, message)
}

func Success(message string) {
	log(SUCCESS, message)
}

// Helper for formatted messages
func Infof(format string, args ...interface{}) {
	Info(fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...interface{}) {
	Warn(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...interface{}) {
	Error(fmt.Sprintf(format, args...))
}

func Fatalf(format string, args ...interface{}) {
	Fatal(fmt.Sprintf(format, args...))
}

func Logf(format string, args ...interface{}) {
	Log(fmt.Sprintf(format, args...))
}

func Successf(format string, args ...interface{}) {
	Success(fmt.Sprintf(format, args...))
}
