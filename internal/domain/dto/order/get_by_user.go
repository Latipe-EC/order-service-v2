package order

import "latipe-order-service-v2/pkg/util/pagable"

type GetByUserIdRequest struct {
	BaseHeader BaseHeader
	UserId     string
	Query      *pagable.Query
}

type GetByUserIdResponse struct {
	pagable.ListResponse
}
