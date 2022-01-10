package log

import (
	"context"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

var defaultLogger *Logger

func init() {
	defaultLogger = New()
}

func WithContext(ctx context.Context) *zap.SugaredLogger {
	var fields []interface{}
	for _, v := range defaultLogger.Field {
		fields = append(fields, zap.String(v, valueFromContext(ctx, v)))
	}
	return defaultLogger.With(fields...)
}

// Default returns a default logger instance.
func Default() *Logger {
	return defaultLogger
}

// SetDefault sets the default logger instance.
func SetDefault(l *Logger) {
	defaultLogger = l
}

// Debug prints a debug-level log by default logger instance.
func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

// Info prints a info-level log by default logger instance.
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Warn prints a warn-level log by default logger instance.
func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// Error prints a error-level log by default logger instance.
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Fatal prints a fatal-level log by default logger instance.
func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

// Panic prints a panic-level log by default logger instance.
func Panic(args ...interface{}) {
	defaultLogger.Panic(args...)
}

// Debugf prints a debug-level log with format by default logger instance.
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// Infof prints a info-level log with format by default logger instance.
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Warnf prints a warn-level log with format by default logger instance.
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

// Errorf prints a error-level log with format by default logger instance.
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// Fatalf prints a fatal-level log with format by default logger instance.
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

// Panicf prints a panic-level log with format by default logger instance.
func Panicf(format string, args ...interface{}) {
	defaultLogger.Panicf(format, args...)
}

// Debugw prints a debug-level log with json-format by default logger instance.
func Debugw(msg string, kvs ...interface{}) {
	defaultLogger.Debugw(msg, kvs...)
}

// Infow prints a info-level log with json-format by default logger instance.
func Infow(msg string, kvs ...interface{}) {
	defaultLogger.Infow(msg, kvs...)
}

// Warnw prints a warn-level log with json-format by default logger instance.
func Warnw(msg string, kvs ...interface{}) {
	defaultLogger.Warnw(msg, kvs...)
}

// Errorw prints a error-level log with json-format by default logger instance.
func Errorw(msg string, kvs ...interface{}) {
	defaultLogger.Errorw(msg, kvs...)
}

// Fatalw prints a fatal-level log with json-format by default logger instance.
func Fatalw(msg string, kvs ...interface{}) {
	defaultLogger.Fatalw(msg, kvs...)
}

// Panicw prints a panicw-level log with json-format by default logger instance.
func Panicw(msg string, kvs ...interface{}) {
	defaultLogger.Panicw(msg, kvs...)
}

func valueFromContext(ctx context.Context, key string) string {
	if ctx == nil {
		return "-"
	}
	key = strings.ToLower(key)
	if md, ok := metadata.FromIncomingContext(ctx); ok && len(md[key]) > 0 {
		return md[key][0]
	}
	return ""
}
