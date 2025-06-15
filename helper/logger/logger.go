package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	cfg    zap.Config
	logger *zap.Logger
)

type LoggerConfig struct {
	Level string
}

type Logger struct {
	Log *zap.Logger
}

func init() {

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding:    "json",
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
	}

	logger = zap.Must(cfg.Build())
	defer logger.Sync()
}

func Log() *zap.Logger {

	return logger

}
