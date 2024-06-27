package statisticRouter

import (
	"github.com/gofiber/fiber/v2"
	"latipe-order-service-v2/internal/api/order"
	"latipe-order-service-v2/internal/middleware"
	"latipe-order-service-v2/internal/middleware/auth"
)

type OrderStatisticRouter interface {
	Init(root *fiber.Router)
}

type orderStatisticRouter struct {
	statisticHandler order.OrderStatisticApiHandler
	middleware       *middleware.Middleware
}

func NewStatisticOrderRouter(statisticHandler order.OrderStatisticApiHandler,
	middleware *middleware.Middleware) OrderStatisticRouter {
	return orderStatisticRouter{
		statisticHandler: statisticHandler,
		middleware:       middleware,
	}
}

func (o orderStatisticRouter) Init(root *fiber.Router) {

	statisticRouter := (*root).Group("/statistic")
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
		statisticRouter.Get("/admin/revenue-distribution",
			o.middleware.Authentication.RequiredRole([]string{auth.ROLE_ADMIN}), o.statisticHandler.AdminGetRevenueDistributionInMonth)
		statisticRouter.Get("/admin/business-report",
			o.middleware.Authentication.RequiredRole([]string{auth.ROLE_ADMIN}), o.statisticHandler.AdminExportOrderData)

		//store
		statisticRouter.Get("/store/total-order/month",
			o.middleware.Authentication.RequiredStoreAuthentication(), o.statisticHandler.GetTotalOrderInMonthOfStore)
		statisticRouter.Get("/store/total-order/year",
			o.middleware.Authentication.RequiredStoreAuthentication(), o.statisticHandler.GetTotalOrderInYearOfStore)
		statisticRouter.Get("/store/total-commission",
			o.middleware.Authentication.RequiredStoreAuthentication(), o.statisticHandler.GetTotalStoreCommissionInYear)
		statisticRouter.Get("/store/list-of-product",
			o.middleware.Authentication.RequiredStoreAuthentication(), o.statisticHandler.ListOfProductSoldOnMonthStore)
		statisticRouter.Get("/store/revenue-distribution",
			o.middleware.Authentication.RequiredStoreAuthentication(), o.statisticHandler.GetStoreRevenueDistributionInMonth)
		statisticRouter.Get("/store/business-report",
			o.middleware.Authentication.RequiredStoreAuthentication(), o.statisticHandler.StoreExportOrderData)
	}
}
