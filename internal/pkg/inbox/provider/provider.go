package inbox

import (
	"github.com/unmei211/notifyme/internal/pkg/inbox"
	"github.com/unmei211/notifyme/internal/pkg/kafka"
	"github.com/unmei211/notifyme/internal/pkg/messaging"
)

func ConsumeHandlerInboxProvider(
	inbox inbox.Inbox,
) kafka.ConsumeHandler {
	return func(msg *messaging.Message, topic kafka.Topic) error {
		err := inbox.Put(msg, string(topic))
		if err != nil {
			return err
		}
		return nil
	}
}
