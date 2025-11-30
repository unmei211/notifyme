package kafka

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	msg "github.com/unmei211/notifyme/internal/pkg/messaging"
	"go.uber.org/zap"
)

type kafkaProducer struct {
	writer *kafka.Writer
}

func (p *kafkaProducer) Produce(ctx context.Context, message *msg.Message, logger *zap.SugaredLogger) error {
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

type ProducerManager struct {
	context   context.Context
	logger    *zap.SugaredLogger
	producers map[msg.RoutingKey]msg.IProducer
}

func (m *ProducerManager) Send(message *msg.Message, routingKey msg.RoutingKey) error {
	producer, exists := m.producers[routingKey]
	if !exists {
		m.logger.Errorf("Not found channel {%s}", routingKey)
		return errors.Errorf("Not found channel {%s}", routingKey)
	}

	err := producer.Produce(m.context, message, m.logger)
	if err != nil {
		m.logger.Errorf("Can't produce message {%s} because err {%+v}", message.MessageId, err)
		return err
	}

	m.logger.Debugf("Produce message {%+v}", message)
	return nil
}

func initProducerManager(cfg *Config, log *zap.SugaredLogger, ctx context.Context) (manager msg.IProducerManager) {
	kafkaLogger := NewKafkaLogger(log)

	mng := &ProducerManager{
		producers: map[msg.RoutingKey]msg.IProducer{},
		logger:    log,
		context:   ctx,
	}
	for key, value := range cfg.Producers {
		producerCfg := value
		producer := kafkaProducer{
			writer: &kafka.Writer{
				Addr:         kafka.TCP(cfg.Addr...),
				Topic:        string(value.Topic),
				Balancer:     &kafka.RoundRobin{},
				MaxAttempts:  100,
				BatchSize:    producerCfg.BatchSize,
				BatchTimeout: time.Duration(producerCfg.BatchTimeout) * time.Millisecond,
				Async:        producerCfg.Async,
				Logger:       kafkaLogger,
				ErrorLogger:  kafkaLogger,
			},
		}
		mng.producers[key] = &producer
	}

	return mng
}
