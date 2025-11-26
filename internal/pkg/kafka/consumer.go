package kafka

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/labstack/gommon/log"
	"github.com/segmentio/kafka-go"
	"github.com/unmei211/notifyme/internal/pkg/messaging"
	"github.com/unmei211/notifyme/internal/pkg/worker"
	"go.uber.org/zap"
)

type kafkaConsumer struct {
	logger     *zap.SugaredLogger
	config     *ConsumerConfig
	handler    ConsumeHandler
	reader     *kafka.Reader
	routingKey RoutingKey
	fallback   func(self *kafkaConsumer)
}

type Consumer interface {
	Consume(ctx context.Context)
	Start(ctx context.Context)
	Stop()
	Fallback()
}

type ConsumeHandler func(msg *messaging.Message, topic Topic) error

// TODO: Implement batch commit and zip messages
func (c *kafkaConsumer) Consume(ctx context.Context) {
	rawMsg, err := c.reader.FetchMessage(ctx)

	if err != nil {
		c.logger.Errorf("Can't fetch message")
		//TODO: dead_letters_queue impl
		return
	}

	msg := messaging.Message{}
	err = sonic.Unmarshal(rawMsg.Value, &msg)

	if err != nil {
		c.logger.Errorf("Can't unmarshal message")
		//TODO: dead_letters_queue impl
		err := c.reader.CommitMessages(ctx, rawMsg)
		if err != nil {
			return
		}
		return
	}

	err = c.handler(&msg, c.config.Topic)
	if err != nil {
		c.logger.Errorf("Can't handle message")
		//TODO: may be infinity cycle. Implement dead_letters_queue
		c.Fallback()
		return
	}

	//TODO: in feature create batch commit and zip it
	err = c.reader.CommitMessages(ctx, rawMsg)
	if err != nil {
		c.logger.Errorf("Can't commit message")
		return
	}
}

type ConsumerManager struct {
	cfg       *Config
	context   context.Context
	logger    *zap.SugaredLogger
	consumers map[RoutingKey]Consumer
}

func (c *kafkaConsumer) Fallback() {
	err := c.reader.Close()
	if err != nil {
		return
	}

	c.fallback(c)
}

func InitConsumers(cfg *Config, log *zap.SugaredLogger, handler ConsumeHandler, ctx context.Context) (manager *ConsumerManager) {
	kafkaLogger := newKafkaLogger(log)

	manager = &ConsumerManager{
		cfg:       cfg,
		consumers: map[RoutingKey]Consumer{},
		logger:    log,
		context:   ctx,
	}

	consumerConfigs := cfg.Consume.Consumers
	for key, consumerConfig := range consumerConfigs {

		consumer := kafkaConsumer{
			logger:     log,
			config:     &consumerConfig,
			handler:    handler,
			reader:     createReader(cfg, &consumerConfig, kafkaLogger),
			routingKey: key,
			fallback: func(self *kafkaConsumer) {
				newReader := createReader(cfg, &consumerConfig, kafkaLogger)
				self.reader = newReader
			},
		}

		manager.consumers[key] = &consumer
	}

	return
}

func createReader(cfg *Config, consumerConfig *ConsumerConfig, kafkaLogger *kafkaLogger) *kafka.Reader {
	readerConf := kafka.ReaderConfig{
		Brokers: cfg.Addr,
		Topic:   string(consumerConfig.Topic),
		Logger:  kafkaLogger,
		GroupID: cfg.Consume.GroupId,
	}
	return kafka.NewReader(readerConf)
}

func (c *kafkaConsumer) Start(ctx context.Context) {
	go func() {
		for {
			c.Consume(ctx)
		}
	}()
}
func (c *kafkaConsumer) Stop() {

}

func LaunchConsumers(consumerManager *ConsumerManager, ctx context.Context) {
	log.Infof("Launch consumers")
	log.Debugf("Consumers count: %d", len(consumerManager.consumers))
	var workers []worker.IWorker
	for key, consumer := range consumerManager.consumers {
		log.Debugf("Add consumer for RoutingKey: %s", key)
		workers = append(workers, consumer)
	}

	runner := worker.NewRunner(workers...)

	runner.Launch(ctx)
}
