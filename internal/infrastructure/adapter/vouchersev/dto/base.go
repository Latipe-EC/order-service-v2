package dto

type AuthorizationHeader struct {
	BearerToken string `reqHeader:"Authorization" json:"bearer_token"`
}

const (
	FREE_SHIP        = 1
	PAYMENT_DISCOUNT = 2
	STORE_DISCOUNT   = 3

	PENDING  = 0
	ACTIVE   = 1
	INACTIVE = -1

	VOUCHER_APPLY_SUCCESS = 1
	VOUCHER_APPLY_FAILED  = -1

	FIXED_DISCOUNT   = 0
	PERCENT_DISCOUNT = 1

	COD_METHOD = 1
)

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
