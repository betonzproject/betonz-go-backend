package timeutils

import (
	"strings"
	"time"
)

// Returns a `time.Time` object representing the start of today in Asia/Yangon timezone.
func StartOfToday() time.Time {
	location, _ := time.LoadLocation("Asia/Yangon")
	now := time.Now().In(location)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
}

// Returns a `time.Time` object representing the end of today in Asia/Yangon timezone.
func EndOfToday() time.Time {
	location, _ := time.LoadLocation("Asia/Yangon")
	now := time.Now().In(location)
	return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, -1, location)
}

// Parses a string in the format "YYYY/MM/DD hh:mm:ss" (24-hour format) or "YYYY/MM/DD h:mm:ss AA" (12-hour format)
// representing a date into a `time.Time` struct.
//
// The timestamp strings are assumed to have Asia/Yangon timezone.
func ParseDate(s string) (time.Time, error) {
	location, _ := time.LoadLocation("Asia/Yangon")
	is12HourFormat := strings.Contains(s, "AM") || strings.Contains(s, "PM")
	var layout string
	if is12HourFormat {
		layout = "2006/01/02 3:04:05 PM"
	} else {
		layout = "2006/01/02 15:04:05"
	}

	t, err := time.ParseInLocation(layout, s, location)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
