package inbox_processor

import (
	"context"

	uuid "github.com/satori/go.uuid"
	"github.com/unmei211/notifyme/internal/pkg/inbox"
	ibxRepository "github.com/unmei211/notifyme/internal/pkg/inbox/repository"
	"github.com/unmei211/notifyme/internal/pkg/messaging"
	"go.uber.org/zap"
)

type inboxConsumer struct {
	logger          *zap.SugaredLogger
	handler         messaging.ConsumeHandler
	inboxRepository ibxRepository.Repository
	workerId        int
	workerCount     int
}

func (c *inboxConsumer) Start(ctx context.Context) {
	go func() {
		for {
			c.Consume(ctx)
		}
	}()
}
func (c *inboxConsumer) Stop() {

}

func (c *inboxConsumer) Fallback() {

}

func (c *inboxConsumer) Consume(ctx context.Context) {
	inboxes, err := c.inboxRepository.FindInboxesForWorker(c.workerId, c.workerCount)

	blacklist := map[uuid.UUID]*ibxRepository.MessageInbox{}
	if err != nil {
		// TODO
	}

	for _, messageInbox := range inboxes {
		if _, exists := blacklist[messageInbox.MessageKey]; exists {
			// TODO: log
			continue
		}
		// TODO:
		//c.handler()
		if err != nil {
			//todo: log
			blacklist[messageInbox.MessageKey] = &messageInbox
			continue
		}
	}

	//rawMsg, err := c.reader.FetchMessage(ctx)
	//
	//if err != nil {
	//	c.logger.Errorf("Can't fetch message")
	//	//TODO: dead_letters_queue impl
	//	return
	//}
	//
	//msg := messaging.Message{}
	//err = sonic.Unmarshal(rawMsg.Value, &msg)
	//
	//if err != nil {
	//	c.logger.Errorf("Can't unmarshal message")
	//	//TODO: dead_letters_queue impl
	//	err := c.reader.CommitMessages(ctx, rawMsg)
	//	if err != nil {
	//		return
	//	}
	//	return
	//}
	//
	//err = c.handler(&msg, c.routingKey)
	//if err != nil {
	//	c.logger.Errorf("Can't handle message")
	//	//TODO: may be infinity cycle. Implement dead_letters_queue
	//	c.Fallback()
	//	return
	//}
	//
	////TODO: in feature create batch commit and zip it
	//err = c.reader.CommitMessages(ctx, rawMsg)
	//if err != nil {
	//	c.logger.Errorf("Can't commit message")
	//	return
	//}
}

type ConsumerManager struct {
	cfg       *inbox.Config
	context   context.Context
	logger    *zap.SugaredLogger
	consumers map[messaging.RoutingKey]messaging.IConsumer
}
