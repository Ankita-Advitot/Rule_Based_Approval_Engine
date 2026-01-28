package utils

import "time"

func CalculateLeaveDays(from, to time.Time) int {
	days := int(to.Sub(from).Hours()/24) + 1
	return days
}
