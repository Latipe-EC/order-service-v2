package custom_entity

type AdminRevenuePer struct {
	PlatformFee     int64 `json:"platform_fee"`
	PlatformVoucher int   `json:"platform_voucher"`
	TotalShipping   int64 `json:"total_shipping"`
	Profit          int64 `json:"profit"` //Profit = PlatformFee - PlatformVoucher
}
