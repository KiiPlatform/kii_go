package kii

import (
	"log"
)

// KiiLogger Logger interface for this library
type KiiLogger interface {
	Debug(format string)
	Debugf(format string, args ...interface{})
	Info(format string)
	Infof(format string, args ...interface{})
	Warn(format string)
	Warnf(format string, args ...interface{})
	Error(format string)
	Errorf(format string, args ...interface{})
}

// DefaultLogger Default implementation for ILogger. This implementation use log.Logger.
type DefaultLogger struct {
	Logger *log.Logger
}

// Debug writes debug message to log.
func (df *DefaultLogger) Debug(message string) {
	df.Logger.Printf("[Debug] " + message)
}

// Debugf formats debug message according to a format specifier and write it to log.
func (df *DefaultLogger) Debugf(format string, v ...interface{}) {
	df.Logger.Printf("[Debug] "+format, v...)
}

// Info writes info message to log.
func (df *DefaultLogger) Info(message string) {
	df.Logger.Printf("[Info] " + message)
}

// Infof formats info message according to a format specifier and write it to log.
func (df *DefaultLogger) Infof(format string, v ...interface{}) {
	df.Logger.Printf("[Info] "+format, v...)
}

// Warn writes warn message to log.
func (df *DefaultLogger) Warn(message string) {
	df.Logger.Printf("[Warn] " + message)
}

// Warnf formats warn message according to a format specifier and write it to log.
func (df *DefaultLogger) Warnf(format string, v ...interface{}) {
	df.Logger.Printf("[Warn] "+format, v...)
}

// Error writes error message to log.
func (df *DefaultLogger) Error(message string) {
	df.Logger.Printf("[Error] " + message)
}

// Errorf formats error message according to a format specifier and write it to log.
func (df *DefaultLogger) Errorf(format string, v ...interface{}) {
	df.Logger.Printf("[Error] "+format, v...)
}
