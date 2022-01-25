package email

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

func NewLogger() *zap.SugaredLogger {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}
	profile := os.Getenv("PROFILE")
	serviceName := os.Getenv("SERVICE_NAME")

	//env.Logger.Infof("Resetting log level to %s", logLevel)

	logger := GetLogger(profile, getLogLevel(logLevel), serviceName) // reset the logger after env load
	//env.Logger = logger
	return logger
}

func GetLogger(env string, level zapcore.Level, serviceName string) *zap.SugaredLogger {
	var slogger *zap.SugaredLogger
	switch env {
	case "test":
		slogger = zap.NewNop().Sugar()
	case "local":
		cfg := zap.Config{
			Level:            zap.NewAtomicLevelAt(level),
			Development:      true,
			Encoding:         "console",
			EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}
		logger, _ := cfg.Build()
		slogger = logger.Sugar()
	default:
		cfg := zap.Config{
			Level:       zap.NewAtomicLevelAt(level),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}

		logger, _ := cfg.Build()
		slogger = logger.With(zap.String("service", serviceName), zap.String("source", "go")).Sugar() // reconfigure with default field
	}

	return slogger
}

func getLogLevel(level string) zapcore.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return zapcore.DebugLevel
	case "INFO":
		return zapcore.InfoLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
