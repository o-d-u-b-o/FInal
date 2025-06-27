package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat = "20060102"
	maxDays    = 400
)

var (
	ErrEmptyRepeat     = errors.New("empty repeat rule")
	ErrInvalidFormat   = errors.New("invalid repeat format")
	ErrInvalidDay      = errors.New("invalid day")
	ErrInvalidMonth    = errors.New("invalid month")
	ErrInvalidWeekday  = errors.New("invalid weekday")
	ErrMaxDaysExceeded = errors.New("max days exceeded")
	ErrUnsupportedRule = errors.New("unsupported repeat rule")
)

func NextDate(now time.Time, dateStr, repeat string) (string, error) {
	if repeat == "" {
		return "", ErrEmptyRepeat
	}

	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return "", fmt.Errorf("invalid date: %w", err)
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", ErrInvalidFormat
	}

	switch parts[0] {
	case "d":
		return handleDailyRule(now, date, parts)
	case "y":
		return handleYearlyRule(now, date), nil
	case "w":
		return handleWeeklyRule(now, date, parts)
	case "m":
		return handleMonthlyRule(now, date, parts)
	default:
		return "", ErrUnsupportedRule
	}
}

func handleDailyRule(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", ErrInvalidFormat
	}

	days, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", ErrInvalidFormat
	}

	if days <= 0 || days > maxDays {
		return "", ErrMaxDaysExceeded
	}

	for {
		date = date.AddDate(0, 0, days)
		if dateAfter(date, now) {
			break
		}
	}

	return date.Format(dateFormat), nil
}

func handleYearlyRule(now, date time.Time) string {
	for {
		date = date.AddDate(1, 0, 0)
		if dateAfter(date, now) {
			break
		}
	}
	return date.Format(dateFormat)
}

func handleWeeklyRule(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", ErrInvalidFormat
	}

	weekdays := make(map[int]bool)
	for _, dayStr := range strings.Split(parts[1], ",") {
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < 1 || day > 7 {
			return "", ErrInvalidWeekday
		}
		weekdays[day] = true
	}

	for {
		date = date.AddDate(0, 0, 1)
		if dateAfter(date, now) {
			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			if weekdays[weekday] {
				break
			}
		}
	}

	return date.Format(dateFormat), nil
}

func handleMonthlyRule(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", ErrInvalidFormat
	}

	daysInput := parts[1]
	if daysInput == "" {
		return "", ErrInvalidDay
	}

	// Парсим дни
	days := strings.Split(daysInput, ",")
	dayNumbers := make([]int, 0, len(days))
	hasMinus1 := false
	hasMinus2 := false

	for _, dayStr := range days {
		if dayStr == "-1" {
			hasMinus1 = true
			continue
		}
		if dayStr == "-2" {
			hasMinus2 = true
			continue
		}
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < 1 || day > 31 {
			return "", ErrInvalidDay
		}
		dayNumbers = append(dayNumbers, day)
	}

	// Парсим месяцы, если указаны
	var months []int
	if len(parts) > 2 {
		monthStrs := strings.Split(parts[2], ",")
		for _, monthStr := range monthStrs {
			month, err := strconv.Atoi(monthStr)
			if err != nil || month < 1 || month > 12 {
				return "", ErrInvalidMonth
			}
			months = append(months, month)
		}
	}

	nextDate := date
	for {
		nextDate = nextDate.AddDate(0, 0, 1)
		if dateAfter(nextDate, now) {
			currentMonth := int(nextDate.Month())
			currentDay := nextDate.Day()

			// Проверяем месяцы, если они указаны
			if len(months) > 0 {
				monthMatch := false
				for _, m := range months {
					if m == currentMonth {
						monthMatch = true
						break
					}
				}
				if !monthMatch {
					continue
				}
			}

			// Проверяем специальные дни
			if hasMinus1 && isLastDayOfMonth(nextDate) {
				return nextDate.Format(dateFormat), nil
			}
			if hasMinus2 && isPenultimateDayOfMonth(nextDate) {
				return nextDate.Format(dateFormat), nil
			}

			// Проверяем обычные дни
			for _, d := range dayNumbers {
				if d == currentDay {
					return nextDate.Format(dateFormat), nil
				}
			}
		}
	}
}

func isLastDayOfMonth(date time.Time) bool {
	return date.AddDate(0, 0, 1).Month() != date.Month()
}

func isPenultimateDayOfMonth(date time.Time) bool {
	return date.AddDate(0, 0, 2).Month() != date.Month()
}

func dateAfter(t1, t2 time.Time) bool {
	return t1.Format(dateFormat) > t2.Format(dateFormat)
}
