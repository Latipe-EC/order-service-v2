package voucherserv

import (
	"context"
	"github.com/stretchr/testify/mock"
	"latipe-order-service-v2/internal/infrastructure/adapter/vouchersev/dto"
)

type VoucherServiceMock struct {
	mock.Mock
}

func (u *VoucherServiceMock) CheckingVoucher(ctx context.Context, req *dto.CheckingVoucherRequest) (*dto.UseVoucherResponse, error) {
	args := u.Called(ctx, req)
	return args.Get(0).(*dto.UseVoucherResponse), args.Error(1)
}
