package messaging

type IConsumer interface {
	Consume(msg *Message, routingKey RoutingKey) error
}
