package kafka

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type kafkaProducer struct {
	writer *kafka.Writer
}

func (p *kafkaProducer) Produce(ctx context.Context) {

}

type Producer interface {
	Produce(ctx context.Context)
}

type ProducerManager struct {
	context   context.Context
	producers map[Topic]Producer
}

func (m *ProducerManager) Send(topic Topic) error {
	producer, exists := m.producers[topic]
	if !exists {
		return errors.Errorf("Not found topic {%s}", topic)
	}

	producer.Produce(m.context)

	return nil
}

func InitProducers(cfg *Config, log *kafkaLogger, ctx context.Context) (manager *ProducerManager) {
	manager = &ProducerManager{
		producers: map[Topic]Producer{},
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
				Logger:       log,
				ErrorLogger:  log,
			},
		}
		manager.producers[producerCfg.Topic] = &producer
	}

	return
}
