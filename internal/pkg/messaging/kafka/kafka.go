package kafka

import (
	"context"

	"github.com/unmei211/notifyme/internal/pkg/messaging"
	"go.uber.org/zap"
)

type Topic string

type kafkaLogger struct {
	logger *zap.SugaredLogger
}

func (l *kafkaLogger) Printf(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func newKafkaLogger(logger *zap.SugaredLogger) *kafkaLogger {
	return &kafkaLogger{logger: logger}
}

func Init(cfg *Config, zap *zap.SugaredLogger, ctx context.Context, consumeHandler messaging.ConsumeHandler) (producerManager messaging.IProducerManager, consumerManager messaging.IConsumerManager) {
	producerManager = initProducerManager(cfg, zap, ctx)
	consumerManager = initConsumerManager(cfg, zap, consumeHandler, ctx)

	return producerManager, consumerManager
}
