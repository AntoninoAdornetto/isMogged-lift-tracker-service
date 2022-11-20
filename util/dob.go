package util

import (
	"strconv"
	"strings"
	"time"
)

type DOB struct {
	month int
	day   int
	year  int
}

func FormatDOB(dateStr string) time.Time {
	dob := &DOB{}
	dates := strings.Split(dateStr, "-")

	for i, val := range dates {
		dateInt, _ := strconv.Atoi(string(val))

		if i == 0 {
			dob.month = dateInt
		}

		if i == 1 {
			dob.day = dateInt
		}

		if i == 2 {
			dob.year = dateInt
		}
	}

	t := time.Date(dob.year, time.Month(dob.month), dob.day, 0, 0, 0, 0, time.UTC)
	return t
}
