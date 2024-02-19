package internalRouter

import (
	"github.com/gofiber/fiber/v2"
	"latipe-order-service-v2/internal/api/order"
	"latipe-order-service-v2/internal/middleware"
)

type InternalOrderRouter interface {
	Init(root *fiber.Router)
}

type internalOrderRouter struct {
	orderHandler     order.OrderApiHandler
	statisticHandler order.OrderStatisticApiHandler
	middleware       *middleware.Middleware
}

func NewInternalOrderRouter(orderHandler order.OrderApiHandler, statisticHandler order.OrderStatisticApiHandler,
	middleware *middleware.Middleware) InternalOrderRouter {
	return &internalOrderRouter{
		orderHandler:     orderHandler,
		statisticHandler: statisticHandler,
		middleware:       middleware,
	}
}

func (o internalOrderRouter) Init(root *fiber.Router) {
	//internal
	internalRouter := (*root).Group("/internal")
	{
		internalRouter.Get("/rating/:id", o.middleware.Authentication.RequiredInternalService(), o.orderHandler.InternalGetOrderByOrderID)
	}

}
