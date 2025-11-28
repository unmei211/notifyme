package inbox

import (
	"github.com/unmei211/notifyme/internal/pkg/inbox"
	"github.com/unmei211/notifyme/internal/pkg/messaging"
)

type inboxConsumerProvider struct {
	inbox inbox.Inbox
}

func (i inboxConsumerProvider) Consume(
	payload *messaging.Message, rawMsg interface{}, messageKey string, routingKey messaging.RoutingKey,
) error {
	err := i.inbox.Put(payload, rawMsg, messageKey, routingKey)
	if err != nil {
		return err
	}
	return nil
}

func InboxConsumerProvider(
	inbox inbox.Inbox,
) messaging.IConsumer {
	return inboxConsumerProvider{
		inbox: inbox,
	}
}
