package errors

import (
	"latipe-order-service-v2/internal/domain/enum"
)

var (
	NotAvailableQuantity = &Error{
		Status:    500,
		Code:      enum.INTERNAL,
		ErrorCode: "PROD001",
		Message:   "Product quantity isn't available",
	}
)
