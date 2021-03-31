package eslgo

import (
	"log"
)

type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
}

type NilLogger struct{}
type NormalLogger struct{}

func (l NormalLogger) Debug(format string, args ...interface{}) {
	log.Print("DEBUG: ")
	log.Printf(format, args...)
}
func (l NormalLogger) Info(format string, args ...interface{}) {
	log.Print("INFO: ")
	log.Printf(format, args...)
}
func (l NormalLogger) Warn(format string, args ...interface{}) {
	log.Print("WARN: ")
	log.Printf(format, args...)
}
func (l NormalLogger) Error(format string, args ...interface{}) {
	log.Print("ERROR: ")
	log.Printf(format, args...)
}

func (l NilLogger) Debug(string, ...interface{}) {}
func (l NilLogger) Info(string, ...interface{})  {}
func (l NilLogger) Warn(string, ...interface{})  {}
func (l NilLogger) Error(string, ...interface{}) {}
