package middleware

import (
	"os"
	"runtime"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	globalLogger  *zap.Logger
	sugaredLogger *zap.SugaredLogger
)

// Config holds logger configuration
type Config struct {
	Development       bool     `yaml:"development"`
	Level             string   `yaml:"level"`
	Encoding          string   `yaml:"encoding"` // json or console
	OutputPaths       []string `yaml:"outputPaths"`
	ErrorOutputPaths  []string `yaml:"errorOutputPaths"`
	DisableCaller     bool     `yaml:"disableCaller"`
	DisableStacktrace bool     `yaml:"disableStacktrace"`
	MaxSizeMB         int      `yaml:"maxSizeMB"`  // Max log file size in MB
	MaxBackups        int      `yaml:"maxBackups"` // Max number of old log files to retain
	MaxAgeDays        int      `yaml:"maxAgeDays"` // Max number of days to retain log files
	Compress          bool     `yaml:"compress"`   // Whether to compress rotated log files
}

// Init initializes the global logger
func Init(config Config) error {
	// Set up log level
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		return err
	}

	// Set up encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if config.Development {
		encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	}

	// Set up output syncers
	var cores []zapcore.Core

	// Add console output
	if config.Development {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(
			consoleEncoder,
			zapcore.Lock(os.Stdout),
			level,
		))
	}

	// Add file output if specified
	for _, path := range config.OutputPaths {
		if path == "stdout" || path == "stderr" {
			continue // already handled
		}

		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   path,
			MaxSize:    config.MaxSizeMB,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAgeDays,
			Compress:   config.Compress,
		})

		cores = append(cores, zapcore.NewCore(
			fileEncoder,
			fileWriter,
			level,
		))
	}

	// Combine cores
	core := zapcore.NewTee(cores...)

	// Create logger options
	opts := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}

	if !config.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	if config.DisableCaller {
		opts = append(opts, zap.WithCaller(false))
	}

	globalLogger = zap.New(core, opts...)
	sugaredLogger = globalLogger.Sugar()

	return nil
}

// Sync flushes any buffered log entries
func Sync() error {
	return globalLogger.Sync()
}

// Logger returns the global zap.Logger instance
func Logger() *zap.Logger {
	return globalLogger
}

// Sugar returns the global zap.SugaredLogger instance
func Sugar() *zap.SugaredLogger {
	return sugaredLogger
}

// WithFields creates a child logger with additional fields
func WithFields(fields ...zap.Field) *zap.Logger {
	return globalLogger.With(fields...)
}

// Helper functions for different log levels
func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

// ErrorWithStack logs an error with stack trace
func ErrorWithStack(err error, fields ...zap.Field) {
	stack := make([]byte, 4096)
	length := runtime.Stack(stack, false)
	stackTrace := strings.TrimSpace(string(stack[:length]))

	fields = append(fields, zap.String("stack", stackTrace))
	globalLogger.Error(err.Error(), fields...)
}

// Panic logs a message at panic level and then panics
func Panic(msg string, fields ...zap.Field) {
	globalLogger.Panic(msg, fields...)
}

// Debugf logs a formatted debug message
func Debugf(template string, args ...interface{}) {
	sugaredLogger.Debugf(template, args...)
}

// Infof logs a formatted info message
func Infof(template string, args ...interface{}) {
	sugaredLogger.Infof(template, args...)
}

// Warnf logs a formatted warning message
func Warnf(template string, args ...interface{}) {
	sugaredLogger.Warnf(template, args...)
}

// Errorf logs a formatted error message
func Errorf(template string, args ...interface{}) {
	sugaredLogger.Errorf(template, args...)
}

// Fatalf logs a formatted fatal message and then calls os.Exit(1)
func Fatalf(template string, args ...interface{}) {
	sugaredLogger.Fatalf(template, args...)
}

// Panicf logs a formatted panic message and then panics
func Panicf(template string, args ...interface{}) {
	sugaredLogger.Panicf(template, args...)
}

// Debugw logs a debug message with additional context
func Debugw(msg string, keysAndValues ...interface{}) {
	sugaredLogger.Debugw(msg, keysAndValues...)
}

// Infow logs an info message with additional context
func Infow(msg string, keysAndValues ...interface{}) {
	sugaredLogger.Infow(msg, keysAndValues...)
}

// Warnw logs a warning message with additional context
func Warnw(msg string, keysAndValues ...interface{}) {
	sugaredLogger.Warnw(msg, keysAndValues...)
}

// Errorw logs an error message with additional context
func Errorw(msg string, keysAndValues ...interface{}) {
	sugaredLogger.Errorw(msg, keysAndValues...)
}

// Fatalw logs a fatal message with additional context and then calls os.Exit(1)
func Fatalw(msg string, keysAndValues ...interface{}) {
	sugaredLogger.Fatalw(msg, keysAndValues...)
}

// Panicw logs a panic message with additional context and then panics
func Panicw(msg string, keysAndValues ...interface{}) {
	sugaredLogger.Panicw(msg, keysAndValues...)
}

// With creates a child logger with structured context
func With(keysAndValues ...interface{}) *zap.SugaredLogger {
	return sugaredLogger.With(keysAndValues...)
}
