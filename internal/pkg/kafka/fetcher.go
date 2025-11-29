package kafka

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/segmentio/kafka-go"
	"github.com/unmei211/notifyme/internal/pkg/messaging"
	"github.com/unmei211/notifyme/internal/pkg/worker"
	"go.uber.org/zap"
)

type Fetcher struct {
	consumer   messaging.IConsumer
	reader     *kafka.Reader
	log        *zap.SugaredLogger
	routingKey messaging.RoutingKey
}

func (f *Fetcher) Fallback() {
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	// TODO: If an error occurs during processing, the message will not be sent again.
	rawMsg, err := f.reader.FetchMessage(ctx)

	if err != nil {
		f.log.Errorf("Can't fetch message")
		//TODO: dead_letters_queue impl
		return nil
	}
	
	msg := messaging.Message{}
	err = sonic.Unmarshal(rawMsg.Value, &msg)

	if err != nil {
		f.log.Errorf("Can't unmarshal message")
		//TODO: dead_letters_queue impl
		err := f.reader.CommitMessages(ctx, rawMsg)
		if err != nil {
			return err
		}
		return err
	}

	err = f.consumer.Consume(&msg, &rawMsg, string(rawMsg.Key[:]), f.routingKey)
	if err != nil {
		f.log.Errorf("Can't handle message")
		//TODO: may be infinity cycle. Implement dead_letters_queue
		f.Fallback()
		return err
	}

	//TODO: in feature create batch commit and zip it
	err = f.reader.CommitMessages(ctx, rawMsg)
	if err != nil {
		f.log.Errorf("Can't commit message")
		return err
	}
	return nil
}

func (f *Fetcher) Start(ctx context.Context) {
	go func() {
		for {
			err := f.Fetch(ctx)
			if err != nil {
				return
			}
		}
	}()
}
func (f *Fetcher) Stop() {

}

type Manager struct {
	fetchingCfg    *FetchingConfig
	routingCfg     *messaging.RoutingConfig
	context        context.Context
	log            *zap.SugaredLogger
	routeToFetcher map[messaging.RoutingKey]*Fetcher
}

func InitFetcher(kafkaConfig *Config,
	routingConfig *messaging.RoutingConfig,
	log *zap.SugaredLogger,
	consumer messaging.IConsumer,
	ctx context.Context) messaging.IFetcherManager {

	kafkaLogger := NewKafkaLogger(log)

	mng := &Manager{
		fetchingCfg:    &kafkaConfig.Fetching,
		routingCfg:     routingConfig,
		context:        ctx,
		log:            log,
		routeToFetcher: map[messaging.RoutingKey]*Fetcher{},
	}

	for routingKey, inputConfig := range routingConfig.Input {

		fetcher := &Fetcher{
			consumer: consumer,
			reader: kafka.NewReader(kafka.ReaderConfig{
				Brokers: kafkaConfig.Addr,
				Topic:   inputConfig.VendorKey,
				Logger:  kafkaLogger,
				GroupID: kafkaConfig.Fetching.GroupId,
				MaxWait: 60 * time.Second,
			}),
			log:        log,
			routingKey: routingKey,
		}

		mng.routeToFetcher[routingKey] = fetcher
	}

	return mng
}

func (m *Manager) Launch(ctx context.Context) {
	m.log.Infof("Launch Fetchers")
	m.log.Debugf("Fetchers count: %d", len(m.routeToFetcher))
	var workers []worker.IWorker
	for routingKey, fetcher := range m.routeToFetcher {
		m.log.Debugf("Add fetcher for RoutingKey: %s", routingKey)
		workers = append(workers, fetcher)
	}

	runner := worker.NewRunner(workers...)

	runner.Launch(ctx)
}

func LaunchFetcher(manager messaging.IFetcherManager, ctx context.Context) {
	println("wtf")
	manager.Launch(ctx)
}
