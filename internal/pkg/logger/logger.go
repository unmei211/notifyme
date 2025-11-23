package logger

import (
	"time"

	"github.com/pkg/errors"
	log "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ConfigurerGetter func(cfg *Config) log.Option

type LogConfigurer struct {
	getters []ConfigurerGetter
}

func (l *LogConfigurer) apply(cfg *Config) []log.Option {
	var options []log.Option
	for i := range l.getters {
		options = append(options, l.getters[i](cfg))
	}
	return options
}

type encoderModifier func(cfg *log.Config) *log.Config

func modifyJson(cfg *log.Config) *log.Config {
	cfg.Encoding = "json"
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return cfg
}
func modifyConsole(cfg *log.Config) *log.Config {
	cfg.Encoding = "console"

	cfg.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:       "T",
		LevelKey:      "L",
		NameKey:       "N",
		CallerKey:     "C",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "M",
		StacktraceKey: "S",
		LineEnding:    zapcore.DefaultLineEnding,

		EncodeLevel: func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			var color string
			switch level {
			case zapcore.DebugLevel:
				color = "\033[36m"
			case zapcore.InfoLevel:
				color = "\033[32m"
			case zapcore.WarnLevel:
				color = "\033[33m"
			case zapcore.ErrorLevel:
				color = "\033[31m"
			case zapcore.FatalLevel:
				color = "\033[31;1m"
			default:
				color = "\033[0m"
			}
			enc.AppendString(color + level.CapitalString() + "\033[0m")
		},

		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("\033[90m" + t.Format("2006-01-02 15:04:05.000") + "\033[0m")
		},

		EncodeDuration: zapcore.StringDurationEncoder,

		EncodeCaller: func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("\033[90m" + caller.TrimmedPath() + "\033[0m")
		},

		EncodeName: func(s string, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("\033[97m" + s + "\033[0m")
		},
	}

	return cfg
}

func InitLogger(
	cfg *Config,
) (*log.SugaredLogger, error) {
	configModeMap := map[string]log.Config{
		"production": log.NewProductionConfig(),
	}

	levelsMap := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
	}

	formatMap := map[string]encoderModifier{
		"json":    modifyJson,
		"console": modifyConsole,
	}

	zapConf := configModeMap["production"]
	zapConf.Level = log.NewAtomicLevelAt(levelsMap[cfg.LogLevel])
	formatMap[cfg.LogFormat](&zapConf)

	build, err := zapConf.Build()
	if err != nil {
		return nil, errors.New("can't create config")
	}
	return build.Sugar(), nil
}

type Config struct {
	LogLevel  string `mapstructure:"level"`
	LogFormat string `mapstructure:"format"`
}
