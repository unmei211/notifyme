package event_router

import (
	"strings"

	msg "github.com/unmei211/notifyme/internal/pkg/messaging"
	"go.uber.org/zap"
)

type RouteFunc func(payload *msg.Message,
	rawMsg interface{},
	messageKey string) error

type IRouter interface {
	Route(payload *msg.Message,
		rawMsg interface{},
		messageKey string,
		routingKey msg.RoutingKey) error
	Register(key msg.RoutingKey, routeFunc RouteFunc)
	RegisterString(key string, routeFunc RouteFunc)
	RegisterAll(keyToFunc map[msg.RoutingKey]RouteFunc)
	msg.IConsumer
}

type SimpleRouter struct {
	routes map[msg.RoutingKey]RouteFunc
	log    *zap.SugaredLogger
}

var _ IRouter = (*SimpleRouter)(nil)

func (s *SimpleRouter) Route(
	payload *msg.Message,
	rawMsg interface{},
	messageKey string,
	routingKey msg.RoutingKey) error {

	routeFunc, exists := s.routes[routingKey]

	if !exists {
		s.log.Fatalf("Fatal error. Has not exists RouterFunc for routingKey: %s", routingKey)
	}

	err := routeFunc(payload, rawMsg, messageKey)

	return err
}

func (s *SimpleRouter) Register(key msg.RoutingKey, routeFunc RouteFunc) {
	s.routes[key] = routeFunc
}

func (s *SimpleRouter) normalizeStr(key string) msg.RoutingKey {
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, ".", "_")

	return msg.RoutingKey(key)
}

func (s *SimpleRouter) RegisterString(key string, routeFunc RouteFunc) {
	s.routes[s.normalizeStr(key)] = routeFunc
}

func (s *SimpleRouter) RegisterAll(keyToFunc map[msg.RoutingKey]RouteFunc) {
	for key, routeFunc := range keyToFunc {
		s.Register(key, routeFunc)
	}
}

func (s *SimpleRouter) Consume(
	payload *msg.Message,
	rawMsg interface{},
	messageKey string,
	routingKey msg.RoutingKey,
) error {
	return s.Route(payload, rawMsg, messageKey, routingKey)
}

func Init(log *zap.SugaredLogger) IRouter {
	router := &SimpleRouter{
		routes: map[msg.RoutingKey]RouteFunc{},
		log:    log,
	}

	return router
}
