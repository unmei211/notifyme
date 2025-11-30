package messaging

import (
	"encoding"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	MessageCommand MessageType = "command"
	MessageEvent   MessageType = "event"
)

type Vendor string

type MessageType string
type RoutingKey string

var _ encoding.TextUnmarshaler = (*RoutingKey)(nil)

func (r *RoutingKey) UnmarshalText(text []byte) error {
	*r = RoutingKey(text)
	return nil
}

type Message struct {
	MessageId     uuid.UUID
	Source        string
	CorrelationId uuid.UUID
	Payload       map[string]interface{}
	Type          MessageType
	Name          string
	CreatedAt     time.Time
}
