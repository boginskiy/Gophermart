package logg

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TakeConfig(fileName string) (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.OutputPaths = []string{"stdout", fileName}
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.TimeKey = "timestamp"

	return config.Build(zap.AddStacktrace(zap.ErrorLevel))
}
