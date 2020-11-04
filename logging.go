package main

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/rs/zerolog"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

/*
Logger is an alias for zerolog.Logger. It overrides the useful logging methods to
match with Printf-like syntax.
*/
type Logger zerolog.Logger

// Info logs a Printf-style message at INFO mode.
func (l Logger) Info(format string, val ...interface{}) {
	zl := zerolog.Logger(l)
	msg := fmt.Sprintf(format, val...)
	zl.Info().Msg(msg)
}

// Debug logs a Printf-style message at DEBUG mode.
func (l Logger) Debug(format string, val ...interface{}) {
	zl := zerolog.Logger(l)
	msg := fmt.Sprintf(format, val...)
	zl.Debug().Msg(msg)
}

// Warn logs a Printf-style message at WARN mode.
func (l Logger) Warn(format string, val ...interface{}) {
	zl := zerolog.Logger(l)
	msg := fmt.Sprintf(format, val...)
	zl.Warn().Msg(msg)
}

// Error logs a Printf-style message at ERROR mode.
func (l Logger) Error(format string, val ...interface{}) {
	zl := zerolog.Logger(l)
	msg := fmt.Sprintf(format, val...)
	zl.Error().Msg(msg)
}

// Fatal logs a Printf-style message at FATAL mode.
func (l Logger) Fatal(format string, val ...interface{}) {
	zl := zerolog.Logger(l)
	msg := fmt.Sprintf(format, val...)
	zl.Fatal().Msg(msg)
}

// BaseLogger returns the base zerlog.Logger instance for fine grained logging.
func (l Logger) BaseLogger() zerolog.Logger {
	return zerolog.Logger(l)
}

// logger is the global logger object
var logger Logger

// logFile is the log file object
var logFile io.WriteCloser

// SetupLogger sets up the logging params
func SetupLogger() error {

	logFile = &lumberjack.Logger{
		Filename:   "logs/redical.log",
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     2,    //days
		Compress:   true, // disabled by default
		LocalTime:  false,
	}
	zl := zerolog.New(zerolog.SyncWriter(logFile)).
		With().
		Caller().
		Stack().
		Logger().
		Level(zerolog.InfoLevel)
	if global.redisDB.debug {
		zl = zl.Level(zerolog.DebugLevel)
	}
	if global.redisDB.prod {
		zl = zl.Level(zerolog.WarnLevel)
	}
	logger = Logger(zl)
	return nil
}

// TearDownLogger tears down the logging params
func TearDownLogger() {
	logger = Logger(zerolog.Nop())
	logFile.Close()
	logFile = nil

}

// LogSafeSlice safely logs a slice with first 5 elements of the slice and size of the slice
func LogSafeSlice(slice interface{}) string {

	switch val := reflect.ValueOf(slice); val.Kind() {
	case reflect.Slice:
		retSlc := make([]string, 0)
		len := val.Len()
		for i := 0; i < 5 && i < val.Len(); i++ {
			retSlc = append(retSlc, fmt.Sprintf("%v", val.Index(i)))
		}
		return strings.Join(retSlc, " ") + fmt.Sprintf(" Length: %d", len)

	default:
		return fmt.Sprintf("%v", val)
	}
}
