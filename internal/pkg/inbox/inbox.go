package inbox

import (
	"github.com/labstack/gommon/log"
	msg "github.com/unmei211/notifyme/internal/pkg/messaging"
	"go.uber.org/zap"
)

type BoxingHandler func(payload *msg.Message, rawMsg interface{}, messageKey string, routingKey msg.RoutingKey) error

type MessageBoxing interface {
	Box(payload *msg.Message, rawMsg interface{}, messageKey string, routingKey msg.RoutingKey) error

	handlerRegistrar(handler BoxingHandler)
}

type SimpleMessageBoxing struct {
	log      *zap.SugaredLogger
	handlers []BoxingHandler
}

func (i *SimpleMessageBoxing) Box(payload *msg.Message, rawMsg interface{}, messageKey string, routingKey msg.RoutingKey) error {
	log.Debugf("Try put message {%s} in inbox", payload.MessageId)
	var err error = nil

	for _, handler := range i.handlers {
		err = handler(payload, rawMsg, messageKey, routingKey)
		if err != nil {
			log.Errorf("Can't handle message {%s} in inbox. Error: {%+v}", payload.MessageId, err)
			return err
		}
	}

	return nil
}
func (i *SimpleMessageBoxing) handlerRegistrar(handler BoxingHandler) {
	i.handlers = append(i.handlers, handler)
}

func InitMessageBoxing(cfg *Config, log *zap.SugaredLogger, handlers []BoxingHandler) MessageBoxing {
	inbx := &SimpleMessageBoxing{
		log: log,
	}

	for _, handler := range handlers {
		inbx.handlerRegistrar(handler)
	}

	return inbx
}

type UnboxingHandler func(payload *msg.Message, rawMsg interface{}, messageKey string, routingKey msg.RoutingKey) error

type MessageUnboxing interface {
	Unbox(payload *msg.Message, rawMsg interface{}, messageKey string, routingKey msg.RoutingKey) error
}

type SimpleMessageUnboxing struct {
	log      *zap.SugaredLogger
	handlers []UnboxingHandler
}

func (u *SimpleMessageUnboxing) Unbox(payload *msg.Message, rawMsg interface{}, messageKey string, routingKey msg.RoutingKey) error {
	return nil
}
