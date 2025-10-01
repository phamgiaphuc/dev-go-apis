package lib

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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

func ParseUnixTime(timestamp string) (time.Time, error) {
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(ts, 0).UTC(), nil
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

func GenerateSHA256(input, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	if _, err := h.Write([]byte(input)); err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

func CompareSHA256(clientSignature, serverSignature string) bool {
	return hmac.Equal([]byte(clientSignature), []byte(serverSignature))
}
