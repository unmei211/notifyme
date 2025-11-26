package inbox

import (
	"github.com/labstack/gommon/log"
	msg "github.com/unmei211/notifyme/internal/pkg/messaging"
	"go.uber.org/zap"
)

type PutHandler func(msg *msg.Message, routingKey msg.RoutingKey) error

type Inbox interface {
	Put(msg *msg.Message, routingKey msg.RoutingKey) error

	regPutHandler(handler PutHandler)
}

type SimpleInbox struct {
	log         *zap.SugaredLogger
	putHandlers []PutHandler
}

func (i *SimpleInbox) Put(msg *msg.Message, routingKey msg.RoutingKey) error {
	log.Debugf("Try put message {%s} in inbox", msg.MessageId)
	var err error = nil

	for _, handler := range i.putHandlers {
		err = handler(msg, routingKey)
		if err != nil {
			log.Errorf("Can't handle message {%s} in inbox. Error: {%+v}", msg.MessageId, err)
			return err
		}
	}

	return nil
}
func (i *SimpleInbox) regPutHandler(handler PutHandler) {
	i.putHandlers = append(i.putHandlers, handler)
}

func InitInbox(cfg *Config, log *zap.SugaredLogger, handlers []PutHandler) Inbox {
	inbx := &SimpleInbox{
		log: log,
	}
	for _, handler := range handlers {
		inbx.regPutHandler(handler)
	}

	return inbx
}
