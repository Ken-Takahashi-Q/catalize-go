package utils

import "time"

func TimeToDate(date time.Time) time.Time {
	return date.Truncate(24 * time.Hour)
}
