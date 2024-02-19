package storeRouter

import (
	"github.com/gofiber/fiber/v2"
	"latipe-order-service-v2/internal/api/order"
	"latipe-order-service-v2/internal/middleware"
)

type StoreOrderRouter interface {
	Init(root *fiber.Router)
}

type storeOrderRouter struct {
	orderHandler order.OrderApiHandler
	middleware   *middleware.Middleware
}

func NewStoreOrderRouter(orderHandler order.OrderApiHandler,
	middleware *middleware.Middleware) StoreOrderRouter {
	return &storeOrderRouter{
		orderHandler: orderHandler,

		middleware: middleware,
	}
}

func (o storeOrderRouter) Init(root *fiber.Router) {
	//store
	storeRouter := (*root).Group("/store")
	{
		storeRouter.Get("", o.middleware.Authentication.RequiredStoreAuthentication(), o.orderHandler.GetMyStoreOrder)
		storeRouter.Get("/search", o.middleware.Authentication.RequiredStoreAuthentication(), o.orderHandler.SearchOrderIdByKeyword)
		storeRouter.Get("/total/count", o.middleware.Authentication.RequiredStoreAuthentication(), o.orderHandler.StoreCountingOrder)
		storeRouter.Get("/:id", o.middleware.Authentication.RequiredStoreAuthentication(), o.orderHandler.GetStoreOrderDetail)
		storeRouter.Patch("/:id/status", o.middleware.Authentication.RequiredStoreAuthentication(), o.orderHandler.UpdateOrderStatusByStore)
	}
}
