package messaging

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type MessageType string

const (
	MessageCommand MessageType = "command"
	MessageEvent   MessageType = "event"
)

type Message struct {
	MessageId     uuid.UUID
	Source        string
	CorrelationId uuid.UUID
	Payload       map[string]interface{}
	Type          MessageType
	Name          string
	CreatedAt     time.Time
}
