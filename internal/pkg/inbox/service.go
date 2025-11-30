package inbox

import (
	"time"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
	msg "github.com/unmei211/notifyme/internal/pkg/messaging"
	"go.uber.org/zap"
	"gorm.io/datatypes"
)

type Service struct {
	repo Repository
	log  *zap.SugaredLogger
}

func (s *Service) HandleMessage(
	payload *msg.Message,
	rawMsg interface{},
	messageKey string,
	routingKey msg.RoutingKey,
) error {
	exists, err := s.repo.ExistsByMessageId(payload.MessageId)
	if err != nil {
		s.log.Errorf("Failed check exists message: {%s} in inbox", payload.MessageId)
		return err
	}
	if exists {
		s.log.Debugf("Message {%s} already exists in inbox, skip handle", payload.MessageId)
		return nil
	}

	payloadByte, err := sonic.Marshal(payload.Payload)

	if err != nil {
		s.log.Errorf("Failed serialize message payload")
		return errors.Wrap(err, "Failed serialize message payload")
	}

	rawMessageByte, err := sonic.Marshal(rawMsg)

	if err != nil {
		s.log.Errorf("Failed serialize raw message")
		return errors.Wrap(err, "Failed serialize raw message")
	}

	payloadJson := datatypes.JSON(payloadByte)
	rawMessageJsonb := datatypes.JSON(rawMessageByte)
	_, err = s.repo.Add(
		&MessageInbox{
			MessageId:     payload.MessageId,
			CorrelationId: payload.CorrelationId,
			MessageKey:    messageKey,
			RoutingKey:    routingKey,
			ReceivedAt:    time.Now().UTC(),
			Payload:       payloadJson,
			RawMessage:    rawMessageJsonb,
		},
	)
	if err != nil {
		s.log.Errorf("Can't add message in inbox. Error: {%+v}", err)
		return errors.Wrap(err, "Can't add message in inbox")
	}
	return nil
}

func InitService(repo Repository, log *zap.SugaredLogger) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}
