package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"dev-go-apis/internal/lib"
	"encoding/hex"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
)

func ApiKeyHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if lib.API_KEY == "" {
			ctx.Next()
			return
		}
		apiKey := ctx.GetHeader("X-ApiKey")
		if apiKey == "" || apiKey != lib.API_KEY {
			lib.SendErrorResponse(ctx, lib.MissingAPIKeyError)
			return
		}
		ctx.Next()
	}
}

func ApiHmacHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		signature := ctx.GetHeader("X-Signature")
		timestamp := ctx.GetHeader("X-Timestamp")

		if signature == "" || timestamp == "" {
			lib.SendErrorResponse(ctx, lib.MissingSignatureAndTimestampError)
			return
		}

		ts, err := time.Parse(time.RFC3339, timestamp)
		if err != nil || time.Since(ts) > 5*time.Minute {
			lib.SendErrorResponse(ctx, lib.ExpiredTimestampError)
			return
		}

		body, _ := io.ReadAll(ctx.Request.Body)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		data := ctx.Request.Method + ctx.Request.URL.Path + string(body) + timestamp
		log.Printf("%v\n", data)

		mac := hmac.New(sha256.New, []byte(lib.HMAC_SECRET_KEY))
		mac.Write([]byte(data))
		expectedMAC := hex.EncodeToString(mac.Sum(nil))

		if !hmac.Equal([]byte(signature), []byte(expectedMAC)) {
			lib.SendErrorResponse(ctx, lib.InvalidSignatureError)
			return
		}

		ctx.Next()
	}
}

func ApiRateLimiterHandler(limiter *redis_rate.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.ClientIP()

		res, err := limiter.Allow(ctx, key, redis_rate.PerMinute(10))
		if err != nil {
			lib.SendErrorResponse(ctx, lib.InternalServerError)
			return
		}

		header := ctx.Writer.Header()
		header.Set("RateLimit-Remaining", strconv.Itoa(res.Remaining))

		if res.Allowed == 0 {
			seconds := int(res.RetryAfter / time.Second)
			header.Set("RateLimit-RetryAfter", strconv.Itoa(seconds))
			lib.SendErrorResponse(ctx, lib.TooManyRequestsError)
			return
		}

		ctx.Next()
	}
}
