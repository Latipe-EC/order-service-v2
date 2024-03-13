package deliveryRouter

import (
	"github.com/gofiber/fiber/v2"
	"latipe-order-service-v2/internal/api/order"
	"latipe-order-service-v2/internal/middleware"
)

type DeliveryOrderRouter interface {
	Init(root *fiber.Router)
}

type deliveryOrderRouter struct {
	orderHandler order.OrderApiHandler
	middleware   *middleware.Middleware
}

func NewDeliveryOrderRouter(orderHandler order.OrderApiHandler,
	middleware *middleware.Middleware) DeliveryOrderRouter {
	return &deliveryOrderRouter{
		orderHandler: orderHandler,
		middleware:   middleware,
	}
}

func (o deliveryOrderRouter) Init(root *fiber.Router) {

	//delivery
	deliveryRouter := (*root).Group("/delivery")
	{
		deliveryRouter.Get("", o.middleware.Authentication.RequiredDeliveryAuthentication(), o.orderHandler.GetOrdersByDelivery)
		deliveryRouter.Get("/:id", o.middleware.Authentication.RequiredDeliveryAuthentication(), o.orderHandler.GetOrderDetailByDelivery)
		deliveryRouter.Get("/total/count", o.middleware.Authentication.RequiredDeliveryAuthentication(), o.orderHandler.DeliveryCountingOrder)
		deliveryRouter.Patch("/:id", o.middleware.Authentication.RequiredDeliveryAuthentication(), o.orderHandler.UpdateOrderStatusByDelivery)
	}

}
