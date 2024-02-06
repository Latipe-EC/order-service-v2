package orderstatistic

import (
	"context"
	"latipe-order-service-v2/internal/domain/dto/order/statistic"
	"latipe-order-service-v2/internal/domain/entities/order"
)

type orderStatisticService struct {
	orderRepo order.OrderRepository
}

func NewOrderStatisicService(orderRepos order.OrderRepository) OrderStatisticUsecase {
	return &orderStatisticService{orderRepo: orderRepos}
}

func (o orderStatisticService) AdminGetTotalOrderInSystemInDay(ctx context.Context, dto *statistic.AdminTotalOrderInDayRequest) (*statistic.AdminTotalOrderInDayResponse, error) {
	items, err := o.orderRepo.GetTotalOrderInSystemInDay(ctx, dto.Date)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.AdminTotalOrderInDayResponse{Items: items}
	dataResp.FilterDate = dto.Date
	return &dataResp, nil
}

func (o orderStatisticService) AdminGetTotalOrderInSystemInMonth(ctx context.Context, dto *statistic.AdminTotalOrderInMonthRequest) (*statistic.AdminTotalOrderInMonthResponse, error) {
	items, err := o.orderRepo.GetTotalOrderInSystemInMonth(ctx, dto.Date)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.AdminTotalOrderInMonthResponse{Items: items}
	dataResp.FilterDate = dto.Date
	return &dataResp, nil
}

func (o orderStatisticService) AdminGetTotalOrderInSystemInYear(ctx context.Context, dto *statistic.AdminGetTotalOrderInYearRequest) (*statistic.AdminGetTotalOrderInYearResponse, error) {
	items, err := o.orderRepo.GetTotalOrderInSystemInYear(ctx, dto.Year)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.AdminGetTotalOrderInYearResponse{Items: items}
	return &dataResp, nil
}

func (o orderStatisticService) AdminGetTotalCommissionOrderInYear(ctx context.Context, dto *statistic.OrderCommissionDetailRequest) (*statistic.OrderCommissionDetailResponse, error) {

	items, err := o.orderRepo.GetTotalCommissionOrderInYear(ctx, dto.Date)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.OrderCommissionDetailResponse{Items: items}
	dataResp.FilterDate = dto.Date
	return &dataResp, nil
}

func (o orderStatisticService) AdminListOfProductSoldOnMonth(ctx context.Context, dto *statistic.ListOfProductSoldRequest) (*statistic.ListOfProductSoldResponse, error) {
	items, err := o.orderRepo.TopOfProductSold(ctx, dto.Date, dto.Count)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.ListOfProductSoldResponse{Items: items}
	dataResp.FilterDate = dto.Date
	return &dataResp, nil
}

func (o orderStatisticService) GetTotalOrderInMonthOfStore(ctx context.Context, dto *statistic.GetTotalStoreOrderInMonthRequest) (*statistic.GetTotalOrderInMonthResponse, error) {
	items, err := o.orderRepo.GetTotalOrderInSystemInMonthOfStore(ctx, dto.Date, dto.StoreId)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.GetTotalOrderInMonthResponse{Items: items}
	dataResp.FilterDate = dto.Date
	return &dataResp, nil
}
func (o orderStatisticService) GetTotalOrderInYearOfStore(ctx context.Context, dto *statistic.GetTotalOrderInYearOfStoreRequest) (*statistic.GetTotalOrderInYearOfStoreResponse, error) {
	items, err := o.orderRepo.GetTotalOrderInSystemInYearOfStore(ctx, dto.Year, dto.StoreID)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.GetTotalOrderInYearOfStoreResponse{Items: items}
	return &dataResp, nil
}
func (o orderStatisticService) GetTotalStoreCommissionInYear(ctx context.Context, dto *statistic.OrderCommissionDetailRequest) (*statistic.OrderCommissionDetailResponse, error) {
	items, err := o.orderRepo.GetTotalCommissionOrderInYearOfStore(ctx, dto.Date, dto.StoreId)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.OrderCommissionDetailResponse{Items: items}
	dataResp.FilterDate = dto.Date
	return &dataResp, nil
}

func (o orderStatisticService) ListOfProductSoldOnMonthStore(ctx context.Context, dto *statistic.ListOfProductSoldRequest) (*statistic.ListOfProductSoldResponse, error) {
	items, err := o.orderRepo.TopOfProductSoldOfStore(ctx, dto.Date, dto.Count, dto.StoreId)
	if err != nil {
		return nil, err
	}

	dataResp := statistic.ListOfProductSoldResponse{Items: items}
	dataResp.FilterDate = dto.Date
	return &dataResp, nil
}
