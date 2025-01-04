package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.SugaredLogger
)

func InitLogger(logPath string, level string) error {
	// Ensure log directory exists
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return err
	}

	// Configure log rotation
	logRotator := &lumberjack.Logger{
		Filename:   filepath.Join(logPath, "app.log"),
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	}

	// Configure log levels
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	default:
		zapLevel = zap.InfoLevel
	}

	// Configure log encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Create core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(logRotator),
		),
		zapLevel,
	)

	// Create logger
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	logger = zapLogger.Sugar()

	return nil
}

func Debug(msg string, fields ...interface{}) {
	if logger != nil {
		logger.Debugw(msg, fields...)
	}
}

func Info(msg string, fields ...interface{}) {
	if logger != nil {
		logger.Infow(msg, fields...)
	}
}

func Warn(msg string, fields ...interface{}) {
	if logger != nil {
		logger.Warnw(msg, fields...)
	}
}

func Error(err error, msg string, fields ...interface{}) {
	if logger != nil {
		fields = append(fields, "error", err)
		logger.Errorw(msg, fields...)
	}
}

func Fatal(msg string, fields ...interface{}) {
	if logger != nil {
		logger.Fatalw(msg, fields...)
	}
}
