package inbox

import (
	"github.com/unmei211/notifyme/internal/pkg/inbox"
	"github.com/unmei211/notifyme/internal/pkg/messaging"
)

func ConsumeHandlerInboxProvider(
	inbox inbox.Inbox,
) messaging.ConsumeHandler {
	return func(msg *messaging.Message, routingKey messaging.RoutingKey) error {
		err := inbox.Put(msg, routingKey)
		if err != nil {
			return err
		}
		return nil
	}
}
