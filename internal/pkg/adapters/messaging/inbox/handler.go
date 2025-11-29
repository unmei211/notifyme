package inbox

import (
	ibx "github.com/unmei211/notifyme/internal/pkg/inbox"
	"go.uber.org/zap"
)

func InitHandlers(
	service *ibx.Service,
	log *zap.SugaredLogger,
) []ibx.BoxingHandler {
	var handlers []ibx.BoxingHandler

	handlers = append(handlers, service.HandleMessage)

	return handlers
}
