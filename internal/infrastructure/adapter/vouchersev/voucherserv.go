package voucherserv

import (
	"context"
	"latipe-order-service-v2/internal/infrastructure/adapter/vouchersev/dto"
)

type Service interface {
	CheckingVoucher(ctx context.Context, req *dto.CheckingVoucherRequest) (*dto.UseVoucherResponse, error)
}
