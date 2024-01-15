package order

import (
	"latipe-order-service-v2/pkg/util/pagable"
)

type GetOrderListRequest struct {
	BaseHeader BaseHeader
	Query      *pagable.Query
}
type GetOrderListResponse struct {
	pagable.ListResponse
}
