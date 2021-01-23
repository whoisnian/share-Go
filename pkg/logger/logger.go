package logger

import (
	"fmt"
	"log"
	"os"
)

var lout *log.Logger
var lerr *log.Logger
var debug bool = false

func init() {
	lout = log.New(os.Stdout, "", log.LstdFlags)
	lerr = log.New(os.Stderr, "", log.LstdFlags)
}

// SetDebug ...
func SetDebug(value bool) {
	debug = value
}

// Info ...
func Info(v ...interface{}) {
	lout.Output(2, fmt.Sprint(v...)+"\n")
}

// Error ...
func Error(v ...interface{}) {
	lerr.Output(2, fmt.Sprint(v...)+"\n")
}

// Debug ...
func Debug(v ...interface{}) {
	if debug {
		lout.Output(2, fmt.Sprint(v...)+"\n")
	}
}

// Panic ...
func Panic(v ...interface{}) {
	msg := fmt.Sprint(v...)
	lerr.Output(2, msg+"\n")
	panic(msg)
}

// Fatal ...
func Fatal(v ...interface{}) {
	lerr.Output(2, fmt.Sprint(v...)+"\n")
	os.Exit(1)
}
