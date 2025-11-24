package kafka

import (
	"context"

	"go.uber.org/zap"
)

type Topic string

type ProducerConfig struct {
	Topic        Topic `mapstructure:"topic"`
	Async        bool  `mapstructure:"async"`
	BatchSize    int   `mapstructure:"batchSize"`
	BatchTimeout int   `mapstructure:"batchTimeout"`
}

type WorkerConfig struct {
	Workers int         `mapstructure:"workers"`
	Group   WorkerGroup `mapstructure:"group"`
}

type WorkerGroup string

type ConsumerConfig struct {
	Topic Topic
}

type ConsumeConfig struct {
	Workers   map[WorkerGroup]WorkerConfig `mapstructure:"workers"`
	Consumers []ConsumerConfig             `mapstructure:"consumers"`
}

type Config struct {
	Addr      []string         `mapstructure:"addr"`
	Producers []ProducerConfig `mapstructure:"producers"`
	Consume   ConsumeConfig    `mapstructure:"consume"`
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

func Init(cfg *Config, zap *zap.SugaredLogger, ctx context.Context) (producer *ProducerManager) {
	producerManager := InitProducers(cfg, zap, ctx)

	return producerManager
}
