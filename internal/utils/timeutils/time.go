package timeutils

import (
	"errors"
	"strings"
	"time"
)

// Returns a `time.Time` object representing the start of today in Asia/Yangon timezone
func StartOfToday() time.Time {
	location, _ := time.LoadLocation("Asia/Yangon")
	now := time.Now().In(location)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
}

// Returns a `time.Time` object representing the end of today in Asia/Yangon timezone
func EndOfToday() time.Time {
	location, _ := time.LoadLocation("Asia/Yangon")
	now := time.Now().In(location)
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, location)
}

// Parses a string in the format "YYYY/MM/DD hh:mm:ss - YYYY/MM/DD hh:mm:ss" representing a start date and
// end date separated by the '-' character and returns time objects for start and end respectively.
//
// The timestamp strings are assumed to have Asia/Yangon timezone
func ParseDateRange(s string) (time.Time, time.Time, error) {
	location, _ := time.LoadLocation("Asia/Yangon")
	is12HourFormat := strings.Contains(s, "AM") || strings.Contains(s, "PM")
	var layout string
	if is12HourFormat {
		layout = "2006/01/02 3:04:05 PM"
	} else {
		layout = "2006/01/02 15:04:05"
	}

	dateParts := strings.Split(s, " - ")
	if len(dateParts) != 2 {
		return time.Time{}, time.Time{}, errors.New("Invalid time format")
	}

	startDate, err := time.ParseInLocation(layout, dateParts[0], location)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endDate, err := time.ParseInLocation(layout, dateParts[1], location)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return startDate, endDate, nil
}
