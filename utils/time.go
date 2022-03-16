package utils

import (
	"time"

	_const "cve-sa-backend/utils/const"
)

func StrToTime(s string) time.Time {
	location, err := time.ParseInLocation(_const.Format, s, time.Local)
	if err != nil {
		return time.Time{}
	}
	return location
}