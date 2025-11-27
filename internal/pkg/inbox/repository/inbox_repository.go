package inbox_repository

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
)

type Repository interface {
	ExistsByMessageId(messageId uuid.UUID) (bool, error)
	Add(msg *MessageInbox) (*MessageInbox, error)
	FindInboxesForWorker(workerId int, workersCount int) ([]MessageInbox, error)
}
type MessageInbox struct {
	MessageId     uuid.UUID `gorm:"primaryKey"`
	CorrelationId uuid.UUID
	MessageKey    uuid.UUID
	RoutingKey    string
	ReceivedAt    time.Time
	ProcessedAt   time.Time
	Payload       datatypes.JSON `gorm:"type:jsonb"`
}

func (e *MessageInbox) TableName() string {
	return "inbox"
}
