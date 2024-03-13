package userRouter

import (
	"github.com/gofiber/fiber/v2"
	"latipe-order-service-v2/internal/api/order"
	"latipe-order-service-v2/internal/middleware"
	"latipe-order-service-v2/internal/middleware/auth"
)

type UserOrderRouter interface {
	Init(root *fiber.Router)
}

type userOrderRouter struct {
	orderHandler order.OrderApiHandler

	middleware *middleware.Middleware
}

func NewUserOrderRouter(orderHandler order.OrderApiHandler,
	middleware *middleware.Middleware) UserOrderRouter {
	return &userOrderRouter{
		orderHandler: orderHandler,
		middleware:   middleware,
	}
}

func (o userOrderRouter) Init(root *fiber.Router) {
	//user
	userRouter := (*root).Group("/user")
	{
		userRouter.Get("", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.GetMyOrder)
		userRouter.Get("/total/count", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.UserCountingOrder)
		userRouter.Get("/:id", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.GetOrderDetailOfUser)
		userRouter.Post("", o.middleware.Authentication.RequiredRole([]string{auth.ROLE_USER, auth.ROLE_STORE}), o.orderHandler.CreateOrder)
		userRouter.Patch("/cancel", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.UserCancelOrder)
		userRouter.Patch("/refund", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.UserCancelOrder)
	}
}
