package main

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// logger is the global logger object
var logger *log.Logger

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

	logger = log.New(logFile, "redical: ",
		log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
	return nil
}

// TearDownLogger tears down the logging params
func TearDownLogger() {
	logger = nil
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
