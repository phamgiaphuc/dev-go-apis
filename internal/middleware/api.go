package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"dev-go-apis/internal/lib"
	"encoding/hex"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func ApiKeyHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
			lib.SendErrorResponse(ctx, lib.ExpiredTimestampError.WithStack(err.Error()))
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
