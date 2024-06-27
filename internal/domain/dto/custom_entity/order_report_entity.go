package custom_entity

type StatisticOrderRecordData struct {
	OrderID       string `json:"order_id"`
	Username      string `json:"username"`
	SubTotal      int64  `json:"sub_total"`
	Amount        int64  `json:"amount"`
	Status        int    `json:"status"`
	DeliveryId    string `json:"delivery_id"`
	DeliveryName  string `json:"delivery_name"`
	StoreId       string `json:"store_id"`
	StoreName     string `json:"store_name"`
	CreatedDate   string `json:"created_date"`
	PaymentMethod int    `json:"payment_method"`
	ShippingCost  int64  `json:"shipping_cost"`
	SVoucherValue int64  `json:"svoucher_value"`
	PVoucherValue int64  `json:"pvoucher_value"`
	PlatformFee   int64  `json:"platform_fee"`
	StoreReceived int64  `json:"store_received"`
}
