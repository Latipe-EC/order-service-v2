package deliveryserv

import (
	"context"
	"github.com/stretchr/testify/mock"
	"latipe-order-service-v2/internal/infrastructure/adapter/deliveryserv/dto"
)

type DeliveryServiceMock struct {
	mock.Mock
}

func (d *DeliveryServiceMock) CalculateShippingCost(ctx context.Context, req *dto.GetShippingCostRequest) (*dto.GetShippingCostResponse, error) {
	args := d.Called(ctx, req)
	return args.Get(0).(*dto.GetShippingCostResponse), args.Error(1)
}

func (d *DeliveryServiceMock) GetDeliveryByToken(ctx context.Context, req *dto.GetDeliveryByTokenRequest) (*dto.GetDeliveryByTokenResponse, error) {
	args := d.Called(ctx, req)
	return args.Get(0).(*dto.GetDeliveryByTokenResponse), args.Error(1)
}
