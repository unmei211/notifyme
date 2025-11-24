package kafka

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"github.com/unmei211/notifyme/internal/pkg/messaging"
	"go.uber.org/zap"
)

type kafkaProducer struct {
	writer *kafka.Writer
}

func (p *kafkaProducer) Produce(ctx context.Context, message *messaging.Message, logger *zap.SugaredLogger) error {
	raw, err := sonic.Marshal(message)

	if err != nil {
		logger.Errorf("Can't marshall message {%s}", message.MessageId)
		return err
	}
	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   nil,
		Value: raw,
	})
	if err != nil {
		logger.Errorf("Can't write message {%s}", message.MessageId)
		return err
	}

	return nil
}

type Producer interface {
	Produce(ctx context.Context, message *messaging.Message, logger *zap.SugaredLogger) error
}

type ProducerManager struct {
	context   context.Context
	logger    *zap.SugaredLogger
	producers map[Topic]Producer
}

func (m *ProducerManager) Send(message *messaging.Message, topic Topic) error {
	producer, exists := m.producers[topic]
	if !exists {
		m.logger.Errorf("Not found topic {%s}", topic)
		return errors.Errorf("Not found topic {%s}", topic)
	}

	err := producer.Produce(m.context, message, m.logger)
	if err != nil {
		m.logger.Errorf("Can't produce message {%s} because err {%+v}", message.MessageId, err)
		return err
	}

	m.logger.Debugf("Produce message {%+v}", message)
	return nil
}

func InitProducers(cfg *Config, log *zap.SugaredLogger, ctx context.Context) (manager *ProducerManager) {
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
