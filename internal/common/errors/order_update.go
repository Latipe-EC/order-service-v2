package errors

import "latipe-order-service-v2/internal/domain/enum"

var (
	OrderStatusNotValid = &Error{
		Status:    400,
		Code:      enum.INVALID_ARGUMENT,
		ErrorCode: "UP001",
		Message:   "Đơn hàng có trạng thái chưa phù hợp để bạn chỉnh sửa",
	}

	OrderCannotCancel = &Error{
		Status:    400,
		Code:      enum.INVALID_ARGUMENT,
		ErrorCode: "UP002",
		Message:   "Đơn hàng đang được vận chuyển bạn không thể hủy",
	}

	OrderCannotUpdate = &Error{
		Status:    400,
		Code:      enum.INVALID_ARGUMENT,
		ErrorCode: "UP003",
		Message:   "Đơn hàng ở trạng thái không phù hợp để cập nhật",
	}

	OrderCannotRefund = &Error{
		Status:    400,
		Code:      enum.INVALID_ARGUMENT,
		ErrorCode: "UP004",
		Message:   "Đơn hàng không thể hủy",
	}

	OrderCannotCreated = &Error{
		Status:    400,
		Code:      enum.INVALID_ARGUMENT,
		ErrorCode: "UP005",
		Message:   "Tạo đơn hàng thất bại",
	}
)
