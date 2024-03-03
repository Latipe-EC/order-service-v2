package errors

import (
	"latipe-order-service-v2/internal/domain/enum"
)

var (
	ErrInternalServer = &Error{
		Status:    500,
		Code:      enum.INTERNAL,
		ErrorCode: "GE001",
		Message:   "Internal server error",
	}

	ErrBadRequest = &Error{
		Status:    400,
		Code:      enum.INVALID_ARGUMENT,
		ErrorCode: "GE002",
		Message:   "Bad request",
	}

	ErrNotChange = &Error{
		Status:    200,
		Code:      enum.OK,
		ErrorCode: "GE003",
		Message:   "Not change",
	}

	ErrPermissionDenied = &Error{
		Status:    403,
		Code:      enum.PERMISSION_DENIED,
		ErrorCode: "GE004",
		Message:   "Permission denied",
	}

	ErrNotFound = &Error{
		Status:    404,
		Code:      enum.NOT_FOUND,
		ErrorCode: "GE005",
		Message:   "Not found",
	}

	ErrAlreadyExists = &Error{
		Status:    409,
		Code:      enum.ALREADY_EXISTS,
		ErrorCode: "GE006",
		Message:   "Already exists",
	}

	ErrUnauthenticated = &Error{
		Status:    401,
		Code:      enum.UNAUTHENTICATED,
		ErrorCode: "GE007",
		Message:   "Unauthorized",
	}

	ErrInvalidCredentials = &Error{
		Status:    401,
		Code:      enum.UNAUTHENTICATED,
		ErrorCode: "GE008",
		Message:   "Invalid credentials",
	}

	ErrNotFoundRecord = &Error{
		Status:    404,
		Code:      enum.NOT_FOUND,
		ErrorCode: "GE009",
		Message:   "Record does not exist",
	}

	ErrInvalidParameters = &Error{
		Status:    400,
		Code:      enum.FAILED_PRECONDITION,
		ErrorCode: "GE010",
		Message:   "Invalid parameters",
	}

	ErrTooManyRequest = &Error{
		Status:    429,
		Code:      enum.RESOURCE_EXHAUSTED,
		ErrorCode: "GE011",
		Message:   "Too Many Requests",
	}

	ErrInvalidFilter = &Error{
		Status:    400,
		Code:      enum.INVALID_ARGUMENT,
		ErrorCode: "GE012",
		Message:   "Invalid filters",
	}
)
