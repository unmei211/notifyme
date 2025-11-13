package logger

import (
	"os"

	log "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func logLevelModifier(cfg *LoggerConfig) log.Option {
	return log.WrapCore(func(core zapcore.Core) zapcore.Core {
		var lvl zapcore.Level
		if err := lvl.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
			lvl = zapcore.InfoLevel
		}
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(log.NewProductionEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			lvl,
		)
	})
}

type ConfigurerGetter func(cfg *LoggerConfig) log.Option

type LogConfigurer struct {
	getters []ConfigurerGetter
}

func (l *LogConfigurer) apply(cfg *LoggerConfig) []log.Option {
	options := []log.Option{}
	for i := range l.getters {
		options = append(options, l.getters[i](cfg))
	}
	return options
}

func InitLogger(
	cfg *LoggerConfig,
) *log.SugaredLogger {

	configurer := LogConfigurer{getters: []ConfigurerGetter{
		logLevelModifier,
	}}

	logger, _ := log.NewProduction(configurer.apply(cfg)...)

	sugarLogger := logger.Sugar()

	return sugarLogger
}

type LoggerConfig struct {
	LogLevel  string `mapstructure:"level"`
	logFormat string `mapstructure:"format"`
}
