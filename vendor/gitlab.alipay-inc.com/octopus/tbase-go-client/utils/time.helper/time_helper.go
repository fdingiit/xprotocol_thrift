package time_helper

import "time"

func NsToMs(ns int64) int64 {
	return ns / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func MsToNs(ms int) int64 {
	return int64(ms) * int64(time.Millisecond) / int64(time.Nanosecond)
}
