package messaging

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	MessageCommand MessageType = "command"
	MessageEvent   MessageType = "event"
)

type MessageType string
type RoutingKey string

type Message struct {
	MessageId     uuid.UUID
	Source        string
	CorrelationId uuid.UUID
	Payload       map[string]interface{}
	Type          MessageType
	Name          string
	CreatedAt     time.Time
}
