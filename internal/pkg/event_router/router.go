package event_router

import msg "github.com/unmei211/notifyme/internal/pkg/messaging"

type RouteFunc func(payload *msg.Message,
	rawMsg interface{},
	messageKey string) error

type IRouter interface {
	msg.IConsumer
	Route(payload *msg.Message,
		rawMsg interface{},
		messageKey string,
		routingKey msg.RoutingKey) error
	Register(key msg.RoutingKey, routeFunc RouteFunc)
	RegisterAll(keyToFunc map[msg.RoutingKey]RouteFunc)
}

type SimpleRouter struct {
	routes map[msg.RoutingKey]RouteFunc
}

var _ IRouter = (*SimpleRouter)(nil)

func (s *SimpleRouter) Route(
	payload *msg.Message,
	rawMsg interface{},
	messageKey string,
	routingKey msg.RoutingKey) error {

	routeFunc := s.routes[routingKey]
	err := routeFunc(payload, rawMsg, messageKey)

	return err
}

func (s *SimpleRouter) Register(key msg.RoutingKey, routeFunc RouteFunc) {
	s.routes[key] = routeFunc
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

func Init() IRouter {
	router := &SimpleRouter{
		routes: map[msg.RoutingKey]RouteFunc{},
	}

	return router
}
