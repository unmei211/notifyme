package inbox

import (
	"github.com/unmei211/notifyme/internal/pkg/inbox"
	"github.com/unmei211/notifyme/internal/pkg/messaging"
)

type inboxConsumerProvider struct {
	inbox inbox.Inbox
}

func (i inboxConsumerProvider) Consume(
	msg *messaging.Message,
	routingKey messaging.RoutingKey) error {
	err := i.inbox.Put(msg, routingKey)
	if err != nil {
		return err
	}
	return nil
}

func ConsumeHandlerInboxProvider(
	inbox inbox.Inbox,
) messaging.IConsumer {
	return inboxConsumerProvider{
		inbox: inbox,
	}
}
