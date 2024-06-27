package order

import (
	"fmt"
	"time"
)

func InitDateValue() string {
	now := time.Now()
	return fmt.Sprintf("%v-%v-%v", now.Year(), int(now.Month()), now.Day())
}

// Helper function to validate the date format (yyyy-mm)
func isValidDateFormat(date string) bool {
	_, err := time.Parse("2006-01", date)
	return err == nil
}
