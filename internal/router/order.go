package router

import (
	"github.com/gofiber/fiber/v2"
	"latipe-order-service-v2/internal/api/order"
	"latipe-order-service-v2/internal/middleware"
	"latipe-order-service-v2/internal/middleware/auth"
)

type OrderRouter interface {
	Init(root *fiber.Router)
}

type orderRouter struct {
	orderHandler     order.OrderApiHandler
	statisticHandler order.OrderStatisticApiHandler
	middleware       *middleware.Middleware
}

func NewOrderRouter(orderHandler order.OrderApiHandler, statisticHandler order.OrderStatisticApiHandler, middleware *middleware.Middleware) OrderRouter {
	return orderRouter{
		orderHandler:     orderHandler,
		statisticHandler: statisticHandler,
		middleware:       middleware,
	}
}

func (o orderRouter) Init(root *fiber.Router) {
	orderRouter := (*root).Group("/orders")
	{
		//admin
		adminRouter := orderRouter.Group("/admin")
		{
			adminRouter.Get("/:id", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.GetByOrderId)
			adminRouter.Get("", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.ListOfOrder)
			adminRouter.Get("/total/count", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.AdminCountingOrder)
			adminRouter.Patch("/status/cancel", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.AdminCancelOrder)
			adminRouter.Patch("/:id/complete", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.UpdateOrderStatus)
		}

		//user
		userRouter := orderRouter.Group("/user")
		{
			userRouter.Get("", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.GetMyOrder)
			userRouter.Get("/total/count", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.UserCountingOrder)
			userRouter.Get("/:id", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.UserGetOrderByID)
			userRouter.Post("", o.middleware.Authentication.RequiredRole([]string{auth.ROLE_USER}), o.orderHandler.CreateOrder)
			userRouter.Patch("/cancel", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.UserCancelOrder)
			userRouter.Patch("/refund", o.middleware.Authentication.RequiredAuthentication(), o.orderHandler.UserCancelOrder)
		}

		//store
		storeRouter := orderRouter.Group("/store")
		{
			storeRouter.Get("", o.middleware.Authentication.RequiredStoreAuthentication(), o.orderHandler.GetMyStoreOrder)
			storeRouter.Get("/search", o.middleware.Authentication.RequiredStoreAuthentication(), o.orderHandler.SearchOrderIdByKeyword)
			storeRouter.Get("/total/count", o.middleware.Authentication.RequiredStoreAuthentication(), o.orderHandler.StoreCountingOrder)
			storeRouter.Get("/:id", o.middleware.Authentication.RequiredStoreAuthentication(), o.orderHandler.GetStoreOrderDetail)
			storeRouter.Patch("/:id/status", o.middleware.Authentication.RequiredStoreAuthentication(), o.orderHandler.UpdateOrderStatusByStore)
		}

		//delivery
		deliveryRouter := orderRouter.Group("/delivery")
		{
			deliveryRouter.Get("", o.middleware.Authentication.RequiredDeliveryAuthentication(), o.orderHandler.GetOrdersByDelivery)
			deliveryRouter.Get("/:id", o.middleware.Authentication.RequiredDeliveryAuthentication(), o.orderHandler.DeliveryGetOrderByID)
			deliveryRouter.Get("/total/count", o.middleware.Authentication.RequiredDeliveryAuthentication(), o.orderHandler.DeliveryCountingOrder)
			deliveryRouter.Patch("/:id", o.middleware.Authentication.RequiredDeliveryAuthentication(), o.orderHandler.UpdateOrderStatusByDelivery)
		}

		//internal
		internalRouter := orderRouter.Group("/internal")
		{
			internalRouter.Get("/rating/:id", o.middleware.Authentication.RequiredInternalService(), o.orderHandler.InternalGetOrderByOrderID)
		}

		statisticRouter := orderRouter.Group("/statistic")
		{
			//admin
			statisticRouter.Get("/admin/total-order/day",
				o.middleware.Authentication.RequiredRole([]string{auth.ROLE_ADMIN}), o.statisticHandler.AdminGetTotalOrderInSystemInDay)
			statisticRouter.Get("/admin/total-order/month",
				o.middleware.Authentication.RequiredRole([]string{auth.ROLE_ADMIN}), o.statisticHandler.AdminGetTotalOrderInSystemInMonth)
			statisticRouter.Get("/admin/total-order/year",
				o.middleware.Authentication.RequiredRole([]string{auth.ROLE_ADMIN}), o.statisticHandler.AdminGetTotalOrderInSystemInYear)
			statisticRouter.Get("/admin/total-commission",
				o.middleware.Authentication.RequiredRole([]string{auth.ROLE_ADMIN}), o.statisticHandler.AdminGetTotalCommissionOrderInYear)
			statisticRouter.Get("/admin/list-of-product",
				o.middleware.Authentication.RequiredRole([]string{auth.ROLE_ADMIN}), o.statisticHandler.AdminListOfProductSoldOnMonth)

			//store
			statisticRouter.Get("/store/total-order/month",
				o.middleware.Authentication.RequiredStoreAuthentication(), o.statisticHandler.GetTotalOrderInMonthOfStore)
			statisticRouter.Get("/store/total-order/year",
				o.middleware.Authentication.RequiredStoreAuthentication(), o.statisticHandler.GetTotalOrderInYearOfStore)
			statisticRouter.Get("/store/total-commission",
				o.middleware.Authentication.RequiredStoreAuthentication(), o.statisticHandler.GetTotalStoreCommissionInYear)
			statisticRouter.Get("/store/list-of-product",
				o.middleware.Authentication.RequiredStoreAuthentication(), o.statisticHandler.ListOfProductSoldOnMonthStore)
		}
	}

}
