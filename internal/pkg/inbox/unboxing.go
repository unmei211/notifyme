package inbox

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/labstack/gommon/log"
	"github.com/unmei211/notifyme/internal/pkg/messaging"
	"github.com/unmei211/notifyme/internal/pkg/worker"
	"go.uber.org/zap"
)

type MessageUnbox struct {
	cfg        *Config
	processors []*Processor
	log        *zap.SugaredLogger
}

var _ messaging.IFetcherManager = (*MessageUnbox)(nil)

func (m MessageUnbox) Launch(ctx context.Context) {
	log.Infof("Startup message unbox with %d processors", len(m.processors))
	for _, processor := range m.processors {
		processor.Start(ctx)
	}
}

func InitMessageUnbox(
	repository Repository,
	consumer messaging.IConsumer,
	cfg *Config,
	log *zap.SugaredLogger,
) *MessageUnbox {
	var processors []*Processor
	manager := &MessageUnbox{
		cfg: cfg,
		log: log,
	}
	log.Debugf("Start initialize %d workers", cfg.Unbox.MaxWorkers)
	for i := 0; i < cfg.Unbox.MaxWorkers; i++ {
		processors = append(processors,
			&Processor{
				repository:  repository,
				consumer:    consumer,
				workerId:    i,
				workerCount: cfg.Unbox.MaxWorkers,
				log:         log,
			})
	}
	manager.processors = processors
	return manager
}

type Processor struct {
	log         *zap.SugaredLogger
	repository  Repository
	consumer    messaging.IConsumer
	workerId    int
	workerCount int
}

var _ worker.IWorker = (*Processor)(nil)
var _ messaging.IFetcher = (*Processor)(nil)

func (p *Processor) Fetch(ctx context.Context) error {
	log.Debugf("Regular fetch inboxes for worker: %d", p.workerId)
	// TODO: pagination
	inboxes, err := p.repository.FindInboxesForWorker(p.workerId, p.workerCount)
	log.Debugf("Founded %d inbox messages for worker %d", len(inboxes), p.workerId)

	if err != nil {
		log.Errorf("Failed inbox messages from repository for worker: %d", p.workerId)
		return err
	}

	blacklist := map[string]*MessageInbox{}

	for _, inbox := range inboxes {

		msg := &messaging.Message{}
		err = sonic.Unmarshal(inbox.Payload, msg)

		if err != nil {
			blacklist[inbox.MessageKey] = inbox
			continue
		}

		var rawMsg interface{}
		err = sonic.Unmarshal(inbox.RawMessage, rawMsg)

		if err != nil {
			blacklist[inbox.MessageKey] = inbox
			continue
		}

		err := p.consumer.Consume(msg, rawMsg, inbox.MessageKey, inbox.RoutingKey)
		if err != nil {
			blacklist[inbox.MessageKey] = inbox
			return err
		}
	}
	return nil
}

func (p *Processor) Fallback() {
	//TODO implement me
	panic("implement me")
}

func (p *Processor) Start(ctx context.Context) {
	go func() {
		for {
			err := p.Fetch(ctx)
			if err != nil {
				return
			}
		}
	}()
}
func (p *Processor) Stop() {}
