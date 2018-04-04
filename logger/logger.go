package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"strings"

	syslog "github.com/RackSec/srslog"
)

const (
	LogLevelTrace   = "TRACE"
	LogLevelDebug   = "DEBUG"
	LogLevelInfo    = "INFO"
	LogLevelWarning = "WARNING"
	LogLevelError   = "ERROR"
)

var logLevel = map[string]int{
	LogLevelTrace:   5,
	LogLevelDebug:   4,
	LogLevelInfo:    3,
	LogLevelWarning: 2,
	LogLevelError:   1,
}

var ptSystemName string

var activeLogLevel = strings.ToUpper(os.Getenv("LOG_LEVEL"))

func parseLogLevel() string {
	switch activeLogLevel {
	case LogLevelTrace, LogLevelDebug, LogLevelInfo, LogLevelWarning, LogLevelError:
	default:
		activeLogLevel = LogLevelTrace
	}
	return activeLogLevel
}

func getActiveLogLevel() int {
	return logLevel[activeLogLevel]
}

// SetupLogger creates logger instance to log to PaperTrail and Console. Should only called once in main function.
func SetupLogger(ptHost string, ptPort string, logLevel int) {
	activeLogLevel = parseLogLevel()
	log.SetPrefix("")
	log.SetFlags(0)

	if ptHost != "" {
		hostname, _ := os.Hostname()
		ptEndpoint := fmt.Sprintf("%s:%s", ptHost, ptPort)
		ptWriter, err := syslog.Dial("udp", ptEndpoint, syslog.LOG_INFO, hostname)

		if err != nil {
			log.Fatal("Can't connect to PaperTrail ...")
		}

		log.SetOutput(io.MultiWriter(os.Stdout, ptWriter))
	} else {
		log.Print("No papertrail transport detected. Logger only use local stdout")
	}
}

func formatterRFC3164(p syslog.Priority, hostname, tag, content string) string {
	return syslog.RFC3164Formatter(p, ptSystemName, hostname, content)
}

// SetupLoggerAuto creates logger instance to log to PaperTrail automatically without specifying PT HOST and PORT
func SetupLoggerAuto(appName string, ptEndpoint string) {
	activeLogLevel = parseLogLevel()
	log.SetPrefix("")
	log.SetFlags(0)

	if appName != "" && ptEndpoint != "" {
		ptSystemName = appName
		hostname, _ := os.Hostname()

		ptWriter, err := syslog.Dial("udp", ptEndpoint, syslog.LOG_INFO, hostname)

		if err != nil {
			log.Fatalf("Can't connect to PaperTrail: %s", err.Error())
		}

		ptWriter.SetFormatter(formatterRFC3164)

		log.SetOutput(io.MultiWriter(os.Stdout, ptWriter))
	} else {
		log.Print("Logger configured to use only local stdout")
	}
}

// Warn prints warning message to logs
func Warn(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelWarning] {
		message := fmt.Sprintf("\033[33mWARN : \033[0m"+format, v...)
		log.Print(message)
	}
}

// Trace prints trace message to logs
func Trace(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelTrace] {
		message := fmt.Sprintf("TRACE: "+format, v...)
		log.Print(message)
	}
}

// Debug prints debug message to logs
func Debug(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelDebug] {
		message := fmt.Sprintf("DEBUG: "+format, v...)
		log.Print(message)
	}
}

// Info prints info message to logs
func Info(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelInfo] {
		message := fmt.Sprintf("\033[32mINFO : \033[0m"+format, v...)
		log.Print(message)
	}
}

// Err prints error message to logs
func Err(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelError] {
		message := []interface{}{fmt.Sprintf("\033[31mERROR: \033[0m"+format, v...)}
		message = append(message, "\n", string(debug.Stack()))
		log.Print(message...)
	}
}

// Fatal calls Err and then os.Exit(1)
func Fatal(format string, v ...interface{}) {
	Err(format, v...)
	os.Exit(1)
}
