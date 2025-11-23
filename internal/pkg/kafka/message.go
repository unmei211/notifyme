package kafka

import "time"

type MessageType string

const (
	MessageCommand MessageType = "command"
	MessageEvent   MessageType = "event"
)

type Message struct {
	MessageId     string
	Source        string
	CorrelationId string
	Payload       interface{}
	Type          MessageType
	Name          string
	CreatedAt     time.Time
}
