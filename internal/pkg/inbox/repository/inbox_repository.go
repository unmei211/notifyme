package inbox

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/unmei211/notifyme/internal/pkg/messaging"
)

type Repository interface {
	ExistsByMessageId(messageId string) (bool, error)
	Add(msg *MessageInbox) error
}

type MessageInbox struct {
	MessageId   uuid.UUID
	Topic       string
	ReceivedAt  time.Time
	ProcessedAt time.Time
	Payload     messaging.Message
}
