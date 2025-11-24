package inbox

import (
	"time"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
	"github.com/unmei211/notifyme/internal/pkg/inbox"
	"github.com/unmei211/notifyme/internal/pkg/messaging"
	"go.uber.org/zap"
	"gorm.io/datatypes"
)

func PutHandlerInboxProvider(
	inboxRepository Repository,
	log *zap.SugaredLogger,
) []inbox.PutHandler {
	var handlers []inbox.PutHandler

	handlers = append(handlers, func(msg *messaging.Message, topic string) error {
		exists, err := inboxRepository.ExistsByMessageId(msg.MessageId)
		if err != nil {
			log.Errorf("Failed check exists message: {%s} in inbox", msg.MessageId)
			return err
		}
		if exists {
			log.Debugf("Message {%s} already exists in inbox, skip handle", msg.MessageId)
			return nil
		}

		payloadByte, err := sonic.Marshal(msg.Payload)

		if err != nil {
			log.Errorf("Failed serialize message payload")
			return errors.Wrap(err, "Failed serialize message payload")
		}

		databaseJSON := datatypes.JSON(payloadByte)

		_, err = inboxRepository.Add(
			&MessageInbox{
				MessageId:  msg.MessageId,
				Topic:      topic,
				ReceivedAt: time.Now().UTC(),
				Payload:    databaseJSON,
			},
		)
		if err != nil {
			log.Errorf("Can't add message in inbox. Error: {%+v}", err)
			return errors.Wrap(err, "Can't add message in inbox")
		}
		return nil
	})

	return handlers
}
