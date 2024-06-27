package statisticQuery

import (
	"context"
	"latipe-order-service-v2/internal/domain/dto/order/statistic"
)

type OrderStatisticUsecase interface {
	//custom_entity - admin
	AdminGetTotalOrderInSystemInDay(ctx context.Context, dto *statistic.AdminTotalOrderInDayRequest) (*statistic.AdminTotalOrderInDayResponse, error)
	AdminGetTotalOrderInSystemInMonth(ctx context.Context, dto *statistic.AdminTotalOrderInMonthRequest) (*statistic.AdminTotalOrderInMonthResponse, error)
	AdminGetTotalOrderInSystemInYear(ctx context.Context, dto *statistic.AdminGetTotalOrderInYearRequest) (*statistic.AdminGetTotalOrderInYearResponse, error)
	AdminGetTotalCommissionOrderInYear(ctx context.Context, dto *statistic.OrderCommissionDetailRequest) (*statistic.OrderCommissionDetailResponse, error)
	AdminListOfProductSoldOnMonth(ctx context.Context, dto *statistic.ListOfProductSoldRequest) (*statistic.ListOfProductSoldResponse, error)
	AdminGetRevenueDistributionInMonth(ctx context.Context, dto *statistic.GetRevenueDistributionRequest) (*statistic.GetRevenueDistributionResponse, error)

	//custom_entity - store
	GetTotalOrderInMonthOfStore(ctx context.Context, dto *statistic.GetTotalStoreOrderInMonthRequest) (*statistic.GetTotalOrderInMonthResponse, error)
	GetTotalOrderInYearOfStore(ctx context.Context, dto *statistic.GetTotalOrderInYearOfStoreRequest) (*statistic.GetTotalOrderInYearOfStoreResponse, error)
	GetTotalStoreCommissionInYear(ctx context.Context, dto *statistic.OrderCommissionDetailRequest) (*statistic.OrderCommissionDetailResponse, error)
	ListOfProductSoldOnMonthStore(ctx context.Context, dto *statistic.ListOfProductSoldRequest) (*statistic.ListOfProductSoldResponse, error)
	GetStoreRevenueDistributionInMonth(ctx context.Context, dto *statistic.GetStoreRevenueDistributionRequest) (*statistic.GetStoreRevenueDistributionResponse, error)

	//export to xls
	AdminExportOrderData(ctx context.Context, dto *statistic.ExportOrderDataForAdminRequest) (*statistic.ExportOrderDataForAdminResponse, error)
	StoreExportOrderData(ctx context.Context, dto *statistic.ExportOrderDataForStoreRequest) (*statistic.ExportOrderDataForStoreResponse, error)
}
