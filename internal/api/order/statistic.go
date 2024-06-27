package order

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"latipe-order-service-v2/internal/common/errors"
	"latipe-order-service-v2/internal/common/responses"
	"latipe-order-service-v2/internal/domain/dto/order/statistic"
	"latipe-order-service-v2/internal/middleware/auth"
	"latipe-order-service-v2/internal/services/queries/statisticQuery"
	"latipe-order-service-v2/pkg/util/valid"
	"strings"
)

type statisticApiHandler struct {
	orderStatisticUsecase statisticQuery.OrderStatisticUsecase
}

func NewStatisticHandler(orderStatisticUsecase statisticQuery.OrderStatisticUsecase) OrderStatisticApiHandler {
	return statisticApiHandler{
		orderStatisticUsecase: orderStatisticUsecase,
	}
}

// @Summary Get total order in system in day
// @Tags Statistic
// @Description Get total order in system in day
// @Accept json
// @Produce json
// @Param date query string false "date"
// @Router /admin/statistic/total-order-in-day [get]
func (s statisticApiHandler) AdminGetTotalOrderInSystemInDay(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.AdminTotalOrderInDayRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Date == "" {
		req.Date = InitDateValue()
	}

	result, err := s.orderStatisticUsecase.AdminGetTotalOrderInSystemInDay(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

// @Summary Get total order in system in month
// @Tags Statistic
// @Description Get total order in system in month
// @Accept json
// @Produce json
// @Param date query string false "date"
// @Router /admin/statistic/total-order-in-month [get]
func (s statisticApiHandler) AdminGetTotalOrderInSystemInMonth(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.AdminTotalOrderInMonthRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Date == "" {
		req.Date = InitDateValue()
	}

	result, err := s.orderStatisticUsecase.AdminGetTotalOrderInSystemInMonth(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

// @Summary Get total order in system in year
// @Tags Statistic
// @Description Get total order in system in year
// @Accept json
// @Produce json
// @Param date query string false "date"
// @Router /admin/statistic/total-order-in-year [get]
func (s statisticApiHandler) AdminGetTotalOrderInSystemInYear(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.AdminGetTotalOrderInYearRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	result, err := s.orderStatisticUsecase.AdminGetTotalOrderInSystemInYear(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

// @Summary Get total commission order in day
// @Tags Statistic
// @Description Get total commission order in day
// @Accept json
// @Produce json
// @Param date query string false "date"
// @Router /admin/statistic/total-commission-order-in-day [get]
func (s statisticApiHandler) AdminGetTotalCommissionOrderInYear(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.OrderCommissionDetailRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Date == "" {
		req.Date = InitDateValue()
	}

	result, err := s.orderStatisticUsecase.AdminGetTotalCommissionOrderInYear(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

// @Summary Get total commission order in month
// @Tags Statistic
// @Description Get total commission order in month
// @Accept json
// @Produce json
// @Param date query string false "date"
// @Router /admin/statistic/total-commission-order-in-month [get]
func (s statisticApiHandler) AdminListOfProductSoldOnMonth(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.ListOfProductSoldRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Date == "" {
		req.Date = InitDateValue()
	}

	result, err := s.orderStatisticUsecase.AdminListOfProductSoldOnMonth(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

// @Summary Get total order in day of store
// @Tags Statistic
// @Description Get total order in day of store
// @Accept json
// @Produce json
// @Param date query string false "date"
// @Router /store/statistic/total-order-in-day [get]
func (s statisticApiHandler) GetTotalOrderInMonthOfStore(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.GetTotalStoreOrderInMonthRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Date == "" {
		req.Date = InitDateValue()
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreId = storeID

	result, err := s.orderStatisticUsecase.GetTotalOrderInMonthOfStore(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

// @Summary Get total order in month of store
// @Tags Statistic
// @Description Get total order in month of store
// @Accept json
// @Produce json
// @Param date query string false "date"
// @Router /store/statistic/total-order-in-month [get]
func (s statisticApiHandler) GetTotalOrderInYearOfStore(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.GetTotalOrderInYearOfStoreRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreID = storeID

	result, err := s.orderStatisticUsecase.GetTotalOrderInYearOfStore(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

// @Summary Get total order in year of store
// @Tags Statistic
// @Description Get total order in year of store
// @Accept json
// @Produce json
// @Param date query string false "date"
// @Router /store/statistic/total-order-in-year [get]
func (s statisticApiHandler) GetTotalStoreCommissionInYear(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.OrderCommissionDetailRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Date == "" {
		req.Date = InitDateValue()
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreId = storeID

	result, err := s.orderStatisticUsecase.GetTotalStoreCommissionInYear(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

// @Summary Get total commission order in month of store
// @Tags Statistic
// @Description Get total commission order in month of store
// @Accept json
// @Produce json
// @Param date query string false "date"
// @Router /store/statistic/total-commission-order-in-month [get]
func (s statisticApiHandler) ListOfProductSoldOnMonthStore(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.ListOfProductSoldRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Date == "" {
		req.Date = InitDateValue()
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreId = storeID

	result, err := s.orderStatisticUsecase.ListOfProductSoldOnMonthStore(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

// @Summary Get store revenue distribution in month
// @Tags Statistic
// @Description Get store revenue distribution in month
// @Accept json
// @Produce json
// @Param date query string false "date"
// @Router /store/statistic/revenue-distribution-in-month [get]
func (s statisticApiHandler) GetStoreRevenueDistributionInMonth(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.GetStoreRevenueDistributionRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	//validate request body date format match yyyy-mm
	if !isValidDateFormat(req.Date) {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreId = storeID

	result, err := s.orderStatisticUsecase.GetStoreRevenueDistributionInMonth(context, &req)
	if err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

// @Summary Get revenue distribution in month
// @Tags Statistic
// @Description Get revenue distribution in month
// @Accept json
// @Produce json
// @Param date query string false "date"
// @Router /admin/statistic/revenue-distribution-in-month [get]
func (s statisticApiHandler) AdminGetRevenueDistributionInMonth(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.GetRevenueDistributionRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	//validate request body date format match yyyy-mm
	if !isValidDateFormat(req.Date) {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	result, err := s.orderStatisticUsecase.AdminGetRevenueDistributionInMonth(context, &req)
	if err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

// @Summary Export order data
// @Tags Statistic
// @Description Export order data
// @Accept json
// @Produce json
// @Param date query string false "date"
// @Router /admin/statistic/export-order-data [get]
func (s statisticApiHandler) AdminExportOrderData(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.ExportOrderDataForAdminRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}
	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	//validate request body date format match yyyy-mm
	if !isValidDateFormat(req.Date) {
		return errors.ErrInvalidParameters
	}

	userName := fmt.Sprintf("%v", ctx.Locals(auth.USERNAME))
	req.UserID = userName

	result, err := s.orderStatisticUsecase.AdminExportOrderData(context, &req)
	if err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	ctx.Attachment(result.FileName)
	return ctx.SendStream(result.FileAttachment)
}

// @Summary Export order data for store
// @Tags Statistic
// @Description Export order data for store
// @Accept json
// @Produce json
// @Param date query string false "date"
// @Router /store/statistic/export-order-data [get]
func (s statisticApiHandler) StoreExportOrderData(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.ExportOrderDataForStoreRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}
	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	//validate request body date format match yyyy-mm
	if !isValidDateFormat(req.Date) {
		return errors.ErrInvalidParameters
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreID = storeID

	username := fmt.Sprintf("%v", ctx.Locals(auth.USERNAME))
	req.Username = username

	result, err := s.orderStatisticUsecase.StoreExportOrderData(context, &req)
	if err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	ctx.Attachment(result.FileName)
	return ctx.SendStream(result.FileAttachment)
}
