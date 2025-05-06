package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	once   sync.Once
)

type Config struct {
	Development bool
	LogLevel    string
}

func New(cfg Config) *zap.Logger {
	once.Do(func() {
		var err error
		logger, err = newZapLogger(cfg)
		if err != nil {
			panic("failed to initialize logger: " + err.Error())
		}
	})
	return logger
}

func newZapLogger(cfg Config) (*zap.Logger, error) {
	level := getLogLevel(cfg.LogLevel)
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create console encoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Create file encoder for production
	var fileEncoder zapcore.Encoder
	if !cfg.Development {
		fileEncoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Create log writers
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Optional: Add file logging in production
	var fileWriter zapcore.WriteSyncer
	if !cfg.Development {
		fileWriter = zapcore.AddSync(getLogFile())
	}

	// Determine cores based on environment
	var cores []zapcore.Core
	if cfg.Development {
		cores = []zapcore.Core{
			zapcore.NewCore(consoleEncoder, consoleWriter, level),
		}
	} else {
		cores = []zapcore.Core{
			zapcore.NewCore(consoleEncoder, consoleWriter, level),
			zapcore.NewCore(fileEncoder, fileWriter, level),
		}
	}

	combinedCore := zapcore.NewTee(cores...)

	zapLogger := zap.New(
		combinedCore,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return zapLogger, nil
}

func getLogLevel(levelStr string) zapcore.Level {
	switch levelStr {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func getLogFile() *os.File {
	os.MkdirAll("./logs", os.ModePerm)
	logFile, err := os.OpenFile(
		"./logs/app.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}
	return logFile
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
