package kafka

import "go.uber.org/zap"

type Topic string

type ProducerConfig struct {
	Topic        Topic `mapstructure:"topic"`
	Async        bool  `mapstructure:"async"`
	BatchSize    int   `mapstructure:"batchSize"`
	BatchTimeout int   `mapstructure:"batchTimeout"`
}

type Config struct {
	Addr      []string         `mapstructure:"addr"`
	Producers []ProducerConfig `mapstructure:"producers"`
}

type kafkaLogger struct {
	logger *zap.SugaredLogger
}

func (l *kafkaLogger) Printf(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func newKafkaLogger(logger *zap.SugaredLogger) *kafkaLogger {
	return &kafkaLogger{logger: logger}
}

func Init(cfg *Config, zap *zap.SugaredLogger) (producer *ProducerManager) {
	logger := newKafkaLogger(zap)
	producerManager := InitProducers(cfg, logger)

	return producerManager
}
