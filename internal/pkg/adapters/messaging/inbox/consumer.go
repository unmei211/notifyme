package inbox

import (
	"github.com/unmei211/notifyme/internal/pkg/inbox"
	"github.com/unmei211/notifyme/internal/pkg/messaging"
)

type inboxConsumerAdapter struct {
	inbox inbox.MessageBoxing
}

func (i inboxConsumerAdapter) Consume(
	payload *messaging.Message, rawMsg interface{}, messageKey string, routingKey messaging.RoutingKey,
) error {
	err := i.inbox.Box(payload, rawMsg, messageKey, routingKey)
	if err != nil {
		return err
	}
	return nil
}

func InitConsumer(
	inbox inbox.MessageBoxing,
) messaging.IConsumer {
	return inboxConsumerAdapter{
		inbox: inbox,
	}
}
