package errors

import "latipe-order-service-v2/internal/domain/enum"

var (
	OrderStatusNotValid = &Error{
		Status:    400,
		Code:      enum.INVALID_ARGUMENT,
		ErrorCode: "ST001",
		Message:   "Đơn hàng có trạng thái chưa phù hợp để bạn chỉnh sửa",
	}

	OrderCannotCancel = &Error{
		Status:    400,
		Code:      enum.INVALID_ARGUMENT,
		ErrorCode: "ST001",
		Message:   "Đơn hàng đang được vận chuyển bạn không thể hủy",
	}
)
