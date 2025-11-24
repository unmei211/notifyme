package inbox

import (
	"github.com/labstack/gommon/log"
	"github.com/unmei211/notifyme/internal/pkg/messaging"
	"go.uber.org/zap"
)

type Inbox interface {
	Put(msg *messaging.Message) error

	regPutHandler(handler PutHandler)
}

type PutHandler func(msg *messaging.Message) error

type inbox struct {
	log         *zap.SugaredLogger
	putHandlers []PutHandler
}

func (i *inbox) Put(msg *messaging.Message) error {
	log.Debugf("Try put message {%s} in inbox", msg.MessageId)
	var err error = nil

	for _, handler := range i.putHandlers {
		err = handler(msg)
		if err != nil {
			log.Errorf("Can't handle message {%s} in inbox. Error: {%+v}", msg.MessageId, err)
			return err
		}
	}

	return nil
}
func (i *inbox) regPutHandler(handler PutHandler) {
	i.putHandlers = append(i.putHandlers, handler)
}

func InitInbox(cfg *Config, log *zap.SugaredLogger) Inbox {
	return &inbox{
		log: log,
	}
}

func RegisterPutHandlers(inbx Inbox, handlers ...PutHandler) {
	for _, handler := range handlers {
		inbx.regPutHandler(handler)
	}
}
