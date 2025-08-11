package logger

import (
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(logDir string) *zap.Logger {

	logFileName := "app.log"
	logFile := filepath.Join(logDir, logFileName)

	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     10,
		Compress:   false,
	}

	writeSyncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(lumberjackLogger),
	)

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		zapcore.DebugLevel,
	)

	logger := zap.New(core)

	return logger
}
