package logg

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
	if you want using Stacktrace
	return config.Build(zap.AddStacktrace(zap.ErrorLevel))
*/

func TakeConfig(fileName string) (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.OutputPaths = []string{"stdout", fileName}
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.TimeKey = "timestamp"

	return config.Build()
}
