package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var (
	levelStrings = []string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}
)

func getCallerInfo() (string, int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown", 0
	}

	startIdx := strings.Index(file, "/cloud")
	return file[startIdx:], line
}

func logPrint(level LogLevel, format string, args ...interface{}) {
	levelStr := levelStrings[level]
	timestamp := time.Now().Format("2024-03-28 21:21:57.521")

	fp, line := getCallerInfo()
	prefix := fmt.Sprintf("[%s]-[%s]-[%s:%d]->", timestamp, levelStr, fp, line)

	fmt.Printf(prefix+" "+format+"\n", args...)

	if level == FATAL {
		os.Exit(1)
	}
}

func Debug(format string, args ...interface{}) {
	logPrint(DEBUG, format, args...)
}

func Info(format string, args ...interface{}) {
	logPrint(DEBUG, format, args...)
}

func Warning(format string, args ...interface{}) {
	logPrint(DEBUG, format, args...)
}

func Error(format string, args ...interface{}) {
	logPrint(DEBUG, format, args...)
}

func Fatal(format string, args ...interface{}) {
	logPrint(DEBUG, format, args...)
}
