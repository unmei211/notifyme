package inbox

import (
	"github.com/labstack/gommon/log"
	msg "github.com/unmei211/notifyme/internal/pkg/messaging"
	"go.uber.org/zap"
)

type Handler func(payload *msg.Message, rawMsg interface{}, messageKey string, routingKey msg.RoutingKey) error

type Inbox interface {
	Put(payload *msg.Message, rawMsg interface{}, messageKey string, routingKey msg.RoutingKey) error

	handlerRegistrar(handler Handler)
}

type SimpleInbox struct {
	log      *zap.SugaredLogger
	handlers []Handler
}

func (i *SimpleInbox) Put(payload *msg.Message, rawMsg interface{}, messageKey string, routingKey msg.RoutingKey) error {
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
func (i *SimpleInbox) handlerRegistrar(handler Handler) {
	i.handlers = append(i.handlers, handler)
}

func InitInbox(cfg *Config, log *zap.SugaredLogger, handlers []Handler) Inbox {
	inbx := &SimpleInbox{
		log: log,
	}
	for _, handler := range handlers {
		inbx.handlerRegistrar(handler)
	}

	return inbx
}
