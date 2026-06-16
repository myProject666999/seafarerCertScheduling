package service

import (
	"time"
)

func parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Now()
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Now()
	}
	return t
}
