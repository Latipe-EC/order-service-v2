package statisticQuery

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"latipe-order-service-v2/internal/domain/dto/custom_entity"
	"latipe-order-service-v2/internal/domain/dto/order/statistic"
	"latipe-order-service-v2/internal/domain/entities/order"
	"latipe-order-service-v2/internal/infrastructure/adapter/storeserv"
	storeDTO "latipe-order-service-v2/internal/infrastructure/adapter/storeserv/dto"
	"latipe-order-service-v2/internal/infrastructure/excel"
)

type orderStatisticService struct {
	orderRepo    order.OrderRepository
	storeAdapter storeserv.Service
	excelExport  excel.ExporterExcelData
}

func NewOrderStatisicService(orderRepos order.OrderRepository, storeAdapter storeserv.Service, exporter excel.ExporterExcelData) OrderStatisticUsecase {
	return &orderStatisticService{orderRepo: orderRepos, storeAdapter: storeAdapter, excelExport: exporter}
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

func (o orderStatisticService) GetStoreRevenueDistributionInMonth(ctx context.Context, dto *statistic.GetStoreRevenueDistributionRequest) (*statistic.GetStoreRevenueDistributionResponse, error) {
	data, err := o.orderRepo.GetStoreRevenueDistributionInMonth(ctx, dto.StoreId, dto.Date)
	if err != nil {
		log.Errorf("error get store revenue distribution in month %v", err)
		return nil, err
	}

	dataResp := statistic.GetStoreRevenueDistributionResponse{
		QueryDate: dto.Date,
		StoreRevenuePer: custom_entity.StoreRevenuePer(struct {
			Revenue      int64
			StoreVoucher int64
			PlatformFee  int64
			Profit       int64
		}{Revenue: data.Revenue, StoreVoucher: data.StoreVoucher, PlatformFee: int64(data.PlatformFee), Profit: data.Profit}),
	}
	return &dataResp, nil
}

func (o orderStatisticService) AdminGetRevenueDistributionInMonth(ctx context.Context, dto *statistic.GetRevenueDistributionRequest) (*statistic.GetRevenueDistributionResponse, error) {
	data, err := o.orderRepo.AdminRevenueDistributionInMonth(ctx, dto.Date)
	if err != nil {
		log.Errorf("error get admin revenue distribution in month %v", err)
		return nil, err
	}

	dataResp := statistic.GetRevenueDistributionResponse{
		QueryDate: dto.Date,
		AdminRevenuePer: custom_entity.AdminRevenuePer(struct {
			PlatformFee     int64
			PlatformVoucher int
			TotalShipping   int64
			Profit          int64
		}{PlatformFee: data.PlatformFee, PlatformVoucher: data.PlatformVoucher, TotalShipping: data.TotalShipping, Profit: data.Profit}),
	}
	return &dataResp, nil
}

func (o orderStatisticService) AdminExportOrderData(ctx context.Context, dto *statistic.ExportOrderDataForAdminRequest) (*statistic.ExportOrderDataForAdminResponse, error) {
	data, err := o.orderRepo.GetAllOrderDataRecordByAdmin(ctx, dto.Date)
	if err != nil {
		log.Errorf("error get order data for export %v", err)
		return nil, err
	}

	path, fileAttch, err := o.excelExport.ExportAdminOrderStatisticInMonth(dto.UserID, dto.Date, data)
	if err != nil {
		log.Errorf("error export order data for admin %v", err)
		return nil, err
	}

	return &statistic.ExportOrderDataForAdminResponse{
		QueryDate:      dto.Date,
		FileAttachment: fileAttch,
		FileName:       path,
	}, nil
}

func (o orderStatisticService) StoreExportOrderData(ctx context.Context, dto *statistic.ExportOrderDataForStoreRequest) (*statistic.ExportOrderDataForStoreResponse, error) {
	storeResp, err := o.storeAdapter.GetStoreByStoreId(ctx, &storeDTO.GetStoreByIdRequest{StoreID: dto.StoreID})
	if err != nil {
		log.Errorf("error get store data for export %v", err)
		return nil, err
	}
	data, err := o.orderRepo.GetAllOrderDataRecordByStore(ctx, dto.StoreID, dto.Date)
	if err != nil {
		log.Errorf("error get order data for export %v", err)
		return nil, err
	}

	path, fileAttch, err := o.excelExport.ExportStoreOrderStatisticInMonth(storeResp.Name, dto.Username, dto.Date, data)
	if err != nil {
		log.Errorf("error export order data for admin %v", err)
		return nil, err
	}

	return &statistic.ExportOrderDataForStoreResponse{
		QueryDate:      dto.Date,
		FileAttachment: fileAttch,
		FileName:       path,
	}, nil
}
