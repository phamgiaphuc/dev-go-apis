package lib

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"time"
)

func ParseTimeDuration(value string, def int, unit time.Duration) time.Duration {
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	return time.Duration(def) * unit
}

func ParseExpiredTime(unit time.Duration) time.Time {
	return time.Now().Add(unit)
}

func GenerateOTP(maxDigits uint32) string {
	bi, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(maxDigits)))),
	)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%0*d", maxDigits, bi)
}
