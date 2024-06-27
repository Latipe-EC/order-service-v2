package excel

import (
	"github.com/xuri/excelize/v2"
	"unicode/utf8"
)

func AutoFitColumn(f *excelize.File, sheetName string) error {
	// Autofit all columns according to their text content
	cols, err := f.GetCols(sheetName)
	if err != nil {
		return err
	}
	for idx, col := range cols {
		largestWidth := 0
		for _, rowCell := range col {
			cellWidth := utf8.RuneCountInString(rowCell) + 2 // + 2 for margin
			if cellWidth > largestWidth {
				largestWidth = cellWidth
			}
		}
		name, err := excelize.ColumnNumberToName(idx + 1)
		if err != nil {
			return err
		}
		err = f.SetColWidth(sheetName, name, name, float64(largestWidth))
		if err != nil {
			return err
		}

	}
	return nil
}

/*
OrderStatusMapping
const (

	ORDER_SYSTEM_PROCESS = 0
	ORDER_CREATED        = 1
	ORDER_PREPARED       = 2
	ORDER_DELIVERY       = 3

	ORDER_SHIPPING_FINISH    = 4
	ORDER_COMPLETED          = 5
	ORDER_REFUND             = 6
	ORDER_CANCEL_BY_USER     = -2
	ORDER_CANCEL_BY_STORE    = -3
	ORDER_CANCEL_BY_ADMIN    = -4
	ORDER_CANCEL_BY_DELI     = -5
	ORDER_CANCEL_USER_REJECT = -7
	ORDER_FAILED             = -1

)
*/
func OrderStatusMapping(status int) string {

	switch status {
	case 0:
		return "Đang xử lý"
	case 1:
		return "Đã tạo"
	case 2:
		return "Đã chuẩn bị"
	case 3:
		return "Đang giao"
	case 4:
		return "Giao hàng xong"
	case 5:
		return "Hoàn thành"
	case 6:
		return "Hoàn tiền"
	case -2:
		return "Hủy bởi người dùng"
	case -3:
		return "Hủy bởi cửa hàng"
	case -4:
		return "Hủy bởi admin"
	case -5:
		return "Hủy bởi shipper"
	case -7:
		return "Hủy bởi người dùng từ chối"
	case -1:
		return "Thất bại"
	default:
		return "#Lỗi"
	}
}

func PaymentMapping(payment int) string {
	switch payment {
	case 1:
		return "Thanh toán khi nhận hàng"
	case 2:
		return "Thanh toán qua Paypal"
	case 3:
		return "Thanh toán qua ví"
	default:
		return "#Lỗi"
	}
}
