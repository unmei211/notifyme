package kafka

import (
	"context"

	"github.com/unmei211/notifyme/internal/pkg/messaging"
	"go.uber.org/zap"
)

type Topic string

type Logger struct {
	logger *zap.SugaredLogger
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func NewKafkaLogger(logger *zap.SugaredLogger) *Logger {
	return &Logger{logger: logger}
}

func Init(
	cfg *Config,
	messagingConfig *messaging.Config,
	zap *zap.SugaredLogger,
	consumer messaging.IConsumer,
	ctx context.Context,
) (producerManager messaging.IProducerManager, fetchingManager messaging.IFetcherManager) {
	routingConfig := messagingConfig.Routing[messaging.Kafka]

	producerManager = initProducerManager(cfg, zap, ctx)
	fetchingManager = InitFetcher(cfg, &routingConfig, zap, consumer, ctx)

	return
}
