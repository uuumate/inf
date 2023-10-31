package logging

import (
	"go.uber.org/zap"
)

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	defaultLogger.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	defaultLogger.Infow(msg, keysAndValues...)
}

func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	defaultLogger.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Debugw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	defaultLogger.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Errorw(msg, keysAndValues...)
}

func For(msg string, keysAndValues ...interface{}) *Logger {
	return &Logger{
		SugaredLogger: defaultLogger.With(keysAndValues...).Desugar().WithOptions(zap.AddCallerSkip(-1)).Sugar(),
	}
}

func Sync() {
	defaultLogger.Sync()
}

func GetDefaultLogger() *Logger {
	return defaultLogger
}
