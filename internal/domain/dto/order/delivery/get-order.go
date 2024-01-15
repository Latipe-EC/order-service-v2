package delivery

import (
	"latipe-order-service-v2/internal/infrastructure/adapter/productserv/dto"
	"latipe-order-service-v2/pkg/util/pagable"
)

type GetOrderListRequest struct {
	BaseHeader dto.BaseHeader
	DeliveryID string
	Query      *pagable.Query
	Keyword    string `query:"keyword"`
}

type GetOrderListResponse struct {
	pagable.ListResponse
}
