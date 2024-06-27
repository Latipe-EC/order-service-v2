package statistic

import "latipe-order-service-v2/internal/domain/dto/custom_entity"

type GetStoreRevenueDistributionRequest struct {
	StoreId string
	Date    string `query:"date" validate:"required"` //yyyy-mm
}

type GetStoreRevenueDistributionResponse struct {
	QueryDate string `json:"query_date"`
	custom_entity.StoreRevenuePer
}
