package handlers

import (
	"github.com/unmei211/notifyme/internal/pkg/event_router"
	"github.com/unmei211/notifyme/internal/pkg/logger"
	msg "github.com/unmei211/notifyme/internal/pkg/messaging"
)

func ConsumeNotificationSentEvent(payload *msg.Message,
	rawMsg interface{},
	messageKey string) error {
	logger.Log.Debugf("Consume notification sent event for message key %s", messageKey)
	return nil
}

var _ event_router.RouteFunc = ConsumeNotificationSentEvent
