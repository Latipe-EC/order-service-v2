package responses

import (
	"latipe-order-service-v2/internal/domain/enum"
)

var (
	DefaultSuccess = General{
		Status:    200,
		Code:      enum.OK,
		ErrorCode: "",
		Message:   "success",
		Data:      nil,
	}

	DefaultError = General{
		Status:    500,
		Code:      enum.INTERNAL,
		ErrorCode: "GENERAL_001",
		Message:   "Internal server error",
		Data:      nil,
	}
)
