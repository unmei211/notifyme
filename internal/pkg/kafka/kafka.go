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

type ConsumerConfig struct {
	Topic Topic `mapstructure:"topic"`
}

type ConsumeConfig struct {
	Consumers []ConsumerConfig `mapstructure:"consumers"`
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

func Init(cfg *Config, zap *zap.SugaredLogger, ctx context.Context, consumeHandler ConsumeHandler) (producer *ProducerManager, manager *ConsumerManager) {
	producerManager := InitProducers(cfg, zap, ctx)
	consumerManager := InitConsumers(cfg, zap, consumeHandler, ctx)

	return producerManager, consumerManager
}
