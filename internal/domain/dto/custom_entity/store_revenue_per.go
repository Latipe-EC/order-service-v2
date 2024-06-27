package custom_entity

type StoreRevenuePer struct {
	Revenue      int64 `json:"revenue"`
	StoreVoucher int64 `json:"store_voucher"`
	PlatformFee  int64 `json:"platform_fee"`
	Profit       int64 `json:"profit"`
}
