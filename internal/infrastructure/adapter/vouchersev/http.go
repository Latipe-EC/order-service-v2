package voucherserv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/wire"
	"latipe-order-service-v2/config"
	"latipe-order-service-v2/internal/common/errors"
	"latipe-order-service-v2/internal/infrastructure/adapter/vouchersev/dto"
	"latipe-order-service-v2/pkg/util/mapper"

	http "latipe-order-service-v2/pkg/internal_http"
)

var Set = wire.NewSet(
	NewUserServHttpAdapter,
)

type httpAdapter struct {
	client http.Client
}

func NewUserServHttpAdapter(config *config.Config) Service {
	restyClient := http.New()
	restyClient.SetRestyClient(
		restyClient.
			Resty().SetBaseURL(config.AdapterService.PromotionService.BaseURL))

	return httpAdapter{
		client: restyClient,
	}
}

func (h httpAdapter) CheckingVoucher(ctx context.Context, req *dto.CheckingVoucherRequest) (*dto.UseVoucherResponse, error) {
	resp, err := h.client.MakeRequest().
		SetContext(ctx).
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", req.BearerToken)).
		SetBody(req).
		Post(req.URL())

	if err != nil {
		log.Errorf("[Checking voucher]: %s", err)
		return nil, errors.ErrBadRequest
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[Checking voucher]: %s", resp.Body())
		return nil, errors.ErrBadRequest
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[Checking voucher]: %s", resp.Body())
		return nil, errors.ErrInternalServer
	}

	var baseResp dto.BaseResponse
	err = json.Unmarshal(resp.Body(), &baseResp)
	if err != nil {
		log.Errorf("[Checking voucher]: %s", err)
		return nil, errors.New("internal server")
	}
	var regResp *dto.UseVoucherResponse

	if err := mapper.BindingStruct(baseResp.Data, &regResp); err != nil {
		log.Errorf("[Checking voucher]: %s", err)
		return nil, err
	}

	return regResp, nil
}
