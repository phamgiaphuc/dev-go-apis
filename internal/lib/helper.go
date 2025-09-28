package lib

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

func ParseTimeDuration(value string, def int, unit time.Duration) time.Duration {
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}

	unitMap := map[string]time.Duration{
		"d":   24 * time.Hour,
		"w":   7 * 24 * time.Hour,
		"mth": 30 * 24 * time.Hour,
		"y":   365 * 24 * time.Hour,
	}

	for suffix, dur := range unitMap {
		if strings.HasSuffix(value, suffix) {
			numStr := strings.TrimSuffix(value, suffix)
			if n, err := strconv.Atoi(numStr); err == nil {
				return time.Duration(n) * dur
			}
		}
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
