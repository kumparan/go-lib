package log

import (
	"strings"

	"github.com/kumparan/go-lib/errors"
	"go.uber.org/zap"
)

// logging library using uber zap

var (
	logger       *zap.Logger
	sugared      *zap.SugaredLogger
	currentLevel level
)

type level int

// level of log
const (
	DebugLevel level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

const (
	DebugLevelString = "debug"
	InfoLevelString  = "info"
	WarnLevelString  = "warn"
	ErrorLevelString = "error"
	FatalLevelString = "fatal"
)

func init() {
	SetLevel(InfoLevel)
}

func newZapConfig() zap.Config {
	config := zap.NewProductionConfig()
	config.DisableCaller = true
	config.DisableStacktrace = true
	return config
}

// SetLevel will set level to logger and create a new logger based on level
func SetLevel(l level) {
	config := newZapConfig()
	switch l {
	case DebugLevel:
		config.Level.SetLevel(zap.DebugLevel)
	case InfoLevel:
		config.Level.SetLevel(zap.InfoLevel)
	case WarnLevel:
		config.Level.SetLevel(zap.WarnLevel)
	case ErrorLevel:
		config.Level.SetLevel(zap.ErrorLevel)
	case FatalLevel:
		config.Level.SetLevel(zap.FatalLevel)
	default:
		config.Level.SetLevel(zap.InfoLevel)
	}
	logger, _ = config.Build()
	defer logger.Sync()
	sugared = logger.Sugar()
	defer sugared.Sync()
	currentLevel = l
}

// GetLevel return log level in string
func GetLevel() string {
	return levelToString(currentLevel)
}

func SetLevelString(l string) {
	SetLevel(stringToLevel(l))
}

func stringToLevel(s string) level {
	switch strings.ToLower(s) {
	case DebugLevelString:
		return DebugLevel
	case InfoLevelString:
		return InfoLevel
	case WarnLevelString:
		return WarnLevel
	case ErrorLevelString:
		return ErrorLevel
	case FatalLevelString:
		return FatalLevel
	default:
		return InfoLevel
	}
}

func levelToString(l level) string {
	switch l {
	case DebugLevel:
		return DebugLevelString
	case InfoLevel:
		return InfoLevelString
	case WarnLevel:
		return WarnLevelString
	case ErrorLevel:
		return ErrorLevelString
	case FatalLevel:
		return FatalLevelString
	default:
		return InfoLevelString
	}
}

func Debug(args ...interface{}) {
	sugared.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	sugared.Debugf(format, args...)
}

func Debugw(msg string, keyAndValues ...interface{}) {
	sugared.Debugw(msg, keyAndValues)
}

func Print(args ...interface{}) {
	sugared.Info(args...)
}

func Println(args ...interface{}) {
	sugared.Info(args...)
}

func Printf(format string, args ...interface{}) {
	sugared.Infof(format, args...)
}

func Printw(msg string, keyAndValues ...interface{}) {
	sugared.Infow(msg, keyAndValues...)
}

func Info(args ...interface{}) {
	sugared.Info(args...)
}

func Infof(format string, args ...interface{}) {
	sugared.Infof(format, args...)
}

func Infow(msg string, keyAndValues ...interface{}) {
	sugared.Infow(msg, keyAndValues...)
}

func Warn(args ...interface{}) {
	sugared.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	sugared.Warnf(format, args...)
}

func Warnw(msg string, keyAndValues ...interface{}) {
	sugared.Warnw(msg, keyAndValues...)
}

func Error(args ...interface{}) {
	sugared.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	sugared.Errorf(format, args...)
}

func Errorw(msg string, keyAndValues ...interface{}) {
	sugared.Errorw(msg, keyAndValues...)
}

func Errors(err error) {
	var (
		errFields errors.Fields
		file      string
		line      int
	)
	switch err.(type) {
	case *errors.Errs:
		errs := err.(*errors.Errs)
		errFields = errs.GetFields()
		file, line = errs.GetFileAndLine()
	}
	if line != 0 {
		errFields["err_file"] = formatFilePath(file)
		errFields["err_line"] = line
	}
	intf := errFields.ToArrayInterface()
	sugared.With(intf...).Error(err.Error())
}

func Fatal(args ...interface{}) {
	sugared.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	sugared.Fatalf(format, args...)
}

func Fatalw(format string, keyAndValues ...interface{}) {
	sugared.Fatalw(format, keyAndValues...)
}

func With(args ...interface{}) *zap.SugaredLogger {
	return sugared.With(args...)
}

func formatFilePath(f string) string {
	slash := strings.LastIndex(f, "/")
	return f[slash+1:]
}
