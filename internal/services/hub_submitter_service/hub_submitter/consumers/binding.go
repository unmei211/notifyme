package consumers

import (
	eventrouter "github.com/unmei211/notifyme/internal/pkg/event_router"
	"github.com/unmei211/notifyme/internal/services/hub_submitter/hub_submitter/consumers/handlers"
)

func Bind(router eventrouter.IRouter) {
	router.RegisterString("notificationSender_sentEvent", handlers.ConsumeNotificationSentEvent)
}
