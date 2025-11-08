package log

import (
	"context"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	loggerKey = "logger"
)

func ContextWithLogger(
	parent context.Context,
	logger *zap.Logger,
) context.Context {
	return context.WithValue(parent, loggerKey, logger)
}

func FromContext(
	ctx context.Context,
) *zap.Logger {
	logger := ctx.Value(loggerKey).(*zap.Logger)
	return logger
}

type LoggerConfig struct {
	Level     string `env:"LOG_LEVEL" env-default:"info"`
	LogFormat string `env:"LOG_FORMAT" env-default:"json"`
	LogPath   string `env:"LOG_PATH" env-default:"var/log/qs"`
	LogFile   bool   `env:"LOG_FILE" env-default:"false"`
	LogStdOut bool   `env:"LOG_STDOUT" env-default:"false"`
	LogStdErr bool   `env:"LOG_STDERR" env-default:"false"`
}

func NewConfig() (*LoggerConfig, error) {
	var cfg LoggerConfig
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func SetupLogger(appID string) (*zap.Logger, error) {
	cfg, err := NewConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to setup logger: %w", err)
	}

	level := getLogLevel(cfg.Level)
	baseCfg := getBaseConfig(cfg.LogFormat)
	baseCfg.Level.SetLevel(level)

	outPaths := make([]string, 0)
	if cfg.LogFile {
		outPaths = append(outPaths, getLogPath(cfg.LogPath, appID))
	}
	if cfg.LogStdOut {
		outPaths = append(outPaths, "stdout")
	}
	if cfg.LogStdErr {
		outPaths = append(outPaths, "stderr")
	}

	baseCfg.OutputPaths = outPaths

	log, err := baseCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to setup logger: %w", err)
	}

	log = log.With(zap.String("app_id", appID))

	return log, nil
}

func getLogPath(path, appID string) string {
	return fmt.Sprintf("%s/%s.log", path, appID)
}

func getBaseConfig(format string) zap.Config {
	switch format {
	case "text":
		return zap.NewDevelopmentConfig()
	case "json":
		return zap.NewProductionConfig()
	}

	return zap.NewProductionConfig()
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	}

	return zap.InfoLevel
}
