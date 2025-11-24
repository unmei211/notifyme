package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/unmei211/notifyme/internal/pkg/messaging"
	"go.uber.org/zap"
)

type kafkaConsumer struct {
	logger  *zap.SugaredLogger
	config  *ConsumerConfig
	handler ConsumeHandler
	reader  *kafka.Reader
}

type Consumer interface {
	Consume(msg *messaging.Message)
}

type ConsumeHandler func(msg *messaging.Message, topic Topic) error

func (c *kafkaConsumer) Consume(msg *messaging.Message) {

}

type ConsumerManager struct {
	context   context.Context
	logger    *zap.SugaredLogger
	consumers map[Topic]Consumer
}

func InitConsumers(cfg *Config, log *zap.SugaredLogger, handler ConsumeHandler, ctx context.Context) (manager *ConsumerManager) {
	kafkaLogger := newKafkaLogger(log)

	manager = &ConsumerManager{
		consumers: map[Topic]Consumer{},
		logger:    log,
		context:   ctx,
	}

	consumerConfigs := cfg.Consume.Consumers
	for _, consumerConfig := range consumerConfigs {

		readerConf := kafka.ReaderConfig{
			Brokers: cfg.Addr,
			Topic:   string(consumerConfig.Topic),
			Logger:  kafkaLogger,
		}
		consumer := kafkaConsumer{
			logger:  log,
			config:  &consumerConfig,
			handler: handler,
			reader:  kafka.NewReader(readerConf),
		}
		manager.consumers[consumerConfig.Topic] = &consumer
	}

	return
}
