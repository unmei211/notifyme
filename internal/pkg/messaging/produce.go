package messaging

import (
	"context"

	"go.uber.org/zap"
)

type IProducerManager interface {
	Send(message *Message, routingKey RoutingKey) error
}

type IProducer interface {
	Produce(ctx context.Context, message *Message, logger *zap.SugaredLogger) error
}
