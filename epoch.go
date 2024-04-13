package main

import (
	"strconv"
	"time"
)

const Year5 = time.Duration(24*time.Hour*365) * 5
const Year15 = time.Duration(24*time.Hour*365) * 15

// SerialToHuman will detect if a number is epoch, or a coded date, ie:
//
// 1712989081 is epoch, because, when converted is less than 15 years ago, and not more than
// 5 years in the future.
//
// If not epoch, we assume a "date" format: 2024041300. Every sequence number 00, 01, is
// assumed to be an hour.
//
// Both are converted to a more human readable string.
func SerialToHuman(s []byte) string {
	// RFC822

	// epoch?
	i, err := strconv.ParseInt(string(s), 10, 64)
	if err == nil {
		// check the time interval
		now := time.Now()
		t := time.Unix(i, 0)
		if now.Sub(t) > Year15 {
			return dateToHuman(s)
		}
		if t.Sub(now) > Year5 {
			return dateToHuman(s)
		}
		return "  " + t.UTC().Format(time.RFC1123)
	}
	return "  " + dateToHuman(s)
}

func dateToHuman(s []byte) string {
	// 2024.04.13.00
	if len(s) != 10 {
		return ""
	}
	year, _ := strconv.ParseInt(string(s[:4]), 10, 64)
	mon, _ := strconv.ParseInt(string(s[4:6]), 10, 64)
	day, _ := strconv.ParseInt(string(s[6:8]), 10, 64)
	hour, _ := strconv.ParseInt(string(s[8:10]), 10, 64)

	println(year, mon, day, hour)

	t := time.Date(int(year), time.Month(mon), int(day), int(hour), 0, 0, 0, time.UTC)
	return t.Format(time.RFC1123)
}
