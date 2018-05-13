package bestbuy

import (
	"time"
)

func EpochDayHour() (int64, int64) {
	now := time.Now().Unix()
	day := now / 86400
	hour := now / 3600
	return day, hour
}
