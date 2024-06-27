package excel

import (
	"bytes"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/gofiber/fiber/v2/log"
	"github.com/xuri/excelize/v2"
	"latipe-order-service-v2/internal/domain/dto/custom_entity"
	"time"
)

type ExporterExcelData interface {
	ExportAdminOrderStatisticInMonth(createdBy string, queryTime string, data []custom_entity.StatisticOrderRecordData) (string, *bytes.Reader, error)
	ExportStoreOrderStatisticInMonth(storeName string, createBy string, queryTime string, data []custom_entity.StatisticOrderRecordData) (string, *bytes.Reader, error)
}

type excelExportClient struct {
}

func NewExcelExportClient() ExporterExcelData {
	return &excelExportClient{}
}

func (e excelExportClient) ExportAdminOrderStatisticInMonth(createdBy string, queryTime string, data []custom_entity.StatisticOrderRecordData) (string, *bytes.Reader, error) {
	f := excelize.NewFile()
	currentTime := time.Now()

	var totalPlatformReceived int64
	var totalStoreReceived int64
	var totalShippingCost int64
	var totalPlatformDiscount int64
	var totalStoreDiscount int64

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	sheetName := "Sheet1"
	_, err := f.NewSheet(sheetName)
	if err != nil {
		log.Error(err)
		return "", nil, err
	}
	header := []string{"No", "Mã đơn hàng", "Người tạo",
		"Mã CH", "Mã ĐVVC",
		"Tên ĐVVC", "Trạng thái", "Ngày tạo",
		"Thanh toán", "Tổng SP", "Phí vận chuyển",
		"Giảm giá (HT)", "Giảm giá (CH)", "Cửa hàng nhận", "Hê thống nhận"}
	// Set value of a cell.

	startC, err := excelize.JoinCellName("A", 3)
	if err != nil {
		log.Error(err)
		return "", nil, err
	}

	if err := f.SetSheetRow(sheetName, startC, &header); err != nil {
		log.Error(err)
		return "", nil, err
	}

	for i, v := range data {
		if v.Status >= 0 {
			totalPlatformReceived += v.PlatformFee
			totalStoreReceived += v.StoreReceived
			totalShippingCost += v.ShippingCost
			totalPlatformDiscount += v.PVoucherValue
			totalStoreDiscount += v.SVoucherValue
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "A", i+2+2), i+1); err != nil {
			log.Error(err)
			return "", nil, err
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "B", i+2+2), v.OrderID); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "C", i+2+2), v.Username); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "D", i+2+2), v.StoreId); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "E", i+2+2), v.DeliveryId); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "F", i+2+2), v.DeliveryName); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "G", i+2+2), OrderStatusMapping(v.Status)); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "H", i+2+2), v.CreatedDate); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "I", i+2+2), PaymentMapping(v.PaymentMethod)); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "J", i+2+2), v.SubTotal); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "K", i+2+2), v.ShippingCost); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "L", i+2+2), v.PVoucherValue); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "M", i+2+2), v.SVoucherValue); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "N", i+2+2), v.StoreReceived); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "O", i+2+2), v.PlatformFee); err != nil {
			log.Error(err)
			return "", nil, err
		}

	}
	err = AutoFitColumn(f, sheetName)
	if err != nil {
		return "", nil, err
	}

	//merge cell
	mergeCellRangesHeader := [][]string{{"A1", "O2"}}
	for _, mergeCellRange := range mergeCellRangesHeader {
		if err := f.MergeCell(sheetName, mergeCellRange[0], mergeCellRange[1]); err != nil {
			log.Error(err)
			return "", nil, err
		}
	}
	if err := f.SetCellValue(sheetName, "A1", fmt.Sprintf("Báo cáo kinh doanh hệ thống [%s] ", queryTime)); err != nil {
		log.Error(err)
		return "", nil, err
	}

	//merge cell
	mergeCellRangesCreator := [][]string{{"P1", "S1"}, {"P2", "S2"}}
	for _, mergeCellRange := range mergeCellRangesCreator {
		if err := f.MergeCell(sheetName, mergeCellRange[0], mergeCellRange[1]); err != nil {
			log.Error(err)
			return "", nil, err
		}
	}

	if err := f.SetCellValue(sheetName, "P1", fmt.Sprintf("Người tạo: %s ", createdBy)); err != nil {
		log.Error(err)
		return "", nil, err
	}

	if err := f.SetCellValue(sheetName, "P2", fmt.Sprintf("Thời gian: %s ", currentTime.Format("2006/01/02 15:04:05"))); err != nil {
		log.Error(err)
		return "", nil, err
	}
	//top header
	createdInfoStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#03B6FC"}, Pattern: 1},
		Font:      &excelize.Font{Italic: true, Size: 12},
	})
	if err := f.SetCellStyle(sheetName, "P1", "S2", createdInfoStyle); err != nil {
		log.Error(err)
		return "", nil, err
	}

	//top header
	topStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#03FC73"}, Pattern: 1},
		Font:      &excelize.Font{Bold: true, Size: 16},
	})

	if err := f.SetCellStyle(sheetName, "A1", "O2", topStyle); err != nil {
		log.Error(err)
		return "", nil, err
	}

	//top header
	titleStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Font:      &excelize.Font{Bold: true},
	})

	if err := f.SetCellStyle(sheetName, "A3", "O3", titleStyle); err != nil {
		log.Error(err)
		return "", nil, err
	}

	//revenue
	//merge cell
	mergeCellRangesRevenue := [][]string{{"P3", "S3"}, {"P4", "S4"}, {"P5", "S5"}, {"P6", "S6"}, {"P7", "S7"}}
	for _, mergeCellRange := range mergeCellRangesRevenue {
		if err := f.MergeCell(sheetName, mergeCellRange[0], mergeCellRange[1]); err != nil {
			log.Error(err)
			return "", nil, err
		}
	}

	if err := f.SetCellValue(sheetName, "P3", fmt.Sprintf("Tổng tiền hệ thống nhận: %s ₫", humanize.Comma(totalPlatformReceived))); err != nil {
		log.Error(err)
		return "", nil, err
	}

	if err := f.SetCellValue(sheetName, "P4", fmt.Sprintf("Tổng tiền CH nhận: %s ₫", humanize.Comma(totalStoreReceived))); err != nil {
		log.Error(err)
		return "", nil, err
	}

	if err := f.SetCellValue(sheetName, "P5", fmt.Sprintf("Tổng phí vận chuyển: %s ₫", humanize.Comma(totalShippingCost))); err != nil {
		log.Error(err)
		return "", nil, err
	}

	if err := f.SetCellValue(sheetName, "P6", fmt.Sprintf("Tổng giảm giá HT: %s ₫", humanize.Comma(totalPlatformDiscount))); err != nil {
		log.Error(err)
		return "", nil, err
	}

	if err := f.SetCellValue(sheetName, "P7", fmt.Sprintf("Tổng giảm giá CH: %s ₫", humanize.Comma(totalStoreDiscount))); err != nil {
		log.Error(err)
		return "", nil, err
	}

	revenueStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#D9D900"}, Pattern: 1},
	})

	if err := f.SetCellStyle(sheetName, "P3", "S7", revenueStyle); err != nil {
		log.Error(err)
		return "", nil, err
	}

	//save data to excel
	format := currentTime.Format("20060102@150405")
	fileName := fmt.Sprintf("business-report-by-%s-at-%s.xlsx", createdBy, format)
	buf, err := f.WriteToBuffer()
	readerIo := bytes.NewReader(buf.Bytes())

	return fileName, readerIo, nil
}

func (e excelExportClient) ExportStoreOrderStatisticInMonth(storeName string, createBy string, queryTime string, data []custom_entity.StatisticOrderRecordData) (string, *bytes.Reader, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	currentTime := time.Now()

	var totalPlatformReceived int64
	var totalStoreReceived int64
	var totalShippingCost int64
	var totalPlatformDiscount int64
	var totalStoreDiscount int64

	// Create a new sheet.
	sheetName := "STORE"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		log.Error(err)
		return "", nil, err
	}
	f.SetActiveSheet(index)

	header := []string{"STT", "Mã đơn hàng", "Người tạo",
		"Mã CH", "Mã ĐVVC",
		"Tên ĐVVC", "Trạng thái", "Ngày tạo",
		"Thanh toán", "Tổng SP", "Phí vận chuyển",
		"Giảm giá (HT)", "Giảm giá (CH)", "Cửa hàng nhận", "Hệ thống nhận"}
	// Set value of a cell.

	startC, err := excelize.JoinCellName("A", 3)
	if err != nil {
		log.Error(err)
		return "", nil, err
	}

	if err := f.SetSheetRow(sheetName, startC, &header); err != nil {
		log.Error(err)
		return "", nil, err
	}

	for i, v := range data {
		if v.Status >= 0 {
			totalPlatformReceived += v.PlatformFee
			totalStoreReceived += v.StoreReceived
			totalShippingCost += v.ShippingCost
			totalPlatformDiscount += v.PVoucherValue
			totalStoreDiscount += v.SVoucherValue
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "A", i+2+2), i+1); err != nil {
			log.Error(err)
			return "", nil, err
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "B", i+2+2), v.OrderID); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "C", i+2+2), v.Username); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "D", i+2+2), v.StoreId); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "E", i+2+2), v.DeliveryId); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "F", i+2+2), v.DeliveryName); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "G", i+2+2), OrderStatusMapping(v.Status)); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "H", i+2+2), v.CreatedDate); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "I", i+2+2), PaymentMapping(v.PaymentMethod)); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "J", i+2+2), v.SubTotal); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "K", i+2+2), v.ShippingCost); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "L", i+2+2), v.PVoucherValue); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "M", i+2+2), v.SVoucherValue); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "N", i+2+2), v.StoreReceived); err != nil {
			log.Error(err)
			return "", nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", "O", i+2+2), v.PlatformFee); err != nil {
			log.Error(err)
			return "", nil, err
		}

	}
	err = AutoFitColumn(f, sheetName)
	if err != nil {
		return "", nil, err
	}

	//merge cell
	mergeCellRangesHeader := [][]string{{"A1", "O2"}}
	for _, mergeCellRange := range mergeCellRangesHeader {
		if err := f.MergeCell(sheetName, mergeCellRange[0], mergeCellRange[1]); err != nil {
			log.Error(err)
			return "", nil, err
		}
	}
	if err := f.SetCellValue(sheetName, "A1", fmt.Sprintf("Báo cáo kinh doanh [%s] - Cửa Hàng: %s", queryTime, storeName)); err != nil {
		log.Error(err)
		return "", nil, err
	}

	//merge cell
	mergeCellRangesCreator := [][]string{{"P1", "S1"}, {"P2", "S2"}}
	for _, mergeCellRange := range mergeCellRangesCreator {
		if err := f.MergeCell(sheetName, mergeCellRange[0], mergeCellRange[1]); err != nil {
			log.Error(err)
			return "", nil, err
		}
	}

	if err := f.SetCellValue(sheetName, "P1", fmt.Sprintf("Người tạo: %s ", createBy)); err != nil {
		log.Error(err)
		return "", nil, err
	}

	if err := f.SetCellValue(sheetName, "P2", fmt.Sprintf("Thời gian: %s ", currentTime.Format("2006/01/02 15:04:05"))); err != nil {
		log.Error(err)
		return "", nil, err
	}
	//top header
	createdInfoStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#03B6FC"}, Pattern: 1},
		Font:      &excelize.Font{Italic: true, Size: 12},
	})
	if err := f.SetCellStyle(sheetName, "P1", "S2", createdInfoStyle); err != nil {
		log.Error(err)
		return "", nil, err
	}

	//top header
	topStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#03FC73"}, Pattern: 1},
		Font:      &excelize.Font{Bold: true, Size: 16},
	})

	if err := f.SetCellStyle(sheetName, "A1", "O2", topStyle); err != nil {
		log.Error(err)
		return "", nil, err
	}

	//top header
	titleStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Font:      &excelize.Font{Bold: true},
	})

	if err := f.SetCellStyle(sheetName, "A3", "O3", titleStyle); err != nil {
		log.Error(err)
		return "", nil, err
	}

	//revenue
	//merge cell
	mergeCellRangesRevenue := [][]string{{"P3", "S3"}, {"P4", "S4"}, {"P5", "S5"}, {"P6", "S6"}, {"P7", "S7"}}
	for _, mergeCellRange := range mergeCellRangesRevenue {
		if err := f.MergeCell(sheetName, mergeCellRange[0], mergeCellRange[1]); err != nil {
			log.Error(err)
			return "", nil, err
		}
	}

	if err := f.SetCellValue(sheetName, "P3", fmt.Sprintf("Tổng số tiền hệ thống nhận: %s ₫", humanize.Comma(totalPlatformReceived))); err != nil {
		log.Error(err)
		return "", nil, err
	}

	if err := f.SetCellValue(sheetName, "P4", fmt.Sprintf("Tổng số tiền CH nhận: %s ₫", humanize.Comma(totalStoreReceived))); err != nil {
		log.Error(err)
		return "", nil, err
	}

	if err := f.SetCellValue(sheetName, "P5", fmt.Sprintf("Tổng phí vận chuyển: %s ₫", humanize.Comma(totalShippingCost))); err != nil {
		log.Error(err)
		return "", nil, err
	}

	if err := f.SetCellValue(sheetName, "P6", fmt.Sprintf("Tổng giảm giá hệ thống: %s ₫", humanize.Comma(totalPlatformDiscount))); err != nil {
		log.Error(err)
		return "", nil, err
	}

	if err := f.SetCellValue(sheetName, "P7", fmt.Sprintf("Tổng giảm giá CH: %s ₫", humanize.Comma(totalStoreDiscount))); err != nil {
		log.Error(err)
		return "", nil, err
	}

	revenueStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#D9D900"}, Pattern: 1},
	})

	if err := f.SetCellStyle(sheetName, "P3", "S7", revenueStyle); err != nil {
		log.Error(err)
		return "", nil, err
	}

	//save data to excel

	format := currentTime.Format("20060102@150405")
	fileName := fmt.Sprintf("store-%s-business-report-by-%s-at-%s.xlsx", storeName, createBy, format)

	buf, err := f.WriteToBuffer()
	readerIo := bytes.NewReader(buf.Bytes())

	return fileName, readerIo, nil
}
