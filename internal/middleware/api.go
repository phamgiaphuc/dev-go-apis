package middleware

import (
	"bytes"
	"io"
	"strconv"
	"time"

	"dev-go-apis/internal/lib"

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
		if lib.HMAC_SECRET_KEY == "" {
			ctx.Next()
			return
		}

		signature := ctx.GetHeader("X-Signature")
		timestamp := ctx.GetHeader("X-Timestamp")

		if signature == "" || timestamp == "" {
			lib.SendErrorResponse(ctx, lib.MissingSignatureAndTimestampError)
			return
		}

		ts, err := lib.ParseUnixTime(timestamp)
		if err != nil || time.Since(ts) > time.Minute {
			lib.SendErrorResponse(ctx, lib.ExpiredTimestampError)
			return
		}

		body, _ := io.ReadAll(ctx.Request.Body)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		data := ctx.Request.Method + ctx.Request.URL.Path + string(body) + timestamp

		expectedSignature := lib.GenerateSHA256(data, lib.HMAC_SECRET_KEY)

		if !lib.CompareSHA256(signature, expectedSignature) {
			lib.SendErrorResponse(ctx, lib.InvalidSignatureError)
			return
		}

		ctx.Next()
	}
}

func ApiRateLimiterHandler(limiter *redis_rate.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.ClientIP()

		res, err := limiter.Allow(ctx, key, redis_rate.PerMinute(20))
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
