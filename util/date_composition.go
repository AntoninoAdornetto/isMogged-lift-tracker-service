package util

import (
	"time"
)

func FormatMSEpoch(n int64) time.Time {
	ms := time.Unix(0, n*int64(time.Millisecond))
	ms.Format("2006-01-02 15:04:05")
	return ms
}
