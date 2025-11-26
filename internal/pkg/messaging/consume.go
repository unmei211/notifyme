package messaging

import (
	"context"
)

type IConsumer interface {
	Consume(ctx context.Context)
	Start(ctx context.Context)
	Stop()
	Fallback()
}
type IConsumerManager interface {
	Launch(ctx context.Context)
}

type ConsumeHandler func(msg *Message, routingKey RoutingKey) error
