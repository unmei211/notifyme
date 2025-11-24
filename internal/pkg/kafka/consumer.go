package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type kafkaConsumer struct {
	reader *kafka.Reader
}

type ConsumerManager struct {
	context   context.Context
	logger    *zap.SugaredLogger
	producers map[Topic]Producer
}

func InitConsumers(cfg *Config, log *zap.SugaredLogger, ctx context.Context) (manager *ProducerManager) {
	kafkaLogger := newKafkaLogger(log)

	manager = &ProducerManager{
		producers: map[Topic]Producer{},
		logger:    log,
		context:   ctx,
	}
	for i := range cfg.Producers {
		producerCfg := cfg.Producers[i]
		producer := kafkaProducer{
			writer: &kafka.Writer{
				Addr:         kafka.TCP(cfg.Addr...),
				Topic:        string(producerCfg.Topic),
				Balancer:     &kafka.RoundRobin{},
				MaxAttempts:  100,
				BatchSize:    producerCfg.BatchSize,
				BatchTimeout: time.Duration(producerCfg.BatchTimeout) * time.Millisecond,
				Async:        producerCfg.Async,
				Logger:       kafkaLogger,
				ErrorLogger:  kafkaLogger,
			},
		}
		manager.producers[producerCfg.Topic] = &producer
	}

	return
}
