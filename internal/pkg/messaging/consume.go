package messaging

type IConsumer interface {
	Consume(
		payload *Message,
		rawMsg interface{},
		messageKey string,
		routingKey RoutingKey,
	) error
}
