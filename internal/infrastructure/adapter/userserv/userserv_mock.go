package userserv

import (
	"context"
	"github.com/stretchr/testify/mock"
	"latipe-order-service-v2/internal/infrastructure/adapter/userserv/dto"
)

type UserServiceMock struct {
	mock.Mock
}

func (u *UserServiceMock) GetAddressDetails(ctx context.Context, req *dto.GetDetailAddressRequest) (*dto.GetDetailAddressResponse, error) {
	args := u.Called(ctx, req)
	return args.Get(0).(*dto.GetDetailAddressResponse), args.Error(1)
}
