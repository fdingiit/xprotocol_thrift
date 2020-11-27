package prometheus

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

func GetTimestampSecond() int64 {
	return int64(float64(time.Now().Unix())/60) * 60
}

func GetTimestampMS() int64 {
	return GetTimestampSecond() * 1000
}

func StringToFloat(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return -1, err
	}
	if math.IsInf(v, 1) {
		return math.MaxFloat64, nil
	}
	if math.IsInf(v, -1) {
		return -math.MaxFloat64, nil
	}
	if math.IsNaN(v) {
		return -1, fmt.Errorf("an IEEE 754 not-a-number value")
	}
	return v, nil
}
