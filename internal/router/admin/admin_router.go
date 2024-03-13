package adminRouter

import (
	"github.com/gofiber/fiber/v2"
	"latipe-order-service-v2/internal/api/order"
	"latipe-order-service-v2/internal/middleware"
)

type AdminOrderRouter interface {
	Init(root *fiber.Router)
}

type adminOrderRouter struct {
	orderHandler order.OrderApiHandler
	middleware   *middleware.Middleware
}

func NewAdminOrderRouter(orderHandler order.OrderApiHandler, middleware *middleware.Middleware) AdminOrderRouter {
	return &adminOrderRouter{
		orderHandler: orderHandler,
		middleware:   middleware,
	}
}

func (o adminOrderRouter) Init(root *fiber.Router) {

	//admin
	adminRouter := (*root).Group("/admin")
	{
		adminRouter.Get("/:id", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.GetOrderDetailByAdmin)
		adminRouter.Get("", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.ListOfOrder)
		adminRouter.Get("/total/count", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.AdminCountingOrder)
		adminRouter.Patch("/status/cancel", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.AdminCancelOrder)
		adminRouter.Patch("/:id/complete", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.UpdateOrderStatus)
	}

}
