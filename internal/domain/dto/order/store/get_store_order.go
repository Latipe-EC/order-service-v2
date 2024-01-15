package store

import (
	"latipe-order-service-v2/internal/domain/dto/order"
	"latipe-order-service-v2/pkg/util/pagable"
)

type GetStoreOrderRequest struct {
	BaseHeader order.BaseHeader
	StoreID    string
	Keyword    string `query:"keyword" validate:"required"`
	Query      *pagable.Query
}

type FindStoreOrderRequest struct {
	Keyword    string `query:"keyword" validate:"required"`
	BaseHeader order.BaseHeader
	StoreID    string
	Query      *pagable.Query
}

type GetStoreOrderResponse struct {
	pagable.ListResponse
}
