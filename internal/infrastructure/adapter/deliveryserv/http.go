package deliveryserv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/wire"
	"latipe-order-service-v2/config"
	"latipe-order-service-v2/internal/common/errors"
	"latipe-order-service-v2/internal/infrastructure/adapter/deliveryserv/dto"
	http "latipe-order-service-v2/pkg/internal_http"
)

var Set = wire.NewSet(
	NewDeliServHttpAdapter,
)

type httpAdapter struct {
	client http.Client
}

func NewDeliServHttpAdapter(config *config.Config) Service {
	restyClient := http.New()
	restyClient.SetRestyClient(
		restyClient.
			Resty().SetBaseURL(config.AdapterService.DeliveryService.BaseURL))

	return httpAdapter{
		client: restyClient,
	}
}

func (h httpAdapter) CalculateShippingCost(ctx context.Context, req *dto.GetShippingCostRequest) (*dto.GetShippingCostResponse, error) {
	resp, err := h.client.MakeRequest().
		SetBody(req).
		SetContext(ctx).
		Post(req.URL())

	if err != nil {
		log.Errorf("[Shipping Cost]: %s", err)
		return nil, err
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[Shipping Cost]: %s", resp.Body())
		return nil, err
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[Shipping Cost]: %s", resp.Body())
		return nil, errors.ErrBadRequest
	}

	var regResp *dto.GetShippingCostResponse

	if err := json.Unmarshal(resp.Body(), &regResp); err != nil {
		log.Errorf("[Shipping Cost]: %s", err)
		return nil, errors.ErrInternalServer
	}

	return regResp, nil
}

func (h httpAdapter) GetDeliveryByToken(ctx context.Context, req *dto.GetDeliveryByTokenRequest) (*dto.GetDeliveryByTokenResponse, error) {
	resp, err := h.client.MakeRequest().
		SetBody(req).
		SetContext(ctx).
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", req.BearerToken)).
		Get(req.URL())

	if err != nil {
		log.Errorf("[get delivery]: %s", err)
		return nil, err
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[get delivery]: %s", resp.Body())
		return nil, errors.ErrInternalServer
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[get delivery]: %s", resp.Body())
		return nil, errors.ErrBadRequest
	}

	var regResp *dto.GetDeliveryByTokenResponse

	if err := json.Unmarshal(resp.Body(), &regResp); err != nil {
		log.Errorf("[get delivery]: %s", err)
		return nil, errors.ErrInternalServer
	}

	return regResp, nil
}
