package inbox

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/unmei211/notifyme/internal/pkg/messaging"
	"github.com/unmei211/notifyme/internal/pkg/worker"
)

type Processor struct {
	worker.IWorker
	messaging.IFetcher
	repository   Repository
	consumersMap map[messaging.RoutingKey]messaging.IConsumer
	workerId     int
	workerCount  int
}

func (p *Processor) Fetch(ctx context.Context) error {
	inboxes, err := p.repository.FindInboxesForWorker(p.workerId, p.workerCount)

	if err != nil {
		return err
	}

	blacklist := map[string]*MessageInbox{}

	for _, inbox := range inboxes {

		consumer := p.consumersMap[inbox.RoutingKey]

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

		err := consumer.Consume(msg, rawMsg, inbox.MessageKey, inbox.RoutingKey)
		if err != nil {
			blacklist[inbox.MessageKey] = inbox
			return err
		}
	}
}

func (p *Processor) Fallback() {
	//TODO implement me
	panic("implement me")
}

func (p *Processor) Start(ctx context.Context) {

}
func (p *Processor) Stop(ctx context.Context) {

}
