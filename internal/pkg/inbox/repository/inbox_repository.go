package inbox

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
)

type Repository interface {
	ExistsByMessageId(messageId uuid.UUID) (bool, error)
	Add(msg *MessageInbox) (*MessageInbox, error)
}
type MessageInbox struct {
	MessageId   uuid.UUID `gorm:"primaryKey"`
	RoutingKey  string
	ReceivedAt  time.Time
	ProcessedAt time.Time
	Payload     datatypes.JSON `gorm:"type:jsonb"`
}

func (e *MessageInbox) TableName() string {
	return "inbox"
}
