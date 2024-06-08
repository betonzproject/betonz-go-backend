package timeutils

import (
	"errors"
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

// Returns a `time.Time` object representing the start of this month in Asia/Yangon timezone.
func StartOfThisMonth() time.Time {
	location, _ := time.LoadLocation("Asia/Yangon")
	now := time.Now().In(location)
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, location)
}

// Returns a `time.Time` object representing the end of this month in Asia/Yangon timezone.
func EndOfThisMonth() time.Time {
	location, _ := time.LoadLocation("Asia/Yangon")
	now := time.Now().In(location)
	return time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, -1, location)
}

// Returns a `time.Time` object representing the start of this year in Asia/Yangon timezone.
func StartOfThisYear() time.Time {
	location, _ := time.LoadLocation("Asia/Yangon")
	now := time.Now().In(location)
	return time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, location)
}

// Returns a int representing the total days of current month.
func DaysInMonth() int {
	currentYear, currentMonth, _ := time.Now().Date()

	firstDayOfNextMonth := time.Date(currentYear, currentMonth+1, 1, 0, 0, 0, 0, time.UTC)

	lastDayOfCurrentMonth := firstDayOfNextMonth.Add(-time.Hour * 24)

	totalDays := lastDayOfCurrentMonth.Day()

	return totalDays
}

// Returns a `time.Time` object representing the end of this year in Asia/Yangon timezone.
func EndOfThisYear() time.Time {
	location, _ := time.LoadLocation("Asia/Yangon")
	now := time.Now().In(location)
	return time.Date(now.Year()+1, time.January, 1, 0, 0, 0, -1, location)
}

// Parses a string in the format "YYYY/MM/DD" representing a date into a `time.Time` struct.
func ParseDate(s string) (time.Time, error) {
	location, _ := time.LoadLocation("Asia/Yangon")
	t, err := time.ParseInLocation("2006/01/02", s, location)
	return t, err
}

// CheckIfAtLeast18 function checks if the user is 18 or older based on their birth date
func CheckIfAtLeast18(birthDate time.Time) error {
	now := time.Now()

	age := now.Year() - birthDate.Year()

	if now.YearDay() < birthDate.YearDay() {
		age--
	}

	if age < 18 {
		return errors.New("user is under 18")
	}
	return nil
}

// Parses a string in the format "YYYY/MM/DD hh:mm:ss" (24-hour format) or "YYYY/MM/DD h:mm:ss AA" (12-hour format)
// representing a datetime into a `time.Time` struct.
//
// The timestamp strings are assumed to have Asia/Yangon timezone.
func ParseDatetime(s string) (time.Time, error) {
	location, _ := time.LoadLocation("Asia/Yangon")
	is12HourFormat := strings.Contains(s, "AM") || strings.Contains(s, "PM")
	var layout string
	if is12HourFormat {
		layout = "2006/01/02 3:04:05 PM"
	} else {
		layout = "2006/01/02 15:04:05"
	}

	t, err := time.ParseInLocation(layout, s, location)
	return t, err
}

// ParseDateTimeInLocation parses a string into the format "YYYY/MM/DD hh:mm:ss" (24-hour format) or "YYYY/MM/DD h:mm:ss AA" (12-hour format),
// considering the given timezone represented by loc.
//
// Parameters:
//   - s: The string to be parsed into a timestamp.
//   - loc: The time.Location representing the timezone for the timestamp string.
//
// Returns:
//   - A time.Time object representing the parsed timestamp.
//   - An error if the string cannot be parsed into a valid timestamp.
func ParseDateTimeInLocation(s string, loc time.Location) (time.Time, error) {
	is12HourFormat := strings.Contains(s, "AM") || strings.Contains(s, "PM")
	var layout string
	if is12HourFormat {
		layout = "2006/01/02 3:04:05 PM"
	} else {
		layout = "2006/01/02 15:04:05"
	}

	t, err := time.ParseInLocation(layout, s, &loc)
	return t, err
}
