package statistic

import "latipe-order-service-v2/internal/domain/dto/custom_entity"

type GetRevenueDistributionRequest struct {
	Date string `query:"date" validate:"required"`
}

type GetRevenueDistributionResponse struct {
	QueryDate string `json:"query_date"`
	custom_entity.AdminRevenuePer
}
